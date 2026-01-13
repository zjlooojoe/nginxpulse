package enrich

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/likaia/nginxpulse/internal/config"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
)

//go:embed data/ip2region.xdb
var ipDataFiles embed.FS

var (
	ipSearcher  *xdb.Searcher
	vectorIndex []byte
	dbPath      = filepath.Join(config.DataDir, "ip2region.xdb")
)

const (
	ipAPIBatchURL  = "http://ip-api.com/batch"
	ipAPIFields    = "status,message,country,countryCode,region,regionName,city,isp,query"
	ipAPILanguage  = "zh-CN"
	ipAPITimeout   = 1200 * time.Millisecond
	maxIPCacheSize = 50000
	ipAPIBatchSize = 100
)

type IPLocation struct {
	Domestic string
	Global   string
}

type ipLocationCacheEntry struct {
	Domestic string
	Global   string
	Updated  time.Time
}

var (
	ipGeoCache   = make(map[string]ipLocationCacheEntry)
	ipGeoCacheMu sync.RWMutex
)

type ipAPIBatchRequest struct {
	Query  string `json:"query"`
	Fields string `json:"fields"`
	Lang   string `json:"lang"`
}

type ipAPIBatchResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Region      string `json:"region"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	Query       string `json:"query"`
}

// ExtractIPRegionDB 从嵌入的文件系统中提取 IP2Region 数据库
func ExtractIPRegionDB() (string, error) {
	// 确保数据目录存在
	if _, err := os.Stat(config.DataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(config.DataDir, 0755); err != nil {
			return "", err
		}
	}

	// 目标文件路径
	dbPath := filepath.Join(config.DataDir, "ip2region.xdb")

	// 检查文件是否已存在
	if _, err := os.Stat(dbPath); err == nil {
		logrus.Info("IP2Region 数据库已存在，跳过提取")
		return dbPath, nil
	}

	// 从嵌入文件系统读取数据
	data, err := fs.ReadFile(ipDataFiles, "data/ip2region.xdb")
	if err != nil {
		return "", err
	}

	// 写入文件
	if err := os.WriteFile(dbPath, data, 0644); err != nil {
		return "", err
	}

	logrus.Info("IP2Region 数据库已成功提取")
	return dbPath, nil
}

// InitIPGeoLocation 初始化 IP 地理位置查询
func InitIPGeoLocation() error {
	// 从嵌入的文件系统中提取数据库文件
	extractedPath, err := ExtractIPRegionDB()
	if err != nil {
		return fmt.Errorf("提取 ip2region 数据库失败: %v", err)
	}

	// 更新数据库路径
	dbPath = extractedPath

	// 加载矢量索引以加速搜索
	vIndex, err := xdb.LoadVectorIndexFromFile(dbPath)
	if err != nil {
		logrus.Warnf("加载 ip2region 矢量索引失败，将使用全量搜索: %v", err)
	} else {
		vectorIndex = vIndex
	}

	// 创建内存搜索器
	searcher, err := xdb.NewWithVectorIndex(dbPath, vectorIndex)
	if err != nil {
		return fmt.Errorf("创建 ip2region 搜索器失败: %v", err)
	}

	ipSearcher = searcher
	logrus.Info("ip2region 初始化成功")
	return nil
}

// GetIPLocation 获取 IP 的地理位置信息
func GetIPLocation(ip string) (string, string, error) {
	// 处理无效 IP
	if ip == "" || ip == "localhost" || ip == "127.0.0.1" || ip == "::1" {
		return "本地", "本地", nil
	}

	if domestic, global, ok := getCachedLocation(ip); ok {
		return domestic, global, nil
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "未知", "未知", fmt.Errorf("无效的 IP 地址")
	}

	// 检查是否是内网 IP
	if isPrivateIP(parsedIP) {
		return "内网", "本地网络", nil
	}

	var domestic, global string
	var err error

	domestic, global, err = queryIPLocationRemote(ip)
	if err != nil || domestic == "未知" {
		if parsedIP.To4() != nil {
			domesticLocal, globalLocal, localErr := queryIPLocationLocal(ip)
			if localErr == nil && domesticLocal != "未知" {
				domestic = domesticLocal
				global = globalLocal
				err = nil
			}
		}
	}

	if err != nil {
		return "未知", "未知", err
	}

	setCachedLocation(ip, domestic, global)

	return domestic, global, nil
}

// GetIPLocationBatch 批量获取 IP 的地理位置信息（优先远端）
func GetIPLocationBatch(ips []string) map[string]IPLocation {
	results := make(map[string]IPLocation, len(ips))
	if len(ips) == 0 {
		return results
	}

	unique := make([]string, 0, len(ips))
	seen := make(map[string]struct{}, len(ips))
	for _, raw := range ips {
		ip := strings.TrimSpace(raw)
		if ip == "" {
			continue
		}
		if _, ok := seen[ip]; ok {
			continue
		}
		seen[ip] = struct{}{}
		unique = append(unique, ip)
	}

	toQuery := make([]string, 0, len(unique))
	for _, ip := range unique {
		if domestic, global, ok := getCachedLocation(ip); ok {
			results[ip] = IPLocation{Domestic: domestic, Global: global}
			continue
		}

		if ip == "localhost" || ip == "127.0.0.1" || ip == "::1" {
			results[ip] = IPLocation{Domestic: "本地", Global: "本地"}
			setCachedLocation(ip, "本地", "本地")
			continue
		}

		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			results[ip] = IPLocation{Domestic: "未知", Global: "未知"}
			setCachedLocation(ip, "未知", "未知")
			continue
		}

		if isPrivateIP(parsedIP) {
			results[ip] = IPLocation{Domestic: "内网", Global: "本地网络"}
			setCachedLocation(ip, "内网", "本地网络")
			continue
		}

		toQuery = append(toQuery, ip)
	}

	if len(toQuery) == 0 {
		return results
	}

	remoteResults, _ := queryIPLocationRemoteBatch(toQuery)
	for _, ip := range toQuery {
		if entry, ok := remoteResults[ip]; ok && entry.Domestic != "" && entry.Domestic != "未知" {
			results[ip] = IPLocation{Domestic: entry.Domestic, Global: entry.Global}
			setCachedLocation(ip, entry.Domestic, entry.Global)
			continue
		}

		parsedIP := net.ParseIP(ip)
		if parsedIP != nil && parsedIP.To4() != nil {
			domestic, global, err := queryIPLocationLocal(ip)
			if err == nil && domestic != "" && domestic != "未知" {
				results[ip] = IPLocation{Domestic: domestic, Global: global}
				setCachedLocation(ip, domestic, global)
				continue
			}
		}

		results[ip] = IPLocation{Domestic: "未知", Global: "未知"}
		setCachedLocation(ip, "未知", "未知")
	}

	return results
}

// 查询 IP 地理位置（本地库）
func queryIPLocationLocal(ip string) (string, string, error) {
	if ipSearcher == nil {
		return "未知", "未知", fmt.Errorf("ip2region 未初始化")
	}

	// 设置 50 毫秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// 使用 channel 处理超时
	resultCh := make(chan struct {
		region string
		err    error
	}, 1)

	go func() {
		var region string
		var err error

		region, err = ipSearcher.SearchByStr(ip)

		resultCh <- struct {
			region string
			err    error
		}{region, err}
	}()

	// 等待结果或超时
	select {
	case <-ctx.Done():
		return "未知", "未知", fmt.Errorf("IP 查询超时")
	case result := <-resultCh:
		if result.err != nil {
			return "未知", "未知", result.err
		}
		return parseIPRegion(result.region)
	}
}

// 查询 IP 地理位置（远程接口）
func queryIPLocationRemote(ip string) (string, string, error) {
	results, err := queryIPLocationRemoteBatch([]string{ip})
	if err != nil {
		return "未知", "未知", err
	}
	entry, ok := results[ip]
	if !ok {
		return "未知", "未知", fmt.Errorf("ip-api 返回为空")
	}
	return entry.Domestic, entry.Global, nil
}

func queryIPLocationRemoteBatch(ips []string) (map[string]ipLocationCacheEntry, error) {
	results := make(map[string]ipLocationCacheEntry, len(ips))
	if len(ips) == 0 {
		return results, nil
	}

	client := &http.Client{Timeout: ipAPITimeout}
	var lastErr error

	for start := 0; start < len(ips); start += ipAPIBatchSize {
		end := start + ipAPIBatchSize
		if end > len(ips) {
			end = len(ips)
		}

		batch := ips[start:end]
		requestPayload := make([]ipAPIBatchRequest, 0, len(batch))
		for _, ip := range batch {
			requestPayload = append(requestPayload, ipAPIBatchRequest{
				Query:  ip,
				Fields: ipAPIFields,
				Lang:   ipAPILanguage,
			})
		}

		requestBody, err := json.Marshal(requestPayload)
		if err != nil {
			lastErr = err
			continue
		}

		req, err := http.NewRequest(http.MethodPost, ipAPIBatchURL, bytes.NewReader(requestBody))
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "nginxpulse/1.0")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("ip-api 响应异常: %s", resp.Status)
			resp.Body.Close()
			continue
		}

		var payload []ipAPIBatchResponse
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			lastErr = err
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		for i, item := range payload {
			query := strings.TrimSpace(item.Query)
			if query == "" && i < len(batch) {
				query = batch[i]
			}
			if query == "" {
				continue
			}

			if item.Status != "" && item.Status != "success" {
				results[query] = ipLocationCacheEntry{Domestic: "未知", Global: "未知"}
				continue
			}

			domestic := formatDomesticLocation(item.Country, item.RegionName, item.City)
			global := formatGlobalLocation(item.Country)
			if domestic == "" {
				domestic = "未知"
			}
			if global == "" {
				global = "未知"
			}
			results[query] = ipLocationCacheEntry{Domestic: domestic, Global: global}
		}
	}

	return results, lastErr
}

// 解析 ip2region 返回的地区信息
func parseIPRegion(region string) (string, string, error) {
	// 返回格式: 国家|区域|省份|城市|ISP
	parts := splitRegion(region)
	var domestic, global string

	// 国内
	if parts[0] == "中国" {
		province := removeSuffixes(parts[2])
		city := removeSuffixes(parts[3])
		if province != "" && province != "0" {
			if city != "" && city != "0" && city != province {
				domestic = fmt.Sprintf("%s·%s", province, city)
			} else {
				domestic = province
			}
		} else if city != "" && city != "0" {
			domestic = city
		} else {
			domestic = "中国"
		}
	} else if parts[0] == "0" || parts[0] == "" {
		domestic = "未知"
	} else {
		domestic = joinLocationParts(parts[0], parts[2], parts[3])
	}

	// 全球
	if parts[0] != "0" && parts[0] != "" {
		global = parts[0]
	} else {
		global = "未知"
	}

	return domestic, global, nil
}

// 解析 ip2region
func splitRegion(region string) []string {
	parts := make([]string, 5)
	fields := bytes.Split([]byte(region), []byte("|"))

	for i := 0; i < len(fields) && i < 5; i++ {
		parts[i] = string(fields[i])
	}

	return parts
}

func formatDomesticLocation(country, regionName, city string) string {
	country = strings.TrimSpace(country)
	if country == "" || country == "0" {
		return "未知"
	}
	if country != "中国" {
		return joinLocationParts(country, regionName, city)
	}
	province := removeSuffixes(strings.TrimSpace(regionName))
	city = removeSuffixes(strings.TrimSpace(city))
	if province == "" && city == "" {
		return "中国"
	}
	if province != "" && city != "" && province == city {
		return province
	}
	return joinLocationParts(province, city)
}

func formatGlobalLocation(country string) string {
	country = strings.TrimSpace(country)
	if country == "" || country == "0" {
		return "未知"
	}
	return country
}

func joinLocationParts(parts ...string) string {
	normalized := make([]string, 0, len(parts))
	for _, part := range parts {
		clean := normalizeLocationPart(part)
		if clean != "" {
			normalized = append(normalized, clean)
		}
	}
	if len(normalized) == 0 {
		return "未知"
	}
	return strings.Join(normalized, "·")
}

func normalizeLocationPart(value string) string {
	clean := strings.TrimSpace(value)
	if clean == "" || clean == "0" || clean == "未知" {
		return ""
	}
	return clean
}

func getCachedLocation(ip string) (string, string, bool) {
	ipGeoCacheMu.RLock()
	entry, ok := ipGeoCache[ip]
	ipGeoCacheMu.RUnlock()
	if !ok {
		return "", "", false
	}
	return entry.Domestic, entry.Global, true
}

func setCachedLocation(ip, domestic, global string) {
	if ip == "" {
		return
	}
	ipGeoCacheMu.Lock()
	defer ipGeoCacheMu.Unlock()
	if len(ipGeoCache) >= maxIPCacheSize {
		ipGeoCache = make(map[string]ipLocationCacheEntry)
	}
	ipGeoCache[ip] = ipLocationCacheEntry{
		Domestic: domestic,
		Global:   global,
		Updated:  time.Now(),
	}
}

// 是否是内网 IP
func isPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}

	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	if v4 := ip.To4(); v4 != nil {
		switch {
		case v4[0] == 10:
			return true
		case v4[0] == 172 && v4[1] >= 16 && v4[1] <= 31:
			return true
		case v4[0] == 192 && v4[1] == 168:
			return true
		case v4[0] == 127:
			return true
		case v4[0] == 169 && v4[1] == 254:
			return true
		default:
			return false
		}
	}

	ip = ip.To16()
	if ip == nil {
		return false
	}

	// IPv6 ULA fc00::/7
	if ip[0]&0xfe == 0xfc {
		return true
	}

	return false
}

// 去掉地区名称后缀
func removeSuffixes(name string) string {
	suffixes := []string{"省", "市", "自治区", "维吾尔自治区", "壮族自治区", "回族自治区", "特别行政区"}
	for _, suffix := range suffixes {
		if len(name) > len(suffix) && name[len(name)-len(suffix):] == suffix {
			return name[:len(name)-len(suffix)]
		}
	}
	return name
}

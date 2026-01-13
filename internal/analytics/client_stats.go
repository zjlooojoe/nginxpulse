package analytics

import (
	"fmt"
	"math"
	"net/url"
	"strings"

	"github.com/likaia/nginxpulse/internal/config"
	"github.com/likaia/nginxpulse/internal/store"
	"github.com/likaia/nginxpulse/internal/timeutil"
)

type ClientStats struct {
	Key       []string `json:"key"`        // 统计项的键
	PV        []int    `json:"pv"`         // 页面浏览量
	UV        []int    `json:"uv"`         // 独立访客数
	PVPercent []int    `json:"pv_percent"` // PV 百分比
	UVPercent []int    `json:"uv_percent"` // UV 百分比
}

func (s ClientStats) GetType() string {
	return "client"
}

type ClientStatsManager struct {
	repo      *store.Repository
	statsType string
}

func NewURLStatsManager(userRepoPtr *store.Repository) *ClientStatsManager {
	return &ClientStatsManager{
		repo:      userRepoPtr,
		statsType: "url",
	}
}

func NewrefererStatsManager(userRepoPtr *store.Repository) *ClientStatsManager {
	return &ClientStatsManager{
		repo:      userRepoPtr,
		statsType: "referer",
	}
}

func NewBrowserStatsManager(userRepoPtr *store.Repository) *ClientStatsManager {
	return &ClientStatsManager{
		repo:      userRepoPtr,
		statsType: "user_browser",
	}
}

func NewOsStatsManager(userRepoPtr *store.Repository) *ClientStatsManager {
	return &ClientStatsManager{
		repo:      userRepoPtr,
		statsType: "user_os",
	}
}

func NewDeviceStatsManager(userRepoPtr *store.Repository) *ClientStatsManager {
	return &ClientStatsManager{
		repo:      userRepoPtr,
		statsType: "user_device",
	}
}

func NewLocationStatsManager(userRepoPtr *store.Repository) *ClientStatsManager {
	return &ClientStatsManager{
		repo:      userRepoPtr,
		statsType: "location",
	}
}

// 实现 StatsManager 接口
func (s *ClientStatsManager) Query(query StatsQuery) (StatsResult, error) {
	result := ClientStats{
		Key:       make([]string, 0),
		PV:        make([]int, 0),
		UV:        make([]int, 0),
		PVPercent: make([]int, 0),
		UVPercent: make([]int, 0),
	}

	statsType := s.statsType
	locationType := ""
	if s.statsType == "location" {
		locationType = query.ExtraParam["locationType"].(string)
		switch locationType {
		case "domestic", "city":
			statsType = "domestic_location"
		case "global":
			statsType = "global_location"
		default:
			statsType = locationType + "_location"
		}
	}
	selectExpr := statsType
	groupExpr := statsType
	if s.statsType == "location" && locationType == "domestic" {
		selectExpr = fmt.Sprintf("CASE WHEN instr(%[1]s, '·') > 0 THEN substr(%[1]s, 1, instr(%[1]s, '·') - 1) ELSE %[1]s END", statsType)
		groupExpr = selectExpr
	}
	if s.statsType == "location" && locationType == "city" {
		selectExpr = fmt.Sprintf("CASE WHEN instr(%[1]s, '·') > 0 THEN substr(%[1]s, instr(%[1]s, '·') + 1) ELSE %[1]s END", statsType)
		groupExpr = selectExpr
	}
	if s.statsType == "referer" {
		internalCond := ""
		if website, ok := config.GetWebsiteByID(query.WebsiteID); ok {
			internalCond = buildInternalRefererCondition(website.Domains)
		}
		if internalCond != "" {
			selectExpr = fmt.Sprintf(
				"CASE WHEN referer = '-' OR referer = '' THEN '直接输入网址访问' WHEN %s THEN '站内访问' ELSE referer END",
				internalCond,
			)
		} else {
			selectExpr = "CASE WHEN referer = '-' OR referer = '' THEN '直接输入网址访问' ELSE referer END"
		}
		groupExpr = selectExpr
	}
	limit, _ := query.ExtraParam["limit"].(int)
	timeRange := query.ExtraParam["timeRange"].(string)
	startTime, endTime, err := timeutil.TimePeriod(timeRange)
	if err != nil {
		return result, err
	}

	extraCondition := ""
	if s.statsType == "location" && (locationType == "domestic" || locationType == "city") {
		extraCondition = " AND global_location = '中国'"
	}

	// 构建、执行查询
	dbQueryStr := fmt.Sprintf(`
        SELECT 
            %[1]s AS url, 
            COUNT(*) AS pv,
            COUNT(DISTINCT ip) AS uv
        FROM "%[2]s_nginx_logs" INDEXED BY idx_%[2]s_pv_ts_ip
        WHERE pageview_flag = 1 AND timestamp >= ? AND timestamp < ?%[4]s
        GROUP BY %[3]s
        ORDER BY uv DESC
        LIMIT ?`,
		selectExpr, query.WebsiteID, groupExpr, extraCondition)

	rows, err := s.repo.GetDB().Query(dbQueryStr, startTime.Unix(), endTime.Unix(), limit)
	if err != nil {
		return result, fmt.Errorf("查询URL统计失败: %v", err)
	}
	defer rows.Close()

	totalPV := 0
	totalUV := 0

	for rows.Next() {
		var url string
		var pv, uv int
		if err := rows.Scan(&url, &pv, &uv); err != nil {
			return result, fmt.Errorf("解析URL统计结果失败: %v", err)
		}
		result.Key = append(result.Key, url)
		result.PV = append(result.PV, pv)
		result.UV = append(result.UV, uv)
		totalPV += pv
		totalUV += uv
	}

	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("遍历URL统计结果失败: %v", err)
	}

	if totalPV > 0 && totalUV > 0 {
		for i := range result.PV {
			result.PVPercent = append(
				result.PVPercent, int(
					math.Round(float64(result.PV[i])/float64(totalPV)*100)))
			result.UVPercent = append(
				result.UVPercent, int(
					math.Round(float64(result.UV[i])/float64(totalUV)*100)))
		}
	}

	return result, nil

}

func buildInternalRefererCondition(domains []string) string {
	conditions := make([]string, 0, len(domains))
	for _, raw := range domains {
		domain := normalizeDomain(raw)
		if domain == "" {
			continue
		}
		domain = strings.ReplaceAll(domain, "'", "''")
		conditions = append(conditions,
			fmt.Sprintf(
				"referer LIKE 'http%%://%s/%%' OR referer LIKE 'http%%://%s' OR referer LIKE 'http%%://%s:%%'",
				domain, domain, domain,
			),
		)
	}
	if len(conditions) == 0 {
		return ""
	}
	return fmt.Sprintf("(%s)", strings.Join(conditions, " OR "))
}

func normalizeDomain(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.Contains(raw, "://") {
		parsed, err := url.Parse(raw)
		if err == nil && parsed.Host != "" {
			return strings.TrimSuffix(parsed.Host, "/")
		}
	}
	raw = strings.TrimPrefix(raw, "//")
	raw = strings.TrimSuffix(raw, "/")
	return raw
}

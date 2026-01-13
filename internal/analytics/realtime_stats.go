package analytics

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/likaia/nginxpulse/internal/config"
	"github.com/likaia/nginxpulse/internal/store"
)

type RealtimeItem struct {
	Name    string  `json:"name"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

type RealtimeStats struct {
	WindowMinutes   int            `json:"windowMinutes"`
	ActiveCount     int            `json:"activeCount"`
	ActiveSeries    []int          `json:"activeSeries"`
	DeviceBreakdown []RealtimeItem `json:"deviceBreakdown"`
	Referers        []RealtimeItem `json:"referers"`
	Pages           []RealtimeItem `json:"pages"`
	EntryPages      []RealtimeItem `json:"entryPages"`
	Browsers        []RealtimeItem `json:"browsers"`
	Locations       []RealtimeItem `json:"locations"`
}

func (s RealtimeStats) GetType() string {
	return "realtime"
}

type RealtimeStatsManager struct {
	repo *store.Repository
}

func NewRealtimeStatsManager(userRepoPtr *store.Repository) *RealtimeStatsManager {
	return &RealtimeStatsManager{
		repo: userRepoPtr,
	}
}

func (m *RealtimeStatsManager) Query(query StatsQuery) (StatsResult, error) {
	result := RealtimeStats{
		WindowMinutes: 30,
	}

	window := 30
	if val, ok := query.ExtraParam["window"].(int); ok && val > 0 {
		window = val
	}
	if window > 60 {
		window = 60
	}
	result.WindowMinutes = window

	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(window) * time.Minute)

	tableName := fmt.Sprintf("%s_nginx_logs", query.WebsiteID)

	activeCount, err := m.activeVisitorCount(tableName, startTime, endTime)
	if err != nil {
		return result, err
	}
	result.ActiveCount = activeCount

	series, err := m.activeSeries(tableName, startTime, endTime, window)
	if err != nil {
		return result, err
	}
	result.ActiveSeries = series

	result.DeviceBreakdown = m.deviceBreakdown(tableName, startTime, endTime)

	refererExpr := buildRealtimeRefererExpr(query.WebsiteID)
	referers, _ := m.queryTopItems(tableName, refererExpr, refererExpr, startTime, endTime, 10, true)
	result.Referers = referers

	pages, _ := m.queryTopItems(tableName, "url", "url", startTime, endTime, 10, false)
	result.Pages = pages

	entryCounts, _ := m.entryPages(tableName, startTime, endTime)
	result.EntryPages = entryCounts

	browsers, _ := m.queryTopItems(tableName, "user_browser", "user_browser", startTime, endTime, 10, true)
	result.Browsers = browsers

	locationExpr := "CASE WHEN instr(domestic_location, '·') > 0 THEN substr(domestic_location, instr(domestic_location, '·') + 1) ELSE domestic_location END"
	locations, _ := m.queryTopItems(
		tableName,
		locationExpr,
		locationExpr,
		startTime,
		endTime,
		10,
		true,
	)
	result.Locations = locations

	return result, nil
}

func (m *RealtimeStatsManager) activeVisitorCount(tableName string, startTime, endTime time.Time) (int, error) {
	query := fmt.Sprintf(`
        SELECT COUNT(DISTINCT ip)
        FROM "%s"
        WHERE pageview_flag = 1 AND timestamp >= ? AND timestamp < ?`,
		tableName)

	row := m.repo.GetDB().QueryRow(query, startTime.Unix(), endTime.Unix())
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (m *RealtimeStatsManager) activeSeries(tableName string, startTime, endTime time.Time, window int) ([]int, error) {
	query := fmt.Sprintf(`
        SELECT (timestamp / 60) as bucket, COUNT(DISTINCT ip) as uv
        FROM "%s"
        WHERE pageview_flag = 1 AND timestamp >= ? AND timestamp < ?
        GROUP BY bucket`,
		tableName)

	rows, err := m.repo.GetDB().Query(query, startTime.Unix(), endTime.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	buckets := make(map[int64]int)
	for rows.Next() {
		var bucket int64
		var uv int
		if err := rows.Scan(&bucket, &uv); err != nil {
			return nil, err
		}
		buckets[bucket] = uv
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	startBucket := startTime.Unix() / 60
	series := make([]int, window)
	for i := 0; i < window; i++ {
		series[i] = buckets[startBucket+int64(i)]
	}
	return series, nil
}

func (m *RealtimeStatsManager) deviceBreakdown(tableName string, startTime, endTime time.Time) []RealtimeItem {
	query := fmt.Sprintf(`
        SELECT user_device, COUNT(DISTINCT ip) as uv
        FROM "%s"
        WHERE pageview_flag = 1 AND timestamp >= ? AND timestamp < ?
        GROUP BY user_device`,
		tableName)

	rows, err := m.repo.GetDB().Query(query, startTime.Unix(), endTime.Unix())
	if err != nil {
		return []RealtimeItem{}
	}
	defer rows.Close()

	var (
		pc     int
		mobile int
		other  int
		total  int
	)

	for rows.Next() {
		var device string
		var count int
		if err := rows.Scan(&device, &count); err != nil {
			return []RealtimeItem{}
		}
		switch device {
		case "桌面设备":
			pc += count
		case "手机", "平板":
			mobile += count
		default:
			other += count
		}
		total += count
	}

	return []RealtimeItem{
		{Name: "PC端", Count: pc, Percent: safePercent(pc, total)},
		{Name: "移动端", Count: mobile, Percent: safePercent(mobile, total)},
		{Name: "其他", Count: other, Percent: safePercent(other, total)},
	}
}

func (m *RealtimeStatsManager) queryTopItems(
	tableName string,
	selectExpr string,
	groupExpr string,
	startTime, endTime time.Time,
	limit int,
	distinctIP bool,
) ([]RealtimeItem, error) {
	if limit <= 0 {
		limit = 10
	}
	countExpr := "COUNT(*)"
	if distinctIP {
		countExpr = "COUNT(DISTINCT ip)"
	}

	query := fmt.Sprintf(`
        SELECT %[1]s as key, %[2]s as cnt
        FROM "%[3]s"
        WHERE pageview_flag = 1 AND timestamp >= ? AND timestamp < ?
        GROUP BY %[4]s
        ORDER BY cnt DESC
        LIMIT ?`,
		selectExpr, countExpr, tableName, groupExpr)

	rows, err := m.repo.GetDB().Query(query, startTime.Unix(), endTime.Unix(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type item struct {
		key   string
		count int
	}
	items := make([]item, 0)
	total := 0
	for rows.Next() {
		var key string
		var count int
		if err := rows.Scan(&key, &count); err != nil {
			return nil, err
		}
		key = strings.TrimSpace(key)
		if key == "" {
			key = "未知"
		}
		items = append(items, item{key: key, count: count})
		total += count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]RealtimeItem, 0, len(items))
	for _, item := range items {
		result = append(result, RealtimeItem{
			Name:    item.key,
			Count:   item.count,
			Percent: safePercent(item.count, total),
		})
	}
	return result, nil
}

func (m *RealtimeStatsManager) entryPages(
	tableName string,
	startTime, endTime time.Time,
) ([]RealtimeItem, error) {
	query := fmt.Sprintf(`
        SELECT timestamp, ip, user_browser, user_os, user_device, url
        FROM "%s"
        WHERE pageview_flag = 1 AND timestamp >= ? AND timestamp < ?
        ORDER BY ip, user_browser, user_os, user_device, timestamp`,
		tableName)

	rows, err := m.repo.GetDB().Query(query, startTime.Unix(), endTime.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entryCounts := make(map[string]int)
	var (
		currentKey    string
		lastTimestamp int64
		initialized   bool
	)

	for rows.Next() {
		var (
			timestamp int64
			ip        string
			browser   string
			os        string
			device    string
			url       string
		)

		if err := rows.Scan(&timestamp, &ip, &browser, &os, &device, &url); err != nil {
			return nil, err
		}

		key := fmt.Sprintf("%s|%s|%s|%s", ip, browser, os, device)
		if !initialized || key != currentKey || timestamp-lastTimestamp > sessionGapSeconds {
			entryCounts[url]++
			currentKey = key
			initialized = true
		}
		lastTimestamp = timestamp
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	type entry struct {
		Key   string
		Count int
	}
	entries := make([]entry, 0, len(entryCounts))
	total := 0
	for key, count := range entryCounts {
		entries = append(entries, entry{Key: key, Count: count})
		total += count
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Count > entries[j].Count
	})

	limit := 10
	if len(entries) < limit {
		limit = len(entries)
	}

	result := make([]RealtimeItem, 0, limit)
	for i := 0; i < limit; i++ {
		item := entries[i]
		result = append(result, RealtimeItem{
			Name:    item.Key,
			Count:   item.Count,
			Percent: safePercent(item.Count, total),
		})
	}
	return result, nil
}

func buildRealtimeRefererExpr(websiteID string) string {
	internalCond := ""
	if website, ok := config.GetWebsiteByID(websiteID); ok {
		internalCond = buildInternalRefererCondition(website.Domains)
	}
	if internalCond != "" {
		return fmt.Sprintf(
			"CASE WHEN referer = '-' OR referer = '' THEN '直接输入网址访问' WHEN %s THEN '站内访问' ELSE referer END",
			internalCond,
		)
	}
	return "CASE WHEN referer = '-' OR referer = '' THEN '直接输入网址访问' ELSE referer END"
}

func safePercent(value, total int) float64 {
	if total <= 0 {
		return 0
	}
	return float64(value) / float64(total)
}

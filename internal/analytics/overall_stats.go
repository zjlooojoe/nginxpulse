package analytics

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/likaia/nginxpulse/internal/store"
	"github.com/likaia/nginxpulse/internal/timeutil"
	"github.com/sirupsen/logrus"
)

type OverallStats struct {
	PV                        int            `json:"pv"`                        // 页面浏览量
	UV                        int            `json:"uv"`                        // 独立访客数
	Traffic                   int64          `json:"traffic"`                   // 流量（字节）
	SessionCount              int            `json:"sessionCount"`              // 会话数
	ActiveVisitorCount        int            `json:"activeVisitorCount"`        // 最近15分钟活跃访客数
	NewVisitorCount           int            `json:"newVisitorCount"`           // 新访客数
	ReturningVisitorCount     int            `json:"returningVisitorCount"`     // 老访客数
	PrevNewVisitorCount       int            `json:"prevNewVisitorCount"`       // 上期新访客数
	PrevReturningVisitorCount int            `json:"prevReturningVisitorCount"` // 上期老访客数
	EntryPages                ClientStats    `json:"entryPages"`                // 入口页统计
	Compare                   OverallCompare `json:"compare"`                   // 对比数据
	StatusCodeHits            StatusCodeHits `json:"statusCodeHits"`            // HTTP 状态码命中次数
	StatusCodeHitsPrevious    StatusCodeHits `json:"statusCodeHitsPrevious"`    // 上一期状态码命中次数
}

type OverallSnapshot struct {
	PV           int `json:"pv"`
	UV           int `json:"uv"`
	SessionCount int `json:"sessionCount"`
}

type OverallCompare struct {
	Previous OverallSnapshot `json:"previous"`
	Forecast OverallSnapshot `json:"forecast"`
	SameTime OverallSnapshot `json:"sameTime"`
}

type StatusCodeHits struct {
	S2xx  int `json:"s2xx"`
	S3xx  int `json:"s3xx"`
	S4xx  int `json:"s4xx"`
	S5xx  int `json:"s5xx"`
	Other int `json:"other"`
}

// OverallStats 实现 StatsResult 接口
func (s OverallStats) GetType() string {
	return "overall"
}

type OverallStatsManager struct {
	repo *store.Repository
}

// NewOverallStatsManager 创建一个新的 OverallStatsManager 实例
func NewOverallStatsManager(userRepoPtr *store.Repository) *OverallStatsManager {
	return &OverallStatsManager{
		repo: userRepoPtr,
	}
}

// 实现 StatsManager 接口
func (s *OverallStatsManager) Query(query StatsQuery) (StatsResult, error) {

	result := OverallStats{
		PV:                        0,
		UV:                        0,
		Traffic:                   0,
		SessionCount:              0,
		ActiveVisitorCount:        0,
		NewVisitorCount:           0,
		ReturningVisitorCount:     0,
		PrevNewVisitorCount:       0,
		PrevReturningVisitorCount: 0,
		EntryPages: ClientStats{
			Key:       make([]string, 0),
			PV:        make([]int, 0),
			UV:        make([]int, 0),
			PVPercent: make([]int, 0),
			UVPercent: make([]int, 0),
		},
		Compare: OverallCompare{
			Previous: OverallSnapshot{},
			Forecast: OverallSnapshot{},
			SameTime: OverallSnapshot{},
		},
		StatusCodeHits:         StatusCodeHits{},
		StatusCodeHitsPrevious: StatusCodeHits{},
	}

	timeRange := query.ExtraParam["timeRange"].(string)
	startTime, endTime, err := timeutil.TimePeriod(timeRange)
	if err != nil {
		return result, err
	}
	prevStart, prevEnd := previousTimeRange(timeRange)
	entryLimit := 10
	if rawLimit, ok := query.ExtraParam["entryLimit"]; ok {
		if limit, ok := rawLimit.(int); ok && limit > 0 {
			entryLimit = limit
		}
	}

	err = s.statsByTimeRangeForWebsite(query.WebsiteID, startTime, endTime, &result)
	if err != nil {
		return result, fmt.Errorf("获取总体统计失败: %v", err)
	}

	statusHits, err := s.statusCodeHitsByTimeRangeForWebsite(query.WebsiteID, startTime, endTime)
	if err != nil {
		logrus.WithError(err).Warn("获取状态码统计失败")
	} else {
		result.StatusCodeHits = statusHits
	}

	if !prevStart.IsZero() && !prevEnd.IsZero() {
		prevStatusHits, err := s.statusCodeHitsByTimeRangeForWebsite(query.WebsiteID, prevStart, prevEnd)
		if err != nil {
			logrus.WithError(err).Warn("获取上一期状态码统计失败")
		} else {
			result.StatusCodeHitsPrevious = prevStatusHits
		}
	}

	metrics, err := collectSessionMetrics(s.repo, query.WebsiteID, startTime, endTime)
	if err != nil {
		logrus.WithError(err).Warn("获取会话统计失败")
	} else {
		result.SessionCount = metrics.SessionCount
		result.EntryPages = buildEntryStats(metrics.EntryCounts, entryLimit)
	}

	activeCount, err := s.activeVisitorCount(query.WebsiteID)
	if err != nil {
		logrus.WithError(err).Warn("获取活跃访客失败")
	} else {
		result.ActiveVisitorCount = activeCount
	}

	newCount, returningCount, err := s.newReturningCounts(query.WebsiteID, startTime, endTime)
	if err != nil {
		logrus.WithError(err).Warn("获取新老访客失败")
	} else {
		result.NewVisitorCount = newCount
		result.ReturningVisitorCount = returningCount
	}

	if !prevStart.IsZero() && !prevEnd.IsZero() {
		prevNew, prevReturning, err := s.newReturningCounts(query.WebsiteID, prevStart, prevEnd)
		if err != nil {
			logrus.WithError(err).Warn("获取上期新老访客失败")
		} else {
			result.PrevNewVisitorCount = prevNew
			result.PrevReturningVisitorCount = prevReturning
		}
	}

	currentSnapshot := snapshotFromOverall(result)
	prevSnapshot, prevSameSnapshot, forecastSnapshot := s.buildCompareSnapshots(
		query.WebsiteID, timeRange, startTime, endTime, currentSnapshot,
	)
	result.Compare = OverallCompare{
		Previous: prevSnapshot,
		Forecast: forecastSnapshot,
		SameTime: prevSameSnapshot,
	}

	return result, nil
}

// StatsByTimePoints 直接使用 db.Query() 方法查询数据库获取指定时间点的统计数据
func (s *OverallStatsManager) statsByTimeRangeForWebsite(
	websiteID string, startTime, endTime time.Time, overall *OverallStats) error {

	// 初始化结果
	overall.PV = 0
	overall.UV = 0
	overall.Traffic = 0

	tableName := fmt.Sprintf("%s_nginx_logs", websiteID)

	// 为更精确的统计，直接在数据库中进行全范围的唯一IP计数
	countQuery := fmt.Sprintf(`
        SELECT 
            COUNT(*) as pv,
            COUNT(DISTINCT ip) as uv,
            COALESCE(SUM(bytes_sent), 0) as traffic
        FROM "%s" INDEXED BY idx_%s_pv_ts_ip
        WHERE pageview_flag = 1 AND timestamp >= ? AND timestamp < ?`,
		tableName, websiteID)

	// 执行全范围查询
	row := s.repo.GetDB().QueryRow(countQuery, startTime.Unix(), endTime.Unix())

	if err := row.Scan(&overall.PV, &overall.UV, &overall.Traffic); err != nil {
		return fmt.Errorf("查询总体统计数据失败: %v", err)
	}

	return nil
}

func (s *OverallStatsManager) statusCodeHitsByTimeRangeForWebsite(
	websiteID string, startTime, endTime time.Time) (StatusCodeHits, error) {

	result := StatusCodeHits{}
	tableName := fmt.Sprintf("%s_nginx_logs", websiteID)

	query := fmt.Sprintf(`
        SELECT
            COALESCE(SUM(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 ELSE 0 END), 0) AS s2xx,
            COALESCE(SUM(CASE WHEN status_code >= 300 AND status_code < 400 THEN 1 ELSE 0 END), 0) AS s3xx,
            COALESCE(SUM(CASE WHEN status_code >= 400 AND status_code < 500 THEN 1 ELSE 0 END), 0) AS s4xx,
            COALESCE(SUM(CASE WHEN status_code >= 500 AND status_code < 600 THEN 1 ELSE 0 END), 0) AS s5xx,
            COALESCE(SUM(CASE WHEN status_code < 200 OR status_code >= 600 THEN 1 ELSE 0 END), 0) AS other
        FROM "%s" INDEXED BY idx_%s_timestamp
        WHERE timestamp >= ? AND timestamp < ?`,
		tableName, websiteID)

	row := s.repo.GetDB().QueryRow(query, startTime.Unix(), endTime.Unix())
	if err := row.Scan(&result.S2xx, &result.S3xx, &result.S4xx, &result.S5xx, &result.Other); err != nil {
		return result, fmt.Errorf("查询状态码统计失败: %v", err)
	}

	return result, nil
}

type sessionMetrics struct {
	SessionCount int
	EntryCounts  map[string]int
}

const sessionGapSeconds = int64(1800)

func collectSessionMetrics(
	repo *store.Repository,
	websiteID string,
	startTime, endTime time.Time,
) (sessionMetrics, error) {
	metrics := sessionMetrics{
		EntryCounts: make(map[string]int),
	}

	query := fmt.Sprintf(`
        SELECT timestamp, ip, user_browser, user_os, user_device, url
        FROM "%s_nginx_logs" INDEXED BY idx_%s_session_key
        WHERE pageview_flag = 1 AND timestamp >= ? AND timestamp < ?
        ORDER BY ip, user_browser, user_os, user_device, timestamp`,
		websiteID, websiteID)

	rows, err := repo.GetDB().Query(query, startTime.Unix(), endTime.Unix())
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

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
			return metrics, err
		}

		key := fmt.Sprintf("%s|%s|%s|%s", ip, browser, os, device)

		if !initialized || key != currentKey || timestamp-lastTimestamp > sessionGapSeconds {
			currentKey = key
			metrics.SessionCount++
			metrics.EntryCounts[url]++
			initialized = true
		}

		lastTimestamp = timestamp
	}

	if err := rows.Err(); err != nil {
		return metrics, err
	}

	return metrics, nil
}

func buildEntryStats(counts map[string]int, limit int) ClientStats {
	result := ClientStats{
		Key:       make([]string, 0),
		PV:        make([]int, 0),
		UV:        make([]int, 0),
		PVPercent: make([]int, 0),
		UVPercent: make([]int, 0),
	}

	if len(counts) == 0 {
		return result
	}

	type entryItem struct {
		Key   string
		Count int
	}

	items := make([]entryItem, 0, len(counts))
	total := 0
	for key, count := range counts {
		items = append(items, entryItem{Key: key, Count: count})
		total += count
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	if limit <= 0 || limit > len(items) {
		limit = len(items)
	}

	for i := 0; i < limit; i++ {
		item := items[i]
		result.Key = append(result.Key, item.Key)
		result.UV = append(result.UV, item.Count)
		if total > 0 {
			percent := int(math.Round(float64(item.Count) / float64(total) * 100))
			result.UVPercent = append(result.UVPercent, percent)
		} else {
			result.UVPercent = append(result.UVPercent, 0)
		}
	}

	return result
}

func (s *OverallStatsManager) activeVisitorCount(websiteID string) (int, error) {
	now := time.Now()
	start := now.Add(-15 * time.Minute)

	query := fmt.Sprintf(`
        SELECT COUNT(DISTINCT ip)
        FROM "%s_nginx_logs"
        WHERE pageview_flag = 1 AND timestamp >= ? AND timestamp < ?`,
		websiteID)

	row := s.repo.GetDB().QueryRow(query, start.Unix(), now.Unix())
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *OverallStatsManager) newReturningCounts(
	websiteID string, startTime, endTime time.Time,
) (int, int, error) {
	query := fmt.Sprintf(`
        WITH active_ips AS (
            SELECT DISTINCT ip
            FROM "%s_nginx_logs"
            WHERE pageview_flag = 1 AND timestamp >= ? AND timestamp < ?
        ),
        first_seen AS (
            SELECT ip, MIN(timestamp) AS first_ts
            FROM "%s_nginx_logs"
            WHERE pageview_flag = 1
            GROUP BY ip
        )
        SELECT
            COALESCE(SUM(CASE WHEN first_ts >= ? AND first_ts < ? THEN 1 ELSE 0 END), 0) AS new_uv,
            COALESCE(SUM(CASE WHEN first_ts < ? THEN 1 ELSE 0 END), 0) AS returning_uv
        FROM first_seen
        INNER JOIN active_ips USING (ip)`,
		websiteID, websiteID)

	row := s.repo.GetDB().QueryRow(
		query,
		startTime.Unix(), endTime.Unix(),
		startTime.Unix(), endTime.Unix(),
		startTime.Unix(),
	)

	var newCount, returningCount int
	if err := row.Scan(&newCount, &returningCount); err != nil {
		return 0, 0, err
	}

	return newCount, returningCount, nil
}

func previousTimeRange(timeRange string) (time.Time, time.Time) {
	now := time.Now()
	if len(timeRange) == 10 {
		if date, err := time.ParseInLocation("2006-01-02", timeRange, now.Location()); err == nil {
			prev := date.AddDate(0, 0, -1)
			start := time.Date(prev.Year(), prev.Month(), prev.Day(), 0, 0, 0, 0, prev.Location())
			end := time.Date(prev.Year(), prev.Month(), prev.Day(), 23, 59, 59, 0, prev.Location())
			return start, end
		}
	}
	switch timeRange {
	case "today":
		day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		start := day.AddDate(0, 0, -1)
		end := time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())
		return start, end
	case "yesterday":
		day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		start := day.AddDate(0, 0, -2)
		end := time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())
		return start, end
	case "last7days":
		end := time.Date(now.Year(), now.Month(), now.Day()-7, 23, 59, 59, 0, now.Location())
		start := time.Date(now.Year(), now.Month(), now.Day()-13, 0, 0, 0, 0, now.Location())
		return start, end
	case "last30days":
		end := time.Date(now.Year(), now.Month(), now.Day()-30, 23, 59, 59, 0, now.Location())
		start := time.Date(now.Year(), now.Month(), now.Day()-59, 0, 0, 0, 0, now.Location())
		return start, end
	case "week":
		start, end, _ := timeutil.TimePeriod("week")
		return start.AddDate(0, 0, -7), end.AddDate(0, 0, -7)
	case "month":
		start, _, _ := timeutil.TimePeriod("month")
		prevEnd := start.Add(-time.Second)
		prevStart := time.Date(prevEnd.Year(), prevEnd.Month(), 1, 0, 0, 0, 0, prevEnd.Location())
		return prevStart, prevEnd
	default:
		day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		start := day.AddDate(0, 0, -1)
		end := time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())
		return start, end
	}
}

func snapshotFromOverall(overall OverallStats) OverallSnapshot {
	return OverallSnapshot{
		PV:           overall.PV,
		UV:           overall.UV,
		SessionCount: overall.SessionCount,
	}
}

func (s *OverallStatsManager) buildCompareSnapshots(
	websiteID string,
	timeRange string,
	startTime, endTime time.Time,
	current OverallSnapshot,
) (OverallSnapshot, OverallSnapshot, OverallSnapshot) {
	prevStart, prevEnd := previousTimeRange(timeRange)
	if prevStart.IsZero() || prevEnd.IsZero() {
		return OverallSnapshot{}, current, current
	}

	prevSnapshot, err := s.snapshotForRange(websiteID, prevStart, prevEnd)
	if err != nil {
		logrus.WithError(err).Warn("获取上一期统计失败")
	}

	now := time.Now()
	currentEnd := now
	if currentEnd.After(endTime) {
		currentEnd = endTime
	}
	if currentEnd.Before(startTime) {
		currentEnd = startTime
	}
	elapsed := currentEnd.Sub(startTime)
	total := endTime.Sub(startTime)
	progress := 1.0
	if total > 0 {
		progress = elapsed.Seconds() / total.Seconds()
		if progress <= 0 {
			progress = 1.0
		}
		if progress > 1 {
			progress = 1.0
		}
	}

	progressForecast := scaleSnapshot(current, progress)
	forecast := s.forecastSnapshot(websiteID, startTime, endTime, currentEnd, progressForecast)

	prevSameEnd := prevStart.Add(elapsed)
	if prevSameEnd.After(prevEnd) {
		prevSameEnd = prevEnd
	}
	prevSameSnapshot := prevSnapshot
	if prevSameEnd.After(prevStart) {
		prevSameSnapshot, err = s.snapshotForRange(websiteID, prevStart, prevSameEnd)
		if err != nil {
			logrus.WithError(err).Warn("获取上一期同期失败")
			prevSameSnapshot = prevSnapshot
		}
	}

	return prevSnapshot, prevSameSnapshot, forecast
}

func scaleSnapshot(current OverallSnapshot, progress float64) OverallSnapshot {
	if progress <= 0 {
		return current
	}
	if progress >= 1 {
		return current
	}

	scale := 1 / progress
	return OverallSnapshot{
		PV:           int(math.Round(float64(current.PV) * scale)),
		UV:           int(math.Round(float64(current.UV) * scale)),
		SessionCount: int(math.Round(float64(current.SessionCount) * scale)),
	}
}

func (s *OverallStatsManager) snapshotForRange(
	websiteID string, startTime, endTime time.Time,
) (OverallSnapshot, error) {
	overall := OverallStats{}
	if err := s.statsByTimeRangeForWebsite(websiteID, startTime, endTime, &overall); err != nil {
		return OverallSnapshot{}, err
	}

	metrics, err := collectSessionMetrics(s.repo, websiteID, startTime, endTime)
	if err != nil {
		return OverallSnapshot{}, err
	}

	snapshot := OverallSnapshot{
		PV:           overall.PV,
		UV:           overall.UV,
		SessionCount: metrics.SessionCount,
	}

	return snapshot, nil
}

func (s *OverallStatsManager) forecastSnapshot(
	websiteID string,
	startTime, endTime, currentEnd time.Time,
	progressForecast OverallSnapshot,
) OverallSnapshot {
	if currentEnd.Before(startTime) {
		return progressForecast
	}

	total := endTime.Sub(startTime)
	if total <= 0 {
		return progressForecast
	}

	remaining := endTime.Sub(currentEnd)
	if remaining <= 0 {
		return progressForecast
	}

	elapsed := currentEnd.Sub(startTime)
	windowDuration := 2 * time.Hour
	if elapsed < windowDuration {
		windowDuration = elapsed
	}
	minWindow := 30 * time.Minute
	if windowDuration < minWindow {
		return progressForecast
	}

	windowStart := currentEnd.Add(-windowDuration)
	if windowStart.Before(startTime) {
		windowStart = startTime
	}

	windowSnapshot, err := s.snapshotForRange(websiteID, windowStart, currentEnd)
	if err != nil {
		logrus.WithError(err).Warn("获取预测窗口数据失败")
		return progressForecast
	}

	windowSeconds := windowDuration.Seconds()
	if windowSeconds <= 0 {
		return progressForecast
	}

	if windowSnapshot.PV == 0 && windowSnapshot.UV == 0 && windowSnapshot.SessionCount == 0 {
		return progressForecast
	}

	rateForecast := OverallSnapshot{
		PV:           progressForecast.PV,
		UV:           progressForecast.UV,
		SessionCount: progressForecast.SessionCount,
	}

	rateForecast.PV = forecastCount(
		progressForecast.PV,
		windowSnapshot.PV,
		windowSeconds,
		remaining.Seconds(),
	)
	rateForecast.UV = forecastCount(
		progressForecast.UV,
		windowSnapshot.UV,
		windowSeconds,
		remaining.Seconds(),
	)
	rateForecast.SessionCount = forecastCount(
		progressForecast.SessionCount,
		windowSnapshot.SessionCount,
		windowSeconds,
		remaining.Seconds(),
	)

	weight := forecastWeight(total)
	return blendSnapshot(rateForecast, progressForecast, weight)
}

func forecastCount(current, windowCount int, windowSeconds, remainingSeconds float64) int {
	if windowSeconds <= 0 || remainingSeconds <= 0 {
		return current
	}
	rate := float64(windowCount) / windowSeconds
	projected := float64(current) + rate*remainingSeconds
	if projected < float64(current) {
		projected = float64(current)
	}
	return int(math.Round(projected))
}

func forecastWeight(total time.Duration) float64 {
	switch {
	case total <= 24*time.Hour:
		return 0.7
	case total <= 7*24*time.Hour:
		return 0.5
	default:
		return 0.3
	}
}

func blendSnapshot(rateSnapshot, progressSnapshot OverallSnapshot, weight float64) OverallSnapshot {
	if weight <= 0 {
		return progressSnapshot
	}
	if weight >= 1 {
		return rateSnapshot
	}

	blend := func(rateVal, progressVal int) int {
		return int(math.Round(float64(rateVal)*weight + float64(progressVal)*(1-weight)))
	}

	result := progressSnapshot
	result.PV = blend(rateSnapshot.PV, progressSnapshot.PV)
	result.UV = blend(rateSnapshot.UV, progressSnapshot.UV)
	result.SessionCount = blend(rateSnapshot.SessionCount, progressSnapshot.SessionCount)
	return result
}

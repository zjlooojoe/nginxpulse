package ingest

import (
	"math"
	"sync"
	"time"
)

type parseProgressState struct {
	TotalBytes     int64
	ProcessedBytes int64
	StartedAt      time.Time
	UpdatedAt      time.Time
}

var (
	parseProgressMu sync.RWMutex
	parseProgress   parseProgressState
)

func resetParsingProgress() {
	now := time.Now()
	parseProgressMu.Lock()
	parseProgress = parseProgressState{
		TotalBytes:     0,
		ProcessedBytes: 0,
		StartedAt:      now,
		UpdatedAt:      now,
	}
	parseProgressMu.Unlock()
}

func setParsingTotalBytes(totalBytes int64) {
	if totalBytes < 0 {
		totalBytes = 0
	}
	parseProgressMu.Lock()
	parseProgress.TotalBytes = totalBytes
	if totalBytes > 0 && parseProgress.ProcessedBytes > totalBytes {
		parseProgress.ProcessedBytes = totalBytes
	}
	parseProgress.UpdatedAt = time.Now()
	parseProgressMu.Unlock()
}

func addParsingProgress(deltaBytes int64) {
	if deltaBytes <= 0 {
		return
	}
	parseProgressMu.Lock()
	parseProgress.ProcessedBytes += deltaBytes
	if parseProgress.TotalBytes > 0 && parseProgress.ProcessedBytes > parseProgress.TotalBytes {
		parseProgress.ProcessedBytes = parseProgress.TotalBytes
	}
	parseProgress.UpdatedAt = time.Now()
	parseProgressMu.Unlock()
}

func finalizeParsingProgress() {
	parseProgressMu.Lock()
	if parseProgress.TotalBytes > 0 {
		parseProgress.ProcessedBytes = parseProgress.TotalBytes
	}
	parseProgress.UpdatedAt = time.Now()
	parseProgressMu.Unlock()
}

func GetIPParsingProgress() int {
	parseProgressMu.RLock()
	total := parseProgress.TotalBytes
	processed := parseProgress.ProcessedBytes
	parseProgressMu.RUnlock()

	if total <= 0 {
		return 0
	}

	progress := float64(processed) / float64(total)
	progress = math.Max(0, math.Min(progress, 1))
	return int(math.Round(progress * 100))
}

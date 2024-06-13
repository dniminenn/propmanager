package service

import "runtime"

// StatsService represents the service for statistics.
type StatsService struct{}

// NewStatsService returns a new StatsService.
func NewStatsService() *StatsService {
	return &StatsService{}
}

// GetStats returns system statistics such as memory usage and performance.
func (s *StatsService) GetStats() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"Alloc":      memStats.Alloc,
		"TotalAlloc": memStats.TotalAlloc,
		"Sys":        memStats.Sys,
		"NumGC":      memStats.NumGC,
	}
}

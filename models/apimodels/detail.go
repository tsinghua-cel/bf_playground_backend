package apimodels

type ProjectDetail struct {
	Stat                          ProjectStat                   `json:"stat"`
	StrategiesWithReorgCount      []StrategyWithReorgCount      `json:"strategies_with_reorg_count"`
	StrategiesWithHonestLose      []StrategyWithHonestLose      `json:"strategies_with_honest_lose"`
	StrategiesWithGreatHonestLose []StrategyWithGreatHonestLose `json:"strategies_with_great_honest_lose"`
}

type ProjectStat struct {
	ProjectId     string `json:"project_id"`
	TotalSlot     int64  `json:"total_slot"`
	TotalStrategy int64  `json:"total_strategy"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
}

type StrategyWithReorgCount struct {
	ReorgCount      string `json:"reorg_count"`
	StrategyId      string `json:"strategy_id"`
	StrategyContent string `json:"strategy_content"`
}

// That HonestLoseRateAvg > 0
type StrategyWithHonestLose struct {
	HonestLoseRateAvg string `json:"honest_lose_rate_avg"`
	StrategyId        string `json:"strategy_id"`
	StrategyContent   string `json:"strategy_content"`
}

// That HonestLoseAvg > MaliciousLoseAvg
type StrategyWithGreatHonestLose struct {
	HonestLoseRateAvg    string `json:"honest_lose"`
	MaliciousLoseRateAvg string `json:"malicious_lose"`
	Ratio                string `json:"ratio"`
	StrategyId           string `json:"strategy_id"`
	StrategyContent      string `json:"strategy_content"`
}

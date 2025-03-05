package apimodels

type HomeProject struct {
	ProjectId     string `json:"project_id"`
	TotalSlot     int64  `json:"total_slot"`
	TotalStrategy int64  `json:"total_strategy"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	Category      string `json:"category"`
}

type HomeStrategy HomeStrategyWithGreatHonestLose

type HomeStrategyWithReorgCount struct {
	ProjectId       string `json:"project_id"`
	ReorgCount      string `json:"reorg_count"`
	StrategyId      string `json:"strategy_id"`
	StrategyContent string `json:"strategy_content"`
}

// That HonestLoseRateAvg > 0
type HomeStrategyWithHonestLose struct {
	ProjectId         string `json:"project_id"`
	HonestLoseRateAvg string `json:"honest_lose_rate_avg"`
	StrategyId        string `json:"strategy_id"`
	StrategyContent   string `json:"strategy_content"`
}

// That HonestLoseAvg > MaliciousLoseAvg
type HomeStrategyWithGreatHonestLose struct {
	ProjectId            string `json:"project_id"`
	HonestLoseRateAvg    string `json:"honest_lose"`
	MaliciousLoseRateAvg string `json:"malicious_lose"`
	Ratio                string `json:"ratio"`
	StrategyId           string `json:"strategy_id"`
	StrategyContent      string `json:"strategy_content"`
	Category             string `json:"category"`
}

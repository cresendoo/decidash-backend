package apiserver

// MarketSentiment 시장 심리 데이터
type MarketSentiment struct {
	LongPercentage  float64 `json:"long_percentage"`
	ShortPercentage float64 `json:"short_percentage"`
}

// TopPerformerMainPosition 최고 실적자 주요 포지션
type TopPerformerMainPosition struct {
	Asset string `json:"asset"`
	ROI   string `json:"roi"`
}

// AssetConcentration 자산 집중도
type AssetConcentration struct {
	HighestOI struct {
		Asset  string  `json:"asset"`
		Amount float64 `json:"amount"`
	} `json:"highest_oi"`
	MostTraded struct {
		Asset   string `json:"asset"`
		Traders int    `json:"traders"`
	} `json:"most_traded"`
	TotalMonitored float64 `json:"total_monitored"`
}

// TraderProfitability 트레이더 수익성
type TraderProfitability struct {
	ProfitablePercentage float64 `json:"profitable_percentage"`
	ProfitableCount      int     `json:"profitable_count"`
	TotalTraders         int     `json:"total_traders"`
	AvgDailyPnL          float64 `json:"avg_daily_pnl"`
}

// Trader 트레이더 정보
type Trader struct {
	Address       string        `json:"address"`
	Avatar        string        `json:"avatar"`
	IsStarred     bool          `json:"is_starred"`
	PerpEquity    float64       `json:"perp_equity"`
	MainPosition  *MainPosition `json:"main_position"`
	DirectionBias DirectionBias `json:"direction_bias"`
	DailyPnL      PnLData       `json:"daily_pnl"`
	WeeklyPnL     PnLData       `json:"weekly_pnl"`
	MonthlyPnL    PnLData       `json:"monthly_pnl"`
	AllTimePnL    PnLData       `json:"all_time_pnl"`
}

// MainPosition 주요 포지션
type MainPosition struct {
	Type   string  `json:"type"` // "LONG" or "SHORT"
	Asset  string  `json:"asset"`
	Amount float64 `json:"amount"`
}

// DirectionBias 방향 편향
type DirectionBias struct {
	LongPercentage  float64 `json:"long_percentage"`
	ShortPercentage float64 `json:"short_percentage"`
}

// PnLData 수익/손실 데이터
type PnLData struct {
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

// DashboardSummary 대시보드 요약 데이터
type DashboardSummary struct {
	MarketSentiment          MarketSentiment          `json:"market_sentiment"`
	TopPerformerMainPosition TopPerformerMainPosition `json:"top_performer_main_position"`
	AssetConcentration       AssetConcentration       `json:"asset_concentration"`
	TraderProfitability      TraderProfitability      `json:"trader_profitability"`
}

// TradersResponse 트레이더 목록 응답
type TradersResponse struct {
	Traders    []Trader `json:"traders"`
	Total      int      `json:"total"`
	Page       int      `json:"page"`
	PerPage    int      `json:"per_page"`
	TotalPages int      `json:"total_pages"`
}

// TradersRequest 트레이더 목록 요청
type TradersRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PerPage  int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search"`
	SortBy   string `form:"sort_by"`
	SortDesc bool   `form:"sort_desc"`
}

type FeePayerRequest struct {
	Signature   []byte `json:"signature"`
	Transaction []byte `json:"transaction"`
}

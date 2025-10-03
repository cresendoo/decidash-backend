package apiserver

import (
	"math/rand"
	"time"
)

// generateMockDashboardSummary 대시보드 요약 데이터 생성
func generateMockDashboardSummary() DashboardSummary {
	return DashboardSummary{
		MarketSentiment: MarketSentiment{
			LongPercentage:  62.0,
			ShortPercentage: 38.0,
		},
		TopPerformerMainPosition: TopPerformerMainPosition{
			Asset: "N/A",
			ROI:   "Infinity%",
		},
		AssetConcentration: AssetConcentration{
			HighestOI: struct {
				Asset  string  `json:"asset"`
				Amount float64 `json:"amount"`
			}{
				Asset:  "BTC",
				Amount: 906700000, // $906.70M
			},
			MostTraded: struct {
				Asset   string `json:"asset"`
				Traders int    `json:"traders"`
			}{
				Asset:   "ETH",
				Traders: 89,
			},
			TotalMonitored: 2250000000, // $2.25B
		},
		TraderProfitability: TraderProfitability{
			ProfitablePercentage: 27.4,
			ProfitableCount:      274,
			TotalTraders:         1000,
			AvgDailyPnL:          4800, // $4.8K
		},
	}
}

// generateMockTraders 트레이더 목록 생성
func generateMockTraders(count int) []Trader {
	rand.Seed(time.Now().UnixNano())

	assets := []string{"BTC", "ETH", "XRP", "SOL", "ADA", "DOT", "MATIC", "AVAX"}
	addresses := []string{
		"0x8af7...fa05", "0x1234...5678", "0xabcd...efgh", "0x9876...5432",
		"0x1111...2222", "0x3333...4444", "0x5555...6666", "0x7777...8888",
		"0x9999...aaaa", "0xbbbb...cccc", "0xdddd...eeee", "0xffff...0000",
	}

	traders := make([]Trader, count)

	for i := 0; i < count; i++ {
		// 랜덤한 포지션 타입과 자산
		positionType := "LONG"
		if rand.Float32() < 0.4 {
			positionType = "SHORT"
		}

		asset := assets[rand.Intn(len(assets))]

		// 랜덤한 수치들
		perpEquity := rand.Float64()*50000000 + 1000000     // $1M - $50M
		positionAmount := rand.Float64() * perpEquity * 0.9 // 포지션은 자산의 90% 이하

		// 방향 편향 (Long 비율)
		longPercentage := rand.Float64() * 100

		// PnL 데이터 생성
		dailyPnL := generateRandomPnL()
		weeklyPnL := generateRandomPnL()
		monthlyPnL := generateRandomPnL()
		allTimePnL := generateRandomPnL()

		traders[i] = Trader{
			Address:    addresses[i%len(addresses)],
			Avatar:     generateAvatarURL(i),
			IsStarred:  rand.Float32() < 0.1, // 10% 확률로 별표
			PerpEquity: perpEquity,
			MainPosition: &MainPosition{
				Type:   positionType,
				Asset:  asset,
				Amount: positionAmount,
			},
			DirectionBias: DirectionBias{
				LongPercentage:  longPercentage,
				ShortPercentage: 100 - longPercentage,
			},
			DailyPnL:   dailyPnL,
			WeeklyPnL:  weeklyPnL,
			MonthlyPnL: monthlyPnL,
			AllTimePnL: allTimePnL,
		}
	}

	return traders
}

// generateRandomPnL 랜덤한 PnL 데이터 생성
func generateRandomPnL() PnLData {
	// -50% ~ +500% 범위의 백분율
	percentage := (rand.Float64() - 0.5) * 11 // -5.5 ~ +5.5
	percentage *= 100                         // -550% ~ +550%

	// 매우 큰 수익률을 가끔 생성 (이미지의 극단적인 값들처럼)
	if rand.Float32() < 0.05 { // 5% 확률
		percentage = rand.Float64() * 1000000000 // 최대 10억%
	}

	// 금액은 백분율에 비례하되 랜덤 요소 추가
	baseAmount := rand.Float64() * 10000000 // 최대 $10M
	amount := baseAmount * (1 + percentage/100)

	return PnLData{
		Amount:     amount,
		Percentage: percentage,
	}
}

// generateAvatarURL 아바타 URL 생성
func generateAvatarURL(index int) string {
	// 실제로는 랜덤한 아바타 서비스 URL을 사용할 수 있음
	return "https://api.dicebear.com/7.x/avataaars/svg?seed=" + string(rune(65+index%26))
}

// getMockTradersWithPagination 페이지네이션을 적용한 트레이더 목록 반환
func getMockTradersWithPagination(page, perPage int, search string) TradersResponse {
	// 전체 트레이더 생성 (1000명)
	allTraders := generateMockTraders(1000)

	// 검색 필터 적용
	var filteredTraders []Trader
	if search != "" {
		for _, trader := range allTraders {
			if contains(trader.Address, search) {
				filteredTraders = append(filteredTraders, trader)
			}
		}
	} else {
		filteredTraders = allTraders
	}

	total := len(filteredTraders)
	totalPages := (total + perPage - 1) / perPage

	// 페이지네이션 적용
	start := (page - 1) * perPage
	end := start + perPage

	if start >= total {
		return TradersResponse{
			Traders:    []Trader{},
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}
	}

	if end > total {
		end = total
	}

	return TradersResponse{
		Traders:    filteredTraders[start:end],
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
}

// contains 문자열 포함 검사
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

package types

type Option[T any] struct {
	Vec []T `json:"vec"`
}

type OrderedMapEntry[K any, V any] struct {
	Key   K `json:"key"`
	Value V `json:"value"`
}

type OrderedMap[K any, V any] struct {
	Vec []OrderedMapEntry[K, V] `json:"vec"`
}

type Object struct {
	Inner string `json:"inner"`
}

type ExtendRef struct {
	Inner string `json:"inner"`
}

type I64 struct {
	IsPositive bool   `json:"is_positive"`
	Amount     Uint64 `json:"amount"`
}

type AccumulativeIndex struct {
	Index Uint128 `json:"index"`
}

type FeeWithDestination struct {
	Address string `json:"address"`
	Fees    Uint64 `json:"fees"`
}

type FeeDistribution struct {
	PositionAddress       string                     `json:"position_address"`
	PositionFeeDelta      I64                        `json:"position_fee_delta"`
	TreasuryFeeDelta      I64                        `json:"treasury_fee_delta"`
	BackstopVaultFees     Uint64                     `json:"backstop_vault_fees"`
	BuilderOrReferrerFees Option[FeeWithDestination] `json:"builder_or_referrer_fees"`
}

type OrderIDType struct {
	OrderID Uint128 `json:"order_id"`
}

type PriceIndexKey struct {
	TriggerPrice    Uint64         `json:"trigger_price"`
	PositionAddress string         `json:"position_address"`
	LimitPrice      Option[Uint64] `json:"limit_price"`
	IsFullSize      bool           `json:"is_full_size"`
}

type PendingTpSlKey struct {
	PriceIndex PriceIndexKey `json:"price_index"`
	OrderID    OrderIDType   `json:"order_id"`
}

type PendingTpSLs struct {
	FullSized  Option[PendingTpSlKey] `json:"full_sized"`
	FixedSized []PendingTpSlKey       `json:"fixed_sized"`
}

type FixedSizedTpSlForEvent struct {
	OrderID      Uint128        `json:"order_id"`
	TriggerPrice Uint64         `json:"trigger_price"`
	LimitPrice   Option[Uint64] `json:"limit_price"`
	Size         Uint64         `json:"size"`
}

type FullSizedTpSlForEvent struct {
	OrderID      Uint128        `json:"order_id"`
	TriggerPrice Uint64         `json:"trigger_price"`
	LimitPrice   Option[Uint64] `json:"limit_price"`
}

type PerpPosition struct {
	Size                                    Uint64            `json:"size"`
	EntryPxTimesSizeSum                     Uint128           `json:"entry_px_times_size_sum"`
	AvgAcquireEntryPx                       Uint64            `json:"avg_acquire_entry_px"`
	UserLeverage                            int               `json:"user_leverage"`
	MaxAllowedLeverage                      int               `json:"max_allowed_leverage"`
	IsLong                                  bool              `json:"is_long"`
	FundingIndexAtLastUpdate                AccumulativeIndex `json:"funding_index_at_last_update"`
	UnrealizedFundingAmountBeforeLastUpdate I64               `json:"unrealized_funding_amount_before_last_update"`
	Market                                  Object            `json:"market"`
	TpReqs                                  PendingTpSLs      `json:"tp_reqs"`
	SlReqs                                  PendingTpSLs      `json:"sl_reqs"`
	ReduceOnlyOrders                        []OrderIDType     `json:"reduce_only_orders"`
}

type CrossedPosition struct {
	Positions []PerpPosition `json:"positions"`
}

type IsolatedPosition struct {
	Position PerpPosition `json:"position"`
}

type IsolatedPositionRefs struct {
	ExtendRefs OrderedMap[Object, ExtendRef] `json:"extend_refs"`
}

type AccountInfo struct {
	FeeTrackingAddr string `json:"fee_tracking_addr"`
}

type AccountStatus struct {
	AccountBalance     I64    `json:"account_balance"`
	InitialMargin      Uint64 `json:"initial_margin"`
	TotalNotionalValue Uint64 `json:"total_notional_value"`
}

type AccountStatusDetailed struct {
	AccountBalance           I64    `json:"account_balance"`
	InitialMargin            Uint64 `json:"initial_margin"`
	LiquidationMargin        Uint64 `json:"liquidation_margin"`
	BackstopLiquidatorMargin Uint64 `json:"backstop_liquidator_margin"`
	TotalNotionalValue       Uint64 `json:"total_notional_value"`
}

type PositionPendingTpSL struct {
	OrderID      OrderIDType    `json:"order_id"`
	TriggerPrice Uint64         `json:"trigger_price"`
	Account      string         `json:"account"`
	LimitPrice   Option[Uint64] `json:"limit_price"`
	Size         Option[Uint64] `json:"size"`
}

type PositionUpdateEvent struct {
	Market                                            Object                        `json:"market"`
	User                                              string                        `json:"user"`
	IsLong                                            bool                          `json:"is_long"`
	Size                                              Uint64                        `json:"size"`
	UserLeverage                                      int                           `json:"user_leverage"`
	MaxAllowedLeverage                                int                           `json:"max_allowed_leverage"`
	EntryPriceTimesSizeSum                            Uint128                       `json:"entry_price_times_size_sum"`
	IsIsolated                                        bool                          `json:"is_isolated"`
	FundingIndexAtLastUpdate                          Uint128                       `json:"funding_index_at_last_update"`
	UnrealizedFundingAmountBeforeLastUpdate           Uint64                        `json:"unrealized_funding_amount_before_last_update"`
	IsUnrealizedFundingAmountBeforeLastUpdatePositive bool                          `json:"is_unrealized_funding_amount_before_last_update_positive"`
	FullSizedTp                                       Option[FullSizedTpSlForEvent] `json:"full_sized_tp"`
	FixedSizedTps                                     []FixedSizedTpSlForEvent      `json:"fixed_sized_tps"`
	FullSizedSl                                       Option[FullSizedTpSlForEvent] `json:"full_sized_sl"`
	FixedSizedSls                                     []FixedSizedTpSlForEvent      `json:"fixed_sized_sls"`
}

type Action string

const (
	ActionOpenLong   Action = "OpenLong"
	ActionCloseLong  Action = "CloseLong"
	ActionOpenShort  Action = "OpenShort"
	ActionCloseShort Action = "CloseShort"
	ActionNet        Action = "Net"
)

type ReduceOnlyValidationResult struct {
	Variant string                            `json:"variant"`
	Fields  *ReduceOnlyValidationResultFields `json:"fields,omitempty"`
}

type ReduceOnlyValidationResultFields struct {
	Size Uint64 `json:"size"`
}

type TradeEvent struct {
	Account               string `json:"account"`
	Market                Object `json:"market"`
	Action                Action `json:"action"`
	Size                  Uint64 `json:"size"`
	Price                 Uint64 `json:"price"`
	IsProfit              bool   `json:"is_profit"`
	RealizedPnlAmount     Uint64 `json:"realized_pnl_amount"`
	IsFundingPositive     bool   `json:"is_funding_positive"`
	RealizedFundingAmount Uint64 `json:"realized_funding_amount"`
	IsRebate              bool   `json:"is_rebate"`
	FeeAmount             Uint64 `json:"fee_amount"`
}

type UpdatePositionResult struct {
	Variant string                      `json:"variant"`
	Fields  *UpdatePositionResultFields `json:"fields,omitempty"`
}

type UpdatePositionResultFields struct {
	Account                       string            `json:"account"`
	Market                        Object            `json:"market"`
	IsIsolated                    bool              `json:"is_isolated"`
	PositionAddress               string            `json:"position_address"`
	MarginDelta                   Option[I64]       `json:"margin_delta"`
	BackstopLiquidatorCoveredLoss Uint64            `json:"backstop_liquidator_covered_loss"`
	FeeDistribution               FeeDistribution   `json:"fee_distribution"`
	RealizedPnl                   Option[I64]       `json:"realized_pnl"`
	RealizedFundingCost           Option[I64]       `json:"realized_funding_cost"`
	UnrealizedFundingCost         I64               `json:"unrealized_funding_cost"`
	UpdatedFundingIndex           AccumulativeIndex `json:"updated_funding_index"`
	VolumeDelta                   Uint128           `json:"volume_delta"`
	IsTaker                       bool              `json:"is_taker"`
}

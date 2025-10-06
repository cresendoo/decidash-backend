package models

import (
	"time"

	"github.com/cresendoo/decidash-backend/internal/application/decibel-indexer/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PerpPosition struct {
	PositionAddress                         string              `gorm:"primaryKey;column:address;type:varchar(66);not null"`
	Market                                  string              `gorm:"primaryKey;column:market;type:varchar(66);not null"`
	IsCrossed                               bool                `gorm:"primaryKey;column:is_crossed;type:bool;not null"`
	Version                                 uint64              `gorm:"column:version;type:numeric;not null"`
	VersionTimestamp                        time.Time           `gorm:"column:version_timestamp;type:timestamp;not null"`
	Owner                                   string              `gorm:"column:owner;type:varchar(66);not null"`
	Size                                    types.Uint64        `gorm:"column:size;type:decimal(20,0);not null"`
	EntryPxTimesSizeSum                     types.Uint128       `gorm:"column:entry_px_times_size_sum;type:decimal(39,0);not null"`
	AvgAcquireEntryPx                       types.Uint64        `gorm:"column:avg_acquire_entry_px;type:decimal(20,0);not null"`
	UserLeverage                            int                 `gorm:"column:user_leverage;type:int;not null"`
	MaxAllowedLeverage                      int                 `gorm:"column:max_allowed_leverage;type:int;not null"`
	IsLong                                  bool                `gorm:"column:is_long;type:bool;not null"`
	FundingIndexAtLastUpdate                types.Uint128       `gorm:"column:funding_index_at_last_update;type:decimal(39,0);not null"`
	UnrealizedFundingAmountBeforeLastUpdate types.I64           `gorm:"column:unrealized_funding_amount_before_last_update;type:json;serializer:json;not null"`
	ReduceOnlyOrders                        []types.OrderIDType `gorm:"column:reduce_only_orders;type:json;serializer:json;not null"`
	SlReqs                                  types.PendingTpSLs  `gorm:"column:sl_reqs;type:json;serializer:json;not null"`
	TpReqs                                  types.PendingTpSLs  `gorm:"column:tp_reqs;type:json;serializer:json;not null"`
}

func (s *PerpPosition) FromPerpPosition(
	positionAddress string,
	version uint64,
	versionTimestamp time.Time,
	owner string,
	isCrossed bool,
	value types.PerpPosition,
) {
	s.PositionAddress = positionAddress
	s.Version = version
	s.VersionTimestamp = versionTimestamp
	s.Owner = owner
	s.IsCrossed = isCrossed
	s.Market = value.Market.Inner
	s.Size = value.Size
	s.EntryPxTimesSizeSum = value.EntryPxTimesSizeSum
	s.AvgAcquireEntryPx = value.AvgAcquireEntryPx
	s.UserLeverage = value.UserLeverage
	s.MaxAllowedLeverage = value.MaxAllowedLeverage
	s.IsLong = value.IsLong
	s.FundingIndexAtLastUpdate = value.FundingIndexAtLastUpdate.Index
	s.UnrealizedFundingAmountBeforeLastUpdate = value.UnrealizedFundingAmountBeforeLastUpdate
	s.ReduceOnlyOrders = value.ReduceOnlyOrders
	s.SlReqs = value.SlReqs
	s.TpReqs = value.TpReqs
}

func (s *PerpPosition) TableName() string {
	return "PERP_POSITIONS"
}

func UpsertPosition(conn *gorm.DB, position PerpPosition) error {
	return conn.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "address"}, {Name: "market"}, {Name: "is_crossed"}},
			Where: clause.Where{
				Exprs: []clause.Expression{
					clause.Lt{
						Column: clause.Column{Name: "PERP_POSITIONS.version"},
						Value:  position.Version,
					},
				},
			},
			DoUpdates: clause.Assignments(position.toAssignments()),
		},
	).Create(&position).Error
}

func UpsertPositions(conn *gorm.DB, positions []PerpPosition) error {
	return conn.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "address"}, {Name: "market"}, {Name: "is_crossed"}},
			Where: clause.Where{
				Exprs: []clause.Expression{
					clause.Expr{
						SQL: "EXCLUDED.version > \"PERP_POSITIONS\".version",
					},
				},
			},
			DoUpdates: clause.AssignmentColumns(positions[0].toAssignmentColumns()),
		},
	).Create(&positions).Error
}

func (p PerpPosition) toAssignmentColumns() []string {
	return []string{
		"owner",
		"market",
		"size",
		"entry_px_times_size_sum",
		"avg_acquire_entry_px",
		"user_leverage",
		"max_allowed_leverage",
		"is_long",
		"funding_index_at_last_update",
		"unrealized_funding_amount_before_last_update",
		"reduce_only_orders",
		"sl_reqs",
		"tp_reqs",
	}
}

func (p PerpPosition) toAssignments() map[string]any {
	return map[string]any{
		"owner":                        p.Owner,
		"market":                       p.Market,
		"size":                         p.Size,
		"entry_px_times_size_sum":      p.EntryPxTimesSizeSum,
		"avg_acquire_entry_px":         p.AvgAcquireEntryPx,
		"user_leverage":                p.UserLeverage,
		"max_allowed_leverage":         p.MaxAllowedLeverage,
		"is_long":                      p.IsLong,
		"funding_index_at_last_update": p.FundingIndexAtLastUpdate,
		"unrealized_funding_amount_before_last_update": p.UnrealizedFundingAmountBeforeLastUpdate,
		"reduce_only_orders":                           p.ReduceOnlyOrders,
		"sl_reqs":                                      p.SlReqs,
		"tp_reqs":                                      p.TpReqs,
	}
}

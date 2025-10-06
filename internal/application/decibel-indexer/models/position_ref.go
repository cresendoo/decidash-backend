package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PerpPositionRef struct {
	Market          string `gorm:"primaryKey;column:market;type:varchar(66);not null"`
	PositionAddress string `gorm:"primaryKey;column:position_address;type:varchar(66);not null"`
	Version         uint64 `gorm:"column:version;type:numeric;not null"`
	Owner           string `gorm:"column:owner;type:varchar(66);not null"`
}

func (s *PerpPositionRef) TableName() string {
	return "PERP_POSITION_REF"
}

func UpsertPositionRef(conn *gorm.DB, positionRef PerpPositionRef) error {
	return conn.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "market"}, {Name: "position_address"}},
			Where: clause.Where{
				Exprs: []clause.Expression{
					clause.Lt{
						Column: clause.Column{Name: "PERP_POSITION_REF.version"},
						Value:  positionRef.Version,
					},
				},
			},
			DoUpdates: clause.Assignments(positionRef.toAssignments()),
		},
	).Create(&positionRef).Error
}

func (s *PerpPositionRef) toAssignments() map[string]any {
	return map[string]any{
		"version": s.Version,
		"owner":   s.Owner,
		"market":  s.Market,
	}
}

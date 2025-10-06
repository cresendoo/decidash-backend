package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IndexerState struct {
	ProcessorName          string    `gorm:"primaryKey;column:processor_name;type:varchar(512);not null"`
	LastProcessedVersion   uint64    `gorm:"column:last_processed_version;type:numeric;not null"`
	LastProcessedTimestamp time.Time `gorm:"column:last_processed_timestamp;type:timestamp;not null"`
	CreatedAt              time.Time `gorm:"autoCreateTime"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime"`
}

func (s *IndexerState) TableName() string {
	return "INDEXER_STATE"
}

func GetIndexerState(conn *gorm.DB, processorName string) (IndexerState, error) {
	var state IndexerState
	if err := conn.Where("processor_name = ?", processorName).First(&state).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return IndexerState{}, nil
		}
		return IndexerState{}, err
	}
	return state, nil
}

func UpsertIndexerState(conn *gorm.DB, state IndexerState) error {
	return conn.
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "processor_name"}},
				Where: clause.Where{
					Exprs: []clause.Expression{
						clause.Lte{
							Column: clause.Column{Name: "INDEXER_STATE.last_processed_version"},
							Value:  state.LastProcessedVersion,
						},
					},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"last_processed_version",
					"last_processed_timestamp",
					"updated_at",
				}),
			},
		).
		Create(&state).
		Error
}

package decibelindexer

import (
	"errors"
	"log/slog"
	"sort"
	"strconv"
	"time"

	"github.com/aptos-labs/aptos-go-sdk/api"
	"github.com/cresendoo/decidash-backend/internal/application/decibel-indexer/models"
	"github.com/cresendoo/decidash-backend/internal/application/decibel-indexer/types"
)

const (
	decibelContract      = "0xb8a5788314451ce4d2fbbad32e1bad88d4184b73943b7fe5166eab93cf1a5a95"
	crossedPosition      = decibelContract + "::perp_positions::CrossedPosition"
	isolatedPosition     = decibelContract + "::perp_positions::IsolatedPosition"
	isolatedPositionRefs = decibelContract + "::perp_positions::IsolatedPositionRefs"
	objectCore           = "0x1::object::ObjectCore"
)

func (a *Application) Process(txs []*api.UserTransaction) error {
	if len(txs) == 0 {
		return nil
	}

	if err := a.ProcessPositions(txs); err != nil {
		return err
	}
	stx := txs[0]
	etx := txs[len(txs)-1]
	if err := models.UpsertIndexerState(a.db, models.IndexerState{
		ProcessorName:          "decibel-indexer",
		LastProcessedVersion:   etx.Version,
		LastProcessedTimestamp: time.UnixMicro(int64(etx.Timestamp)),
	}); err != nil {
		return err
	}
	slog.Info("processed transactions", "start", stx.Version, "end", etx.Version, "count", len(txs))
	return nil
}

func (a *Application) ProcessPositions(txs []*api.UserTransaction) error {
	for _, tx := range txs {
		_, writeResources, _, _ := types.ExtractWriteSetChange(tx)
		var exist bool
		var objectCoreIndexs []int
		crossedPositions := make(map[string]types.CrossedPosition)
		isolatedPositions := make(map[string]types.IsolatedPosition)

		for idx, writeResource := range writeResources {
			switch writeResource.Data.Type {
			case crossedPosition:
				var crossedPosition types.CrossedPosition
				if err := MapToStructJSON(writeResource.Data.Data, &crossedPosition); err != nil {
					return err
				}
				crossedPositions[writeResource.Address.StringLong()] = crossedPosition
				exist = true
			case isolatedPosition:
				var isolatedPosition types.IsolatedPosition
				if err := MapToStructJSON(writeResource.Data.Data, &isolatedPosition); err != nil {
					return err
				}
				isolatedPositions[writeResource.Address.StringLong()] = isolatedPosition
				exist = true
			case objectCore:
				objectCoreIndexs = append(objectCoreIndexs, idx)
			}
		}
		if !exist {
			continue
		}

		addressOwnerMapping := make(map[string]string, len(objectCoreIndexs))
		for _, objectCoreIndex := range objectCoreIndexs {
			wr := writeResources[objectCoreIndex]
			var objectCore types.ObjectCore
			if err := MapToStructJSON(wr.Data.Data, &objectCore); err != nil {
				return err
			}
			addressOwnerMapping[wr.Address.StringLong()] = objectCore.Owner
		}

		positions := make(map[string]map[string]models.PerpPosition)

		for positionAddress, crossedPosition := range crossedPositions {
			for _, position := range crossedPosition.Positions {
				var row models.PerpPosition
				row.FromPerpPosition(positionAddress, tx.Version, time.UnixMicro(int64(tx.Timestamp)), positionAddress, true, position)
				if _, ok := positions[positionAddress]; !ok {
					positions[positionAddress] = make(map[string]models.PerpPosition)
				}
				if old, ok := positions[positionAddress][position.Market.Inner]; ok {
					if old.Version < tx.Version {
						positions[positionAddress][position.Market.Inner] = row
					}
				} else {
					positions[positionAddress][position.Market.Inner] = row
				}
			}
		}

		for positionAddress, isolatedPosition := range isolatedPositions {
			owner, ok := addressOwnerMapping[positionAddress]
			if !ok {
				return errors.New("owner not found, " + positionAddress + ", " + strconv.FormatUint(tx.Version, 10))
			}
			var row models.PerpPosition
			row.FromPerpPosition(positionAddress, tx.Version, time.UnixMicro(int64(tx.Timestamp)), owner, false, isolatedPosition.Position)
			if _, ok := positions[positionAddress]; !ok {
				positions[positionAddress] = make(map[string]models.PerpPosition)
			}
			if old, ok := positions[positionAddress][isolatedPosition.Position.Market.Inner]; ok {
				if old.Version < tx.Version {
					positions[positionAddress][isolatedPosition.Position.Market.Inner] = row
				}
			} else {
				positions[positionAddress][isolatedPosition.Position.Market.Inner] = row
			}
		}
		if len(positions) == 0 {
			continue
		}

		positionArray := make([]models.PerpPosition, 0, len(positions))
		for _, arr := range positions {
			for _, p := range arr {
				positionArray = append(positionArray, p)
			}
		}
		sort.Slice(positionArray, func(i, j int) bool {
			if positionArray[i].PositionAddress == positionArray[j].PositionAddress {
				return positionArray[i].Market < positionArray[j].Market
			}
			return positionArray[i].PositionAddress < positionArray[j].PositionAddress
		})
		if err := models.UpsertPositions(a.db, positionArray); err != nil {
			return err
		}
	}
	return nil
}

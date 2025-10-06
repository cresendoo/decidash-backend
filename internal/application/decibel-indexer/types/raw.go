package types

import (
	"time"

	"github.com/aptos-labs/aptos-go-sdk/api"
)

type Event struct {
	Version    uint64
	Timestamp  time.Time
	EventIndex int
	*api.Event
}

type WriteTableItem struct {
	Version   uint64
	Timestamp time.Time
	*api.WriteSetChangeWriteTableItem
}

type WriteResource struct {
	Version   uint64
	Timestamp time.Time
	*api.WriteSetChangeWriteResource
}

type DeleteTableItem struct {
	Version   uint64
	Timestamp time.Time
	*api.WriteSetChangeDeleteTableItem
}

type DeleteResource struct {
	Version   uint64
	Timestamp time.Time
	*api.WriteSetChangeDeleteResource
}

func ExtractWriteSetChange(tx *api.UserTransaction) (
	[]*WriteTableItem,
	[]*WriteResource,
	[]*DeleteResource,
	[]*DeleteTableItem,
) {
	var (
		deleteTableItems []*DeleteTableItem
		writeTableItems  []*WriteTableItem
		deleteResources  []*DeleteResource
		writeResources   []*WriteResource
	)
	for _, writeSet := range tx.Changes {
		version := tx.Version
		timestamp := time.UnixMicro(int64(tx.Timestamp))
		switch writeSet.Type {
		case api.WriteSetChangeVariantDeleteTableItem:
			deleteTableItems = append(deleteTableItems,
				&DeleteTableItem{
					Version:                       version,
					Timestamp:                     timestamp,
					WriteSetChangeDeleteTableItem: writeSet.Inner.(*api.WriteSetChangeDeleteTableItem),
				},
			)
		case api.WriteSetChangeVariantWriteTableItem:
			writeTableItems = append(writeTableItems,
				&WriteTableItem{
					Version:                      version,
					Timestamp:                    timestamp,
					WriteSetChangeWriteTableItem: writeSet.Inner.(*api.WriteSetChangeWriteTableItem),
				},
			)
		case api.WriteSetChangeVariantDeleteResource:
			deleteResources = append(deleteResources,
				&DeleteResource{
					Version:                      version,
					Timestamp:                    timestamp,
					WriteSetChangeDeleteResource: writeSet.Inner.(*api.WriteSetChangeDeleteResource),
				},
			)
		case api.WriteSetChangeVariantWriteResource:
			writeResources = append(writeResources,
				&WriteResource{
					Version:                     version,
					Timestamp:                   timestamp,
					WriteSetChangeWriteResource: writeSet.Inner.(*api.WriteSetChangeWriteResource),
				},
			)
		}
	}
	return writeTableItems, writeResources, deleteResources, deleteTableItems
}

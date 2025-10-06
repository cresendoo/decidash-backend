package types

type ObjectCore struct {
	AllowUngatedTransfer bool   `json:"allow_ungated_transfer"`
	GuidCreationNum      Uint64 `json:"guid_creation_num"`
	Owner                string `json:"owner"`
	TransferEvents       struct {
		Counter Uint64 `json:"counter"`
		Guid    struct {
			Id struct {
				Addr        string `json:"addr"`
				CreationNum Uint64 `json:"creation_num"`
			} `json:"id"`
		} `json:"guid"`
	} `json:"transfer_events"`
}

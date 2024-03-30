package http

type GetMintsResponse struct {
	TotalCount int    `json:"total_count"`
	Mints      []Mint `json:"minted"`
}

type Mint struct {
	TxHash       string `json:"tx_hash"`
	BlockNumber  int    `json:"block_number"`
	Owner        string `json:"owner"`
	Amount       int    `json:"amount"`
	BlobGasPrice uint   `json:"blob_gas_price"`
	GasFee       uint   `json:"gas_fee"`
	Time         int    `json:"time"`
}

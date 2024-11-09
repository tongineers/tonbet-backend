package models

type (
	AccountState struct {
		Status   string `json:"status"`
		Balance  int64  `json:"balance"`
		Data     string `json:"data"`
		LastHash []byte `json:"lastHash"`
		LastLt   uint64 `json:"lastLt"`
	}
)

package models

type (
	JettonData struct {
		TotalSupply int64  `json:"totalSupply"`
		UserAddress string `json:"userAddress"`
		Balance     int64  `json:"balance"`
	}
)

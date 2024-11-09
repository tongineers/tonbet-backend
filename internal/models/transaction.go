package models

type (
	Transaction struct {
		Hash     string `json:"hash"`
		LT       uint64 `json:"lt"`
		Fee      int64  `json:"fee"`
		DestAddr string `json:"dest_addr"`
		Comment  string `json:"comment"`
	}

	// Message struct {
	// 	BodyHash    string `json:"body_hash"`
	// 	CreatedLt   int64  `json:"created_lt"`
	// 	Destination string `json:"destination"`
	// 	FwdFee      int64  `json:"fwd_fee"`
	// 	IhrFee      int64  `json:"ihr_fee"`
	// 	Message     string `json:"message"`
	// 	Source      string `json:"source"`
	// 	Value       int64  `json:"value"`
	// }

	// TransactionID struct {
	// 	Hash string `json:"hash"`
	// 	Lt   int64  `json:"lt"`
	// }
)

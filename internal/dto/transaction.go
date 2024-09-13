package dto

type (
	Transaction struct {
		Data       string         `json:"data"`
		Fee        int64          `json:"fee"`
		InMsg      *Message       `json:"in_msg"`
		OtherFee   int64          `json:"other_fee"`
		OutMsgs    []*Message     `json:"out_msgs"`
		StorageFee int64          `json:"storage_fee"`
		TxnID      *TransactionID `json:"txn_id"`
	}

	Message struct {
		BodyHash    string `json:"body_hash"`
		CreatedLt   int64  `json:"created_lt"`
		Destination string `json:"destination"`
		FwdFee      int64  `json:"fwd_fee"`
		IhrFee      int64  `json:"ihr_fee"`
		Message     string `json:"message"`
		Source      string `json:"source"`
		Value       int64  `json:"value"`
	}

	TransactionID struct {
		Hash string `json:"hash"`
		Lt   int64  `json:"lt"`
	}
)

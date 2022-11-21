package types

type BlockData struct {
	Result Result `json:"result"`
	Error  Error  `json:"error"`
}

type Result struct {
	Number       string        `json:"number"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

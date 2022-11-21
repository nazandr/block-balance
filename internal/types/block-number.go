package types

type BlockNumber struct {
	Result string `json:"result"`
	Error  Error  `json:"error"`
}
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

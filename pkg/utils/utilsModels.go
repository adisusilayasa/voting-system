package utils

type Result struct {
	Response Response    `json:"response"`
	Data     interface{} `json:"data"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

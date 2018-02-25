package library

type Res struct{
	Code int          `json:"code"`
	Data interface{}  `json:"data"`
	Error string      `json:"error"`
	Message string    `json:"message"`
}
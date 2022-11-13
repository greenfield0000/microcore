package microcore

// swagger:model BaseResponse BaseResponse
type BaseResponse struct {
	ErrorMessage string      `json:"errorMessage"`
	Result       interface{} `json:"result"`
}

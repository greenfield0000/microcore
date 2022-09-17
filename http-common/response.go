package http_common

type BaseResponseWrapper struct {
	Response BaseResponse `json:"response"`
}

type BaseResponse struct {
	Data BaseData `json:"data"`
}

type BaseData struct {
	ErrorMessage string      `json:"errorMessage"`
	Result       interface{} `json:"result"`
}

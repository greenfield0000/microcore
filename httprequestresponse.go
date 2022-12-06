package main

// swagger:model BaseResponse BaseResponse
type BaseResponse struct {
	ErrorMessage string      `json:"errorMessage"`
	Result       interface{} `json:"result"`
}

type HttpRequestResponse interface {
	CreateResponseResult(result interface{}) BaseResponse
	CreateErrorMessage(errorMessage string) BaseResponse
}

type CommonHttpRequestResponse struct{}

func NewHttpRequestResponse() HttpRequestResponse {
	return CommonHttpRequestResponse{}
}

func (c CommonHttpRequestResponse) CreateResponseResult(result interface{}) BaseResponse {
	return BaseResponse{
		ErrorMessage: "",
		Result:       result,
	}
}

func (c CommonHttpRequestResponse) CreateErrorMessage(errorMessage string) BaseResponse {
	return BaseResponse{
		ErrorMessage: errorMessage,
		Result:       nil,
	}
}

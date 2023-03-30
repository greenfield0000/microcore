package common

type BaseResponse struct {
	ErrorMessage string      `json:"errorMessage"`
	Result       interface{} `json:"result"`
}

func CreateSuccessResult(result interface{}) BaseResponse {
	return BaseResponse{
		ErrorMessage: "",
		Result:       result,
	}
}

func CreateErrorMessage(errorMessage string) BaseResponse {
	return BaseResponse{
		ErrorMessage: errorMessage,
		Result:       nil,
	}
}

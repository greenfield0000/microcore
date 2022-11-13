package microcore

func CreateResponseResult(result interface{}) BaseResponse {
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

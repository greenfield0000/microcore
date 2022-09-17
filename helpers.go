package microcore

func CreateResponseResult(result interface{}) BaseResponse {
	return BaseResponse{
		Data: BaseData{Result: result},
	}
}

func CreateErrorMessage(errorMessage string) BaseResponseWrapper {
	return BaseResponseWrapper{
		Response: BaseResponse{
			Data: BaseData{
				ErrorMessage: errorMessage,
			},
		},
	}
}

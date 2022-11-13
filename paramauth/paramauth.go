package paramauth

import "github.com/greenfield0000/microcore"

// swagger:parameters LoginReq
type LoginReq struct {
	// in: body
	// required: true
	LoginParam AccountLoginParam
}

// AccountLoginParam
//
// Сущность для входа в приложение
//
// swagger:model AccountLoginParam
type AccountLoginParam struct {
	// Email адрес пользователя
	//
	// required: true
	// example: user@gmail.com
	Email string `json:"email"`
	// Пароль пользователя
	//
	// required: true
	// example: superPasswordInTheWorld@1123!!
	Password string `json:"password"`
}

// swagger:model AccountLoginResponse
type AccountLoginResponse struct {
	microcore.BaseResponse
	Result AccountLoginResponseData `json:"result"`
}

// swagger:model AccountLoginResponseData
type AccountLoginResponseData struct {
	Token string `json:"token"`
}

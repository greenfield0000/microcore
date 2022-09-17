package middleware_auth

// swagger:parameters LoginReq
type LoginReq struct {
	// in: body
	// required: true
	LoginParam LoginParam
}

// Сущность для входа в приложение
// swagger:model LoginParam
type LoginParam struct {
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

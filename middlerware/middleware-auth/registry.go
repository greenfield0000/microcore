package middleware_auth

// swagger:parameters RegistryReq
type RegistryReq struct {
	// in: body
	// required: true
	RegistryParam RegistryParam
}

// Сущность для входа в приложение
// swagger:parameters LoginParam
type RegistryParam struct {
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
	// Мобтльный телефон пользователя
	//
	// required: true
	// example: +79999999999
	Phone string `json:"phone"`
}

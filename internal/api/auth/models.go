package auth

// Константы путей API
const (
	LoginPath        = "/api/v1/auth/login"
	GoogleInitPath   = "/api/v1/auth/google/device/init"
	GoogleCheckPath  = "/api/v1/auth/google/device/check"
	YandexInitPath   = "/api/v1/auth/yandex/device/init"
	YandexCheckPath  = "/api/v1/auth/yandex/device/check"
	RefreshTokenPath = "/api/v1/auth/refresh"
	ProfilePath      = "/api/v1/users/me"
	LogoutPath       = "/api/v1/auth/logout"
	UpdateUserPath   = "/api/v1/users/me"
	DeleteUserPath   = "/api/v1/users/me"
	RegistrationPath = "/api/v1/auth/registration/"
)

// UserRegRequest - запрос на регистрацию пользователя
type UserRegRequest struct {
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest - запрос на вход с логином и паролем
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenResponse - ответ с токенами доступа и обновления
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ProfileResponse - ответ с профилем пользователя
type ProfileResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Currency struct {
		Gold      int `json:"gold"`
		GuildRage int `json:"guild_rage"`
	} `json:"currencies"`
}

// ErrorResponse - ответ об ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}

// DeviceAuthResponse - ответ на инициализацию OAuth
type DeviceAuthResponse struct {
	UserCode        string `json:"user_code"`          // код для пользователя
	DeviceCode      string `json:"device_code"`        // код устройства
	VerificationURL string `json:"verification_url"`   // предпочтительный URL для верификации
	VerificationURI string `json:"verification_uri"`   // альтернативное название (для совместимости)
	ExpiresIn       int    `json:"expires_in"`         // время жизни кода (сек)
	Interval        int    `json:"interval,omitempty"` // интервал опроса (сек)
}

// DeviceCheckResponse - ответ на проверку статуса
type DeviceCheckResponse struct {
	Status        string           `json:"status"`                   // Статус: "pending", "authenticated", "expired", "denied"
	TokenResponse *TokenResponse   `json:"token_response,omitempty"` // Токены, если авторизация успешна
	User          *ProfileResponse `json:"user,omitempty"`           // Данные пользователя, если есть
	Error         string           `json:"error,omitempty"`          // Описание ошибки, если есть
}

// UpdateUserRequest - запрос на изменение данных пользователя
// Поля не обязательные, обновляются только измененные
type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// SuccessResponse - успешный ответ
type SuccessResponse struct {
	Message string `json:"message"`
}

package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client - клиент для взаимодействия с API
type Client struct {
	baseURL      *url.URL
	httpClient   *http.Client
	accessToken  string
	refreshToken string
	userID       int
}

// NewClient - создание клиента для работы с API
func NewClient(baseURL string) (*Client, error) {

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("некорректный базовый URL: %w", err)
	}

	return &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}, nil
}

// SetTokens - установка access и refresh токенов в клиенте
func (c *Client) SetTokens(accessToken, refreshToken string) {
	c.accessToken = accessToken
	c.refreshToken = refreshToken
}

// doRequest HTTP запрос с заданным методом, путем и телом
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	reqURL := c.baseURL.ResolveReference(&url.URL{Path: path})

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("ошибка кодирования тела запроса: %w", err)
		}
	}

	// создание HTTP запроса с контекстом
	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), &buf)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// установка заголовков
	req.Header.Set("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	// выполнение запроса
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// обработка сетевых ошибок (недоступность сервера)
		return nil, fmt.Errorf("сетевая ошибка: %w", err)
	}
	defer resp.Body.Close()

	// чтение тела ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	// обработка HTTP ошибок (статус >= 400)
	if resp.StatusCode >= 400 {
		// Попытка распарсить как стандартную ошибку
		var serviceErr ErrorResponse
		if json.Unmarshal(responseBody, &serviceErr) == nil && serviceErr.Error != "" {
			return nil, fmt.Errorf("ошибка сервиса: %s", serviceErr.Error)
		}

		// Попытка распарсить как ошибку Gateway
		var gatewayErr struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(responseBody, &gatewayErr) == nil && gatewayErr.Message != "" {
			return nil, fmt.Errorf("ошибка шлюза: %s", gatewayErr.Message)
		}

		// обработка специфичных статусов
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, fmt.Errorf("не авторизован")
		case http.StatusServiceUnavailable: // 503
			return nil, fmt.Errorf("сервис временно недоступен")
		case http.StatusGatewayTimeout: // 504
			return nil, fmt.Errorf("таймаут шлюза")
		default:
			// уменьшение ответа для удобства
			errorBody := string(responseBody)
			if len(errorBody) > 200 {
				errorBody = errorBody[:200] + "..."
			}
			return nil, fmt.Errorf("HTTP ошибка %d: %s", resp.StatusCode, errorBody)
		}
	}

	return responseBody, nil
}

// Register - регистрация нового пользователя
func (c *Client) Register(ctx context.Context, req UserRegRequest) (*TokenResponse, *ProfileResponse, error) {
	body, err := c.doRequest(ctx, "POST", RegistrationPath, req)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка регистрации: %w", err)
	}

	var resp TokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	// установка токенов в клиент
	c.SetTokens(resp.AccessToken, resp.RefreshToken)

	profile, err := c.GetProfile(ctx)
	if err == nil {
		c.userID = profile.ID
	}

	return &resp, profile, nil
}

// Login - вход по логину и паролю
func (c *Client) Login(ctx context.Context, req LoginRequest) (*TokenResponse, *ProfileResponse, error) {
	body, err := c.doRequest(ctx, "POST", LoginPath, req)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка входа: %w", err)
	}

	var resp TokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	// установка токенов в клиент
	c.SetTokens(resp.AccessToken, resp.RefreshToken)

	profile, err := c.GetProfile(ctx)
	if err == nil {
		c.userID = profile.ID
	}

	return &resp, profile, nil
}

// RefreshToken - обновление access token с помощью refresh token
func (c *Client) RefreshToken(ctx context.Context) (*TokenResponse, error) {
	if c.refreshToken == "" {
		return nil, fmt.Errorf("отсутствует refresh token")
	}

	// использование refresh token как временный access token
	tempToken := c.accessToken
	c.accessToken = c.refreshToken

	body, err := c.doRequest(ctx, "POST", RefreshTokenPath, nil)

	// восстановление оригинальный access token
	c.accessToken = tempToken

	if err != nil {
		return nil, fmt.Errorf("ошибка обновления токена: %w", err)
	}

	var resp TokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	c.SetTokens(resp.AccessToken, c.refreshToken)
	return &resp, nil
}

// GetProfile - получение профиля текущего пользователя
func (c *Client) GetProfile(ctx context.Context) (*ProfileResponse, error) {
	path := fmt.Sprintf(GetProfilePath, c.userID)
	if c.userID == 0 {
		return nil, fmt.Errorf("user ID not set")
	}

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения профиля: %w", err)
	}

	var profile ProfileResponse
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("ошибка декодирования профиля: %w", err)
	}
	c.userID = profile.ID
	return &profile, nil
}

// InitOAuthDeviceFlow - инициация процесса OAuth
// provider - "google" или "yandex"
func (c *Client) InitOAuthDeviceFlow(ctx context.Context, provider string) (*DeviceAuthResponse, error) {
	var initPath string
	switch provider {
	case "google":
		initPath = GoogleInitPath
	case "yandex":
		initPath = YandexInitPath
	default:
		return nil, fmt.Errorf("неподдерживаемый провайдер: %s", provider)
	}

	body, err := c.doRequest(ctx, "POST", initPath, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации OAuth: %w", err)
	}

	var resp DeviceAuthResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа инициализации: %w", err)
	}

	// обработка разных вариантов названий полей
	if resp.VerificationURL == "" && resp.VerificationURI != "" {
		resp.VerificationURL = resp.VerificationURI
	}

	return &resp, nil
}

// CheckOAuthDeviceFlow - проверка статуса авторизации
func (c *Client) CheckOAuthDeviceFlow(ctx context.Context, provider, deviceCode string) (*DeviceCheckResponse, error) {
	var checkPath string
	switch provider {
	case "google":
		checkPath = GoogleCheckPath
	case "yandex":
		checkPath = YandexCheckPath
	default:
		return nil, fmt.Errorf("неподдерживаемый провайдер: %s", provider)
	}

	requestBody := struct {
		DeviceCode string `json:"device_code"`
	}{DeviceCode: deviceCode}

	body, err := c.doRequest(ctx, "POST", checkPath, requestBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки статуса OAuth: %w", err)
	}

	var resp DeviceCheckResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа проверки: %w", err)
	}

	// Автоматическое определение статуса, если не задан
	if resp.Status == "" {
		if resp.Error != "" {
			resp.Status = "error"
		} else if resp.TokenResponse != nil {
			resp.Status = "authenticated"
		} else {
			resp.Status = "pending"
		}
	}

	return &resp, nil
}

// CompleteOAuthPolling - полный цикл опроса для завершения авторизации через OAuth
// возвращает токены и профиль пользователя, если все гуд
func (c *Client) CompleteOAuthPolling(
	ctx context.Context,
	provider,
	deviceCode string,
	expiresIn,
	interval int,
) (*TokenResponse, *ProfileResponse, error) {
	if interval <= 0 {
		interval = 5
	}

	pollInterval := time.Duration(interval) * time.Second
	timeout := time.After(time.Duration(expiresIn) * time.Second)
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkResp, err := c.CheckOAuthDeviceFlow(ctx, provider, deviceCode)
			if err != nil {
				return nil, nil, err
			}

			switch checkResp.Status {
			case "authenticated":
				tokens := checkResp.TokenResponse
				if tokens == nil {
					return nil, nil, fmt.Errorf("токены отсутствуют в ответе")
				}

				// сохраняем токены в клиенте
				c.SetTokens(tokens.AccessToken, tokens.RefreshToken)

				var profile *ProfileResponse
				if checkResp.User != nil {
					// берем профиль из ответа, если он есть
					profile = checkResp.User
					c.userID = profile.ID
				} else {
					// если профиль не пришел, запрашиваем отдельно
					profile, err = c.GetProfile(ctx)
					if err != nil {
						return tokens, nil, fmt.Errorf("ошибка получения профиля: %w", err)
					}
					c.userID = profile.ID
				}

				return tokens, profile, nil

			case "expired":
				return nil, nil, fmt.Errorf("код устройства истек")
			case "denied":
				return nil, nil, fmt.Errorf("пользователь отклонил авторизацию")
			case "pending", "":
				// продолжаем опрос
			default:
				return nil, nil, fmt.Errorf("неожиданный статус: %s", checkResp.Status)
			}

		case <-timeout:
			return nil, nil, fmt.Errorf("время ожидания авторизации истекло (%d секунд)", expiresIn)
		case <-ctx.Done():
			return nil, nil, fmt.Errorf("авторизация отменена: %w", ctx.Err())
		}
	}
}

// Logout - выход из системы
func (c *Client) Logout(ctx context.Context) error {
	_, err := c.doRequest(ctx, "POST", LogoutPath, nil)
	if err != nil {
		return fmt.Errorf("ошибка выхода: %w", err)
	}

	c.clearSession()
	return nil
}

// UpdateUser - обновление данных пользователя (имя и/или пароль)
func (c *Client) UpdateUser(ctx context.Context, req UpdateUserRequest) (*ProfileResponse, error) {
	path := fmt.Sprintf(UpdateUserPath, c.userID)
	if c.userID == 0 {
		return nil, fmt.Errorf("user ID not set")
	}

	// проверка на наличие изменений
	if req.Username == "" && req.Password == "" {
		return nil, fmt.Errorf("не указаны данные для обновления")
	}

	body, err := c.doRequest(ctx, "PATCH", path, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления пользователя: %w", err)
	}

	var profile ProfileResponse
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("ошибка декодирования профиля: %w", err)
	}
	c.userID = profile.ID
	return &profile, nil
}

// DeleteUser - удаление текущего пользователя
func (c *Client) DeleteUser(ctx context.Context) error {
	path := fmt.Sprintf(DeleteUserPath, c.userID)
	if c.userID == 0 {
		return fmt.Errorf("user ID not set")
	}

	_, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя: %w", err)
	}

	c.clearSession()
	return nil
}

// clearSession очищение состояния клиента (токены, куки, userID)
func (c *Client) clearSession() {
	// сброс токенов и ID
	c.accessToken = ""
	c.refreshToken = ""
	c.userID = 0
}

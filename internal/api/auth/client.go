package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
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
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать cookie jar: %w", err)
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("некорректный базовый URL: %w", err)
	}

	return &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Jar:     jar,
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

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), &buf)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
			return nil, fmt.Errorf("ошибка API: %s", errResp.Error)
		}
		return nil, fmt.Errorf("HTTP ошибка %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// Login - вход по логину и паролю
func (c *Client) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	body, err := c.doRequest(ctx, "POST", LoginPath, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка входа: %w", err)
	}

	var resp TokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	c.SetTokens(resp.AccessToken, resp.RefreshToken)
	return &resp, nil
}

// RefreshToken - обновление access token с помощью refresh token
func (c *Client) RefreshToken(ctx context.Context) (*TokenResponse, error) {
	if c.refreshToken == "" {
		return nil, fmt.Errorf("отсутствует refresh token")
	}

	body, err := c.doRequest(ctx, "POST", RefreshTokenPath, map[string]string{
		"refresh_token": c.refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления токена: %w", err)
	}

	var resp TokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	c.SetTokens(resp.AccessToken, resp.RefreshToken)
	return &resp, nil
}

// GetProfile - получение профиля текущего пользователя
func (c *Client) GetProfile(ctx context.Context) (*ProfileResponse, error) {
	body, err := c.doRequest(ctx, "GET", ProfilePath, nil)
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
	// проверка на наличие изменений
	if req.Username == "" && req.Password == "" {
		return nil, fmt.Errorf("не указаны данные для обновления")
	}

	body, err := c.doRequest(ctx, "PATCH", UpdateUserPath, req)
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
	_, err := c.doRequest(ctx, "DELETE", DeleteUserPath, nil)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя: %w", err)
	}

	c.clearSession()
	return nil
}

// clearSession очищение состояния клиента (токены, куки, userID)
func (c *Client) clearSession() {
	// Сбрасываем токены и ID
	c.accessToken = ""
	c.refreshToken = ""
	c.userID = 0

	// Очищаем куки
	c.httpClient.Jar, _ = cookiejar.New(nil)
}

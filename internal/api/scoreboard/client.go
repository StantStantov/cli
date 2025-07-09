package scoreboard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lesta-start-battleship/cli/internal/api/token"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	basePath = "/"
)

// Client - клиент для работы с Scoreboard
type Client struct {
	baseURL      *url.URL
	httpClient   *http.Client
	accessToken  string
	refreshToken string
}

// NewClient - создание нового клиента
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

// SetAccessToken - установка токенов доступа для авторизации
func (c *Client) SetAccessToken(accessToken, refreshToken string) {
	c.accessToken = token.AccessToken
	c.refreshToken = token.RefreshToken
}

// GetUserStats - получение статистики пользователей
func (c *Client) GetUserStats(
	ctx context.Context,
	userID *int,
	nameFilter string,
	orderBy string,
	reverse bool,
	limit int,
	page int,
) (*UserListResponse, error) {
	endpoint := fmt.Sprintf("%s%s/users", c.baseURL, basePath)
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("ошибка формирования URL: %w", err)
	}

	// Подготовка параметров запроса
	q := u.Query()
	if userID != nil {
		q.Add("id_like", strconv.Itoa(*userID))
	}
	if nameFilter != "" {
		q.Add("name_ilike", nameFilter)
	}
	if orderBy != "" {
		q.Add("order_by", orderBy)
	}
	if reverse {
		q.Add("reverse", "true")
	}
	q.Add("limit", strconv.Itoa(limit))
	q.Add("page", strconv.Itoa(page))

	u.RawQuery = q.Encode()

	body, err := c.doRequest(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var response UserListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return &response, nil
}

// GetCurrentUserStats - получение статистики текущего пользователя
func (c *Client) GetCurrentUserStats(ctx context.Context, userID int) (*UserStat, error) {
	resp, err := c.GetUserStats(ctx, &userID, "", "", false, 1, 1)
	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("статистика пользователя с id %d не найдена", userID)
	}

	return &resp.Items[0], nil
}

// GetChestLeaders - получение лидеров по открытию сундуков
func (c *Client) GetChestLeaders(ctx context.Context, limit int) ([]UserStat, error) {
	resp, err := c.GetUserStats(
		ctx,
		nil,
		"",             // без фильтра по имени
		"chest_opened", // сортировка по сундукам
		true,           // по убыванию (от большего к меньшему)
		limit,          // запрошенное количество
		1,              // страница
	)

	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// GetGuildStats - получение статистики гильдий
func (c *Client) GetGuildStats(
	ctx context.Context,
	guildID *int,
	nameFilter string,
	orderBy string,
	reverse bool,
	limit int,
	page int,
) (*GuildListResponse, error) {
	endpoint := fmt.Sprintf("%s%s/guilds", c.baseURL, basePath)
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("ошибка формирования URL: %w", err)
	}

	q := u.Query()
	if guildID != nil {
		q.Add("id_like", strconv.Itoa(*guildID))
	}
	if nameFilter != "" {
		q.Add("name_ilike", nameFilter)
	}
	if orderBy != "" {
		q.Add("order_by", orderBy)
	}
	if reverse {
		q.Add("reverse", "true")
	}
	q.Add("limit", strconv.Itoa(limit))
	q.Add("page", strconv.Itoa(page))

	u.RawQuery = q.Encode()

	body, err := c.doRequest(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var response GuildListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return &response, nil
}

// doRequest - шаблон для GET запросов и получения ответов
func (c *Client) doRequest(ctx context.Context, url string) ([]byte, error) {
	// Создаем HTTP-запрос
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// установка токена в хедер
	req.Header.Set("Accept", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", c.accessToken)
		req.Header.Set("Refresh-Token", c.refreshToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()
	token.AccessToken = resp.Header.Get("Authorization")
	token.RefreshToken = resp.Header.Get("Refresh-Token")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	// обработка HTTP-ошибок
	if resp.StatusCode >= 400 {
		errorMsg := string(body)

		switch resp.StatusCode {
		case 400:
			return nil, fmt.Errorf("неверный запрос: %s", errorMsg)
		case 401:
			return nil, fmt.Errorf("требуется авторизация: %s", errorMsg)
		case 403:
			return nil, fmt.Errorf("доступ запрещен: %s", errorMsg)
		case 404:
			return nil, fmt.Errorf("ресурс не найден: %s", errorMsg)
		case 500:
			return nil, fmt.Errorf("внутренняя ошибка сервера: %s", errorMsg)
		default:
			return nil, fmt.Errorf("ошибка HTTP %d: %s", resp.StatusCode, errorMsg)
		}
	}

	return body, nil
}

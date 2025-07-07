package guilds

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

// Client - клиент для взаимодействия с API гильдий
type Client struct {
	baseURL     *url.URL
	httpClient  *http.Client
	accessToken string
}

// NewClient создает новый клиент для работы с API гильдий
func NewClient(baseURL string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("некорректный базовый URL: %w", err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать cookie jar: %w", err)
	}

	return &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: 15 * time.Second,
		},
	}, nil
}

// SetAccessToken устанавливает Access token
func (c *Client) SetAccessToken(token string) {
	c.accessToken = token
}

// doRequest выполняет HTTP запрос
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

// JoinGuild - отправка запроса на вступление в гильдию
func (c *Client) JoinGuild(ctx context.Context, guildID int) (*BaseResponse, error) {
	path := fmt.Sprintf(JoinGuildPath, guildID)
	body, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка вступления в гильдию: %w", err)
	}

	var resp BaseResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	return &resp, nil
}

// CreateGuild - создание новой гильдии
func (c *Client) CreateGuild(ctx context.Context, req CreateGuildRequest) (*BaseResponse, error) {
	body, err := c.doRequest(ctx, "POST", CreateGuildPath, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания гильдии: %w", err)
	}

	var resp BaseResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	return &resp, nil
}

// DeleteGuild - удаление гильдии
func (c *Client) DeleteGuild(ctx context.Context, guildID int) (*BaseResponse, error) {
	path := fmt.Sprintf(DeleteGuildPath, guildID)
	body, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка удаления гильдии: %w", err)
	}

	var resp BaseResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	return &resp, nil
}

// GetGuildMembers - получение списка участников гильдии
func (c *Client) GetGuildMembers(ctx context.Context, guildID int) (*BaseResponse, error) {
	path := fmt.Sprintf(GetGuildMembersPath, guildID)
	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка участников: %w", err)
	}

	var resp BaseResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	return &resp, nil
}

// UpdateMemberRole - изменение роли участника гильдии
func (c *Client) UpdateMemberRole(
	ctx context.Context,
	guildID int,
	memberID int,
	req UpdateRoleRequest,
) (*BaseResponse, error) {
	path := fmt.Sprintf(UpdateMemberRolePath, guildID, memberID)
	body, err := c.doRequest(ctx, "PUT", path, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка изменения роли участника: %w", err)
	}

	var resp BaseResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	return &resp, nil
}

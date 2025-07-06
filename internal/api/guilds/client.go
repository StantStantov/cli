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
	"strconv"
	"time"
)

// Client - клиент для работы с API гильдий
type Client struct {
	baseURL     *url.URL
	httpClient  *http.Client
	accessToken string
}

// NewClient создает новый клиент
func NewClient(baseURL string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	return &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: 15 * time.Second,
		},
	}, nil
}

// SetAccessToken устанавливает токен доступа
func (c *Client) SetAccessToken(token string) {
	c.accessToken = token
}

// doRequest HTTP запрос с заданным методом, путем и телом и с учётом query-параметров
func (c *Client) doRequest(
	ctx context.Context,
	method, path string,
	queryParams map[string]string,
	body interface{},
) ([]byte, error) {
	reqURL := c.baseURL.ResolveReference(&url.URL{Path: path})
	q := reqURL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	reqURL.RawQuery = q.Encode()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("error encoding request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), &buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SeaBattle-CLI/1.0")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// GetMemberByUserID - получить инфо об участнике по user_id
func (c *Client) GetMemberByUserID(ctx context.Context, userID int) (*MemberResponse, error) {
	path := fmt.Sprintf(PathGetMemberByUserID, userID)
	body, err := c.doRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp ResponseMember
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// GetGuildByTag - получить инфо о гильдии по тегу
func (c *Client) GetGuildByTag(ctx context.Context, tag string) (*GuildResponse, error) {
	path := fmt.Sprintf(PathGetGuildByTag, tag)
	body, err := c.doRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp ResponseGuild
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// SendJoinRequest - отправить запрос на вступление в гильдию
func (c *Client) SendJoinRequest(ctx context.Context, guildTag string, userID int) error {
	path := fmt.Sprintf(PathSendJoinRequest, guildTag)
	params := map[string]string{"user_id": strconv.Itoa(userID)}
	_, err := c.doRequest(ctx, "POST", path, params, nil)
	return err
}

// GetJoinRequests - получить список заявок на вступление (для owner/officer)
func (c *Client) GetJoinRequests(ctx context.Context, guildTag string, userID int) (*RequestPagination, error) {
	path := fmt.Sprintf(PathGetJoinRequests, guildTag)
	params := map[string]string{"user_id": strconv.Itoa(userID)}
	body, err := c.doRequest(ctx, "GET", path, params, nil)
	if err != nil {
		return nil, err
	}
	var resp ResponseRequestPagination
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// ApplyJoinRequest - принять заявку на вступление (owner/officer)
func (c *Client) ApplyJoinRequest(ctx context.Context, guildTag string, userID int, guildMemberID int) error {
	path := fmt.Sprintf(PathApplyJoinRequest, guildTag, userID)
	params := map[string]string{"guild_member_id": strconv.Itoa(guildMemberID)}
	_, err := c.doRequest(ctx, "POST", path, params, nil)
	return err
}

// CancelJoinRequest - отклонить заявку на вступление (owner/officer)
func (c *Client) CancelJoinRequest(ctx context.Context, guildTag string, userID int, guildMemberID int) error {
	path := fmt.Sprintf(PathCancelJoinRequest, guildTag, userID)
	params := map[string]string{"guild_member_id": strconv.Itoa(guildMemberID)}
	_, err := c.doRequest(ctx, "DELETE", path, params, nil)
	return err
}

// CreateGuild - создать новую гильдию
func (c *Client) CreateGuild(ctx context.Context, userID int, req CreateGuildRequest) (*GuildResponse, error) {
	path := PathCreateGuild
	params := map[string]string{"user_id": strconv.Itoa(userID)}
	body, err := c.doRequest(ctx, "POST", path, params, &req)
	if err != nil {
		return nil, err
	}
	var resp ResponseGuild
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// DeleteGuild - удалить свою гильдию (owner)
func (c *Client) DeleteGuild(ctx context.Context, tag string, userID int) error {
	path := fmt.Sprintf(PathDeleteGuild, tag)
	params := map[string]string{"user_id": strconv.Itoa(userID)}
	_, err := c.doRequest(ctx, "DELETE", path, params, nil)
	return err
}

// GetGuildMembers - получить список участников гильдии (с пагинацией)
func (c *Client) GetGuildMembers(ctx context.Context, tag string, offset, limit int) (*MemberPagination, error) {
	path := fmt.Sprintf(PathGetGuildMembers, tag)
	params := map[string]string{
		"offset": strconv.Itoa(offset),
		"limit":  strconv.Itoa(limit),
	}
	body, err := c.doRequest(ctx, "GET", path, params, nil)
	if err != nil {
		return nil, err
	}
	var resp ResponseMemberPagination
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// DeleteMember - удалить участника из гильдии (owner/officer)
func (c *Client) DeleteMember(ctx context.Context, tag string, userID, guildMemberID int) error {
	path := fmt.Sprintf(PathDeleteMember, tag, userID)
	params := map[string]string{"guild_member_id": strconv.Itoa(guildMemberID)}
	_, err := c.doRequest(ctx, "DELETE", path, params, nil)
	return err
}

// EditMember - изменить роль или имя участника (owner/officer)
func (c *Client) EditMember(ctx context.Context, tag string, userID, guildMemberID int, req EditMemberRequest) error {
	path := fmt.Sprintf(PathEditMember, tag, userID)
	params := map[string]string{"guild_member_id": strconv.Itoa(guildMemberID)}
	_, err := c.doRequest(ctx, "PATCH", path, params, &req)
	return err
}

// ExitGuild - выйти из гильдии (любой участник)
func (c *Client) ExitGuild(ctx context.Context, tag string, userID int) error {
	path := fmt.Sprintf(PathExitGuild, tag, userID)
	_, err := c.doRequest(ctx, "DELETE", path, nil, nil)
	return err
}

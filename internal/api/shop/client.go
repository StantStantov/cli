package shop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lesta-start-battleship/cli/internal/api/token"
	"net/http"
	"net/url"
	"time"
)

// Client - клиент для работы с Shop
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
		baseURL:    parsedURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}, nil
}

// SetAccessToken - установка токенов доступа для авторизации
func (c *Client) SetAccessToken(accessToken, refreshToken string) {
	c.accessToken = accessToken
	c.refreshToken = refreshToken
}

// doRequest - универсальный метод для выполнения запросов с обработкой тела
func (c *Client) doRequest(
	ctx context.Context,
	method, path string,
	body interface{},
) ([]byte, error) {
	// Формируем полный URL
	fullURL := c.baseURL.ResolveReference(&url.URL{Path: path}).String()

	// Подготавливаем тело запроса
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("ошибка кодирования тела запроса: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
		req.Header.Set("Refresh-Token", c.refreshToken)
	}

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()
	token.AccessToken = resp.Header.Get("Authorization")
	token.RefreshToken = resp.Header.Get("Refresh-Token")

	// Читаем тело ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	// Обрабатываем HTTP-ошибки
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// GetProducts - получение списка предметов
func (c *Client) GetProducts(ctx context.Context) ([]Product, error) {
	body, err := c.doRequest(ctx, "GET", "/item/", nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения продуктов: %w", err)
	}

	var products []Product
	if err := json.Unmarshal(body, &products); err != nil {
		return nil, fmt.Errorf("ошибка декодирования продуктов: %w", err)
	}

	return products, nil
}

// GetChests - получение списка сундуков
func (c *Client) GetChests(ctx context.Context) ([]Chest, error) {
	body, err := c.doRequest(ctx, "GET", "/chest/chest/", nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения сундуков: %w", err)
	}

	var chests []Chest
	if err := json.Unmarshal(body, &chests); err != nil {
		return nil, fmt.Errorf("ошибка декодирования сундуков: %w", err)
	}

	return chests, nil
}

// GetPromotions - получение списка акций
func (c *Client) GetPromotions(ctx context.Context) ([]Promotion, error) {
	body, err := c.doRequest(ctx, "GET", "/promotion/", nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения акций: %w", err)
	}

	var promotions []Promotion
	if err := json.Unmarshal(body, &promotions); err != nil {
		return nil, fmt.Errorf("ошибка декодирования акций: %w", err)
	}

	return promotions, nil
}

// BuyProduct - покупка предмета
func (c *Client) BuyProduct(ctx context.Context, itemID int) error {
	path := fmt.Sprintf("/item/%d/buy/", itemID)
	_, err := c.doRequest(ctx, "POST", path, nil)
	return err
}

// BuyChest - покупка сундука
func (c *Client) BuyChest(ctx context.Context, chestID int) error {
	path := fmt.Sprintf("/chest/chest/%d/buy/", chestID)
	_, err := c.doRequest(ctx, "POST", path, nil)
	return err
}

// OpenChest - открытие сундука
func (c *Client) OpenChest(ctx context.Context, chestID, amount int) error {
	requestBody := OpenChestRequest{
		ChestID: chestID,
		Amount:  amount,
	}

	_, err := c.doRequest(ctx, "POST", "/chest/chest/open/", requestBody)
	return err
}

// TODO:
// 	GetUserPurchases - получить историю покупок пользователя ( GET /purchase/ )
//  Вернет список всех покупок текущего пользователя в разделе "История покупок"

func (c *Client) GetUserPurchases(ctx context.Context) ([]Purchase, error) {
	// Заглушка для будущей реализации
	return nil, fmt.Errorf("not implemented yet")
}

// TODO:
//  GetPromotionDetails - получить детали акции ( GET /promotion/{id}/ )
//  Вернет полную информацию о конкретной акции на её странице

func (c *Client) GetPromotionDetails(ctx context.Context, promotionID int) (*Promotion, error) {
	// Заглушка для будущей реализации
	return nil, fmt.Errorf("not implemented yet")
}

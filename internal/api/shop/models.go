package shop

import "time"

// Product - игровой предмет для покупки
type Product struct {
	ID          int               `json:"item_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Currency    string            `json:"currency_type"`
	Cost        int               `json:"cost"`
	Promotion   *ProductPromotion `json:"promotion"`
	DailyLimit  *int              `json:"daily_purchase_limit"`
}

type ProductPromotion struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Chest - игровой сундук
type Chest struct {
	ID              int    `json:"item_id"`
	Name            string `json:"name"`
	Gold            int    `json:"gold"`
	PromotionID     *int   `json:"promotion"`
	ItemProbability int    `json:"item_probability"`
	Currency        string `json:"currency_type"`
	Cost            int    `json:"cost"`
	Experience      int    `json:"experience"`
	Products        string `json:"products"`
	SpecialProducts string `json:"special_products"`
}

// Promotion - активная акция
type Promotion struct {
	ID        int                `json:"id"`
	Name      string             `json:"name"`
	StartDate time.Time          `json:"start_date"`
	EndDate   time.Time          `json:"end_date"`
	Duration  time.Duration      `json:"duration"`
	IsActive  string             `json:"is_active"`
	Chests    []Chest            `json:"chests"`
	Product   []ProductPromotion `json:"product"`
}

// Purchase - информация о покупке
type Purchase struct {
	ID          int       `json:"id"`
	UserID      int       `json:"owner"`
	Quantity    int       `json:"quantity"`
	Date        time.Time `json:"date"`
	ItemID      *int      `json:"item"`
	ChestID     *int      `json:"chest"`
	PromotionID *int      `json:"promotion"`
}

// OpenChestRequest - запрос на открытие сундука
type OpenChestRequest struct {
	ChestID int `json:"item_id"`
	Amount  int `json:"amount"`
}

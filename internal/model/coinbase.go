package model

import (
	"time"
)

type CoinbaseSubscribeMessage struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type CoinbaseResponseMessage struct {
	Type         string    `json:"type"`
	TradeID      int       `json:"trade_id"`
	Sequence     int64     `json:"sequence"`
	MakerOrderID string    `json:"maker_order_id"`
	TakerOrderID string    `json:"taker_order_id"`
	Time         time.Time `json:"time"`
	ProductID    string    `json:"product_id"`
	Size         string    `json:"size"`
	Price        string    `json:"price"`
	Side         string    `json:"side"`
}

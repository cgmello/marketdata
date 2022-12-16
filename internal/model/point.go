package model

// New price and quantity for a trading pair
type Point struct {
	TradingPair string
	Price       float64
	Quantity    float64
}

func NewPoint(t string, p, q float64) Point {
	return Point{
		TradingPair: t,
		Price:       p,
		Quantity:    q,
	}
}

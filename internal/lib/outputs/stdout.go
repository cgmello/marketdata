package outputs

import (
	"log"

	"github.com/cgmello/marketdata/model"
)

type Stdout struct {
	Point model.Point
	Name  string
	Value float64
}

func (out *Stdout) Send() {
	log.Printf("%s Price:%9.2f Qty:%5.2f %s:%9.2f\n", out.Point.TradingPair, out.Point.Price, out.Point.Quantity, out.Name, out.Value)
}

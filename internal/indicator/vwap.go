package indicator

import (
	"errors"
	"sync"

	"github.com/cgmello/marketdata/model"
	"github.com/shopspring/decimal"
)

/*
	Calculate the VWAP per trading pair using a sliding window of 200 data
	points. Meaning, when a new data point arrives through the websocket feed
	the oldest data point will fall off and the new one will be added such that
	no more than 200 data points are included in the calculation. The first 200
	updates will have less than 200 data points included.
*/

// Maps values to each trading pair
// Handles concurrency
// VWAP equation = sum(price*quantity) / sum(quantity)
type VWAP struct {
	mu            sync.Mutex
	SlidingWindow int64                      // FIFO sliding window size
	Points        map[string][]model.Point   // list of points for each trading pair for the sliding window
	SumPriceQty   map[string]decimal.Decimal // last sum(price*quantity) for each trading pair
	SumQty        map[string]decimal.Decimal // last sum(quantity) for each trading pair
	VWAP          map[string]decimal.Decimal // last VWAP value for each trading pair
}

func NewVWAP(w int64) VWAP {
	return VWAP{
		SlidingWindow: w,
		Points:        make(map[string][]model.Point),
		SumPriceQty:   make(map[string]decimal.Decimal),
		SumQty:        make(map[string]decimal.Decimal),
		VWAP:          make(map[string]decimal.Decimal),
	}
}

// Get the current VWAP value for a trading pair
func (d *VWAP) Get(tradingPair string) (float64, error) {
	v, ok := d.VWAP[tradingPair]
	if !ok {
		return 0.00, errors.New("VWAP value undefined for trading pair")
	}

	f, _ := v.Float64()

	return f, nil
}

// Push a new point to the data series, returns the VWAP value
func (d *VWAP) Process(p model.Point) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	tradingPair := p.TradingPair

	// check if we have points for this trading pair
	if _, ok := d.Points[tradingPair]; !ok {
		// create a new empty series for this trading pair
		d.Points[tradingPair] = []model.Point{}
	} else {
		// check if we reached the limit of the sliding window size
		if int64(len(d.Points[tradingPair])) == d.SlidingWindow {
			// remove the first one (FIFO)
			d.remove(tradingPair)
		}
	}

	// add the new point and recalculate VWAP value
	d.add(tradingPair, p)

	// return VWAP value
	v := d.VWAP[tradingPair]
	f, _ := v.Float64()

	return f
}

// Removes the first point (FIFO)
func (d *VWAP) remove(tradingPair string) {
	firstPoint := d.Points[tradingPair][0]

	price := decimal.NewFromFloat(firstPoint.Price)
	qty := decimal.NewFromFloat(firstPoint.Quantity)
	m := price.Mul(qty)

	// remove first point values from computation
	d.SumPriceQty[tradingPair] = d.SumPriceQty[tradingPair].Sub(m)
	d.SumQty[tradingPair] = d.SumQty[tradingPair].Sub(qty)

	// VWAP equation = sum(price*quantity) / sum(quantity)
	if d.SumQty[tradingPair].GreaterThanOrEqual(decimal.NewFromInt(0)) {
		d.VWAP[tradingPair] = d.SumPriceQty[tradingPair].Div(d.SumQty[tradingPair])
	}

	// remove from list of points
	d.Points[tradingPair] = d.Points[tradingPair][1:]
}

func (d *VWAP) add(tradingPair string, p model.Point) {
	price := decimal.NewFromFloat(p.Price)
	qty := decimal.NewFromFloat(p.Quantity)
	m := price.Mul(qty)

	// add new point values to computation
	d.SumPriceQty[tradingPair] = d.SumPriceQty[tradingPair].Add(m)
	d.SumQty[tradingPair] = d.SumQty[tradingPair].Add(qty)

	// VWAP equation = sum(price*volume) / sum(volume)
	if d.SumQty[tradingPair].GreaterThanOrEqual(decimal.NewFromInt(0)) {
		d.VWAP[tradingPair] = d.SumPriceQty[tradingPair].Div(d.SumQty[tradingPair])
	}

	// add new point
	points := d.Points[tradingPair]
	d.Points[tradingPair] = append(points, p)
}

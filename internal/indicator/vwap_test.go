package indicator

import (
	"github.com/cgmello/marketdata/model"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

var points = []model.Point{
	model.NewPoint("XXX-YYY", 100, 10),
	model.NewPoint("XXX-YYY", 150, 15),
	model.NewPoint("XXX-YYY", 200, 10),
}

func TestVWAP(t *testing.T) {
	slidingWindow := int64(200)

	vwap := NewVWAP(slidingWindow)

	// VWAP equation = sum(price*quantity) / sum(quantity)

	// Point0 -> WVAP = 100
	v := vwap.Process(points[0])
	require.Equal(t, 1, int(len(vwap.Points["XXX-YYY"])))  // size
	require.Equal(t, points[0], vwap.Points["XXX-YYY"][0]) // check data series
	require.Equal(t, float64(100), v)

	// Point1 -> WVAP = (100*10 + 150*15) / (10 + 15) = 130
	v = vwap.Process(points[1])
	require.Equal(t, 2, int(len(vwap.Points["XXX-YYY"])))  // size
	require.Equal(t, points[1], vwap.Points["XXX-YYY"][1]) // check data series
	require.Equal(t, float64(130), v)

	// Point2 -> WVAP = 150
	v = vwap.Process(points[2])
	require.Equal(t, 3, int(len(vwap.Points["XXX-YYY"])))  // size
	require.Equal(t, points[2], vwap.Points["XXX-YYY"][2]) // check data series
	require.Equal(t, float64(150), v)
}

func TestConcurrency(t *testing.T) {
	slidingWindow := int64(200)

	vwap := NewVWAP(slidingWindow)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		vwap.Process(points[0])
		wg.Done()
	}()

	go func() {
		vwap.Process(points[1])
		wg.Done()
	}()

	go func() {
		vwap.Process(points[2])
		wg.Done()
	}()

	wg.Wait()

	require.Equal(t, 3, int(len(vwap.Points["XXX-YYY"])))

	v, err := vwap.Get("XXX-YYY")
	require.NoError(t, err)
	require.Equal(t, float64(150), v)
}

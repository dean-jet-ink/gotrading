package algo

import (
	"math"

	"github.com/markcheno/go-talib"
)

func HV(closes []float64, period int) []float64 {
	var changes []float64

	for i := range closes {
		if i == 0 {
			continue
		}

		dayChange := math.Log(closes[i] / closes[i-1])
		changes = append(changes, dayChange)
	}

	return talib.StdDev(changes, period, math.Sqrt(1)*100)
}

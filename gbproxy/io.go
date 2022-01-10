package gbproxy

import (
	"math"
)

func writeToPins(value uint, pins []GameBoyPin) {
	for i, p := range pins {
		ab := uint(math.Pow(2, float64(i)))
		p.SetState((value & ab) >= uint(1))
	}
}

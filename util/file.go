package util

import (
	"math"
	"strconv"
)

var (
	suffixes = []string{"B", "KB", "MB", "GB", "TB"}
)

func GetReadableSize(b int64) string {
	if b == 0 {
		return "0B"
	}

	fb := float64(b)

	base := math.Log(fb) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + string(getSuffix)
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow

	return
}

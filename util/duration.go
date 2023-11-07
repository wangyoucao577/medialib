package util

// DurationInSeconds returns duration in milliseconds calculated by given timescale.
func DurationInSeconds(duration uint64, timescale uint64) float64 {
	if timescale == 0 {
		return 0
	}

	return float64(duration) / float64(timescale)
}

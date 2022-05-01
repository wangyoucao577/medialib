// Package time1904 provides help utils to calculate time since 1904, by contrast time normally represents since 1970.
package time1904

import "time"

var baseTime = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)

// Unix is similar with time.Unix, but calculates start from 1904 instead of 1970.
func Unix(sec int64, nsec int64) time.Time {
	sec1970 := int64(sec) + baseTime.Unix()
	return time.Unix(sec1970, nsec)
}

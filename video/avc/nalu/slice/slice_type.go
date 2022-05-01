package slice

// Type returns slice type human-readable representation.
func Type(t int) string {
	switch t {
	case 0, 5:
		return "P"
	case 1, 6:
		return "B"
	case 2, 7:
		return "I"
	case 3, 8:
		return "SP"
	case 4, 9:
		return "SI"
	}
	return "unknown"
}

package sps

// profile_idc definition, ISO/IEC-14496-10 Annex A
const (
	ProfileIDCBaseLine = 66
	ProfileIDCMain     = 77
	ProfileIDCExtended = 88
	ProfileIDCHigh     = 100
	ProfileIDCHigh10   = 110
	ProfileIDCHigh422  = 122
	ProfileIDCHigh444  = 244
)

var profileNames = map[int]string{
	ProfileIDCBaseLine: "Baseline",
	ProfileIDCMain:     "Main",
	ProfileIDCExtended: "Extended",
	ProfileIDCHigh:     "High",
	ProfileIDCHigh10:   "High10",
	ProfileIDCHigh422:  "High422",
	ProfileIDCHigh444:  "High444",
}

// ProfileName returns name of the profile_idc.
func ProfileName(t uint8) string {
	n, ok := profileNames[int(t)]
	if !ok {
		return "unknown"
	}
	return n
}

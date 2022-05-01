package esds

// ClassTagInfo contains basic information of class tag, such as name, etc.
type ClassTagInfo struct {
	Name string `json:"name"`
}

// ISO/IEC-14496-1 7.2.2.1 Table-1 List of Class Tags for Descriptors
const (
	ClassTagForbidden uint8 = iota
	ClassTagObjectDescrTag
	ClassTagInitialObjectDescrTag
	ClassTagES_DescrTag
	ClassTagDecoderConfigDescrTag
	ClassTagDecSpecificInfoTag
	ClassTagSLConfigDescrTag

	//TODO: more class tags for descriptors
)

var classTags = map[uint8]ClassTagInfo{
	ClassTagForbidden: {Name: "Forbidden"},

	ClassTagObjectDescrTag:        {Name: "ObjectDescrTag"},
	ClassTagInitialObjectDescrTag: {Name: "InitialObjectDescrTag"},
	ClassTagES_DescrTag:           {Name: "ES_DescrTag"},
	ClassTagDecoderConfigDescrTag: {Name: "DecoderConfigDescrTag"},
	ClassTagDecSpecificInfoTag:    {Name: "DecSpecificInfoTag"},
	ClassTagSLConfigDescrTag:      {Name: "SLConfigDescrTag"},
}

func classTagName(t uint8) string {
	info, ok := classTags[t]
	if !ok {
		return "unknown"
	}
	return info.Name
}

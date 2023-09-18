package dump

import "fmt"

// Format represents marshal supported format type.
type Format string

// supported marshal formats
const (
	FormatJSON          = "json"
	FormatJSONFormatted = "json_formatted"
	FormatYAML          = "yaml"
	FormatYML           = "yml"
	FormatCSV           = "csv"
)

var supportedFormats = map[Format]struct{}{
	FormatJSON:          {},
	FormatJSONFormatted: {},
	FormatYAML:          {},
	FormatYML:           {},
	FormatCSV:           {},
}

// String implements Stringer.
func (f Format) String() string {
	return string(f)
}

// FormatsHelper returns formats help information, including supported formats as well as extra notes.
func FormatsHelper() string {
	var s string

	for k := range supportedFormats {
		if len(s) != 0 {
			s += ","
		}
		s += k.String()
	}
	s += ". "
	s += "\nNote that 'csv' only available for 'no parse' content."

	return s
}

// GetFormat parses string to get Format.
func GetFormat(s string) (Format, error) {
	_, ok := supportedFormats[Format(s)]
	if !ok {
		return "unknown", fmt.Errorf("unknown format %s", s)
	}
	return Format(s), nil
}

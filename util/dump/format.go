package dump

import "fmt"

// Format represents marshal supported format type.
type Format int

// supported marshal formats
const (
	FormatUnkonwn = iota
	FormatJSON
	FormatJSONNewLines
	FormatYAML
	FormatCSV
)

var supportedFormats = map[string]Format{
	"json":          FormatJSON,
	"json_newlines": FormatJSONNewLines,
	"yaml":          FormatYAML,
	"csv":           FormatCSV,
}

// FormatsHelper returns formats help information, including supported formats as well as extra notes.
func FormatsHelper() string {
	var s string

	for k := range supportedFormats {
		if len(s) != 0 {
			s += ","
		}
		s += k
	}
	s += ". "
	s += "\nNote that 'csv' only available for 'no parse' content."

	return s
}

// GetFormat parses string to get Format.
func GetFormat(s string) (Format, error) {
	f, ok := supportedFormats[s]
	if !ok {
		return FormatUnkonwn, fmt.Errorf("unknown format %s", s)
	}
	return f, nil
}

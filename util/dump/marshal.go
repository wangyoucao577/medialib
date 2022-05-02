// Package dump provides dump utilities by various format.
package dump

import (
	"fmt"
)

// Marshaler defines serveral Marshal interfaces in one place.
type Marshaler interface {

	// JSON marshals object to JSON represtation.
	JSON() ([]byte, error)

	// JSONIndent marshals object to JSON with indent represent, it has same parameters with json.MarshalIndent.
	JSONIndent(prefix, indent string) ([]byte, error)

	// YAML marshals object to YAML representation.
	YAML() ([]byte, error)

	// CSV marshals object to CSV representation.
	CSV() ([]byte, error)
}

// Marshal marshals data by format.
func Marshal(m Marshaler, f Format) ([]byte, error) {
	switch f {
	case FormatYAML:
		return m.YAML()
	case FormatJSON:
		return m.JSON()
	case FormatJSONNewLines:
		return m.JSONIndent("", "\t")
	case FormatCSV:
		return m.CSV()
	}

	// default
	return nil, fmt.Errorf("unknown format %d", f)
}

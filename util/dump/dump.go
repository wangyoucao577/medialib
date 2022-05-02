package dump

import (
	"fmt"
	"io"
)

// DumpToWriter marshals data and dump to io.Writer.
func DumpToWriter(m Marshaler, f Format, w io.Writer) error {
	d, err := Marshal(m, f)
	if err != nil {
		return err
	}

	n, err := w.Write(d)
	if err != nil {
		return err
	}
	if n != len(d) {
		return fmt.Errorf("write incompleted, expect %d but write %d bytes", len(d), n)
	}
	return nil
}

// Dump marshals data and dump to output.
func Dump(m Marshaler, f Format, output string) error {
	w, closer, err := CreateOutput(output)
	if err != nil {
		return err
	}
	if closer != nil {
		defer closer.Close()
	}

	return DumpToWriter(m, f, w)
}

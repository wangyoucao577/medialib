// Package exit provides utility to exit process with expected behavior, e.g., fail the shell.
package exit

import (
	"os"
)

// Fail exits the process with return code `-1` to fail shell.
func Fail() {
	os.Exit(-1)
}

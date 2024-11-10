package server

import (
	"fmt"

	"github.com/Wa4h1h/memdb/pkg"
)

func formatString(value string) string {
	return fmt.Sprintf("%s%d\r\n%s\r\n", pkg.StringIdentifier, len(value), value)
}

func formatError(errtype string, err string) string {
	return fmt.Sprintf("%s%s %s\r\n", pkg.ErrorIdentifier, errtype, err)
}

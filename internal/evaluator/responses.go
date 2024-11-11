package evaluator

import (
	"fmt"

	"github.com/Wa4h1h/memdb/internal/utils"
)

func FormatString(value string) string {
	return fmt.Sprintf("%s%d\r\n%s\r\n", utils.StringIdentifier, len(value), value)
}

func FormatError(errtype string, err string) string {
	return fmt.Sprintf("%s%s %s\r\n", utils.ErrorIdentifier, errtype, err)
}

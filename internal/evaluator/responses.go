package evaluator

import (
	"fmt"
	"strconv"

	"github.com/Wa4h1h/memdb/internal/utils"
)

func FormatString(value string) string {
	integer, err := strconv.Atoi(value)
	if err == nil {
		return fmt.Sprintf("%s%d\r\n%d\r\n",
			utils.IntegerIdentifier, len(value), integer)
	}

	if utils.StringIsSlice(value) {
		return fmt.Sprintf("%s%d\r\n%s\r\n",
			utils.ListIdentifier, len(value)-2, value)
	}

	boolean, err := strconv.ParseBool(value)
	if err == nil {
		return fmt.Sprintf("%s%d\r\n%v\r\n",
			utils.BooleanIdentifier, len(value), boolean)
	}

	return fmt.Sprintf("%s%d\r\n%s\r\n",
		utils.StringIdentifier, len(value), value)
}

func FormatError(errtype string, err string) string {
	return fmt.Sprintf("%s%s %s\r\n", utils.ErrorIdentifier, errtype, err)
}

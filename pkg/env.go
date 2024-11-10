package pkg

import (
	"fmt"
	"os"
	"strconv"
)

type Env interface {
	uint | string | bool
}

func GetEnv[T Env](defaultValue string, required bool, key string) T {
	var val T

	envVal, ok := os.LookupEnv(key)
	if !ok {
		if required {
			panic(fmt.Sprintf("env %s is required", key))
		}

		envVal = defaultValue
	}

	switch ptr := any(&val).(type) {
	case *string:
		*ptr = envVal
	case *bool:
		tmp, err := strconv.ParseBool(envVal)
		if err != nil {
			panic(fmt.Sprintf("can not parse evn variable to bool %s=%s", key, envVal))
		}

		*ptr = tmp
	case *uint:
		tmp, err := strconv.ParseUint(envVal, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("can not parse evn variable to uint %s=%s", key, envVal))
		}

		*ptr = uint(tmp)
	}

	return val
}

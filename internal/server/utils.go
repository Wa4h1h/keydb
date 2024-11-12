package server

import (
	"errors"
	"fmt"
	"net"

	"github.com/Wa4h1h/keydb/internal/evaluator"

	"github.com/Wa4h1h/keydb/internal/utils"
)

func (s *Server) simpleWrite(conn net.Conn, msg string) {
	if _, err := conn.Write([]byte(evaluator.FormatError("ERROR", msg))); err != nil {
		s.Logger.Error(fmt.Sprintf("error while writing to conn:%s", err.Error()))
	}
}

func (s *Server) checkError(err error) string {
	errStr := err.Error()

	switch {
	case errors.Is(err, utils.ErrUnknownCommand):
		return evaluator.FormatError("ERROR", errStr)
	case errors.Is(err, utils.ErrNotFoundItem):
		return evaluator.FormatError("NOT_FOUND", errStr)
	case errors.Is(err, utils.ErrParsingTTL) || errors.Is(err, utils.ErrParsingToInt):
		return evaluator.FormatError("PARSE_INTEGER", errStr)
	case errors.Is(err, utils.ErrItemNotRemoved):
		return evaluator.FormatError("ITEM_NOT_REMOVED", errStr)
	case errors.Is(err, utils.ErrMissingOptions):
		return evaluator.FormatError("MISSING_OPTIONS", errStr)
	case errors.Is(err, utils.ErrMalformedSlice):
		return evaluator.FormatError("MALFORMED_LIST", errStr)
	case errors.Is(err, utils.ErrElementNotinList):
		return evaluator.FormatError("ITEM_NOT_IN_LIST", errStr)
	case errors.Is(err, utils.ErrParsingToBool):
		return evaluator.FormatError("PARSE_BOOL", errStr)
	default:
		return ""
	}
}

package server

import (
	"errors"
	"fmt"
	"net"

	"github.com/Wa4h1h/memdb/pkg"
)

func (s *Server) simpleWrite(conn net.Conn, msg string) {
	if _, err := conn.Write([]byte(formatError("ERROR", msg))); err != nil {
		s.Logger.Error(fmt.Sprintf("error while writing to conn:%s", err.Error()))
	}
}

func (s *Server) checkError(err error) string {
	errStr := err.Error()

	switch {
	case errors.Is(err, pkg.ErrUnknownCommand):
		return formatError("ERROR", errStr)
	case errors.Is(err, pkg.ErrNotFoundItem):
		return formatError("NOT_FOUND", errStr)
	case errors.Is(err, pkg.ErrParsingTTL) || errors.Is(err, pkg.ErrParsingToInt):
		return formatError("PARSE_INTEGER", errStr)
	case errors.Is(err, pkg.ErrItemNotRemoved):
		return formatError("ITEM_NOT_REMOVED", errStr)
	case errors.Is(err, pkg.ErrMissingOptions):
		return formatError("MISSING_OPTIONS", errStr)
	case errors.Is(err, pkg.ErrMalformedSlice):
		return formatError("MALFORMED_LIST", errStr)
	case errors.Is(err, pkg.ErrElementNotinList):
		return formatError("ITEM_NOT_IN_LIST", errStr)
	case errors.Is(err, pkg.ErrParsingToBool):
		return formatError("PARSE_BOOL", errStr)
	default:
		return ""
	}
}

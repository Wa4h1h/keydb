package pkg

import "errors"

var (
	ErrNotFoundItem          = errors.New("item not present")
	ErrItemNotRemoved        = errors.New("item could not be removed")
	ErrUnknownCommand        = errors.New("unknown command")
	ErrParsingTTL            = errors.New("parsing ttl error")
	ErrParsingToInt          = errors.New("parsing string to int error")
	ErrParsingToBool         = errors.New("parsing string to bool error")
	ErrMissingOptions        = errors.New("missing options")
	ErrMalformedSlice        = errors.New("slice is malformed")
	ErrElementNotinList      = errors.New("element not in the list")
	ErrPartOfBodyWentMissing = errors.New("error body not fully sent")
	ErrReading               = errors.New("error while reading from connection")
	ErrWriting               = errors.New("error while writing from connection")
)

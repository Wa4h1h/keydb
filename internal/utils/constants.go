package utils

type Identifier string

const (
	StringIdentifier  Identifier = "+"
	IntegerIdentifier Identifier = ":"
	BooleanIdentifier Identifier = "#"
	ListIdentifier    Identifier = "*"
	ErrorIdentifier   Identifier = "-"
)

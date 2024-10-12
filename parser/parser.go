package parser

type Parser interface {
	Parse(response []byte) (interface{}, error)
}

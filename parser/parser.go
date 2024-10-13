package parser

type Parser[T any] interface {
	Parse(response []byte) (T, error)
}

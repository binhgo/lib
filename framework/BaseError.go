package framework

type BaseError interface {
	Error() string
	GetErrNo() int
	ErrNo() string
}

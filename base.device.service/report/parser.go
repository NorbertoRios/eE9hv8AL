package report

//IParser interface to describe parser
type IParser interface {
	Parse(packet []byte) ([]IMessage, error)
	GetUnknownAck(packet []byte) []byte
}

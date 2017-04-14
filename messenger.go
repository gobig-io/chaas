package bot

// Messenger is the interface to implement for all messengers
type Messenger interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Listen(*Message) error
	Respond([]*Result) error
}

package messaging

type Messaging interface {
	Send(title, body string) error
}

func New() Messaging {
	return NewHipChat("jqRJ2CO8ZGEMoZiTEHcYEdc0k6vaMVQZpgCYmsP9")
}

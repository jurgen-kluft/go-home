package messaging

import (
	"fmt"
	"github.com/jurgen-kluft/go-home/messaging/hipchat"
)

type HipChat struct {
	client *hipchat.Client
}

func NewHipChat(token string) Messaging {
	im = &HipChat{}
	hipchat.AuthTest = false
	im.client = hipchat.NewClient(token)
	return im
}

func (im *HipChat) Send(title, body string) error {
	return fmt.Errorf("Not implemented")
}

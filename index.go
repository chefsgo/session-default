package session_default

import (
	"github.com/chefsgo/session"
)

func Driver() session.Driver {
	return &defaultDriver{}
}

func init() {
	session.Register("default", Driver())
}

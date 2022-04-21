package session

import (
	"github.com/chefsgo/chef"
)

func Driver() chef.SessionDriver {
	return &defaultSessionDriver{}
}

func init() {
	chef.Register("default", Driver())
}

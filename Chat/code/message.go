package code

import (
	"time"
)
// message represents a single message
type Message struct {

	Name string
	Message string
	//timestamp of when msg was send
	When time.Time
}

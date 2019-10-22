package queue

import (
	"fmt"
	"time"
)

type event struct {
	timestamp time.Time
}

func (e event) String() string {
	return fmt.Sprintf("Event[%s]", e.timestamp)
}

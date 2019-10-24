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

func (e event) MarshalBinary() ([]byte, error) {
	return e.timestamp.MarshalBinary()
}

func (e *event) UnmarshalBinary(b []byte) error {
	var ts time.Time
	if err := ts.UnmarshalBinary(b); err != nil {
		return err
	}
	e.timestamp = ts
	return nil
}

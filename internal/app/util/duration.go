package util

import (
	"fmt"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshallText(t []byte) error {
	duration, err := time.ParseDuration(string(t))
	if err != nil {
		return fmt.Errorf("duration parse failed: %s", t)
	}
	d.Duration = duration
	return nil
}

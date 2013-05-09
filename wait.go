package timefind

import (
	"time"
)

func (q *TQuery) WaitNext() <-chan time.Time {
	now := time.Now()
	d := q.Next(now).Sub(now)

	/* add 10 second to make sure that the timer is fired at the
	expected minute boundary */
	return time.NewTimer(d + 10*time.Second).C
}

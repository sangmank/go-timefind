/*
 Copyright (C) 2013 Sangman Kim

 This Source Code Form is subject to the terms of the Mozilla Public
 License, v. 2.0. If a copy of the MPL was not distributed with this
 file, You can obtain one at

 http://mozilla.org/MPL/2.0/. 
*/

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

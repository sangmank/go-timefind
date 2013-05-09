package timefind

import (
	"log"
	"time"
)

func (q *TQuery) Match(t time.Time) bool {
	idxs := []int{int(t.Month()), t.Day(), t.Hour(), t.Minute(), int(t.Weekday())}
	for i, idx := range idxs {
		if !q.fields[i].Test(uint(idx)) {
			return false
		}
	}
	return true
}

func (q *TQuery) Next(t time.Time) time.Time {
	var startT time.Time

	/* start from the next minute */
	startT = time.Unix(((t.Unix()/60)+1)*60, 0)

	startUnixSec := startT.Unix()
	startMin := startT.Minute()

	/* check remaining minutes in the given hour */
	for min := startMin; min < 60; min++ {
		t1 := time.Unix(startUnixSec+int64((min-startMin)*60), 0)
		if q.Match(t1) {
			return t1
		}
	}

	/* at this point, the next minute would be minMatchMin.
	only check time by increasing hours
	*/
	minMatchMin, err := q.fields[MINUTE].LowestSetIndex()
	if err != nil {
		log.Fatalf("minute has no element in it: %v", err)
	}
	startT = time.Unix(startUnixSec+int64((60-startMin+int(minMatchMin))*60), 0)
	startUnixSec = startT.Unix()

	startHour := startT.Hour()
	for hour := startHour; hour < 24; hour++ {
		t1 := time.Unix(startUnixSec+int64((hour-startHour)*60*60), 0)
		if q.Match(t1) {
			return t1
		}
	}

	/* similarly, check by increasing days */
	minMatchHour, err := q.fields[HOUR].LowestSetIndex()
	if err != nil {
		log.Fatalf("hour has no element in it: %v", err)
	}

	startT = time.Unix(startUnixSec+int64((24-startHour+int(minMatchHour))*60*60), 0)
	startUnixSec = startT.Unix()

	for i := 0; i < 366; i++ {
		t1 := time.Unix(startUnixSec+int64(i*60*60*24), 0)
		if q.Match(t1) {
			return t1
		}
	}

	log.Fatal("ERROR: no date was selected")
	return time.Unix(0, 0)
}

func (q *TQuery) NextNow() time.Time {
	return q.Next(time.Now())
}

func (q *TQuery) NextList(t time.Time, n int) []time.Time {
	result := make([]time.Time, n)
	for i := 0; i < n; i++ {
		t1 := q.Next(t)
		result[i] = t1
		t = t1
	}
	return result
}

func (q *TQuery) NextNowList(n int) []time.Time {
	return q.NextList(time.Now(), n)
}

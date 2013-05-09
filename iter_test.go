package timefind

import (
	"testing"
	"time"
)

func TestNext(t *testing.T) {
	q, err := Hour(3)
	if err != nil {
		t.Error(err)
	}
	q, err = q.Minute(4)
	if err != nil {
		t.Error(err)
	}

	tv, err := time.Parse(time.ANSIC, "Mon Jan 2 15:04:05 2006")
	if err != nil {
		t.Fatal(err)
	}

	t1 := q.Next(tv)
	if t1.Hour() != 3 {
		t.Errorf("hour should be 3, not %d", t1.Hour())
	}
	if t1.Day() != 3 {
		t.Errorf("day should be 3, not %d", t1.Day())
	}
	if t1.Month() != 1 {
		t.Errorf("month should be 1, not %d", t1.Month())
	}
	if t1.Year() != 2006 {
		t.Errorf("year should be 2006, not %d", t1.Year())
	}
	if t1.Minute() != 4 {
		t.Errorf("minute should be 0, not %d", t1.Minute())
	}

	tq := q.Next(q.Next(tv))
	if tq.Hour() != 3 {
		t.Errorf("hour should be 3, not %d", t1.Hour())
	}
	if tq.Day() != 4 {
		t.Errorf("day should be 4, not %d", t1.Day())
	}
	if tq.Minute() != 4 {
		t.Errorf("minute should be 0, not %d", t1.Minute())
	}

	times := q.NextList(tv, 5)
	for i := 0; i < 4; i++ {
		if times[i+1].Unix()-times[i].Unix() != 24*60*60 {
			t.Errorf("time difference should be %d, not %v",
				24*60*60, times[i+1].Unix()-times[i].Unix())
		}
	}
}

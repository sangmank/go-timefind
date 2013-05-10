package timefind

import (
	"testing"
)

func TestPattern(t *testing.T) {
	if patternRegular.FindString("*/2") == "" {
		t.Error("patternRegular not matched")
	}

	if patternRegular.FindString("*") != "" {
		t.Error("patternRegular not matched")
	}

	if patternRanges.FindString("0,20,30,40-50") == "" {
		t.Error("patternRegular not matched")
	}

	if patternRanges.FindString("0") == "" {
		t.Error("patternRegular not matched")
	}

	if patternRanges.FindString("0-10") == "" {
		t.Error("patternRegular not matched")
	}

	if patternRanges.FindString("0-10,20") == "" {
		t.Error("patternRegular not matched")
	}

	if patternRanges.FindString("0-10,20-30,40") == "" {
		t.Error("patternRegular not matched")
	}

	if patternRanges.FindString("0-10,20,40-50") == "" {
		t.Error("patternRegular not matched")
	}

	if patternRanges.FindString("0-10,20,40-50") == "" {
		t.Error("patternRegular not matched")
	}
}

func TestTQuery(t *testing.T) {
	q, err := Months("1-3,5-7,9-11")
	if err != nil {
		t.Error(err.Error())
	}
	if q.fields[MONTH].Count() != 9 {
		t.Error("Months setting error")
	}

	q, err = Days("1-8")
	if err != nil {
		t.Error(err.Error())
	}
	if q.fields[DAY].Count() != 8 {
		t.Error("days setting error")
	}

	q, err = WeekDays("0-7")
	if err != nil {
		t.Error(err.Error())
	}
	if q.fields[DAY_WEEK].Count() != 7 {
		t.Error("days setting error")
	}

	q1, _ := Months("1-3")
	q2, _ := Months("4-6")
	_, err = And(q1, q2)
	if err == nil {
		t.Error("empty set should report an error")
	}

	q3, _ := Days("4-6")
	_, err = And(q1, q3)
	if err != nil {
		t.Error("This should be non-empty set")
	}

	q4, _ := Months("2,4-12")
	_, err = And(q1, q4)
	if err != nil {
		t.Error("2 should be overlapped")
	}

	q, err = NewFromString("* * * 1 *")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	q.Hour(3)
	q.Minute(0, 30)
	if q.ToString() != "0,30 3 * 1 *" {
		t.Errorf("Expected: 0,30 3 * 1 *, Give: %s", q.ToString())		
	}

	q, err = NewFromString("30,60 * * 1 *")
	if err == nil {
		t.Errorf("No error, but we expected an error %s", q.ToString())
	}
	
}

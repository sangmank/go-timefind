package timefind

import (
	"fmt"
	"github.com/sangmank/bitset"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TQueryField int

const (
	MONTH TQueryField = iota
	DAY
	HOUR
	MINUTE
	DAY_WEEK
	FIELD_END
)

func (f TQueryField) GetLimit() (min, max int) {
	switch f {
	case MONTH:
		return 1, 12
	case DAY:
		return 1, 31
	case HOUR:
		return 0, 23
	case MINUTE:
		return 0, 59
	case DAY_WEEK:
		return 0, 7
	default:
		log.Fatal("getLimit out of range")
	}
	return
}

func (f TQueryField) Index(n int) int {
	/* 0 and 7 in DAY_WEEK both represent Sunday */
	if f == DAY_WEEK {
		if n == 7 {
			return 0
		}
		return n
	}

	min, _ := f.GetLimit()
	return n - min
}

var fieldNames = []string{"month", "day", "hour", "minute", "dayofweek"}

func (f TQueryField) ToString() string {
	i := int(f)
	if i < len(fieldNames) {
		return fieldNames[i]
	} else {
		log.Fatal("ToString out of range")
	}
	return ""
}

/* intermediate node for logical relationship */
type TQuery struct {
	fields     []*bitset.BitSet
	minT, maxT *time.Time
}

var ErrEmpty = fmt.Errorf("No date is available")

func newQuery() *TQuery {
	q := new(TQuery)
	q.fields = make([]*bitset.BitSet, FIELD_END)

	for f := MONTH; f < FIELD_END; f++ {
		min, max := f.GetLimit()
		set := bitset.New(uint(max - min + 1))
		set.SetAll()
		q.fields[f] = set
	}

	q.minT = nil
	q.maxT = nil

	return q
}

func initQuery(tq TQueryField, data ...int) (*TQuery, error) {
	q := newQuery()
	return q.intersection(tq, data...)
}

func initQueryStr(tq TQueryField, data string) (*TQuery, error) {
	q := newQuery()
	return q.parseStr(tq, data)
}

/* int-version */

func Month(months ...int) (*TQuery, error) {
	return initQuery(MONTH, months...)
}

func Day(days ...int) (*TQuery, error) {
	return initQuery(DAY, days...)
}

func Hour(hours ...int) (*TQuery, error) {
	return initQuery(HOUR, hours...)
}

func Minute(minutes ...int) (*TQuery, error) {
	return initQuery(MINUTE, minutes...)
}

func WeekDay(days ...int) (*TQuery, error) {
	return initQuery(DAY_WEEK, days...)
}

/* string versions */

func Months(months string) (*TQuery, error) {
	return initQueryStr(MONTH, months)
}

func Days(days string) (*TQuery, error) {
	return initQueryStr(DAY, days)
}

func Hours(hours string) (*TQuery, error) {
	return initQueryStr(HOUR, hours)
}

func Minutes(minutes string) (*TQuery, error) {
	return initQueryStr(MINUTE, minutes)
}

func WeekDays(days string) (*TQuery, error) {
	return initQueryStr(DAY_WEEK, days)
}

/* TQuery-based int version */

func (q *TQuery) intersection(f TQueryField, data ...int) (*TQuery, error) {
	min, max := f.GetLimit()
	newset := bitset.New(uint(max - min + 1))

	for _, datum := range data {
		if datum > max || datum < min {
			return nil, fmt.Errorf("%s should be between %d and %d. Given: %d", f.ToString(), min, max, datum)
		}
		newset.Set(uint(f.Index(datum)))
	}

	q.fields[f] = q.fields[f].Intersection(newset)
	if q.fields[f].Count() == 0 {
		return nil, ErrEmpty
	}

	return q, nil
}

func (q *TQuery) Month(months ...int) (*TQuery, error) {
	return q.intersection(MONTH, months...)
}

func (q *TQuery) Day(days ...int) (*TQuery, error) {
	return q.intersection(DAY, days...)
}

func (q *TQuery) Hour(hours ...int) (*TQuery, error) {
	return q.intersection(HOUR, hours...)
}

func (q *TQuery) Minute(minutes ...int) (*TQuery, error) {
	return q.intersection(MINUTE, minutes...)
}

func (q *TQuery) WeekDay(days ...int) (*TQuery, error) {
	return q.intersection(DAY_WEEK, days...)
}

/* TQuery-based string versions */

var patternRegular, _ = regexp.Compile(`^[*]/([0-9]+)$`)
var patternRanges, _ = regexp.Compile(`^(?P<term>[0-9]+(-[0-9]+)?)(,(?P<term>[0-9]+(-[0-9]+)?))*$`)
var patternRange, _ = regexp.Compile(`^([0-9]+)-([0-9]+)$`)

func (q *TQuery) parseStr(f TQueryField, selector string) (*TQuery, error) {
	min, max := f.GetLimit()
	newset := bitset.New(uint(max - min + 1))

	if selector == "*" {
		return q, nil
	}

	if patternRegular.FindString(selector) != "" {
		interval, err := strconv.ParseInt(selector[2:], 10, 32)
		if err != nil {
			return nil, err
		}

		for i := min; i <= max; i += int(interval) {
			newset.Set(uint(f.Index(i)))
		}
	} else if patternRanges.FindString(selector) != "" {
		ranges := strings.Split(selector, ",")
		for _, r := range ranges {
			values := patternRange.FindStringSubmatch(r)
			if values != nil {
				// range
				minRange, _ := strconv.ParseInt(values[1], 10, 32)
				maxRange, _ := strconv.ParseInt(values[2], 10, 32)
				for i := minRange; i <= maxRange; i++ {
					newset.Set(uint(f.Index(int(i))))
				}
			} else {
				// single value
				val, _ := strconv.ParseInt(r, 10, 32)
				newset.Set(uint(f.Index(int(val))))
			}
		}
	} else {
		return nil, fmt.Errorf("Not supported selector %s", selector)
	}

	q.fields[f] = q.fields[f].Intersection(newset)
	return q, nil
}

func (q *TQuery) Months(months string) (*TQuery, error) {
	return q.parseStr(MONTH, months)
}

func (q *TQuery) Days(days string) (*TQuery, error) {
	return q.parseStr(DAY, days)
}

func (q *TQuery) Hours(hours string) (*TQuery, error) {
	return q.parseStr(HOUR, hours)
}

func (q *TQuery) Minutes(minutes string) (*TQuery, error) {
	return q.parseStr(MINUTE, minutes)
}

func (q *TQuery) WeekDays(days string) (*TQuery, error) {
	return q.parseStr(DAY_WEEK, days)
}

func (q *TQuery) Before(t time.Time) (*TQuery, error) {
	q.maxT = &t
	if q.IsEmpty() {
		return nil, ErrEmpty
	}
	return q, nil
}

func (q *TQuery) After(t time.Time) (*TQuery, error) {
	q.minT = &t
	if q.IsEmpty() {
		return nil, ErrEmpty
	}
	return q, nil
}

func (q *TQuery) Between(t1, t2 time.Time) *TQuery {
	if t1.Before(t2) {
		q.minT = &t1
		q.maxT = &t2
	} else {
		q.minT = &t2
		q.maxT = &t1
	}
	return q
}

func (q *TQuery) IsEmpty() bool {
	for _, set := range q.fields {
		if set.Count() == 0 {
			return true
		}
	}

	if q.maxT != nil && q.minT != nil && q.maxT.Before(*q.minT) {
		return true
	}

	return false
}

func andMinMaxT(q1, q2 *TQuery) (minT, maxT *time.Time) {
	/* setting minT */
	if q1.minT == nil {
		minT = q2.minT
	} else if q2.minT == nil {
		minT = q1.minT
	} else {
		if q1.minT.Before(*q2.minT) {
			minT = q2.minT
		} else {
			minT = q1.minT
		}
	}

	/* setting maxT */
	if q1.maxT == nil {
		maxT = q2.maxT
	} else if q2.maxT == nil {
		maxT = q1.maxT
	} else {
		if q1.maxT.After(*q2.maxT) {
			maxT = q2.maxT
		} else {
			maxT = q1.maxT
		}
	}
	return
}

func And(q1, q2 *TQuery) (*TQuery, error) {
	q := newQuery()

	for i, _ := range q.fields {
		q.fields[i] = q1.fields[i].Intersection(q2.fields[i])
	}

	q.minT, q.maxT = andMinMaxT(q1, q2)

	if q.IsEmpty() {
		return nil, ErrEmpty
	}

	return q, nil
}

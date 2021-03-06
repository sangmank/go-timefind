/*
 Copyright (C) 2013 Sangman Kim

 This Source Code Form is subject to the terms of the Mozilla Public
 License, v. 2.0. If a copy of the MPL was not distributed with this
 file, You can obtain one at

 http://mozilla.org/MPL/2.0/. 
*/

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
	MINUTE TQueryField = iota
	HOUR
	DAY
	MONTH 
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

func (f TQueryField) Name() string {
	i := int(f)
	if i < len(fieldNames) {
		return fieldNames[i]
	} else {
		log.Fatal("Name out of range")
	}
	return ""
}

/* intermediate node for logical relationship */
type TQuery struct {
	fields  []*bitset.BitSet
	minT    *time.Time
}

var ErrEmpty = fmt.Errorf("No date is available")

func newQuery() *TQuery {
	q := new(TQuery)
	q.fields = make([]*bitset.BitSet, FIELD_END)

	for f := TQueryField(0); f < FIELD_END; f++ {
		min, max := f.GetLimit()
		set := bitset.New(uint(max - min + 1))
		set.SetAll()
		q.fields[f] = set
	}

	q.minT = nil

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
			return nil, fmt.Errorf("%s should be between %d and %d. Given: %d", f.Name(), min, max, datum)
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
		// patterns like */2, */3
		// interval becomes 2 or 3 in these cases
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
				if int(minRange) < min || int(minRange) > max {
					return nil, fmt.Errorf("Min out of range %d. expected [%d, %d]", minRange, min, max)
				}
				if int(maxRange) < min || int(maxRange) > max {
					return nil, fmt.Errorf("Max out of range %d. expected [%d, %d]", minRange, min, max)
				}
				
				for i := minRange; i <= maxRange; i++ {
					newset.Set(uint(f.Index(int(i))))
				}
			} else {
				// single value
				val, _ := strconv.ParseInt(r, 10, 32)
				if int(val) < min || int(val) > max {
					return nil, fmt.Errorf("Out of range %d. expected [%d, %d]", val, min, max)
				}
				
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

func (q *TQuery) After(t time.Time) (*TQuery, error) {
	q.minT = &t
	if q.IsEmpty() {
		return nil, ErrEmpty
	}
	return q, nil
}

func (q *TQuery) IsEmpty() bool {
	for _, set := range q.fields {
		if set.Count() == 0 {
			return true
		}
	}

	return false
}

func andMinT(q1, q2 *TQuery) (minT *time.Time) {
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
	
	return
}

func And(q1, q2 *TQuery) (*TQuery, error) {
	q := newQuery()

	for i, _ := range q.fields {
		q.fields[i] = q1.fields[i].Intersection(q2.fields[i])
	}

	q.minT = andMinT(q1, q2)

	if q.IsEmpty() {
		return nil, ErrEmpty
	}

	return q, nil
}

func bitsetToString(set *bitset.BitSet, f TQueryField) string {
	min, _ := f.GetLimit()

	if set.Len() == 0 {
		return ""
	}
	
	if set.Len() == set.Count() {
		return "*"
	}

	if set.Len() == 1 {
		index, _ := set.LowestSetIndex()
		return string(index)
	}
	
	indices := set.SetIndices()
	idxstrs := make([]string, len(indices))
	for i := range idxstrs  {
		idxstrs[i] = fmt.Sprintf("%d", indices[i] + uint(min))
	}
	return strings.Join(idxstrs, ",")
}

func (q *TQuery) ToString() string {
	strs := make([]string, len(q.fields))
	
	for i, _ := range q.fields {
		strs[i] = bitsetToString(q.fields[i], TQueryField(i))
	}
	
	return strings.Join(strs, " ")
}

func New(s string) (*TQuery, error) {
	q := newQuery()
	strs := strings.Split(s, " ")
	if len(strs) != int(FIELD_END) {
		return nil, fmt.Errorf("There should be five entries (minute, hour, day of month, month, day of week)")
	}
	
	for i, field := range strs {
		_, err := q.parseStr(TQueryField(i), field)
		if err != nil {
			return nil, err
		}
	}
	
	return q, nil
}

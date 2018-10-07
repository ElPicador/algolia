package datetree

import (
	"sort"
	"sync"
	"time"
)

// Tree is a structure that allows to insert and retrieve date-indexed hits from
// Algolia's HN Search.
type Tree struct {
	Years      map[int]*YearNode
	TotalCount int
}

// NewTree returns an initialized tree.
func NewTree() *Tree {
	return &Tree{Years: map[int]*YearNode{}}
}

// Insert inserts a new hit from HN search into a Tree.
func (t *Tree) Insert(address string, ti time.Time) {

	year := ti.Year()
	yn, ok := t.Years[year]
	if !ok {
		yn = new(YearNode)
		yn.Hits = map[string]int{}
		t.Years[year] = yn
	}
	yn.Insert(address, ti)
	t.TotalCount++
}

// Count returns the number of hits for a specific date.
func (t *Tree) Count(s Search) int {
	var ret int
	if y, ok := t.Years[s.Year]; ok && y != nil {
		ret = y.Count(s)
	}

	return ret
}

// IndexPopularity creates an index for each node containing information about
// hit popularity.
func (t *Tree) IndexPopularity() {
	var wg sync.WaitGroup
	for _, v := range t.Years {
		wg.Add(1)
		go func(y *YearNode) {
			defer wg.Done()
			y.IndexPopularity()
		}(v)
	}

	wg.Wait()
}

// Popular returns the most popular hits for a specific date.
func (t *Tree) Popular(s Search) []Popularity {
	var ret []Popularity
	if y, ok := t.Years[s.Year]; ok && y != nil {
		ret = y.Popular(s)
	}

	return ret
}

// YearNode is a node representing a year.
type YearNode struct {
	Months   [12]*MonthNode
	Hits     Hits
	PopIndex []Popularity
}

// Insert inserts a new node in a YearNode.
func (y *YearNode) Insert(address string, t time.Time) {
	month := t.Month()
	m := y.Months[month-1]
	if m == nil {
		m = new(MonthNode)
		m.Hits = map[string]int{}
		y.Months[month-1] = m
	}

	y.Hits[address]++
	m.Insert(address, t)
}

// Count returns the number of hits for a specific date.
func (y *YearNode) Count(s Search) int {
	if !s.Month.Valid {
		return len(y.Hits)
	}

	var ret int
	if m := y.Months[s.Month.Int-1]; m != nil {
		ret = m.Count(s)
	}

	return ret
}

// IndexPopularity creates an index for each node containing information about
// hit popularity.
func (y *YearNode) IndexPopularity() {
	var wg sync.WaitGroup
	for _, v := range y.Months {
		if v != nil {
			wg.Add(1)
			go func(m *MonthNode) {
				defer wg.Done()
				m.IndexPopularity()
			}(v)
		}
	}

	y.PopIndex = y.Hits.IndexPopularity()
	wg.Wait()
}

// Popular returns the most popular hits for a specific date.
func (y *YearNode) Popular(s Search) []Popularity {
	if !s.Month.Valid {
		return queryPopularity(s.Popularity, y.PopIndex)
	}

	var ret []Popularity
	if m := y.Months[s.Month.Int-1]; m != nil {
		ret = m.Popular(s)
	}

	return ret
}

// MonthNode is a node representing a month.
type MonthNode struct {
	Days     [31]*DayNode
	Hits     Hits
	PopIndex []Popularity
}

// Insert inserts a new hit from HN search into a Tree.
func (m *MonthNode) Insert(address string, t time.Time) {
	day := t.Day()
	d := m.Days[day-1]
	if d == nil {
		d = new(DayNode)
		d.Hits = map[string]int{}
		m.Days[day-1] = d
	}

	m.Hits[address]++
	d.Insert(address, t)
}

// Count returns the number of hits for a specific date.
func (m *MonthNode) Count(s Search) int {
	if !s.Day.Valid {
		return len(m.Hits)
	}

	var ret int
	if d := m.Days[s.Day.Int-1]; d != nil {
		ret = d.Count(s)
	}

	return ret
}

// IndexPopularity creates an index for each node containing information about
// hit popularity.
func (m *MonthNode) IndexPopularity() {
	var wg sync.WaitGroup
	for _, v := range m.Days {
		if v != nil {
			wg.Add(1)
			go func(d *DayNode) {
				defer wg.Done()
				d.IndexPopularity()
			}(v)
		}
	}

	m.PopIndex = m.Hits.IndexPopularity()
	wg.Wait()
}

// Popular returns the most popular hits for a specific date.
func (m *MonthNode) Popular(s Search) []Popularity {
	if !s.Day.Valid {
		return queryPopularity(s.Popularity, m.PopIndex)
	}

	var ret []Popularity
	if d := m.Days[s.Day.Int-1]; d != nil {
		ret = d.Popular(s)
	}

	return ret
}

// DayNode is a node representing a day.
type DayNode struct {
	Hours    [24]*HourNode
	Hits     Hits
	PopIndex []Popularity
}

// Insert insert a new node in a DayNode.
func (d *DayNode) Insert(address string, t time.Time) {
	hour := t.Hour()
	h := d.Hours[hour]
	if h == nil {
		h = new(HourNode)
		h.Hits = map[string]int{}
		d.Hours[hour] = h
	}
	d.Hits[address]++
	h.Insert(address, t)
}

// Count returns the number of hits for a specific date.
func (d *DayNode) Count(s Search) int {
	if !s.Hour.Valid {
		return len(d.Hits)
	}

	var ret int
	if h := d.Hours[s.Hour.Int]; h != nil {
		ret = h.Count(s)
	}

	return ret
}

// IndexPopularity creates an index for each node containing information about
// hit popularity.
func (d *DayNode) IndexPopularity() {
	var wg sync.WaitGroup
	for _, v := range d.Hours {
		if v != nil {
			wg.Add(1)
			go func(h *HourNode) {
				defer wg.Done()
				h.IndexPopularity()
			}(v)
		}
	}

	d.PopIndex = d.Hits.IndexPopularity()
	wg.Wait()
}

// Popular returns the most popular hits for a specific date.
func (d *DayNode) Popular(s Search) []Popularity {
	if !s.Hour.Valid {
		return queryPopularity(s.Popularity, d.PopIndex)
	}

	var ret []Popularity
	if h := d.Hours[s.Hour.Int]; h != nil {
		ret = h.Popular(s)
	}

	return ret
}

// HourNode is a node representing an hour.
type HourNode struct {
	Minutes  [60]*MinuteNode
	Hits     Hits
	PopIndex []Popularity
}

// Insert inserts a new node in a HourNode.
func (h *HourNode) Insert(address string, t time.Time) {
	minute := t.Minute()
	m := h.Minutes[minute]
	if m == nil {
		m = new(MinuteNode)
		m.Hits = map[string]int{}
		h.Minutes[minute] = m
	}
	h.Hits[address]++
	m.Insert(address, t)
}

// Count returns the number of hits for a specific date.
func (h *HourNode) Count(s Search) int {
	if !s.Minute.Valid {
		return len(h.Hits)
	}

	var ret int
	if m := h.Minutes[s.Minute.Int]; m != nil {
		ret = m.Count(s)
	}

	return ret
}

// IndexPopularity creates an index for each node containing information about
// hit popularity.
func (h *HourNode) IndexPopularity() {
	var wg sync.WaitGroup
	for _, v := range h.Minutes {
		if v != nil {
			wg.Add(1)
			go func(m *MinuteNode) {
				defer wg.Done()
				m.IndexPopularity()
			}(v)
		}
	}

	h.PopIndex = h.Hits.IndexPopularity()
	wg.Wait()
}

// Popular returns the most popular hits for a specific date.
func (h *HourNode) Popular(s Search) []Popularity {
	if !s.Minute.Valid {
		return queryPopularity(s.Popularity, h.PopIndex)
	}

	var ret []Popularity
	if m := h.Minutes[s.Minute.Int]; m != nil {
		ret = m.Popular(s)
	}

	return ret
}

// MinuteNode is a node representing a minute.
type MinuteNode struct {
	Seconds  [60]*SecondNode
	Hits     Hits
	PopIndex []Popularity
}

// Insert inserts a new node in a MinuteNode.
func (m *MinuteNode) Insert(address string, t time.Time) {
	second := t.Second()
	s := m.Seconds[second]
	if s == nil {
		s = new(SecondNode)
		s.Hits = map[string]int{}
		m.Seconds[second] = s
	}
	m.Hits[address]++
	s.Hits[address]++
}

// Count returns the number of hits for a specific date.
func (m *MinuteNode) Count(s Search) int {
	if !s.Second.Valid {
		return len(m.Hits)
	}

	var ret int
	if sec := m.Seconds[s.Second.Int]; sec != nil {
		ret = len(sec.Hits)
	}

	return ret
}

// IndexPopularity creates an index for each node containing information about
// hit popularity.
func (m *MinuteNode) IndexPopularity() {
	var wg sync.WaitGroup
	for _, v := range m.Seconds {
		if v != nil {
			wg.Add(1)
			go func(s *SecondNode) {
				defer wg.Done()
				s.PopIndex = s.Hits.IndexPopularity()
			}(v)
		}
	}

	m.PopIndex = m.Hits.IndexPopularity()
	wg.Wait()
}

// Popular returns the most popular hits for a specific date.
func (m *MinuteNode) Popular(s Search) []Popularity {
	if !s.Second.Valid {
		return queryPopularity(s.Popularity, m.PopIndex)
	}

	var ret []Popularity
	if sec := m.Seconds[s.Second.Int]; sec != nil {
		ret = queryPopularity(s.Popularity, sec.PopIndex)
	}

	return ret
}

// SecondNode is a node representing a second.
type SecondNode struct {
	Hits     Hits
	PopIndex []Popularity
}

// Popularity represents the popularity of an address in HN Search.
type Popularity struct {
	Query string
	Count int
}

// Hits represents a aggregate of HN Search hits indexed by their address.
type Hits map[string]int

// IndexPopularity takes the hits of a node and order them in the returned array
// so they can be retrieved easily.
func (h Hits) IndexPopularity() []Popularity {
	if len(h) == 0 {
		return nil
	}

	ret := make([]Popularity, 0, len(h))
	for k, v := range h {
		ret = append(ret, Popularity{Count: v, Query: k})
	}

	sort.Slice(ret, func(i int, j int) bool {
		// if counts are equal, order strings asc
		if ret[i].Count == ret[j].Count {
			return ret[i].Query < ret[j].Query
		}
		return ret[i].Count > ret[j].Count
	})

	return ret
}

// queryPopularity returns the n most popular queries from the Popularity slice
// passed as parameter. If n is bigger than the length of the array, the whole array
// is returned.
func queryPopularity(n int, p []Popularity) []Popularity {
	if n >= len(p) {
		return p
	}

	return p[:n]
}

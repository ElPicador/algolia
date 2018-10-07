package datetree

import "tpaulmyer/algolia/null"

// Search is a structure that can be used to perform searches in a DateTree.
// The values are nullable in order to search for the results of a specific
// time range. The values being set in this structure should be obtained from
// a time.Time variable in order to ensure date's validity. Querying an unexisting
// time value (for example hour 44) will result in a panic.
type Search struct {
	Year       int
	Month      null.Int
	Day        null.Int
	Hour       null.Int
	Minute     null.Int
	Second     null.Int
	Popularity int
}

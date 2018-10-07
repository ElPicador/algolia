// Package datetree provide a B-tree like data structure specialized for dates up
// to the second. It allows to insert new data and to peform searches efficiently.
// A date tree must be first filled with the data to be indexed before being used.
// The tree can also index queries by popularity by calling the IndexPopularity
// method.
package datetree

package datetree_test

import (
	"fmt"
	"time"
	"tpaulmyer/algolia/datetree"
	"tpaulmyer/algolia/null"
)

func ExampleTree() {

	// Create a tree and insert data in it.
	tree := datetree.NewTree()
	d1 := time.Date(2018, time.October, 12, 21, 0, 0, 0, time.UTC)
	d2 := time.Date(2018, time.October, 12, 18, 10, 5, 0, time.UTC)
	tree.Insert("https://www.algolia.com/", d1)
	tree.Insert("https://www.algolia.com/", d2)
	tree.IndexPopularity()
	fmt.Println(tree.TotalCount)

	// Perform count and popularity searches in it.
	s := datetree.Search{
		Year:       2018,
		Month:      null.Int{Valid: true, Int: int(time.October)},
		Popularity: 10,
	}
	fmt.Println(tree.Count(s))
	fmt.Println(tree.Popular(s))
	// Output: 2
	// 1
	// [{https://www.algolia.com/ 2}]
}

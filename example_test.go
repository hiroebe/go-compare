package compare_test

import (
	"fmt"

	"github.com/hiroebe/go-compare"
)

func Example() {
	type Item struct {
		Key   string
		Value string
	}
	type Object struct {
		ID    int
		Items []*Item
	}

	obj1 := Object{
		ID: 1,
		Items: []*Item{
			{Key: "item1", Value: "value1"},
			{Key: "item2", Value: "value2"},
		},
	}
	obj2 := Object{
		ID: 2,
		Items: []*Item{
			{Key: "item1", Value: "value1"},
			{Key: "item3", Value: "value3"},
		},
	}
	diff := compare.Diff(obj1, obj2, func(obj Object) compare.C {
		return compare.C{
			// compare `ID` by `==`.
			compare.Comparable("ID", obj.ID),
			// compare `Items` as slice.
			compare.Slice("Items", obj.Items, func(item *Item) compare.C {
				// define nested `compare.C` to compare each `Item`.
				return compare.C{
					compare.Comparable("Key", item.Key),
					compare.Comparable("Value", item.Value),
				}
			}),
		}
	})

	fmt.Println(diff)

	// Output:
	// ID: 1 != 2
	// Items[1].Key: "item2" != "item3"
	// Items[1].Value: "value2" != "value3"
}

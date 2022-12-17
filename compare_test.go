package compare_test

import (
	"fmt"
	"testing"

	"github.com/hiroebe/go-compare"
)

func TestDiff(t *testing.T) {
	t.Parallel()

	type Value struct {
		Data []byte
	}
	type Item struct {
		Key   string
		Value Value
	}
	type Object struct {
		ID    int
		Items []*Item
	}

	obj1 := Object{
		ID: 1,
		Items: []*Item{
			{Key: "item1", Value: Value{Data: []byte("value1")}},
			{Key: "item2", Value: Value{Data: []byte("value2")}},
		},
	}
	obj2 := Object{
		ID: 2,
		Items: []*Item{
			{Key: "item1", Value: Value{Data: []byte("value1")}},
			{Key: "item3", Value: Value{Data: []byte("value3")}},
		},
	}
	diff := compare.Diff(obj1, obj2, func(obj Object) compare.C {
		return compare.C{
			compare.Comparable("ID", obj.ID),
			compare.Slice("Items", obj.Items, func(item *Item) compare.C {
				return compare.C{
					compare.Comparable("Key", item.Key),
					compare.Nest("Value", item.Value, func(value Value) compare.C {
						return compare.C{
							compare.Func("Data", value.Data, func(v1, v2 []byte) error {
								if s1, s2 := string(v1), string(v2); s1 != s2 {
									return fmt.Errorf("%s != %s", s1, s2)
								}
								return nil
							}),
						}
					}),
				}
			}),
		}
	})
	want := `ID: 1 != 2
Items[1].Key: "item2" != "item3"
Items[1].Value.Data: value2 != value3`
	if diff != want {
		t.Errorf("got:\n%s\nwant:\n%s", diff, want)
	}
}

package compare_test

import (
	"fmt"
	"testing"

	"github.com/hiroebe/go-compare"
)

func TestComparable(t *testing.T) {
	t.Parallel()

	type S struct {
		F1 int
		F2 string
	}

	tests := map[string]struct {
		diff string
		want string
	}{
		"diff: int": {
			diff: compare.Diff(1, 2, func(v int) compare.C {
				return compare.C{compare.Comparable("v", v)}
			}),
			want: `v: 1 != 2`,
		},
		"diff: string": {
			diff: compare.Diff("s1", "s2", func(v string) compare.C {
				return compare.C{compare.Comparable("v", v)}
			}),
			want: `v: "s1" != "s2"`,
		},
		"diff: struct": {
			diff: compare.Diff(S{1, "s1"}, S{2, "s2"}, func(v S) compare.C {
				return compare.C{compare.Comparable("v", v)}
			}),
			want: `v: compare_test.S{F1:1, F2:"s1"} != compare_test.S{F1:2, F2:"s2"}`,
		},
		"diff: pointer": {
			diff: compare.Diff(&S{}, &S{}, func(v *S) compare.C {
				return compare.C{compare.Comparable("v", v)}
			}),
			want: `v: &compare_test.S{F1:0, F2:""} != &compare_test.S{F1:0, F2:""}`,
		},
		"no label": {
			diff: compare.Diff(1, 2, func(v int) compare.C {
				return compare.C{compare.Comparable("", v)}
			}),
			want: `1 != 2`,
		},
		"no diff": {
			diff: compare.Diff(1, 1, func(v int) compare.C {
				return compare.C{compare.Comparable("v", v)}
			}),
			want: ``,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.diff != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", tt.diff, tt.want)
			}
		})
	}
}

func TestComparableSlice(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		diff string
		want string
	}{
		"diff: slice of int": {
			diff: compare.Diff([]int{1, 2}, []int{3, 4}, func(v []int) compare.C {
				return compare.C{compare.ComparableSlice("v", v)}
			}),
			want: `v[0]: 1 != 3
v[1]: 2 != 4`,
		},
		"no diff": {
			diff: compare.Diff([]int{1, 2}, []int{1, 2}, func(v []int) compare.C {
				return compare.C{compare.ComparableSlice("v", v)}
			}),
			want: ``,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.diff != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", tt.diff, tt.want)
			}
		})
	}
}

func TestComparablePointer(t *testing.T) {
	t.Parallel()

	type S struct {
		F int
	}

	tests := map[string]struct {
		diff string
		want string
	}{
		"diff: pointer": {
			diff: compare.Diff(&S{1}, &S{2}, func(v *S) compare.C {
				return compare.C{compare.ComparablePointer("v", v)}
			}),
			want: `v: compare_test.S{F:1} != compare_test.S{F:2}`,
		},
		"diff: nil pointer": {
			diff: compare.Diff(&S{1}, nil, func(v *S) compare.C {
				return compare.C{compare.ComparablePointer("v", v)}
			}),
			want: `v: &compare_test.S{F:1} != (*compare_test.S)(nil)`,
		},
		"no diff: pointer": {
			diff: compare.Diff(&S{1}, &S{1}, func(v *S) compare.C {
				return compare.C{compare.ComparablePointer("v", v)}
			}),
			want: ``,
		},
		"no diff: nil pointer": {
			diff: compare.Diff(nil, nil, func(v *S) compare.C {
				return compare.C{compare.ComparablePointer("v", v)}
			}),
			want: ``,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.diff != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", tt.diff, tt.want)
			}
		})
	}
}

func TestFunc(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		diff string
		want string
	}{
		"diff": {
			diff: compare.Diff(1, 1, func(v int) compare.C {
				return compare.C{compare.Func("v", v, func(v1, v2 int) error {
					return fmt.Errorf("custom diff")
				})}
			}),
			want: `v: custom diff`,
		},
		"no diff": {
			diff: compare.Diff(1, 2, func(v int) compare.C {
				return compare.C{compare.Func("v", v, func(v1, v2 int) error {
					return nil
				})}
			}),
			want: ``,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.diff != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", tt.diff, tt.want)
			}
		})
	}
}

func TestNest(t *testing.T) {
	t.Parallel()

	type S struct {
		F1 int
		F2 string
	}
	type P struct {
		S S
	}

	tests := map[string]struct {
		diff string
		want string
	}{
		"diff": {
			diff: compare.Diff(P{S{1, "s1"}}, P{S{2, "s2"}}, func(p P) compare.C {
				return compare.C{
					compare.Nest("S", p.S, func(s S) compare.C {
						return compare.C{
							compare.Comparable("F1", s.F1),
							compare.Comparable("F2", s.F2),
						}
					}),
				}
			}),
			want: `S.F1: 1 != 2
S.F2: "s1" != "s2"`,
		},
		"no diff": {
			diff: compare.Diff(P{S{1, "s1"}}, P{S{1, "s1"}}, func(v P) compare.C {
				return compare.C{
					compare.Nest("S", v.S, func(s S) compare.C {
						return compare.C{
							compare.Comparable("F1", s.F1),
							compare.Comparable("F2", s.F2),
						}
					}),
				}
			}),
			want: ``,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.diff != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", tt.diff, tt.want)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	t.Parallel()

	type S struct {
		F int
	}

	tests := map[string]struct {
		diff string
		want string
	}{
		"diff": {
			diff: compare.Diff([]S{{1}, {2}}, []S{{3}, {4}}, func(vs []S) compare.C {
				return compare.C{
					compare.Slice("vs", vs, func(v S) compare.C {
						return compare.C{
							compare.Comparable("F", v.F),
						}
					}),
				}
			}),
			want: `vs[0].F: 1 != 3
vs[1].F: 2 != 4`,
		},
		"no diff": {
			diff: compare.Diff([]S{{1}, {2}}, []S{{1}, {2}}, func(vs []S) compare.C {
				return compare.C{
					compare.Slice("vs", vs, func(v S) compare.C {
						return compare.C{
							compare.Comparable("F", v.F),
						}
					}),
				}
			}),
			want: ``,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.diff != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", tt.diff, tt.want)
			}
		})
	}
}

package listing

import (
	"math"
	"reflect"
	"testing"
)

func fixture() []Item {
	names := []string{
		"apple", "banana", "avocado", "blueberry", "apricot",
		"cherry", "acai", "cranberry", "almond", "date",
	}
	items := make([]Item, len(names))
	for i, n := range names {
		items[i] = Item{ID: i + 1, Name: n}
	}
	return items
}

func ids(items []Item) []int {
	out := make([]int, len(items))
	for i, it := range items {
		out[i] = it.ID
	}
	return out
}

func TestPage(t *testing.T) {
	tests := []struct {
		name       string
		page, size int
		want       []int
	}{
		{"first page", 1, 3, []int{1, 2, 3}},
		{"middle page", 2, 3, []int{4, 5, 6}},
		{"last page is partial", 4, 3, []int{10}},
		{"exact fit", 2, 5, []int{6, 7, 8, 9, 10}},
		{"size larger than catalogue", 1, 100, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{"beyond end", 99, 3, []int{}},
		{"first page past end", 5, 3, []int{}},
		{"page zero", 0, 3, []int{}},
		{"negative page", -1, 3, []int{}},
		{"size zero", 1, 0, []int{}},
		{"negative size", 1, -5, []int{}},
		{"huge page", math.MaxInt, 3, []int{}},
		{"huge size", 2, math.MaxInt, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(fixture())
			got := ids(s.Page(tt.page, tt.size))
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Page(%d, %d) = %v; want %v", tt.page, tt.size, got, tt.want)
			}
		})
	}
}

func TestTotalPages(t *testing.T) {
	s := NewService(fixture())
	tests := []struct{ size, want int }{
		{3, 4}, {5, 2}, {10, 1}, {100, 1}, {1, 10}, {0, 0}, {-1, 0},
	}
	for _, tt := range tests {
		if got := s.TotalPages(tt.size); got != tt.want {
			t.Errorf("TotalPages(%d) = %d; want %d", tt.size, got, tt.want)
		}
	}
}

func TestSearchDoesNotChangeCatalogue(t *testing.T) {
	s := NewService(fixture())
	got := s.Search("a")
	if want := []int{1, 3, 5, 7, 9}; !reflect.DeepEqual(ids(got), want) {
		t.Fatalf("Search(a) = %v; want %v", ids(got), want)
	}
	if want := ids(fixture()); !reflect.DeepEqual(ids(s.Items()), want) {
		t.Fatalf("catalogue changed after search: %v; want %v", ids(s.Items()), want)
	}
}

func TestSearchNoMatchIsEmpty(t *testing.T) {
	s := NewService(fixture())
	if got := s.Search("zzz"); got == nil || len(got) != 0 {
		t.Fatalf("Search(zzz) = %#v; want empty non-nil slice", got)
	}
}

func TestCallerCannotMutateCatalogue(t *testing.T) {
	s := NewService(fixture())
	s.Items()[0].Name = "hacked"
	s.Page(1, 3)[1].Name = "hacked"
	if got := s.Items()[0].Name; got != "apple" {
		t.Fatalf("first item name = %q; want apple", got)
	}
	if got := s.Items()[1].Name; got != "banana" {
		t.Fatalf("second item name = %q; want banana", got)
	}
}

func TestAppendToPageDoesNotCorruptNextPage(t *testing.T) {
	s := NewService(fixture())
	page := s.Page(1, 3)
	_ = append(page, Item{ID: 99, Name: "intruder"})
	if got := ids(s.Page(2, 3)); !reflect.DeepEqual(got, []int{4, 5, 6}) {
		t.Fatalf("Page(2, 3) = %v; want [4 5 6]", got)
	}
}

func TestServiceCopiesInput(t *testing.T) {
	in := fixture()
	s := NewService(in)
	in[0].Name = "changed by caller"
	if got := s.Items()[0].Name; got != "apple" {
		t.Fatalf("first item name = %q; want apple", got)
	}
}

func TestTagging(t *testing.T) {
	s := NewService(fixture())
	s.Tag(1, "red")
	s.Tag(6, "red")
	got := ids(s.ByTag("red"))
	if want := []int{1, 6}; !reflect.DeepEqual(got, want) {
		t.Fatalf("ByTag(red) = %v; want %v", got, want)
	}
	if got := s.ByTag("blue"); len(got) != 0 {
		t.Fatalf("ByTag(blue) = %v; want empty", got)
	}
}

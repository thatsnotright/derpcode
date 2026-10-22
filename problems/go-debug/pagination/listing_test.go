package listing

import (
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

func TestPageFirst(t *testing.T) {
	s := NewService(fixture())
	got := ids(s.Page(1, 3))
	if want := []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Page(1, 3) = %v; want %v", got, want)
	}
}

func TestPageLastIsPartial(t *testing.T) {
	s := NewService(fixture())
	got := ids(s.Page(4, 3))
	if want := []int{10}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Page(4, 3) = %v; want %v", got, want)
	}
}

func TestPageBeyondEnd(t *testing.T) {
	s := NewService(fixture())
	if got := s.Page(99, 3); len(got) != 0 {
		t.Fatalf("Page(99, 3) = %v; want empty", got)
	}
}

func TestPageInvalidArgs(t *testing.T) {
	s := NewService(fixture())
	for _, c := range [][2]int{{0, 3}, {-1, 3}, {1, 0}, {1, -5}} {
		if got := s.Page(c[0], c[1]); len(got) != 0 {
			t.Errorf("Page(%d, %d) = %v; want empty", c[0], c[1], got)
		}
	}
}

func TestTotalPages(t *testing.T) {
	s := NewService(fixture())
	if got := s.TotalPages(3); got != 4 {
		t.Fatalf("TotalPages(3) = %d; want 4", got)
	}
	if got := s.TotalPages(5); got != 2 {
		t.Fatalf("TotalPages(5) = %d; want 2", got)
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

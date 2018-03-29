package mist

import (
	"strings"
	"testing"
)

// BenchmarkAddRemoveSimple
func BenchmarkAddRemoveSimpleFilter(b *testing.B) {
	node := newFilter()
	keys := []string{"a b & c | d & e ^ f & g | h &", "a ^ b"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Add(keys)
		node.Remove(keys)
	}
}

// BenchmarkMatchSimple
func BenchmarkMatchSimpleFilter(b *testing.B) {
	node := newFilter()
	node.Add([]string{"a b & c | & d ^ e"})
	b.ResetTimer()
	keys := []string{"a", "b", "c", "d"}
	for i := 0; i < b.N; i++ {
		node.Match(keys)
	}
}

// BenchmarkAddRemoveComplex benchmarks to see how fast mist can add/remove keys to
// a subscription
func BenchmarkAddRemoveComplexFilter(b *testing.B) {
	node := newFilter()

	// create a giant slice of random keys
	keys := [][]string{}
	keys = append(keys, []string{randKey(), randKey()})
	for i := 0; i < b.N; i++ {
		keys = append(keys, []string{randKey()})
	}

	b.ResetTimer()

	// add/remove keys
	for _, k := range keys {
		node.Add(k)
		node.Remove(k)
	}
}

// BenchmarkMatchComplex benchmarks to see how fast mist can match a set of keys on a
// subscription
func BenchmarkMatchComplexFilter(b *testing.B) {

}

// TestEmptySubscription
func TestEmptySubscriptionFilter(t *testing.T) {
	node := newFilter()
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Unexpected tags in new subscription!")
	}
}

func TestCompileFilter(t *testing.T) {
	node := newFilter()
	filterExpression := "a b & c &"
	node.Add([]string{filterExpression})

	if len(node.filters) != 1 {
		t.Fatalf("Failed to compile filter!")
	}
	if len(node.varToIndex) != 3 {
		t.Fatalf("Failed to compile filter!")
	}
	if len(node.varValues) != 3 {
		t.Fatalf("Failed to compile filter!")
	}
	expr := node.filters[filterExpression]
	if len(expr) != 5 {
		t.Fatalf("Failed to compile filter correctly!")
	}
}

// TestAddRemoveSimple
func TestAddRemoveSimpleFilter(t *testing.T) {
	node := newFilter()

	node.Add([]string{"a"})
	if len(node.ToSlice()) != 1 {
		t.Fatalf("Failed to add filter")
	}

	node.Remove([]string{"a"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove node")
	}
}

// TestAddRemoveComplex
func TestAddRemoveComplexFilter(t *testing.T) {
	node := newFilter()

	// add/remove unordered keys; should remove
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"c", "b", "a"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove node")
	}

	// add/remove incomplete keys; should not remove
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"a"})
	node.Remove([]string{"b"})
	node.Remove([]string{"c"})
	node.Remove([]string{"a", "b"})
	node.Remove([]string{"b", "c"})
	node.Remove([]string{"a", "c"})
	node.Remove([]string{"b", "c", "d"})
	node.Remove([]string{"a", "b", "c", "d"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Filter should be empty")
	}

	// add duplicate keys; should only add once
	node.Add([]string{"a", "b", "c"})
	node.Add([]string{"a", "b", "c"})
	if len(node.ToSlice()) != 3 {
		t.Fatalf("Duplicate filters added")
	}
	node.Remove([]string{"a", "b", "c"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove filters")
	}

	// remove duplicate keys; should only remove once
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"a", "b", "c"})
	node.Remove([]string{"c", "b", "a"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove nodes")
	}

	// add duplicate remote one; should leave no nodes
	node.Add([]string{"a", "b", "c"})
	node.Add([]string{"a", "b", "c"})
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"a", "b", "c"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove nodes")
	}
}

// TestList
func TestListFilter(t *testing.T) {
	node := newFilter()

	// test simple list; length should be 1 and value should be "a"
	node.Add([]string{"a"})
	list := node.ToSlice()
	if len(list) != 1 {
		t.Fatalf("Wrong number of keys - Expecting 1 got %d", len(list))
	}
	if len(list[0]) != 1 {
		t.Fatalf("Wrong number of keys - Expecing 2 got %d", len(list[0]))
	}
	if strings.Join(list[0], ",") != "a" {
		t.Fatalf("Wrong tags - Expecing 'a' got %s", list[0])
	}

	node.Add([]string{"a", "b"})
	list = node.ToSlice()
	if len(list) != 2 {
		t.Fatalf("Wrong number of keys - Expecting 2 got %d", len(list))
	}
	if len(list[1]) != 1 {
		t.Fatalf("Wrong number of keys - Expecing 2 got %d", len(list[1]))
	}

	node.Add([]string{"a", "b", "c"})
	list = node.ToSlice()
	if len(list) != 3 {
		t.Fatalf("wrong length of list. Expecting 3 got %d", len(list))
	}
	if len(list[2]) != 1 {
		t.Fatalf("Wrong number of keys - Expecing 3 got %d", len(list[2]))
	}
}

// TestMatchSimple
func TestMatchSimpleFilter(t *testing.T) {
	node := newFilter()

	// simple match
	node.Add([]string{"a"})
	if !node.Match([]string{"a"}) {
		t.Fatalf("Expected match!")
	}

	node.Add([]string{"a", "b"})
	if !node.Match([]string{"a", "b"}) {
		t.Fatalf("Expected match!")
	}

	node.Add([]string{"a", "b", "c"})
	if !node.Match([]string{"a", "b", "c"}) {
		t.Fatalf("Expected match!")
	}
}

// TestMatchComplex
func TestMatchComplexFilter(t *testing.T) {
	node := newFilter()

	// match unordered keys; should match
	node.Add([]string{"a", "b", "c"})
	if !node.Match([]string{"c", "b", "a"}) {
		t.Fatalf("Expected match!")
	}
	node.Remove([]string{"a", "b", "c"})

	// match multiple subs with single match; should match
	node.Add([]string{"a", "b", "e"})
	node.Add([]string{"c"})
	if !node.Match([]string{"a", "b", "c", "d"}) {
		t.Fatalf("Expected match!")
	}
	node.Remove([]string{"a", "b", "e"})
	node.Remove([]string{"c"})

	node = newFilter()

	// match expression filter; should not match
	node.Add([]string{"a b & c &"})
	if len(node.filters["a b & c &"]) != 5 {
		t.Fatalf("Unexpected match!")
	}
	if len(node.varToIndex) != 3 {
		t.Fatalf("Unexpected match!")
	}
	if len(node.varValues) != 3 {
		t.Fatalf("Unexpected match!")
	}
	if node.Match([]string{"a"}) {
		t.Fatalf("Unexpected match!")
	}
	if node.Match([]string{"a", "b"}) {
		t.Fatalf("Unexpected match!")
	}
	if node.Match([]string{"b", "c"}) {
		t.Fatalf("Unexpected match!")
	}
	if node.Match([]string{"a", "c"}) {
		t.Fatalf("Unexpected match!")
	}
	if !node.Match([]string{"a", "b", "c"}) {
		t.Fatalf("Should match!")
	}
	node.Remove([]string{"a b & c &"})

	// match more complex expression
	node.Add([]string{"a b & c |"})
	if !node.Match([]string{"b", "c"}) {
		t.Fatalf("Should match!")
	}
	if !node.Match([]string{"c"}) {
		t.Fatalf("Should match!")
	}
	if !node.Match([]string{"a", "b"}) {
		t.Fatalf("Should match!")
	}
	node.Remove([]string{"a b & | c"})
}

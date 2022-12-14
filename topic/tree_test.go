package topic

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTreeAdd(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)

	assert.Equal(t, 1, tree.root.children["foo"].children["bar"].values[0])
}

func TestTreeAddDuplicate(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)
	tree.Add("foo/bar", 1)

	assert.Equal(t, 1, len(tree.root.children["foo"].children["bar"].values))
}

func TestTreeSet(t *testing.T) {
	tree := NewStandardTree()

	tree.Set("foo/bar", 1)

	assert.Equal(t, 1, tree.root.children["foo"].children["bar"].values[0])
}

func TestTreeSetReplace(t *testing.T) {
	tree := NewStandardTree()

	tree.Set("foo/bar", 1)
	tree.Set("foo/bar", 2)

	assert.Equal(t, 2, tree.root.children["foo"].children["bar"].values[0])
}

func TestTreeGet(t *testing.T) {
	tree := NewStandardTree()

	tree.Set("foo/#", 1)

	assert.Equal(t, 1, tree.Get("foo/#")[0])
}

func TestTreeRemove(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)
	tree.Remove("foo/bar", 1)

	assert.Equal(t, 0, len(tree.root.children))
}

func TestTreeRemoveMissing(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)
	tree.Remove("bar/baz", 1)

	assert.Equal(t, 1, len(tree.root.children))
}

func TestTreeEmpty(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)
	tree.Add("foo/bar", 2)
	tree.Empty("foo/bar")

	assert.Equal(t, 0, len(tree.root.children))
}

func TestTreeClear(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)
	tree.Add("foo/bar/baz", 1)
	tree.Clear(1)

	assert.Equal(t, 0, len(tree.root.children))
}

func TestTreeMatchExact(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)

	assert.Equal(t, 1, tree.Match("foo/bar")[0])
}

func TestTreeMatchWildcard1(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/+", 1)

	assert.Equal(t, 1, tree.Match("foo/bar")[0])
}

func TestTreeMatchWildcard2(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/#", 1)

	assert.Equal(t, 1, tree.Match("foo/bar")[0])
}

func TestTreeMatchWildcard3(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/#", 1)

	assert.Equal(t, 1, tree.Match("foo/bar/baz")[0])
}

func TestTreeMatchWildcard4(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar/#", 1)

	assert.Equal(t, 1, tree.Match("foo/bar")[0])
}

func TestTreeMatchWildcard5(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/#", 1)

	assert.Equal(t, 1, tree.Match("foo/bar/#")[0])
}

func TestTreeMatchMultiple(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)
	tree.Add("foo/+", 2)
	tree.Add("foo/#", 3)

	assert.Equal(t, 3, len(tree.Match("foo/bar")))
}

func TestTreeMatchNoDuplicates(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)
	tree.Add("foo/+", 1)
	tree.Add("foo/#", 1)

	assert.Equal(t, 1, len(tree.Match("foo/bar")))
}

func TestTreeMatchFirst(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/+", 1)

	assert.Equal(t, 1, tree.MatchFirst("foo/bar"))
}

func TestTreeMatchFirstNone(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/+", 1)

	assert.Nil(t, tree.MatchFirst("baz/qux"))
}

func TestTreeSearchExact(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)

	assert.Equal(t, 1, tree.Search("foo/bar")[0])
}

func TestTreeSearchWildcard1(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)

	assert.Equal(t, 1, tree.Search("foo/+")[0])
}

func TestTreeSearchWildcard2(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)

	assert.Equal(t, 1, tree.Search("foo/#")[0])
}

func TestTreeSearchWildcard3(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar/baz", 1)

	assert.Equal(t, 1, tree.Search("foo/#")[0])
}

func TestTreeSearchWildcard4(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)

	assert.Equal(t, 1, tree.Search("foo/bar/#")[0])
}

func TestTreeSearchWildcard5(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar/#", 1)

	assert.Equal(t, 1, tree.Search("foo/#")[0])
}

func TestTreeSearchMultiple(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo", 1)
	tree.Add("foo/bar", 2)
	tree.Add("foo/bar/baz", 3)

	assert.Equal(t, 3, len(tree.Search("foo/#")))
}

func TestTreeSearchNoDuplicates(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo", 1)
	tree.Add("foo/bar", 1)
	tree.Add("foo/bar/baz", 1)

	assert.Equal(t, 1, len(tree.Search("foo/#")))
}

func TestTreeSearchFirst(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)

	assert.Equal(t, 1, tree.SearchFirst("foo/+"))
}

func TestTreeSearchFirstNone(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)

	assert.Nil(t, tree.SearchFirst("baz/qux"))
}

func TestTreeCount(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo", 1)
	tree.Add("foo/bar", 2)
	tree.Add("foo/bar/baz", 3)
	tree.Add("foo/bar/baz", 4)
	tree.Add("quz/bar/baz", 4)

	assert.Equal(t, 5, tree.Count())
}

func TestTreeAll(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo", 1)
	tree.Add("foo/bar", 2)
	tree.Add("foo/bar/baz", 3)

	assert.Equal(t, 3, len(tree.All()))
}

func TestTreeAllNoDuplicates(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo", 1)
	tree.Add("foo/bar", 1)
	tree.Add("foo/bar/baz", 1)

	assert.Equal(t, 1, len(tree.All()))
}

func TestTreeReset(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("foo/bar", 1)
	tree.Reset()

	assert.Equal(t, 0, len(tree.root.children))
}

func TestTreeString(t *testing.T) {
	tree := NewStandardTree()

	tree.Add("", 1)
	tree.Add("/foo", 7)
	tree.Add("/foo/bar", 42)

	assert.Equal(t, "topic.Tree:\n| '' => 1\n|   'foo' => 1\n|     'bar' => 1", tree.String())
}

func BenchmarkTreeAddSame(b *testing.B) {
	tree := NewStandardTree()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Add("foo/bar", 1)
	}
}

func BenchmarkTreeAdd(b *testing.B) {
	tree := NewStandardTree()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Add("foo/bar", i)
	}
}

func BenchmarkTreeAddUnique(b *testing.B) {
	tree := NewStandardTree()

	strings := make([]string, 0, b.N)

	for i := 0; i < b.N; i++ {
		strings = append(strings, fmt.Sprintf("foo/%d", i))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Add(strings[i], 1)
	}
}

func BenchmarkTreeSetSame(b *testing.B) {
	tree := NewStandardTree()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Set("foo/bar", 1)
	}
}

func BenchmarkTreeSetUnique(b *testing.B) {
	tree := NewStandardTree()

	strings := make([]string, 0, b.N)

	for i := 0; i < b.N; i++ {
		strings = append(strings, fmt.Sprintf("foo/%d", i))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Set(strings[i], 1)
	}
}

func TestTreeMatch10Milion(t *testing.T) {
	tree := NewStandardTree()
	type Obj struct {
		Id int
	}
	for i := 0; i < 100000; i++ {
		tree.Add("foo/bar", Obj{Id: i})
	}

	start := time.Now()
	res := tree.Match("foo/bar")
	t.Log("1", len(res), time.Now().Sub(start))
}

func BenchmarkTreeMatchExact(b *testing.B) {
	tree := NewStandardTree()
	tree.Add("foo/bar", 0)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Match("foo/bar")
	}
}

func BenchmarkTreeMatchWildcardOne(b *testing.B) {
	tree := NewStandardTree()
	tree.Add("foo/+", 1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Match("foo/bar")
	}
}

func BenchmarkTreeMatchWildcardSome(b *testing.B) {
	tree := NewStandardTree()
	tree.Add("#", 1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Match("foo/bar")
	}
}

func BenchmarkTreeSearchExact(b *testing.B) {
	tree := NewStandardTree()
	tree.Add("foo/bar", 1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Search("foo/bar")
	}
}

func BenchmarkTreeSearchWildcardOne(b *testing.B) {
	tree := NewStandardTree()
	tree.Add("foo/bar", 1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Search("foo/+")
	}
}

func BenchmarkTreeSearchWildcardSome(b *testing.B) {
	tree := NewStandardTree()
	tree.Add("foo/bar", 1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Search("#")
	}
}

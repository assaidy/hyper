package hyper

import (
	"bytes"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIf(t *testing.T) {
	t.Run("picks between the two values", func(t *testing.T) {
		require.Equal(t, "yes", If(true, "yes", "no"))
		require.Equal(t, "no", If(false, "yes", "no"))
		require.Equal(t, 42, If(true, 42, 0))
		require.Equal(t, 0, If(false, 42, 0))
	})

	t.Run("returns the zero value when no alternative is given", func(t *testing.T) {
		require.Equal(t, "x", If(true, "x"))
		require.Equal(t, "", If(false, "x"))
		require.Equal(t, 0, If(false, 42))

		var ptr *int
		require.Nil(t, If(false, ptr))

		items := []string{"a"}
		require.Nil(t, If(false, items))
		require.Equal(t, items, If(true, items))
	})
}

func TestConditionalNode(t *testing.T) {
	t.Run("renders the body when the condition holds", func(t *testing.T) {
		require.Equal(t, "<p>yes</p>", mustRender(t, IfNode(true, P("yes"))))
	})

	t.Run("renders nothing when the condition fails and there is no else", func(t *testing.T) {
		require.Equal(t, "", mustRender(t, IfNode(false, P("no"))))
	})

	t.Run("renders the else fallback", func(t *testing.T) {
		node := IfNode(false, P("first")).Else(P("fallback"))
		require.Equal(t, "<p>fallback</p>", mustRender(t, node))
	})

	t.Run("first true branch wins", func(t *testing.T) {
		node := IfNode(true, P("one")).ElseIf(true, P("two")).Else(P("three"))
		require.Equal(t, "<p>one</p>", mustRender(t, node))
	})

	t.Run("skips false branches until a true one", func(t *testing.T) {
		node := IfNode(false, P("one")).ElseIf(false, P("two")).ElseIf(true, P("three")).Else(P("four"))
		require.Equal(t, "<p>three</p>", mustRender(t, node))
	})

	t.Run("all false falls through to else", func(t *testing.T) {
		node := IfNode(false, P("one")).ElseIf(false, P("two")).Else(P("three"))
		require.Equal(t, "<p>three</p>", mustRender(t, node))
	})

	t.Run("panics on a nil body", func(t *testing.T) {
		require.PanicsWithValue(t, "nil node passed to IfNode", func() {
			IfNode(true, nil)
		})
		require.PanicsWithValue(t, "nil node passed to ElseIf", func() {
			IfNode(false, P("x")).ElseIf(false, nil)
		})
		require.PanicsWithValue(t, "nil node passed to Else", func() {
			IfNode(false, P("x")).Else(nil)
		})
	})

	t.Run("propagates branch render errors", func(t *testing.T) {
		_, err := tryRender(t, IfNode(true, errNode{errBoom}))
		require.ErrorIs(t, err, errBoom)
	})

	t.Run("branches derived from a shared chain stay independent", func(t *testing.T) {
		base := IfNode(false, P("base"))
		first := base.ElseIf(true, P("first"))
		second := base.ElseIf(true, P("second"))
		require.Equal(t, "<p>first</p>", mustRender(t, first))
		require.Equal(t, "<p>second</p>", mustRender(t, second))
	})
}

func TestConditionalAttr(t *testing.T) {
	render := func(t *testing.T, attr Attribute) string {
		t.Helper()
		out, err := renderAttribute(t, attr)
		require.NoError(t, err)
		return out
	}

	t.Run("renders the attribute when the condition holds", func(t *testing.T) {
		require.Equal(t, " hidden", render(t, IfAttr(true, AttrHidden(true))))
	})

	t.Run("renders nothing when the condition fails and there is no else", func(t *testing.T) {
		require.Empty(t, render(t, IfAttr(false, AttrHidden(true))))
	})

	t.Run("renders the else fallback", func(t *testing.T) {
		attr := IfAttr(false, AttrHidden(true)).Else(AttrRequired(true))
		require.Equal(t, " required", render(t, attr))
	})

	t.Run("first true branch wins", func(t *testing.T) {
		attr := IfAttr(true, AttrHidden(true)).ElseIf(true, AttrRequired(true)).Else(AttrReadOnly(true))
		require.Equal(t, " hidden", render(t, attr))
	})

	t.Run("skips false branches until a true one", func(t *testing.T) {
		attr := IfAttr(false, AttrHidden(true)).ElseIf(false, AttrRequired(true)).ElseIf(true, AttrChecked(true)).Else(AttrReadOnly(true))
		require.Equal(t, " checked", render(t, attr))
	})

	t.Run("all false falls through to else", func(t *testing.T) {
		attr := IfAttr(false, AttrHidden(true)).ElseIf(false, AttrRequired(true)).Else(AttrReadOnly(true))
		require.Equal(t, " readonly", render(t, attr))
	})

	t.Run("branches derived from a shared chain stay independent", func(t *testing.T) {
		base := IfAttr(false, AttrHidden(true))
		first := base.ElseIf(true, AttrRequired(true))
		second := base.ElseIf(false, AttrRequired(true))
		require.Equal(t, " required", render(t, first))
		require.Empty(t, render(t, second))
	})

	t.Run("propagates errors from the matched branch", func(t *testing.T) {
		_, err := renderAttribute(t, IfAttr(true, PairAttribute{Key: " "}))
		require.EqualError(t, err, "empty/whitespace attribute key not allowed.")
	})

	t.Run("propagates errors from the else branch", func(t *testing.T) {
		_, err := renderAttribute(t, IfAttr(false, AttrHidden(true)).Else(PairAttribute{Key: " "}))
		require.EqualError(t, err, "empty/whitespace attribute key not allowed.")
	})

	t.Run("panics on a nil attribute", func(t *testing.T) {
		require.PanicsWithValue(t, "nil attribute passed to IfAttr", func() {
			IfAttr(true, nil)
		})
		require.PanicsWithValue(t, "nil attribute passed to ElseIf", func() {
			IfAttr(true, AttrHidden(true)).ElseIf(false, nil)
		})
		require.PanicsWithValue(t, "nil attribute passed to Else", func() {
			IfAttr(true, AttrHidden(true)).Else(nil)
		})
	})
}

func TestRepeat(t *testing.T) {
	t.Run("calls generate exactly n times and renders in order", func(t *testing.T) {
		calls := 0
		node := Repeat(3, func() any {
			calls++
			return LI(strconv.Itoa(calls))
		})
		require.Equal(t, "<li>1</li><li>2</li><li>3</li>", mustRender(t, node))
		require.Equal(t, 3, calls)
	})

	t.Run("zero count renders nothing and never calls generate", func(t *testing.T) {
		calls := 0
		node := Repeat(0, func() any {
			calls++
			return P("x")
		})
		require.Equal(t, "", mustRender(t, node))
		require.Zero(t, calls)
	})

	t.Run("converts generated values to nodes", func(t *testing.T) {
		values := []any{"text", 42, P("elem")}
		i := 0
		node := Repeat(len(values), func() any {
			v := values[i]
			i++
			return v
		})
		require.Equal(t, "text42<p>elem</p>", mustRender(t, node))
	})

	t.Run("negative count panics with a clear message", func(t *testing.T) {
		require.PanicsWithValue(t, "negative count passed to Repeat", func() {
			Repeat(-1, func() any { return P("x") })
		})
	})
}

func TestRange(t *testing.T) {
	t.Run("maps every element in order and calls generate per element", func(t *testing.T) {
		items := []string{"Apple", "Banana"}
		calls := 0
		node := Range(items, func(item string) any {
			calls++
			return LI(item)
		})
		require.Equal(t, "<li>Apple</li><li>Banana</li>", mustRender(t, node))
		require.Equal(t, 2, calls)
	})

	t.Run("empty input renders nothing without calling generate", func(t *testing.T) {
		calls := 0
		node := Range([]string{}, func(string) any {
			calls++
			return P("x")
		})
		require.Equal(t, "", mustRender(t, node))
		require.Zero(t, calls)
	})

	t.Run("nil input behaves like empty input", func(t *testing.T) {
		node := Range(nil, func(string) any { return P("x") })
		require.Equal(t, "", mustRender(t, node))
	})
}

func TestFragment(t *testing.T) {
	t.Run("renders children back to back without a wrapper", func(t *testing.T) {
		node := Fragment(P("Item 1"), H1("Item 2"), "Item 3")
		require.Equal(t, "<p>Item 1</p><h1>Item 2</h1>Item 3", mustRender(t, node))
	})

	t.Run("empty fragment renders nothing", func(t *testing.T) {
		require.Equal(t, "", mustRender(t, Fragment()))
	})

	t.Run("nested fragments flatten", func(t *testing.T) {
		node := Fragment(Fragment(P("a"), P("b")), P("c"))
		require.Equal(t, "<p>a</p><p>b</p><p>c</p>", mustRender(t, node))
	})

	t.Run("propagates child errors and writes nothing", func(t *testing.T) {
		out, err := tryRender(t, Fragment(P("fine"), errNode{errBoom}))
		require.ErrorIs(t, err, errBoom)
		require.Empty(t, out)
	})

	t.Run("panics on a nil argument", func(t *testing.T) {
		require.PanicsWithValue(t, "nil argument passed to Fragment", func() {
			Fragment(P("ok"), nil)
		})
	})
}

func TestOnce(t *testing.T) {
	resetOnceCache(t)

	t.Run("caches the render for a call site", func(t *testing.T) {
		calls := 0
		node := Once(func() HyperNode {
			calls++
			return DIV("once")
		})

		first := mustRender(t, node)
		second := mustRender(t, node)

		require.Equal(t, "<div>once</div>", first)
		require.Equal(t, first, second)
		require.Equal(t, 1, calls)
	})

	t.Run("all iterations of a loop share one cache entry", func(t *testing.T) {
		calls := 0
		var outputs []string
		for i := range 3 {
			node := Once(func() HyperNode {
				calls++
				return DIV(strconv.Itoa(i))
			})
			outputs = append(outputs, mustRender(t, node))
		}

		require.Equal(t, 1, calls)
		require.Equal(t, []string{"<div>0</div>", "<div>0</div>", "<div>0</div>"}, outputs)
	})
}

func TestOnceKey(t *testing.T) {
	resetOnceCache(t)

	t.Run("generate runs only on the first render", func(t *testing.T) {
		calls := 0
		node := OnceKey("test-once-key-first-render", func() HyperNode {
			calls++
			return DIV("built")
		})

		first := mustRender(t, node)
		second := mustRender(t, node)

		require.Equal(t, "<div>built</div>", first)
		require.Equal(t, first, second)
		require.Equal(t, 1, calls)
	})

	t.Run("nodes sharing a key replay the first result", func(t *testing.T) {
		firstCalls, secondCalls := 0, 0
		a := OnceKey("test-once-key-shared", func() HyperNode {
			firstCalls++
			return DIV("from-a")
		})
		b := OnceKey("test-once-key-shared", func() HyperNode {
			secondCalls++
			return DIV("from-b")
		})

		require.Equal(t, "<div>from-a</div>", mustRender(t, a))
		require.Equal(t, "<div>from-a</div>", mustRender(t, b))
		require.Equal(t, 1, firstCalls)
		require.Zero(t, secondCalls)
	})

	t.Run("distinct keys are cached independently", func(t *testing.T) {
		a := OnceKey("test-once-key-a", func() HyperNode { return DIV("a") })
		b := OnceKey("test-once-key-b", func() HyperNode { return DIV("b") })
		require.Equal(t, "<div>a</div>", mustRender(t, a))
		require.Equal(t, "<div>b</div>", mustRender(t, b))
	})

	t.Run("a failed render is not cached", func(t *testing.T) {
		calls := 0
		node := OnceKey("test-once-key-error", func() HyperNode {
			calls++
			if calls == 1 {
				return errNode{errBoom}
			}
			return DIV("recovered")
		})

		_, err := tryRender(t, node)
		require.ErrorIs(t, err, errBoom)

		require.Equal(t, "<div>recovered</div>", mustRender(t, node))
		require.Equal(t, 2, calls)
	})

	t.Run("propagates writer errors when replaying the cache", func(t *testing.T) {
		node := OnceKey("test-once-key-writer-error", func() HyperNode { return DIV("x") })
		require.Equal(t, "<div>x</div>", mustRender(t, node)) // seed the cache

		err := RenderNode(failingWriter{errBoom}, node)
		require.ErrorIs(t, err, errBoom)
	})

	t.Run("concurrent renders produce identical output", func(t *testing.T) {
		const goroutines = 8
		node := OnceKey("test-once-key-concurrent", func() HyperNode { return DIV("static") })

		outputs := make([]string, goroutines)
		errs := make([]error, goroutines)
		var wg sync.WaitGroup
		for i := range goroutines {
			wg.Go(func() {
				var buf bytes.Buffer
				errs[i] = RenderNode(&buf, node)
				outputs[i] = buf.String()
			})
		}
		wg.Wait()

		for i := range goroutines {
			require.NoError(t, errs[i])
			require.Equal(t, "<div>static</div>", outputs[i])
		}
	})
}

func TestClasses(t *testing.T) {
	require.Equal(t, "", Classes())
	require.Equal(t, "btn primary", Classes("btn", "primary"))
	require.Equal(t, "", Classes("", "   "))
	require.Equal(t, "a", Classes(" a ", "a")) // trims before dedup
	require.Equal(t, "b a", Classes("b", "a", "b"))
	require.Equal(t, "a b", Classes("a", "", "b", "a"))
	require.Equal(t, "a  b", Classes("  a  b  ")) // only the ends are trimmed
}

func TestJson(t *testing.T) {
	t.Run("marshals with sorted keys", func(t *testing.T) {
		require.Equal(t, `{"active":true,"role":"admin"}`, Json(Object{"role": "admin", "active": true}))
	})

	t.Run("supports nested values", func(t *testing.T) {
		in := Object{
			"name":  "hyper",
			"tags":  []string{"fast", "safe"},
			"count": 2,
		}
		require.Equal(t, `{"count":2,"name":"hyper","tags":["fast","safe"]}`, Json(in))
	})

	t.Run("panics on values JSON cannot represent", func(t *testing.T) {
		require.PanicsWithError(t, "json: unsupported type: chan int", func() {
			Json(make(chan int))
		})
	})
}

func TestGrowSliceCapacity(t *testing.T) {
	tests := []struct {
		name     string
		s        []int
		n        int
		wantCap  int
		wantSame bool
	}{
		{"returns the slice untouched when capacity suffices", []int{7, 8}, 2, 2, true},
		{"grows exactly one past capacity", []int{7, 8}, 3, 4, false},
		{"jumps to the next power of two", make([]int, 8), 9, 16, false},
		{"grows an empty slice to one slot", nil, 1, 1, false},
		{"ignores a zero target", []int{7, 8}, 0, 2, true},
		{"ignores a negative target", []int{7, 8}, -5, 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := growSliceCapacity(tt.s, tt.n)

			require.Equal(t, tt.wantCap, cap(got))
			require.Equal(t, len(tt.s), len(got))
			if tt.s != nil {
				require.Equal(t, tt.s, got) // elements preserved
			}
			if tt.wantSame {
				require.Same(t, &tt.s[0], &got[0])
			}
		})
	}
}

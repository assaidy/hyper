package hyper

import (
	"bytes"
	"encoding/json"
	"io"
	"math/bits"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// If returns the appropriate value based on a boolean condition. When no
// alternative is given, a false condition yields the zero value of T.
//
// This generic function is useful for inline conditional expressions when
// building elements, letting you choose between two values in a single expression.
// For conditional attributes use [IfAttr] instead: a false condition without an
// alternative yields a nil attribute, which panics in element constructors.
//
// Example:
//
//	class := If(err != nil, "text-red", "text-black")
//	class := If(err != nil, "text-red")
func If[T any](condition bool, ifTrue T, ifFalse ...T) T {
	if condition {
		return ifTrue
	}
	if len(ifFalse) > 0 {
		return ifFalse[0]
	}
	var zero T
	return zero
}

// IfAttr returns a conditional attribute chain starting with a condition.
//
// The attribute is rendered only when its condition is true. When no branch
// matches and no else branch was set, nothing is rendered (the attribute is
// omitted). Extend the chain with ElseIf and Else.
//
// Example:
//
//	IfAttr(isHidden, AttrHidden(true)).
//		Else(AttrRequired(true))
func IfAttr(condition bool, attr Attribute) conditionalAttr {
	if attr == nil {
		panic("nil attribute passed to IfAttr")
	}
	return conditionalAttr{
		ifBranches: []attrIfBranch{{condition: condition, body: attr}},
	}
}

// ElseIf adds an additional condition to the conditional chain.
//
// Example:
//
//	IfAttr(isDark, AttrHidden(true)).
//		ElseIf(isLight, AttrRequired(true))
func (me conditionalAttr) ElseIf(condition bool, attr Attribute) conditionalAttr {
	if attr == nil {
		panic("nil attribute passed to ElseIf")
	}
	me.ifBranches = append(me.ifBranches, attrIfBranch{condition: condition, body: attr})
	return me
}

// Else provides a fallback attribute when no conditions match.
//
// Example:
//
//	IfAttr(isAdmin, AttrRequired(true)).
//		Else(AttrReadOnly(true))
func (me conditionalAttr) Else(attr Attribute) Attribute {
	if attr == nil {
		panic("nil attribute passed to Else")
	}
	me.elseBranch = attr
	return me
}

// conditionalAttr represents a chain of if-else conditions for attributes.
// It is created by [IfAttr] and can be extended with ElseIf() and Else().
type conditionalAttr struct {
	ifBranches []attrIfBranch
	elseBranch Attribute
}

func (me conditionalAttr) RenderAttribute(w io.Writer) error {
	for _, n := range me.ifBranches {
		if n.condition == true {
			return n.body.RenderAttribute(w)
		}
	}

	if me.elseBranch != nil {
		return me.elseBranch.RenderAttribute(w)
	}

	// don't render anything if the condition is false and else branche was not specified.
	return nil
}

type attrIfBranch struct {
	condition bool
	body      Attribute
}

// IfNode creates a conditional node chain starting with a condition.
//
// The body is rendered only if the condition is true, otherwise
// an empty fragment is rendered (preventing nil pointer issues).
//
// Example:
//
//	IfNode(isAuthenticated, HEADER("Welcome")).
//		ElseIf(isTrial, HEADER("Try Premium")).
//		Else(BUTTON("Login"))
func IfNode(condition bool, body HyperNode) conditionalNode {
	if body == nil {
		panic("nil node passed to IfNode")
	}
	return conditionalNode{
		ifBranches: []nodeIfBranch{{condition: condition, body: body}},
		elseBranch: Fragment(),
	}
}

// ElseIf adds an additional condition to the conditional chain.
//
// Example:
//
//	IfNode(isLoggedIn, DIV("Welcome")).
//		ElseIf(isAdmin, DIV("Admin Panel"))
func (me conditionalNode) ElseIf(condition bool, body HyperNode) conditionalNode {
	if body == nil {
		panic("nil node passed to ElseIf")
	}
	me.ifBranches = append(me.ifBranches, nodeIfBranch{condition: condition, body: body})
	return me
}

// Else provides a fallback body when no conditions match.
//
// Example:
//
//	IfNode(isAdmin, DIV("Admin")).
//		Else(DIV("User"))
func (me conditionalNode) Else(body HyperNode) HyperNode {
	if body == nil {
		panic("nil node passed to Else")
	}
	me.elseBranch = body
	return me
}

// conditionalNode represents a chain of if-else conditions.
// It is created by [IfNode] and can be extended with ElseIf() and Else().
type conditionalNode struct {
	ifBranches []nodeIfBranch
	elseBranch HyperNode
}

func (me conditionalNode) RenderNode(w io.Writer) error {
	for _, n := range me.ifBranches {
		if n.condition == true {
			return RenderNode(w, n.body)
		}
	}
	return RenderNode(w, me.elseBranch)
}

// nodeIfBranch represents a single condition-body pair within a conditionalNode.
type nodeIfBranch struct {
	condition bool
	body      HyperNode
}

// Repeat generates multiple Nodes by calling a function n times.
//
// The provided function is called exactly n times, and each resulting value
// is converted to a [HyperNode] and aggregated into a single container Node.
// Using a function ensures each Node instance is unique (important for elements
// with mutable state).
//
// Example:
//
//	UL(
//		Repeat(5, func() any {
//			return LI("List item")
//		}),
//	)
func Repeat(n int, generate func() any) HyperNode {
	if n < 0 {
		panic("negative count passed to Repeat")
	}
	fragment := fragmentNode{nodes: make([]HyperNode, 0, n)}
	for range n {
		fragment.nodes = append(fragment.nodes, toHyperNode(generate()))
	}
	return fragment
}

// Range transforms a slice of items into Nodes by applying a function to each element.
//
// Each element in the input slice is transformed using the provided function, and
// all resulting values are converted to [HyperNode] and aggregated into a single
// container Node.
//
// Example:
//
//	items := []string{"Apple", "Banana", "Cherry"}
//	UL(
//		Range(items, func(item string) any {
//			return LI(item)
//		}),
//	)
func Range[T any](input []T, generate func(T) any) HyperNode {
	fragment := fragmentNode{nodes: make([]HyperNode, 0, len(input))}
	for _, item := range input {
		fragment.nodes = append(fragment.nodes, toHyperNode(generate(item)))
	}
	return fragment
}

// Fragment groups multiple nodes into a single [HyperNode] without wrapping
// them in an HTML tag. Each argument is converted to a [HyperNode] and the
// children are rendered back-to-back, so the output is identical to rendering
// each node individually.
//
// This is useful when you need to pass several nodes as one value, e.g. as the
// body of a conditional or as a single argument to another component.
//
// Example:
//
//	Fragment(
//		P("Item 1"),
//		H1("Item 2"),
//		"Item 3",
//	)
func Fragment(args ...any) HyperNode {
	fragment := fragmentNode{nodes: make([]HyperNode, 0, len(args))}
	for _, arg := range args {
		if arg == nil {
			panic("nil argument passed to Fragment")
		}
		fragment.nodes = append(fragment.nodes, toHyperNode(arg))
	}
	return fragment
}

type fragmentNode struct {
	nodes []HyperNode
}

func (me fragmentNode) RenderNode(w io.Writer) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufferPool.Put(buf)
	}()

	for _, node := range me.nodes {
		if err := node.RenderNode(buf); err != nil {
			return err
		}
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// Once is like [OnceKey] but derives the cache key automatically from the
// caller's program counter. This guarantees uniqueness without manual key management.
//
// Note: When Once is called inside a loop (for, [Repeat], [Range]), all iterations
// share the same call site and therefore the same cache key. Only the first
// iteration renders; subsequent ones reuse the cached HTML. Use [OnceKey]
// with a distinguishing value (e.g., the loop index) when each iteration needs
// its own cache entry.
//
// Example:
//
//	page := Once(func() HyperNode {
//	    return Fragment(
//	        DOCTYPE(),
//	        HTML(
//	            HEAD(TITLE("Dashboard")),
//	            BODY(H1("Welcome")),
//	        ),
//	    )
//	})
//
//go:noinline
func Once(generate func() HyperNode) HyperNode {
	var pc [1]uintptr
	if runtime.Callers(2, pc[:]) == 0 {
		panic("failed to get caller PC")
	}
	return OnceKey(strconv.FormatUint(uint64(pc[0]), 10), generate)
}

// OnceKey caches the rendered output of a component under an explicit key.
//
// The first time the returned node is rendered, generate is called to build
// the component, its output is rendered and cached. Subsequent renders replay
// the cached output without calling generate. This is useful for expensive
// static components whose tree is rebuilt per request.
//
// The key must be unique across all OnceKey calls in your application.
// Two calls with the same key share the same cache entry.
//
// Example:
//
//	page := OnceKey("dashboard-page", func() HyperNode {
//	    return Fragment(
//	        DOCTYPE(),
//	        HTML(
//	            HEAD(TITLE("Dashboard")),
//	            BODY(H1("Welcome")),
//	        ),
//	    )
//	})
func OnceKey(key string, generate func() HyperNode) HyperNode {
	return onceNode{nodeFunc: generate, key: key}
}

type onceNode struct {
	key      string
	nodeFunc func() HyperNode
}

// I benchmarked against using a map[string][]byte with a sync.RWMutex
// and found no tangible performance difference.
// I decided to use sync.Map for simplicity.
var onceCache sync.Map

func (me onceNode) RenderNode(w io.Writer) error {
	value, ok := onceCache.Load(me.key)
	if ok {
		_, err := w.Write(value.([]byte))
		return err
	}

	node := me.nodeFunc()
	var buffer bytes.Buffer
	if err := node.RenderNode(&buffer); err != nil {
		return err
	}

	result, _ := onceCache.LoadOrStore(me.key, buffer.Bytes())

	_, err := w.Write(result.([]byte))
	return err
}

// Classes joins multiple CSS class names into a single space-separated string.
// Duplicate or empty classes are filtered out.
//
// Example:
//
//	BUTTON(
//		AttrClass(Classes(
//			"btn",
//			If(err != nil, "btn-error", "btn-primary"),
//			If(isHidden, "hidden"),
//		)),
//	)
func Classes(classes ...string) string {
	if len(classes) == 0 {
		return ""
	}

	taken := make(map[string]struct{}, len(classes))
	toRender := make([]string, 0, len(classes))

	for _, c := range classes {
		trimmed := strings.TrimSpace(c)
		if trimmed == "" {
			continue
		}
		if _, ok := taken[trimmed]; ok {
			continue
		}
		taken[trimmed] = struct{}{}
		toRender = append(toRender, trimmed)
	}

	return strings.Join(toRender, " ")
}

// Json marshals v to a JSON string, panicking on error.
// Useful for embedding static JSON in templates where the value is known
// to be valid at compile time, avoiding error-handling boilerplate.
//
// Example:
//
//	FORM(Attr("hx-vals", Json(Object{"role": "admin", "active": true})))
func Json(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// Object is a shortcut for map[string]any, useful for JSON strings.
//
// Example:
//
//	Json(Object{"role": "admin", "active": true})
type Object map[string]any

// growSliceCapacity grows s's capacity to the smallest power of two
// greater than or equal to n, if necessary, while preserving its length
// and elements. If cap(s) is already sufficient, s is returned unchanged.
func growSliceCapacity[T any](s []T, n int) []T {
	if cap(s) >= n {
		return s
	}

	newCap := 1 << bits.Len(uint(n-1))
	newSlice := make([]T, len(s), newCap)
	copy(newSlice, s)
	return newSlice
}

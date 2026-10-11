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

// If starts a conditional value chain for plain values such as strings,
// numbers, or structs.
//
// Use If for inline conditional expressions when building elements, letting
// you choose between values in a single expression. Extend the chain with
// ElseIf and terminate it with Else (explicit fallback) or Value
// (zero-value fallback).
//
// Do not use If with [HyperNode] or [Attribute] values. Use [IfNode] and
// [IfAttr] instead: they are render/nil-safe — a false condition without an
// else renders nothing, while a nil node or attribute panics in element
// constructors.
//
// Example:
//
//	AttrClass(Classes(
//		"btn",
//		If(hasError, "btn-error").Else("btn-primary"),
//		If(isHidden, "hidden").Value(),
//	))
func If[T any](condition bool, value T) conditionalValue[T] {
	return conditionalValue[T]{
		ifBranches: []valueIfBranch[T]{{condition: condition, value: value}},
	}
}

// conditionalValue represents a chain of if-else conditions for plain values.
// It is created by [If] and can be extended with ElseIf, then terminated
// with Else or Value.
type conditionalValue[T any] struct {
	ifBranches []valueIfBranch[T]
	elseBranch T
}

type valueIfBranch[T any] struct {
	condition bool
	value     T
}

// ElseIf adds an additional condition to the conditional value chain.
//
// Example:
//
//	If(isAdmin, "admin").ElseIf(isMember, "member").Else("guest")
func (me conditionalValue[T]) ElseIf(condition bool, value T) conditionalValue[T] {
	me.ifBranches = append(me.ifBranches, valueIfBranch[T]{condition: condition, value: value})
	return me
}

// Else terminates the chain and returns the first matching value, or the
// fallback value when no condition matches.
//
// Example:
//
//	If(isAdmin, "admin").Else("user")
func (me conditionalValue[T]) Else(value T) T {
	me.elseBranch = value
	return me.Value()
}

// Value terminates the chain and returns the first matching value, or the
// zero value of T when no condition matches.
//
// Example:
//
//	If(isHidden, "hidden").Value() // "" when isHidden is false
func (me conditionalValue[T]) Value() T {
	for _, b := range me.ifBranches {
		if b.condition {
			return b.value
		}
	}
	return me.elseBranch
}

// IfAttr starts a conditional attribute chain.
//
// The attribute is rendered only when its condition is true. When no branch
// matches and no else branch was set, nothing is rendered (the attribute is
// omitted). Extend the chain with ElseIf and terminate it with Else.
//
// It panics if attr is nil. Use IfAttr instead of passing a nil attribute to
// an element constructor (which panics).
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

// ElseIf adds an additional condition to the conditional attribute chain.
//
// It panics if attr is nil.
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

// Else terminates the chain and provides a fallback attribute when no
// conditions match.
//
// It panics if attr is nil.
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
// It is created by [IfAttr] and can be extended with ElseIf and terminated
// with Else.
type conditionalAttr struct {
	ifBranches []attrIfBranch
	elseBranch Attribute
}

type attrIfBranch struct {
	condition bool
	body      Attribute
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

// IfNode starts a conditional node chain.
//
// The body is rendered only if the condition is true, otherwise an empty
// fragment is rendered (preventing nil pointer issues). Extend the chain with
// ElseIf and terminate it with Else.
// Without an Else, a fully-false chain renders nothing.
//
// It panics if body is nil. Use IfNode instead of passing a nil node to an
// element constructor (which panics).
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

// ElseIf adds an additional condition to the conditional node chain.
//
// It panics if body is nil.
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

// Else terminates the chain and provides a fallback body when no conditions
// match.
//
// It panics if body is nil.
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

// conditionalNode represents a chain of if-else conditions for nodes.
// It is created by [IfNode] and can be extended with ElseIf and terminated
// with Else.
type conditionalNode struct {
	ifBranches []nodeIfBranch
	elseBranch HyperNode
}

// nodeIfBranch represents a single condition-body pair within a conditionalNode.
type nodeIfBranch struct {
	condition bool
	body      HyperNode
}

func (me conditionalNode) RenderNode(w io.Writer) error {
	for _, n := range me.ifBranches {
		if n.condition == true {
			return RenderNode(w, n.body)
		}
	}
	return RenderNode(w, me.elseBranch)
}

// Repeat generates multiple nodes by calling generate exactly n times with
// the zero-based index.
//
// Each resulting value is converted to a [HyperNode] and
// aggregated into a single container node. Using a function ensures each node
// instance is unique (important for elements with mutable state).
// A zero or negative count renders nothing without calling generate.
// Ignore the index with _ when it is not needed.
//
// Example:
//
//	UL(
//		Repeat(5, func(i int) any {
//			return LI(i)
//		}),
//	)
func Repeat(n int, generate func(int) any) HyperNode {
	if n <= 0 {
		return Fragment()
	}
	fragment := fragmentNode{nodes: make([]HyperNode, 0, n)}
	for i := range n {
		fragment.nodes = append(fragment.nodes, toHyperNode(generate(i)))
	}
	return fragment
}

// Range transforms a slice of items into nodes by applying generate to each
// element with its index.
//
// Each element is transformed using the provided function, converted to a
// [HyperNode], and aggregated into a single container node.
// Nil and empty inputs render nothing without calling generate.
// Ignore the index with _ when it is not needed.
//
// Example:
//
//	items := []string{"Apple", "Banana", "Cherry"}
//	UL(
//		Range(items, func(i int, item string) any {
//			return LI(item)
//		}),
//	)
func Range[T any](input []T, generate func(int, T) any) HyperNode {
	if len(input) == 0 {
		return Fragment()
	}
	fragment := fragmentNode{nodes: make([]HyperNode, 0, len(input))}
	for i, item := range input {
		fragment.nodes = append(fragment.nodes, toHyperNode(generate(i, item)))
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
// An empty Fragment renders nothing.
//
// It panics on a nil argument.
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

// Once caches the rendered output of a component, deriving the cache key
// automatically from the caller's program counter. This guarantees uniqueness
// without manual key management. It behaves like [OnceKey] otherwise: the
// first render calls generate, caches its output, and later renders replay
// the cache without calling generate.
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
// the cached output without calling generate. A failed render is not cached.
// This is useful for expensive static components whose tree is rebuilt per request.
//
// The key must be unique across all OnceKey (and [Once]) calls in your
// application. Two calls with the same key share the same cache entry.
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
//
// Empty and whitespace-only entries are dropped, surrounding whitespace is
// trimmed, and duplicates are removed (after trimming). With no classes it
// returns "".
//
// Pair it with [If] for conditional classes: terminate single-branch cases
// with Value (yields "" when false, which Classes drops) and two-way cases
// with Else.
//
// Example:
//
//	BUTTON(
//		AttrClass(Classes(
//			"btn",
//			If(hasError, "btn-error").Else("btn-primary"),
//			If(isHidden, "hidden").Value(),
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
// Pair it with [Object] for inline maps.
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

// Object is a shortcut for map[string]any, useful with [Json] for inline
// JSON objects.
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

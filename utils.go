package hyper

import (
	"bytes"
	"encoding/json"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// IfElse returns the appropriate value based on a boolean condition.
//
// This generic function is useful for inline conditional expressions in
// builder-style code where you need to choose between two values without
// breaking the chain of method calls.
//
// Example:
//
// class := IfElse(err != nil, "text-red", "text-black")
func IfElse[T any](condition bool, result, alternative T) T {
	if condition {
		return result
	}
	return alternative
}

// IfElseZero returns the result if the condition is true, otherwise returns the zero value.
//
// This generic function is useful when you need a default/empty value when
// a condition is false without explicitly specifying the alternative.
//
// Example:
//
// class += IfElseZero(err != nil, " text-red")
func IfElseZero[T any](condition bool, result T) T {
	if condition {
		return result
	}
	var zero T
	return zero
}

// If creates a conditional node chain starting with a condition.
//
// The body is rendered only if the condition is true, otherwise
// an empty group is rendered (preventing nil pointer issues).
//
// Example:
//
//	If(isAuthenticated, HEADER()("Welcome")).
//		ElseIf(isTrial, HEADER()("Try Premium")).
//		Else(BUTTON()("Login"))
func If(condition bool, body HyperNode) conditionalNode {
	return conditionalNode{
		ifBranches: []ifBranch{{condition: condition, body: body}},
		elseBranch: Group(),
	}
}

// ElseIf adds an additional condition to the conditional chain.
//
// Example:
//
//	If(isLoggedIn, DIV()("Welcome")).
//		ElseIf(isAdmin, DIV()("Admin Panel"))
func (me conditionalNode) ElseIf(condition bool, body HyperNode) conditionalNode {
	me.ifBranches = append(me.ifBranches, ifBranch{condition: condition, body: body})
	return me
}

// Else provides a fallback body when no conditions match.
//
// Example:
//
//	If(isAdmin, DIV()("Admin")).
//		Else(DIV()("User"))
func (me conditionalNode) Else(body HyperNode) HyperNode {
	me.elseBranch = body
	return me
}

// conditionalNode represents a chain of if-else conditions.
// It is created by If() and can be extended with ElseIf() and Else().
type conditionalNode struct {
	ifBranches []ifBranch
	elseBranch HyperNode
}

func (me conditionalNode) Render(w io.Writer) error {
	for _, n := range me.ifBranches {
		if n.condition == true {
			return Render(w, n.body)
		}
	}
	return Render(w, me.elseBranch)
}

// ifBranch represents a single condition-body pair within a conditionalNode.
type ifBranch struct {
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
//	UL()(
//		Repeat(5, func() any {
//			return LI()("List item")
//		}),
//	)
func Repeat(n int, f func() any) HyperNode {
	result := Element{Tag: "", Children: make([]HyperNode, 0, n)}
	for range n {
		result.Children = append(result.Children, toHyperNode(f()))
	}
	return result
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
//	UL()(
//		Range(items, func(item string) any {
//			return LI()(item)
//		}),
//	)
func Range[T any](input []T, f func(T) any) HyperNode {
	result := Element{Tag: "", Children: make([]HyperNode, 0, len(input))}
	for _, item := range input {
		result.Children = append(result.Children, toHyperNode(f(item)))
	}
	return result
}

// Group groups multiple children without wrapping them in a tag.
// It creates a container Element with an empty Tag, which renders only its children.
//
// Example:
//
//	Group(
//		P()("Item 1"),
//		H1()("Item 2"),
//		"Item 3",
//	)
func Group(children ...any) Element {
	element := Element{Tag: ""}
	element.InsertChildren(children...)
	return element
}

// Once is like [OnceWithKey] but derives the cache key automatically from the
// caller's program counter. This guarantees uniqueness without manual key management.
//
// Note: When Once is called inside a loop (for, [Repeat], [Range]), all iterations
// share the same call site and therefore the same cache key. Only the first
// iteration renders; subsequent ones reuse the cached HTML. Use [OnceWithKey]
// with a distinguishing value (e.g., the loop index) when each iteration needs
// its own cache entry.
//
// Example:
//
//	page := Once(func() HyperNode {
//	    return Group(
//	        DOCTYPE(),
//	        HTML()(
//	            HEAD()(TITLE()("Dashboard")),
//	            BODY()(H1()("Welcome")),
//	        ),
//	    )
//	})
//
//go:noinline
func Once(f func() HyperNode) HyperNode {
	var pc [1]uintptr
	if runtime.Callers(2, pc[:]) == 0 {
		panic("failed to get caller PC")
	}
	return OnceWithKey(strconv.FormatUint(uint64(pc[0]), 10), f)
}

// OnceWithKey caches the rendered output of a component under an explicit key.
//
// The first time the returned node is rendered, f is called to build the component,
// its output is rendered and cached. Subsequent renders replay the cached output
// without calling f. This is useful for expensive static components whose tree
// is rebuilt per request.
//
// The key must be unique across all OnceWithKey calls in your application.
// Two calls with the same key share the same cache entry.
//
// Example:
//
//	page := OnceWithKey("dashboard-page", func() HyperNode {
//	    return Group(
//	        DOCTYPE(),
//	        HTML()(
//	            HEAD()(TITLE()("Dashboard")),
//	            BODY()(H1()("Welcome")),
//	        ),
//	    )
//	})
func OnceWithKey(key string, f func() HyperNode) HyperNode {
	return onceNode{nodeFunc: f, key: key}
}

type onceNode struct {
	key      string
	nodeFunc func() HyperNode
}

// I benchmarked against using a map[string][]byte with a sync.RWMutex
// and found no tangible performance difference.
// I decided to use sync.Map for simplicity.
var onceCache sync.Map

func (me onceNode) Render(w io.Writer) error {
	value, ok := onceCache.Load(me.key)
	if ok {
		_, err := w.Write(value.([]byte))
		return err
	}

	node := me.nodeFunc()
	var buffer bytes.Buffer
	if err := node.Render(&buffer); err != nil {
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
//			IfElse(err != nil, "btn-error", btn-primary),
//			IfElseZero(isHidden, "hidden"),
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

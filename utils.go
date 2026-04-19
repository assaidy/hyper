package hyper

import (
	"bytes"
	"io"
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
//	BODY()(
//		IfElse(isAdmin,
//			DIV()("Admin content"),
//			P()("Regular user content"),
//		),
//	)
//
// class := IfElse(err != nil, "text-red", "text-black")
func IfElse[T any](condition bool, result, alternative T) T {
	if condition {
		return result
	}
	return alternative
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
// The provided function is called exactly n times, and each resulting Node
// is aggregated into a single container Node. Using a function ensures each
// Node instance is unique (important for elements with mutable state).
//
// Example:
//
//	UL()(
//		Repeat(5, func() HyperNode {
//			return LI()("List item")
//		}),
//	)
func Repeat(n int, f func() HyperNode) HyperNode {
	result := Element{Tag: "", Children: make([]HyperNode, 0, n)}
	for range n {
		result.Children = append(result.Children, f())
	}
	return result
}

// Range transforms a slice of items into Nodes by applying a function to each element.
//
// Each element in the input slice is transformed using the provided function, and
// all resulting Nodes are aggregated into a single container Node.
//
// Example:
//
//	items := []string{"Apple", "Banana", "Cherry"}
//	UL()(
//		Range(items, func(item string) HyperNode {
//			return LI()(item)
//		}),
//	)
func Range[T any](input []T, f func(T) HyperNode) HyperNode {
	result := Element{Tag: "", Children: make([]HyperNode, 0, len(input))}
	for _, item := range input {
		result.Children = append(result.Children, f(item))
	}
	return result
}

// Group groups multiple children without wrapping them in a tag.
// It creates a container Element with an empty Tag, which renders only its children.
//
// Example:
//
//	Group(P()("Item 1"), H1()("Item 2"), "Item 3")
func Group(children ...any) HyperNode {
	element := Element{Tag: "", Children: make([]HyperNode, 0, len(children))}
	InsertChildren(&element, children...)
	return element
}

// Cache renders a Node once and caches the output for subsequent renders.
//
// The first time the cached node is rendered, the underlying node is rendered
// and its output is stored in the cache. Subsequent renders simply replay
// the cached output without re-rendering the node. This is useful for expensive
// static content like headers, footers, or complex components that don't change.
//
// A NodeCache must be provided to hold the cached output. The same cache instance
// should be used across all renders for the same content.
//
// Example:
//
//	var headerCache hyper.NodeCache
//
//	// Cache the header - renders once on first access
//	header := hyper.Cache(&headerCache,
//		DIV()(
//			NAV()(
//				A(AttrHref("/"))("Home"),
//				A(AttrHref("/about"))("About"),
//			),
//		),
//	)
//
//	// Render the same cached node multiple times - only renders once
//	BODY()(header, mainContent, header)
func Cache(cache *NodeCache, node HyperNode) HyperNode {
	return cachedNode{cache: cache, node: node}
}

// NodeCache stores rendered output for a cached node.
// It must be reused across renders for the same content.
type NodeCache struct {
	buffer bytes.Buffer
	once   sync.Once
	err    error
}

// cachedNode wraps a node with a cache to render once and replay on subsequent renders.
type cachedNode struct {
	cache *NodeCache
	node  HyperNode
}

func (me cachedNode) Render(w io.Writer) error {
	me.cache.once.Do(func() {
		me.cache.err = me.node.Render(&me.cache.buffer)
	})
	if me.cache.err != nil {
		return me.cache.err
	}
	_, err := w.Write(me.cache.buffer.Bytes())
	return err
}

package h

// IfElse returns the appropriate value based on a boolean condition.
//
// This generic function is useful for inline conditional expressions in
// builder-style code where you need to choose between two values without
// breaking the chain of method calls.
//
// Example:
//
//	div := DIV(KV{"class": IfElse(isActive, "active", "inactive")})
//
//	Body(
//		IfElse(isAdmin,
//			DIV("Admin content"),
//			P("Regular user content"),
//		),
//	)
func IfElse[T any](condition bool, result, alternative T) T {
	if condition {
		return result
	}
	return alternative
}

// Conditionally returns a Node based on a boolean condition.
//
// This function returns an empty Node (not nil) when the
// condition is false, which prevents nil pointer issues when building
// DOM trees.
//
// Example:
//
//	Body(
//		If(showHeader, HEADER(...)),
//		MAIN(...),
//	)
func If(condition bool, result HyperNode) HyperNode {
	if condition {
		return result
	}
	return EMPTY()
}

// Repeat generates multiple Nodes by calling a function n times.
//
// The provided function is called exactly n times, and each resulting Node
// is aggregated into a single container Node. Using a function ensures each
// Node instance is unique (important for elements with mutable state).
//
// Example:
//
//	Ul(
//		Repeat(5, func() HyperNode {
//			return LI("List item")
//		}),
//	)
func Repeat(n int, f func() HyperNode) HyperNode {
	result := newElem("")
	result.Children = make([]HyperNode, 0, n)
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
//	UL(
//		Range(items, func(item string) HyperNode {
//			return LI(item)
//		}),
//	)
func Range[T any](input []T, f func(T) HyperNode) HyperNode {
	result := newElem("")
	result.Children = make([]HyperNode, 0, len(input))
	for _, item := range input {
		result.Children = append(result.Children, f(item))
	}
	return result
}

package h

// IfElse returns the appropriate value based on a boolean condition.
//
// This generic function is useful for inline conditional expressions in
// builder-style code where you need to choose between two values without
// breaking the chain of method calls.
//
// Example:
//
//	div := Div(KV{"class": IfElse(isActive, "active", "inactive")})
//
//	Body(
//		IfElse(isAdmin,
//			Div("Admin content"),
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
//		If(showHeader, Header(...)),
//		Main(...),
//	)
func If(condition bool, result Node) Node {
	if condition {
		return result
	}
	return Empty()
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
//		Repeat(5, func() Node {
//			return Li("List item")
//		}),
//	)
func Repeat(n int, f func() Node) Node {
	result := newElem("")
	for range n {
		result.Children = append(result.Children, f())
	}
	return result
}

// MapSlice transforms a slice of items into Nodes by applying a function to each element.
//
// Each element in the input slice is transformed using the provided function, and
// all resulting Nodes are aggregated into a single container Node.
//
// Example:
//
//	items := []string{"Apple", "Banana", "Cherry"}
//	Ul(
//		MapSlice(items, func(item string) Node {
//			return Li(item)
//		}),
//	)
func MapSlice[T any](input []T, f func(T) Node) Node {
	result := newElem("")
	for _, item := range input {
		result.Children = append(result.Children, f(item))
	}
	return result
}

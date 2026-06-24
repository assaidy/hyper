# hyper

A fast, type-safe HTML generator for Go.

## Features

- **Auto-escaping** - Strings are HTML-escaped automatically for security
- **Type-safe** - Compile-time checking of your HTML structure
- **Zero dependencies** - Pure Go standard library
- **Fast** - Minimal allocations, direct writer output
- **Composable** - Build complex layouts from simple components

## Installation

```bash
go get github.com/assaidy/hyper/v2
```

## Quick Start

```go
package main

import (
    "os"

    . "github.com/assaidy/hyper/v2"
)

func main() {
    users := []string{"Alice", "Bob", "Charlie"}
    isAdmin := true

    page := Group(
        DOCTYPE(),
        HTML(AttrLang("en"))(
            HEAD()(
                TITLE()("Dashboard"),
                SCRIPT(AttrSrc("https://cdn.tailwindcss.com")),
            ),
            BODY(AttrClass("bg-gray-100 p-8"))(
                DIV(AttrClass("max-w-2xl mx-auto"))(
                    H1(AttrClass("text-3xl font-bold mb-4"))("Dashboard"),
                    
                    // Conditional admin panel
                    If(isAdmin,
                        DIV(AttrClass("bg-blue-50 p-4 rounded mb-4"))(
                            P(AttrClass("font-semibold"))("Admin Panel"),
                        ),
                    ).ElseIf(isTrial,
                        DIV(AttrClass("bg-yellow-50 p-4 rounded mb-4"))(
                            P(AttrClass("font-semibold"))("Try Premium"),
                        ),
                    ).Else(
                        DIV(AttrClass("p-4"))(
                            P()("Welcome Guest"),
                        ),
                    ),
                    
                    // User count
                    P()("Total users: ", len(users)),

                    // Standard form submission to refresh users
                    FORM(AttrMethod(MethodPost), AttrAction("/api/users/refresh"))(
                        BUTTON(
                            AttrClass("px-4 py-2 bg-blue-500 text-white rounded mt-4"),
                            AttrType(TypeSubmit),
                        )("Refresh Users"),
                    ),
                    
                    // User list
                    UL(AttrClass("space-y-2 mt-4"), AttrId("users-list"))(
                        Range(users, func(name string) any {
                            return LI(AttrClass("p-2 bg-white rounded shadow"))(name)
                        }),
                    ),
                ),
            ),
        ),
    )

    if err := page.Render(os.Stdout); err != nil {
        panic(err)
    }
}
```

## API Explained

All HTML element functions use **ALL_CAPS** names (DIV, P, H1, etc.) to avoid conflicts with standard library types and functions.

### Element Constructor Pattern

Elements follow a consistent pattern for attribute and child handling:

```go
ELEMENT(attrs)(children)
```

When an element has no children, the trailing `()` can be omitted:

```go
// No attributes, no children
DIV()
DIV()() // keeping trailing ()

// With attributes, no children  
DIV(AttrClass("container"))

// No attributes, with children
DIV()("Hello")

// With attributes and children
DIV(AttrClass("container"))("Hello")
```

Void elements (self-closing tags like `<br>`, `<img>`, `<input>`) cannot have children, so they only have one set of parentheses:

```go
BR()                              // <br>
IMG(AttrSrc("image.jpg"))        // <img src="image.jpg">
INPUT(AttrType("text"))          // <input type="text">
```

### Attributes
Use `Attr(key, value)` for any attribute:
- `key` is a string (the attribute name)
- `value` can be a `string` (creates `key="value"`) or `bool` (creates `key` when true, omitted when false)

Alternatively, use specific attribute functions for common attributes:

```go
// Using Attr function
DIV(Attr("class", "container"), Attr("id", "main"))("Content")
// <div class="container" id="main">Content</div>

// Using specific attribute functions (recommended for clarity)
DIV(AttrClass("container"), AttrId("main"))("Content")
// <div class="container" id="main">Content</div>

// Boolean attributes (present when true, absent when false)
INPUT(AttrDisabled(true))   // <input disabled>
INPUT(AttrDisabled(false))  // <input>
```

`nil` attributes are silently skipped during rendering. This is useful for conditional attributes with `IfElseZero`:

```go
DIV(
    AttrClass(Classes("btn", IfElseZero(isPrimary, "btn-primary"))),
    IfElseZero(isHidden, AttrHidden(true)),
)("Content")
// When isHidden=false, no hidden attribute is rendered
```

Custom attributes can be created with `MakePairAttribute` and `MakeBooleanAttribute`:

```go
// Pair attribute — takes a string value
var AttrHxGet = MakePairAttribute("hx-get")
var AttrHxPost = MakePairAttribute("hx-post")
var AttrHxTarget = MakePairAttribute("hx-target")

// Boolean attribute — takes a bool value
var AttrHxPreserve = MakeBooleanAttribute("hx-preserve")

BUTTON(
    AttrHxGet("/api/data"),
    AttrHxTarget("#result"),
    AttrHxPreserve(true),
)("Click me")
// Renders: <button hx-get="/api/data" hx-target="#result" hx-preserve>Click me</button>
```

### Children
The second set of parentheses accepts children. It accepts `HyperNode` values, strings (converted to `Text`), and other values (converted to `Text` via fmt.Sprint).

When no children are needed, the element function alone (`DIV()`, `P()`, etc.) already implements `HyperNode`, so the trailing `()` can be dropped. This is what enables the simplified patterns above.

> **Performance note:** On hot paths, keeping the trailing `()` is slightly faster because it resolves to an `Element` directly, avoiding the closure call in `ChildrenInserter.Render`. Measure both to decide what matters for your use case.

```go
// Strings are auto-escaped
DIV()("Hello", " ", "World")  // <div>Hello World</div>

P()("<script>alert('xss')</script>")
// <p>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</p>

// Raw HTML (not escaped. use with caution)
DIV()(RawText("<svg>...</svg>")) // <svg>...</svg>

// Numbers and booleans are auto-converted
P()("Count: ", 42)           // <p>Count: 42</p>
P()("Active: ", true)        // <p>Active: true</p>
```

### Conditional Rendering

```go
// Show element only if condition is true
If(isLoggedIn, DIV()("Welcome back!"))

// Choose between two options
IfElse(isAdmin, DIV()("Admin"), DIV()("User"))
```

### Lists and Iteration

```go
items := []string{"Apple", "Banana"}

// Map over slice
UL()(
    Range(items, func(item string) any {
        return LI()(item)
    }),
)

// Repeat N times
DIV()(
    Repeat(3, func() any {
        return P()("Repeated")
    }),
)
```

### Grouping
Use `Group()` to group multiple children without wrapping them in a tag. This is useful for fragments that don't have a common ancestor:

```go
Group(
    H1()("Title"),
    P()("Description"),
)
// Renders: <h1>Title</h1><p>Description</p>
```

### CSS Classes

`Classes()` joins multiple class names into a single space-separated string, filtering out empty and duplicate entries:

```go
BUTTON(
    AttrClass(Classes(
        "btn",
        IfElse(err != nil, "btn-error", "btn-primary"),
        IfElseZero(isHidden, "hidden"),
    )),
)
```

### JSON

`Json()` marshals a value to a JSON string, panicking on error. `Object` is a shorthand for `map[string]any`:

```go
FORM(Attr("hx-vals", Json(Object{"role": "admin", "active": true})))
```

### Caching

`Once()` and `OnceWithKey()` cache the rendered output of a component so the underlying tree is only built and rendered once — even when the tree is reconstructed on every request.

When you rebuild the component tree per request (common in web apps), wrap the expensive static parts with `Once` or `OnceWithKey` to avoid reconstructing and re-rendering them on every call.

- `Once` derives the cache key from the caller's program counter — zero-config, guaranteed uniqueness.
- `OnceWithKey` uses an explicit key — useful when the same component is built from multiple call sites. It's also slightly faster (skips the PC capture).

> **Note:** When `Once` is called inside a loop (a raw `for`, `Repeat()`, or `Range()`), all iterations share the same call site and therefore **the same cache key**. Only the first iteration renders; subsequent ones reuse the cached HTML. Use `OnceWithKey` with a distinguishing value (e.g., the loop index) when each iteration needs its own cache entry.

```go
// Once — auto-keyed by caller PC (recommended)
page := Once(func() HyperNode {
    return Group(
        DOCTYPE(),
        HTML()(
            HEAD()(TITLE()("Dashboard")),
            BODY()(H1()("Welcome")),
        ),
    )
})

// OnceWithKey — explicit key
page := OnceWithKey("dashboard", func() HyperNode {
    return Group(
        DOCTYPE(),
        HTML()(
            HEAD()(TITLE()("Dashboard")),
            BODY()(H1()("Welcome")),
        ),
    )
})
```

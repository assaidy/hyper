# hyper

A fast, type-safe HTML generator for Go.

## Features

- **Auto-escaping** - Strings are HTML-escaped automatically for security
- **Zero runtime dependencies** - Pure Go standard library
- **Fast** - Minimal allocations, direct writer output
- **Composable** - Build complex layouts from simple components

## Installation

```bash
go get github.com/assaidy/hyper/v3
```

## Quick Start

```go
package main

import (
    "os"

    . "github.com/assaidy/hyper/v3"
)

// Custom attributes via constructors.
var (
    AttrHxGet      = MakePairAttributeConstructor("hx-get")
    AttrHxTarget   = MakePairAttributeConstructor("hx-target")
    AttrHxPreserve = MakeBooleanAttributeConstructor("hx-preserve")
)

func main() {
    users := []string{"Alice", "Bob", "Charlie"}
    isAdmin := true
    isTrial := false
    isHidden := false
    hasError := false

    // Static parts are cached: Once is auto-keyed, OnceKey uses an explicit key.
    header := Once(func() HyperNode {
        return HEADER(
            AttrClass("p-4 border-b"),
            H1(AttrClass("text-3xl font-bold"), "Dashboard"),
        )
    })
    footer := OnceKey("footer", func() HyperNode {
        return FOOTER(AttrClass("p-4 text-sm"), "Static footer")
    })

    // Fragment groups nodes without a wrapper; DOCTYPE is a void node.
    page := Fragment(
        DOCTYPE(),
        HTML(
            AttrLang("en"),
            HEAD(
                TITLE("Dashboard"),
                META(AttrCharset("utf-8")), // Void elements take only attributes.
                LINK(AttrRel("stylesheet"), AttrHref("/app.css")),
                SCRIPT(AttrSrc("https://cdn.tailwindcss.com")),
            ),
            BODY(
                AttrClass("bg-gray-100 p-8"),
                header,
                DIV(
                    AttrClass("max-w-2xl mx-auto"),
                    // IfNode/ElseIf/Else for conditional nodes.
                    IfNode(isAdmin,
                        DIV(
                            AttrClass("bg-blue-50 p-4 rounded mb-4"),
                            P(AttrClass("font-semibold"), "Admin Panel"),
                        ),
                    ).ElseIf(isTrial,
                        DIV(
                            AttrClass("bg-yellow-50 p-4 rounded mb-4"),
                            P(AttrClass("font-semibold"), "Try Premium"),
                        ),
                    ).Else(
                        DIV(AttrClass("p-4"), P("Welcome Guest")),
                    ),

                    // If/Else/Value for plain values; Classes joins conditional classes.
                    P(
                        AttrClass(Classes(
                            "text-sm",
                            If(hasError, "text-red-500").Else("text-gray-500"),
                            If(isHidden, "hidden").Value(),
                        )),
                        "Total users: ", len(users), // Ints auto-convert to Text.
                    ),

                    // Strings escape by default; RawText does not.
                    P("<b>escaped</b>"),
                    DIV(RawText("<hr>")),

                    // Void elements; Attr works for any key; IfAttr omits when false.
                    IMG(AttrSrc("/logo.png"), AttrAlt("logo")),
                    BR(),
                    INPUT(
                        AttrType(TypeSubmit),
                        AttrDisabled(!isAdmin), // Boolean attr: present when true.
                        Attr("placeholder", "Search"),
                        IfAttr(isHidden, AttrHidden(true)),
                    ),

                    // Json/Object for inline JSON; custom attrs from constructors.
                    FORM(
                        AttrMethod(MethodPost),
                        AttrAction("/api/users/refresh"),
                        Attr("hx-vals", Json(Object{"role": "admin", "active": true})),
                        BUTTON(
                            AttrClass("px-4 py-2 bg-blue-500 text-white rounded mt-4"),
                            AttrType(TypeSubmit),
                            // Attrs we defined before.
                            AttrHxGet("/api/users"),
                            AttrHxTarget("#users-list"),
                            AttrHxPreserve(true),
                            "Refresh Users",
                        ),
                    ),

                    // Range maps a slice; Repeat repeats N times (both pass the index).
                    UL(
                        AttrClass("space-y-2 mt-4"),
                        AttrId("users-list"),
                        Range(users, func(_ int, name string) any {
                            return LI(AttrClass("p-2 bg-white rounded shadow"), name)
                        }),
                        Repeat(2, func(_ int) any {
                            return LI(AttrClass("p-2 opacity-50"), "Placeholder")
                        }),
                    ),
                ),
                footer,
            ),
        ),
    )

    // RenderNode writes to any io.Writer.
    if err := RenderNode(os.Stdout, page); err != nil {
        panic(err)
    }
}
```

## API Explained

All HTML element functions use **ALL_CAPS** names (DIV, P, H1, etc.) to avoid conflicts with standard library types and functions.

### Element Constructor Pattern

Elements take attributes and children in a single call:

```go
ELEMENT(attr1, attr2, child1, child2, ...)
```

```go
// No attributes, no children
DIV()

// With attributes, no children
DIV(AttrClass("container"))

// No attributes, with children
DIV("Hello")

// With attributes and children
DIV(AttrClass("container"), "Hello")
```

Void elements (self-closing tags like `<br>`, `<img>`, `<input>`) cannot have children, so they only accept attributes:

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
DIV(Attr("class", "container"), Attr("id", "main"), "Content")
// <div class="container" id="main">Content</div>

// Using specific attribute functions (recommended for clarity)
DIV(AttrClass("container"), AttrId("main"), "Content")
// <div class="container" id="main">Content</div>

// Boolean attributes (present when true, absent when false)
INPUT(AttrDisabled(true))   // <input disabled>
INPUT(AttrDisabled(false))  // <input>
```

Passing a `nil` argument to an element constructor panics. For conditional attributes use `IfAttr`, which renders nothing when no branch matches and no else branch was set (see [IfAttr](#ifattr)):

```go
DIV(
    AttrClass(Classes("btn", If(isPrimary, "btn-primary").Value())),
    IfAttr(isHidden, AttrHidden(true)),
    "Content",
)
// When isHidden=false, no hidden attribute is rendered
```

Custom attributes can be created with `MakePairAttributeConstructor` and `MakeBooleanAttributeConstructor`:

```go
// Pair attribute — takes a string value
var AttrHxGet = MakePairAttributeConstructor("hx-get")
var AttrHxPost = MakePairAttributeConstructor("hx-post")
var AttrHxTarget = MakePairAttributeConstructor("hx-target")

// Boolean attribute — takes a bool value
var AttrHxPreserve = MakeBooleanAttributeConstructor("hx-preserve")

BUTTON(
    AttrHxGet("/api/data"),
    AttrHxTarget("#result"),
    AttrHxPreserve(true),
    "Click me",
)
// Renders: <button hx-get="/api/data" hx-target="#result" hx-preserve>Click me</button>
```

### Children
Element constructors accept children after the attributes. Children can be `HyperNode` values, strings (converted to `Text`), and other values (converted to `Text` via fmt.Sprint).

An element with no children (e.g. `DIV()`, `P()`) already implements `HyperNode`, so it can be used directly anywhere a node is expected.

```go
// Strings are auto-escaped
DIV("Hello", " ", "World")  // <div>Hello World</div>

P("<script>alert('xss')</script>")
// <p>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</p>

// Raw HTML (not escaped. use with caution)
DIV(RawText("<svg>...</svg>")) // <svg>...</svg>

// Numbers and booleans are auto-converted
P("Count: ", 42)           // <p>Count: 42</p>
P("Active: ", true)        // <p>Active: true</p>
```

### Utils

Helpers for conditionals, lists, grouping, classes, JSON, and caching.

#### If

`If` is a generic if expression for plain values (strings, numbers, structs). Chain `ElseIf` calls and terminate with `Else` (explicit fallback) or `Value` (zero-value fallback).

Do not use `If` with nodes or attributes — use [IfNode](#ifnode) and [IfAttr](#ifattr) instead, as they are render/nil-safe.

```go
// Two-way choice
If(hasError, "btn-error").Else("btn-primary") // "btn-error" or "btn-primary"

// Single-branch (falls back to "" when false — dropped by Classes)
If(isHidden, "hidden").Value()

// Multi-branch
If(isAdmin, "admin").ElseIf(isMember, "member").Else("guest")

// Typical use with Classes
BUTTON(
    AttrClass(Classes(
        "btn",
        If(hasError, "btn-error").Else("btn-primary"),
        If(isHidden, "hidden").Value(),
    )),
)
```

#### IfAttr

`IfAttr` starts a conditional attribute chain. The attribute renders only when its condition is true; when no branch matches and no else was set, nothing renders (the attribute is omitted). Chain `ElseIf` and terminate with `Else`. Passing `nil` to `IfAttr`/`ElseIf`/`Else` panics.

```go
// Rendered only when isHidden
IfAttr(isHidden, AttrHidden(true))

// With fallback
IfAttr(isAdmin, AttrRequired(true)).
    Else(AttrReadOnly(true))

// Multi-branch
IfAttr(isDark, AttrHidden(true)).
    ElseIf(isLight, AttrRequired(true))
```

#### IfNode

`IfNode` starts a conditional node chain. The body renders only when its condition is true, otherwise an empty fragment renders (avoiding nil panics). Without an `Else`, a fully-false chain renders nothing. Chain `ElseIf` and terminate with `Else`. Passing `nil` to `IfNode`/`ElseIf`/`Else` panics.

```go
// Show only when logged in (renders "" otherwise)
IfNode(isLoggedIn, DIV("Welcome back!"))

// Two-way choice
IfNode(isAdmin, DIV("Admin")).Else(DIV("User"))

// Multi-branch
IfNode(isAdmin,
    DIV(
        AttrClass("bg-blue-50 p-4 rounded mb-4"),
        P(AttrClass("font-semibold"), "Admin Panel"),
    ),
).ElseIf(isTrial,
    DIV(
        AttrClass("bg-yellow-50 p-4 rounded mb-4"),
        P(AttrClass("font-semibold"), "Try Premium"),
    ),
).Else(
    DIV(
        AttrClass("p-4"),
        P("Welcome Guest"),
    ),
)
```

#### Repeat

`Repeat` calls `generate` exactly `n` times with the zero-based index and aggregates each result (converted to a node) into a single container node. Using a function keeps each instance unique. Zero or negative renders nothing without calling `generate`. Ignore the index with `_` when it is not needed.

```go
UL(
    Repeat(5, func(i int) any {
        return LI(i)
    }),
)
```

#### Range

`Range` transforms each slice element with `generate` (receiving the index and the element) and aggregates the results (converted to nodes) into a single container node. Nil and empty inputs render nothing without calling `generate`. Ignore the index with `_` when it is not needed.

```go
items := []string{"Apple", "Banana", "Cherry"}
UL(
    Range(items, func(i int, item string) any {
        return LI(item)
    }),
)
```

#### Fragment

`Fragment` groups multiple nodes into a single node without a wrapping tag; children render back-to-back, identical to rendering each node individually. Useful for passing several nodes as one value (e.g. a conditional body, or a single argument to another component). Empty renders nothing; `nil` arguments panic.

```go
Fragment(
    H1("Title"),
    P("Description"),
    "Plain text",
)
// Renders: <h1>Title</h1><p>Description</p>Plain text
```

#### Classes

`Classes` joins class names into a single space-separated string, trimming surrounding whitespace and dropping empty and duplicate entries. With no classes it returns `""`. Pair with `If` (see [If](#if)): terminate single-branch cases with `Value` and two-way cases with `Else`. See [Attributes](#attributes).

```go
Classes("btn", "btn", "", " primary ") // "btn primary"

BUTTON(
    AttrClass(Classes(
        "btn",
        If(hasError, "btn-error").Else("btn-primary"),
        If(isHidden, "hidden").Value(),
    )),
)
```

#### Json and Object

`Json` marshals a value to a JSON string, panicking on error — for static values known valid at compile time. `Object` is a shorthand for `map[string]any` for inline maps.

```go
FORM(Attr("hx-vals", Json(Object{"role": "admin", "active": true})))

Json(make(chan int)) // panics: unsupported type
```

#### Once and OnceKey

`Once` and `OnceKey` cache a component's rendered output so the tree is built and rendered once, even when rebuilt per request. The first render calls `generate` and caches; later renders replay the cache. A failed render is not cached.

- `Once` derives the key from the caller's program counter — zero-config and unique per call site.
- `OnceKey` uses an explicit key — for shared components and slightly faster (skips PC capture). Keys must be unique; same-key calls share the entry.

> **Note:** When `Once` is called inside a loop (a raw `for`, `Repeat()`, or `Range()`), all iterations share the same call site and therefore **the same cache key**. Only the first iteration renders; subsequent ones reuse the cached HTML. Use `OnceKey` with a distinguishing value (e.g., the loop index) when each iteration needs its own cache entry.

```go
// Once — auto-keyed by caller PC (recommended)
page := Once(func() HyperNode {
    return Fragment(
        DOCTYPE(),
        HTML(
            HEAD(TITLE("Dashboard")),
            BODY(H1("Welcome")),
        ),
    )
})

// OnceKey — explicit key
page := OnceKey("dashboard", func() HyperNode {
    return Fragment(
        DOCTYPE(),
        HTML(
            HEAD(TITLE("Dashboard")),
            BODY(H1("Welcome")),
        ),
    )
})
```

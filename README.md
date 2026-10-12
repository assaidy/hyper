# hyper

A fast, type-safe HTML generator for Go.

## Features

- **Auto-escaping** - Strings are HTML-escaped automatically for security
- **Zero runtime dependencies** - Pure Go standard library
- **No codegen, no CLI, no build step** - Just Go: `go get` and write components, nothing to install or generate
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

## Benchmarks

### Scenarios

All three engines render **byte-identical HTML** (enforced by `TestParity`):

| Scenario | What it measures |
|---|---|
| `Simple` | `<div class="greeting">Hello, {name}</div>`, name needs escaping |
| `Dashboard` | Conditionals, conditional classes, boolean attr, `Range` over 3 users, count int |
| `List100` / `List1000` | `<ul>` with 100 / 1000 `<li>` rows — loop throughput |
| `Escape` | `<p>` with `<script>` + quotes payload — escaping cost |
| `Once` (fresh) | Static page rendered fresh: hyper built once, `Template.Execute`, `Templ.Render` |
| `Once` (cached) | Same page replayed from cache: hyper `Once`/`OnceKey`, pre-rendered bytes for template/templ |

Hyper reports two variants: **render-only** (tree built once outside the
loop — the expected per-request pattern) and **build+render** (tree rebuilt
every iteration, for transparency).

### Run

```bash
cd benchmarks
go test -bench=. -benchmem -benchtime=1s -count=3 .
```

### Regenerate templ code

`components_templ.go` is checked in so no tooling is needed to run the
benchmarks. To regenerate after editing `components.templ`, use the pinned
toolchain:

```bash
# pinned: templ v0.3.1070
go install github.com/a-h/templ/cmd/templ@v0.3.1070
cd benchmarks
templ generate
```
### Results

(`go1.27.2 linux/amd64`, i7-1165G7, `templ v0.3.1070`).
Time is ns/op medians of the 3 runs below, with allocs/op in brackets:

| Scenario | hyper | `html/template` | `templ` |
|---|---|---|---|
| Simple | ~225ns (1) | ~1100ns (9) | ~390ns (6) |
| Dashboard | ~1.35µs (0) | ~5.0µs (30) | ~1.7µs (16) |
| List100 | ~14.2µs (0) | ~83µs (603) | ~12.8µs (203) |
| List1000 | ~142µs (0) | ~850µs (6747) | ~132µs (2003) |
| Escape | ~332ns (2) | ~1.3µs (9) | ~630ns (7) |
| Static page | fresh ~535ns / `Once` ~23ns / `OnceKey` ~21ns (0) | ~117ns (2) | ~195ns (3) |

Full raw output of `go test -bench=. -benchmem -benchtime=1s -count=3 .`:

```
goos: linux
goarch: amd64
pkg: github.com/assaidy/hyper/v3/benchmarks
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkSimple_Hyper-8                 	 5451915	       211.8 ns/op	      24 B/op	       1 allocs/op
BenchmarkSimple_Hyper-8                 	 5418240	       228.4 ns/op	      24 B/op	       1 allocs/op
BenchmarkSimple_Hyper-8                 	 5031178	       230.1 ns/op	      24 B/op	       1 allocs/op
BenchmarkSimple_HyperBuildRender-8      	 2142985	       561.5 ns/op	     248 B/op	       9 allocs/op
BenchmarkSimple_HyperBuildRender-8      	 2155891	       571.9 ns/op	     248 B/op	       9 allocs/op
BenchmarkSimple_HyperBuildRender-8      	 1988871	       603.5 ns/op	     248 B/op	       9 allocs/op
BenchmarkSimple_Template-8              	  947815	      1104 ns/op	     288 B/op	       9 allocs/op
BenchmarkSimple_Template-8              	  916033	      1098 ns/op	     288 B/op	       9 allocs/op
BenchmarkSimple_Template-8              	 1000000	      2063 ns/op	     288 B/op	       9 allocs/op
BenchmarkSimple_Templ-8                 	 2288954	       621.3 ns/op	     200 B/op	       6 allocs/op
BenchmarkSimple_Templ-8                 	 3221280	       388.1 ns/op	     200 B/op	       6 allocs/op
BenchmarkSimple_Templ-8                 	 3144090	       381.7 ns/op	     200 B/op	       6 allocs/op
BenchmarkDashboard_Hyper-8              	  821268	      1352 ns/op	       0 B/op	       0 allocs/op
BenchmarkDashboard_Hyper-8              	  909792	      1372 ns/op	       0 B/op	       0 allocs/op
BenchmarkDashboard_Hyper-8              	  908710	      1275 ns/op	       0 B/op	       0 allocs/op
BenchmarkDashboard_HyperBuildRender-8   	  243452	      4841 ns/op	    2250 B/op	      67 allocs/op
BenchmarkDashboard_HyperBuildRender-8   	  262700	      4792 ns/op	    2250 B/op	      67 allocs/op
BenchmarkDashboard_HyperBuildRender-8   	  261814	      4765 ns/op	    2250 B/op	      67 allocs/op
BenchmarkDashboard_Template-8           	  239091	      5345 ns/op	     752 B/op	      30 allocs/op
BenchmarkDashboard_Template-8           	  224031	      5039 ns/op	     752 B/op	      30 allocs/op
BenchmarkDashboard_Template-8           	  238149	      5016 ns/op	     752 B/op	      30 allocs/op
BenchmarkDashboard_Templ-8              	  697069	      1697 ns/op	     592 B/op	      16 allocs/op
BenchmarkDashboard_Templ-8              	  741787	      1697 ns/op	     592 B/op	      16 allocs/op
BenchmarkDashboard_Templ-8              	  768928	      1643 ns/op	     592 B/op	      16 allocs/op
BenchmarkList100_Hyper-8                	   87747	     14954 ns/op	       0 B/op	       0 allocs/op
BenchmarkList100_Hyper-8                	   86704	     14234 ns/op	       0 B/op	       0 allocs/op
BenchmarkList100_Hyper-8                	   75361	     13845 ns/op	       0 B/op	       0 allocs/op
BenchmarkList100_HyperBuildRender-8     	   31137	     39112 ns/op	   19667 B/op	     608 allocs/op
BenchmarkList100_HyperBuildRender-8     	   31360	     38956 ns/op	   19662 B/op	     608 allocs/op
BenchmarkList100_HyperBuildRender-8     	   30390	     40675 ns/op	   19665 B/op	     608 allocs/op
BenchmarkList100_Template-8             	   14394	     83412 ns/op	   12965 B/op	     603 allocs/op
BenchmarkList100_Template-8             	   14774	     82823 ns/op	   12965 B/op	     603 allocs/op
BenchmarkList100_Template-8             	   14146	     83351 ns/op	   12965 B/op	     603 allocs/op
BenchmarkList100_Templ-8                	   95119	     12058 ns/op	    2545 B/op	     203 allocs/op
BenchmarkList100_Templ-8                	   83366	     12831 ns/op	    2545 B/op	     203 allocs/op
BenchmarkList100_Templ-8                	   94368	     13039 ns/op	    2545 B/op	     203 allocs/op
BenchmarkList1000_Hyper-8               	    8659	    143599 ns/op	       6 B/op	       0 allocs/op
BenchmarkList1000_Hyper-8               	    7755	    141711 ns/op	       7 B/op	       0 allocs/op
BenchmarkList1000_Hyper-8               	    7660	    140776 ns/op	       7 B/op	       0 allocs/op
BenchmarkList1000_HyperBuildRender-8    	    2643	    389185 ns/op	  195513 B/op	    6008 allocs/op
BenchmarkList1000_HyperBuildRender-8    	    3106	    378893 ns/op	  195413 B/op	    6008 allocs/op
BenchmarkList1000_HyperBuildRender-8    	    2889	    387311 ns/op	  195415 B/op	    6008 allocs/op
BenchmarkList1000_Template-8            	    1395	    849805 ns/op	  134167 B/op	    6747 allocs/op
BenchmarkList1000_Template-8            	    1420	    833930 ns/op	  134168 B/op	    6747 allocs/op
BenchmarkList1000_Template-8            	    1398	    873191 ns/op	  134168 B/op	    6747 allocs/op
BenchmarkList1000_Templ-8               	    8568	    130918 ns/op	   24169 B/op	    2003 allocs/op
BenchmarkList1000_Templ-8               	    7797	    135787 ns/op	   24169 B/op	    2003 allocs/op
BenchmarkList1000_Templ-8               	    9382	    131919 ns/op	   24169 B/op	    2003 allocs/op
BenchmarkEscape_Hyper-8                 	 3549889	       334.4 ns/op	     192 B/op	       2 allocs/op
BenchmarkEscape_Hyper-8                 	 3716401	       332.0 ns/op	     192 B/op	       2 allocs/op
BenchmarkEscape_Hyper-8                 	 3485166	       322.5 ns/op	     192 B/op	       2 allocs/op
BenchmarkEscape_HyperBuildRender-8      	 2682072	       852.7 ns/op	     240 B/op	       5 allocs/op
BenchmarkEscape_HyperBuildRender-8      	 2308149	       518.1 ns/op	     240 B/op	       5 allocs/op
BenchmarkEscape_HyperBuildRender-8      	 2679180	       437.8 ns/op	     240 B/op	       5 allocs/op
BenchmarkEscape_Template-8              	  921966	      1314 ns/op	     432 B/op	       9 allocs/op
BenchmarkEscape_Template-8              	  830755	      1340 ns/op	     432 B/op	       9 allocs/op
BenchmarkEscape_Template-8              	  741736	      1398 ns/op	     432 B/op	       9 allocs/op
BenchmarkEscape_Templ-8                 	 1832640	       631.5 ns/op	     416 B/op	       7 allocs/op
BenchmarkEscape_Templ-8                 	 2012838	       603.8 ns/op	     416 B/op	       7 allocs/op
BenchmarkEscape_Templ-8                 	 1962873	       636.4 ns/op	     416 B/op	       7 allocs/op
BenchmarkOnce_HyperNoCache-8            	 2143704	       533.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperNoCache-8            	 2138282	       543.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperNoCache-8            	 2142007	       534.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperCached-8             	55774915	        20.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperCached-8             	55502074	        21.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperCached-8             	55982004	        21.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperOnceCached-8         	56760690	        22.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperOnceCached-8         	53884869	        22.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperOnceCached-8         	51631347	        24.29 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperOnceBuild-8          	 4508143	       256.7 ns/op	      32 B/op	       2 allocs/op
BenchmarkOnce_HyperOnceBuild-8          	 4364444	       271.7 ns/op	      32 B/op	       2 allocs/op
BenchmarkOnce_HyperOnceBuild-8          	 4237484	       275.7 ns/op	      32 B/op	       2 allocs/op
BenchmarkOnce_HyperOnceKeyBuild-8       	43479802	        23.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperOnceKeyBuild-8       	46906054	        23.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_HyperOnceKeyBuild-8       	50845287	        21.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_TemplateExecute-8         	10504894	       114.8 ns/op	     128 B/op	       2 allocs/op
BenchmarkOnce_TemplateExecute-8         	 9147585	       120.2 ns/op	     128 B/op	       2 allocs/op
BenchmarkOnce_TemplateExecute-8         	 9732920	       116.7 ns/op	     128 B/op	       2 allocs/op
BenchmarkOnce_TemplateCachedBytes-8     	609283545	         2.034 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_TemplateCachedBytes-8     	615087121	         1.913 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_TemplateCachedBytes-8     	612976748	         1.920 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_TemplRender-8             	 6476883	       196.1 ns/op	     144 B/op	       3 allocs/op
BenchmarkOnce_TemplRender-8             	 6054528	       192.5 ns/op	     144 B/op	       3 allocs/op
BenchmarkOnce_TemplRender-8             	 6334855	       197.3 ns/op	     144 B/op	       3 allocs/op
BenchmarkOnce_TemplCachedBytes-8        	577302494	         1.938 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_TemplCachedBytes-8        	656310472	         1.958 ns/op	       0 B/op	       0 allocs/op
BenchmarkOnce_TemplCachedBytes-8        	574639908	         2.104 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/assaidy/hyper/v3/benchmarks	127.031s
```

Takeaways:

- **vs `html/template`**: hyper is several times faster on every scenario with
  real work (roughly 4-6x on this machine), with zero allocations on
  dashboard/list renders vs thousands. The one exception is a trivial
  fully-static page, where `html/template` wins — hyper's fixed per-element
  overhead dominates when there is nothing to do. Wrap statics in `OnceKey`
  (~21ns replay vs ~530ns fresh).
- **vs `templ`**: hyper is faster on small and branchy pages with fewer
  allocations. `templ` is faster on the 100-row list and the static page. The
  1000-row list is about tied in time — but `templ` allocates ~24KB / 2000
  allocs there vs ~0 for hyper.
- **Build cost**: even rebuilding the hyper tree on every request beats
  `html/template` on every scenario with real work; Dashboard is about tied.
- **`OnceKey` vs `Once`**: cached replay is identical (~21-23ns). Per-request
  construction is about 12x cheaper with `OnceKey` (~24ns / 0 allocs vs ~270ns
  / 2 allocs) — it skips caller-PC capture. Prefer `OnceKey` in hot paths.

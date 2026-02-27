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
go get github.com/assaidy/hyper
```

## Quick Start

```go
package main

import (
    "os"

    "github.com/assaidy/hyper"
)

func main() {
    page := h.Empty(
        h.DoctypeHtml(),
        h.Html(
            h.Head(
                h.Title("My Page"),
            ),
            h.Body(
                h.H1("Hello, World!"),
                h.P("Auto-escaped: <script>alert('xss')</script>"),
            ),
        ),
    )
    
    if err := page.Render(os.Stdout); err != nil {
        panic(err)
    }
}
```

## Usage

### Basic Elements

```go
// Strings are auto-escaped
h.Div("Hello", " ", "World")  // <div>Hello World</div>

h.P("<script>alert('xss')</script>")
// <p>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</p>

// Raw HTML (not escaped. use with caution)
h.Div(h.RawText("<svg>...</svg>")) // <svg>...</svg>

// Numbers and booleans are auto-converted
h.P("Count: ", 42)           // <p>Count: 42</p>
h.P("Active: ", true)        // <p>Active: true</p>
```

### Attributes

```go
h.Div(h.KV{h.AttrClass: "container", h.AttrId: "main"}, "Content")
// <div class="container" id="main">Content</div>
```

### Conditional Rendering

```go
// Show element only if condition is true
h.If(isLoggedIn, h.Div("Welcome back!"))

// Choose between two options
h.IfElse(isAdmin, h.Div("Admin"), h.Div("User"))
```

### Lists and Iteration

```go
items := []string{"Apple", "Banana"}

// Map over slice
h.Ul(
    h.MapSlice(items, func(item string) h.Node {
        return h.Li(item)
    }),
)

// Repeat N times
h.Div(
    h.Repeat(3, func() h.Node {
        return h.P("Repeated")
    }),
)
```

### With Tailwind CSS

```go
h.Div(h.KV{h.AttrClass: "bg-gray-100 min-h-screen p-8"},
    h.Div(h.KV{h.AttrClass: "max-w-4xl mx-auto"},
        h.H1(h.KV{h.AttrClass: "text-4xl font-bold text-gray-800"}, "Title"),
        h.P(h.KV{h.AttrClass: "text-gray-600 mt-2"}, "Description"),
        h.Button(h.KV{h.AttrClass: "px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"},
            "Click Me",
        ),
    ),
)
```

### With HTMX

```go
// Include htmx library from CDN
h.Script(h.KV{h.AttrSrc: "https://unpkg.com/htmx.org@2.0.8"})

// HTMX button that loads content
h.Button(h.KV{
    h.AttrClass:     "px-4 py-2 bg-blue-500 text-white rounded",
    hx.AttrHxGet:      "/api/users",
    hx.AttrHxTarget:   "#users-list",
    hx.AttrHxSwap:     hx.SwapOuterHtml,
},
    "Load Users",
)

// HTMX form
h.Form(h.KV{
    hx.AttrHxPost:     "/api/submit",
    hx.AttrHxTarget:   "#result",
    h.AttrClass:     "space-y-4",
},
    h.Input(h.KV{
        h.AttrType:  h.TypeText,
        h.AttrName:  "message",
        h.AttrClass: "border rounded px-3 py-2",
    }),
    h.Button(h.KV{h.AttrType: h.TypeSubmit, h.AttrClass: "bg-blue-500 text-white px-4 py-2 rounded"},
        "Submit",
    ),
)
```

### Complete Example

```go
package main

import (
    "os"

    "github.com/assaidy/hyper"
    "github.com/assaidy/hyper/htmx"
)

func main() {
    users := []string{"Alice", "Bob", "Charlie"}
    isAdmin := true

    page := h.Empty(
        h.DoctypeHtml(),
        h.Html(h.KV{h.AttrLang: "en"},
            h.Head(
                h.Title("Dashboard"),
                h.Script(h.KV{h.AttrSrc: "https://unpkg.com/htmx.org@2.0.8"}),
                h.Script(h.KV{h.AttrSrc: "https://cdn.tailwindcss.com"}),
            ),
            h.Body(h.KV{h.AttrClass: "bg-gray-100 p-8"},
                h.Div(h.KV{h.AttrClass: "max-w-2xl mx-auto"},
                    h.H1(h.KV{h.AttrClass: "text-3xl font-bold mb-4"}, "Dashboard"),
                    
                    // Conditional admin panel
                    h.If(isAdmin,
                        h.Div(h.KV{h.AttrClass: "bg-blue-50 p-4 rounded mb-4"},
                            h.P(h.KV{h.AttrClass: "font-semibold"}, "Admin Panel"),
                        ),
                    ),
                    
                    // User count
                    h.P("Total users: ", len(users)),

                    // HTMX button to refresh users
                    h.Button(h.KV{
                        h.AttrClass:      "px-4 py-2 bg-blue-500 text-white rounded mt-4",
                        hx.AttrHxGet:      "/api/users",
                        hx.AttrHxTarget:   "#users-list",
                        hx.AttrHxSwap:     hx.SwapOuterHtml,
                    },
                        "Refresh Users",
                    ),
                    
                    // User list
                    h.Ul(h.KV{h.AttrClass: "space-y-2 mt-4", h.AttrId: "users-list"},
                        h.MapSlice(users, func(name string) h.Node {
                            return h.Li(h.KV{h.AttrClass: "p-2 bg-white rounded shadow"},
                                name,
                            )
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

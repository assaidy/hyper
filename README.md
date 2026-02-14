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
        h.DoctypeHTML(),
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
h.Div(h.RawHTML("<svg>...</svg>")) // <svg>...</svg>

// Numbers and booleans are auto-converted
h.P("Count: ", 42)           // <p>Count: 42</p>
h.P("Active: ", true)        // <p>Active: true</p>
```

### Attributes

```go
h.Div(h.KV{"class": "container", "id": "main"}, "Content")
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
h.Div(h.KV{"class": "bg-gray-100 min-h-screen p-8"},
    h.Div(h.KV{"class": "max-w-4xl mx-auto"},
        h.H1(h.KV{"class": "text-4xl font-bold text-gray-800"}, "Title"),
        h.P(h.KV{"class": "text-gray-600 mt-2"}, "Description"),
        h.Button(
            h.KV{"class": "px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"},
            "Click Me",
        ),
    ),
)
```

### With HTMX

```go
// HTMX button that loads content
h.Button(
    h.KV{
        "class":     "px-4 py-2 bg-blue-500 text-white rounded",
        "hx-get":    "/api/users",
        "hx-target": "#users-list",
        "hx-swap":   "outerHTML",
    },
    "Load Users",
)

// HTMX form
h.Form(
    h.KV{
        "hx-post":   "/api/submit",
        "hx-target": "#result",
        "class":     "space-y-4",
    },
    h.Input(h.KV{
        "type":  "text",
        "name":  "message",
        "class": "border rounded px-3 py-2",
    }),
    h.Button(
        h.KV{"type": "submit", "class": "bg-blue-500 text-white px-4 py-2 rounded"},
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
)

func main() {
    users := []string{"Alice", "Bob", "Charlie"}
    isAdmin := true

    page := h.Empty(
        h.DoctypeHTML(),
        h.Html(h.KV{"lang": "en"},
            h.Head(
                h.Title("Dashboard"),
                h.Script(h.KV{"src": "https://cdn.tailwindcss.com"}),
            ),
            h.Body(h.KV{"class": "bg-gray-100 p-8"},
                h.Div(h.KV{"class": "max-w-2xl mx-auto"},
                    h.H1(h.KV{"class": "text-3xl font-bold mb-4"}, "Dashboard"),
                    
                    // Conditional admin panel
                    h.If(isAdmin,
                        h.Div(h.KV{"class": "bg-blue-50 p-4 rounded mb-4"},
                            h.P(h.KV{"class": "font-semibold"}, "Admin Panel"),
                        ),
                    ),
                    
                    // User count
                    h.P("Total users: ", len(users)),
                    
                    // User list
                    h.Ul(h.KV{"class": "space-y-2 mt-4"},
                        h.MapSlice(users, func(name string) h.Node {
                            return h.Li(
                                h.KV{"class": "p-2 bg-white rounded shadow"},
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

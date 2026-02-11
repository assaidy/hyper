# g

A fast, type-safe HTML generator for Go.

## Features

- **Auto-escaping** - Strings are HTML-escaped automatically for security
- **Type-safe** - Compile-time checking of your HTML structure
- **Zero dependencies** - Pure Go standard library
- **Fast** - Minimal allocations, direct writer output
- **Composable** - Build complex layouts from simple components

## Installation

```bash
go get github.com/assaidy/g
```

## Quick Start

```go
package main

import (
    "os"
    "github.com/assaidy/g"
)

func main() {
    page := g.Empty(
        g.DoctypeHTML(),
        g.Html(
            g.Head(
                g.Title("My Page"),
            ),
            g.Body(
                g.H1("Hello, World!"),
                g.P("Auto-escaped: <script>alert('xss')</script>"),
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
g.Div("Hello", " ", "World")  // <div>Hello World</div>

g.P("<script>alert('xss')</script>")
// <p>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</p>

// Raw HTML (not escaped. use with caution)
g.Div(g.RawHTML("<svg>...</svg>")) // <svg>...</svg>

// Numbers and booleans are auto-converted
g.P("Count: ", 42)           // <p>Count: 42</p>
g.P("Active: ", true)        // <p>Active: true</p>
```

### Attributes

```go
g.Div(g.KV{"class": "container", "id": "main"}, "Content")
// <div class="container" id="main">Content</div>
```

### Conditional Rendering

```go
// Show element only if condition is true
g.If(isLoggedIn, g.Div("Welcome back!"))

// Choose between two options
g.IfElse(isAdmin, g.Div("Admin"), g.Div("User"))
```

### Lists and Iteration

```go
items := []string{"Apple", "Banana"}

// Map over slice
g.Ul(
    g.MapSlice(items, func(item string) g.Node {
        return g.Li(item)
    }),
)

// Repeat N times
g.Div(
    g.Repeat(3, func() g.Node {
        return g.P("Repeated")
    }),
)
```

### With Tailwind CSS

```go
g.Div(g.KV{"class": "bg-gray-100 min-h-screen p-8"},
    g.Div(g.KV{"class": "max-w-4xl mx-auto"},
        g.H1(g.KV{"class": "text-4xl font-bold text-gray-800"}, "Title"),
        g.P(g.KV{"class": "text-gray-600 mt-2"}, "Description"),
        g.Button(
            g.KV{"class": "px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"},
            "Click Me",
        ),
    ),
)
```

### With HTMX

```go
// HTMX button that loads content
g.Button(
    g.KV{
        "class":     "px-4 py-2 bg-blue-500 text-white rounded",
        "hx-get":    "/api/users",
        "hx-target": "#users-list",
        "hx-swap":   "outerHTML",
    },
    "Load Users",
)

// HTMX form
g.Form(
    g.KV{
        "hx-post":   "/api/submit",
        "hx-target": "#result",
        "class":     "space-y-4",
    },
    g.Input(g.KV{
        "type":  "text",
        "name":  "message",
        "class": "border rounded px-3 py-2",
    }),
    g.Button(
        g.KV{"type": "submit", "class": "bg-blue-500 text-white px-4 py-2 rounded"},
        "Submit",
    ),
)
```

### Complete Example

```go
package main

import (
    "os"
    "github.com/assaidy/g"
)

func main() {
    users := []string{"Alice", "Bob", "Charlie"}
    isAdmin := true

    page := g.Empty(
        g.DoctypeHTML(),
        g.Html(g.KV{"lang": "en"},
            g.Head(
                g.Title("Dashboard"),
                g.Script(g.KV{"src": "https://cdn.tailwindcss.com"}),
            ),
            g.Body(g.KV{"class": "bg-gray-100 p-8"},
                g.Div(g.KV{"class": "max-w-2xl mx-auto"},
                    g.H1(g.KV{"class": "text-3xl font-bold mb-4"}, "Dashboard"),
                    
                    // Conditional admin panel
                    g.If(isAdmin,
                        g.Div(g.KV{"class": "bg-blue-50 p-4 rounded mb-4"},
                            g.P(g.KV{"class": "font-semibold"}, "Admin Panel"),
                        ),
                    ),
                    
                    // User count
                    g.P("Total users: ", len(users)),
                    
                    // User list
                    g.Ul(g.KV{"class": "space-y-2 mt-4"},
                        g.MapSlice(users, func(name string) g.Node {
                            return g.Li(
                                g.KV{"class": "p-2 bg-white rounded shadow"},
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

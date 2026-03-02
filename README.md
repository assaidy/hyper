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

    . "github.com/assaidy/hyper"
)

func main() {
    page := Empty(
        DoctypeHtml(),
        Html(
            Head(
                Title("My Page"),
            ),
            Body(
                H1("Hello, World!"),
                P("Auto-escaped: <script>alert('xss')</script>"),
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
Div("Hello", " ", "World")  // <div>Hello World</div>

P("<script>alert('xss')</script>")
// <p>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</p>

// Raw HTML (not escaped. use with caution)
Div(RawText("<svg>...</svg>")) // <svg>...</svg>

// Numbers and booleans are auto-converted
P("Count: ", 42)           // <p>Count: 42</p>
P("Active: ", true)        // <p>Active: true</p>
```

### Attributes

```go
Div(KV{AttrClass: "container", AttrId: "main"}, "Content")
// <div class="container" id="main">Content</div>
```

### Conditional Rendering

```go
// Show element only if condition is true
If(isLoggedIn, Div("Welcome back!"))

// Choose between two options
IfElse(isAdmin, Div("Admin"), Div("User"))
```

### Lists and Iteration

```go
items := []string{"Apple", "Banana"}

// Map over slice
Ul(
    MapSlice(items, func(item string) HyperNode {
        return Li(item)
    }),
)

// Repeat N times
Div(
    Repeat(3, func() HyperNode {
        return P("Repeated")
    }),
)
```

### With Tailwind CSS

```go
Div(KV{AttrClass: "bg-gray-100 min-h-screen p-8"},
    Div(KV{AttrClass: "max-w-4xl mx-auto"},
        H1(KV{AttrClass: "text-4xl font-bold text-gray-800"}, "Title"),
        P(KV{AttrClass: "text-gray-600 mt-2"}, "Description"),
        Button(KV{AttrClass: "px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"},
            "Click Me",
        ),
    ),
)
```

### With HTMX

```go
// Include htmx library from CDN
Script(KV{AttrSrc: "https://unpkg.com/htmx.org@2.0.8"})

// HTMX button that loads content
Button(KV{
    AttrClass:     "px-4 py-2 bg-blue-500 text-white rounded",
    AttrHxGet:      "/api/users",
    AttrHxTarget:   "#users-list",
    AttrHxSwap:     SwapOuterHtml,
},
    "Load Users",
)

// HTMX form
Form(KV{
    AttrHxPost:     "/api/submit",
    AttrHxTarget:   "#result",
    AttrClass:     "space-y-4",
},
    Input(KV{
        AttrType:  TypeText,
        AttrName:  "message",
        AttrClass: "border rounded px-3 py-2",
    }),
    Button(KV{AttrType: TypeSubmit, AttrClass: "bg-blue-500 text-white px-4 py-2 rounded"},
        "Submit",
    ),
)
```

### Complete Example

```go
package main

import (
    "os"

    . "github.com/assaidy/hyper"
    . "github.com/assaidy/hyper/htmx"
)

func main() {
    users := []string{"Alice", "Bob", "Charlie"}
    isAdmin := true

    page := Empty(
        DoctypeHtml(),
        Html(KV{AttrLang: "en"},
            Head(
                Title("Dashboard"),
                Script(KV{AttrSrc: "https://unpkg.com/htmx.org@2.0.8"}),
                Script(KV{AttrSrc: "https://cdn.tailwindcss.com"}),
            ),
            Body(KV{AttrClass: "bg-gray-100 p-8"},
                Div(KV{AttrClass: "max-w-2xl mx-auto"},
                    H1(KV{AttrClass: "text-3xl font-bold mb-4"}, "Dashboard"),
                    
                    // Conditional admin panel
                    If(isAdmin,
                        Div(KV{AttrClass: "bg-blue-50 p-4 rounded mb-4"},
                            P(KV{AttrClass: "font-semibold"}, "Admin Panel"),
                        ),
                    ),
                    
                    // User count
                    P("Total users: ", len(users)),

                    // HTMX button to refresh users
                    Button(KV{
                        AttrClass:      "px-4 py-2 bg-blue-500 text-white rounded mt-4",
                        AttrHxGet:      "/api/users",
                        AttrHxTarget:   "#users-list",
                        AttrHxSwap:     SwapOuterHtml,
                    },
                        "Refresh Users",
                    ),
                    
                    // User list
                    Ul(KV{AttrClass: "space-y-2 mt-4", AttrId: "users-list"},
                        MapSlice(users, func(name string) HyperNode {
                            return Li(KV{AttrClass: "p-2 bg-white rounded shadow"},
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

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
    page := EMPTY(
        DOCTYPE(),
        HTML(
            HEAD(
                TITLE("My Page"),
            ),
            BODY(
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
DIV("Hello", " ", "World")  // <div>Hello World</div>

P("<script>alert('xss')</script>")
// <p>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</p>

// Raw HTML (not escaped. use with caution)
DIV(RawText("<svg>...</svg>")) // <svg>...</svg>

// Numbers and booleans are auto-converted
P("Count: ", 42)           // <p>Count: 42</p>
P("Active: ", true)        // <p>Active: true</p>
```

### Attributes

```go
DIV(KV{AttrClass: "container", AttrID: "main"}, "Content")
// <div class="container" id="main">Content</div>
```

### Conditional Rendering

```go
// Show element only if condition is true
If(isLoggedIn, DIV("Welcome back!"))

// Choose between two options
IfElse(isAdmin, DIV("Admin"), DIV("User"))
```

### Lists and Iteration

```go
items := []string{"Apple", "Banana"}

// Map over slice
UL(
    Range(items, func(item string) HyperNode {
        return LI(item)
    }),
)

// Repeat N times
DIV(
    Repeat(3, func() HyperNode {
        return P("Repeated")
    }),
)
```

### With Tailwind CSS

```go
DIV(KV{AttrClass: "bg-gray-100 min-h-screen p-8"},
    DIV(KV{AttrClass: "max-w-4xl mx-auto"},
        H1(KV{AttrClass: "text-4xl font-bold text-gray-800"}, "Title"),
        P(KV{AttrClass: "text-gray-600 mt-2"}, "Description"),
        BUTTON(KV{AttrClass: "px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"},
            "Click Me",
        ),
    ),
)
```

### Complete Example

```go
package main

import (
    "os"

    . "github.com/assaidy/hyper/v2"
)

func main() {
    users := []string{"Alice", "Bob", "Charlie"}
    isAdmin := true

    page := EMPTY(
        DOCTYPE(),
        HTML(KV{AttrLang: "en"},
            HEAD(
                TITLE("Dashboard"),
                SCRIPT(KV{AttrSrc: "https://cdn.tailwindcss.com"}),
            ),
            BODY(KV{AttrClass: "bg-gray-100 p-8"},
                DIV(KV{AttrClass: "max-w-2xl mx-auto"},
                    H1(KV{AttrClass: "text-3xl font-bold mb-4"}, "Dashboard"),
                    
                    // Conditional admin panel
                    If(isAdmin,
                        DIV(KV{AttrClass: "bg-blue-50 p-4 rounded mb-4"},
                            P(KV{AttrClass: "font-semibold"}, "Admin Panel"),
                        ),
                    ),
                    
                    // User count
                    P("Total users: ", len(users)),

                    // Standard form submission to refresh users
                    FORM(KV{
                        AttrMethod: MethodPost,
                        AttrAction: "/api/users/refresh",
                    },
                        BUTTON(KV{
                            AttrClass: "px-4 py-2 bg-blue-500 text-white rounded mt-4",
                            AttrType:  TypeSubmit,
                        },
                            "Refresh Users",
                        ),
                    ),
                    
                    // User list
                    UL(KV{AttrClass: "space-y-2 mt-4", AttrID: "users-list"},
                        Range(users, func(name string) HyperNode {
                            return LI(KV{AttrClass: "p-2 bg-white rounded shadow"},
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

package main

import (
	"fmt"
	"os"

	"github.com/assaidy/h"
)

func main() {
	// Sample data for demonstration
	items := []string{"Apple", "Banana", "Cherry", "Date"}
	isLoggedIn := true
	userName := "John Doe"

	// Demonstrate all utility functions in a single example
	page := h.Empty(
		h.DoctypeHTML(),
		h.Html(h.KV{"lang": "en"},
			h.Head(
				h.Meta(h.KV{"charset": "UTF-8"}),
				h.Title("G Utils Example"),
			),
			h.Body(
				// Using IfElse to show conditional content
				h.IfElse(isLoggedIn,
					h.Div(h.KV{"class": "welcome"},
						h.H1(fmt.Sprintf("Welcome back, %s", userName)),
						h.P("You are logged in!"),
					),
					h.Div(h.KV{"class": "login-prompt"},
						h.H1("Please log in"),
						h.P("You need to authenticate to continue."),
					),
				),

				// Using If for optional content
				h.Hr(),
				h.If(isLoggedIn, // Try to toggle this, and see the result
					h.Div(h.KV{"class": "user-actions"},
						h.Button("Profile"),
						" ", // Add a whitespace between the two buttons. not needed if using css styles
						h.Button("Settings"),
					),
				),

				// Using Repeat to generate repeated elements
				h.Hr(),
				h.H2("Repeated Elements"),
				h.Repeat(3, func() h.Node {
					return h.Div(h.KV{"class": "repeated-item"},
						h.Li("This is a repeated item"),
					)
				}),

				// Using Map to transform data into elements
				h.Hr(),
				h.H2("Mapped List"),
				h.Ul(
					h.MapSlice(items, func(item string) h.Node {
						if item == "Apple" {
							return h.Li(item, h.Span(h.KV{"class": "badge"}, " (Popular)"))
						}
						return h.Li(item)
					}),
				),

				// Combining utilities
				h.Hr(),
				h.H2("Combined Example"),
				h.Div(
					"Total items: ", h.Strong(len(items)),
					h.If(len(items) > 2,
						h.P("There are many items to display!"),
					),
				),
			),
		),
	)

	if err := page.Render(os.Stdout); err != nil {
		panic(err)
	}
	// the same as:
	// if err := h.Render(os.Stdout, page); err != nil {
	// }
}

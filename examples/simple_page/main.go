package main

import (
	"fmt"
	"os"

	"github.com/assaidy/g"
)

func main() {
	// Sample data for demonstration
	items := []string{"Apple", "Banana", "Cherry", "Date"}
	isLoggedIn := true
	userName := "John Doe"

	// Demonstrate all utility functions in a single example
	page := gg.Empty(
		gg.DoctypeHTML(),
		gg.Html(gg.KV{"lang": "en"},
			gg.Head(
				gg.Meta(gg.KV{"charset": "UTF-8"}),
				gg.Title("G Utils Example"),
			),
			gg.Body(
				// Using IfElse to show conditional content
				gg.IfElse(isLoggedIn,
					gg.Div(gg.KV{"class": "welcome"},
						gg.H1(fmt.Sprintf("Welcome back, %s", userName)),
						gg.P("You are logged in!"),
					),
					gg.Div(gg.KV{"class": "login-prompt"},
						gg.H1("Please log in"),
						gg.P("You need to authenticate to continue."),
					),
				),

				// Using If for optional content
				gg.Hr(),
				gg.If(isLoggedIn, // Try to toggle this, and see the result
					gg.Div(gg.KV{"class": "user-actions"},
						gg.Button("Profile"),
						" ", // Add a whitespace between the two buttons. not needed if using css styles
						gg.Button("Settings"),
					),
				),

				// Using Repeat to generate repeated elements
				gg.Hr(),
				gg.H2("Repeated Elements"),
				gg.Repeat(3, func() gg.Node {
					return gg.Div(gg.KV{"class": "repeated-item"},
						gg.Li("This is a repeated item"),
					)
				}),

				// Using Map to transform data into elements
				gg.Hr(),
				gg.H2("Mapped List"),
				gg.Ul(
					gg.MapSlice(items, func(item string) gg.Node {
						if item == "Apple" {
							return gg.Li(item, gg.Span(gg.KV{"class": "badge"}, " (Popular)"))
						}
						return gg.Li(item)
					}),
				),

				// Combining utilities
				gg.Hr(),
				gg.H2("Combined Example"),
				gg.Div(
					"Total items: ", gg.Strong(len(items)),
					gg.If(len(items) > 2,
						gg.P("There are many items to display!"),
					),
				),
			),
		),
	)

	if err := page.Render(os.Stdout); err != nil {
		panic(err)
	}
	// the same as:
	// if err := gg.Render(os.Stdout, page); err != nil {
	// }
}

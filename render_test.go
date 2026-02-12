package h

import (
	"bytes"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple text in element",
			node:     Div("Hello World"),
			expected: "<div>Hello World</div>",
			wantErr:  false,
		},
		{
			name:     "Simple element",
			node:     Div(),
			expected: "<div></div>",
			wantErr:  false,
		},
		{
			name:     "Element with children",
			node:     Div("Hello", P("World")),
			expected: "<div>Hello<p>World</p></div>",
			wantErr:  false,
		},
		{
			name:     "Void element",
			node:     Br(),
			expected: "<br>",
			wantErr:  false,
		},
		{
			name:     "Empty element with children",
			node:     Empty("test"),
			expected: "test",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Render(&buf, tt.node)

			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && buf.String() != tt.expected {
				t.Errorf("Render() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}

func TestRender_ErrorHandling(t *testing.T) {
	// Test with a node that will cause an error during rendering
	element := Div(KV{"invalid": nil})

	var buf bytes.Buffer
	err := Render(&buf, element)

	if err == nil {
		t.Error("Render() should return error for invalid attribute")
	}
}

func TestRender_WriteError(t *testing.T) {
	// Create a writer that will return an error on write
	errorWriter := &errorWriter{}
	node := Div("test")

	err := Render(errorWriter, node)
	if err == nil {
		t.Error("Render() should return error when writer fails")
	}

	if !strings.Contains(err.Error(), "write error") {
		t.Errorf("Render() should return writer error, got: %v", err)
	}
}

// errorWriter is a test helper that always returns an error on Write
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (n int, error error) {
	return 0, &writeError{"write error"}
}

type writeError struct {
	msg string
}

func (e *writeError) Error() string {
	return e.msg
}

func TestRender_ComplexStructure(t *testing.T) {
	// Test with a complex nested structure to ensure it handles correctly
	node := Html(KV{"lang": "en"},
		Head(
			Title("Test Page"),
		),
		Body(
			Div(KV{"class": "container"},
				H1("Welcome"),
				P("This is a test."),
				Ul(
					Li("Item 1"),
					Li("Item 2"),
				),
			),
		),
	)

	expected := `<html lang="en"><head><title>Test Page</title></head><body><div class="container"><h1>Welcome</h1><p>This is a test.</p><ul><li>Item 1</li><li>Item 2</li></ul></div></body></html>`

	var buf bytes.Buffer
	err := Render(&buf, node)

	if err != nil {
		t.Errorf("Render() unexpected error: %v", err)
		return
	}

	if buf.String() != expected {
		t.Errorf("Render() complex structure = %q, want %q", buf.String(), expected)
	}
}

func BenchmarkRender_DensePage(b *testing.B) {
	// Create a dense page with many nested elements and attributes
	node := Html(KV{"lang": "en", "data-theme": "light"},
		Head(
			Meta(KV{"charset": "utf-8"}),
			Meta(KV{"name": "viewport", "content": "width=device-width, initial-scale=1"}),
			Title("Dense Page Benchmark"),
			Style(KV{"type": "text/css"}, "body{margin:0;padding:0}"),
			Script(KV{"src": "/app.js", "defer": true}),
		),
		Body(
			Header(KV{"class": "header", "role": "banner"},
				Nav(KV{"class": "navigation", "aria-label": "main"},
					Ul(
						Li(A(KV{"href": "#home"}, "Home")),
						Li(A(KV{"href": "#about"}, "About")),
						Li(A(KV{"href": "#contact"})),
					),
					Main(KV{"class": "main-content", "role": "main"},
						Section(KV{"class": "hero", "id": "hero"},
							Div(KV{"class": "container"},
								H1("Welcome to Our Site"),
								P("This is a dense page for benchmarking purposes."),
								Button(KV{"class": "btn btn-primary", "type": "button"}, "Get Started"),
							),
						),
						Section(KV{"class": "features", "id": "features"},
							Div(KV{"class": "container"},
								H2("Features"),
								Div(KV{"class": "grid"},
									Div(KV{"class": "card"},
										H3("Feature 1"),
										P("Description of feature 1 with lots of content."),
										A(KV{"href": "#", "class": "learn-more"}, "Learn More"),
									),
									Div(KV{"class": "card"},
										H3("Feature 2"),
										P("Description of feature 2 with lots of content."),
										A(KV{"href": "#", "class": "learn-more"}, "Learn More"),
									),
									Div(KV{"class": "card"},
										H3("Feature 3"),
										P("Description of feature 3 with lots of content."),
										A(KV{"href": "#", "class": "learn-more"}, "Learn More"),
									),
								),
							),
						),
					),
					Footer(KV{"class": "footer", "role": "contentinfo"},
						Div(KV{"class": "container"},
							P("© 2024 Dense Page. All rights reserved."),
							Div(KV{"class": "links"},
								A(KV{"href": "#privacy"}, "Privacy")),
							A(KV{"href": "#terms"}, "Terms")),
					),
				),
			),
		),
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := Render(&buf, node)
		if err != nil {
			b.Fatalf("Render() error: %v", err)
		}
	}
}

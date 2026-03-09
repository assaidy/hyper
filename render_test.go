package h

import (
	"bytes"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		node     HyperNode
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple text in element",
			node:     DIV("Hello World"),
			expected: "<div>Hello World</div>",
			wantErr:  false,
		},
		{
			name:     "Simple element",
			node:     DIV(),
			expected: "<div></div>",
			wantErr:  false,
		},
		{
			name:     "Element with children",
			node:     DIV("Hello", P("World")),
			expected: "<div>Hello<p>World</p></div>",
			wantErr:  false,
		},
		{
			name:     "Void element",
			node:     BR(),
			expected: "<br>",
			wantErr:  false,
		},
		{
			name:     "Empty element with children",
			node:     EMPTY("test"),
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
	element := DIV(KV{"invalid": nil})

	var buf bytes.Buffer
	err := Render(&buf, element)

	if err == nil {
		t.Error("Render() should return error for invalid attribute")
	}
}

func TestRender_WriteError(t *testing.T) {
	// Create a writer that will return an error on write
	errorWriter := &errorWriter{}
	node := DIV("test")

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
	node := HTML(KV{"lang": "en"},
		HEAD(
			TITLE("Test Page"),
		),
		BODY(
			DIV(KV{"class": "container"},
				H1("Welcome"),
				P("This is a test."),
				UL(
					LI("Item 1"),
					LI("Item 2"),
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
	node := HTML(KV{"lang": "en", "data-theme": "light"},
		HEAD(
			META(KV{"charset": "utf-8"}),
			META(KV{"name": "viewport", "content": "width=device-width, initial-scale=1"}),
			TITLE("Dense Page Benchmark"),
			STYLE(KV{"type": "text/css"}, "body{margin:0;padding:0}"),
			SCRIPT(KV{"src": "/app.js", "defer": true}),
		),
		BODY(
			HEADER(KV{"class": "header", "role": "banner"},
				NAV(KV{"class": "navigation", "aria-label": "main"},
					UL(
						LI(A(KV{"href": "#home"}, "Home")),
						LI(A(KV{"href": "#about"}, "About")),
						LI(A(KV{"href": "#contact"})),
					),
					MAIN(KV{"class": "main-content", "role": "main"},
						SECTION(KV{"class": "hero", "id": "hero"},
							DIV(KV{"class": "container"},
								H1("Welcome to Our Site"),
								P("This is a dense page for benchmarking purposes."),
								BUTTON(KV{"class": "btn btn-primary", "type": "button"}, "Get Started"),
							),
						),
						SECTION(KV{"class": "features", "id": "features"},
							DIV(KV{"class": "container"},
								H2("Features"),
								DIV(KV{"class": "grid"},
									DIV(KV{"class": "card"},
										H3("Feature 1"),
										P("Description of feature 1 with lots of content."),
										A(KV{"href": "#", "class": "learn-more"}, "Learn More"),
									),
									DIV(KV{"class": "card"},
										H3("Feature 2"),
										P("Description of feature 2 with lots of content."),
										A(KV{"href": "#", "class": "learn-more"}, "Learn More"),
									),
									DIV(KV{"class": "card"},
										H3("Feature 3"),
										P("Description of feature 3 with lots of content."),
										A(KV{"href": "#", "class": "learn-more"}, "Learn More"),
									),
								),
							),
						),
					),
					FOOTER(KV{"class": "footer", "role": "contentinfo"},
						DIV(KV{"class": "container"},
							P("© 2024 Dense Page. All rights reserved."),
							DIV(KV{"class": "links"},
								A(KV{"href": "#privacy"}, "Privacy")),
							A(KV{"href": "#terms"}, "Terms")),
					),
				),
			),
		),
	)

	
	b.ReportAllocs()

	for b.Loop() {
		var buf bytes.Buffer
		err := Render(&buf, node)
		if err != nil {
			b.Fatalf("Render() error: %v", err)
		}
	}
}

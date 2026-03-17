package hyper

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
			node:     DIV()("Hello World"),
			expected: "<div>Hello World</div>",
			wantErr:  false,
		},
		{
			name:     "Simple element",
			node:     DIV()(),
			expected: "<div></div>",
			wantErr:  false,
		},
		{
			name:     "Element with children",
			node:     DIV()("Hello", P()("World")),
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
			node:     Group("test"),
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

func TestRender_WriteError(t *testing.T) {
	errorWriter := &errorWriter{}
	node := DIV()("test")

	err := Render(errorWriter, node)
	if err == nil {
		t.Error("Render() should return error when writer fails")
	}

	if !strings.Contains(err.Error(), "write error") {
		t.Errorf("Render() should return writer error, got: %v", err)
	}
}

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
	node := HTML(AttrLang("en"))(
		HEAD()(
			TITLE()("Test Page"),
		),
		BODY()(
			DIV(AttrClass("container"))(
				H1()("Welcome"),
				P()("This is a test."),
				UL()(
					LI()("Item 1"),
					LI()("Item 2"),
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
	node := HTML(Attr("lang", "en"), Attr("data-theme", "light"))(
		HEAD()(
			META(AttrCharset("utf-8")),
			META(AttrName("viewport"), Attr("content", "width=device-width, initial-scale=1")),
			TITLE()("Dense Page Benchmark"),
			STYLE(AttrType("text/css"))("body{margin:0;padding:0}"),
			SCRIPT(AttrSrc("/app.js"), AttrDefer(true)),
		),
		BODY()(
			HEADER(AttrClass("header"), AttrRole("banner"))(
				NAV(AttrClass("navigation"), Attr("aria-label", "main"))(
					UL()(
						LI()(A(AttrHref("#home"))("Home")),
						LI()(A(AttrHref("#about"))("About")),
						LI()(A(AttrHref("#contact"))()),
					),
					MAIN(AttrClass("main-content"), AttrRole("main"))(
						SECTION(AttrClass("hero"), AttrID("hero"))(
							DIV(AttrClass("container"))(
								H1()("Welcome to Our Site"),
								P()("This is a dense page for benchmarking purposes."),
								BUTTON(AttrClass("btn btn-primary"), AttrType("button"))("Get Started"),
							),
						),
						SECTION(AttrClass("features"), AttrID("features"))(
							DIV(AttrClass("container"))(
								H2()("Features"),
								DIV(AttrClass("grid"))(
									DIV(AttrClass("card"))(
										H3()("Feature 1"),
										P()("Description of feature 1 with lots of content."),
										A(AttrHref("#"), AttrClass("learn-more"))("Learn More"),
									),
									DIV(AttrClass("card"))(
										H3()("Feature 2"),
										P()("Description of feature 2 with lots of content."),
										A(AttrHref("#"), AttrClass("learn-more"))("Learn More"),
									),
									DIV(AttrClass("card"))(
										H3()("Feature 3"),
										P()("Description of feature 3 with lots of content."),
										A(AttrHref("#"), AttrClass("learn-more"))("Learn More"),
									),
								),
							),
						),
					),
					FOOTER(AttrClass("footer"), AttrRole("contentinfo"))(
						DIV(AttrClass("container"))(
							P()("© 2024 Dense Page. All rights reserved."),
							DIV(AttrClass("links"))(
								A(AttrHref("#privacy"))("Privacy"),
							),
							A(AttrHref("#terms"))("Terms"),
						),
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

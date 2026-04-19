package hyper

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestIfElse(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		result    string
		alt       string
		expected  string
	}{
		{
			name:      "Condition true returns result",
			condition: true,
			result:    "yes",
			alt:       "no",
			expected:  "yes",
		},
		{
			name:      "Condition false returns alternative",
			condition: false,
			result:    "yes",
			alt:       "no",
			expected:  "no",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfElse(tt.condition, tt.result, tt.alt)
			if result != tt.expected {
				t.Errorf("IfElse() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIfElse_Nodes(t *testing.T) {
	trueNode := DIV()("true")
	falseNode := P()("false")

	tests := []struct {
		name      string
		condition bool
		expected  string
	}{
		{
			name:      "Condition true returns node",
			condition: true,
			expected:  "<div>true</div>",
		},
		{
			name:      "Condition false returns node",
			condition: false,
			expected:  "<p>false</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := IfElse(tt.condition, trueNode, falseNode)
			var buf bytes.Buffer
			err := Render(&buf, node)
			if err != nil {
				t.Errorf("IfElse() node render error: %v", err)
				return
			}
			if buf.String() != tt.expected {
				t.Errorf("IfElse() node render = %v, want %v", buf.String(), tt.expected)
			}
		})
	}
}

func TestIf(t *testing.T) {
	t.Run("Basic If", func(t *testing.T) {
		node := DIV()("content")
		resultNode := If(true, node)
		var buf bytes.Buffer
		err := Render(&buf, resultNode)
		if err != nil {
			t.Errorf("If() render error: %v", err)
			return
		}
		if buf.String() != "<div>content</div>" {
			t.Errorf("If() = %v, want <div>content</div>", buf.String())
		}
	})

	t.Run("Condition false returns empty", func(t *testing.T) {
		node := DIV()("content")
		resultNode := If(false, node)
		var buf bytes.Buffer
		err := Render(&buf, resultNode)
		if err != nil {
			t.Errorf("If() render error: %v", err)
			return
		}
		if buf.String() != "" {
			t.Errorf("If() = %v, want empty", buf.String())
		}
	})

	t.Run("If with ElseIf and Else", func(t *testing.T) {
		resultNode := If(false, DIV()("first")).
			ElseIf(true, SPAN()("second")).
			Else(DIV()("default"))
		var buf bytes.Buffer
		err := Render(&buf, resultNode)
		if err != nil {
			t.Errorf("If().ElseIf().Else() render error: %v", err)
			return
		}
		if buf.String() != "<span>second</span>" {
			t.Errorf("If().ElseIf().Else() = %v, want <span>second</span>", buf.String())
		}
	})

	t.Run("Multiple ElseIf branches", func(t *testing.T) {
		resultNode := If(false, DIV()("first")).
			ElseIf(false, SPAN()("second")).
			ElseIf(true, P()("third")).
			Else(DIV()("default"))
		var buf bytes.Buffer
		err := Render(&buf, resultNode)
		if err != nil {
			t.Errorf("If().ElseIf().ElseIf().Else() render error: %v", err)
			return
		}
		if buf.String() != "<p>third</p>" {
			t.Errorf("If().ElseIf().ElseIf().Else() = %v, want <p>third</p>", buf.String())
		}
	})

	t.Run("No condition matches falls through to Else", func(t *testing.T) {
		resultNode := If(false, DIV()("first")).
			ElseIf(false, SPAN()("second")).
			Else(DIV()("default"))
		var buf bytes.Buffer
		err := Render(&buf, resultNode)
		if err != nil {
			t.Errorf("If().ElseIf().Else() render error: %v", err)
			return
		}
		if buf.String() != "<div>default</div>" {
			t.Errorf("If().ElseIf().Else() = %v, want <div>default</div>", buf.String())
		}
	})

	t.Run("If with ElseIf without Else", func(t *testing.T) {
		resultNode := If(false, DIV()("first")).
			ElseIf(false, SPAN()("second"))
		var buf bytes.Buffer
		err := Render(&buf, resultNode)
		if err != nil {
			t.Errorf("If().ElseIf() render error: %v", err)
			return
		}
		if buf.String() != "" {
			t.Errorf("If().ElseIf() = %v, want empty", buf.String())
		}
	})

	t.Run("If with ElseIf without Else matches", func(t *testing.T) {
		resultNode := If(false, DIV()("first")).
			ElseIf(true, DIV()("trial"))
		var buf bytes.Buffer
		err := Render(&buf, resultNode)
		if err != nil {
			t.Errorf("If().ElseIf() render error: %v", err)
			return
		}
		if buf.String() != "<div>trial</div>" {
			t.Errorf("If().ElseIf() = %v, want <div>trial</div>", buf.String())
		}
	})
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		f        func() HyperNode
		expected string
	}{
		{
			name:     "Repeat zero times",
			n:        0,
			f:        func() HyperNode { return DIV()() },
			expected: "",
		},
		{
			name:     "Repeat once",
			n:        1,
			f:        func() HyperNode { return DIV()("item") },
			expected: "<div>item</div>",
		},
		{
			name:     "Repeat multiple times",
			n:        3,
			f:        func() HyperNode { return DIV()("item") },
			expected: "<div>item</div><div>item</div><div>item</div>",
		},
		{
			name: "Repeat with different content",
			n:    2,
			f: func() HyperNode {
				static := 0
				static++
				return DIV()(string(rune('a' + static)))
			},
			expected: "<div>b</div><div>b</div>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultNode := Repeat(tt.n, tt.f)
			var buf bytes.Buffer
			err := Render(&buf, resultNode)
			if err != nil {
				t.Errorf("Repeat() node render error: %v", err)
				return
			}
			if buf.String() != tt.expected {
				t.Errorf("Repeat() node render = %v, want %v", buf.String(), tt.expected)
			}
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		f        func(string) HyperNode
		expected string
	}{
		{
			name:     "Range empty slice",
			input:    []string{},
			f:        func(s string) HyperNode { return LI()(s) },
			expected: "",
		},
		{
			name:     "Range single item",
			input:    []string{"apple"},
			f:        func(s string) HyperNode { return LI()(s) },
			expected: "<li>apple</li>",
		},
		{
			name:     "Range multiple items",
			input:    []string{"apple", "banana", "cherry"},
			f:        func(s string) HyperNode { return LI()(s) },
			expected: "<li>apple</li><li>banana</li><li>cherry</li>",
		},
		{
			name:  "Range with conditional logic",
			input: []string{"apple", "banana"},
			f: func(s string) HyperNode {
				if s == "apple" {
					return LI()(s, SPAN()(" (popular)"))
				}
				return LI()(s)
			},
			expected: "<li>apple<span> (popular)</span></li><li>banana</li>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultNode := Range(tt.input, tt.f)
			var buf bytes.Buffer
			err := Render(&buf, resultNode)
			if err != nil {
				t.Errorf("Range() node render error: %v", err)
				return
			}
			if buf.String() != tt.expected {
				t.Errorf("Range() node render = %v, want %v", buf.String(), tt.expected)
			}
		})
	}
}

func TestRange_Integers(t *testing.T) {
	numbers := []int{1, 2, 3}
	resultNode := Range(numbers, func(n int) HyperNode {
		return DIV()(string(rune('0' + n)))
	})

	var buf bytes.Buffer
	err := Render(&buf, resultNode)
	if err != nil {
		t.Errorf("Range() integers node render error: %v", err)
		return
	}
	expected := "<div>1</div><div>2</div><div>3</div>"
	if buf.String() != expected {
		t.Errorf("Range() integers node render = %v, want %v", buf.String(), expected)
	}
}

func TestCache(t *testing.T) {
	t.Run("first render caches output", func(t *testing.T) {
		var cache NodeCache
		renderCount := 0

		makeDynamicNode := func() HyperNode {
			return Group(
				DIV()("Static content"),
				func() HyperNode {
					renderCount++
					return SPAN()(renderCount)
				}(),
			)
		}

		cached := Cache(&cache, makeDynamicNode())

		var buf1 bytes.Buffer
		if err := Render(&buf1, cached); err != nil {
			t.Errorf("first render error: %v", err)
		}

		if renderCount != 1 {
			t.Errorf("renderCount = %d, want 1", renderCount)
		}

		var buf2 bytes.Buffer
		if err := Render(&buf2, cached); err != nil {
			t.Errorf("second render error: %v", err)
		}

		if renderCount != 1 {
			t.Errorf("renderCount after second render = %d, want 1 (cached)", renderCount)
		}

		if buf1.String() != buf2.String() {
			t.Errorf("buf1 = %v, buf2 = %v, want same output", buf1.String(), buf2.String())
		}
	})

	t.Run("multiple renders produce same output", func(t *testing.T) {
		var cache NodeCache
		complexNode := DIV(AttrClass("header"), AttrID("main"))(
			NAV()(
				A(AttrHref("/"))("Home"),
				A(AttrHref("/about"))("About"),
			),
			HEADER()(
				H1()("Welcome"),
				P()("Static content"),
			),
		)

		cached := Cache(&cache, complexNode)

		var outputs []string
		for i := 0; i < 5; i++ {
			var buf bytes.Buffer
			if err := Render(&buf, cached); err != nil {
				t.Errorf("render %d error: %v", i, err)
			}
			outputs = append(outputs, buf.String())
		}

		for i, out := range outputs {
			if out != outputs[0] {
				t.Errorf("output %d = %v, want %v", i, out, outputs[0])
			}
		}
	})

	t.Run("different caches produce independent results", func(t *testing.T) {
		var cache1, cache2 NodeCache
		counter := 0

		makeNode := func() HyperNode {
			return Group(
				DIV()(func() string {
					counter++
					return fmt.Sprintf("%d", counter)
				}()),
			)
		}

		cached1 := Cache(&cache1, makeNode())
		cached2 := Cache(&cache2, makeNode())

		var buf1a, buf1b, buf2a, buf2b bytes.Buffer
		Render(&buf1a, cached1)
		Render(&buf1b, cached1)
		Render(&buf2a, cached2)
		Render(&buf2b, cached2)

		if buf1a.String() != "<div>1</div>" {
			t.Errorf("buf1a = %v, want <div>1</div>", buf1a.String())
		}
		if buf1b.String() != "<div>1</div>" {
			t.Errorf("buf1b = %v, want <div>1</div> (cached)", buf1b.String())
		}
		if buf2a.String() != "<div>2</div>" {
			t.Errorf("buf2a = %v, want <div>2</div>", buf2a.String())
		}
		if buf2b.String() != "<div>2</div>" {
			t.Errorf("buf2b = %v, want <div>2</div> (cached)", buf2b.String())
		}
		if counter != 2 {
			t.Errorf("counter = %d, want 2 (each cache renders once)", counter)
		}
	})

	t.Run("cache with complex nested structure", func(t *testing.T) {
		var cache NodeCache
		renderCount := 0

		makeComplexNode := func() HyperNode {
			return Group(
				DIV(AttrClass("wrapper"))(
					HEADER()(
						H1()("Site Title"),
						NAV()(
							UL()(
								Range([]string{"Home", "About", "Contact"}, func(item string) HyperNode {
									return LI()(A(AttrHref("/" + strings.ToLower(item)))(item))
								}),
							),
						),
					),
					MAIN()(
						ARTICLE()(
							H2()("Article Title"),
							P()("Lorem ipsum dolor sit amet."),
							IMG(AttrSrc("/image.jpg"), AttrAlt("Image")),
						),
					),
					FOOTER()(
						P()("Copyright 2024"),
					),
				),
				func() HyperNode {
					renderCount++
					return SPAN()(renderCount)
				}(),
			)
		}

		cached := Cache(&cache, makeComplexNode())

		var buf bytes.Buffer
		Render(&buf, cached)
		Render(&buf, cached)

		if renderCount != 1 {
			t.Errorf("renderCount = %d, want 1", renderCount)
		}
	})
}

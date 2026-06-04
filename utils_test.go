package hyper

import (
	"bytes"
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

func TestIfElseZero(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		result    string
		expected  string
	}{
		{
			name:      "Condition true returns result",
			condition: true,
			result:    "text-red",
			expected:  "text-red",
		},
		{
			name:      "Condition false returns zero value",
			condition: false,
			result:    "text-red",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfElseZero(tt.condition, tt.result)
			if result != tt.expected {
				t.Errorf("IfElseZero() = %q, want %q", result, tt.expected)
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
		f        func() any
		expected string
	}{
		{
			name:     "Repeat zero times",
			n:        0,
			f:        func() any { return DIV()() },
			expected: "",
		},
		{
			name:     "Repeat once",
			n:        1,
			f:        func() any { return DIV()("item") },
			expected: "<div>item</div>",
		},
		{
			name:     "Repeat multiple times",
			n:        3,
			f:        func() any { return DIV()("item") },
			expected: "<div>item</div><div>item</div><div>item</div>",
		},
		{
			name: "Repeat with different content",
			n:    2,
			f: func() any {
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
		run      func() HyperNode
		expected string
	}{
		{
			name: "empty slice",
			run: func() HyperNode {
				return Range([]string{}, func(s string) any { return LI()(s) })
			},
			expected: "",
		},
		{
			name: "single string item",
			run: func() HyperNode {
				return Range([]string{"apple"}, func(s string) any { return LI()(s) })
			},
			expected: "<li>apple</li>",
		},
		{
			name: "multiple string items",
			run: func() HyperNode {
				return Range([]string{"apple", "banana", "cherry"}, func(s string) any { return LI()(s) })
			},
			expected: "<li>apple</li><li>banana</li><li>cherry</li>",
		},
		{
			name: "string items with conditional logic",
			run: func() HyperNode {
				return Range([]string{"apple", "banana"}, func(s string) any {
					if s == "apple" {
						return LI()(s, SPAN()(" (popular)"))
					}
					return LI()(s)
				})
			},
			expected: "<li>apple<span> (popular)</span></li><li>banana</li>",
		},
		{
			name: "integer items",
			run: func() HyperNode {
				return Range([]int{1, 2, 3}, func(n int) any { return n })
			},
			expected: "123",
		},
		{
			name: "bool items",
			run: func() HyperNode {
				return Range([]bool{true, false}, func(b bool) any { return b })
			},
			expected: "truefalse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultNode := tt.run()
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

func TestClasses(t *testing.T) {
	tests := []struct {
		name     string
		classes  []string
		expected string
	}{
		{
			name:     "no arguments",
			classes:  nil,
			expected: "",
		},
		{
			name:     "empty slice",
			classes:  []string{},
			expected: "",
		},
		{
			name:     "single class",
			classes:  []string{"btn"},
			expected: "btn",
		},
		{
			name:     "multiple classes",
			classes:  []string{"btn", "btn-primary", "active"},
			expected: "btn btn-primary active",
		},
		{
			name:     "empty strings skipped",
			classes:  []string{"btn", "", "active"},
			expected: "btn active",
		},
		{
			name:     "duplicates removed",
			classes:  []string{"btn", "btn-primary", "btn", "active"},
			expected: "btn btn-primary active",
		},
		{
			name:     "whitespace trimmed",
			classes:  []string{" btn ", "active"},
			expected: "btn active",
		},
		{
			name:     "all whitespace",
			classes:  []string{"  ", "\t", "\n"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Classes(tt.classes...)
			if result != tt.expected {
				t.Errorf("Classes() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestJson(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "number",
			input:    42,
			expected: `42`,
		},
		{
			name:     "map",
			input:    Object{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			name:     "nested Object",
			input:    Object{"a": Object{"b": 1}},
			expected: `{"a":{"b":1}}`,
		},
		{
			name:     "slice",
			input:    []int{1, 2, 3},
			expected: `[1,2,3]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Json(tt.input)
			if result != tt.expected {
				t.Errorf("Json() = %s, want %s", result, tt.expected)
			}
		})
	}

	t.Run("panic on unmarshalable value", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic")
			}
		}()
		Json(make(chan int))
	})
}

func TestObject(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected string
	}{
		{
			name:     "empty",
			input:    Object{},
			expected: `{}`,
		},
		{
			name:     "single key",
			input:    Object{"a": 1},
			expected: `{"a":1}`,
		},
		{
			name:     "mixed types",
			input:    Object{"name": "test", "count": 3, "active": true},
			expected: `{"active":true,"count":3,"name":"test"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Json(tt.input)
			if result != tt.expected {
				t.Errorf("Object via Json() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestOnceWithKey(t *testing.T) {
	t.Run("renders and caches", func(t *testing.T) {
		counter := 0
		n := OnceWithKey("test-1", func() HyperNode {
			counter++
			return DIV()("hello")
		})

		var buf1 bytes.Buffer
		err := Render(&buf1, n)
		if err != nil {
			t.Fatal(err)
		}
		if buf1.String() != "<div>hello</div>" {
			t.Errorf("got %q", buf1.String())
		}
		if counter != 1 {
			t.Errorf("node func called %d times, want 1", counter)
		}

		var buf2 bytes.Buffer
		err = Render(&buf2, n)
		if err != nil {
			t.Fatal(err)
		}
		if buf2.String() != "<div>hello</div>" {
			t.Errorf("got %q", buf2.String())
		}
		if counter != 1 {
			t.Errorf("node func called %d times, want 1", counter)
		}
	})

	t.Run("different keys are independent", func(t *testing.T) {
		counterA := 0
		counterB := 0
		a := OnceWithKey("key-a", func() HyperNode {
			counterA++
			return DIV()("a")
		})
		b := OnceWithKey("key-b", func() HyperNode {
			counterB++
			return P()("b")
		})

		var buf bytes.Buffer
		Render(&buf, a)
		Render(&buf, b)
		if counterA != 1 || counterB != 1 {
			t.Errorf("counterA=%d counterB=%d, want both 1", counterA, counterB)
		}

		Render(&buf, a)
		Render(&buf, b)
		if counterA != 1 || counterB != 1 {
			t.Errorf("counterA=%d counterB=%d, want both 1 (cached)", counterA, counterB)
		}
	})

	t.Run("same key shares cache across instances", func(t *testing.T) {
		counter := 0
		fn := func() HyperNode {
			counter++
			return DIV()("shared")
		}

		a := OnceWithKey("shared-key", fn)
		var buf1 bytes.Buffer
		Render(&buf1, a)

		b := OnceWithKey("shared-key", fn)
		var buf2 bytes.Buffer
		Render(&buf2, b)

		if counter != 1 {
			t.Errorf("node func called %d times, want 1 (shared cache)", counter)
		}
	})
}

func TestOnce(t *testing.T) {
	t.Run("renders and caches", func(t *testing.T) {
		counter := 0
		n := Once(func() HyperNode {
			counter++
			return DIV()("hello")
		})

		var buf1 bytes.Buffer
		err := Render(&buf1, n)
		if err != nil {
			t.Fatal(err)
		}
		if buf1.String() != "<div>hello</div>" {
			t.Errorf("got %q", buf1.String())
		}
		if counter != 1 {
			t.Errorf("node func called %d times, want 1", counter)
		}

		var buf2 bytes.Buffer
		err = Render(&buf2, n)
		if err != nil {
			t.Fatal(err)
		}
		if buf2.String() != "<div>hello</div>" {
			t.Errorf("got %q", buf2.String())
		}
		if counter != 1 {
			t.Errorf("node func called %d times, want 1", counter)
		}
	})

	t.Run("different call sites produce different keys", func(t *testing.T) {
		counterA := 0
		counterB := 0

		onceA := Once(func() HyperNode {
			counterA++
			return DIV()("a")
		})
		onceB := Once(func() HyperNode {
			counterB++
			return DIV()("b")
		})

		var buf bytes.Buffer
		Render(&buf, onceA)
		Render(&buf, onceB)
		Render(&buf, onceA)
		Render(&buf, onceB)

		if counterA != 1 || counterB != 1 {
			t.Errorf("counterA=%d counterB=%d, want both 1", counterA, counterB)
		}
	})
}

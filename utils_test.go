package h

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
	trueNode := Div("true")
	falseNode := P("false")

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
	node := Div("content")

	tests := []struct {
		name      string
		condition bool
		expected  string
	}{
		{
			name:      "Condition true returns node",
			condition: true,
			expected:  "<div>content</div>",
		},
		{
			name:      "Condition false returns empty",
			condition: false,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultNode := If(tt.condition, node)
			var buf bytes.Buffer
			err := Render(&buf, resultNode)
			if err != nil {
				t.Errorf("If() node render error: %v", err)
				return
			}
			if buf.String() != tt.expected {
				t.Errorf("If() node render = %v, want %v", buf.String(), tt.expected)
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		f        func() Node
		expected string
	}{
		{
			name:     "Repeat zero times",
			n:        0,
			f:        func() Node { return Div() },
			expected: "",
		},
		{
			name:     "Repeat once",
			n:        1,
			f:        func() Node { return Div("item") },
			expected: "<div>item</div>",
		},
		{
			name:     "Repeat multiple times",
			n:        3,
			f:        func() Node { return Div("item") },
			expected: "<div>item</div><div>item</div><div>item</div>",
		},
		{
			name: "Repeat with different content",
			n:    2,
			f: func() Node {
				static := 0
				static++
				return Div(string(rune('a' + static)))
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

func TestMapSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		f        func(string) Node
		expected string
	}{
		{
			name:     "MapSlice empty slice",
			input:    []string{},
			f:        func(s string) Node { return Li(s) },
			expected: "",
		},
		{
			name:     "MapSlice single item",
			input:    []string{"apple"},
			f:        func(s string) Node { return Li(s) },
			expected: "<li>apple</li>",
		},
		{
			name:     "MapSlice multiple items",
			input:    []string{"apple", "banana", "cherry"},
			f:        func(s string) Node { return Li(s) },
			expected: "<li>apple</li><li>banana</li><li>cherry</li>",
		},
		{
			name:  "MapSlice with conditional logic",
			input: []string{"apple", "banana"},
			f: func(s string) Node {
				if s == "apple" {
					return Li(s, Span(" (popular)"))
				}
				return Li(s)
			},
			expected: "<li>apple<span> (popular)</span></li><li>banana</li>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultNode := MapSlice(tt.input, tt.f)
			var buf bytes.Buffer
			err := Render(&buf, resultNode)
			if err != nil {
				t.Errorf("MapSlice() node render error: %v", err)
				return
			}
			if buf.String() != tt.expected {
				t.Errorf("MapSlice() node render = %v, want %v", buf.String(), tt.expected)
			}
		})
	}
}

func TestMapSlice_Integers(t *testing.T) {
	numbers := []int{1, 2, 3}
	resultNode := MapSlice(numbers, func(n int) Node {
		return Div(string(rune('0' + n)))
	})

	var buf bytes.Buffer
	err := Render(&buf, resultNode)
	if err != nil {
		t.Errorf("MapSlice() integers node render error: %v", err)
		return
	}
	expected := "<div>1</div><div>2</div><div>3</div>"
	if buf.String() != expected {
		t.Errorf("MapSlice() integers node render = %v, want %v", buf.String(), expected)
	}
}

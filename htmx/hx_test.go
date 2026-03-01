package hx

import "testing"

func TestAttrHxOn(t *testing.T) {
	tests := []struct {
		name     string
		event    string
		expected string
	}{
		{
			name:     "simple click event",
			event:    "click",
			expected: "hx-on:click",
		},
		{
			name:     "htmx event",
			event:    "htmx:before-swap",
			expected: "hx-on:htmx:before-swap",
		},
		{
			name:     "mouseover event",
			event:    "mouseover",
			expected: "hx-on:mouseover",
		},
		{
			name:     "load event",
			event:    "load",
			expected: "hx-on:load",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AttrHxOn(tt.event)
			if result != tt.expected {
				t.Errorf("AttrHxOn(%q) = %q, want %q", tt.event, result, tt.expected)
			}
		})
	}
}

func TestSwapOob(t *testing.T) {
	tests := []struct {
		name     string
		swap     string
		target   string
		expected string
	}{
		{
			name:     "simple id selector",
			swap:     SwapBeforeEnd,
			target:   "#notifications",
			expected: "beforeend:#notifications",
		},
		{
			name:     "complex selector with space",
			swap:     SwapBeforeEnd,
			target:   "#table tbody",
			expected: "beforeend:#table tbody",
		},
		{
			name:     "outerHTML with id",
			swap:     SwapOuterHtml,
			target:   "#circle1",
			expected: "outerHTML:#circle1",
		},
		{
			name:     "innerHTML with id",
			swap:     SwapInnerHtml,
			target:   "#content",
			expected: "innerHTML:#content",
		},
		{
			name:     "beforebegin with id",
			swap:     SwapBeforeBegin,
			target:   "#element",
			expected: "beforebegin:#element",
		},
		{
			name:     "afterend with class selector",
			swap:     SwapAfterEnd,
			target:   ".container",
			expected: "afterend:.container",
		},
		{
			name:     "delete swap",
			swap:     SwapDelete,
			target:   "#row",
			expected: "delete:#row",
		},
		{
			name:     "none swap",
			swap:     SwapNone,
			target:   "#placeholder",
			expected: "none:#placeholder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SwapOob(tt.swap, tt.target)
			if result != tt.expected {
				t.Errorf("SwapOob(%q, %q) = %q, want %q", tt.swap, tt.target, result, tt.expected)
			}
		})
	}
}

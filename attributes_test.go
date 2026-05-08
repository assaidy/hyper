package hyper

import "testing"

func TestAttrReflect(t *testing.T) {
	type derivedString string
	type derivedBool bool

	tests := []struct {
		name     string
		key      string
		value    any
		expected Attribute
	}{
		{
			name:  "direct string creates PairAttribute",
			key:   "class",
			value: "container",
			expected: PairAttribute{
				Key:   "class",
				Value: "container",
			},
		},
		{
			name:  "derived string creates PairAttribute",
			key:   "class",
			value: derivedString("container"),
			expected: PairAttribute{
				Key:   "class",
				Value: "container",
			},
		},
		{
			name:  "direct bool true creates active BooleanAttribute",
			key:   "hidden",
			value: true,
			expected: BooleanAttribute{
				Key:      "hidden",
				IsActive: true,
			},
		},
		{
			name:  "direct bool false creates inactive BooleanAttribute",
			key:   "disabled",
			value: false,
			expected: BooleanAttribute{
				Key:      "disabled",
				IsActive: false,
			},
		},
		{
			name:  "derived bool true creates active BooleanAttribute",
			key:   "hidden",
			value: derivedBool(true),
			expected: BooleanAttribute{
				Key:      "hidden",
				IsActive: true,
			},
		},
		{
			name:  "derived bool false creates inactive BooleanAttribute",
			key:   "disabled",
			value: derivedBool(false),
			expected: BooleanAttribute{
				Key:      "disabled",
				IsActive: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := attrReflect(tt.key, tt.value)
			switch expected := tt.expected.(type) {
			case PairAttribute:
				pair, ok := result.(PairAttribute)
				if !ok {
					t.Errorf("expected PairAttribute, got %T", result)
					return
				}
				if pair.Key != expected.Key || pair.Value != expected.Value {
					t.Errorf("PairAttribute = {Key: %q, Value: %q}, want {Key: %q, Value: %q}",
						pair.Key, pair.Value, expected.Key, expected.Value)
				}
			case BooleanAttribute:
				boolAttr, ok := result.(BooleanAttribute)
				if !ok {
					t.Errorf("expected BooleanAttribute, got %T", result)
					return
				}
				if boolAttr.Key != expected.Key || boolAttr.IsActive != expected.IsActive {
					t.Errorf("BooleanAttribute = {Key: %q, IsActive: %v}, want {Key: %q, IsActive: %v}",
						boolAttr.Key, boolAttr.IsActive, expected.Key, expected.IsActive)
				}
			}
		})
	}
}

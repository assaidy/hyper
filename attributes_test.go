package hyper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPairAttribute(t *testing.T) {
	t.Run("renders a leading space, key, and quoted value", func(t *testing.T) {
		out, err := renderAttribute(t, PairAttribute{Key: "class", Value: "box"})
		require.NoError(t, err)
		require.Equal(t, ` class="box"`, out)
	})

	t.Run("trims whitespace around the key", func(t *testing.T) {
		out, err := renderAttribute(t, PairAttribute{Key: "  class  ", Value: "box"})
		require.NoError(t, err)
		require.Equal(t, ` class="box"`, out)
	})

	t.Run("escapes quotes in the value", func(t *testing.T) {
		out, err := renderAttribute(t, PairAttribute{Key: "title", Value: `say "hi"`})
		require.NoError(t, err)
		require.Equal(t, ` title="say &quot;hi&quot;"`, out)
	})

	t.Run("html-escapes the key", func(t *testing.T) {
		out, err := renderAttribute(t, PairAttribute{Key: "da<ta>", Value: "v"})
		require.NoError(t, err)
		require.Equal(t, ` da&lt;ta&gt;="v"`, out)
	})

	t.Run("rejects an empty key without writing anything", func(t *testing.T) {
		out, err := renderAttribute(t, PairAttribute{Key: "", Value: "v"})
		require.EqualError(t, err, "empty/whitespace attribute key not allowed.")
		require.Empty(t, out)
	})

	t.Run("rejects a whitespace-only key without writing anything", func(t *testing.T) {
		out, err := renderAttribute(t, PairAttribute{Key: "   ", Value: "v"})
		require.EqualError(t, err, "empty/whitespace attribute key not allowed.")
		require.Empty(t, out)
	})
}

func TestBooleanAttribute(t *testing.T) {
	t.Run("renders only the key when active", func(t *testing.T) {
		out, err := renderAttribute(t, BooleanAttribute{Key: "disabled", IsActive: true})
		require.NoError(t, err)
		require.Equal(t, " disabled", out)
	})

	t.Run("renders nothing when inactive", func(t *testing.T) {
		out, err := renderAttribute(t, BooleanAttribute{Key: "disabled", IsActive: false})
		require.NoError(t, err)
		require.Empty(t, out)
	})

	t.Run("trims whitespace around the key", func(t *testing.T) {
		out, err := renderAttribute(t, BooleanAttribute{Key: " checked ", IsActive: true})
		require.NoError(t, err)
		require.Equal(t, " checked", out)
	})

	t.Run("rejects an empty key even when inactive", func(t *testing.T) {
		out, err := renderAttribute(t, BooleanAttribute{Key: " ", IsActive: false})
		require.EqualError(t, err, "empty/whitespace attribute key not allowed.")
		require.Empty(t, out)
	})
}

func TestAttr(t *testing.T) {
	t.Run("string values become pair attributes", func(t *testing.T) {
		require.Equal(t, PairAttribute{Key: "class", Value: "box"}, Attr("class", "box"))
	})

	t.Run("bool values become boolean attributes", func(t *testing.T) {
		require.Equal(t, BooleanAttribute{Key: "hidden", IsActive: true}, Attr("hidden", true))
		require.Equal(t, BooleanAttribute{Key: "hidden", IsActive: false}, Attr("hidden", false))
	})

	t.Run("accepts named string and bool types", func(t *testing.T) {
		type cssClass string
		type flag bool
		require.Equal(t, PairAttribute{Key: "class", Value: "box"}, Attr("class", cssClass("box")))
		require.Equal(t, BooleanAttribute{Key: "hidden", IsActive: true}, Attr("hidden", flag(true)))
	})
}

func TestAttrReflect(t *testing.T) {
	t.Run("rejects unsupported value types", func(t *testing.T) {
		require.PanicsWithValue(t, "unexpected value type for attribute", func() {
			attrReflect("key", 42)
		})
	})
}

func TestMakePairAttributeConstructor(t *testing.T) {
	attrOf := MakePairAttributeConstructor("hx-get")
	var attr PairAttribute = attrOf("/api/data") // concrete return type, not Attribute
	require.Equal(t, PairAttribute{Key: "hx-get", Value: "/api/data"}, attr)
	require.Equal(t, PairAttribute{Key: "hx-get", Value: ""}, attrOf(""))
}

func TestMakeBooleanAttributeConstructor(t *testing.T) {
	attrOf := MakeBooleanAttributeConstructor("hx-preserve")
	var attr BooleanAttribute = attrOf(true) // concrete return type, not Attribute
	require.Equal(t, BooleanAttribute{Key: "hx-preserve", IsActive: true}, attr)
	require.Equal(t, BooleanAttribute{Key: "hx-preserve", IsActive: false}, attrOf(false))
}

package hyper

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderNode(t *testing.T) {
	t.Run("writes the node output to the writer", func(t *testing.T) {
		out, err := tryRender(t, DIV("hello"))
		require.NoError(t, err)
		require.Equal(t, "<div>hello</div>", out)
	})

	t.Run("nil node panics instead of rendering", func(t *testing.T) {
		require.PanicsWithValue(t, "nil node passed to RenderNode", func() {
			RenderNode(&bytes.Buffer{}, nil)
		})
	})

	t.Run("propagates errors from the node without writing anything", func(t *testing.T) {
		out, err := tryRender(t, errNode{errBoom})
		require.ErrorIs(t, err, errBoom)
		require.Empty(t, out)
	})

	t.Run("propagates errors from the writer", func(t *testing.T) {
		err := RenderNode(failingWriter{errBoom}, DIV("x"))
		require.ErrorIs(t, err, errBoom)
	})
}

func TestRenderNodeToString(t *testing.T) {
	t.Run("returns the rendered output", func(t *testing.T) {
		html, err := RenderNodeToString(DIV("hi"))
		require.NoError(t, err)
		require.Equal(t, "<div>hi</div>", html)
	})

	t.Run("returns the render error with an empty string", func(t *testing.T) {
		html, err := RenderNodeToString(errNode{errBoom})
		require.ErrorIs(t, err, errBoom)
		require.Empty(t, html)
	})

	t.Run("nil node panics", func(t *testing.T) {
		require.PanicsWithValue(t, "nil node passed to RenderNodeToString", func() {
			_, _ = RenderNodeToString(nil)
		})
	})
}

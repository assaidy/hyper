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

func TestRenderNodeThen(t *testing.T) {
	t.Run("passes the rendered bytes to the callback", func(t *testing.T) {
		var got string
		err := RenderNodeThen(DIV("hi"), func(data []byte) error {
			got = string(data)
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, "<div>hi</div>", got)
	})

	t.Run("returns the callback error", func(t *testing.T) {
		err := RenderNodeThen(DIV("hi"), func([]byte) error { return errBoom })
		require.ErrorIs(t, err, errBoom)
	})

	t.Run("does not call the callback when rendering fails", func(t *testing.T) {
		called := false
		err := RenderNodeThen(errNode{errBoom}, func([]byte) error {
			called = true
			return nil
		})
		require.ErrorIs(t, err, errBoom)
		require.False(t, called)
	})

	t.Run("nil node panics before the callback", func(t *testing.T) {
		called := false
		require.PanicsWithValue(t, "nil node passed to RenderNodeThen", func() {
			RenderNodeThen(nil, func([]byte) error {
				called = true
				return nil
			})
		})
		require.False(t, called)
	})
}

package hyper

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// errBoom is the sentinel error returned by errNode and failingWriter.
var errBoom = errors.New("boom")

// errNode is a HyperNode that always fails, used to exercise error
// propagation through the render pipeline.
type errNode struct{ err error }

func (e errNode) RenderNode(io.Writer) error { return e.err }

// failingWriter is an io.Writer that always fails, used to exercise
// write-error propagation.
type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

// countingNode counts how often it is rendered and writes "counted".
type countingNode struct{ calls int }

func (n *countingNode) RenderNode(w io.Writer) error {
	n.calls++
	_, err := io.WriteString(w, "counted")
	return err
}

// tryRender renders node and returns both the output and any error.
func tryRender(t *testing.T, node HyperNode) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	err := RenderNode(&buf, node)
	return buf.String(), err
}

// mustRender renders node, failing the test on error.
func mustRender(t *testing.T, node HyperNode) string {
	t.Helper()
	out, err := tryRender(t, node)
	require.NoError(t, err)
	return out
}

// renderAttribute renders a single attribute and returns (output, error).
func renderAttribute(t *testing.T, attr Attribute) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	err := attr.RenderAttribute(&buf)
	return buf.String(), err
}

// resetOnceCache wipes the package-global once cache so the Once/OnceKey
// tests are independent of execution order and reruns (e.g. go test -count=2).
func resetOnceCache(t *testing.T) {
	t.Helper()
	onceCache = sync.Map{}
	t.Cleanup(func() { onceCache = sync.Map{} })
}

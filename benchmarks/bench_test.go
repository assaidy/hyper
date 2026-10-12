package benchmarks

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	. "github.com/assaidy/hyper/v4"
)

// ---- render helpers (not timed) ----

func mustHyperString(t *testing.T, n HyperNode) string {
	t.Helper()
	s, err := RenderNodeToString(n)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func mustTemplString(t *testing.T, c interface {
	Render(context.Context, io.Writer) error
},
) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// ---- parity: all three engines must produce identical HTML ----

func TestParity(t *testing.T) {
	d := fixtureDashboard()

	cases := []struct {
		name     string
		hyper    string
		template string
		templ    string
	}{
		{
			"simple",
			mustHyperString(t, hyperSimple(benchSimpleName)),
			renderTemplateToString(t, tmplSimple, benchSimpleName),
			mustTemplString(t, Simple(benchSimpleName)),
		},
		{
			"dashboard",
			mustHyperString(t, hyperDashboard(d)),
			renderTemplateToString(t, tmplDashboard, d),
			mustTemplString(t, Dashboard(d)),
		},
		{
			"dashboard-guest-hidden-error",
			mustHyperString(t, hyperDashboard(DashboardData{
				IsAdmin: false, HasError: true, IsHidden: true,
				Users: []string{"X <&> Y"},
			})),
			renderTemplateToString(t, tmplDashboard, DashboardData{
				IsAdmin: false, HasError: true, IsHidden: true,
				Users: []string{"X <&> Y"},
			}),
			mustTemplString(t, Dashboard(DashboardData{
				IsAdmin: false, HasError: true, IsHidden: true,
				Users: []string{"X <&> Y"},
			})),
		},
		{
			"list",
			mustHyperString(t, hyperUserList([]string{"Alice", "Bob"})),
			renderTemplateToString(t, tmplUserList, []string{"Alice", "Bob"}),
			mustTemplString(t, UserList([]string{"Alice", "Bob"})),
		},
		{
			"escape",
			mustHyperString(t, hyperEscape(benchEscapeText)),
			renderTemplateToString(t, tmplEscape, benchEscapeText),
			mustTemplString(t, Escape(benchEscapeText)),
		},
		{
			"once",
			mustHyperString(t, hyperOncePage()),
			renderTemplateToString(t, tmplOncePage, nil),
			mustTemplString(t, OncePage()),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.hyper != c.template {
				t.Errorf("hyper vs html/template mismatch:\n hyper:    %q\n template: %q", c.hyper, c.template)
			}
			if c.hyper != c.templ {
				t.Errorf("hyper vs templ mismatch:\n hyper: %q\n templ: %q", c.hyper, c.templ)
			}
		})
	}
}

func renderTemplateToString(t *testing.T, tmpl interface {
	Execute(io.Writer, any) error
}, data any,
) string {
	t.Helper()
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// ---- benchmark harness ----

// discard render: pure generation cost, no result retention.
func benchHyper(b *testing.B, n HyperNode) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RenderNode(io.Discard, n); err != nil {
			b.Fatal(err)
		}
	}
}

func benchTemplate(b *testing.B, tmpl interface {
	Execute(io.Writer, any) error
}, data any,
) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := tmpl.Execute(io.Discard, data); err != nil {
			b.Fatal(err)
		}
	}
}

func benchTempl(b *testing.B, c interface {
	Render(context.Context, io.Writer) error
},
) {
	b.Helper()
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := c.Render(ctx, io.Discard); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- Simple ----

func BenchmarkSimple_Hyper(b *testing.B) { benchHyper(b, hyperSimple(benchSimpleName)) }

func BenchmarkSimple_HyperBuildRender(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RenderNode(io.Discard, hyperSimple(benchSimpleName)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSimple_Template(b *testing.B) { benchTemplate(b, tmplSimple, benchSimpleName) }

func BenchmarkSimple_Templ(b *testing.B) { benchTempl(b, Simple(benchSimpleName)) }

// ---- Dashboard ----

func BenchmarkDashboard_Hyper(b *testing.B) {
	d := fixtureDashboard()
	benchHyper(b, hyperDashboard(d))
}

func BenchmarkDashboard_HyperBuildRender(b *testing.B) {
	d := fixtureDashboard()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RenderNode(io.Discard, hyperDashboard(d)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDashboard_Template(b *testing.B) {
	benchTemplate(b, tmplDashboard, fixtureDashboard())
}

func BenchmarkDashboard_Templ(b *testing.B) { benchTempl(b, Dashboard(fixtureDashboard())) }

// ---- List100 ----

func BenchmarkList100_Hyper(b *testing.B) { benchHyper(b, hyperUserList(benchUsers100)) }

func BenchmarkList100_HyperBuildRender(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RenderNode(io.Discard, hyperUserList(benchUsers100)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkList100_Template(b *testing.B) { benchTemplate(b, tmplUserList, benchUsers100) }

func BenchmarkList100_Templ(b *testing.B) { benchTempl(b, UserList(benchUsers100)) }

// ---- List1000 ----

func BenchmarkList1000_Hyper(b *testing.B) { benchHyper(b, hyperUserList(benchUsers1000)) }

func BenchmarkList1000_HyperBuildRender(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RenderNode(io.Discard, hyperUserList(benchUsers1000)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkList1000_Template(b *testing.B) { benchTemplate(b, tmplUserList, benchUsers1000) }

func BenchmarkList1000_Templ(b *testing.B) { benchTempl(b, UserList(benchUsers1000)) }

// ---- Escape ----

func BenchmarkEscape_Hyper(b *testing.B) { benchHyper(b, hyperEscape(benchEscapeText)) }

func BenchmarkEscape_HyperBuildRender(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RenderNode(io.Discard, hyperEscape(benchEscapeText)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEscape_Template(b *testing.B) { benchTemplate(b, tmplEscape, benchEscapeText) }

func BenchmarkEscape_Templ(b *testing.B) { benchTempl(b, Escape(benchEscapeText)) }

// ---- Once / static page ----

func BenchmarkOnce_HyperNoCache(b *testing.B) { benchHyper(b, hyperOncePage()) }

func BenchmarkOnce_HyperCached(b *testing.B) {
	// Warm the OnceKey cache so the loop measures steady-state replay.
	warm, err := RenderNodeToString(hyperOncePageCached())
	if err != nil || warm == "" {
		b.Fatalf("warmup failed: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RenderNode(io.Discard, hyperOncePageCached()); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkOnce_HyperOnceCached is the Once (PC-keyed) equivalent of
// BenchmarkOnce_HyperCached: node built once, steady-state replay.
func BenchmarkOnce_HyperOnceCached(b *testing.B) {
	n := hyperOncePageAuto()
	warm, err := RenderNodeToString(n)
	if err != nil || warm == "" {
		b.Fatalf("warmup failed: %v", err)
	}
	benchHyper(b, n)
}

// BenchmarkOnce_HyperOnceBuild rebuilds via Once on every iteration:
// measures per-request construction (caller-PC capture) + cached replay.
func BenchmarkOnce_HyperOnceBuild(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RenderNode(io.Discard, Once(hyperOncePage)); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkOnce_HyperOnceKeyBuild rebuilds via OnceKey on every iteration:
// the explicit-key counterpart to BenchmarkOnce_HyperOnceBuild.
func BenchmarkOnce_HyperOnceKeyBuild(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RenderNode(io.Discard, OnceKey("bench-once-page-build", hyperOncePage)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOnce_TemplateExecute(b *testing.B) { benchTemplate(b, tmplOncePage, nil) }

func BenchmarkOnce_TemplateCachedBytes(b *testing.B) {
	var buf bytes.Buffer
	if err := tmplOncePage.Execute(&buf, nil); err != nil {
		b.Fatal(err)
	}
	cached := bytes.Clone(buf.Bytes())
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := io.Discard.Write(cached); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOnce_TemplRender(b *testing.B) { benchTempl(b, OncePage()) }

func BenchmarkOnce_TemplCachedBytes(b *testing.B) {
	var buf bytes.Buffer
	if err := OncePage().Render(context.Background(), &buf); err != nil {
		b.Fatal(err)
	}
	cached := bytes.Clone(buf.Bytes())
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := io.Discard.Write(cached); err != nil {
			b.Fatal(err)
		}
	}
}

// Sanity: cached replay really equals fresh render.
func TestOnceCacheParity(t *testing.T) {
	fresh := mustHyperString(t, hyperOncePage())
	for _, tc := range []struct {
		name string
		node HyperNode
	}{
		{"oncekey", hyperOncePageCached()},
		{"once", hyperOncePageAuto()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := RenderNode(&buf, tc.node); err != nil {
				t.Fatal(err)
			}
			// Second render must be a cache hit with identical bytes.
			buf.Reset()
			if err := RenderNode(&buf, tc.node); err != nil {
				t.Fatal(err)
			}
			if !strings.EqualFold(buf.String(), fresh) {
				t.Errorf("once cache mismatch:\n fresh: %q\n cached: %q", fresh, buf.String())
			}
		})
	}
}

package benchmarks

import (
	. "github.com/assaidy/hyper/v3"
)

// This file holds the hyper implementations. Each builder is called once
// outside the benchmark loop; only RenderNode is timed (except the
// explicit Build+Render variants).

func hyperSimple(name string) HyperNode {
	return DIV(AttrClass("greeting"), "Hello, ", name)
}

func hyperDashboard(d DashboardData) HyperNode {
	return DIV(
		AttrClass("max-w-2xl mx-auto"),
		IfNode(d.IsAdmin,
			DIV(
				AttrClass("bg-blue-50 p-4 rounded mb-4"),
				P(AttrClass("font-semibold"), "Admin Panel"),
			),
		).Else(
			DIV(
				AttrClass("p-4"),
				P("Welcome Guest"),
			),
		),
		P(
			AttrClass(Classes(
				"text-sm",
				If(d.HasError, "text-red-500").Else("text-gray-500"),
				If(d.IsHidden, "hidden").Value(),
			)),
			"Total users: ", len(d.Users),
		),
		UL(
			AttrClass("space-y-2 mt-4"),
			AttrId("users-list"),
			Range(d.Users, func(_ int, name string) any {
				return LI(AttrClass("p-2 bg-white rounded shadow"), name)
			}),
		),
		INPUT(
			AttrType(TypeSubmit),
			AttrDisabled(d.IsAdmin),
			Attr("placeholder", "Search"),
		),
	)
}

func hyperUserList(users []string) HyperNode {
	return UL(
		AttrClass("space-y-2 mt-4"),
		AttrId("users-list"),
		Range(users, func(_ int, name string) any {
			return LI(AttrClass("p-2 bg-white rounded shadow"), name)
		}),
	)
}

func hyperEscape(payload string) HyperNode {
	return P(payload)
}

func hyperOncePage() HyperNode {
	return Fragment(
		DOCTYPE(),
		HTML(
			AttrLang("en"),
			HEAD(
				TITLE("Dashboard"),
			),
			BODY(
				AttrClass("bg-gray-100 p-8"),
				H1(AttrClass("text-3xl font-bold"), "Welcome"),
				P("Static content"),
			),
		),
	)
}

// hyperOncePageCached wraps the static page in OnceKey so re-renders
// replay cached bytes (sync.Map load + write).
func hyperOncePageCached() HyperNode {
	return OnceKey("bench-once-page", hyperOncePage)
}

// hyperOncePageAuto wraps the static page in Once (auto-keyed by caller PC).
func hyperOncePageAuto() HyperNode {
	return Once(hyperOncePage)
}

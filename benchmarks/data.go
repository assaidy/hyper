package benchmarks

import (
	"fmt"

	. "github.com/assaidy/hyper/v4"
)

// DashboardData is the shared fixture for the dashboard scenario.
// Same shape is fed to hyper, html/template and templ.
type DashboardData struct {
	IsAdmin  bool
	HasError bool
	IsHidden bool
	Users    []string
}

// dashboardClass mirrors hyper's Classes(...) computation for use
// by templ (which can't call the generic If-chain inside an attribute
// as ergonomically). Logic is identical: base + error branch +
// optional hidden.
func dashboardClass(d DashboardData) string {
	errCls := "text-gray-500"
	if d.HasError {
		errCls = "text-red-500"
	}
	return Classes("text-sm", errCls, If(d.IsHidden, "hidden").Value())
}

// fixtureDashboard returns a representative dashboard input.
func fixtureDashboard() DashboardData {
	return DashboardData{
		IsAdmin:  true,
		HasError: false,
		IsHidden: false,
		Users:    []string{"Alice", "Bob", "Charlie"},
	}
}

// makeUsers builds a deterministic N-user list: "User 001" ... .
func makeUsers(n int) []string {
	users := make([]string, n)
	for i := range users {
		users[i] = fmt.Sprintf("User %03d", i+1)
	}
	return users
}

var (
	benchSimpleName = "Alice & <Bob>"
	benchEscapeText = `<script>alert('xss')</script> & "quotes" <b>bold</b>`
	benchUsers100   = makeUsers(100)
	benchUsers1000  = makeUsers(1000)
)

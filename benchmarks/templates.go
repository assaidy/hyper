package benchmarks

import (
	"html/template"
)

// html/template sources. Kept compact (no extra whitespace) so output
// is byte-identical to the hyper builders.

const (
	tmplSimpleSrc = `<div class="greeting">Hello, {{.}}</div>`

	tmplDashboardSrc = `<div class="max-w-2xl mx-auto">{{if .IsAdmin}}<div class="bg-blue-50 p-4 rounded mb-4"><p class="font-semibold">Admin Panel</p></div>{{else}}<div class="p-4"><p>Welcome Guest</p></div>{{end}}<p class="text-sm {{if .HasError}}text-red-500{{else}}text-gray-500{{end}}{{if .IsHidden}} hidden{{end}}">Total users: {{len .Users}}</p><ul class="space-y-2 mt-4" id="users-list">{{range .Users}}<li class="p-2 bg-white rounded shadow">{{.}}</li>{{end}}</ul><input type="submit"{{if .IsAdmin}} disabled{{end}} placeholder="Search"></div>`

	tmplUserListSrc = `<ul class="space-y-2 mt-4" id="users-list">{{range .}}<li class="p-2 bg-white rounded shadow">{{.}}</li>{{end}}</ul>`

	tmplEscapeSrc = `<p>{{.}}</p>`

	tmplOncePageSrc = `<!DOCTYPE html><html lang="en"><head><title>Dashboard</title></head><body class="bg-gray-100 p-8"><h1 class="text-3xl font-bold">Welcome</h1><p>Static content</p></body></html>`
)

var (
	tmplSimple    = template.Must(template.New("simple").Parse(tmplSimpleSrc))
	tmplDashboard = template.Must(template.New("dashboard").Parse(tmplDashboardSrc))
	tmplUserList  = template.Must(template.New("userlist").Parse(tmplUserListSrc))
	tmplEscape    = template.Must(template.New("escape").Parse(tmplEscapeSrc))
	tmplOncePage  = template.Must(template.New("once").Parse(tmplOncePageSrc))
)

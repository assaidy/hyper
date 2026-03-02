package h_vs_templ

import "html/template"

var (
	simpleElementTmpl  *template.Template
	deepNestingTmpl    *template.Template
	manyAttributesTmpl *template.Template
	largeTextTmpl      *template.Template
	listTempl10Tmpl    *template.Template
	listTempl100Tmpl   *template.Template
	conditionalsTmpl   *template.Template
	mixedContentTmpl   *template.Template
	voidElementsTmpl   *template.Template
	htmlEscapingTmpl   *template.Template
	tableTmpl          *template.Template
	formTempl          *template.Template
	emptyPageTmpl      *template.Template
	rawTextTmpl        *template.Template
	svgTempl           *template.Template
	realWorldTempl     *template.Template
)

func init() {
	simpleElementTmpl = template.Must(template.New("simpleElement").Parse(`<div>Hello World</div>`))

	deepNestingTmpl = template.Must(template.New("deepNesting").Parse(`<div><div><div><div><div><p>Deep content</p></div></div></div></div></div>`))

	manyAttributesTmpl = template.Must(template.New("manyAttributes").Parse(`<div id="main" class="container wrapper" data-role="content" data-value="12345" aria-label="Main content" hidden></div>`))

	largeTextTmpl = template.Must(template.New("largeText").Parse(`<p>{{.}}</p>`))

	listTempl10Tmpl = template.Must(template.New("list10").Parse(`<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>`))

	listTempl100Tmpl = template.Must(template.New("list100").Parse(`<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>`))

	conditionalsTmpl = template.Must(template.New("conditionals").Parse(`<div>{{if .First}}<span>First</span> {{end}}{{if .Second}}<span>Second</span> {{end}}{{if .Third}}<span>Third</span> {{end}}{{if .True}}<strong>True</strong>{{else}}<em>False</em>{{end}}</div>`))

	mixedContentTmpl = template.Must(template.New("mixedContent").Parse(`<div><h1>Title</h1><p>Paragraph with <strong>bold</strong> and <em>italic</em> text.</p><ul><li>Item 1</li><li><a href="#">Link</a></li></ul><div class="footer"><small>Copyright 2024</small></div></div>`))

	voidElementsTmpl = template.Must(template.New("voidElements").Parse(`<div><img src="image.jpg" alt="Image"><br><hr><input type="text" value="input"><meta charset="UTF-8"><link rel="stylesheet" href="style.css"></div>`))

	htmlEscapingTmpl = template.Must(template.New("htmlEscaping").Parse(`<div>{{.}}</div>`))

	tableTmpl = template.Must(template.New("table").Parse(`<table><thead><tr><th>Name</th><th>Value</th><th>Action</th></tr></thead> <tbody>{{range .}}<tr><td>Cell 1</td><td>Cell 2</td><td><button>Click</button></td></tr>{{end}}</tbody></table>`))

	formTempl = template.Must(template.New("form").Parse(`<form action="/submit" method="POST"><fieldset><legend>User Form</legend> <label for="name">Name:</label> <input type="text" id="name" name="name"><br><label for="email">Email:</label> <input type="email" id="email" name="email"><br><button type="submit">Submit</button></fieldset></form>`))

	emptyPageTmpl = template.Must(template.New("emptyPage").Parse(`<html><body></body></html>`))

	rawTextTmpl = template.Must(template.New("rawText").Parse(`<div>{{.}}</div>`))

	svgTempl = template.Must(template.New("svg").Parse(`<svg width="100" height="100"><circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red"></circle></svg>`))

	realWorldTempl = template.Must(template.New("realWorld").Parse(`<!doctype html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Dashboard - User Management</title><link rel="stylesheet" href="/css/main.css"><link rel="icon" href="/favicon.ico"></head><body><header class="site-header"><nav class="main-nav"><a href="/" class="nav-link">Home</a> <a href="/users" class="nav-link active">Users</a> <a href="/settings" class="nav-link">Settings</a> <a href="/logout" class="nav-link">Logout</a></nav></header><main class="main-content"><h1>User Management Dashboard</h1><p>Welcome to the admin dashboard. Manage users and permissions below.</p>{{if .Users}}<section class="users-section"><h2>Active Users</h2><table class="users-table"><thead><tr><th>ID</th><th>Name</th><th>Role</th><th>Status</th><th>Actions</th></tr></thead> <tbody>{{range .Users}}<tr><td><strong>#</strong></td><td>{{.Name}}</td><td>{{if .Admin}}<span class="badge admin">Administrator</span>{{else}}<span class="badge user">User</span>{{end}}</td><td><span class="status active">Active</span></td><td><button class="btn-edit">Edit</button> <button class="btn-delete">Delete</button></td></tr>{{end}}</tbody></table></section>{{end}}{{if not .Users}}<div class="empty-state"><p>No users found. Add your first user to get started.</p></div>{{end}}<section class="quick-stats"><h3>Quick Stats</h3><div class="stats-grid"><div class="stat-card"><strong>{{len .Users}}</strong> <span>Total Users</span></div><div class="stat-card"><strong>{{if .Users}}{{len .Users}}{{else}}0{{end}}</strong> <span>Active Now</span></div></div></section></main><footer class="site-footer"><p>&copy; 2025 Company Inc. All rights reserved.</p></footer></body></html>`))
}

type ConditionalsData struct {
	First  bool
	Second bool
	Third  bool
	True   bool
}

type RealWorldData struct {
	Users []User
}

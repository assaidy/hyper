package benchmarks

import "html/template"

var (
	simpleElementTempl  *template.Template
	deepNestingTempl    *template.Template
	manyAttributesTempl *template.Template
	largeTextTempl      *template.Template
	listTempl10Templ    *template.Template
	listTempl100Templ   *template.Template
	conditionalsTempl  *template.Template
	mixedContentTempl  *template.Template
	voidElementsTempl  *template.Template
	htmlEscapingTempl  *template.Template
	tableTempl          *template.Template
	formTempl          *template.Template
	emptyPageTempl     *template.Template
	rawTextTempl        *template.Template
	svgTempl           *template.Template
	staticContentTempl  *template.Template
	realWorldTempl     *template.Template
)

func init() {
	simpleElementTempl = template.Must(template.New("simpleElement").Parse(`<div>Hello World</div>`))

	deepNestingTempl = template.Must(template.New("deepNesting").Parse(`<div><div><div><div><div><p>Deep content</p></div></div></div></div></div>`))

	manyAttributesTempl = template.Must(template.New("manyAttributes").Parse(`<div id="main" class="container wrapper" data-role="content" data-value="12345" aria-label="Main content" hidden></div>`))

	largeTextTempl = template.Must(template.New("largeText").Parse(`<p>{{.}}</p>`))

	listTempl10Templ = template.Must(template.New("list10").Parse(`<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>`))

	listTempl100Templ = template.Must(template.New("list100").Parse(`<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>`))

	conditionalsTempl = template.Must(template.New("conditionals").Parse(`<div>{{if .First}}<span>First</span> {{end}}{{if .Second}}<span>Second</span> {{end}}{{if .Third}}<span>Third</span> {{end}}{{if .True}}<strong>True</strong>{{else}}<em>False</em>{{end}}</div>`))

	mixedContentTempl = template.Must(template.New("mixedContent").Parse(`<div><h1>Title</h1><p>Paragraph with <strong>bold</strong> and <em>italic</em> text.</p><ul><li>Item 1</li><li><a href="#">Link</a></li></ul><div class="footer"><small>Copyright 2024</small></div></div>`))

	voidElementsTempl = template.Must(template.New("voidElements").Parse(`<div><img src="image.jpg" alt="Image"><br><hr><input type="text" value="input"><meta charset="UTF-8"><link rel="stylesheet" href="style.css"></div>`))

	htmlEscapingTempl = template.Must(template.New("htmlEscaping").Parse(`<div>{{.}}</div>`))

	tableTempl = template.Must(template.New("table").Parse(`<table><thead><tr><th>Name</th><th>Value</th><th>Action</th></tr></thead> <tbody>{{range .}}<tr><td>Cell 1</td><td>Cell 2</td><td><button>Click</button></td></tr>{{end}}</tbody></table>`))

	formTempl = template.Must(template.New("form").Parse(`<form action="/submit" method="POST"><fieldset><legend>User Form</legend> <label for="name">Name:</label> <input type="text" id="name" name="name"><br><label for="email">Email:</label> <input type="email" id="email" name="email"><br><button type="submit">Submit</button></fieldset></form>`))

	emptyPageTempl = template.Must(template.New("emptyPage").Parse(`<html><body></body></html>`))

	rawTextTempl = template.Must(template.New("rawText").Parse(`<div>{{.}}</div>`))

	svgTempl = template.Must(template.New("svg").Parse(`<svg width="100" height="100"><circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red"></circle></svg>`))

	staticContentTempl = template.Must(template.New("staticContent").Parse(`<header class="site-header"><nav class="main-nav"><a href="/" class="nav-link">Home</a> <a href="/users" class="nav-link">Users</a> <a href="/settings" class="nav-link">Settings</a> <a href="/logout" class="nav-link">Logout</a></nav></header><main class="main-content"><h1>User Management Dashboard</h1><p>Welcome to the admin dashboard. Manage users and permissions below.</p><section class="users-section"><h2>Active Users</h2><table class="users-table"><thead><tr><th>ID</th><th>Name</th><th>Role</th><th>Status</th><th>Actions</th></tr></thead></table></section></main><footer class="site-footer"><p>2025 Company Inc. All rights reserved.</p></footer>`))

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

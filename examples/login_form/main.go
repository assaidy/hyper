package main

import (
	"log"
	"net/http"

	"github.com/assaidy/hyper"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		if err := h.Render(w, loginPage()); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Printf("couldn't render html: %v", err)
		}
	}))

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}

	log.Println("starting server at localhost:8000")
	log.Fatal(server.ListenAndServe())
}

func pageLayout(title string, content h.Node) h.Node {
	return h.Html(
		h.Head(
			h.Meta(h.KV{"charset": "UTF-8"}),
			h.Meta(h.KV{"name": "viewport", "content": "width=device-width, initial-scale=1"}),
			h.Title(title),
		),
		h.Body(content),
	)
}

func loginPage() h.Node {
	return pageLayout("login", h.Empty(
		loginPageStyle(),

		h.Form(h.KV{"method": "post"},
			h.H1("Login"),
			h.Div(
				h.Label("Username:"),
				h.Input(h.KV{"type": "text", "name": "username", "required": true, "placeholder": "Enter your username"}),
			),
			h.Div(
				h.Label("Password:"),
				h.Input(h.KV{"type": "password", "name": "password", "required": true, "placeholder": "Enter your password"}),
			),
			h.Div(
				h.Button(h.KV{"type": "submit"}, "Login"),
			),
		),
	))
}

func loginPageStyle() h.Node {
	return h.Style(`
			body {
				font-family: Arial, sans-serif;
				background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
				margin: 0;
				padding: 0;
				display: flex;
				justify-content: center;
				align-items: center;
				min-height: 100vh;
			}
			form {
				background: white;
				padding: 2rem;
				border-radius: 10px;
				box-shadow: 0 10px 25px rgba(0,0,0,0.2);
				width: 100%;
				max-width: 400px;
			}
			div {
				margin-bottom: 1rem;
			}
			label {
				display: block;
				margin-bottom: 0.5rem;
				font-weight: bold;
				color: #333;
			}
			input {
				width: 100%;
				padding: 0.75rem;
				border: 2px solid #ddd;
				border-radius: 5px;
				font-size: 1rem;
				box-sizing: border-box;
				transition: border-color 0.3s;
			}
			input:focus {
				outline: none;
				border-color: #667eea;
			}
			button {
				width: 100%;
				padding: 0.75rem;
				background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
				color: white;
				border: none;
				border-radius: 5px;
				font-size: 1rem;
				font-weight: bold;
				cursor: pointer;
				transition: transform 0.2s;
			}
			button:hover {
				transform: translateY(-2px);
			}
	`)
}

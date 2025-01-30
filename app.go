package main

import (
	"fmt"
	"local/router"
	"local/template"
	"net/http"
)

func serve() {

	myRouter := router.New()

	groupRouter := myRouter.Prefix("nested").Prefix("well")

	myRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("home boy"))
	})

	groupRouter.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok buddy"))
	})

	groupRouter.Get("/test-2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok buddy 2"))
	})

	myRouter.Get("/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		queryParameters := r.URL.Query()

		view, err := template.New("home.index", template.TemplateData{
			"test": queryParameters.Get("test"),
		}).Parse()

		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Problem loading the file")
			return
		}

		fmt.Fprintf(w, "%s", view)
	})

	myRouter.Serve(router.RouterServeOptions{
		PublicDir: false,
	})
}

func main() {
	serve()
}

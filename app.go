package main

import (
	"fmt"
	"local/template"
	"log"
	"net/http"
)

func serve() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		queryParameters := r.URL.Query()

		view, err := template.New("home.index", template.TemplateData{
			"test": queryParameters.Get("test"),
		}).Parse()

		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Problem loading the file")
		}

		fmt.Fprintf(w, "%s", view)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	serve()
}

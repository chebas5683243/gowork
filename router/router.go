package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Router struct {
	Path string
}

type RouterHandler = func(http.ResponseWriter, *http.Request)

type RouterServeOptions struct {
	PublicDir bool
}

func New() *Router {
	return &Router{
		Path: "/",
	}
}

func (router *Router) Prefix(prefix string) *Router {
	return &Router{
		Path: router.Path + "/" + prefix,
	}
}

func (router *Router) Get(resource string, callback RouterHandler) {
	fullResource := router.Path + resource
	fullResource = strings.ReplaceAll(fullResource, "//", "/")

	fmt.Println(fullResource)
	http.HandleFunc(fullResource, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			callback(w, r)
		}
	})
}

func (router *Router) Serve(opts RouterServeOptions) {
	if opts.PublicDir {
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

// Simple function handler
func funcHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from funcHandler")
}

// Complex handler
type appHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		default:
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func myAppHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Fprintf(w, "Hello from myAppHandler")
	return http.StatusOK, nil
}

// Complex handler with context
type appContext struct {
	greeting string
}

type appCtxHandler struct {
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func (ah appCtxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := ah.h(ah.appContext, w, r); err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		default:
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func myAppCtxHandler(ctx *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Fprintf(w, "Hello from %s", ctx.greeting)
	return http.StatusOK, nil
}

func main() {
	port := flag.String("port", "8000", "Port to listen on.")

	// Simple function
	http.HandleFunc("/func", funcHandler)

	// Handler with refactored error handling
	http.Handle("/app", appHandler(myAppHandler))

	// Handler with complex context
	context := &appContext{greeting: "Context"}
	http.Handle("/ctx", appCtxHandler{context, myAppCtxHandler})

	log.Printf("Listening on :%s\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}

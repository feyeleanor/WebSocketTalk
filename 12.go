package main
import "fmt"
import "html/template"
import "io"
import "net/http"
import "os"
import "strings"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

var VERSION, ADDRESS string

func init() {
	s := strings.Split(os.Args[0], "/")
	VERSION = s[len(s) - 1]

	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}
}

type WebHandler func(http.ResponseWriter, *http.Request)
type CallBridge map[template.JS] WebHandler
type PageConfiguration struct {
	CallBridge
}
type Template interface {
	Execute(io.Writer, interface{}) error
}

func main() {
	p := PageConfiguration{
		CallBridge {
			"A": AJAX_handler("A"),
			"B": AJAX_handler("B"),
			"C": AJAX_handler("C"),
		},
	}

	html, e := template.ParseFiles(VERSION + ".html")
	Abort(FILE_READ, e)

	js, e := template.ParseFiles(VERSION + ".js")
	Abort(FILE_READ, e)

	http.HandleFunc("/", ServeTemplate(html, "text/html", Tap(p)))
	http.HandleFunc("/js", ServeTemplate(js, "application/javascript", Tap(p)))

 	for c, f := range p.CallBridge {
		http.HandleFunc("/" + string(c), f)
	}
	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(n)
	}
}

func Tap(v interface{}) func() interface{} {
	return func() interface{} {
		return v
	}
}

func AJAX_handler(c string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		switch r.Method {
		case "GET":
			fmt.Fprintf(w, "GET (%v) %v", c, r.URL)
		case "POST":
			r.ParseForm()
			fmt.Fprintf(w, "POST (%v) %v {%v}", c, r.URL, r.Form)
		}
	}
}

func ServeTemplate(t Template, mime_type string, f func() interface{}) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime_type)
		Abort(BAD_TEMPLATE, t.Execute(w, f()))
	}
}
package main
import "fmt"
import "html/template"
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

type Handler func(http.ResponseWriter, *http.Request)
type CallBridge map[template.JS] Handler
type PageConfiguration struct {
	Version string
	CallBridge
}

func main() {
	p :=  PageConfiguration{ VERSION, CallBridge {
			"A": AJAX_handler("A"),
			"B": AJAX_handler("B"),
			"C": AJAX_handler("C"),
		},
	}

	html := LoadTemplate(VERSION + ".html")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		Abort(BAD_TEMPLATE, html.Execute(w, p))
	})

	js_file := VERSION + ".js"
	js := LoadTemplate(js_file)
	http.HandleFunc("/" + js_file, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		Abort(BAD_TEMPLATE, js.Execute(w, p))
	})

 	for c, f := range p.CallBridge {
		http.HandleFunc(fmt.Sprint("/", c), f)
	}
	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(n)
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

func LoadTemplate(s string) (r *template.Template) {
	var e error
	r, e = template.ParseFiles(s)
	Abort(FILE_READ, e)
	return
}
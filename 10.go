package main
import "fmt"
import "io"
import "net/http"
import "os"
import "strings"
import "text/template"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const ADDRESS = ":3000"

type WebHandler func(http.ResponseWriter, *http.Request)
type Commands map[string] func(http.ResponseWriter, *http.Request)
type PageConfiguration struct {
	Commands
}
type Template interface {
	Execute(io.Writer, interface{}) error
}

func main() {
	html, e := template.ParseFiles(BaseName() + ".html")
	Abort(FILE_READ, e)

	js, e := template.ParseFiles(BaseName() + ".js")
	Abort(FILE_READ, e)

	p := PageConfiguration{
		Commands {
			"A": AJAX_handler("A"),
			"B": AJAX_handler("B"),
			"C": AJAX_handler("C"),
		},
	}

	http.HandleFunc("/", ServeTemplate(html, "text/html", Tap(p)))
	http.HandleFunc("/js", ServeTemplate(js, "application/javascript", Tap(p)))

 	for c, f := range p.Commands {
		http.HandleFunc("/" + c, f)
	}
	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(n)
	}
}

func BaseName() string {
	s := strings.Split(os.Args[0], "/")
	return s[len(s) - 1]	
}

func Tap(v interface{}) func() interface{} {
	return func() interface{} {
		return v
	}
}

func AJAX_handler(c string) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, c)
	}
}

func ServeTemplate(t Template, mime_type string, f func() interface{}) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime_type)
		Abort(BAD_TEMPLATE, t.Execute(w, f()))
	}
}
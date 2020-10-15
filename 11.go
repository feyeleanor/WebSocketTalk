package main
import "encoding/json"
import "fmt"
import "html/template"
import "io"
import "net/http"
import "os"
import "strings"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const ADDRESS = ":3000"

type WebHandler func(http.ResponseWriter, *http.Request)
type CallBridge map[template.JS] WebHandler
type PageConfiguration struct {
	CallBridge
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
		CallBridge {
			"A": AJAX_handler("A"),
			"B": AJAX_handler("B"),
			"C": AJAX_handler("C"),
		},
	}

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
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(c)
		fmt.Fprint(w, string(b))
	}
}

func ServeTemplate(t Template, mime_type string, f func() interface{}) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime_type)
		Abort(BAD_TEMPLATE, t.Execute(w, f()))
	}
}
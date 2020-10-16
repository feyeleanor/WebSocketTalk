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

	http.HandleFunc("/", ServeTemplate(html, "text/html", p))
	http.HandleFunc("/js", ServeTemplate(js, "application/javascript", p))

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

type Message struct {
	Method string
	Command string
	URL string
	Values string
}

func AJAX_handler(c string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			b, _ := json.Marshal(Message { "GET", c, r.URL.String(), "" })
			fmt.Fprint(w, string(b))

		case "POST":
			r.ParseForm()
			b, _ := json.Marshal(Message { "POST", c, r.URL.String(), r.Form.Encode() })
			fmt.Fprintf(w, string(b))
		}
	}
}

func ServeTemplate(t Template, mime_type string, v interface{}) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime_type)
		Abort(BAD_TEMPLATE, t.Execute(w, v))
	}
}
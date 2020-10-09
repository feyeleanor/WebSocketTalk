package main
import "fmt"
import "html/template"
import "net/http"
import "os"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

var ADDRESS string

func init() {
	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}
}

type CallBridge struct {
	JS template.JS
	Go func(http.ResponseWriter, *http.Request)
}

type PageConfiguration struct {
	Commands map[string] CallBridge
}

func main() {
	html, e1 := template.ParseFiles("11.html")
	halt_on_error(FILE_READ, e1)

	js_temp, e2 := template.ParseFiles("11.js")
	halt_on_error(FILE_READ, e2)

	p :=  PageConfiguration{
		map[string] CallBridge {
			"A": CallBridge { "A()", AJAX_handler("A") },
			"B": CallBridge { "B()", AJAX_handler("B") },
			"C": CallBridge { "C()", AJAX_handler("C") },
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		halt_on_error(BAD_TEMPLATE, html.Execute(w, p))
	})

	http.HandleFunc("/10.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		halt_on_error(BAD_TEMPLATE, js_temp.Execute(w, p))
	})

 	for c, f := range p.Commands {
		http.HandleFunc(fmt.Sprint("/", c), f.Go)
	}
	halt_on_error(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func halt_on_error(n int, e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(n)
	}
}

func AJAX_handler(c string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, c)
	}
}
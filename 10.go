package main
import "fmt"
import "text/template"
import "net/http"
import "os"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const RESOURCES = "10"

var ADDRESS string

func init() {
	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}
}

type PageConfiguration struct {
	Commands map[string] func(http.ResponseWriter, *http.Request)
}

func main() {
	html, e1 := template.ParseFiles(fmt.Sprintf("%v.html", RESOURCES))
	halt_on_error(FILE_READ, e1)

	js_temp, e2 := template.ParseFiles(fmt.Sprintf("%v.js", RESOURCES))
	halt_on_error(FILE_READ, e2)

	p :=  PageConfiguration{
		map[string] func(http.ResponseWriter, *http.Request) {
			"A": AJAX_handler("A"),
			"B": AJAX_handler("B"),
			"C": AJAX_handler("C"),
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
		http.HandleFunc(fmt.Sprint("/", c), f)
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
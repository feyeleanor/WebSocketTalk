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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		Abort(BAD_TEMPLATE, html.Execute(w, p))
	})

 	http.HandleFunc("/js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		Abort(BAD_TEMPLATE, js.Execute(w, p))
 	})

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

func AJAX_handler(c string) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, c)
	}
}
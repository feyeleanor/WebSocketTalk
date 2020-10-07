package main
import "fmt"
import "text/template"
import "net/http"
import "os"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

var ADDRESS string
var AJAX_COMMANDS map[string] func(http.ResponseWriter, *http.Request)

func init() {
	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}

	AJAX_COMMANDS = map[string] func(http.ResponseWriter, *http.Request) {
		"A": make_AJAX_command("A"),
		"B": make_AJAX_command("B"),
		"C": make_AJAX_command("C"),
	}
}

type PageContent struct {
	Commands map[string] func(http.ResponseWriter, *http.Request)
}

func main() {
	html, e := template.ParseFiles("10.html")
	halt_on_error(FILE_READ, e)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		halt_on_error(BAD_TEMPLATE, html.Execute(w, PageContent{ AJAX_COMMANDS }))
	})
	for c, f := range AJAX_COMMANDS {
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

func make_AJAX_command(c string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
fmt.Println("calling", c)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, c)
	}
}
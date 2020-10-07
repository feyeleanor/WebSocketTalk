package main
import "fmt"
import "html/template"
import "net/http"
import "os"
import "net/url"

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

type PageContent struct {
	Title, Heading string
	Query url.Values
}

func main() {
	html, e := template.ParseFiles("09.html")
	halt_on_error(FILE_READ, e)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		halt_on_error(BAD_TEMPLATE, html.Execute(w, PageContent{
			"SERVING HTML",
			"PARAMETERS RECEIVED",
			r.URL.Query(),
		}))
	})
	halt_on_error(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func halt_on_error(n int, e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(n)
	}
}

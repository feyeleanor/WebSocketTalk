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

type PageContent struct {
	Title, Heading, Message string
}

func main() {
	html, e := template.ParseFiles("08.html")
	Abort(FILE_READ, e)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		Abort(BAD_TEMPLATE, html.Execute(w, PageContent{
			"SERVING HTML",
			"SERVING HTML",
			r.URL.String(),
		}))
	})
	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
	if e != nil {
    	fmt.Println(e)
		os.Exit(n)
	}
}

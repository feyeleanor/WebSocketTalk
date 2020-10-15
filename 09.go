package main
import "fmt"
import "html/template"
import "net/http"
import "os"
import "net/url"
import "strings"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const ADDRESS = ":3000"

type PageContent struct {
	Title, Heading string
	Query url.Values
}

func main() {
	html, e := template.ParseFiles(BaseName() + ".html")
	Abort(FILE_READ, e)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		Abort(BAD_TEMPLATE, html.Execute(w, PageContent{
			"SERVING HTML",
			"PARAMETERS RECEIVED",
			r.URL.Query(),
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

func BaseName() string {
	s := strings.Split(os.Args[0], "/")
	return s[len(s) - 1]	
}
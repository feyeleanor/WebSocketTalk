package main
import . "fmt"
import "html/template"
import "net/http"
import "os"

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
	html, e := template.ParseFiles("07.html")
	halt_on_error(1, e)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		halt_on_error(2, html.Execute(w, PageContent{
			"SERVING HTML",
			"SERVING HTML",
			r.URL.String(),
		}))
	})
	if e := http.ListenAndServe(ADDRESS, nil); e != nil {
		Println(e)
	}
}

func halt_on_error(n int, e error) {
	if e != nil {
		os.Exit(n)
	}
}

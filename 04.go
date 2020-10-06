package main
import . "fmt"
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

func main() {
	http.HandleFunc("/", print_url)
	if e := http.ListenAndServe(ADDRESS, nil); e != nil {
		Println(e)
	}
}

func print_url(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	Fprint(w, r.URL)
}

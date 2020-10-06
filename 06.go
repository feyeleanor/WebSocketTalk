package main
import . "fmt"
import "io/ioutil"
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
	html, e := ioutil.ReadFile("06.html")
	halt_on_error(1, e)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		Fprint(w, string(html))
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
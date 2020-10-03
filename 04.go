package main
import . "fmt"
import "io/ioutil"
import "net/http"
import "os"

var ADDRESS, MESSAGE string

func init() {
	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}

	html, e := ioutil.ReadFile("04.html")
	halt_on_error(1, e)
	MESSAGE = string(html)
}

func main() {
  http.HandleFunc("/", billboard)
  if e := http.ListenAndServe(ADDRESS, nil); e != nil {
    Println(e)
  }
}

func billboard(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  Fprintf(w, MESSAGE)
}

func halt_on_error(n int, e error) {
    if e != nil {
        os.Exit(n)
    }	
}
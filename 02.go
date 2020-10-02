package main
import . "fmt"
import "net/http"
import "os"

const MESSAGE = "hello world"
var ADDRESS string

func init() {
	if len(os.Args) > 1 {
		ADDRESS = os.Args[1]
	} else {
		ADDRESS = ":3000"
	}
}

func main() {
  http.HandleFunc("/", billboard)
  if e := http.ListenAndServe(ADDRESS, nil); e != nil {
    Println(e)
  }
}

func billboard(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain")
  Fprintf(w, MESSAGE)
}

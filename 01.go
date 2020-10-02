package main
import . "fmt"
import "net/http"

const MESSAGE = "hello world"
const ADDRESS = ":3000"

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

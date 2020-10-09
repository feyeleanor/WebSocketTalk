package main
import "fmt"
import "net/http"
import "os"

const MESSAGE = "hello world"
const LAUNCH_FAILED = 1

var ADDRESS string

func init() {
	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}
}

func main() {
	http.HandleFunc("/", billboard)
	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func billboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, MESSAGE)
}

func Abort(n int, e error) {
    if e != nil {
    	fmt.Println(e)
        os.Exit(n)
    }	
}
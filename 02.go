package main
import "fmt"
import "net/http"

const MESSAGE = "hello world"
const ADDRESS = ":3000"
const LAUNCH_FAILED = 1

func main() {
	http.HandleFunc("/", billboard)
	halt_on_error(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func billboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, MESSAGE)
}

func halt_on_error(n int, e error) {
    if e != nil {
    	fmt.Println(e)
        os.Exit(n)
    }	
}
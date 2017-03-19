package main

import (
	"net/http"
	"io"
	"fmt"
	"gopkg.in/couchbase/gocb.v1"
	"github.com/gorilla/mux"
	"os"
)

type User struct {
	Id string `json:"uid"`
	Email string `json:"email"`
	Location []string `json:"location"`
}

type thedoc map[string]interface{}
type curDoc string
type theid string
var err error
//var bucket *gocb.Bucket
var cluster, _ = gocb.Connect("couchbase://192.168.61.101")
var bucket, _ = cluster.OpenBucket("travel-sample", "")

func main() {

	rtr := mux.NewRouter()  //.StrictSlash(true)
	rtr.HandleFunc("/hello", hello)
	//rtr.HandleFunc("/POSTit/{id}/{doc}", postit) //.Methods("POST")
	rtr.HandleFunc("/GETit/{id}", getit) //.Methods("GET")

	http.ListenAndServe(":8888", rtr)
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func getit(w http.ResponseWriter, r1 *http.Request) {
	var1 := mux.Vars(r1)
	theid := var1["id"]

	// get the doc from bucket
	_, err = bucket.Get(theid, &thedoc{})
	if err != nil {
		fmt.Println(err)
		os.Exit(51)
	}

	//var tmp map[string]interface{}
	//_ = json.NewDecoder(thedoc{}).Decode(&tmp)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "GotDoc")
	fmt.Fprintln(w, thedoc{})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

}

/*func postit(w http.ResponseWriter, r1 *http.Request, r2 *http.Request) {
	theid := r1.Body
	json.NewEncoder(r2).Encode((&thedoc{}))
	//post the doc to bucket
	bucket.Upsert(theid(), thedoc{}, 0)
	fmt.Fprintln(w, "PutDoc")
}

/* Some notes for later ...
func someDoc(bucket *gocb.Bucket) {
	var theDoc map[string]interface{}
	tmpDoc, err := os.OpenFile("/Users/justin/GOstuff/a.json", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(50)
	}
	err = json.NewDecoder(tmpDoc).Decode(&theDoc)
	if err != nil {
		fmt.Println(err)
		os.Exit(54)
	}

	//create a key to use
	tmpA, err := exec.Command("uuidgen").Output()
	if err != nil {
		panic(err)
		os.Exit(77)
	}
	akey := ("u::" + string(tmpA))

	bucket.Upsert(akey,
		User{
			Id: "someone",
			Email: "someone@this.com",
			Location: []string{"CA", "VA"},
		}, 0)

	// Get the user back
	var getUser User
	bucket.Get(akey, &getUser)
	fmt.Printf("User: %v\n", getUser)

	// Use query
	query := gocb.NewN1qlQuery("SELECT * FROM travel-sample WHERE $1 IN location")
	rows, _ := bucket.ExecuteN1qlQuery(query, []interface{}{"VA"})
	var row interface{}
	for rows.Next(&row) {
		fmt.Printf("Row: %v", row)
	}
}
*/


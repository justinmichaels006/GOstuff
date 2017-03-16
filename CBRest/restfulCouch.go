package CBRest

import (
	"github.com/couchbase/gocb"
	"fmt"
	"os"
	"encoding/json"
	"net/http"
	"io"
	"os/exec"
)

type User struct {
	Id string `json:"uid"`
	Email string `json:"email"`
	Location []string `json:"location"`
}

type thedoc map[string]interface{}
type theid string
var cluster gocb
var bucket gocb.Bucket

func main() {
	var err error
	cluster, err = gocb.Connect("couchbase://192.168.61.101")
	if err != nil {
		fmt.Println(err)
		os.Exit(22)
	}

	bucket, err = cluster.OpenBucket("travel-sample", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(23)
	}

	go someDoc(bucket)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/POSTit/{id}/{doc}", postit)
	mux.HandleFunc("/GETit/{id}", getit)
	http.ListenAndServe(":8888", mux)
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func postit(w http.ResponseWriter, r1 *http.Request, r2 *http.Request) {
	theid := r1.GetBody
	thedoc := r2.GetBody

	for theid {
		json.NewEncoder(w).Encode((&thedoc()))
		//post the doc to bucket
	}

}

func getit(w http.ResponseWriter, r1 *http.Request) {
	theid := r1.GetBody

	for theid {
		// get the doc from bucket
		json.NewEncoder(w).Encode(&thedoc{})
	}

}

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

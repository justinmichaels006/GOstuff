package main

import (
	"net/http"
	"io"
	"fmt"
	"github.com/couchbase/gocb"
	"github.com/gorilla/mux"
	"os"
	"encoding/json"
	"sync"
	"strconv"
	"time"
)

type User struct {
	Id string `json:"uid"`
	Email string `json:"email"`
	Location []string `json:"location"`
}

var thedoc map[string]interface{}
var err error
//Establish the environment to use
var cluster, _ = gocb.Connect("couchbase://192.168.61.102/")
var bucket, _ = cluster.OpenBucket("testload", "testload")
var waitGroup sync.WaitGroup
//Number is the number of bulk operations to perform
var number int = 10
var jsonset map[string]interface{}
var jsonget map[string]interface{}

func testit(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "(Y)")
}

func closeBucket(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Closed Bucket and Shutdown")
	bucket.Close()
	os.Exit(88)
}

func getit(w http.ResponseWriter, r1 *http.Request) {
	var1 := mux.Vars(r1)
	theid := var1["id"]

	// get the doc from bucket
	_, err = bucket.Get(theid, &thedoc)
	if err != nil {
		fmt.Println(err.Error())
		io.WriteString(w, err.Error())
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&thedoc)
}

func postit(w http.ResponseWriter, r2 *http.Request) {
	var2 := mux.Vars(r2)
	theid := var2["id"]
	thedoc := var2["doc"]

	fmt.Fprintln(w, thedoc)
	//post the doc to bucket
	_, err = bucket.Upsert(theid, thedoc, 0)
	if err != nil {
		fmt.Println(err.Error())
		io.WriteString(w, err.Error())
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&thedoc)
}

func bulkloader(w http.ResponseWriter, r *http.Request) {

	fmt.Println("loop check")
	// Retrieve Document
	_, err = bucket.Get("dummy", &jsonset)
	if err != nil {
		fmt.Println("Doh!")
	}

	waitGroup.Add(number)
	// Channel our data so it can be picked up by an available goroutine
	for i := 0; i < number; i++ {
		go aworker(strconv.Itoa(i), jsonset)
	}

	waitGroup.Wait()

	fmt.Println("The loader application has completed!")
}

func aworker(id string, adoc map[string]interface{}) {
	defer waitGroup.Done()
	fmt.Println("debug check")
	// Create the document
	adoc["upStamp"] = nil
	bucket.Upsert(id, adoc, 0)

	// Append to an array or create if it doesn't exist, using a subdocument operation
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	bucket.MutateIn(id, 0, 0).Upsert("upStamp", millis, true).Execute()

}

func bulkreader(w http.ResponseWriter, r *http.Request) {

	waitGroup.Add(number)
	// Channel our data so it can be picked up by an available goroutine
	for i := 0; i < number; i++ {
		go bworker(i, jsonget)
	}

	waitGroup.Wait()

	fmt.Println("The reader application has completed!")
}

func bworker(thenum int, bdoc map[string]interface{}) {

	defer waitGroup.Done()
	fmt.Println("debug check")

	id := strconv.Itoa(thenum)

	for {
		// Retrieve Document
		_, err = bucket.Get(id, &bdoc)
		if err != nil {
			fmt.Println("Doh! ", id)
			//time.Sleep(100000)
			//go bworker(thenum, bdoc)
		} else {
			now := time.Now()
			nanos := now.UnixNano()
			millis := nanos / 1000000
			bucket.MutateIn(id, 0, 0).Upsert("getStamp", millis, true).Execute()
			fmt.Println("Got It:", id)
			break
		}
	}
}

func main() {

	rtr := mux.NewRouter()  //?.StrictSlash(true)
	rtr.HandleFunc("/test", testit)
	rtr.HandleFunc("/done", closeBucket)
	rtr.HandleFunc("/POSTit/{id}/{doc}", postit) //?.Methods("POST")
	rtr.HandleFunc("/GETit/{id}", getit) //?.Methods("GET")
	rtr.HandleFunc("/POSTbulk/", bulkloader)
	rtr.HandleFunc("/GETbulk/", bulkreader)

	http.ListenAndServe(":8888", rtr)
}



package main

import (
	"net/http"
	"io"
	"fmt"
	"gopkg.in/couchbase/gocb.v1"
	"github.com/gorilla/mux"
	"os"
	"encoding/json"
)

type User struct {
	Id string `json:"uid"`
	Email string `json:"email"`
	Location []string `json:"location"`
}

var thedoc map[string]interface{}
var err error
var cluster, _ = gocb.Connect("couchbase://10.112.151.101")
var bucket, _ = cluster.OpenBucket("beers", "")

func main() {

	rtr := mux.NewRouter()  //?.StrictSlash(true)
	rtr.HandleFunc("/done", closeBucket)
	rtr.HandleFunc("/POSTit/{id}/{doc}", postit) //?.Methods("POST")
	rtr.HandleFunc("/GETit/{id}", getit) //?.Methods("GET")

	http.ListenAndServe(":8888", rtr)
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



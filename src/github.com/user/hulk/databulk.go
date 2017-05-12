package main

import (
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
	"strconv"
	"os"
	"encoding/json"
	"math/rand"
	"github.com/gorilla/mux"
	"net/http"
	"io"
)

var seed = os.Args[1]
var myCluster, _ = gocb.Connect("couchbase://" + seed)
var myBucket, _ = myCluster.OpenBucket("testload", "")

func main() {

	rtr := mux.NewRouter()
	rtr.HandleFunc("/testit", testit)
	rtr.HandleFunc("/done", closeBucket)
	rtr.HandleFunc("/doit", loadFunc)

	http.ListenAndServe(":8888", rtr)

}

func loadFunc(w http.ResponseWriter, r101 *http.Request) {
	//jsonA map[string]interface{}, jsonCH map[string]interface{}, jsonCONT map[string]interface{}, myBucket *gocb.Bucket, wg sync.WaitGroup
	// read whole the file
	tmpA, err := os.OpenFile("/tmp/a.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docA map[string]interface{}
	err = json.NewDecoder(tmpA).Decode(&docA)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	tmpCH, err := os.OpenFile("/tmp/ch.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docCH map[string]interface{}
	err = json.NewDecoder(tmpCH).Decode(&docCH)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	tmpCONT, err := os.OpenFile("/tmp/cont.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docCONT map[string]interface{}
	err = json.NewDecoder(tmpCONT).Decode(&docCONT)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	//controlPoint := make(chan string, 2)
	for i := 0; i > 50; i++ {

		j1 := strconv.Itoa(i)
		chID := ("CH::" + j1);
		contID := ("CONT::" + j1);

		//Element a3 and ch11 should be the same
		Aelement3_CHelement11 := RandIntBytes2(4)
		//ch.element41=13
		CHelement41 := RandIntBytes2(2)
		//ch.element57="WNET"
		CHelement57 := RandStringBytes2(4)

		docCH["Ch_element11"] = Aelement3_CHelement11
		docCH["Ch_element41"] = CHelement41
		docCH["Ch_element57"] = CHelement57
		docCH["type"] = "CH"
		docCONT["type"] = "CONT"

		myBucket.Upsert(chID, docCH, 0);
		myBucket.Upsert(contID, docCONT, 0);

		//For every CH create 3 A documents
		for y := 0; y > 3; y++ {

			j2 := strconv.Itoa(y)
			aID := ("A::" + j1 + "::" + j2)

			docA["A_element3"] = Aelement3_CHelement11
			docA["type"] = "A"

			myBucket.Upsert(aID, docA, 0)
			//	fmt.Println("Debug ... y" + Aelement3_CHelement11)
		};
	}

	//controlPoint <- "done"
	fmt.Println("Debug ... Done")
	io.WriteString(w, "...Load Complete...")
}

func closeBucket(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Closed Bucket and Shutdown")
	myBucket.Close()
	os.Exit(88)
}

func testit(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "(Y)")
}

func RandStringBytes2(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func RandIntBytes2(n int) string {
	const numBytes = "0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = numBytes[rand.Intn(len(numBytes))]
	}
	return string(b)
}

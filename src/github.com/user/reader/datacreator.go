package main

import (
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"time"
	"sync"
)

func FlowControl(controller bool, number int, jsonACCT map[string]interface{}, jsonCUST map[string]interface{}, myBucket *gocb.Bucket) int {

	if number == 0 {
		fmt.Println("Done: ", number)
		go finish()
	}

	//go upsertOne(number, jsonACCT, jsonCUST, myBucket)
	str := strconv.Itoa(number)
	fmt.Println("Upsert: ", str)

	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	jsonACCT["upStamp"] = millis
	jsonCUST["upStamp"] = millis
	myBucket.Upsert("ACCT::"+str, jsonACCT, 0)
	myBucket.Upsert("CUST::"+str, jsonCUST, 0)

	return number + FlowControl(controller, number-1, jsonACCT, jsonCUST, myBucket)
}

/*func upsertOne(num int, jsonACCT map[string]interface{}, jsonCUST map[string]interface{}, myBucket *gocb.Bucket) {

	// Upsert One Document
	str := strconv.Itoa(num)
	fmt.Println("Upsert: ", str)
	myBucket.Upsert("ACCT::"+str, jsonACCT, 0)
	myBucket.Upsert("CUST::"+str, jsonCUST, 0)
}*/

func main() {

	// Configuration
	//var getSetPercentage = 0.99;
	//var totalDocs = 10000;
	var currentGroup = 0
	var opsGroups = 50000
	var runningLoad = false
	var seedNode string
	var wg sync.WaitGroup
	// holds the arguments for Couchbase seed node
	seedNode = ("couchbase://" + os.Args[1])
	//seedNode = ("couchbase://192.168.61.101")

	// Connect to Couchbase
	myC, _ := gocb.Connect(seedNode)
	myB, _ := myC.OpenBucket("test", "")

	// read whole the file
	tmpACCT, err := os.OpenFile("/tmp/a.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docACCT map[string]interface{}
	err = json.NewDecoder(tmpACCT).Decode(&docACCT)
	if err != nil {
		fmt.Println(err)
		os.Exit(52)
	}

	tmpCUST, err := os.OpenFile("/tmp/ch.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docCUST map[string]interface{}
	err = json.NewDecoder(tmpCUST).Decode(&docCUST)
	if err != nil {
		fmt.Println(err)
		os.Exit(53)
	}

	if (runningLoad == false && currentGroup == 0) {
		FlowControl(runningLoad, opsGroups, docACCT, docCUST, myB)
	}

	wg.Wait()
}

// Shut down.
func finish() {
	fmt.Println("Good Bye")
	os.Exit(101)
}



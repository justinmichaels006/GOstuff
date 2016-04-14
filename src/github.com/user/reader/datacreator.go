package main

import (
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"time"
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
	var currentGroup = 0;
	var opsGroups = 50;
	var runningLoad = false;
	var seedNode string
	// holds the arguments for Couchbase seed node
	seedNode = ("couchbase://" + os.Args[1])

	// Connect to Couchbase
	myC, _ := gocb.Connect(seedNode)
	myB, _ := myC.OpenBucket("testload", "")

	// read whole the file
	tmpACCT, err := os.OpenFile("/tmp/ACCT.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docACCT map[string]interface{}
	err = json.NewDecoder(tmpACCT).Decode(&docACCT)
	if err != nil {
		fmt.Println(err)
		os.Exit(52)
	}

	tmpCUST, err := os.OpenFile("/tmp/CUST.json", os.O_RDONLY, 0644)
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
		// flowOutput :=
		// fmt.Println("Recursive: ", flowOutput)
	}
}

// Shut down.
func finish() {
	fmt.Println("Good Bye")
	os.Exit(101)
}



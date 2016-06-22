package main

import (
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"math/rand"
	"github.com/nu7hatch/gouuid"
	"time"
)

func main() {

	// Configuration
	// 1 Customer -> 800 to 1200 Groups -> 40 to 80 Devices -> 70 to 200 App
	// 10,000 Customers -> 1M Groups -> 50M Devices -> 2B App
	//TODO: Create ranges during simulation
	var cusomterTotal = 100 //10000;
	var groupTotal = cusomterTotal * 1000
	var deviceTotal = 50
	var appTotal = 200
	var appCatalog = 2000
	var controller = false
	fControl := make(chan string)
	var seedNode string
	var itemApp []gocb.BulkOp

	// holds the arguments for Couchbase seed node
	//seedNode = ("couchbase://" + os.Args[1])
	seedNode = ("couchbase://" + "192.168.61.101")

	// Connect to Couchbase
	myC, _ := gocb.Connect(seedNode)
	myB, _ := myC.OpenBucket("testload", "")

	// Read the Group file
	tmpGROUP, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/GROUP.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
		os.Exit(50)
	}
	var docGROUP map[string]interface{}
	err = json.NewDecoder(tmpGROUP).Decode(&docGROUP)
	if err != nil {
		fmt.Println(err)
		os.Exit(51)
	}

	tmpDEVICE, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/DEVICE.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
		os.Exit(52)
	}
	var docDEVICE map[string]interface{}
	err = json.NewDecoder(tmpDEVICE).Decode(&docDEVICE)
	if err != nil {
		fmt.Println(err)
		os.Exit(53)
	}

	tmpAPP, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/APP.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docAPP map[string]interface{}
	err = json.NewDecoder(tmpAPP).Decode(&docAPP)
	if err != nil {
		fmt.Println(err)
		os.Exit(54)
	}

	var docCUST map[string]interface{}
	str := `{"TYPE": "CUSTONMER", "ID": "uuid", "NUM": "x"}`
	json.Unmarshal([]byte(str), &docCUST)

	//Create Simulated App Catalog
	for y := 0; y < appCatalog; y++ {
		docAPP["TYPE"] = "APP"
		itemApp = append(itemApp, &gocb.InsertOp{Key: "APP::" + strconv.Itoa(y), Value: &docAPP})
	}
	// Perform the bulk operation
	err = myB.Do(itemApp)
	if err != nil {
		fmt.Println("ERRROR PERFORMING CATALOG INSERT:", err)
	}

	//Create The Customer
	go FlowControl(controller, cusomterTotal, "CUSTOMER", docCUST, myB, fControl)

	//TODO: Add the catalog_app_version document

	if fControl == true {
		finish2(myB)
	}
}

//Customer (controller, cusomterTotal, theType, docGROUP, myB, fControl)
func FlowControl(controller bool, number int, theTYPE string,
	jsonDOC map[string]interface{}, myBucket *gocb.Bucket, msg chan string) int {

	if number == 0 {
		fmt.Println("Done: ", number)
		msg := <- "DONE"
		return msg
	}

	// Create the Customer that will anchor the rest of the relationships
	uuid, err := newUUID2()
	if err != nil {
		fmt.Println("uuid issue")
	}
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	theID := uuid

	jsonDOC["upStamp"] = millis
	jsonDOC["CUST_id"] = uuid
	jsonDOC["CUST_num"] = strconv.Itoa(number)
	jsonDOC["TYPE"] = theTYPE
	myBucket.Upsert(theID, jsonDOC, 0)

	//Group (num uuid.UUID, number int, theTYPE string,jsonDOC map[string]interface{}, myBucket *gocb.Bucket, msg chan string)
	return number + FlowGroups(uuid, 1000, "GROUP", , myBucket, msg)
}


//Groups
func FlowGroups (num uuid.UUID, number int, theTYPE string,
	jsonDOC map[string]interface{}, myBucket *gocb.Bucket, msg chan string) int {

	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000

	jsonDOC["upStamp"] = millis
	jsonDOC["CUST_id"] = uuid
	jsonDOC["CUST_num"] = strconv.Itoa(number)
	jsonDOC["TYPE"] = theTYPE
	docGROUP["TYPE"] = theType
	go FlowControl(controller, groupTotal, theType, docGROUP, myB, fControl)
	if fControl == true {

	}
	return
}

//Devices
func FlowDevices (num uuid.UUID, number int, theTYPE string,
jsonDOC map[string]interface{}, myBucket *gocb.Bucket, msg chan string) int {
	if
	//var appArray [appTotal]string
	appArray := make([]string, number)
	//fmt.Printf("%s\n", uuid) // Debug
	for m := 0; m < number; m++ {
		var a = rand.Intn(number)
		// fmt.Println(a) //debug
		appArray[m] = "APP::" + strconv.Itoa(a)
	}
	jsonDOC["APP_install"] = appArray

	return
}


// newUUID generates a random UUID according to RFC 4122
func newUUID2() (string, error) {
	/*This won't work on Windows
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}*/
	out, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return out.String(), err
	}
	//fmt.Println(out)
	return out.String(), err
}

// Shut down.
func finish2(cBucket *gocb.Bucket) {
	cBucket.Close()
	fmt.Println("Good Bye")
	os.Exit(101)
}

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

// newUUID generates a random UUID according to RFC 4122
func newUUID3() (string, error) {
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
func finish3(cBucket *gocb.Bucket) {
	cBucket.Close()
	fmt.Println("Good Bye")
}

func main() {

	// Configuration
	// 1 Customer -> 800 to 1200 Groups -> 40 to 80 Devices -> 70 to 200 App
	// 10,000 Customers -> 1M Groups -> 50M Devices -> 2B App
	//TODO: Create ranges during simulation
	var cusomterTotal = 100 //10000;
	var appCatalog = 2000
	fControl := make(chan bool)
	// Create an Array of BulkOps for the App Catalog
	var itemApp []gocb.BulkOp
	var seedNode string
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

	flowControlCustomer(cusomterTotal, "CUSTOMER", docCUST, docGROUP, docDEVICE, myB)

	<-fControl

}

func flowControlCustomer(customerTotal int, theType string,
docCUST map[string]interface{}, docGROUP map[string]interface{}, docDEVICE map[string]interface{},
myBucket *gocb.Bucket) int {

	if customerTotal == 0 {
		fmt.Println("Last Customer: ", customerTotal)
		return customerTotal //finish3(myBucket)
	}

	// Create the Customer that will anchor the rest of the relationships
	uuid, err := newUUID3()
	if err != nil {
		fmt.Println("UUID Issue")
	}

	docCUST["TYPE"] = "CUSTOMER"
	docCUST["ID"] = "CUSTOMER::" + uuid + "::" + strconv.Itoa(customerTotal)

	myBucket.Upsert(theType + "::" + uuid + "::" + strconv.Itoa(customerTotal), docCUST, 0)

	//number = 1000 GROUPS
	flowControlGroup(uuid, 100, "GROUP", docCUST, docGROUP, docDEVICE, myBucket)
	return customerTotal + flowControlCustomer(customerTotal-1, "CUSTOMER", docCUST, docGROUP, docDEVICE, myBucket)
}

func flowControlGroup(uuid string, number int, theType string,
docCUST map[string]interface{}, docGROUP map[string]interface{}, docDEVICE map[string]interface{},
myBucket *gocb.Bucket) int {

	if number == 0 {
		fmt.Println("Done Group: ", number)
		return number
	}

	//Create Groups
	docGROUP["TYPE"] = "GROUP"
	docGROUP["GROUP_id"] = "GROUP::" + uuid + "::" + strconv.Itoa(number)

	myBucket.Upsert(theType + "::" + uuid + "::" + strconv.Itoa(number), docGROUP, 0)

	//number = 50 DEVICES
	flowControlDevice(uuid, 50, "DEVICE", docCUST, docGROUP, docDEVICE, myBucket)
	return number + flowControlGroup(uuid, number-1, "GROUP", docCUST, docGROUP, docDEVICE, myBucket)
}

func flowControlDevice(uuid string, number int, theType string,
docCUST map[string]interface{}, docGROUP map[string]interface{}, docDEVICE map[string]interface{},
myBucket *gocb.Bucket) int {

	rand.Seed(time.Now().Unix())

	if number == 0 {
		fmt.Println("Done Device: ", number)
		return number
	}

	docDEVICE["TYPE"] = "DEVICE"
	docDEVICE["DEVICE_id"] = "DEVICE::" + uuid + "::" + strconv.Itoa(number)

	//Array of apps from the catalog
	appArray := make([]string, 50)

	for k := 0; k < 50; k++ {
		//match appCatalog 2000
		var a = rand.Intn(2000)
		// fmt.Println(a) //debug
		appArray[k] = "APP::" + strconv.Itoa(a) + "::" + strconv.Itoa(rand.Intn(10))
	}

	docDEVICE["APP_install"] = appArray

	myBucket.Upsert(theType + "::" + uuid + "::" + strconv.Itoa(number), docDEVICE, 0)

	return number + flowControlDevice(uuid, number-1, "DEVICE", docCUST, docGROUP, docDEVICE, myBucket)

}

/*tmpCUST, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/stocks.json", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("error opening file: %v\n",err)
		os.Exit(55)
	}
	sc := bufio.NewScanner(tmpCUST)
	for sc.Scan() {
		fmt.Println(sc.Text())
	}
	if err := sc.Err(); err != nil {
		fmt.Printf("error opening file: %v\n",err)
		os.Exit(55)
	}*/

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
	var cusomterTotal = 1 //100 //10000;
	var groupTotal = cusomterTotal * 10 //1000
	var deviceTotal = 5 //50
	var appTotal = 2 //200
	var appCatalog = 200 //2000
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

	for x := cusomterTotal; x != 0; x-- {
		// Create the Customer that will anchor the rest of the relationships
		uuid, err := newUUID3()

		docCUST["TYPE"] = "CUSTOMER"
		docCUST["ID"] = "CUSTOMER::" + uuid + "::" + strconv.Itoa(x)

		//(int, uuid, type, jsonDOC, myBucket, chan)
		go FlowControl3(x, uuid, "CUSTOMER", docCUST, myB)
		//FlowControl3(x, uuid, "CUSTOMER", docCUST, myB)

		//Create Groups
		for i := 0; i < groupTotal; i++ {
			docGROUP["TYPE"] = "GROUP"
			docGROUP["GROUP_id"] = "GROUP::" + uuid + "::" + strconv.Itoa(i)

			//(int, uuid, type, jsonDOC, myBucket, chan)
			go FlowControl3(i, uuid, "GROUP", docGROUP, myB)
			//FlowControl3(i, uuid, "GROUP", docGROUP, myB)

			//Create Devices
			for j := 0; j < deviceTotal; j++ {
				docDEVICE["TYPE"] = "DEVICE"
				docDEVICE["DEVICE_id"] = "DEVICE::" + uuid + "::" + strconv.Itoa(j)

				//Array of apps from the catalog
				appArray := make([]string, appTotal)
				if err != nil {
					fmt.Printf("error: %v\n", err)
				}

					for k := 0; k < appTotal; k++ {
						var a = rand.Intn(appCatalog)
						// fmt.Println(a) //debug
						appArray[k] = "APP::" + strconv.Itoa(a)
					}

				docDEVICE["APP_install"] = appArray

				//(int, uuid, type, jsonDOC, myBucket, chan)
				go FlowControl3(j, uuid, "DEVICE", docDEVICE, myB)
				//FlowControl3(j, uuid, "DEVICE", docDEVICE, myB)

				//TODO: Add the catalog_app_version document
			}
		}
	fControl <- true
	}

	if <-fControl {
		fmt.Println("Done")
		go finish3(myB)
	}

}

func FlowControl3(number int, uuid string, theType string,
jsonDOC map[string]interface{}, myBucket *gocb.Bucket) {

	str := strconv.Itoa(number)
	theID := theType + "::" + uuid + "::" + str
	//fmt.Println("Upsert: ", str)

	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	jsonDOC["upStamp"] = millis

	myBucket.Upsert(theID, jsonDOC, 0)

	return
}

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
	os.Exit(101)
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

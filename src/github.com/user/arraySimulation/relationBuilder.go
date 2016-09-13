package main

import (
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
	//"encoding/json"
	"time"
	"math/rand"
	"strconv"
)

var bucket *gocb.Bucket

// Create a struct for strongly typed query results
type getID struct {
	ID string "id"
}

type row interface {}

type theDOC struct {
	Device string `json:"Device"`
	App string `json:"App"`
	Type string `json:"TYPE"`
}

func main() {

	var seedNode2 string

	// holds the arguments for Couchbase seed node
	//seedNode = ("couchbase://" + os.Args[1])
	seedNode2 = ("couchbase://" + "192.168.61.101")

	// Connect to Couchbase
	myCo, _ := gocb.Connect(seedNode2)
	myBu, _ := myCo.OpenBucket("testload", "")

	aQuery := gocb.NewN1qlQuery("select meta().id from testload where TYPE = \"DEVICE\"")
	rows, err := myBu.ExecuteN1qlQuery(aQuery, nil)

	// Interfaces for handling streaming return values
	var theID string
	var deviceID string
	var appID string
	var row map[string]interface{}

	for rows.Next(&row) {
		deviceID = row["id"].(string)

		s3 := rand.NewSource(time.Now().UnixNano())
		r3 := rand.New(s3)

		appID = "APP::" + strconv.Itoa(r3.Intn(5))
		theID = deviceID + "@@" + appID

		d := &theDOC{Device: deviceID, App: appID, Type: "REL"}
		myBu.Upsert(theID, d, 0)
	}
	if err = rows.Close(); err != nil {
		fmt.Printf("Couldn't get all the rows: %s\n", err);
	}

	fmt.Println("Let me query the relationships")

	bQuery := gocb.NewN1qlQuery("select meta().id,App,Device,TYPE from testload where TYPE = \"REL\"")
	rowsb, err := myBu.ExecuteN1qlQuery(bQuery, nil)
	for rowsb.Next(&row) {
		fmt.Printf("Row: %+v\n", row["id"].(string), row["App"].(string), row["Device"].(string), row["TYPE"].(string))
	}
	if err = rows.Close(); err != nil {
		fmt.Printf("Couldn't get all the rows: %s\n", err);
	}
}






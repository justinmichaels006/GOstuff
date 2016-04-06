package main

import (
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
	//"strconv"
	"sync"
	"os"
	"math/rand"
	"os/exec"
	"encoding/json"
	"bufio"
)

func main() {

	myCluster, _ := gocb.Connect("couchbase://10.111.94.102")
	myBucket, _ := myCluster.OpenBucket("testload", "")

	// read whole the file
	tmpACCT, err := os.OpenFile("/Users/justin/Documents/SVB/ACCT.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docACCT map[string]interface{}
	err = json.NewDecoder(tmpACCT).Decode(&docACCT)
	if err != nil {
		fmt.Println(err)
		os.Exit(52)
	}

	tmpCUST, err := os.OpenFile("/Users/justin/Documents/SVB/CUST.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docCUST map[string]interface{}
	err = json.NewDecoder(tmpCUST).Decode(&docCUST)
	if err != nil {
		fmt.Println(err)
		os.Exit(53)
	}

	tmpTRANS, err := os.OpenFile("/Users/justin/Documents/SVB/TRANS.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docTRANS map[string]interface{}
	err = json.NewDecoder(tmpTRANS).Decode(&docTRANS)
	if err != nil {
		fmt.Println(err)
		os.Exit(54)
	}

	//messages := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	// Create customers and accounts to use
	fmt.Println("Account and Customer creation ...")
	go SVBLoad(docACCT, docCUST, myBucket, )

	wg.Wait()

	/*wg.Add(1)
	// Simulate customers using their accounts
	fmt.Println("Transaction creation")
	go SVBtransaction(docTRANS, myBucket, wg)
	wg.Wait()*/

	fmt.Println("Good Bye")
	os.Exit(0)
}

func checker(e error) {
	if e != nil {
		panic(e)
	}
}

func SVBLoad(jsonACCT map[string]interface{}, jsonCUST map[string]interface{}, myBucket *gocb.Bucket) {

	//For use when transactions are generated
	f, err := os.OpenFile("/Users/justin/Documents/SVB/SVBaccounts.json", os.O_APPEND|os.O_WRONLY, 0600)
	checker(err)

	for i := 0; i < 9; i++ {

		//Element CUST_id used to create elements ACCT_Account1,2,3
		tmpA, err := exec.Command("uuidgen").Output()
		checker(err)

		CUSTID_ACCTAccount := string(tmpA)
		f.WriteString(CUSTID_ACCTAccount)
		jsonCUST["CUST_id"] = CUSTID_ACCTAccount

		//tmpA1, err := exec.Command("uuidgen").Output()
		//A1 := string(tmpA1)
		//f.WriteString(A1)
		jsonACCT["ACCT_Account1"] = CUSTID_ACCTAccount + ("::CHECKING")
		jsonACCT["ACCT_Account2"] = CUSTID_ACCTAccount + ("::CREDITCARD")
		jsonACCT["ACCT_Account3"] = CUSTID_ACCTAccount + ("::LOC")

		//i := strconv.Itoa(i)
		custID := ("CUST::" + CUSTID_ACCTAccount);
		acctID := ("ACCT::" + CUSTID_ACCTAccount);

		myBucket.Upsert(custID, jsonCUST, 0);
		myBucket.Upsert(acctID, jsonACCT, 0);

		//fmt.Println("Debug ... i" + i)
	}

	fmt.Println("Wait Group Done ...")
	defer f.Close()
	//defer sync.WaitGroup.Done()
	return
}

/*func SVBtransaction(jsonTRANS map[string]interface{}, myBucket *gocb.Bucket, wg sync.WaitGroup) {

	f, err := os.OpenFile("/Users/justin/Documents/SVB/SVBaccounts.json", os.O_RDONLY, 0600)
	checker(err)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		AccountArray = readLines(f)
	}

	for y := 0; y < 10; y++ {

		tmpT, err := exec.Command("uuidgen").Output()
		checker(err)
		TRANS := string(tmpT)
		transID := ("TRANS::" + TRANS)

		*//**//*x := RandAcct(1)*//**//*

		myBucket.Get("ACCT::" + TRANS, theAccount)
		// Set the account used during the customer transaction

		jsonTRANS["TRANS_element1"] = theAccount["ACCT_Account" + x]
		jsonTRANS["TRANS_Amount"] = RandIntBytes(5)*//**//*

		myBucket.Upsert(transID, jsonTRANS, 0)
		//fmt.Println("Debug ... y" + Aelement3_CHelement11)
	}

	defer f.Close()
	fmt.Println("Wait Group Done ...")
	defer wg.Done()
}*/

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func RandIntBytes(n int) string {
	const numBytes = "0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = numBytes[rand.Intn(len(numBytes))]
	}
	return string(b)
}

func RandAcct(n int) string {
	const numBytes = "23"

	b := make([]byte, n)
	for i := range b {
		b[i] = numBytes[rand.Intn(len(numBytes))]
	}
	return string(b)
}

package main

import (
	"github.com/couchbaselabs/gocb"
	"fmt"
	"os"
	"bufio"
)

func checker(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

myC, _ := gocb.Connect("couchbase://10.111.94.102")
myB, _ := myC.OpenBucket("testload", "")

	//Account Reader
	f, err := os.OpenFile("/Users/justin/Documents/SVB/SVBaccounts.json", os.O_RDONLY, 0600)
	checker(err)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		//fmt.Println(scanner.Text()) // Println will add back the final '\n'
		var myDoc interface{}
		custID := scanner.Text()
		myB.Get(custID, &myDoc)
		fmt.Println(custID)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}



	fmt.Println("Good Bye")
	os.Exit(0)
}




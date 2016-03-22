package main

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"github.com/couchbaselabs/gocb"
	"crypto/rand"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func main() {

	myCluster, _ := gocb.Connect("couchbase://192.168.61.101")
	myBucket, _ := myCluster.OpenBucket("testload", "")

	// read whole the file
	tmpPARK, err := os.Open("/Users/justin/Documents/Demo/samples/parks.json")
	check(err)

	r := bufio.NewReaderSize(tmpPARK, 4*1024)
	line, isPrefix, err := r.ReadLine()

	for err == nil && !isPrefix {
		s := string(line)

		theID, err := newUUID()
		check(err)

		fmt.Println(s)
		fmt.Println(theID)

		myBucket.Upsert(theID, s, 0);
		line, isPrefix, err = r.ReadLine()
	}

	if isPrefix {
		fmt.Println("buffer size to small")
		return
	}

	if err != io.EOF {
		fmt.Println(err)
		return
	}

}
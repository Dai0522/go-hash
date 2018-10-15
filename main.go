package main

import (
	"fmt"
	"io/ioutil"

	"github.com/Dai0522/go-hash/hash/bloomfilter"
)

func main() {
	b, err := ioutil.ReadFile("/Users/daiwei/Desktop/bf.txt")
	if err != nil {
		fmt.Print(err)
	}
	bf, _ := bloomfilter.Load(&b)

	// bf, _ := bloomfilter.New(100000, 0.0001)

	// for i := 0; i < 100000; i++ {
	// 	bf.PutUint64(uint64(i))
	// }

	total := 0
	fail := 0
	for i := 100001; i < 1100001; i++ {
		if bf.MightContainUint64(uint64(i)) {
			fmt.Println(i)
			fail++
		}
		total++
	}
	fmt.Println(total, fail)
}

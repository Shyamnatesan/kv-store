package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/Shyamnatesan/kv-store/btree"
)

// const M = 4

// type Data struct {
// 	keys   []int32
// 	values []string
// }

// func NewData() *Data {
// 	return &Data{
// 		keys:   make([]int32, M),
// 		values: make([]string, M),
// 	}
// }

// func NewNode() *Node {
// 	return &Node{
// 		numKeys:  0,
// 		data:     NewData(),
// 		children: [M + 1]*Node{},
// 		parent:   nil,
// 		isLeaf:   true,
// 	}
// }

// type Node struct {
// 	numKeys  int
// 	data     *Data
// 	children [M + 1]*Node
// 	parent   *Node
// 	isLeaf   bool
// }

// type Tree struct {
// 	root    *Node
// 	maxKeys int
// }

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 							   			|200:monish|
// 							  			 /        \
// 									   /            \
//					 |100:shyam, 150:johny|     |300:natesan, 400:suganthi|


func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}


func main() {
		
	btree := btree.NewTree()
	for i := 100000; i >= 0; i-- {
	    btree.Put(int32(i), RandString(5))
	}

	result := btree.Print()
	fmt.Println(result[200: 300])
	



	file, err := os.OpenFile("example.bin", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	result = btree.SerializeAndDeserialize(file)
	fmt.Println(result[200: 300])

}
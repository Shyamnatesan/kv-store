package main

import "github.com/Shyamnatesan/kv-store/btree"

func main() {
	btree := btree.NewTree()
	btree.Put(10)

	btree.Print()
}
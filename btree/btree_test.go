package btree

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
)


func generateUniqueRandomNumbers(n, min, max int) []int {
	if max - min < n {
		fmt.Println("Error: Cannot generate unique numbers with the given range and count.")
		return nil
	}
	nums := make(map[int]bool)
	result := make([]int, 0, n)
	for len(nums) < n {
		num := rand.Intn(max-min) + min
		if !nums[num] {
			nums[num] = true
			result = append(result, num)
		}
	}
	return result
}


func GenerateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]rune, length)

	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func TestBtreeInsertAndSearch(t *testing.T) {
	n := 100
	btree := NewTree()
	// generating random intergers(keys) from 0 to n for testing
	keys := generateUniqueRandomNumbers(n, 0, n)

	// generating random strings(values) from 0 to n for testing
	values := []string{}
	for i := 0; i < n; i++ {
		values = append(values, GenerateRandomString(5))
	}

	// inserting the key-value into the btree, in random order
	for i := 0; i < n; i++ {
		t.Logf("putting key: %d, value: %s", keys[i], values[keys[i]])
		btree.Put(keys[i], values[keys[i]])
	}

	// searching for the inserted keys test
	for i := 0; i < n; i++ {
		node, pos := btree.Find(keys[i])
		if node == nil || pos == -1 || node.data.keys[pos] != keys[i] || node.data.values[pos] != values[keys[i]] {
			t.Errorf("Find(%d) returned unexpected result", keys[i])
		}
	}

	// getting the inorder traversal of the btree
	result := btree.Print()
	for i := 0; i < n; i++ {
		if i != result[i].key {
			t.Errorf("keys are not in a sorted manner")
		}
	}

	log.Println("Insert and Search test passed")
}

func TestBtreeDeletionTest(t *testing.T) {
	n := 1000000
	btree := NewTree()

	// Generate unique random keys
	keys := generateUniqueRandomNumbers(n, 0, n)
	values := make([]string, n)

	// Generate random values
	for i := 0; i < n; i++ {
		values[i] = GenerateRandomString(5)
	}

	// Insert keys and values into the B-tree
	for i := 0; i < n; i++ {
		btree.Put(keys[i], values[keys[i]])
	}

	// Delete half of the keys from the B-tree
	for i := 0; i < n/2; i++ {
		btree.Del(i)
	}

	// Check if the remaining keys still exist in the B-tree
	for i := n / 2; i < n; i++ {
		node, pos := btree.Find(i)
		if node == nil || pos == -1 {
			t.Errorf("Find(%d) returned unexpected result after deletion", keys[i])
		}
	}

	log.Println("deletion test passed")
}
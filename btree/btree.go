package btree


const (
	M = 4 // ORDER OF A TREE
	MAX_NUM_OF_KEYS = M - 1 // MAXIMUM NUMBER OF KEYS ALLOWED IN A NODE
	ceilOfM = M / 2
	MIN_NUM_OF_KEYS = ceilOfM - 1 // MINIMUM NUMBER OF KEYS ALLOWED IN A NODE EXCEPT ROOT
)


type Data struct {
	keys []int
	values []string
}

func NewData() *Data {
	return &Data{
		keys: make([]int, M),
		values: make([]string, M),
	}
}

type Node struct {
	numKeys int
	data *Data
	children [M + 1]*Node
	parent *Node
	isLeaf bool
}

func (node *Node) searchNode(key int, pos *int) (*Node, int) {
	// if found return the node
	// else return nil
	// log.Println("searching for key ", key, " in the node ", node.data.keys)
	for ; *pos < node.numKeys; {
		if key > node.data.keys[*pos] {
			*pos++;
		}else if key ==  node.data.keys[*pos]{
			// log.Println("found a match at position ", *pos, "in node ", node.data.keys)
			return node, *pos;
		}else{
			break
		}
	}
	return nil, -1
}

func (node *Node) search(key int) (*Node, int) {
	if node == nil {
		return nil, -1
	}
	pos := 0

	found, indexOfKey := node.searchNode(key, &pos)
	if found != nil {
		return found, indexOfKey
	}

	return node.children[pos].search(key)
}

func (node *Node) insertIntoNode(key int, value string) int {
	// find the position "pos" => (which index) where to insert
	// log.Println("inserting key", key, " in node", node.data.keys)
	pos := 0
	for ; pos < node.numKeys;  {
		if key > node.data.keys[pos] {
			pos++;
		}else if key == node.data.keys[pos] {
			// log.Println("key already exists. No duplicate allowed")
			return -1
		}else{
			break
		}
	}


	// shift the elements to the right based on pos
	for i := node.numKeys - 1; i >= pos; i-- {
		node.data.keys[i + 1] = node.data.keys[i]
		node.data.values[i + 1] = node.data.values[i]

	}
	// insert
	node.data.keys[pos] = key
	node.data.values[pos] = value
	node.numKeys++
	// log.Println("successfully inserted key = ", key, "at position ", pos)
	return pos
}

func (node *Node) maxKeyThresholdReached(key int, value string, rightChild *Node, tree *Tree, pos *int) {
	// this is a recursive function

	// rightChild is initially nil
	// when recursively called for parents, we pass the rightnode

	// if rightnode is not nil, we have to add it to the node.children
	// if rightnode is nil, we do nothing

	// initially, check if the node.numKeys != maxnumberofkeys, 
	// if so return from this function
	if node.numKeys != MAX_NUM_OF_KEYS {
		index := node.insertIntoNode(key, value)
		if rightChild != nil {
            for i := node.numKeys; i > index + 1; i-- {
                node.children[i] = node.children[i - 1]
            }
            node.children[index + 1] = rightChild
            rightChild.parent = node
        }
		return
	}else{
		// else, insert the key in the keys array in a sorted manner
		index := node.insertIntoNode(key, value)
		if rightChild != nil {
            for i := node.numKeys; i > index + 1; i-- {
                node.children[i] = node.children[i - 1]
            }
            node.children[index + 1] = rightChild
            rightChild.parent = node
        }
		// split it. now you'll have leftnode, rightnode, and median(to pass to the parent node)
		rightNode, medianKey, medianValue := node.splitNode()
		if node.parent == nil {
			// log.Println("no parent")
			// if this current node's parent is nil,
			// the currentnode does not have a parent, so we create a new node
			parentNode := NewNode()
			// set the isLeaf to false, since parentNode has leftnode and rightnode as children
			parentNode.isLeaf = false
			// then create a new node and update its children(add leftnode and righnode as children)
			parentNode.insertIntoNode(medianKey, medianValue)
			// add the leftnode(node) and rightnode as children to this parent node
			parentNode.children[0] = node
			parentNode.children[1] = rightNode
			// also update the parent field in the leftnode(node) and the rightnode
			node.parent = parentNode
			rightNode.parent = parentNode

			// since we create a new node and add a key,
			// we have to make this parentNode as root, and then return.
			// or somehow make this parent node a root
			tree.root = parentNode
			return;
		}else{
			// else, if current node's parent is not nil,
			// then after splitting, we'll have a median right,
			// call this function recursively for the current node's parent and pass the median as key
			node.parent.maxKeyThresholdReached(medianKey, medianValue, rightNode, tree, pos)

		}
	}
}

func (node *Node) splitNode() (*Node, int, string) {
	// log.Println("splitting the node ", node.keys)
	medianIndex := node.numKeys / 2
	if medianIndex % 2 == 0 {
		medianIndex--
	}
	mediankey := node.data.keys[medianIndex]
	medianValue := node.data.values[medianIndex]

	// create a new node called rightnode
	rightNode := NewNode()
	// rightnode.isLeaf = node.isLeaf (because only if node has children,
	//  they will be copied to the rightnode)
	rightNode.isLeaf = node.isLeaf

	// remaining keys after median, add it to the rightnode
	j := 0
	for i := medianIndex + 1; i < node.numKeys; i++ {
		rightNode.data.keys[j] = node.data.keys[i]
		rightNode.data.values[j] = node.data.values[i]
 		rightNode.numKeys++
		j++
	}

	// keep all the keys in node.keys till the median,
	for i := medianIndex; i < node.numKeys; i++ {
		node.data.keys[i] = 0 // Assuming keys are int type; otherwise, use zero value of the key type
		node.data.values[i] = ""
	}
	node.numKeys = medianIndex


	// remaining children after median + 1 send it to the rightnode
	if !node.isLeaf {
		// copy(rightNode.children[:], node.children[medianIndex+1:node.numKeys+1])
		// log.Println("moving children to rightnode", rightNode)
		j := 0
		for i := medianIndex + 1; i < M + 1; i++ {
			// log.Println(node.children[i].keys)
			rightNode.children[j] = node.children[i]
			node.children[i] = nil
			rightNode.children[j].parent = rightNode
			j++
		}
	}
	// if children in node, then keep all children in node.children till, median(index)
	// copy(node.children[:], node.children[:medianIndex + 1])
	
	// else no children, then do nothing
	return rightNode, mediankey, medianValue
}

func (node *Node) insert(key int, value string, tree *Tree) {
	// if the node has space, 
	// insert into the node by shifting elements accordingly
	pos := 0
	if node.isLeaf {
		if node.numKeys == MAX_NUM_OF_KEYS {
			// log.Println("max limit reached")
			// this is a leaf node with maximum number of keys
			// we have to insert and split and check them recursively for maxNumberOfKeys
			// so, every node should have access to its parent node
			node.maxKeyThresholdReached(key, value, nil, tree, &pos)
		}else{
			// this is a leaf node with space to add another key, so we just insert it
			node.insertIntoNode(key, value);
			return
		}
	}else{
		node.searchNode(key, &pos)
		// log.Println("going to child node at index", pos)
		node.children[pos].insert(key, value, tree)
	}
}

type Pair struct {
	key int
	value string
}

func NewPair(key int, value string) *Pair{
	return &Pair{
		key: key,
		value: value,
	}
}


func (node *Node) inorder(result *[]Pair) {
	if node == nil {
	  return
	}
	if node.isLeaf {
	  // Print the keys in the leaf node
	  for i := 0; i < node.numKeys; i++ {
		pair := NewPair(node.data.keys[i], node.data.values[i])
		*result = append(*result, *pair)
	  }
	  return
	}
	// Traverse child nodes and keys
	for i := 0; i <= node.numKeys; i++ { // process all children (0 to numKeys+1)
	  // Recursively traverse the i-th child
	  node.children[i].inorder(result)
	  // Append the i-th key (if applicable)
	  if i < node.numKeys {
		pair := NewPair(node.data.keys[i], node.data.values[i])
		*result = append(*result, *pair)
	  }
	}
}




func NewNode() *Node {
	return &Node{
		numKeys: 0,
		data: NewData(),
		children: [M + 1]*Node{},
		parent: nil,
		isLeaf: true,
	}
}

type Tree struct {
	root *Node
	maxKeys int
}



func (tree *Tree) Find(key int) (*Node, int) {
	currentNode := tree.root
	result, posOfKey := currentNode.search(key)
	return result, posOfKey
}

func (tree *Tree) Put(key int, value string) {
	if tree.root == nil {
		tree.root = NewNode()
	}
	currentNode := tree.root
	// pos := 0
	currentNode.insert(key, value, tree)
}

// [1, 2, 3, 4, 5]
func (node *Node) deleteKeyInNode(pos int) {
	// log.Println("deleting key in position ", pos, " in node", node.data.keys)
	n := node.numKeys
	i := pos + 1
	for ; i < n; i++ {
		node.data.keys[i - 1] = node.data.keys[i]
		node.data.values[i - 1] = node.data.values[i]
	}
	node.data.keys[i - 1] = 0
	node.data.values[i - 1] = ""
	node.numKeys--
}

func (node *Node) getSiblings() (*Node, *Node, int, int) {
	// log.Println("getting siblings of node ", node.data.keys)
	parentNode := node.parent
	var leftSibling *Node
	var rightSibling *Node
	var leftSeparaterIndex int
	var rightSeparaterIndex int
	var nodeIndexInParent int
	// Find the index of the node in the parent's children
	for i := 0; i < parentNode.numKeys + 1; i++ {
		if parentNode.children[i] == node {
			nodeIndexInParent = i
			break
		}
	}
	// If the node is not the first child, then left sibling exists
	if nodeIndexInParent > 0 {
		// left sibling exists
		leftSibling = parentNode.children[nodeIndexInParent - 1]
		leftSeparaterIndex = nodeIndexInParent - 1
		// log.Println("leftSibling of ", node.data.keys, " is ", leftSibling.data.keys)
	}
	// If the node is not the last child, then right sibling exists
	if nodeIndexInParent < parentNode.numKeys {
		// right sibling exists
		rightSibling = parentNode.children[nodeIndexInParent + 1]
		rightSeparaterIndex = nodeIndexInParent
		// log.Println("rightSibling of ", node.data.keys, " is ", rightSibling.data.keys)
	}
	return leftSibling, rightSibling, leftSeparaterIndex, rightSeparaterIndex
}

func (node *Node) borrowFromLeftSibling(leftSibling *Node, separaterIndex int) {
	parentNode := node.parent
	rightMostkeyInLeftSiblingKey := leftSibling.data.keys[leftSibling.numKeys - 1]
	rightMostkeyInLeftSiblingValue := leftSibling.data.values[leftSibling.numKeys - 1]
	separaterKey := parentNode.data.keys[separaterIndex]
	separaterValue := parentNode.data.values[separaterIndex]
	node.insertIntoNode(separaterKey, separaterValue)
	if !node.isLeaf {
		// if internal node,
		// handle shifting the *ptr(rightChild of rightMostKeyInLeftSibling) to
		// the leftChild of the separater key in the node
		rightMostChildInLeftSibling := leftSibling.children[leftSibling.numKeys]
		// to insert this child to the node's 0th child,
		// first you have to right shift them by 1
		// log.Println("right shifting the children of node.children by 1", node.data.keys)
		i := node.numKeys
		for ; i > 0; i-- {
			node.children[i] = node.children[i - 1] // right shift by 1
		}
		// after right shifting elements, now put the key as the left most child in the node
		node.children[i] = rightMostChildInLeftSibling
		// updating the parent pointer
		rightMostChildInLeftSibling.parent = node
		// now delete the rightMostChildInLeftSibling in leftSibling
		leftSibling.children[leftSibling.numKeys] = nil
	}
	parentNode.data.keys[separaterIndex] = rightMostkeyInLeftSiblingKey
	parentNode.data.values[separaterIndex] = rightMostkeyInLeftSiblingValue
	leftSibling.deleteKeyInNode(leftSibling.numKeys - 1)
	// log.Println("borrowing from left sibling successfull")
}

func (node *Node) borrowFromRightSibling(rightSibling *Node, separaterIndex int) {
	parentNode := node.parent
	leftMostKeyInRightSiblingKey := rightSibling.data.keys[0]
	leftMostKeyInRightSiblingValue := rightSibling.data.values[0]
	separaterKey := parentNode.data.keys[separaterIndex]
	separaterValue := parentNode.data.values[separaterIndex]
	node.insertIntoNode(separaterKey, separaterValue)
	if !node.isLeaf {
		// if internal node, 
		// handle shifting the *ptr(leftChild of leftMostKeyInRightSibling) to
		// the rightChild of the separater key in the node
		leftMostChildInRightSibling := rightSibling.children[0]
		// adding the child at the end. no need to shift
		node.children[node.numKeys] = leftMostChildInRightSibling
		// updating the parent pointer
		leftMostChildInRightSibling.parent = node
		// since we copied the leftMostChildInRightSibling, we should
		// handle shifting the children by 1 in the rightSibling
		i := 1
		for ; i <= rightSibling.numKeys; i++ {
			rightSibling.children[i - 1] = rightSibling.children[i]
		}
		// when shifting above, the last key will be shifted (but also duplicated, so we delete that also)
		rightSibling.children[i - 1] = nil
	}
	parentNode.data.keys[separaterIndex] = leftMostKeyInRightSiblingKey
	parentNode.data.values[separaterIndex] = leftMostKeyInRightSiblingValue
	rightSibling.deleteKeyInNode(0)
	// log.Println("borrowing from right sibling successfull")
}

func (node *Node) mergeNodes(separaterIndex int, node2 *Node, tree *Tree)  *Node {

	initialNumOfNodes := node.numKeys
	parentNode := node.parent
	separaterKey := parentNode.data.keys[separaterIndex]
	separaterValue := parentNode.data.values[separaterIndex]
	// inserting separater and the keys in the underflow nodes to the node
	// log.Println("merging nodes", node.data.keys, separaterKey, node2.data.keys)
	node.insertIntoNode(separaterKey, separaterValue)
	for i := 0; i < node2.numKeys; i++ {
		node.insertIntoNode(node2.data.keys[i], node.data.values[i])
	}

	if !node.isLeaf {
		// if it is not a leaf node,
		// then this is an internal node
		// so, we must insert/copy the children as well from node2 to the node
		// log.Println(node.data.keys, " is not a leaf node, so shifting children")
		j := initialNumOfNodes + 1
		for i := 0; i <= node2.numKeys; i++ {
			// updating the parent node
			node2.children[i].parent = node
			node.children[j] = node2.children[i]
			j++
		}
	}
	// handle shifting children to the left by one in the parent node
	i := separaterIndex + 2
	for ; i <= parentNode.numKeys; i++ {
		parentNode.children[i - 1] = parentNode.children[i]
	}
	// when shifting above, the last key will be shifted (but also duplicated, so we delete that also)
	parentNode.children[i - 1] = nil 
	// we delete the separater key now. here, the deleteKeyInNode handles the shifting of keys as well
	parentNode.deleteKeyInNode(separaterIndex)
	if parentNode.numKeys == 0 && parentNode == tree.root {
		tree.root = node
		return node
	}
	// log.Println("merging successfull")
	return parentNode	
}

func (node *Node) rebalancing(tree *Tree) {
	if node.numKeys >= MIN_NUM_OF_KEYS {
		return
	}
	// log.Println("reblancing node ", node.data.keys)
	leftSibling, rightSibling, leftSeparaterIndex, rightSeparaterIndex := node.getSiblings()
	if leftSibling != nil && leftSibling.numKeys > MIN_NUM_OF_KEYS {
		// leftSibling exists and has keys to spare
		node.borrowFromLeftSibling(leftSibling, leftSeparaterIndex)
		return
	}else if rightSibling != nil && rightSibling.numKeys > MIN_NUM_OF_KEYS {
		// rightSibling exists and has keys to spare
		node.borrowFromRightSibling(rightSibling, rightSeparaterIndex)
		return
	}else if leftSibling != nil && leftSibling.numKeys <= MIN_NUM_OF_KEYS {
		// leftSibling exists and no keys to spare, so we merge with leftSibling
		parentNode := leftSibling.mergeNodes(leftSeparaterIndex, node, tree)
		parentNode.rebalancing(tree)
	}else if rightSibling != nil && rightSibling.numKeys <= MIN_NUM_OF_KEYS {
		// rightSibling exists and no keys to spare, so we merge with rightSibling
		parentNode := node.mergeNodes(rightSeparaterIndex, rightSibling, tree)
		parentNode.rebalancing(tree)
	}
}

func (node *Node) deleteKeyFromLeafNode(posOfKey int, tree *Tree) {
	// 1. delete the key
	node.deleteKeyInNode(posOfKey)
	if node.numKeys >= MIN_NUM_OF_KEYS {
		// no need to rebalance, since the noOfKeys is >= minimum number of keys in a node threshold
		// log.Println("no need to rebalance, threshold is available")
		return
	}
	// log.Println("node underflow, rebalancing", node.data.keys)
	node.rebalancing(tree)
	// log.Println("deletion and rebalancing successfull")
}

func (node *Node) copyPredecessor(key int) (*Node, int, string) {
	if node.isLeaf {
		predecessorKey := node.data.keys[node.numKeys - 1]
		predecessorValue := node.data.values[node.numKeys - 1]
		node.data.keys[node.numKeys - 1] = key
		return node, predecessorKey, predecessorValue
	}
	return node.children[node.numKeys].copyPredecessor(key)
}

func (node *Node) copySuccessor(key int) (*Node, int, string) {
	if node.isLeaf {
		successorKey := node.data.keys[0]
		successorValue := node.data.values[0]
		node.data.keys[0] = key
		return node, successorKey, successorValue
	}
	return node.children[0].copySuccessor(key)
}


func (tree *Tree) Del(key int) {
	node, posOfKey := tree.root.search(key)
	if node == nil {
		// log.Println("key does not exist in the btree")
		return
	}
	if node.isLeaf {
		// key is in a leaf node
		node.deleteKeyFromLeafNode(posOfKey, tree)
		return
		
	}
	// find the leftChild and rightChild
	leftChild := node.children[posOfKey]
	rightChild := node.children[posOfKey + 1]
	if leftChild != nil {
		// find the inorder predecessor
		// we get the max key from leftSibling to swap the key to be deleted in the node
		leafNode, predecessorKey, predecessorValue := leftChild.copyPredecessor(key)
		node.data.keys[posOfKey] = predecessorKey
		node.data.values[posOfKey] = predecessorValue
		// after replacing/swapping, we call leftChild.deleteKeyFromLeafNode(index of the maxKey)
		leafNode.deleteKeyFromLeafNode(leafNode.numKeys - 1, tree)
		return
	}
	// otherwise if leftChild is nil or does not have keys to spare, we go for rightChild
	// find the successor
	// we get the minKey from rightSibling to swap the key to be deleted in the node
	leafNode, successorKey, successorValue := rightChild.copySuccessor(key)
	node.data.keys[posOfKey] = successorKey
	node.data.values[posOfKey] = successorValue
	// after replacing/swapping, we call rightSibling.deleteKeyFromLeafNode(index of the minKey)
	leafNode.deleteKeyFromLeafNode(0, tree)
	// log.Println()
}

func (tree *Tree) Print() []Pair {
    if tree.root == nil {
        return []Pair{}
    }

    result := []Pair{}
    tree.root.inorder(&result)
    return result
}

func NewTree() *Tree {
	return &Tree{
		root: nil,
		maxKeys: MAX_NUM_OF_KEYS,
	}
}
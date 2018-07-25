package BLC

import (
	"crypto/sha256"
)


type Jyh_MerkleTree struct {
	RootNode *Jyh_MerkleNode
}


// Block  [tx1 tx2 tx3 tx3]


//MerkleNode{nil,nil,tx1Bytes}
//MerkleNode{nil,nil,tx2Bytes}
//MerkleNode{nil,nil,tx3Bytes}
//MerkleNode{nil,nil,tx3Bytes}
//
//

//
//MerkleNode:
//	left: MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
//
//	right: MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
//
//	sha256(sha256(tx1Bytes,tx2Bytes)+sha256(tx3Bytes,tx3Bytes))





type Jyh_MerkleNode struct {
	Left  *Jyh_MerkleNode
	Right *Jyh_MerkleNode
	Data  []byte
}


func NewMerkleTree(data [][]byte) *Jyh_MerkleTree {

	//[tx1,tx2,tx3]

	var nodes []Jyh_MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
		//[tx1,tx2,tx3,tx3]
	}

	// 创建叶子节点
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}


	//MerkleNode{nil,nil,tx1Bytes}
	//MerkleNode{nil,nil,tx2Bytes}
	//MerkleNode{nil,nil,tx3Bytes}
	//MerkleNode{nil,nil,tx3Bytes}



	// 　循环两次
	for i := 0; i < len(data)/2; i++ {

		var newLevel []Jyh_MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		//MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
		//
		//MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
		//

		if len(newLevel) % 2 != 0 {

			newLevel = append(newLevel,newLevel[len(newLevel) - 1])
		}
		nodes = newLevel
	}

	//MerkleNode:
	//	left: MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
	//
	//	right: MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
	//
	//	sha256(sha256(tx1Bytes,tx2Bytes)+sha256(tx3Bytes,tx3Bytes))

	mTree := Jyh_MerkleTree{&nodes[0]}

	return &mTree
}


func NewMerkleNode(left, right *Jyh_MerkleNode, data []byte) *Jyh_MerkleNode {
	mNode := Jyh_MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right

	return &mNode
}
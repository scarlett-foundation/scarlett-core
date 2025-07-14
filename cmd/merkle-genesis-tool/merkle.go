package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"scarlett-core/x/proofofdegen/types"
)

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Hash  string
	Left  *MerkleNode
	Right *MerkleNode
	Data  string // For leaf nodes, contains address:weight
}

// GenerateMerkleTree creates a Merkle tree from wallet entries
func GenerateMerkleTree(wallets []WalletEntry) (string, map[string]types.MerkleProof, uint64, error) {
	if len(wallets) == 0 {
		return "", nil, 0, fmt.Errorf("no wallets provided")
	}

	// Sort wallets by address for deterministic tree
	sort.Slice(wallets, func(i, j int) bool {
		return wallets[i].Address < wallets[j].Address
	})

	// Calculate total weight
	var totalWeight uint64
	for _, wallet := range wallets {
		totalWeight += wallet.Weight
	}

	// Create leaf nodes
	var leafNodes []*MerkleNode
	for _, wallet := range wallets {
		leafData := fmt.Sprintf("%s:%d", wallet.Address, wallet.Weight)
		hash := hashData(leafData)
		leafNodes = append(leafNodes, &MerkleNode{
			Hash: hash,
			Data: leafData,
		})
	}

	// Build tree from bottom up
	root := buildTree(leafNodes)
	if root == nil {
		return "", nil, 0, fmt.Errorf("failed to build Merkle tree")
	}

	// Generate proofs for each wallet
	proofs := make(map[string]types.MerkleProof)
	for _, wallet := range wallets {
		proof := generateProof(root, wallet.Address, wallet.Weight)
		proofs[wallet.Address] = types.MerkleProof{
			Proof:  proof,
			Weight: wallet.Weight,
		}
	}

	return root.Hash, proofs, totalWeight, nil
}

// hashData creates a SHA-256 hash of the input data
func hashData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// buildTree constructs the Merkle tree from leaf nodes
func buildTree(nodes []*MerkleNode) *MerkleNode {
	if len(nodes) == 0 {
		return nil
	}

	if len(nodes) == 1 {
		return nodes[0]
	}

	var nextLevel []*MerkleNode

	// Pair up nodes and create parent nodes
	for i := 0; i < len(nodes); i += 2 {
		left := nodes[i]
		var right *MerkleNode

		if i+1 < len(nodes) {
			right = nodes[i+1]
		} else {
			// If odd number of nodes, duplicate the last one
			right = left
		}

		// Create parent node
		combinedHash := hashData(left.Hash + right.Hash)
		parent := &MerkleNode{
			Hash:  combinedHash,
			Left:  left,
			Right: right,
		}

		nextLevel = append(nextLevel, parent)
	}

	// Recursively build the next level
	return buildTree(nextLevel)
}

// generateProof creates a Merkle proof for a specific address and weight
func generateProof(root *MerkleNode, address string, weight uint64) []string {
	targetData := fmt.Sprintf("%s:%d", address, weight)
	targetHash := hashData(targetData)

	var proof []string
	findProof(root, targetHash, &proof)
	return proof
}

// findProof recursively finds the Merkle proof path
func findProof(node *MerkleNode, targetHash string, proof *[]string) bool {
	if node == nil {
		return false
	}

	// If this is a leaf node, check if it matches our target
	if node.Left == nil && node.Right == nil {
		return node.Hash == targetHash
	}

	// Check left subtree
	if findProof(node.Left, targetHash, proof) {
		// Target is in left subtree, add right sibling to proof
		if node.Right != nil {
			*proof = append(*proof, node.Right.Hash)
		}
		return true
	}

	// Check right subtree
	if findProof(node.Right, targetHash, proof) {
		// Target is in right subtree, add left sibling to proof
		if node.Left != nil {
			*proof = append(*proof, node.Left.Hash)
		}
		return true
	}

	return false
}

package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"scarlett-core/x/proofofdegen/types"
)

// HashLeaf creates a hash of address and weight for Merkle tree leaf
func (k Keeper) HashLeaf(address string, weight uint64) string {
	// Create leaf data: address + weight
	leafData := fmt.Sprintf("%s:%d", address, weight)

	// Hash the leaf data
	hash := sha256.Sum256([]byte(leafData))
	return hex.EncodeToString(hash[:])
}

// ValidateMerkleProof validates a Merkle proof against the stored root
func (k Keeper) ValidateMerkleProof(address string, proof types.MerkleProof, merkleRoot string) (uint64, bool) {
	if merkleRoot == "" {
		return 0, false
	}

	// Start with the leaf hash
	currentHash := k.HashLeaf(address, proof.Weight)

	// Apply each proof element
	for _, proofElement := range proof.Proof {
		// Combine current hash with proof element
		// Order matters: we need to determine whether current hash goes left or right
		// For simplicity, use lexicographic ordering
		var combined string
		if currentHash < proofElement {
			combined = currentHash + proofElement
		} else {
			combined = proofElement + currentHash
		}

		// Hash the combined value
		hash := sha256.Sum256([]byte(combined))
		currentHash = hex.EncodeToString(hash[:])
	}

	// Check if final hash matches the merkle root
	return proof.Weight, currentHash == merkleRoot
}

// GetMerkleProofForAddress retrieves the stored Merkle proof for an address
func (k Keeper) GetMerkleProofForAddress(ctx context.Context, address string) (types.MerkleProof, error) {
	proof, err := k.MerkleProofs.Get(ctx, address)
	if err != nil {
		return types.MerkleProof{}, err
	}
	return proof, nil
}

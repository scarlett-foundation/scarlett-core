package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"scarlett-core/x/proofofdegen/types"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run cmd/merkle-genesis-tool/main.go <csv-file>")
		fmt.Println("CSV format: address,weight")
		os.Exit(1)
	}

	csvFile := os.Args[1]

	// Read CSV file
	eligibleWallets, err := readCSVFile(csvFile)
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		os.Exit(1)
	}

	if len(eligibleWallets) == 0 {
		fmt.Println("No eligible wallets found in CSV file")
		os.Exit(1)
	}

	fmt.Printf("Read %d eligible wallets from CSV\n", len(eligibleWallets))

	// Generate Merkle tree
	merkleRoot, proofs, totalWeight, err := GenerateMerkleTree(eligibleWallets)
	if err != nil {
		fmt.Printf("Error generating Merkle tree: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated Merkle tree with root: %s\n", merkleRoot)
	fmt.Printf("Total weight: %d\n", totalWeight)

	// Generate genesis configuration
	genesis := generateGenesisState(merkleRoot, totalWeight, proofs)

	// Write genesis JSON
	genesisJSON, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling genesis: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile("genesis_proofofdegen.json", genesisJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing genesis file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Generated genesis_proofofdegen.json")
	fmt.Printf("✅ Merkle root: %s\n", merkleRoot)
	fmt.Printf("✅ Total addresses: %d\n", len(eligibleWallets))
	fmt.Printf("✅ Total weight: %d\n", totalWeight)
}

// WalletEntry represents a single wallet entry from CSV
type WalletEntry struct {
	Address string
	Weight  uint64
}

// readCSVFile reads the CSV file and returns wallet entries
func readCSVFile(filename string) ([]WalletEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var wallets []WalletEntry
	for i, record := range records {
		// Skip header row if it exists
		if i == 0 && (record[0] == "address" || record[0] == "Address") {
			continue
		}

		if len(record) != 2 {
			return nil, fmt.Errorf("invalid CSV format at line %d: expected 2 columns, got %d", i+1, len(record))
		}

		address := record[0]
		weightStr := record[1]

		weight, err := strconv.ParseUint(weightStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid weight at line %d: %v", i+1, err)
		}

		wallets = append(wallets, WalletEntry{
			Address: address,
			Weight:  weight,
		})
	}

	return wallets, nil
}

// generateGenesisState creates the genesis state with Merkle tree data
func generateGenesisState(merkleRoot string, totalWeight uint64, proofs map[string]types.MerkleProof) types.GenesisState {
	// Create campaign with Merkle root
	campaign := types.Campaign{
		Index:           "genesis",
		Name:            "Genesis",
		Active:          true,
		TotalAllocation: totalWeight, // TODO: This should be total_weight field
		Creator:         "",
		MerkleRoot:      merkleRoot,
		ClaimedCount:    0,
		ClaimedWeight:   0,
	}

	// Convert proofs map to slice for genesis
	var merkleProofsList []types.MerkleProof
	for _, proof := range proofs {
		// Note: This is a workaround since MerkleProof doesn't have address field
		// In production, we might want to modify the proto or use a different approach
		merkleProofsList = append(merkleProofsList, proof)
	}

	return types.GenesisState{
		Params:            types.DefaultParams(),
		CampaignMap:       []types.Campaign{campaign},
		EligibleWalletMap: []types.EligibleWallet{}, // Empty - using Merkle proofs instead
		// TODO: Add MerkleProofMap when genesis support is implemented
	}
}

#!/bin/bash

# Emergency Stop Testing Script for Dynamic Governance Emissions
# This script tests the Phase 1 emergency controls implementation

echo "🚨 TESTING EMERGENCY STOP FUNCTIONALITY 🚨"
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CHAIN_ID="scarlett"
ALICE_KEY="alice"
PROPOSAL_DEPOSIT="10000000sclt"

echo -e "${BLUE}📋 Test Plan:${NC}"
echo "1. Start chain with fresh state"
echo "2. Submit emergency stop proposal"
echo "3. Vote on proposal"
echo "4. Verify emissions are halted"
echo "5. Submit resume proposal"
echo "6. Vote on resume proposal"
echo "7. Verify emissions are restored"
echo ""

# Step 1: Start the chain
echo -e "${YELLOW}🔄 Step 1: Starting chain with fresh state...${NC}"
echo "Please run: ignite chain serve --reset-once --verbose"
echo "Wait for chain to start, then press ENTER to continue..."
read -p ""

# Step 2: Check current emissions parameters
echo -e "${YELLOW}🔍 Step 2: Checking current emission parameters...${NC}"
echo "scarlett-cored query emissions params --output json | jq ."
echo ""

# Step 3: Submit emergency stop proposal
echo -e "${RED}🛑 Step 3: Submitting EMERGENCY STOP proposal...${NC}"
echo "scarlett-cored tx gov submit-proposal proposals/tests/emergency-stop-proposal.json --from $ALICE_KEY --chain-id $CHAIN_ID --yes"
echo ""

# Step 4: List proposals to get proposal ID
echo -e "${BLUE}📋 Step 4: Checking proposal status...${NC}"
echo "scarlett-cored query gov proposals --output json | jq '.proposals[] | {id: .id, title: .title, status: .status}'"
echo ""

# Step 5: Vote on proposal (replace PROPOSAL_ID)
echo -e "${GREEN}🗳️  Step 5: Voting YES on emergency stop proposal...${NC}"
echo "Replace PROPOSAL_ID with actual ID from step 4:"
echo "scarlett-cored tx gov vote PROPOSAL_ID yes --from $ALICE_KEY --chain-id $CHAIN_ID --yes"
echo ""

# Step 6: Wait for proposal to pass
echo -e "${YELLOW}⏳ Step 6: Waiting for proposal to pass...${NC}"
echo "Monitor proposal status:"
echo "scarlett-cored query gov proposal PROPOSAL_ID"
echo ""

# Step 7: Verify emergency stop is active
echo -e "${RED}🔍 Step 7: Verifying emergency stop is ACTIVE...${NC}"
echo "scarlett-cored query emissions params --output json | jq '.emergency_stop, .emergency_reason'"
echo ""

# Step 8: Check chain logs for emergency stop messages
echo -e "${RED}📋 Step 8: Expected log messages during emergency stop:${NC}"
echo "Look for these messages in chain logs:"
echo "  🛑 EMERGENCY STOP ACTIVE - HALTING ALL EMISSIONS 🛑"
echo "  🔥🔥🔥 GOVERNANCE-CONTROLLED MINT FUNCTION CALLED 🔥🔥🔥"
echo ""

# Step 9: Submit resume proposal
echo -e "${GREEN}🔄 Step 9: Submitting RESUME proposal...${NC}"
echo "scarlett-cored tx gov submit-proposal proposals/tests/emergency-resume-proposal.json --from $ALICE_KEY --chain-id $CHAIN_ID --yes"
echo ""

# Step 10: Vote on resume proposal
echo -e "${GREEN}🗳️  Step 10: Voting YES on resume proposal...${NC}"
echo "Replace RESUME_PROPOSAL_ID with actual ID:"
echo "scarlett-cored tx gov vote RESUME_PROPOSAL_ID yes --from $ALICE_KEY --chain-id $CHAIN_ID --yes"
echo ""

# Step 11: Verify emissions are restored
echo -e "${GREEN}✅ Step 11: Verifying emissions are RESTORED...${NC}"
echo "scarlett-cored query emissions params --output json | jq '.emergency_stop, .emergency_reason'"
echo ""

echo -e "${BLUE}🎉 Emergency Stop Test Complete!${NC}"
echo "=============================================="
echo "Expected Results:"
echo "- Emergency stop proposal should halt all emissions"
echo "- Chain logs should show 🛑 emergency stop messages"
echo "- No tokens should be minted/distributed during emergency"
echo "- Resume proposal should restore normal emissions"
echo "- Chain should resume normal token distribution" 
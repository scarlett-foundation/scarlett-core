# Scarlett: The Decentralised AI Operating System

**Scarlett** is a sovereign, proof-of-stake blockchain purpose-built to serve as the foundational layer for a new generation of decentralised, transparent, and community-owned AI applications—a true operating system for the AI age.

By integrating a provably fair launch, a community-directed economy, and a multi-layered development environment, Scarlett delivers the essential infrastructure for decentralised AI: verifiable compute, programmable incentives, and democratic governance.

> **Read the full vision in the [Scarlett Lightpaper](Scarlett-Lightpaper.md).**

## The Four Pillars of Scarlett

Scarlett rests on four core design principles that together form a self-sustaining, decentralised AI ecosystem.

### Pillar 1: A Provably Fair & Decentralised Genesis
Scarlett's legitimacy begins at genesis. It implements a **Satoshi-style fair launch**: no venture capital, no pre-mine, no insiders. The founding stake is programmatically burned over time through on-chain transactions, ensuring credible neutrality and sustainable scarcity.

### Pillar 2: A Community-Directed Economy
Scarlett's monetary policy is not hard-coded—it's a living system governed on-chain. The community collectively determines how token emissions are allocated, turning inflation into fuel for innovation.

### Pillar 3: A Multi-Layered Application Environment
Scarlett supports two development tracks: high-performance, governance-gated **Native Modules (Go)** for core logic and permissionless **Smart Contracts (CosmWasm)** for rapid, sandboxed innovation.

### Pillar 4: Stake-Based Application Onboarding
Scarlett introduces a novel **Stake-to-Register** model. Before any module or contract is eligible for community funding, its developers must stake `sclt`, creating an economic bond and ensuring long-term alignment.

## Get Started

To get a local development node running:

```bash
ignite chain serve
```

This command installs dependencies, builds, initializes, and starts the blockchain. The first time you run it, a wallet and genesis configuration will be created for you.

If you encounter startup errors after modifying the code, you may need to reset the chain's state to ensure compatibility with the new binary:

```bash
ignite chain serve --reset-once
```

## Key Functionality: Dynamic Emissions Governance

The core of Scarlett's economic model is managed by the `x/emissions` module, which is controlled by on-chain governance. Here's how you can participate.

### 1. Query Current Emission Parameters

To see the active emission split, run:

```bash
scarlett-cored query emissions params -o json
```

### 2. Create a Governance Proposal

To change the parameters, you must submit a governance proposal. Create a JSON file (e.g., `proposal.json`) with the new parameters. The `authority` must be the governance module's account address.

**Example `proposal.json` (60% to Validators, 40% to AI Rewards):**

```json
{
  "title": "Activate Governance-Controlled Emissions: 60% Validators, 40% AI Inference",
  "description": "Enable dynamic emissions with 60% to validators (fee_collector) and 40% to AI inference rewards (inferencerewards)",
  "summary": "Activate governance-controlled token emissions with 60/40 split",
  "messages": [
    {
      "@type": "/scarlettcore.emissions.v1.MsgUpdateParams",
      "authority": "scarlett10d07y265gmmuvt4z0w9aw880jnsr700j4l5sjv",
      "params": {
        "emission_destinations": "[{\"module_name\":\"fee_collector\",\"weight\":\"0.60\",\"description\":\"Validator and delegator rewards\",\"enabled\":true},{\"module_name\":\"inferencerewards\",\"weight\":\"0.40\",\"description\":\"AI inference provider rewards\",\"enabled\":true}]",
        "enabled": true
      }
    }
  ],
  "deposit": "10000000sclt"
}
```

### 3. Submit the Proposal

Submit the proposal transaction from a funded account (e.g., `alice`):

```bash
scarlett-cored tx gov submit-proposal proposal.json --from=alice --yes --gas=auto --gas-adjustment=1.5
```

### 4. Vote on the Proposal

Once a proposal is in the voting period, validators and token holders can vote on it.

```bash
# Query proposals to find the proposal ID
scarlett-cored query gov proposals

# Vote 'yes' on proposal #1 from the 'alice' account
scarlett-cored tx gov vote 1 yes --from=alice --yes
```

If the proposal passes, the new emission parameters will take effect immediately.

## Useful Queries

-   **Check Wallet Balance**:
    `scarlett-cored q bank balances [address]`
-   **Check Claimable Staking Rewards**:
    `scarlett-cored q distribution rewards [delegator_address]`
-   **Check Community Pool Balance**:
    `scarlett-cored q distribution community-pool`

## Learn more

-   [Ignite CLI](https://ignite.com/cli)
-   [Tutorials](https://docs.ignite.com/guide)
-   [Ignite CLI docs](https://docs.ignite.com)
-   [Cosmos SDK docs](https://docs.cosmos.network)
-   [Developer Chat](https://discord.gg/ignite)

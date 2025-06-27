# Scarlett Core: A Purpose-Built L1 Blockchain for Decentralized AI

**Scarlett Core** is a sovereign, proof-of-stake blockchain built with the [Cosmos SDK](https://docs.cosmos.network) and [Tendermint](https://tendermint.com/). It is designed from the ground up to provide a robust, decentralized foundation for AI-driven applications, featuring a dynamic and community-controlled economic model.

## Core Features

-   **Dynamic Governance-Controlled Emissions**: Scarlett Core's primary innovation is its flexible token emission system. Instead of a hardcoded inflation model, the Scarlett community can dynamically control how new tokens are distributed through on-chain governance. This allows the network to adapt and incentivize different behaviors as the AI ecosystem evolves. Key destinations include staking rewards (`fee_collector`) and AI inference provider rewards (`inferencerewards`).

-   **Modular & Extensible**: Built on the Cosmos SDK, Scarlett Core is fully modular. This allows for the rapid development and integration of new features and custom modules, ensuring the chain can keep pace with the rapid evolution of AI.

-   **High Performance & Fast Finality**: Leveraging the Tendermint consensus engine, Scarlett Core provides the high throughput and fast transaction finality required for real-time AI applications and services.

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

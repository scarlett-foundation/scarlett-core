# Dynamic Emissions Governance - Acceptance Criteria

## 🎯 Project Goal

The primary goal is to **replace the rigid, hardcoded token emission system with a flexible, transparent, and democratic system controlled entirely by on-chain governance.**

This transforms the static 50/50 split into a system that empowers the Scarlett Core community to make critical economic decisions by allowing them to:

- **Control Distribution:** Dynamically change the percentage of newly minted tokens allocated to different parts of the ecosystem (e.g., validator rewards, AI inference rewards, a future grants program, etc.).
- **Add/Remove Destinations:** Introduce new reward destinations or decommission old ones through governance proposals.
- **Ensure Transparency:** Make every change to the emission policy a public, auditable event on the blockchain.

Ultimately, we are building a core economic engine that allows the protocol to incentivize the most valuable activities as determined by its stakeholders.

---

## ✅ Acceptance Criteria

### 1. Default Behavior (Fresh Genesis)
- [ ] When the chain starts with no governance parameters set, emissions **must** default to the hardcoded **50/50 split** between the `fee_collector` (validators) and `inferencerewards` modules.
- [ ] The standard `distribution` module's community pool **must receive zero** tokens from this custom minting process. We must verify its balance does not increase from emissions.

### 2. Governance Control
- [ ] It must be possible to submit a governance proposal to change the emission splits (e.g., to 80/20).
- [ ] After the proposal passes, the new parameters **must** take effect immediately.
- [ ] The system must correctly distribute tokens according to the new, governance-defined percentages.
- [ ] Only the governance module is authorized to make these changes.

### 3. Verification & Testing
- [ ] We must be able to query the module balances (`fee_collector`, `inferencerewards`, `distribution` community pool) to prove that the token distribution is working exactly as expected in both the default and governance-controlled states.
- [ ] The CLI must be fully functional for querying the current emission parameters.
- [ ] The entire system must pass unit and integration tests covering the full lifecycle: default minting -> proposal submission -> voting -> parameter update -> dynamic minting. 
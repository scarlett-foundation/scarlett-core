---
title: "Alignment & Capital: How Scarlett V1’s Economy Differs from Bittensor"
date: "2025-07-23"
author: "The Scarlett Team"
tags: ["Scarlett", "Bittensor", "Economics", "Mechanism Design", "Decentralized AI"]
---

## The Challenge: Building Sustainable Decentralized Economies

In the world of decentralized AI, Bittensor stands as a pioneer. Its success has proven there is a massive appetite for community-owned AI infrastructure. However, pioneering systems also provide crucial lessons. One of the most significant challenges in any decentralized protocol is designing an economic engine that is sustainable, aligned, and fosters long-term growth.

When we designed Scarlett, we studied Bittensor’s model deeply. Its core economic loop, centered around "subnets" and the alpha token `dTAO`, revealed a fundamental incentive misalignment that creates significant headwinds for growth.

This post is an exercise in transparency. We want to clearly articulate the failings we observed in Bittensor’s model and explain how Scarlett V1’s economic architecture is specifically designed to solve them.

## Bittensor’s Economic Engine: The dTAO Subnet Model

To understand the problem, we first need to understand how Bittensor’s economy works at a high level.

1.  **Subnets**: These are specialized AI services on the Bittensor network (e.g., a text generation subnet, an image analysis subnet).
2.  **Miners**: The operators who provide the AI computation for a subnet. They are rewarded in the subnet's alpha token (e.g., `dTAO`).
3.  **Validators & Stakers**: TAO holders stake their TAO into a subnet's pool. This registers the subnet to the main chain and directs a portion of TAO emissions to it. Validators help run this process.
4.  **The Flow of Value**:
    *   The Bittensor root network emits TAO tokens to the subnets.
    *   Miners perform work and receive the subnet's alpha token (`dTAO`).
    *   **Crucially, miners must consistently sell their `dTAO` for TAO to cover operational costs (hardware, electricity) and to realize their profits.**

This creates a constant, downward pressure on the price of the subnet's token.

### The Problem: A Persistent Down-Only Feedback Loop

Imagine you are an external investor or a community member who believes in a specific subnet. You want to buy its `dTAO` token to support its growth. However, you are constantly fighting a losing battle against a torrent of sell pressure from two powerful sources:

1.  **Miners are forced to sell.** Their business model depends on it.
2.  **The root network itself sells.** The TAO emissions given to the subnet are effectively sold into the `dTAO`/TAO liquidity pool, further diluting external holders.

This creates what we call a **prohibitive down-only feedback loop**. For an external investor, buying the subnet's token is like trying to fill a bucket with a large hole in the bottom. The economics are structurally misaligned against anyone other than the insiders (miners and the root network itself). This model makes it incredibly difficult for individual subnets to cultivate their own healthy, independent economies and attract long-term investment.

## The Scarlett V1 Solution: Governance-Directed Capital Formation

Scarlett is built on the lessons learned from this model. We recognized that the most valuable and scarce resource a protocol has is its **emissions**. Emissions are the fuel for innovation, and how they are allocated is the single most important economic decision the network makes.

Instead of a complex and speculative staking game, Scarlett V1 implements a simple, transparent, and direct system: **Community Governance-Directed Emissions**.

1.  **A Single Source of Truth**: The `x/emissions` module controls the distribution of all new `SCLT` tokens.
2.  **Community as Venture Capitalist**: The community of `SCLT` token holders, through on-chain governance, acts as a decentralized venture capital fund. They collectively decide which projects, modules, or services should receive funding in the form of emissions.
3.  **Direct, Aligned Funding**: A developer building a new AI service on Scarlett can write a governance proposal outlining their vision and requesting a share of protocol emissions. If the community believes in the project, they vote to approve it. The developer then receives a direct stream of `SCLT` to fund their work.

### How This Solves the Problem

This model directly corrects the flaws of the Bittensor system:

1.  **Incentive Alignment**: There is no secondary token with misaligned incentives. The interests of the core protocol (`SCLT` holders) and the application-layer projects are one and the same: to increase the value and utility of the entire Scarlett ecosystem.
2.  **Transparent Capital Allocation**: The decision to fund a project is made in the open through a public governance vote. Every token holder can see why a project is being funded and can participate in the decision. This is a stark contrast to the opaque and speculative nature of subnet staking.
3.  **Sustainable Growth**: By directly funding promising projects with the native protocol token, Scarlett empowers them to build sustainable operations without being bled dry by a flawed economic model. This creates a positive feedback loop:
    *   Good projects get funded.
    *   They build valuable services on Scarlett.
    *   This drives demand and value for `SCLT`.
    *   A more valuable `SCLT` gives the community even more resources to fund the next wave of innovation.

To make the contrast clear, here is a summary of the key differences:

| Feature                | Bittensor (dTAO Model)                                                       | Scarlett (V1 Model)                                                              |
| :--------------------- | :--------------------------------------------------------------------------- | :------------------------------------------------------------------------------- |
| **Capital Allocation** | Speculative Staking into Subnets                                             | Direct Governance Vote for Projects                                              |
| **Mechanism**          | Open market for `dTAO` tokens directs `TAO` emissions.                       | On-chain proposals to the `x/emissions` module.                                  |
| **Incentive Alignment**| Aligned with subnet miners who consistently sell `dTAO` to cover costs.      | Aligned with `SCLT` holders who vote to maximize network value.                  |
| **Value Accrual**      | Value is captured in secondary `dTAO` tokens, which face constant dilution.    | Value accrues directly to the native `SCLT` token through increased utility.     |
| **Primary Risk**       | "Down-only feedback loop" stifles subnet growth and external investment.     | Potential for voter apathy or plutocracy in the governance process.              |

### Acknowledging V1's Limitations: The Need for Evolution

However, transparency also demands that we acknowledge the limitations of our own V1 design. While it solves the critical issues present in its predecessors, no protocol design is perfect, and every system has potential failure modes.

-   **Reliance on Formal Models**: Our governance model, while provably incentive-compatible *in theory*, rests on assumptions that all participants are rational and have perfect information. The real world is far messier. A truly robust system must be resilient against a spectrum of human behaviors, not just the idealized ones.
-   **The Threat of Stagnation**: Any fixed set of rules is a snapshot in time. A mechanism that is secure today can become a target for exploitation tomorrow as new strategies for gaming the system are discovered. A protocol that cannot evolve faster than its adversaries is destined to fail.
-   **The Risk of Plutocracy**: While our governance is more direct, it is still a coin-weighted voting system. We must remain vigilant that the network doesn't devolve into a plutocracy, where the largest token holders can dictate outcomes, stifling the innovation that comes from the broader community.

These are not trivial challenges. They are the fundamental problems that all decentralized networks must confront in the long run.

## Holding Ourselves to Account

We are not just building a protocol; we are building an economy. Being transparent about our economic design choices—including their potential weaknesses—is a core part of our commitment to the community.

The Bittensor model was a necessary and important first step. Scarlett V1 is the logical next step—an evolution toward a more aligned, transparent, and sustainable economic engine designed to fuel a truly decentralized AI ecosystem for the long term. And acknowledging its limitations is why we are already designing Scarlett V2. The goal is to build a protocol that doesn't just start with an aligned economy but has the built-in capacity to monitor itself, learn, and adapt, hardening its defenses against these very challenges over time.

## A Vision For The Future: The Evolution to Scarlett V2

Our journey doesn't end with V1. The ultimate goal is to create a protocol that is not just robust, but truly **anti-fragile**—a system that learns from stress and becomes stronger over time. This is the vision for Scarlett V2.

Where V1 was about fixing the foundational economic alignment, V2 is about embedding the ability to evolve directly into the protocol's DNA. It moves from a static design to a dynamic, self-improving system.

The table below shows the full evolutionary path:

| Feature                   | Bittensor                                     | Scarlett V1 (The Present)                                   | Scarlett V2 (The Future)                                              |
| :------------------------ | :-------------------------------------------- | :---------------------------------------------------------- | :-------------------------------------------------------------------- |
| **Core Paradigm**         | A market of competing subnets                 | A decentralized operating system                            | A self-improving, adaptive protocol                                   |
| **Capital Allocation**    | Speculative market-based staking            | Transparent governance voting                               | Data-driven portfolio optimization                                    |
| **Primary Defense**       | Complex, fixed rules                          | Formally proven (but static) incentive alignment            | **Adaptation & Learning**; the system evolves its own rules.          |
| **Governance Model**      | Opaque, plutocratic validator set             | Transparent coin-voting                                     | **Reputation-weighted** and multi-tiered (Constitutional, Legislative)  |
| **Key Failure Mode**      | Incentive misalignment & feedback loops       | Voter apathy & risk of plutocracy                           | Oracle manipulation & the inherent complexity of adaptive systems   |

This is our commitment: to continuously learn, adapt, and build in the open. We believe this is the only way to create a credibly neutral foundation for the future of decentralized intelligence that can stand the test of time. 
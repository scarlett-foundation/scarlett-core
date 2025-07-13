# 🎯 Proofofdegen Airdrop Module

## What is Proofofdegen?

Proofofdegen is a **patience-rewarded airdrop system** that creates a game where waiting longer gets you exponentially bigger rewards. Instead of getting a fixed amount of tokens, you get a growing share of an ever-expanding pool.

## 🏦 How It Works: The Shared Pool Model

### The Simple Version
Think of it like a **pizza that keeps growing** while people leave the table:
- Everyone starts with a fair slice based on their "weight" (achievements)  
- Every block, the pizza gets bigger (more emissions)
- When someone takes their slice and leaves, your slice gets bigger
- The longer you wait, the bigger your slice becomes!

### The Technical Version
1. **Emissions flow** from the main system into one shared module account
2. **All tokens accumulate** in this single pool (no individual allocations)
3. **Your claimable amount** = `(Total Pool × Your Weight) ÷ Total Remaining Weight`
4. **When you claim**, you get your current fair share and stop accumulating
5. **When others claim**, the remaining pool gets concentrated among fewer people

## 🎮 The Patience Game: Real Example

Let's say you're Alice with **weight = 100**:

### Week 1: Early Days
```
📊 Pool Status:
Total Pool: 10,000 SCLT
Unclaimed Wallets: 100 people
Total Weight: 10,000
Your Weight: 100 (1% of remaining)

💰 Your Claimable: (10,000 × 100) ÷ 10,000 = 100 SCLT
```

### Week 4: Some People Claim
```
📊 Pool Status:
Total Pool: 35,000 SCLT (more emissions!)
Unclaimed Wallets: 60 people (40 claimed and left)
Total Weight: 6,000 (claimed wallets removed)
Your Weight: 100 (now 1.67% of remaining!)

💰 Your Claimable: (35,000 × 100) ÷ 6,000 = 583 SCLT
🚀 Growth: 5.83x bigger than Week 1!
```

### Week 12: Diamond Hands
```
📊 Pool Status:
Total Pool: 80,000 SCLT (even more emissions!)
Unclaimed Wallets: 20 people (80 claimed and left)
Total Weight: 2,000 (most weight gone)
Your Weight: 100 (now 5% of remaining!)

💰 Your Claimable: (80,000 × 100) ÷ 2,000 = 4,000 SCLT  
🚀 Growth: 40x bigger than Week 1!
```

## 🔥 Why This Creates Amazing Incentives

### The Compound Effect
- **More emissions arrive** → Pool grows larger
- **People claim and leave** → Your percentage increases
- **Both effects multiply** → Exponential rewards for patience

### Game Theory Psychology
- **FOMO**: "Should I claim now or wait for even more?"
- **Diamond Hands**: True believers get rewarded most
- **Viral Growth**: You want MORE people to be eligible (grows emission allocation)
- **Network Effects**: Project success = bigger rewards

## 🎯 Key Insights

### It's NOT a Traditional Airdrop
❌ **Traditional**: "Here's your 100 tokens, done"  
✅ **Proofofdegen**: "Here's your evolving share of a growing, concentrating pool"

### It's NOT Individual Vesting  
❌ **Vesting**: "You earn 10 tokens per day regardless of others"  
✅ **Shared Pool**: "Your fair share depends on current pool size and who's left"

### It IS a Game of Chicken
- Everyone wants to be patient, but not TOO patient
- Creates natural tension and community engagement  
- Rewards genuine long-term conviction

## 🛠️ How to Interact

### Check Your Current Claimable Amount
```bash
scarlett-cored query proofofdegen eligible-amount [your-address]
```

**Response includes:**
- `claimable_amount`: How much you can claim right now
- `wallet_weight`: Your weight in the system
- `weight_percentage`: Your current percentage of remaining pool
- `total_unclaimed_wallets`: How many people haven't claimed yet

### Claim Your Share
```bash
scarlett-cored tx proofofdegen claim [your-address] --from [your-key]
```

**What happens:**
- You receive your current fair share immediately
- You stop accumulating future emissions (game over for you)
- The pool concentrates for remaining players

### Check Campaign Statistics
```bash
scarlett-cored query proofofdegen campaign-info
```

**See overall stats:**
- Total emissions received by the module
- Number of claimed vs unclaimed wallets
- Current module balance
- Total remaining weight

## 🏆 Strategy Guide

### For Maximum Rewards
- **Monitor the pool**: Use queries to track growth
- **Watch claim patterns**: See how your percentage is growing
- **Time your exit**: Balance greed vs risk
- **Diamond hands**: The most patient get the most rewards

### For Community Growth  
- **Share the vision**: More participants = larger emission allocation
- **Educate others**: Help people understand the game mechanics
- **Build together**: Project success increases everyone's rewards

## ⚠️ Important Notes

### One-Time Decision
- Once you claim, you're out of the game forever
- No second chances or re-entry
- Make sure you're ready before claiming

### Gas Required
- You need gas (transaction fees) to claim
- Queries are free and can be done anytime
- Plan ahead for gas costs

### Fair Launch Principles
- Weights based on achievements, not money
- No whales can buy bigger shares
- Community contributors get rewarded most

## 🚀 The Vision

Proofofdegen creates a **fair launch ecosystem** where:
- **Patient community members** get life-changing rewards
- **Active contributors** get recognized with higher weights  
- **Network effects** benefit everyone who participates
- **Long-term thinking** is rewarded over short-term greed

**Welcome to the patience game!** 🎯

---

*Built with ❤️ for the Scarlett community* 
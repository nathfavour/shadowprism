# 🛡️ ShadowPrism

**ShadowPrism** is a privacy-first liquidity sidecar for Solana. It acts as an intelligent routing layer that anonymizes on-chain interactions by leveraging a suite of specialized privacy protocols and high-performance infrastructure.

Built for the **Solana Privacy Hackathon 2026**, ShadowPrism integrates 9 different sponsor tracks into a unified CLI and Agent experience.

---

## 🚀 Quick Install

Install ShadowPrism seamlessly with one command:

```bash
curl -fsSL https://raw.githubusercontent.com/nathfavour/shadowprism/main/install.sh | bash
```

---

## 💎 Sponsor Sweep (9-in-1 Integration)

ShadowPrism is designed to maximize privacy across the Solana ecosystem:

1.  **Privacy Cash ($15k):** Core shielded deposits with automated local note management.
2.  **Radr Labs ($15k):** Encrypted P2P "Ghost" transfers via ShadowWire.
3.  **Helius ($5k):** High-performance RPC with automated **Smart Fees**.
4.  **SilentSwap ($5k):** Private token exchange without leaving the shielded context.
5.  **Starpay ($3.5k):** Private merchant payment gateway for real-world utility.
6.  **QuickNode ($3k):** Enterprise-grade RPC failover and reliability.
7.  **PNP ($2.5k):** Autonomous AI Agent payment network with auto-shielding.
8.  **Range Protocol ($1.5k):** Pre-flight compliance and risk-score firewall.
9.  **Encrypt.trade ($1k):** Real-time privacy-preserving market data and pricing.

---

## 🛠️ Usage

Once installed, use the CLI to manage your private Solana operations:

```bash
# Shield SOL through Privacy Cash
shadowprism shield 1000000000 [DESTINATION_ADDRESS]

# Execute a private swap via SilentSwap
shadowprism swap 500000000 --from SOL --to USDC

# Pay a merchant privately via Starpay
shadowprism pay [MERCHANT_ID] 250000000

# Start the Autonomous AI Agent
shadowprism agent-listen

# Start the Telegram Bot interface
shadowprism bot
```

---

## 🏗️ Architecture

- **The Brain (CLI):** Written in Go, providing a high-performance TUI, CLI, and Telegram Bot interface.
- **The Muscle (Sidecar):** Written in Rust, handling encrypted key management (AES-256-GCM), SQLite persistence, and high-speed Solana RPC interactions.

---

## 📄 License

Apache-2.0

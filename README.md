# 🛡️ ShadowPrism

**ShadowPrism** is a privacy-first liquidity sidecar for Solana. It acts as an intelligent routing layer that anonymizes on-chain interactions by leveraging a suite of specialized privacy protocols and high-performance infrastructure.

ShadowPrism integrates a suite of specialized privacy protocols into a unified CLI and Agent experience for the Solana ecosystem.

---

## 🚀 Quick Install

Install ShadowPrism seamlessly with one command:

```bash
curl -fsSL https://raw.githubusercontent.com/nathfavour/shadowprism/main/install.sh | bash
```

---

## 💎 Sponsor Sweep (9-in-1 Integration)

ShadowPrism is designed to maximize privacy across the Solana ecosystem:

1.  **Privacy Cash:** Core shielded deposits with automated local note management.
2.  **Radr Labs:** Encrypted P2P "Ghost" transfers via ShadowWire.
3.  **Helius:** High-performance RPC with automated **Smart Fees**.
4.  **SilentSwap:** Private token exchange without leaving the shielded context.
5.  **Starpay:** Private merchant payment gateway for real-world utility.
6.  **QuickNode:** Enterprise-grade RPC failover and reliability.
7.  **PNP:** Autonomous AI Agent payment network with auto-shielding.
8.  **Range Protocol:** Pre-flight compliance and risk-score firewall.
9.  **Encrypt.trade:** Real-time privacy-preserving market data and pricing.

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

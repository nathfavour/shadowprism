# 🛡️ ShadowPrism

**ShadowPrism** is a privacy-first liquidity sidecar for Solana. It acts as an intelligent routing layer that anonymizes on-chain interactions by leveraging a suite of specialized privacy protocols and high-performance infrastructure.

ShadowPrism integrates a suite of specialized privacy protocols into a unified CLI and Agent experience for the Solana ecosystem.

---

## 🤖 AI & Autonomous Agent Integration

ShadowPrism is more than a CLI; it's an AI-native privacy layer.

### **Current Capabilities**
*   **Conversational AI:** Integrated ShadowPrism AI Assistant (powered by `vibeaura`) providing real-time privacy advice, market sentiment, and technical guidance.
*   **Context-Aware Hints:** Every CLI command (`shield`, `swap`, `pay`) is followed by an intelligent AI insight tailored to the operation performed.
*   **Autonomous PNP Agent:** A dedicated background mode (`agent-listen`) that scans for and fulfills autonomous settlement requests without human intervention.
*   **Agent-to-Agent Portal:** Secure Telegram interface (`/agent`) for managing and simulating cross-agent liquidity flows via the PNP protocol.

### **Future Roadmap**
*   **Neural Routing:** AI-driven path optimization to select the most private and cost-effective mix of protocols based on network congestion.
*   **Zero-Knowledge Proof Generation:** Offloading ZK-proof computation to autonomous sidecar agents for mobile performance.
*   **Predictive Privacy:** Pre-emptive shielding of funds based on predicted user interaction patterns.

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

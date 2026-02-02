# ShadowPrism System Architecture

## 1. High-Level Overview

**ShadowPrism** is a privacy-first liquidity aggregator and infrastructure sidecar for the Solana ecosystem. It acts as a "Privacy Proxy" for developers, AI Agents, and power users, enabling them to route transactions through multiple privacy protocols programmatically.

The system utilizes a **Polyglot Sidecar Architecture**:
1.  **The Brain (CLI/Bot):** A Go process handling user intent, routing orchestration, and interface management.
2.  **The Muscle (Core Engine):** A Rust daemon executing on-chain transactions, managing encrypted keys, and interfacing with privacy protocols.

---

## 2. System Components

### A. Interface Layer (Go)
The entry point for all user interactions, providing three primary modes:
*   **Interactive TUI:** A terminal interface built with BubbleTea for manual operations.
*   **Telegram Bot:** An automated interface for remote interaction.
*   **Autonomous Agent:** A PNP-compatible agent for automated privacy-preserving settlements.

### B. Engine Layer (Rust)
The core execution environment responsible for:
*   **Protocol Adapters:** Specialized modules for Privacy Cash, Radr Labs, SilentSwap, and Starpay.
*   **State Management:** SQLite persistence for transaction tracking and privacy notes.
*   **Secure Keystore:** AES-256-GCM encrypted local storage for Solana keypairs.

---

## 3. Communication & Security

### Secure IPC
Communication between the Go and Rust processes is handled via **Unix Domain Sockets (UDS)** on supported platforms, falling back to localhost TCP.
*   **Authentication:** All requests require a Bearer Token generated at runtime.
*   **Isolation:** Sockets are created with restricted permissions (0700) to prevent local cross-user access.

### Compliance Firewall
All outbound transactions are routed through a pre-flight check integrated with **Range Protocol**. Risk scores are cached locally to minimize latency and maximize user privacy.

---

## 4. Internal API Contract

### `POST /v1/shield`
Initiates an anonymization transaction.
```json
{
  "amount_lamports": 1000000000,
  "destination_addr": "Pubkey...",
  "strategy": "privacy_cash | radr_p2p"
}
```

### `POST /v1/swap`
Executes a private token exchange.
```json
{
  "amount_lamports": 1000000000,
  "from_token": "SOL",
  "to_token": "USDC"
}
```

---

## 5. Security Model

1.  **Zero-Leak Compliance:** Compliance checks are proxied through the sidecar to prevent direct address harvesting by third-party providers.
2.  **Encrypted Secrets:** Private keys are never stored in plain text and are only decrypted in memory when required for signing.
3.  **Process Lifecycle:** The Go process monitors the Rust core daemon and ensures clean memory teardown upon exit.

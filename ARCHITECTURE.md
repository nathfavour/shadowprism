# ShadowPrism System Architecture

## 1. High-Level Overview

**ShadowPrism** is a privacy-first liquidity aggregator and infrastructure sidecar for the Solana ecosystem. It is designed to act as a "Privacy Proxy" for developers, AI Agents, and advanced power users, allowing them to route transactions through multiple privacy protocols (Privacy Cash, ShadowWire) programmatically.

The system utilizes a **Polyglot Sidecar Architecture**, leveraging **Go** for high-level orchestration, UI, and Agent logic, and **Rust** for low-level cryptographic operations and Solana program interactions.

### The "Sidecar" Pattern
To ensure maximum compatibility and performance, ShadowPrism runs as two distinct processes:
1.  **The Brain (CLI/Bot):** A lightweight Go process that handles user intent, routing strategies, and process management.
2.  **The Muscle (Core Engine):** A high-performance Rust daemon that executes the actual on-chain transactions via SDK integrations.

---

## 2. System Components

### A. The Interface Layer (Go)
*Located in: `/cli`*

The entry point for the user. It is built using the **Cobra** framework and serves three modes:
1.  **Interactive CLI:** A Terminal User Interface (TUI) built with **BubbleTea** for manual operations.
2.  **Telegram Bot:** An automated agent interface built with `telebot` for mobile/remote interaction.
3.  **Process Manager:** Responsible for locating, spawning, and managing the lifecycle of the Rust Core daemon.

**Key Libraries:**
* `spf13/cobra`: CLI structure.
* `charmbracelet/bubbletea`: TUI components.
* `gopkg.in/telebot.v3`: Telegram Bot API.
* `go-resty`: HTTP Client for communicating with Core.

### B. The Engine Layer (Rust)
*Located in: `/core`*

The heavy lifter. It runs as a local HTTP server and holds the specific logic for interacting with various Solana privacy protocols. It isolates the "dependency hell" of multiple Rust crates from the user's application logic.

**Key Libraries:**
* `axum`: High-performance HTTP server.
* `anchor-client` / `solana-sdk`: Blockchain interaction.
* `privacy-cash-sdk`: Integration for shielding funds.
* `shadow-wire`: Integration for encrypted P2P transfers (Radr Labs).
* `reqwest`: Interface for Range Protocol compliance checks.

---

## 3. Data Flow & Communication

Communication between the **Brain** (Go) and **Muscle** (Rust) happens via HTTP over `localhost`.

**Port:** `42069` (Default)
**Protocol:** JSON / REST

### Flow: The "Shielding" Transaction
1.  **User Input:** User triggers `/shield 1 SOL` via Telegram Bot.
2.  **Orchestration (Go):**
    * Bot parses the command.
    * Sidecar Manager ensures `shadowprism-core` is running.
    * Bot sends `POST /v1/shield` to `localhost:42069`.
3.  **Compliance (Rust):**
    * Middleware calls **Range Protocol** API to check the destination address risk score.
    * If `risk > threshold`, transaction is aborted.
4.  **Execution (Rust):**
    * Engine selects the provider (Privacy Cash default).
    * Constructs the specific instruction (CPI) for the Privacy Cash program.
    * Signs and broadcasts the transaction via **Helius RPC**.
5.  **Feedback (Go):**
    * Rust returns `tx_hash` and `privacy_route`.
    * Bot displays success animation to the user.

---

## 4. API Contract (Internal)

The Rust Core exposes the following endpoints for the CLI/Bot to consume.

### `GET /health`
Used by the Go Process Manager to verify the daemon is ready.
```json
{
  "status": "ready",
  "engine": "rust",
  "block_height": 2450123
}

POST /v1/shield
The primary instruction to anonymize funds.
Request:
{
  "amount_lamports": 1000000000,
  "destination_addr": "BuX...7z",
  "strategy": "mix_standard",
  "provider_override": "radr_shadow_wire" // Optional: "privacy_cash" | "radr"
}

Response:
{
  "status": "success",
  "tx_hash": "5xG...9z",
  "provider_used": "radr_shadow_wire",
  "risk_score": 0
}

5. Sponsor Integrations (Technical Implementation)
ShadowPrism uses an Adapter pattern (core/src/adapters/) to modularize sponsor tracks.
| Sponsor | Integration Point | Implementation Details |
|---|---|---|
| Privacy Cash | adapters::privacy_cash | Wraps the native Rust SDK to perform deposit/shielding actions. Used as the default routing provider. |
| Radr Labs | adapters::radr | Uses the ShadowWire SDK for P2P encrypted transfers. Activated via the p2p flag or strategy. |
| Helius | Core Configuration | Used as the primary RPC provider for transaction submission and DAS API for history lookups. |
| Range Protocol | middleware::compliance | Pre-flight check. Queries Range API to validate wallet reputation before interacting with privacy pools. |
6. Directory Structure
shadow-prism/
├── README.md               # Documentation
├── install.sh              # Unified build script
├── docker-compose.yml      # Container orchestration
│
├── core/                   # [RUST] Privacy Engine
│   ├── src/
│   │   ├── main.rs         # Axum Entrypoint
│   │   ├── api.rs          # HTTP Routes
│   │   └── adapters/       # Protocol Integrations
│   │       ├── mod.rs
│   │       ├── privacy_cash.rs
│   │       ├── radr.rs
│   │       └── range.rs
│   └── Cargo.toml
│
└── cli/                    # [GO] CLI & Bot
    ├── main.go             # Entrypoint
    ├── cmd/                # Cobra Commands
    │   ├── bot.go          # Telegram Listener
    │   └── send.go         # Terminal logic
    ├── internal/
    │   └── sidecar/        # Daemon Management
    │       ├── client.go   # HTTP Client
    │       └── manager.go  # Process Spawner
    └── go.mod

7. Security Considerations
 * Localhost Binding: The Rust API listens strictly on 127.0.0.1. It does not expose ports to the public internet to prevent unauthorized remote control.
 * Key Management: For the hackathon reference implementation, keys are loaded via standard Solana CLI config or Environment Variables (SOLANA_PRIVATE_KEY).
 * Ephemeral Binaries: When using the embedding strategy, the Rust binary is extracted to a secure temporary location with strict 0700 permissions.
<!-- end list -->


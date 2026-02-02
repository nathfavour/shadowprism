# NineLives: The "Magnificent 9" Integration Strategy

> **Mission:** Transform ShadowPrism into the ultimate "Sponsor Sweep" machine.
> **Architecture:** Go CLI (The Brain) + Rust Sidecar (The Muscle).
> **Status:** APPROVED for immediate execution.

---

## 1. Privacy Cash ($15k) - The "Shield"
**Role:** Core Privacy Provider
**File:** `core/src/adapters/privacy_cash.rs`

### Integration Logic
The default routing engine. When a user executes `/shield`, the Rust Sidecar must perform a Cross-Program Invocation (CPI) to the Privacy Cash program.

1.  **Input:** `amount` (u64), `recipient` (Pubkey).
2.  **Process:**
    * Load the Privacy Cash SDK (`privacy-cash-sdk`).
    * Construct a `Deposit` instruction targeting their on-chain Mixer.
    * Generate a "Note" (the private receipt) and return it to the Go CLI.
3.  **Output:** A transaction hash and a serialized Note string (for future withdrawals).

### Creativity Zone
* **"Note Management":** Instead of just dumping the note in the logs, store it locally in an encrypted SQLite file (`~/.shadowprism/wallet.db`) so the user can "Auto-Withdraw" later without copying/pasting secrets.

---

## 2. Radr Labs ($15k) - The "Ghost"
**Role:** P2P Transfer (ShadowWire)
**File:** `core/src/adapters/radr.rs`

### Integration Logic
Activated via the flag `/transfer --p2p` or `/shield --mode=ghost`. This bypasses the mixer and uses Radr's encrypted transfer protocol.

1.  **Input:** `amount`, `destination_address`.
2.  **Process:**
    * Check if `strategy == "radr"`.
    * Initialize the `shadow-wire` crate.
    * Execute an encrypted transfer where the link between sender and receiver is obfuscated on-chain.
3.  **Output:** A "Ghost Receipt" proving the funds moved without revealing the source.

### Creativity Zone
* **"Contact Book":** Map Solana addresses to human-readable aliases in the Go CLI. Allow the user to type `/shield @bob --p2p` and have the Rust engine resolve `@bob` to a Radr-compatible address.

---

## 3. Helius ($5k) - The "Eyes"
**Role:** Primary Infrastructure (RPC & DAS)
**File:** `core/src/adapters/rpc.rs`

### Integration Logic
Helius is the backbone. We do not just use it for sending transactions; we use it to "see."

1.  **Transaction Submission:** All `send_transaction` calls must route through the `HELIUS_RPC_URL` env var.
2.  **History Lookup:** Implement a `get_history(address)` function in Rust that queries the **Helius DAS API** (Digital Asset Standard) to find past shielded interactions or compressed NFTs (if used for notes).

### Creativity Zone
* **"Smart Fees":** Use Helius's Priority Fee API to automatically calculate the optimal compute unit price so the user's shield transaction never fails during congestion.

---

## 4. SilentSwap ($5k) - The "Swap"
**Role:** Private Token Exchange
**File:** `core/src/adapters/silent_swap.rs`

### Integration Logic
Enables the `/swap <amount> <from> <to>` command. Users shouldn't have to exit privacy to trade SOL for USDC.

1.  **Input:** `10 SOL`, `USDC`.
2.  **Process:**
    * Build a CPI instruction to the SilentSwap program.
    * **Crucial:** Do not route through Jupiter or Raydium public pools.
    * Execute the swap within the shielded context (if supported) or atomically (swap-and-shield).
3.  **Output:** The new token balance appearing in the user's shadow wallet.

### Creativity Zone
* **"Route Optimization":** If SilentSwap liquidity is low, implement a "Hybrid Route": Shield SOL -> Withdraw to Ephemeral Key -> Swap on Jupiter -> Re-shield USDC. (Complex, but impressive).

---

## 5. Starpay ($3.5k) - The "Merchant"
**Role:** Payment Gateway
**File:** `core/src/adapters/starpay.rs`

### Integration Logic
Enables the `/pay <merchant> <amount>` command. ShadowPrism acts as the "Private Bank" for Starpay merchants.

1.  **Input:** `merchant_id` (or wallet), `amount`.
2.  **Process:**
    * User holds funds in Shielded Pool.
    * Rust Engine performs a `Withdraw` action directly to the Starpay Settlement Address.
    * Sidecar sends a webhook/ping to Starpay API: "Payment Sent: <tx_hash>".
3.  **Output:** A "Payment Confirmed" receipt in the Telegram Bot.

### Creativity Zone
* **"Recurring Privacy":** Allow users to set up a `cron` job in the Go CLI that pays a Starpay merchant every month (e.g., for a VPN subscription) automatically from shielded funds.

---

## 6. QuickNode ($3k) - The "Backup"
**Role:** Redundancy & High Availability
**File:** `core/src/adapters/rpc.rs`

### Integration Logic
Demonstrate "Enterprise-Grade Reliability."

1.  **Logic:** Define `Primary = Helius`, `Secondary = QuickNode`.
2.  **Implementation:** Wrap the RPC client in a `ReliableClient` struct.
    * `client.send_tx()` -> Try Helius.
    * `if Error::Timeout` -> Log "Failover Active" -> Try QuickNode.
3.  **Proof:** Add a `--force-failover` flag to the CLI to demonstrate this switching capability live to judges.

### Creativity Zone
* **"Race Mode":** Send the transaction to *both* RPCs simultaneously. Whichever lands the block first wins. (Aggressive, but guarantees speed).

---

## 7. PNP ($2.5k) - The "Agent"
**Role:** Autonomous Payment Network
**File:** `cli/cmd/agent.go`

### Integration Logic
ShadowPrism isn't just a tool for humans; it's a wallet for AI Agents.

1.  **New Mode:** `shadowprism agent-listen`.
2.  **Process:**
    * Starts a server listening for "Payment Requests" (simulated or real PNP protocol).
    * When another bot "pings" your bot, ShadowPrism automatically:
        1.  Validates the request.
        2.  Shields the payment.
        3.  Sends it back.
3.  **Narrative:** "The Operating System for Agent Economies."

### Creativity Zone
* **"The Handshake":** Create a mock "Merchant Agent" script that spams your main bot with payment requests. Show the logs of your bot handling them autonomously without human input.

---

## 8. Range Protocol ($1.5k) - The "Firewall"
**Role:** Compliance & Security
**File:** `core/src/adapters/range.rs`

### Integration Logic
The "Pre-Flight Check." No transaction leaves the engine without passing this gate.

1.  **Hook:** Inside `execute_transaction()`.
2.  **Process:**
    * Take the `destination_address`.
    * Call Range Protocol API: `GET /risk/{address}`.
    * **Rule:** If `risk_score > 75` (High Risk), abort the transaction and warn the user.
3.  **Output:** "🛡️ Security Check Passed" or "⚠️ Transaction Blocked by Range Protocol".

### Creativity Zone
* **"Override Key":** Allow the user to provide a `--force` flag (with a scary warning) to bypass Range. This shows you care about both safety *and* user sovereignty.

---

## 9. Encrypt.trade ($1k) - The "Oracle"
**Role:** Data & Pricing
**File:** `core/src/adapters/market.rs`

### Integration Logic
Users think in Dollars, not Lamports.

1.  **Hook:** Before displaying any UI output (Balance, Shield Amount).
2.  **Process:**
    * Call Encrypt.trade API to get `SOL/USD`.
    * Convert the `lamports` amount to formatted USD.
3.  **Display:** "Shielding 1.5 SOL (~$210.45)" instead of just "1.5 SOL".

### Creativity Zone
* **"Arbitrage Alert":** If Encrypt.trade shows a significant price difference between the Shielded Pool implied rate and the market rate, alert the user (rare, but cool feature).

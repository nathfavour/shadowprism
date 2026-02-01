use rusqlite::{params, Connection, Result};
use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::{DateTime, Utc};
use std::path::Path;

pub mod watchdog;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct TransactionRecord {
    pub id: String,
    pub amount_lamports: u64,
    pub destination: String,
    pub status: String, // Pending, Broadcast, Confirmed, Failed
    pub tx_hash: Option<String>,
    pub provider: String,
    pub created_at: DateTime<Utc>,
}

pub struct TransactionStore {
    conn: Connection,
}

impl TransactionStore {
    pub fn new(path: &Path) -> Result<Self> {
        let conn = Connection::open(path)?;
        
        // Initialize Schema
        conn.execute(
            "CREATE TABLE IF NOT EXISTS transactions (
                id TEXT PRIMARY KEY,
                amount_lamports INTEGER NOT NULL,
                status TEXT NOT NULL,
                destination TEXT NOT NULL,
                tx_hash TEXT,
                provider TEXT NOT NULL,
                created_at TEXT NOT NULL
            )",
            [],
        )?;

        Ok(Self { conn })
    }

    pub fn create_transaction(&self, amount: u64, dest: &str, provider: &str) -> Result<String> {
        let id = Uuid::new_v4().to_string();
        let now = Utc::now();
        
        self.conn.execute(
            "INSERT INTO transactions (id, amount_lamports, destination, status, provider, created_at)
             VALUES (?1, ?2, ?3, ?4, ?5, ?6)",
            params![id, amount as i64, dest, "Pending", provider, now.to_rfc3339()],
        )?;

        Ok(id)
    }

    pub fn update_status(&self, id: &str, status: &str, hash: Option<&str>) -> Result<()> {
        self.conn.execute(
            "UPDATE transactions SET status = ?1, tx_hash = ?2 WHERE id = ?3",
            params![status, hash, id],
        )?;
        Ok(())
    }

    pub fn get_transaction(&self, id: &str) -> Result<TransactionRecord> {
        self.conn.query_row(
            "SELECT id, amount_lamports, destination, status, tx_hash, provider, created_at FROM transactions WHERE id = ?1",
            params![id],
            |row| {
                let amount_i64: i64 = row.get(1)?;
                let created_at_str: String = row.get(6)?;
                Ok(TransactionRecord {
                    id: row.get(0)?,
                    amount_lamports: amount_i64 as u64,
                    destination: row.get(2)?,
                    status: row.get(3)?,
                    tx_hash: row.get(4)?,
                    provider: row.get(5)?,
                    created_at: DateTime::parse_from_rfc3339(&created_at_str).unwrap().with_timezone(&Utc),
                })
            },
        )
    }

    pub fn list_transactions(&self) -> Result<Vec<TransactionRecord>> {
        let mut stmt = self.conn.prepare(
            "SELECT id, amount_lamports, destination, status, tx_hash, provider, created_at FROM transactions ORDER BY created_at DESC LIMIT 50"
        )?;
        
        let rows = stmt.query_map([], |row| {
            let amount_i64: i64 = row.get(1)?;
            let created_at_str: String = row.get(6)?;
            Ok(TransactionRecord {
                id: row.get(0)?,
                amount_lamports: amount_i64 as u64,
                destination: row.get(2)?,
                status: row.get(3)?,
                tx_hash: row.get(4)?,
                provider: row.get(5)?,
                created_at: DateTime::parse_from_rfc3339(&created_at_str).unwrap().with_timezone(&Utc),
            })
        })?;

        let mut results = Vec::new();
        for row in rows {
            results.push(row?);
        }
        Ok(results)
    }
}

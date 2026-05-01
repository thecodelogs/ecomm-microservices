use std::env;

/// Application configuration loaded from environment variables.
#[derive(Debug, Clone)]
pub struct Config {
    /// Full Postgres connection string, e.g. `postgres://user:pass@host/db`
    pub database_url: String,

    /// Port on which the gRPC server listens (default: 50054)
    pub grpc_port: u16,
}

impl Config {
    pub fn from_env() -> anyhow::Result<Self> {
        let database_url = env::var("DATABASE_URL")
            .map_err(|_| anyhow::anyhow!("DATABASE_URL environment variable is required"))?;

        let grpc_port = env::var("GRPC_PORT")
            .unwrap_or_else(|_| "50054".to_string())
            .parse::<u16>()
            .map_err(|e| anyhow::anyhow!("Invalid GRPC_PORT: {}", e))?;

        Ok(Self {
            database_url,
            grpc_port,
        })
    }
}

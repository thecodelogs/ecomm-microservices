use deadpool_postgres::{Config as PgConfig, ManagerConfig, Pool, RecyclingMethod, Runtime};
use tokio_postgres::NoTls;
use tracing::info;

/// Build an async connection pool from a Postgres connection string.
pub async fn create_pool(database_url: &str) -> anyhow::Result<Pool> {
    // Parse the URL into deadpool's config
    let mut cfg = PgConfig::new();
    cfg.url = Some(database_url.to_string());
    cfg.manager = Some(ManagerConfig {
        recycling_method: RecyclingMethod::Fast,
    });

    let pool = cfg
        .create_pool(Some(Runtime::Tokio1), NoTls)
        .map_err(|e| anyhow::anyhow!("Failed to create DB pool: {}", e))?;

    // Eagerly verify connectivity
    let _client = pool
        .get()
        .await
        .map_err(|e| anyhow::anyhow!("Failed to connect to database: {}", e))?;

    info!("Database pool created successfully");
    Ok(pool)
}

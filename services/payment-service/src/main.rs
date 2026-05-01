mod config;
mod db;
mod grpc;
mod proto;

use std::net::SocketAddr;

use dotenvy::dotenv;
use tonic::transport::Server;
use tracing::{info, warn};

use proto::payment::payment_service_server::PaymentServiceServer;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // ── Logging ──────────────────────────────────────────────────────────────
    tracing_subscriber::fmt()
        .with_env_filter(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| tracing_subscriber::EnvFilter::new("info")),
        )
        .init();

    // ── Config ───────────────────────────────────────────────────────────────
    if dotenv().is_err() {
        warn!("No .env file found, reading environment variables directly");
    }
    let cfg = config::Config::from_env()?;

    // ── Database ─────────────────────────────────────────────────────────────
    let pool = db::create_pool(&cfg.database_url).await?;
    info!("Connected to database");

    // ── gRPC server ──────────────────────────────────────────────────────────
    let addr: SocketAddr = format!("0.0.0.0:{}", cfg.grpc_port).parse()?;
    let payment_handler = grpc::PaymentHandler::new(pool);

    info!("payment-service gRPC server starting on {}", addr);

    let reflection_service = tonic_reflection::server::Builder::configure()
        .register_encoded_file_descriptor_set(proto::payment::FILE_DESCRIPTOR_SET)
        .build_v1()
        .expect("Failed to build reflection service");

    Server::builder()
        .add_service(reflection_service)
        .add_service(PaymentServiceServer::new(payment_handler))
        .serve_with_shutdown(addr, shutdown_signal())
        .await?;

    info!("payment-service shut down gracefully");
    Ok(())
}

/// Listens for SIGINT / SIGTERM and resolves when either is received.
async fn shutdown_signal() {
    use tokio::signal;

    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    };

    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install SIGTERM handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {}
        _ = terminate => {}
    }
}

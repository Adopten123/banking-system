use sqlx::postgres::PgPoolOptions;
use std::sync::Arc;
use tonic::transport::Server;

pub mod config;
pub mod domain;
pub mod handler;
pub mod infrastructure;
pub mod service;

use crate::handler::grpc_server::CardVaultGrpcServer;
use crate::handler::pb::v1::card_vault_service_server::CardVaultServiceServer;
use crate::infrastructure::crypto_service::CryptoService;
use crate::infrastructure::postgres_repo::PostgresCardRepository;
use crate::service::use_cases::CardVaultUseCase;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // 1. Load cfg
    let settings = config::Settings::load()
        .expect("Failed to load configuration. Check your YAML file or environment variables.");

    println!("Starting Card Vault...");

    // 2. Init DB
    let pool = PgPoolOptions::new()
        .max_connections(settings.database.max_connections)
        .connect(&settings.database.url)
        .await
        .expect("Failed to connect to the Vault database");

    // 3. Building repo
    let repo = Arc::new(PostgresCardRepository::new(pool));

    let crypto = Arc::new(
        CryptoService::new(&settings.security.master_key)
            .expect("Cryptography initialization error: Check the master key length (32 bytes)")
    );

    let use_case = Arc::new(CardVaultUseCase::new(repo, crypto));
    let grpc_handler = CardVaultGrpcServer::new(use_case);

    // 4. Starting gRPC-server
    let addr = format!("0.0.0.0:{}", settings.server.grpc_port).parse()?;
    println!("Card Vault server has been successfully started on port {}", settings.server.grpc_port);

    Server::builder()
        .add_service(CardVaultServiceServer::new(grpc_handler))
        .serve(addr)
        .await?;

    Ok(())
}
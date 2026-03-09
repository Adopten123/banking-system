use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Settings {
    pub server: ServerSettings,
    pub database: DatabaseSettings,
    pub security: SecuritySettings,
}

#[derive(Debug, Deserialize)]
pub struct ServerSettings {
    pub grpc_port: u16,
}

#[derive(Debug, Deserialize)]
pub struct DatabaseSettings {
    pub url: String,
    pub max_connections: u32,
}

#[derive(Debug, Deserialize)]
pub struct SecuritySettings {
    pub master_key: String,
}

impl Settings {
    pub fn load() -> Result<Self, config::ConfigError> {
        let config_path = std::env::var("CONFIG_PATH").unwrap_or_else(|_| "config/local.yaml".into());

        let settings = config::Config::builder()

            .add_source(config::File::with_name(&config_path))
            .add_source(config::Environment::with_prefix("APP").separator("__"))
            .build()?;

        settings.try_deserialize()
    }
}
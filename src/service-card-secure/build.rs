fn main() -> Result<(), Box<dyn std::error::Error>> {
    let proto_file = "../../proto/card_vault/v1/card_vault.proto";

    tonic_build::compile_protos(proto_file)?;

    println!("cargo:rerun-if-changed={}", proto_file);

    Ok(())
}
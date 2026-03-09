use std::sync::Arc;
use tonic::{Request, Response, Status};
use uuid::Uuid;

use crate::domain::errors::VaultError;
use crate::service::use_cases::CardVaultUseCase;

use crate::handler::pb::v1::card_vault_service_server::CardVaultService;
use crate::handler::pb::v1::*;

pub struct CardVaultGrpcServer {
    use_case: Arc<CardVaultUseCase>,
}

impl CardVaultGrpcServer {
    pub fn new(use_case: Arc<CardVaultUseCase>) -> Self {
        Self { use_case }
    }

    fn map_error(err: VaultError) -> Status {
        match err {
            VaultError::CardNotFound => Status::not_found(err.to_string()),
            VaultError::InvalidPin => Status::permission_denied(err.to_string()),
            VaultError::CardBlocked => Status::permission_denied(err.to_string()),
            _ => Status::internal(err.to_string()),
        }
    }

    fn parse_uuid(token: &str) -> Result<Uuid, Status> {
        Uuid::parse_str(token).map_err(|_| Status::invalid_argument("Некорректный формат token_id (ожидается UUID)"))
    }
}

#[tonic::async_trait]
impl CardVaultService for CardVaultGrpcServer {
    async fn issue_card(
        &self,
        request: Request<IssueCardRequest>,
    ) -> Result<Response<IssueCardResponse>, Status> {
        let req = request.into_inner();

        let result = self.use_case
            .issue_card(&req.payment_system, req.is_virtual)
            .await
            .map_err(Self::map_error)?;

        Ok(Response::new(IssueCardResponse {
            token_id: result.token_id.to_string(),
            pan_mask: result.pan_mask,
            expiry_month: result.expiry_month,
            expiry_year: result.expiry_year,
        }))
    }

    async fn get_card_details(
        &self,
        request: Request<GetCardDetailsRequest>,
    ) -> Result<Response<GetCardDetailsResponse>, Status> {
        let token_id = Self::parse_uuid(&request.into_inner().token_id)?;

        let details = self.use_case
            .get_card_details(&token_id)
            .await
            .map_err(Self::map_error)?;

        Ok(Response::new(GetCardDetailsResponse {
            pan: details.pan,
            cvv: details.cvv,
            expiry_month: details.expiry_month,
            expiry_year: details.expiry_year,
        }))
    }

    async fn verify_card(
        &self,
        request: Request<VerifyCardRequest>,
    ) -> Result<Response<VerifyCardResponse>, Status> {
        let req = request.into_inner();

        let (is_valid, token_id_opt) = self.use_case
            .verify_card(&req.pan, &req.cvv, req.expiry_month, req.expiry_year)
            .await
            .map_err(Self::map_error)?;

        Ok(Response::new(VerifyCardResponse {
            is_valid,
            token_id: token_id_opt.map(|u| u.to_string()).unwrap_or_default(),
        }))
    }

    async fn update_card_status(
        &self,
        request: Request<UpdateCardStatusRequest>,
    ) -> Result<Response<UpdateCardStatusResponse>, Status> {
        let req = request.into_inner();
        let token_id = Self::parse_uuid(&req.token_id)?;

        self.use_case
            .update_status(&token_id, &req.status)
            .await
            .map_err(Self::map_error)?;

        Ok(Response::new(UpdateCardStatusResponse { success: true }))
    }

    async fn delete_card_data(
        &self,
        request: Request<DeleteCardDataRequest>,
    ) -> Result<Response<DeleteCardDataResponse>, Status> {
        let token_id = Self::parse_uuid(&request.into_inner().token_id)?;

        self.use_case
            .delete_card(&token_id)
            .await
            .map_err(Self::map_error)?;

        Ok(Response::new(DeleteCardDataResponse { success: true }))
    }

    async fn set_pin(
        &self,
        request: Request<SetPinRequest>,
    ) -> Result<Response<SetPinResponse>, Status> {
        let req = request.into_inner();
        let token_id = Self::parse_uuid(&req.token_id)?;

        self.use_case
            .set_pin(&token_id, &req.pin)
            .await
            .map_err(Self::map_error)?;

        Ok(Response::new(SetPinResponse { success: true }))
    }

    async fn verify_pin(
        &self,
        request: Request<VerifyPinRequest>,
    ) -> Result<Response<VerifyPinResponse>, Status> {
        let req = request.into_inner();
        let token_id = Self::parse_uuid(&req.token_id)?;

        let is_valid = self.use_case
            .verify_pin(&token_id, &req.pin)
            .await
            .map_err(Self::map_error)?;

        Ok(Response::new(VerifyPinResponse { is_valid }))
    }
}
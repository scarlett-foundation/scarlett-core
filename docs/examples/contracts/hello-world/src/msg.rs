use cosmwasm_schema::{cw_serde, QueryResponses};

#[cw_serde]
pub struct InstantiateMsg {
    pub greeting: String,
}

#[cw_serde]
pub enum ExecuteMsg {
    UpdateGreeting { greeting: String },
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// Get the current greeting
    #[returns(GreetingResponse)]
    GetGreeting {},
    
    /// Get the contract state
    #[returns(StateResponse)]
    GetState {},
}

#[cw_serde]
pub struct GreetingResponse {
    pub greeting: String,
}

#[cw_serde]
pub struct StateResponse {
    pub greeting: String,
    pub owner: String,
}

#[cw_serde]
pub struct MigrateMsg {} 
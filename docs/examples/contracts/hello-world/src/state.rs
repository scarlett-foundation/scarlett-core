use cosmwasm_std::Addr;
use cw_storage_plus::Item;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct State {
    pub greeting: String,
    pub owner: Addr,
}

// Storage key for the main state
pub const STATE: Item<State> = Item::new("state"); 
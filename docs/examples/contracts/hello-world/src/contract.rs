use cosmwasm_std::{
    entry_point, to_json_binary, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult,
};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg, GreetingResponse, StateResponse, MigrateMsg};
use crate::state::{State, STATE};

// Contract name and version used for migration
const CONTRACT_NAME: &str = "crates.io:hello-world";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[entry_point]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    let state = State {
        greeting: msg.greeting.clone(),
        owner: info.sender.clone(),
    };

    // Validate greeting is not empty
    if msg.greeting.trim().is_empty() {
        return Err(ContractError::EmptyGreeting {});
    }

    // Store the state
    STATE.save(deps.storage, &state)?;

    // Set contract version for migration
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("owner", info.sender.as_str())
        .add_attribute("greeting", msg.greeting))
}

#[entry_point]
pub fn execute(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::UpdateGreeting { greeting } => execute_update_greeting(deps, info, greeting),
    }
}

pub fn execute_update_greeting(
    deps: DepsMut,
    info: MessageInfo,
    greeting: String,
) -> Result<Response, ContractError> {
    // Load the current state
    let mut state = STATE.load(deps.storage)?;

    // Check if sender is the owner
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {});
    }

    // Validate greeting is not empty
    if greeting.trim().is_empty() {
        return Err(ContractError::EmptyGreeting {});
    }

    // Update the greeting
    state.greeting = greeting.clone();
    STATE.save(deps.storage, &state)?;

    Ok(Response::new()
        .add_attribute("method", "update_greeting")
        .add_attribute("sender", info.sender.as_str())
        .add_attribute("greeting", greeting))
}

#[entry_point]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetGreeting {} => to_json_binary(&query_greeting(deps)?),
        QueryMsg::GetState {} => to_json_binary(&query_state(deps)?),
    }
}

fn query_greeting(deps: Deps) -> StdResult<GreetingResponse> {
    let state = STATE.load(deps.storage)?;
    Ok(GreetingResponse {
        greeting: state.greeting,
    })
}

fn query_state(deps: Deps) -> StdResult<StateResponse> {
    let state = STATE.load(deps.storage)?;
    Ok(StateResponse {
        greeting: state.greeting,
        owner: state.owner.to_string(),
    })
}

#[entry_point]
pub fn migrate(_deps: DepsMut, _env: Env, _msg: MigrateMsg) -> Result<Response, ContractError> {
    Ok(Response::default())
}

#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env};
    use cosmwasm_std::{Addr, coins, from_json, MessageInfo};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            greeting: "Hello, CosmWasm!".to_string(),
        };
        let info = MessageInfo {
            sender: Addr::unchecked("creator"),
            funds: coins(1000, "earth"),
        };

        // We can just call .unwrap() to assert this was a success
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        // It worked, let's query the state
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetGreeting {}).unwrap();
        let value: GreetingResponse = from_json(&res).unwrap();
        assert_eq!("Hello, CosmWasm!", value.greeting);
    }

    #[test]
    fn update_greeting() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            greeting: "Hello, CosmWasm!".to_string(),
        };
        let info = MessageInfo {
            sender: Addr::unchecked("creator"),
            funds: coins(2, "token"),
        };
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Update greeting
        let info = MessageInfo {
            sender: Addr::unchecked("creator"),
            funds: coins(2, "token"),
        };
        let msg = ExecuteMsg::UpdateGreeting {
            greeting: "Hello, World!".to_string(),
        };

        let res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        // Check that greeting was updated
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetGreeting {}).unwrap();
        let value: GreetingResponse = from_json(&res).unwrap();
        assert_eq!("Hello, World!", value.greeting);
    }

    #[test]
    fn unauthorized_update() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            greeting: "Hello, CosmWasm!".to_string(),
        };
        let info = MessageInfo {
            sender: Addr::unchecked("creator"),
            funds: coins(2, "token"),
        };
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Try to update greeting with different sender
        let info = MessageInfo {
            sender: Addr::unchecked("anyone"),
            funds: coins(2, "token"),
        };
        let msg = ExecuteMsg::UpdateGreeting {
            greeting: "Hello, World!".to_string(),
        };

        let res = execute(deps.as_mut(), mock_env(), info, msg);
        match res {
            Err(ContractError::Unauthorized {}) => {}
            _ => panic!("Must return unauthorized error"),
        }
    }

    #[test]
    fn empty_greeting_fails() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            greeting: "".to_string(),
        };
        let info = MessageInfo {
            sender: Addr::unchecked("creator"),
            funds: coins(1000, "earth"),
        };

        let res = instantiate(deps.as_mut(), mock_env(), info, msg);
        match res {
            Err(ContractError::EmptyGreeting {}) => {}
            _ => panic!("Must return empty greeting error"),
        }
    }
} 
use cosmwasm_std::{
    entry_point, Binary, Deps, DepsMut, Empty, Env, MessageInfo, Response, StdResult, to_json_binary,
};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct QueryResponse {
    pub message: String,
}

#[entry_point]
pub fn instantiate(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: Empty,
) -> StdResult<Response> {
    Ok(Response::new())
}

#[entry_point]
pub fn query(_deps: Deps, _env: Env, _msg: Empty) -> StdResult<Binary> {
    let response = QueryResponse {
        message: "Hello from CosmWasm contract".to_string(),
    };
    to_json_binary(&response)
}

#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env};
    use cosmwasm_std::{from_json, Empty, Addr};
    use cw_multi_test::{App, ContractWrapper, Executor};

    #[test]
    fn test_query() {
        let deps = mock_dependencies();
        let env = mock_env();
        let msg = Empty {};

        let result = query(deps.as_ref(), env, msg).unwrap();
        let response: QueryResponse = from_json(result).unwrap();
        
        assert_eq!(response.message, "Hello from CosmWasm contract");
    }

    #[test]
    fn test_multitest_query() {
        let mut app = App::default();

        // Create the contract wrapper
        let code = ContractWrapper::new(
            |_, _, _, _: Empty| -> Result<_, cosmwasm_std::StdError> {
                Ok(Response::new())
            },
            |_, _, _, _: Empty| -> Result<_, cosmwasm_std::StdError> {
                Ok(Response::new())
            },
            query,
        );

        // Store the contract code
        let code_id = app.store_code(Box::new(code));

        // Instantiate the contract
        let contract_addr = app
            .instantiate_contract(
                code_id,
                Addr::unchecked("owner"),
                &Empty {},
                &[],
                "Tutorial Contract",
                None,
            )
            .unwrap();

        // Query the contract
        let response: QueryResponse = app
            .wrap()
            .query_wasm_smart(contract_addr, &Empty {})
            .unwrap();

        assert_eq!(response.message, "Hello from CosmWasm contract");
    }
}

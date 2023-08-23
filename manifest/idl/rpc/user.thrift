namespace go user

include "../base/common.thrift"

struct registryRequest {
    1: required string username,
    2: required string password,
}

struct registryResponse {
    1: common.BaseResponse base_resp,
}

struct AuthReq {
    1: required string username
    2: required string password
}

struct AuthResp {
    1: common.BaseResponse base_resp,
    2: required string token
}



service UserService {
    AuthResp Auth(1:AuthReq req)
    registryResponse registry(1: registryRequest req),
}
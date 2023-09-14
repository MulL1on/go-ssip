namespace go user

include "../base/common.thrift"

struct RegisterReq {
    1: required string username,
    2: required string password,
}

struct RegisterResp {
    1: common.BaseResponse base_resp,
}

struct AuthReq {
    1: required string username
    2: required string password
}

struct AuthResp {
    1: common.BaseResponse base_resp
    2: required string token
}



service UserService {
    AuthResp Auth(1:AuthReq req)
    RegisterResp Register(1: RegisterReq req),
}
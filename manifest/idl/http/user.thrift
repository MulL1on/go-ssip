namespace go user

include "../base/common.thrift"


struct RegisterReq {
    1: required string username (api.raw="username" api.vd="len($)>0 && lne($)<33")
    2: required string password (api.raw="password" api.vd="len($)>0 && lne($)<33")
}

struct LoginReq {
    1: required string username(api.raw="username" api.vd="len($)>0 && lne($)<33")
    2: required string password(api.raw="username" api.vd="len($)>0 && lne($)<33")
}

service UserService {
    common.NilResponse Register(1: RegisterReq req) (api.POST="/user/register")
    common.NilResponse Login(1: LoginReq req)(api.POST="/user/login")
}
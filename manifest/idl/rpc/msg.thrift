namespace go msg

include "../base/common.thrift"
include "../base/msg.thrift"

struct SendMsgReq {
    1: msg.Msg msg
}

struct SendMsgResp {
    1: common.BaseResponse base_resp
}

struct GetMsgReq {
    1: i64 user_id
    2: i64 seq
}

struct GetMsgResp {
    1: common.BaseResponse base_resp
}

service MsgService {
    SendMsgResp sendMsg (1: SendMsgReq req)
    GetMsgResp getMsg(1: GetMsgReq req)
}
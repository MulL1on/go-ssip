namespace go msg

include "../base/common.thrift"
include "../base/msg.thrift"

struct SendMsgReq {
    1: msg.Msg msg
}

struct SendMsgResp {
    1: common.BaseResponse base_resp
}

struct SendGroupMsgReq{
    1: msg.GroupMsg group_msg

}

struct SendGroupMsgResp {
    1: common.BaseResponse base_resp
}

struct GetMsgReq {
    1: i64 from_user
    2: i64 to_user
}

struct GetMsgResp {
    1: list<msg.Msg> msgs
    2: common.BaseResponse base_resp
}

struct GetGroupMsgReq{
    1: i64 group
}

struct GetGroupMsgResp{
    1: list<msg.GroupMsg> group_msgs
    2: common.BaseResponse base_resp
}

service MsgService {
    SendMsgResp sendMsg (1: SendMsgReq req)
    SendGroupMsgResp sendGroupMsg(1: SendGroupMsgReq req)
    GetMsgResp getMsg(1: GetMsgReq req)
    GetGroupMsgResp getGroupMsg(1:GetGroupMsgReq req)
}
namespace go group

include "../base/common.thrift"

struct CreateGroupReq {
    1: string group_name
}

struct CreateGroupResp{
    1: common.BaseResponse base_resp
}

struct  JoinGroupReq{
  1: i64 group_id
  2: i64 user_id
}

struct JoinGroupResp{
  1: common.BaseResponse base_resp
}

struct QuitGroupReq {
    1: i64 group_id
    2: i64 user_id
}

struct QuitGroupResp{
  1: common.BaseResponse base_resp
}

service GroupService{
  CreateGroupResp CreateGroup(1: CreateGroupReq req)
  JoinGroupResp JoinGroup(1: JoinGroupReq req)
  QuitGroupResp QuitGroup(1: QuitGroupReq req)
}




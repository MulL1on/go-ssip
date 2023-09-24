namespace go group

include "../base/common.thrift"


struct CreateGroupReq {
    1: required string group_name (api.raw="group_name" api.vd="len($)>0 && len($)<33")
}

struct JoinGroupReq {
    1: required i64 group_id(api.raw="group_id" )
}

struct QuitGroupReq {
    1: required i64 group_id(api.raw="group_id" )
}

service GroupService {
    common.NilResponse CreateGroup(1: CreateGroupReq req) (api.POST="/group")
    common.NilResponse JoinGroup(1: JoinGroupReq req)(api.POST="/group/member")
    common.NilResponse QuitGroup(1: QuitGroupReq req)(api.DELETE="/group/member")
}
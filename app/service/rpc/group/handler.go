package main

import (
	"context"
	"go-ssip/app/common/errno"
	"go-ssip/app/common/kitex_gen/group"
	"go-ssip/app/common/tools"
	g "go-ssip/app/service/rpc/group/global"
	"go-ssip/app/service/rpc/group/pkg/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GroupServiceImpl implements the last service interface defined in the IDL.
type GroupServiceImpl struct {
	MysqlManager
	IDGenerator
}

type IDGenerator interface {
	CreateUUID() int64
}

type MysqlManager interface {
	CreateGroup(g *mysql.Group) error
	GetGroupByName(name string) (*mysql.Group, error)
	GetGroupById(id int64) (*mysql.Group, error)
	CreateGroupMember(gm *mysql.GroupMember) error
	DeleteGroupMember(groupID, userID int64) error
	GetGroupMember(groupID, userID int64) error
}

// CreateGroup implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) CreateGroup(ctx context.Context, req *group.CreateGroupReq) (resp *group.CreateGroupResp, err error) {
	resp = new(group.CreateGroupResp)
	_, err = s.MysqlManager.GetGroupByName(req.GroupName)
	if err == nil {
		resp.BaseResp = tools.BuildBaseResp(errno.RecordAlreadyExist)
		return resp, nil
	}

	err = s.MysqlManager.CreateGroup(&mysql.Group{
		GroupID:   s.IDGenerator.CreateUUID(),
		GroupName: req.GroupName,
	})
	if err != nil {
		g.Logger.Error("create group error", zap.Error(err))
		resp.BaseResp = tools.BuildBaseResp(errno.GroupSrvErr)
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(errno.Success)
	return resp, nil
}

// JoinGroup implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) JoinGroup(ctx context.Context, req *group.JoinGroupReq) (resp *group.JoinGroupResp, err error) {
	// check group
	resp = new(group.JoinGroupResp)

	var gm = &mysql.GroupMember{
		GroupID: req.GroupId,
		UserID:  req.UserId,
	}

	_, err = s.MysqlManager.GetGroupById(req.GroupId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.BaseResp = tools.BuildBaseResp(errno.RecordNotFound)
			return resp, nil
		}
		g.Logger.Error("get group by id error", zap.Error(err))
		resp.BaseResp = tools.BuildBaseResp(errno.GroupSrvErr)
		return resp, nil
	}

	// check is group member
	err = s.MysqlManager.GetGroupMember(gm.GroupID, gm.UserID)
	if err == nil {
		resp.BaseResp = tools.BuildBaseResp(errno.RecordAlreadyExist)
		return resp, nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		g.Logger.Error("get group member error", zap.Error(err))
		resp.BaseResp = tools.BuildBaseResp(errno.GroupSrvErr)
		return resp, nil
	}

	// insert

	err = s.MysqlManager.CreateGroupMember(gm)
	if err != nil {
		g.Logger.Error("create group member error", zap.Error(err))
		resp.BaseResp = tools.BuildBaseResp(errno.GroupSrvErr)
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(errno.Success)
	return resp, nil
}

// QuitGroup implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) QuitGroup(ctx context.Context, req *group.QuitGroupReq) (resp *group.QuitGroupResp, err error) {
	// check is in group
	resp = new(group.QuitGroupResp)

	err = s.MysqlManager.GetGroupMember(req.GroupId, req.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.BaseResp = tools.BuildBaseResp(errno.RecordNotFound)
			return resp, nil
		}
		g.Logger.Error("get group member error", zap.Error(err))
		resp.BaseResp = tools.BuildBaseResp(errno.GroupSrvErr)
		return resp, nil
	}

	err = s.MysqlManager.DeleteGroupMember(req.GroupId, req.UserId)
	if err != nil {
		g.Logger.Error("delete group member error", zap.Error(err))
		resp.BaseResp = tools.BuildBaseResp(errno.GroupSrvErr)
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(errno.Success)
	return resp, nil
}

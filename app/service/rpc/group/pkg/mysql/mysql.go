package mysql

import (
	"go-ssip/app/common/errno"
	"gorm.io/gorm"
	"time"
)

type Group struct {
	GroupID      int64         `gorm:"column:group_id;primaryKey" json:"group_id"`
	GroupName    string        `gorm:"group_name" json:"group_name"`
	CreateTime   time.Time     `gorm:"create_time;autoCreateTime" `
	GroupMembers []GroupMember `gorm:"foreignKey:group_id;references:group_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GroupMember struct {
	ID         int       `gorm:"column:id;primaryKey" json:"id"`
	GroupID    int64     `gorm:"column:group_id" json:"group_id"`
	UserID     int64     `gorm:"column:user_id" json:"user_id"`
	CreateTime time.Time `gorm:"create_time;autoCreateTime"`
}

func NewGroupManager(db *gorm.DB) *GroupManager {
	return &GroupManager{
		db: db,
	}
}

type GroupManager struct {
	db *gorm.DB
}

func (m *GroupManager) CreateGroup(g *Group) error {
	err := m.db.Create(&g).Error
	return err
}

func (m *GroupManager) GetGroupByName(name string) (*Group, error) {
	var g Group
	err := m.db.Where(&Group{GroupName: name}).First(&g).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.RecordNotFound
		}
		return nil, err
	}
	return &g, nil
}

func (m *GroupManager) GetGroupById(id int64) (*Group, error) {
	var g Group
	err := m.db.Where(&Group{GroupID: id}).First(&g).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.RecordNotFound
		}
		return nil, err
	}
	return &g, nil
}

func (m *GroupManager) CreateGroupMember(gm *GroupMember) error {
	err := m.db.Create(&gm).Error
	return err
}

func (m *GroupManager) DeleteGroupMember(groupID, userID int64) error {
	err := m.db.Where(&GroupMember{GroupID: groupID, UserID: userID}).Delete(&GroupMember{}).Error
	return err
}

func (m *GroupManager) GetGroupMember(groupID, userID int64) error {
	var gm GroupMember
	err := m.db.Where(&GroupMember{GroupID: groupID, UserID: userID}).First(&gm).Error
	return err
}

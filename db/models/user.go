package models

import "github.com/google/uuid"

type Organization struct {
	ID   uuid.UUID `gorm:"primaryKey"`
	Name string
}

type User struct {
	ID             uuid.UUID    `gorm:"primaryKey"`
	OrganizationID uuid.UUID    //FK to Organization.id
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	WalletAddress  string       `gorm:"uniqueIndex"`
	Email          string       `gorm:"uniqueIndex"`
	UserGroupID    uint64       //FK to UserGroup.id
	UserGroup      LvUserGroup  `gorm:"foreignKey:UserGroupID"`
}

type LvUserGroup struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	GroupName string
}

type LvPermission struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	Permission string
}

type PermissionToUserGroup struct {
	UserGroupID  uint64       //FK to UserGroup.id
	UserGroup    LvUserGroup  `gorm:"foreignKey:UserGroupID"`
	PermissionID uint64       //FK to Permission.id
	Permission   LvPermission `gorm:"foreignKey:PermissionID"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string
}

type User struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey"`
	OrganizationID uuid.UUID    //FK to Organization.id
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	WalletAddress  *string      `gorm:"uniqueIndex"`
	Email          *string      `gorm:"uniqueIndex"`
	UserGroupID    uint64       //FK to UserGroup.id
	UserGroup      LvUserGroup  `gorm:"foreignKey:UserGroupID"`
	IsValidated    bool
	IsActive       bool
	CreatedAt      time.Time `gorm:"type:TIMESTAMPTZ;"`
	LastLoginAt    time.Time `gorm:"type:TIMESTAMPTZ;"`
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

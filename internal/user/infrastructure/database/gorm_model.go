package database

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Name      string           `gorm:"column:name;not null"`
	Document  string           `gorm:"column:document;not null"`
	BirthDate string           `gorm:"column:birth_date;not null"`
	Users     []Users          `gorm:"foreignKey:PersonId"`
	Contacts  []PersonContacts `gorm:"foreignKey:PersonId"`
}

type PersonContacts struct {
	gorm.Model
	Contact       string `gorm:"column:contact;not null"`
	Main          bool   `gorm:"column:main;not null; default:false"`
	PersonId      uint   `gorm:"column:person_id;not null"`
	Person        Person
	ContactTypeId uint `gorm:"column:person_contact_type_id;not null"`
	ContactType   PersonContactsType
}

type PersonContactsType struct {
	gorm.Model
	Name           string           `gorm:"column:name;not null"`
	PersonContacts []PersonContacts `gorm:"foreignKey:ContactTypeId"`
}

type Users struct {
	gorm.Model
	Username         string        `gorm:"column:username;not null"`
	Password         string        `gorm:"column:password;not null"`
	VerificationCode string        `gorm:"column:verification_code"`
	UserGroups       []UsersGroups `gorm:"many2many:user_group_users;joinForeignKey:user_id;joinReferences:user_group_id"`
	PersonId         uint          `gorm:"column:person_id;not null"`
	Person           Person
	Devices          []UsersDevice `gorm:"foreignKey:UserId"`
}

type UsersDevice struct {
	gorm.Model
	UserId   uint `gorm:"column:user_id;not null"`
	User     Users
	DeviceId string `gorm:"column:device_id;not null"`
}

type UserGroupsPermissions struct {
	gorm.Model
	ApiRoutApplicationId uint   `gorm:"column:api_route_application_id;not null"`
	Action               string `gorm:"column:action;not null"`
	UserGroupId          uint   `gorm:"column:user_group_id;not null"`
	UserGroup            UsersGroups
}

type UserGroupUser struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Main          bool           `gorm:"column:main;not null; default:false"`
	UsersID       uint           `gorm:"column:user_id;primaryKey"`
	UsersGroupsID uint           `gorm:"column:user_group_id;primaryKey"`
}

type UsersGroups struct {
	gorm.Model
	Name        string                  `gorm:"column:name;not null"`
	Users       []Users                 `gorm:"many2many:user_group_users;joinForeignKey:user_group_id;joinReferences:user_id"`
	Permissions []UserGroupsPermissions `gorm:"foreignKey:UserGroupId"`
}

func (Users) TableName() string {
	return "hex-api-go.users"
}

func (Person) TableName() string {
	return "hex-api-go.persons"
}

func (UserGroupUser) TableName() string {
	return "hex-api-go.user_group_users"
}

func (UsersGroups) TableName() string {
	return "hex-api-go.users_groups"
}

func (UserGroupsPermissions) TableName() string {
	return "hex-api-go.user_groups_permissions"
}

func (PersonContacts) TableName() string {
	return "hex-api-go.person_contacts"
}

func (PersonContactsType) TableName() string {
	return "hex-api-go.person_contacts_types"
}

func (UsersDevice) TableName() string {
	return "hex-api-go.users_devices"
}

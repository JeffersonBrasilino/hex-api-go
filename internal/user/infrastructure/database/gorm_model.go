package database

import (
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Uuid      string           `gorm:"column:uuid;not null"`
	Name      string           `gorm:"column:name;not null"`
	Document  string           `gorm:"column:document;not null"`
	BirthDate string           `gorm:"column:birth_date;not null"`
	Users     []Users          `gorm:"foreignKey:PersonId"`
	Contacts  []PersonContacts `gorm:"foreignKey:PersonId"`
	Status    int              `gorm:"column:status;not null; default:1"`
}

type PersonContacts struct {
	gorm.Model
	Uuid          string `gorm:"column:uuid;not null"`
	Contact       string `gorm:"column:contact;not null"`
	Main          bool   `gorm:"column:main;not null; default:false"`
	PersonId      uint   `gorm:"column:person_id;not null"`
	Person        Person
	ContactTypeId uint `gorm:"column:person_contact_type_id;not null"`
	ContactType   PersonContactsType
	Status        int `gorm:"column:status;not null; default:1"`
}

type PersonContactsType struct {
	gorm.Model
	Uuid           string           `gorm:"column:uuid;not null"`
	Name           string           `gorm:"column:name;not null"`
	PersonContacts []PersonContacts `gorm:"foreignKey:ContactTypeId"`
	Status         int              `gorm:"column:status;not null; default:1"`
}

type Users struct {
	gorm.Model
	Uuid             string `gorm:"column:uuid;not null"`
	Username         string `gorm:"column:username;not null"`
	Password         string `gorm:"column:password;not null"`
	VerificationCode string `gorm:"column:verification_code"`
	PersonId         uint   `gorm:"column:person_id;not null"`
	Person           Person
	Status           int `gorm:"column:status;not null; default:1"`
	//UserGroups       []UsersGroups `gorm:"many2many:user_group_users;joinForeignKey:user_id;joinReferences:user_group_id"`
	Devices         []UsersDevice   `gorm:"foreignKey:UserId"`
	UserGroupsUsers []UserGroupUser `gorm:"foreignKey:UserId"` //for create/update only need change custon fields.
}

type UsersGroups struct {
	gorm.Model
	Uuid string `gorm:"column:uuid;not null"`
	Name string `gorm:"column:name;not null"`
	//Users       []Users                 `gorm:"many2many:user_group_users;joinForeignKey:user_group_id;joinReferences:user_id"`
	Permissions     []UserGroupsPermissions `gorm:"foreignKey:UserGroupId"`
	Status          int                     `gorm:"column:status;not null; default:1"`
	UserGroupsUsers []UserGroupUser         `gorm:"foreignKey:UserGroupId"`
}

type UserGroupUser struct {
	gorm.Model
	Uuid        string `gorm:"column:uuid;not null"`
	Main        bool   `gorm:"column:main;not null; default:false"`
	UserId      uint   `gorm:"column:user_id"`
	UserGroupId uint   `gorm:"column:user_group_id"`
	Status      int    `gorm:"column:status;not null; default:1"`
}

type UsersDevice struct {
	gorm.Model
	Uuid     string `gorm:"column:uuid;not null"`
	UserId   uint   `gorm:"column:user_id;not null"`
	User     Users
	DeviceId string `gorm:"column:device_id;not null"`
	Status   int    `gorm:"column:status;not null; default:1"`
}

type UserGroupsPermissions struct {
	gorm.Model
	Uuid        string `gorm:"column:uuid;not null"`
	ApiRouteId  uint   `gorm:"column:api_route_id;not null"`
	ApiRoutes   ApiRoutes
	Action      string `gorm:"column:action;not null"`
	UserGroupId uint   `gorm:"column:user_group_id;not null"`
	UserGroup   UsersGroups
	Status      int `gorm:"column:status;not null; default:1"`
}

type ApiRoutes struct {
	gorm.Model
	Uuid   string `gorm:"column:uuid;not null"`
	Route  string `gorm:"column:route;not null"`
	Status int    `gorm:"column:status;not null; default:1"`
	UserGroupsPermissions []UserGroupsPermissions `gorm:"foreignKey:ApiRouteId"`
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
func (ApiRoutes) TableName() string {
	return "hex-api-go.api_routes"
}

// alteracoes
// ALTER TABLE "hex-api-go".user_groups_permissions RENAME COLUMN api_route_application_id TO api_route_id;
// ALTER TABLE "user_groups_permissions" ADD CONSTRAINT "FK_05def5b3fde2080239fd27575bd" FOREIGN KEY ("api_route_id") REFERENCES "api_routes"("id") ON DELETE NO ACTION ON UPDATE NO ACTION

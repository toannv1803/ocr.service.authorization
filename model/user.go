package model

type User struct {
	Id          string `json:"id,omitempty" bson:"id,omitempty" form:"id,omitempty"`
	Username    string `json:"username,omitempty" bson:"username,omitempty" form:"username,omitempty"`
	Password    string `json:"password,omitempty" bson:"password,omitempty" form:"password,omitempty"`
	FullName    string `json:"full_name,omitempty" bson:"full_name,omitempty" form:"full_name,omitempty"`
	Email       string `json:"email,omitempty" bson:"email,omitempty" form:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty" bson:"phone_number,omitempty" form:"phone_number,omitempty"`
	Roles       string `json:"roles,omitempty" bson:"roles,omitempty" form:"roles,omitempty"`
	Status      string `json:"status,omitempty" bson:"status,omitempty" form:"status,omitempty"`
	CreateAt    string `json:"create_at,omitempty" bson:"create_at,omitempty" form:"create_at,omitempty"`
}

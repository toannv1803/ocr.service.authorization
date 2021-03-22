package model

type UserLog struct {
	UserId      string `json:"user_id,omitempty" bson:"user_id,omitempty" form:"user_id,omitempty"`
	Ip          string `json:"ip,omitempty" bson:"ip,omitempty" form:"ip,omitempty"`
	Mac         string `json:"mac,omitempty" bson:"mac,omitempty" form:"mac,omitempty"`
	Status      string `json:"status,omitempty" bson:"status,omitempty" form:"status,omitempty"`
	CreateAt    string `json:"create_at,omitempty" bson:"create_at,omitempty" form:"create_at,omitempty"`
	ExpiredTime string `json:"expired_time,omitempty" bson:"expired_time,omitempty" form:"expired_time,omitempty"`
}

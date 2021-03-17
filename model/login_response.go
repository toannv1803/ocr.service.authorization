package model

type LoginResponse struct {
	Code   int    `json:"code,omitempty" bson:"code,omitempty" form:"code,omitempty"`
	UserId string `json:"user_id,omitempty" bson:"user_id,omitempty" form:"user_id,omitempty"`
	Expire string `json:"expire,omitempty" bson:"expire,omitempty" form:"expire,omitempty"`
	Token  string `json:"token,omitempty" bson:"token,omitempty" form:"token,omitempty"`
}

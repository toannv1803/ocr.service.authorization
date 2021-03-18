package model

type Claim struct {
	UserId string `json:"user_id,omitempty" bson:"user_id,omitempty" form:"user_id,omitempty"`
	Role   string `json:"role,omitempty" bson:"role,omitempty" form:"role,omitempty"`
}
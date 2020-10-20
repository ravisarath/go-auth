package Models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Username    string    `json:"username"`
	Company     string    `json:"company"`
	Group       string    `json:"group"`
	Admin       bool      `json:"admin"`
	CreatedDate time.Time `json:"createdDate"`
	CreatedBy   string    `json:"createdby"`
}
type Accountuserdetails struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Username string             `json:"username"`
	Company  string             `json:"company"`
	Group    string             `json:"group"`
	Admin    bool               `json:"admin"`
}

type AccountAccesstoken struct {
	UserID       string `json:"userid"`
	Accesstoken  string `json:"accesstoken"`
	Refreshtoken string `json:"refreshtoken"`
}

type LoginTime struct {
	UserID string    `json:"userid"`
	Login  time.Time `json:"login"`
	Logout time.Time `json:"logout"`
}

type UserLoginTime struct {
	UserID string    `json:"userid"`
	Login  time.Time `json:"login"`
	Logout time.Time `json:"logout"`
}

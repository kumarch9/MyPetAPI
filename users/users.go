package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty" validate:"required"`
	Email        string             `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
	MobileNum    string             `bson:"mobilenum,omitempty" json:"mobilenum,omitempty" validate:"required,len=10"`
	Organisation string             `bson:"organisation,omitempty" json:"organisation,omitempty" validate:"required"`
	Designation  string             `bson:"designation,omitempty" json:"designation,omitempty" validate:"required"`
	Address      string             `bson:"address,omitempty" json:"address,omitempty" validate:"required"`
	AadhaarNum   string             `bson:"aadhaarnum,omitempty" json:"aadhaarnum,omitempty" validate:"required"`
	Password     string             `bson:"password" json:"password" validate:"required"`
	CreatedAt    time.Time          `bson:"createdAt," json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    time.Time          `bson:"deletedAt" json:"deletedAt"`
}

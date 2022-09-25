package pets

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// use it when admin get request of pet data
type Pet struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	GenInfo      PetGeneralInfo     `bson:"geninfo" json:"geninfo" validate:"required" `
	ActivityInfo PetActivities      `bson:"activityinfo" json:"activityinfo" validate:"required" `
	FeedInfo     PetFeed            `bson:"feedinfo" json:"feedinfo" validate:"required" `
	VetInfo      PetVeterinary      `bson:"vetinfo" json:"vetinfo" validate:"required" `
	CreatorInfo  Creator            `bson:"creatorinfo,omitempty" json:"creatorinfo,omitempty"`
	UpdaterInfo  Updater            `bson:"updaterinfo,omitempty" json:"updaterinfo,omitempty"`
	DeleterInfo  Deleter            `bson:"deleterinfo,omitempty" json:"deleterinfo,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    time.Time          `bson:"deletedAt" json:"deletedAt"`
}

// use it when user get request of pet data
type PetBind_User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	GenInfo      PetGeneralInfo     `bson:"geninfo" json:"geninfo" validate:"required" `
	ActivityInfo PetActivities      `bson:"activityinfo" json:"activityinfo" validate:"required" `
	FeedInfo     PetFeed            `bson:"feedinfo" json:"feedinfo" validate:"required" `
	VetInfo      PetVeterinary      `bson:"vetinfo" json:"vetinfo" validate:"required" `
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    time.Time          `bson:"deletedAt" json:"deletedAt"`
}

type PetGeneralInfo struct {
	Name                  string `bson:"name,omitempty" json:"name,omitempty" validate:"required" `
	PetType               string `bson:"pettype,omitempty" json:"pettype,omitempty" validate:"required" `
	Age                   int    `bson:"age,omitempty" json:"age,omitempty" validate:"required" `
	Gender                string `bson:"gender,omitempty" json:"gender,omitempty" validate:"required" `
	DOB                   string `bson:"DOB,omitempty" json:"DOB,omitempty" validate:"required" `
	Weight                string `bson:"weight,omitempty" json:"weight,omitempty" validate:"required" `
	Breed                 string `bson:"breed,omitempty" json:"breed,omitempty" validate:"required" `
	Card_Number           string `bson:"card_number,omitempty" json:"card_number,omitempty" validate:"required" `
	Rabies_Vaccine_Number string `bson:"rabies_vaccine_number,omitempty" json:"rabies_vaccine_number,omitempty" validate:"required" `
	Licence_Number        string `bson:"licence_number,omitempty" json:"licence_number,omitempty" validate:"required" `
	Any_Allergies         string `bson:"any_allergies,omitempty" json:"any_allergies,omitempty" validate:"required" `
	Contact_Number        string `bson:"contact_number,omitempty" json:"contact_number,omitempty" validate:"required" `
	Address               string `bson:"address,omitempty" json:"address,omitempty" validate:"required" `
}

type PetActivities struct {
	Like          string `bson:"like,omitempty" json:"like,omitempty" validate:"required" `
	Dislike       string `bson:"dislike,omitempty" json:"dislike,omitempty" validate:"required"`
	Place_To_Play string `bson:"place_to_play,omitempty" json:"place_to_play,omitempty" validate:"required" `
}

type PetFeed struct {
	Brand_Name     string `bson:"brand_name,omitempty" json:"brand_name,omitempty" validate:"required" `
	Non_Brand_Name string `bson:"non_brand_name,omitempty" json:"non_brand_name,omitempty" validate:"required" `
	Morning_Amount string `bson:"morning_amount,omitempty" json:"morning_amount,omitempty" validate:"required" `
	Noon_Amount    string `bson:"noon_amount,omitempty" json:"noon_amount,omitempty" validate:"required" `
	Night_Amount   string `bson:"night_amount,omitempty" json:"night_amount,omitempty" validate:"required" `
}

type PetVeterinary struct {
	Regular_Vet         string `bson:"regular_vet,omitempty" json:"regular_vet,omitempty" validate:"required" `
	Regular_Vet_Contact string `bson:"regular_vet_contact,omitempty" json:"regular_vet_contact,omitempty" validate:"required" `
}

type Creator struct {
	CreatorName  string `bson:"creatorname,omitempty" json:"creatorname,omitempty"`
	CreatorEmail string `bson:"creatoremail,omitempty" json:"creatoremail,omitempty"`
}

type Updater struct {
	UpdaterName  string `bson:"updatername,omitempty" json:"UpdaterName,omitempty"`
	UpdaterEmail string `bson:"updateremail,omitempty" json:"updateremail,omitempty"`
}

type Deleter struct {
	DeleterName  string `bson:"deletername,omitempty" json:"deletername,omitempty"`
	DeleterEmail string `bson:"deleteremail,omitempty" json:"deleteremail,omitempty"`
}

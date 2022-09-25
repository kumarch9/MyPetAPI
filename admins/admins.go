package admins

type Admin struct {
	AdminName string `bson:"adminname,omitempty" json:"adminname,omitempty"`
	Email     string `bson:"email,omitempty" json:"email,omitempty" validate:"email"`
	Password  string `bson:"password" json:"password"`
}

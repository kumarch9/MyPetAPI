package logics

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func KeyValueIsExists(collectionName *mongo.Collection, inputKey string, inputValue any) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	var keysModel = bson.M{}
	err := collectionName.FindOne(ctx, bson.M{inputKey: inputValue}).Decode(&keysModel)
	if err != nil {
		//log.Println("err in  CollectionValueIsExists :", err)
		return false
	}
	if keysModel != nil {
		//log.Println("keymodel ! nil :", keysModel)
		return true
	}
	//log.Println("keymodel  nil :", keysModel)
	return false
}

// func FindValueInDB(UserIForBind interface{},
// 	dbCollection *mongo.Collection, findKey string, findValue string) (findOk bool, er error) {

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	errFind := dbCollection.FindOne(ctx, bson.M{findKey: findValue}).Decode(&UserIForBind)
// 	if errFind != nil {
// 		return false, errFind
// 	}
// 	return true, nil
// }

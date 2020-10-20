package Models

import (
	"context"
	"fmt"
	"jwt-todo/auth-server/Config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Timestamp time.Time

func UpdateLoginTime(userId string) (err error) {
	collection := Config.CLI.Database("test").Collection("LoginTime")
	ctx := context.Background()
	count, err := collection.CountDocuments(ctx, bson.M{"UserID": userId})
	fmt.Println(count)
	if count == 0 {

		logdetails := map[string]interface{}{
			"UserID": userId,
			"Login":  time.Now(),
		}
		result, err := collection.InsertOne(ctx, logdetails)
		if err != nil {

			return err
		}
		objectID := result.InsertedID.(primitive.ObjectID)
		fmt.Println(objectID)

		return nil

	} else {
		resultUpdate, err := collection.UpdateOne(
			ctx,
			bson.M{"UserID": userId},
			bson.M{
				"$set": bson.M{

					"Login": time.Now(),
				},
			},
		)
		if err != nil {

			return err
		}
		fmt.Println(resultUpdate.ModifiedCount)
		return nil

	}

}
func UpdateLogoutTime(userId string) (err error) {
	collection := Config.CLI.Database("test").Collection("LoginTime")
	ctx := context.Background()

	resultUpdate, err := collection.UpdateOne(
		ctx,
		bson.M{"UserID": userId},
		bson.M{
			"$set": bson.M{
				"Logout": time.Now(),
			},
		},
	)
	if err != nil {

		return err
	}

	findResult := collection.FindOne(ctx, bson.M{"UserID": userId})

	if err := findResult.Err(); err != nil {

		return err
	}
	account := &LoginTime{}
	err = findResult.Decode(&account)
	if err != nil {

		return err

	}
	err = UpdateLogdetailsTime(account)
	if err != nil {
		return err
	}
	fmt.Println(resultUpdate.ModifiedCount)
	return nil

}
func UpdateLogdetailsTime(login *LoginTime) (err error) {
	collection := Config.CLI.Database("test").Collection("userLoginTime")
	ctx := context.Background()

	logdetails := map[string]interface{}{
		"UserID": login.UserID,
		"Login":  login.Login,
		"Logout": login.Logout,
	}
	result, err := collection.InsertOne(ctx, logdetails)
	if err != nil {

		return err
	}
	objectID := result.InsertedID.(primitive.ObjectID)
	fmt.Println(objectID)

	return nil
}

package Models

import (
	"context"
	"fmt"
	"jwt-todo/auth-server/Config"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (account *Account) Validate() bool {

	if !strings.Contains(account.Email, "@") {
		return false
	}

	if len(account.Password) < 6 {
		return false
	}

	temp := &Account{}
	collection := Config.CLI.Database("test").Collection("notesCollection")
	ctx := context.Background()
	err := collection.FindOne(ctx, bson.M{"email": account.Email}).Decode(&temp)

	fmt.Println(err)

	if temp.Email != "" {

		return false
	}

	return true
}
func CreateAccount(account *Account) bool {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	resp := account.Validate()
	if !resp {

		return false
	}

	Collection := Config.CLI.Database("test").Collection("notesCollection")

	ctx := context.Background()
	result, err := Collection.InsertOne(ctx, account)
	if err != nil {

		return false
	}

	objectID := result.InsertedID.(primitive.ObjectID)
	fmt.Println(objectID)

	return true
}

func Login(email string, password string) (bool, *Accountuserdetails) {

	account := &Accountuserdetails{}
	collection := Config.CLI.Database("test").Collection("notesCollection")
	ctx := context.Background()
	findResult := collection.FindOne(ctx, bson.M{"email": email})

	if err := findResult.Err(); err != nil {
		fmt.Println(err)
		return false, account
	}

	err := findResult.Decode(&account)
	if err != nil {
		fmt.Println(err)
		return false, account

	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		fmt.Println("password missmatch")
		return false, account
	}
	//Worked! Logged In
	account.Password = ""

	return true, account
}

func GetDetails(account *Account, email string) (err error) {

	collection := Config.CLI.Database("test").Collection("notesCollection")
	ctx := context.Background()
	findResult := collection.FindOne(ctx, bson.M{"email": email})

	if err := findResult.Err(); err != nil {

		return err
	}

	err = findResult.Decode(&account)
	if err != nil {

		return err

	}
	return nil
}
func UpdateToken(userId string, accestoken string, refreshtoken string) (err error) {
	collection := Config.CLI.Database("test").Collection("userAccessToken")
	ctx := context.Background()
	count, err := collection.CountDocuments(ctx, bson.M{"uniqueid": userId})
	fmt.Println(userId)
	fmt.Println(count)
	if count == 0 {

		fmt.Println("count")
		accesstoken := map[string]string{
			"uniqueid":     userId,
			"accesstoken":  accestoken,
			"refreshtoken": refreshtoken,
		}
		result, err := collection.InsertOne(ctx, accesstoken)
		if err != nil {

			return err
		}
		objectID := result.InsertedID.(primitive.ObjectID)
		fmt.Println(objectID)

		return nil

	} else {
		resultUpdate, err := collection.UpdateOne(
			ctx,
			bson.M{"uniqueid": userId},
			bson.M{
				"$set": bson.M{

					"accesstoken":  accestoken,
					"refreshtoken": refreshtoken,
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
func RefreshUpdateToken(userId string, AccessToken string, RefreshToken string) (err error) {
	collection := Config.CLI.Database("test").Collection("userAccessToken")
	ctx := context.Background()
	fmt.Println(userId)
	resultUpdate, err := collection.UpdateOne(
		ctx,
		bson.M{"uniqueid": userId},
		bson.M{
			"$set": bson.M{

				"accesstoken":  AccessToken,
				"refreshtoken": RefreshToken,
			},
		},
	)
	if err != nil {

		return err
	}
	fmt.Println(resultUpdate.ModifiedCount)
	return nil
}
func RemoveToken(userId string) (err error) {
	collection := Config.CLI.Database("test").Collection("userAccessToken")
	ctx := context.Background()
	resultUpdate, err := collection.UpdateOne(
		ctx,
		bson.M{"uniqueid": userId},
		bson.M{
			"$set": bson.M{

				"accesstoken":  "",
				"refreshtoken": "",
			},
		},
	)
	if err != nil {

		return err
	}
	fmt.Println(resultUpdate.ModifiedCount)
	return nil
}

package database

import (
	"context"
	"github.com/lotteryjs/ten-minutes-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

// GetUsers returns all users.
// start, end int, order, sort string
func (d *TenDatabase) GetUsers(paging *model.Paging) []*model.User {
	var users []*model.User
	cursor, err := d.DB.Collection("users").
		Find(context.Background(), bson.D{},
			&options.FindOptions{
				Skip:  paging.Skip,
				Sort:  bson.D{bson.E{Key: paging.SortKey, Value: paging.SortVal}},
				Limit: paging.Limit,
			})
	if err != nil {
		return nil
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		user := &model.User{}
		if err := cursor.Decode(user); err != nil {
			return nil
		}
		users = append(users, user)
	}

	return users
}

// CreateUser creates a user.
func (d *TenDatabase) CreateUser(user *model.User) error {
	if _, err := d.DB.Collection("users").
		InsertOne(context.Background(), user); err != nil {
		return err
	}
	return nil
}

// GetUserByName returns the user by the given name or nil.
func (d *TenDatabase) GetUserByName(name string) *model.User {
	var user *model.User
	err := d.DB.Collection("users").
		FindOne(nil, bson.D{{Key: "name", Value: name}}).
		Decode(&user)
	if err != nil {
		return nil
	}
	return user
}

// GetUserByID returns the user by the given id or nil.
func (d *TenDatabase) GetUserByID(id primitive.ObjectID) *model.User {
	var user *model.User
	err := d.DB.Collection("users").
		FindOne(nil, bson.D{{Key: "_id", Value: id}}).
		Decode(&user)
	if err != nil {
		return nil
	}
	return user
}

// CountUser returns the user count
func (d *TenDatabase) CountUser() string {
	total, err := d.DB.Collection("users").CountDocuments(nil, bson.D{{}}, nil)
	if err != nil {
		return "0"
	}
	return strconv.Itoa(int(total))
}

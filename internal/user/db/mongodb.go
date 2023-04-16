package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nikitapopov/Habbit/internal/apperror"
	userPkg "github.com/Nikitapopov/Habbit/internal/user"
	"github.com/Nikitapopov/Habbit/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewRepository(collection *mongo.Collection, logger *logging.Logger) userPkg.Repository {
	return &db{
		collection: collection,
		logger:     logger,
	}
}

func (d *db) Create(ctx context.Context, user *userPkg.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, &userPkg.User{
		ID:           primitive.NewObjectID(),
		Username:     user.Username,
		PasswordHash: fmt.Sprintf("%s-hashed", user.PasswordHash),
		Email:        user.Email,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}

	d.logger.Debug("convert InsertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to hex. probably oid: %s", oid)
}

func (d *db) FindAll(ctx context.Context) (u []userPkg.User, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return u, fmt.Errorf("failed to find users due to error: %v", err)
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read users due to error: %v", err)
	}

	return u, nil
}

func (d *db) FindOne(ctx context.Context, id string) (u userPkg.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}

	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrNotFound
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, result.Err())
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user (id:%s) from DB due to error: %v", id, err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, userId string, user userPkg.User) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectid. hex: %s", user.ID)
	}

	filter := bson.M{"_id": oid}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user. error: %v", err)
	}

	var updateUserObj bson.M
	if err = bson.Unmarshal(userBytes, &updateUserObj); err != nil {
		return fmt.Errorf("failed to unmarshal user bytes. error: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{"$set": updateUserObj}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query. error: %v", err)
	}

	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}
	d.logger.Tracef("Matched %d documents and Modified &d documents", result.MatchedCount, result.ModifiedCount)
	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}

	filter := bson.M{"_id": oid}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}
	d.logger.Tracef("Deleted %d documents", result.DeletedCount)
	return nil
}

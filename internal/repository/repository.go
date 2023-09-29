package repository

import (
	"errors"

	"github.com/demig00d/auth-service/internal/model"
	. "github.com/demig00d/auth-service/internal/model"
	"github.com/demig00d/auth-service/pkg/logger"
	"github.com/demig00d/auth-service/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrCantUpsertRefreshToken = errors.New("can't upsert refresh token")
	ErrCantFindRefreshToken   = errors.New("can't find refresh token")
)

type Repository interface {
	UpsertRefreshToken(GUID, Encrypted[RefreshToken]) error
	FindRefreshToken(GUID) (Encrypted[RefreshToken], error)
}

type RepositoryImpl struct {
	*mongodb.MongoDB
	*mongo.Collection
	logger logger.Logger
}

func NewRepository(
	m *mongodb.MongoDB,
	dbName,
	collectionName string,
	logger logger.Logger,
) RepositoryImpl {
	collection := m.Database(dbName).Collection(collectionName)
	logger.SetPrefix("repository ")

	return RepositoryImpl{m, collection, logger}
}

func (r RepositoryImpl) UpsertRefreshToken(guid GUID, refreshToken Encrypted[RefreshToken]) error {
	doUpsert := true

	res, err := r.UpdateOne(
		r.Ctx,
		bson.M{"_id": guid.String()},
		bson.M{"$set": bson.M{"refresh_token": refreshToken.String()}},
		&options.UpdateOptions{Upsert: &doUpsert},
	)

	if err != nil {
		r.logger.Debug(err)
		return err
	}

	if res.ModifiedCount != 1 && res.UpsertedCount != 1 {
		return errors.New("storage: token has not been saved to the database")
	}

	return nil
}

func (r RepositoryImpl) FindRefreshToken(guid GUID) (Encrypted[RefreshToken], error) {
	var result struct {
		GUID                  string `bson:"_id"`
		EncryptedRefreshToken string `bson:"refresh_token"`
	}
	err := r.FindOne(
		r.Ctx,
		bson.M{"_id": guid.String()},
	).Decode(&result)

	if err != nil {
		r.logger.Debug(err)
		return Encrypted[RefreshToken]{}, ErrCantFindRefreshToken
	}

	encrypted, err := model.EncryptedFromString[model.RefreshToken](
		result.EncryptedRefreshToken,
	)

	if err != nil {
		r.logger.Debug(err)
		return Encrypted[RefreshToken]{}, ErrCantFindRefreshToken
	}

	return encrypted, nil
}

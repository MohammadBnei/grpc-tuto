package person

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Server struct {
	Db *mongo.Database
}

func (s *Server) CreateUser(ctx context.Context, person *Person) (*Id, error) {
	inserted, err := s.Db.Collection("person").InsertOne(ctx, person)
	if err != nil {
		log.Fatalf("Error inserting person: %v", err)
		return nil, err
	}
	id := &Id{
		Id: fmt.Sprint(inserted.InsertedID),
	}
	return id, nil
}

func (s *Server) GetUsers(ctx context.Context, _ *Empty) (*AddressBook, error) {
	cursor, err := s.Db.Collection("person").Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Error finding all persons: %v", err)
		return nil, err
	}

	addressBook := &AddressBook{}
	if err := cursor.All(ctx, &addressBook.People); err != nil {
		log.Fatalf("Error decoding persons: %v", err)
		return nil, err
	}

	return addressBook, nil
}

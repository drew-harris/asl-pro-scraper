package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func SaveTag(tag string) error {
	db := GetDatabase()
	collection := db.Collection("tags")
	doc := bson.M{"tag": tag}
	_, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}
	return nil
}

func SaveWord(word string, tag string) error {
	db := GetDatabase()
	collection := db.Collection("cards")
	doc := bson.M{"word": word, "filename": word + ".mp4", "tag": tag, "created": time.Now().UnixMilli()}
	_, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}
	return nil

}

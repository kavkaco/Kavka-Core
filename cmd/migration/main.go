package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type messageDoc struct {
	ChatID   model.ChatID              `bson:"chat_id"`
	Messages []*model.MessageGetter    `bson:"messages"`
}

func main() {
	fmt.Println("Kavka Message Migration Tool")
	fmt.Println("Migrates from embedded messages array to separate documents")
	fmt.Println()

	cfg := config.Read()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	uri := database.NewMongoDBConnectionString(
		cfg.Mongo.Host,
		cfg.Mongo.Port,
		cfg.Mongo.Username,
		cfg.Mongo.Password,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(cfg.Mongo.DBName)
	oldCollection := db.Collection("messages")
	newCollection := db.Collection("messages_v2")

	_, err = newCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "chat_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "chat_id", Value: 1}, {Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "chat_id", Value: 1}}},
	})
	if err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	totalMigrated := int64(0)
	totalSkipped := int64(0)

	cursor, err := oldCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to read old messages: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc messageDoc
		err := cursor.Decode(&doc)
		if err != nil {
			log.Printf("Error decoding document: %v", err)
			continue
		}

		existing, err := newCollection.CountDocuments(ctx, bson.M{"chat_id": doc.ChatID})
		if err == nil && existing > 0 {
			totalSkipped++
			fmt.Printf("Skipping chat %s (already migrated, %d messages)\n", doc.ChatID.Hex(), existing)
			continue
		}

		if len(doc.Messages) == 0 {
			totalSkipped++
			continue
		}

		var newDocs []interface{}
		for _, msgGetter := range doc.Messages {
			if msgGetter == nil || msgGetter.Message == nil {
				continue
			}

			msg := msgGetter.Message
			sender := msgGetter.Sender

			newDoc := model.MessageDocumentV2{
				ID:             msg.MessageID,
				ChatID:         doc.ChatID,
				SenderID:       msg.SenderID,
				CreatedAt:      msg.CreatedAt,
				Edited:         msg.Edited,
				Seen:           msg.Seen,
				Type:           msg.Type,
				Content:        msg.Content,
			}

			if sender != nil {
				newDoc.SenderName = sender.Name
				newDoc.SenderLastName = sender.LastName
				newDoc.SenderUsername = sender.Username
			}

			newDocs = append(newDocs, newDoc)
		}

		if len(newDocs) > 0 {
			_, err := newCollection.InsertMany(ctx, newDocs)
			if err != nil {
				log.Printf("Error inserting messages for chat %s: %v", doc.ChatID.Hex(), err)
				continue
			}

			totalMigrated += int64(len(newDocs))
			fmt.Printf("Migrated chat %s: %d messages\n", doc.ChatID.Hex(), len(newDocs))
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
	}

	fmt.Println()
	fmt.Println("Migration complete!")
	fmt.Printf("Total messages migrated: %d\n", totalMigrated)
	fmt.Printf("Total chats skipped: %d\n", totalSkipped)
}

package main

import (
	"Nexus/app/database"
	"Nexus/app/models"
	"Nexus/internal/configs"
	"Nexus/pkg/auth"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMakeVerificCode(t *testing.T) {
	t.Logf("Random Verific Code : %d\n", auth.MakeVerificCode())
}

func TestEmailWithoutAt(t *testing.T) {
	t.Logf("Email Without At : %s\n", auth.GetEmailWithoutAt("NexusChat@gmail.com"))
}

func TestIsVerificCodeExpired(t *testing.T) {
	wd, _ := os.Getwd()
	mongoConfigs, mongoConfigsErr := configs.ParseMongoConfigs(wd + "/../configs/mongo.yml")
	if mongoConfigsErr != nil {
		fmt.Println("Error in parsing mongodb configs")
	}

	email := "NexusChat@gmail.com"
	database.EstablishMongoDBConnection(*mongoConfigs)
	time.Sleep(1 * time.Second)

	var user *models.User = models.FindUserByEmail(email)
	if user == nil {
		t.Fatalf("User %s not found", email)
		return
	}

	now := time.Now().Unix()
	expire := user.VerificLimitDate
	for i := 0; i < 500; i++ {
		expired := !(now < expire)
		t.Logf("Expired? %t", expired)
		time.Sleep(1 * time.Second)
	}
}

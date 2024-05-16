package chat

import (
	"github.com/kavkaco/Kavka-Core/internal/model/message"
	"github.com/kavkaco/Kavka-Core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TypeChannel = "channel"
	TypeGroup   = "group"
	TypeDirect  = "direct"
)

type Chat struct {
	ChatID     primitive.ObjectID `bson:"chat_id"     json:"chatId"`
	ChatType   string             `bson:"chat_type"   json:"chatType"`
	ChatDetail interface{}        `bson:"chat_detail" json:"chatDetail"`
}

// Chat struct that includes messages.
// ChatC (Complete Chat) created because `Chat` struct does not contain `Messages` field.
type ChatC struct {
	ChatID     primitive.ObjectID `bson:"chat_id"     json:"chatId"`
	ChatType   string             `bson:"chat_type"   json:"chatType"`
	ChatDetail interface{}        `bson:"chat_detail" json:"chatDetail"`
	Messages   []*message.Message `bson:"messages"    json:"messages"`
}

type Member struct {
	StaticID primitive.ObjectID `bson:"id"        json:"id"`
	Name     string             `bson:"name" json:"name"`
	LastName string             `bson:"last_name"  json:"lastName"`

	// We will add a profile photo here later =)
}

// Chat Detail

type ChannelChatDetail struct {
	Title        string               `bson:"title" json:"title"`
	Members      []primitive.ObjectID `bson:"members,omitempty" json:"members"`
	Admins       []primitive.ObjectID `bson:"admins,omitempty" json:"admins"`
	Owner        *primitive.ObjectID  `bson:"owner,omitempty"         json:"owner"`
	RemovedUsers []primitive.ObjectID `bson:"removed_users,omitempty" json:"removedUsers"`
	Username     string               `bson:"username,omitempty" json:"username"`
	Description  string               `bson:"description,omitempty" json:"description"`
}

type GroupChatDetail struct {
	Title        string               `bson:"title" json:"title"`
	Members      []primitive.ObjectID `bson:"members,omitempty" json:"members"`
	Admins       []primitive.ObjectID `bson:"admins,omitempty" json:"admins"`
	Owner        *primitive.ObjectID  `bson:"owner,omitempty"         json:"owner"`
	RemovedUsers []primitive.ObjectID `bson:"removed_users,omitempty" json:"removedUsers"`
	Username     string               `bson:"username,omitempty" json:"username"`
	Description  string               `bson:"description,omitempty" json:"description"`
}

type DirectChatDetail struct {
	// Chat partners
	Sides [2]primitive.ObjectID `json:"sides"`
}

func (c *Chat) IsMember(staticID primitive.ObjectID) bool {
	d, _ := utils.TypeConverter[ChannelChatDetail](c.ChatDetail)
	for _, member := range d.Members {
		if member.Hex() == staticID.Hex() {
			return true
		}
	}

	return false
}

func (c *Chat) IsAdmin(staticID primitive.ObjectID) bool {
	d, _ := utils.TypeConverter[ChannelChatDetail](c.ChatDetail)
	for _, admin := range d.Admins {
		if admin.Hex() == staticID.Hex() {
			return true
		}
	}

	return false
}

func (d *DirectChatDetail) HasSide(staticID primitive.ObjectID) bool {
	has := false
	for _, v := range d.Sides {
		if v == staticID {
			has = true
			break
		}
	}
	return has
}

// FIXME
// DetectTarget determines the appropriate chat partner for the user identified by userStaticID,
// considering a list of potential users and assuming only two participants are involved.
// It returns a pointer to the target user's struct.
func DetectTarget(userIDs [2]primitive.ObjectID, userStaticID primitive.ObjectID) *primitive.ObjectID {
	if userIDs[0].Hex() == userStaticID.Hex() {
		return &userIDs[1]
	} else {
		return &userIDs[0]
	}
}

func NewChat(chatType string, chatDetail interface{}) *Chat {
	m := &Chat{}

	m.ChatType = chatType
	m.ChatDetail = chatDetail
	m.ChatID = primitive.NewObjectID()

	return m
}

//  Interfaces

type Repository interface {
	GetChatMembers(chatID primitive.ObjectID) []Member
	Create(newChat Chat) (*Chat, error)
	FindMany(filter bson.M) ([]Chat, error)
	FindOne(filter bson.M) (*Chat, error)
	Destroy(chatID primitive.ObjectID) error
	FindByID(staticID primitive.ObjectID) (*Chat, error)
	FindChatOrSidesByStaticID(staticID primitive.ObjectID) (*ChatC, error)
	FindBySides(sides [2]primitive.ObjectID) (*Chat, error)
}

type Service interface {
	GetChat(staticID primitive.ObjectID) (*ChatC, error)
	GetUserChats(userStaticID primitive.ObjectID) ([]ChatC, error)
	CreateDirect(userStaticID primitive.ObjectID, targetStaticID primitive.ObjectID) (*Chat, error)
	CreateGroup(userStaticID primitive.ObjectID, title string, username string, description string) (*Chat, error)
	CreateChannel(userStaticID primitive.ObjectID, title string, username string, description string) (*Chat, error)
}

package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/trevorprater/youfie-api-2/core/etc"
)

type Conversation struct {
	ID      string `json:"id" form:"id" db:"id"`
	OwnerID string `json:"owner_id" form:"owner_id" db:"owner_id"`
	PhotoID string `json:"photo_id" form:"photo_id" db:"photo_id"`

	CreatedAt time.Time `json:"created_at" form:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" db:"updated_at"`

	FaceIDs      []string                   `json:"face_ids" form:"face_ids" db:"-"`
	Participants []*ConversationParticipant `json:"participants" form:"participants" db:"-"`
	Messages     []*ConversationMessage     `json:"messages" form:"messages" db:"-"`
	MessageText  string                     `json:"message_text" form:"message_text" db:"-"`
}

type ConversationParticipant struct {
	ConversationID string `json:"-" form:"-" db:"conversation_id"`
	FaceID         string `json:"-" form:"-" db:"face_id"`
	UserApproved   bool   `json:"-" form:"-" db:"user_approved"`

	CreatedAt time.Time `json"created_at" form:"created_at" db:"created_at"`
}

type ConversationMessage struct {
	ID             string `json:"id" form:"id" db:"id"`
	ConversationID string `json:"-" form:"-" db:"conversation_id"`
	OwnerID        string `json:"-" form:"-" db:"owner_id"`
	MessageText    string `json:"message_text" form:"message_text" db:"message_text"`
}

func (m *ConversationMessage) Insert(db sqlx.Ext) error {
	m.ID = uuid.New()
	q := `INSERT INTO conversation_messages(id, conversation_id, face_id, message_text) VALUES (
		    :id,
			:conversation_id,
			:owner_id,
			:message_text`

	_, err := sqlx.NamedExec(db, q, m)
	if etc.Duperr(err) {
		log.Println(err)
		return err
	}
	if err != nil {
		log.Println(err)
	}
	return err
}

func getMessage(messageID string, db sqlx.Ext) (*ConversationMessage, error) {
	var message *ConversationMessage
	err := sqlx.Get(db, &message, "SELECT * FROM conversation_messages WHERE id = '"+messageID+"'")
	return message, err
}

func (c *Conversation) populate(db sqlx.Ext) error {
	err := c.populateParticipants(db)
	if err != nil {
		log.Println(err)
		return err
	}
	err = c.populateMessages(db)
	if err != nil {
		log.Println(err)
		return err
	}
	return c.populateOwnerID(db)
}

func (c *Conversation) populateMessages(db sqlx.Ext) error {
	var messages []*ConversationMessage
	rows, err := db.Queryx("SELECT * FROM conversation_messages WHERE conversation_id='" + c.ID + "'")
	if err != nil {
		log.Println(err)
		return err
	}
	for rows.Next() {
		var m ConversationMessage
		err = rows.StructScan(&m)
		if err != nil {
			log.Println(err)
		}
		messages = append(messages, &m)
	}
	c.Messages = messages
	return err
}

func (c *Conversation) populateParticipants(db sqlx.Ext) error {
	var participants []*ConversationParticipant
	rows, err := db.Queryx("SELECT * FROM conversation_participants WHERE conversation_id='" + c.ID + "'")
	if err != nil {
		log.Println(err)
		return err
	}
	for rows.Next() {
		var p ConversationParticipant
		err = rows.StructScan(&p)
		if err != nil {
			log.Println(err)
		}
		participants = append(participants, &p)
	}
	c.Participants = participants
	return err
}

func (c *Conversation) populateOwnerID(db sqlx.Ext) error {
	photo, err := GetPhotoByID(c.PhotoID, db)
	if err != nil {
		log.Println(err)
	}
	c.OwnerID = photo.ID
	return err
}

func GetConversationsForUser(userID string, db sqlx.Ext) ([]*Conversation, error) {
	// Get all conversations created by the requesting user.
	var conversations []*Conversation
	rows, err := db.Queryx("SELECT * FROM conversations WHERE owner_id='" + userID + "'")
	if err != nil {
		log.Println(err)
		return conversations, err
	}
	for rows.Next() {
		var c Conversation
		err = rows.StructScan(&c)
		if err != nil {
			log.Println(err)
		}

		err = c.populateParticipants(db)
		if err != nil {
			log.Println(err)
		}
		conversations = append(conversations, &c)
	}

	// Get all conversations the user's face appears in
	rows, err = db.Queryx("SELECT * FROM matches WHERE user_id = '" + userID + "' AND is_match=true")
	if err != nil {
		log.Println(err)
		return conversations, err
	}
	for rows.Next() {
		var m Match
		err = rows.StructScan(&m)
		if err != nil {
			log.Println(err)
			return conversations, err
		}
		rows, err = db.Queryx("SELECT * FROM conversation_participants WHERE face_id='" + m.FaceID + "'")
		if err != nil {
			log.Println(err)
			return conversations, err
		}
		for rows.Next() {
			var p ConversationParticipant
			err = rows.StructScan(&p)
			if err != nil {
				log.Println(err)
				return conversations, err
			}
			c, err := GetConversationByID(p.ConversationID, db)
			if err != nil {
				log.Println(err)
				return conversations, err
			}
			err = c.populateParticipants(db)
			if err != nil {
				log.Println(err)
				return conversations, err
			}
			conversations = append(conversations, c)
		}
	}

	return conversations, err
}

func GetConversationByID(id string, db sqlx.Ext) (*Conversation, error) {
	var conversation Conversation
	err := sqlx.Get(db, &conversation, "SELECT * FROM conversations WHERE id = '"+id+"'")
	if err != nil {
		log.Println(err)
		return &conversation, err
	}
	err = conversation.populateParticipants(db)
	if err != nil {
		log.Println(err)
		return &conversation, err
	}
	err = conversation.populateMessages(db)
	if err != nil {
		log.Println(err)
	}
	return &conversation, err
}

func (c *Conversation) Update(db sqlx.Ext, messageText string) ([]byte, int) {
	var m *ConversationMessage
	m.ConversationID = c.ID
	m.OwnerID = c.OwnerID
	m.MessageText = messageText

	err := m.Insert(db)
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}

	_, err = sqlx.NamedExec(db, `
		UPDATE conversations SET updated_at = CURRENT_TIMESTAMP WHERE id = :id`, c.ID)
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}

	conversation, err := GetConversationByID(c.ID, db)
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	err = conversation.populate(db)
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	conversationJson, err := json.MarshalIndent(&conversation, "", "    ")
	return conversationJson, http.StatusCreated
}

func (c *Conversation) Insert(db sqlx.Ext) ([]byte, int) {
	c.ID = uuid.New()
	err := c.populateOwnerID(db)
	if err != nil {
		return []byte("internal server error"), http.StatusInternalServerError
	}
	_, err = sqlx.NamedExec(db, `
		INSERT INTO conversations
		(id, photo_id, owner_id)
		VALUES (:id, :photo_id, :owner_id)`, c)

	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	faces, err := GetFacesForPhoto(c.PhotoID, db)
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}

	for _, face_id_in_request := range c.FaceIDs {
		for _, face_in_photo := range faces {
			if face_id_in_request == face_in_photo.ID {
				var p ConversationParticipant
				p.UserApproved = false
				p.ConversationID = c.ID
				p.FaceID = face_in_photo.ID

				_, err := sqlx.NamedExec(db, `INSERT INTO conversation_participants(conversation_id, face_id, user_approved)
				VALUES (:conversation_id, :face_id, :user_approved)`, p)

				if err != nil {
					log.Println(err)
					return []byte("internal server error"), http.StatusInternalServerError
				}
			}
		}
	}

	conversation, err := GetConversationByID(c.ID, db)
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	conversationJson, err := json.MarshalIndent(&conversation, "", "    ")
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	return conversationJson, http.StatusCreated
}

func (c *Conversation) Delete(db sqlx.Ext) ([]byte, int) {
	if uuid.Parse(c.ID) == nil {
		log.Println("conversation not found: " + c.ID)
		return []byte("conversation not found"), http.StatusNotFound
	}

	_, err := db.Exec(`
		DELETE FROM conversation_messages WHERE conversation_id = $1`, c.ID,
	)
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}

	_, err = db.Exec(`
		DELETE FROM conversation_participants WHERE conversation_id = $1`, c.ID,
	)
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}

	res, err := db.Exec(`
		DELETE FROM conversations WHERE id = $1`, c.ID,
	)
	count, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	if count < 1 {
		log.Println(err)
		return []byte("conversation not found"), http.StatusNotFound
	}
	return []byte("conversation deleted"), http.StatusCreated
}

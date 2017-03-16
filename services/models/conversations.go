package models

type Conversation struct {
	ID string `json:"id" form:"id" db:"id"`
	OwnerID string `json:"owner_id" form:"owner_id" db:"owner_id"`
	PhotoID string `json:"photo_id" form:"photo_id" db:"photo_id"`

	CreatedAt time.Time `json:"created_at" form:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" db:"updated_at"`
}

type ConversationParticipant struct {
	ConversationID string `json:"-" form:"-" db:"conversation_id"`
	UserID string `json:"-" form:"-" db:"user_id"`
	UserApproved bool `json:"-" form:"-" db:"user_approved"`

	CreatedAt time.Time `json"created_at" form:"created_at" db:"created_at"`
}

type ConversationMessage struct {
	ID string `json:"id" form:"id" db:"id"`
	ConversationID string `json:"-" form:"-" db:"conversation_id"`
	OwnerID string `json:"-" form:"-" db:"owner_id"`
	MessageText string `json:"message_text" form:"message_text" db:"message_text"`
}

type ConversationHttpResp struct {
	CoreConversation Conversation `json:"conversation"`
	Messages []*Message `json:"messages"`
}


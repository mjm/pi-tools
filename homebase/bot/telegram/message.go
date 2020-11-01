package telegram

type Message struct {
	MessageID int    `json:"message_id"`
	From      *User  `json:"from"`
	Date      int    `json:"date"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

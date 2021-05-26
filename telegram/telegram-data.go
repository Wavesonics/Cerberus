package telegram

///////////////////////////////////////////////////
// Request / Response objects
///////////////////////////////////////////////////

// WebhookUpdateBody https://core.telegram.org/bots/api#update
type WebhookUpdateBody struct {
	UpdateId      int64          `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	InlineQuery   *InlineQuery   `json:"inline_query,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

// SendMessageBody https://core.telegram.org/bots/api#sendmessage
type SendMessageBody struct {
	ChatID              int64                 `json:"chat_id"`
	Text                string                `json:"text"`
	ParseMode           string                `json:"parse_mode"`
	DisableNotification bool                  `json:"disable_notification"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// AnswerCallbackQueryBody https://core.telegram.org/bots/api#answercallbackquery
type AnswerCallbackQueryBody struct {
	CallbackQueryId string  `json:"callback_query_id"`
	Text            *string `json:"text,omitempty"`
	ShowAlert       *bool   `json:"show_alert,omitempty"`
	Url             *string `json:"url,omitempty"`
	CacheTime       *int    `json:"cache_time,omitempty"`
}

// EditMessageTextBody https://core.telegram.org/bots/api#editmessagetext
type EditMessageTextBody struct {
	ChatID                *int64                `json:"chat_id,omitempty"`
	MessageID             *int                  `json:"message_id,omitempty"`
	InlineMessageId       *string               `json:"inline_message_id,omitempty"`
	Text                  string                `json:"text"`
	ParseMode             string                `json:"parse_mode"`
	DisableWebPagePreview bool                  `json:"disable_web_page_preview"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageReplyMarkupBody https://core.telegram.org/bots/api#editmessagereplymarkup
type EditMessageReplyMarkupBody struct {
	ChatID          *int64                `json:"chat_id,omitempty"`
	MessageID       *int                  `json:"message_id,omitempty"`
	InlineMessageId *string               `json:"inline_message_id,omitempty"`
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// DeleteMessageBody https://core.telegram.org/bots/api#deletemessage
type DeleteMessageBody struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int   `json:"message_id"`
}

///////////////////////////////////////////////////
// Date Entities
///////////////////////////////////////////////////

type InlineQuery struct {
	Id     string `json:"id"`
	From   User   `json:"from"`
	Query  string `json:"query"`
	Offset string `json:"offset"`
}

type CallbackQuery struct {
	Id              string   `json:"id"`
	From            User     `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageId *string  `json:"inline_message_id,omitempty"`
	ChatInstance    string   `json:"chat_instance"`
	Data            *string  `json:"data,omitempty"`
	GameShortName   *string  `json:"game_short_name,omitempty"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	Text      string `json:"text"`
	From      *User  `json:"from,omitempty"`
	Chat      *Chat  `json:"chat,omitempty"`
	Date      int    `json:"date"`
}

type Chat struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type User struct {
	ID                    int64  `json:"id"`
	IsBot                 bool   `json:"is_bot"`
	SupportsInlineQueries bool   `json:"supports_inline_queries"`
	Username              string `json:"username"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string  `json:"text"`
	Url          *string `json:"url,omitempty"`
	CallbackData *string `json:"callback_data,omitempty"`
}

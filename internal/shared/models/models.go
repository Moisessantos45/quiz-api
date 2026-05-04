package models

import (
	"time"
)

type User struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Nickname  string    `json:"nickname" gorm:"not null"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Question struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	Text       string    `json:"text" gorm:"not null"`
	MediaType  string    `json:"media_type" gorm:"default:'text'"`
	CategoryID uint64    `json:"category_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`

	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
	Answers  []Answer `json:"options" gorm:"foreignKey:QuestionID"`
}

type Answer struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	Content    string    `json:"content" gorm:"not null"`
	IsCorrect  bool      `json:"is_correct" gorm:"not null"`
	QuestionID uint64    `json:"question_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`

	Question Question `json:"question" gorm:"foreignKey:QuestionID"`
}

type Game struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	Key        string    `json:"key" gorm:"unique;not null;index"`
	Name       string    `json:"name" gorm:"not null"`
	Status     string    `json:"status" gorm:"default:'waiting'"`
	UserID     uint64    `json:"user_id" gorm:"not null"`
	CategoryID uint64    `json:"category_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`

	User     User     `json:"user" gorm:"foreignKey:UserID"`
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
}

type Departure struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	GameID    uint64    `json:"game_id" gorm:"not null"`
	UserID    uint64    `json:"user_id" gorm:"not null"`
	Score     int64     `json:"score" gorm:"default:0"`
	Hits      int       `json:"hits" gorm:"default:0"`
	TotalTime int       `json:"total_time" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`

	Game    Game         `json:"game" gorm:"foreignKey:GameID"`
	User    User         `json:"user" gorm:"foreignKey:UserID"`
	Details []GameDetail `json:"details" gorm:"foreignKey:DepartureID"`
}

type GameDetail struct {
	ID           uint64 `json:"id" gorm:"primaryKey"`
	DepartureID  uint64 `json:"departure_id" gorm:"not null"`
	QuestionID   uint64 `json:"question_id" gorm:"not null"`
	AnswerID     uint64 `json:"answer_id" gorm:"not null"`
	IsCorrect    bool   `json:"is_correct"`
	ResponseTime int    `json:"response_time"`

	Departure Departure `json:"departure" gorm:"foreignKey:DepartureID"`
	Question  Question  `json:"question" gorm:"foreignKey:QuestionID"`
	Answer    Answer    `json:"option" gorm:"foreignKey:AnswerID"`
}

var Models = []any{
	&User{},
	&Category{},
	&Game{},
	&Departure{},
	&Question{},
	&Answer{},
	&GameDetail{},
}

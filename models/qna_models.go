package models

import (
// "database/sql"
// "time"
)

type QnA struct {
	QnAUUID   string `json:"qna_uuid" db:"qna_uuid"`
	Question  string `json:"question" db:"question"`
	Answer    string `json:"answer" db:"answer"`
	CreatedAt string `json:"created_at" db:"created_at"`
	CreatedBy string `json:"created_by" db:"created_by"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
	UpdatedBy string `json:"updated_by" db:"updated_by"`
	DeletedAt string `json:"deleted_at" db:"deleted_at"`
	DeletedBy string `json:"deleted_by" db:"deleted_by"`
}

type QnAResponse struct {
	QnAUUID   string `json:"qna_uuid" db:"qna_uuid"`
	Question  string `json:"question" db:"question"`
	Answer    string `json:"answer" db:"answer"`
	CreatedAt string `json:"created_at" db:"created_at"`
	CreatedBy string `json:"created_by" db:"created_by"`
}

type QnARequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

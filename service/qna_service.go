package service

import (
	"database/sql"
	"document/models"
	"log"
	"time"

	"github.com/google/uuid"
)

func GetAllQnA() ([]models.QnAResponse, error) {
	var qnaList []models.QnAResponse

	rows, err := db.Query(`
		SELECT 
			qna_uuid, 
			question, 
			answer,
			created_at,
			created_by
		FROM qna
		WHERE deleted_at IS NULL;
	`)
	if err != nil {
		log.Println("Error query GetAllQnA:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var qna models.QnAResponse
		err := rows.Scan(
			&qna.QnAUUID,
			&qna.Question,
			&qna.Answer,
			&qna.CreatedAt,
			&qna.CreatedBy,
		)
		if err != nil {
			log.Println("Error scan row:", err)
			return nil, err
		}
		qnaList = append(qnaList, qna)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, err
	}

	return qnaList, nil
}

func GetSpecQnA(id string) (*models.QnAResponse, error) {
	var qna models.QnAResponse

	err := db.QueryRow(`
		SELECT 
			qna_uuid, 
			question, 
			answer,
			created_at,
			created_by
		FROM qna
		WHERE qna_uuid = $1 AND deleted_at IS NULL;
	`, id).Scan(
		&qna.QnAUUID,
		&qna.Question,
		&qna.Answer,
		&qna.CreatedAt,
		&qna.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("QnA tidak ditemukan dengan UUID:", id)
			return nil, nil // Tidak error, tapi data kosong
		}
		log.Println("Error query GetSpecQnA:", err)
		return nil, err
	}

	return &qna, nil
}

func AddQnA(qna models.QnARequest, createdBy string) error {
	qnauuid := uuid.New()
	createdAt := time.Now()
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()
	qnaid := currentTimestamp + int64(uniqueID)

	query := `
		INSERT INTO qna (qna_id, qna_uuid, question, answer, created_at, created_by) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(query, qnaid, qnauuid.String(), qna.Question, qna.Answer, createdAt, createdBy)
	if err != nil {
		log.Println("Error saat menambahkan QnA:", err)
		return err
	}
	return nil
}

func UpdateQnA(id string, qna models.QnARequest, updatedBy string) error {
	updatedAt := time.Now()

	query := `
		UPDATE qna 
		SET question = $1, answer = $2, updated_at = $3, updated_by = $4 
		WHERE qna_uuid = $5
	`
	_, err := db.Exec(query, qna.Question, qna.Answer, updatedAt, updatedBy, id)
	if err != nil {
		log.Println("Error saat memperbarui QnA:", err)
		return err
	}
	return nil
}

func DeleteQnA(id string, deletedBy string) error {
	deletedAt := time.Now()

	query := `
		UPDATE qna 
		SET deleted_at = $1, deleted_by = $2 
		WHERE qna_uuid = $3
	`
	_, err := db.Exec(query, deletedAt, deletedBy, id)
	if err != nil {
		log.Println("Error saat menghapus QnA:", err)
		return err
	}
	return nil
}

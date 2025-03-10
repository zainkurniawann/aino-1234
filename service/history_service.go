package service

import (
	"database/sql"
	"document/models"

	// "fmt"
	"log"
	// "time"
	// "time"
	// "github.com/google/uuid"
)

func GetRecentTimelineHistorySuperAdmin(db *sql.DB) ([]models.TimelineHistory, error) {
	rows, err := db.Query(`
		SELECT 
			f.form_uuid, f.form_number, f.form_ticket, f.form_status,
			d.document_uuid, d.document_name,
			COALESCE(p.project_uuid, '') AS project_uuid, 
			COALESCE(p.project_name, '') AS project_name,
			f.created_by, f.created_at, 
			COALESCE(f.updated_by, '') AS updated_by, 
			COALESCE(f.updated_at, '1970-01-01 00:00:00') AS updated_at
		FROM form_ms f
		LEFT JOIN document_ms d ON f.document_id = d.document_id
		LEFT JOIN project_ms p ON f.project_id = p.project_id
		WHERE f.deleted_at IS NULL
		ORDER BY f.created_at DESC
		LIMIT 3;
	`)
	if err != nil {
		log.Println("Error executing recent query:", err)
		return nil, err
	}
	defer rows.Close()

	return scanTimelineHistory(rows)
}

func GetRecentTimelineHistoryAdmin(db *sql.DB, divisionCode string) ([]models.TimelineHistory, error) {
	rows, err := db.Query(`
		SELECT 
			f.form_uuid, 
			f.form_number, 
			f.form_ticket, 
			f.form_status,
			COALESCE(d.document_uuid, '') AS document_uuid, 
			COALESCE(d.document_name, '') AS document_name,
			COALESCE(p.project_uuid, '') AS project_uuid, 
			COALESCE(p.project_name, '') AS project_name,
			f.created_by, 
			f.created_at, 
			f.updated_by, 
			f.updated_at
		FROM form_ms f
		LEFT JOIN document_ms d ON f.document_id = d.document_id
		LEFT JOIN project_ms p ON f.project_id = p.project_id
		WHERE f.deleted_at IS NULL
		AND SPLIT_PART(f.form_number, '/', 2) = $1
		ORDER BY f.created_at DESC
		LIMIT 3;
	`, divisionCode)
	if err != nil {
		log.Println("Error executing recent query:", err)
		return nil, err
	}
	defer rows.Close()

	return scanTimelineHistory(rows)
}

func GetOlderTimelineHistorySuperAdmin(db *sql.DB, limit int, offset int) ([]models.TimelineHistory, error) {
	log.Printf("Fetching older timeline with limit %d and offset %d\n", limit, offset)

	// Query untuk mengambil data yang lebih lama dari 3 data terbaru
	rows, err := db.Query(`
		WITH RecentForms AS (
			SELECT form_uuid
			FROM form_ms
			WHERE deleted_at IS NULL
			ORDER BY created_at DESC
			LIMIT 3
		)
		SELECT 
			f.form_uuid, f.form_number, f.form_ticket, f.form_status,
			d.document_uuid, d.document_name,
			COALESCE(p.project_uuid, '') AS project_uuid, 
			COALESCE(p.project_name, '') AS project_name,
			f.created_by, f.created_at, 
			COALESCE(f.updated_by, '') AS updated_by, 
			COALESCE(f.updated_at, '1970-01-01 00:00:00') AS updated_at
		FROM form_ms f
		LEFT JOIN document_ms d ON f.document_id = d.document_id
		LEFT JOIN project_ms p ON f.project_id = p.project_id
		WHERE f.deleted_at IS NULL 
		AND f.form_uuid NOT IN (SELECT form_uuid FROM RecentForms)
		ORDER BY f.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		log.Println("Error executing older query:", err)
		return nil, err
	}
	defer rows.Close()

	// Memindai hasil query dan mengembalikan data
	return scanTimelineHistory(rows)
}

func GetOlderTimelineHistoryAdmin(db *sql.DB, divisionCode string, limit int, offset int) ([]models.TimelineHistory, error) {
	log.Printf("Fetching older timeline with limit %d and offset %d\n", limit, offset)

	// Query untuk mengambil data yang lebih lama dari 3 data terbaru
	rows, err := db.Query(`
       WITH RecentForms AS (
			SELECT form_uuid
			FROM form_ms
			WHERE deleted_at IS NULL
			AND SPLIT_PART(form_number, '/', 2) = $1
			ORDER BY created_at DESC
			LIMIT 3
		)
		SELECT 
			f.form_uuid, 
			f.form_number, 
			f.form_ticket, 
			f.form_status,
			COALESCE(d.document_uuid, '') AS document_uuid, 
			COALESCE(d.document_name, '') AS document_name,
			COALESCE(p.project_uuid, '') AS project_uuid, 
			COALESCE(p.project_name, '') AS project_name,
			f.created_by, 
			f.created_at, 
			f.updated_by, 
			f.updated_at
		FROM form_ms f
		LEFT JOIN document_ms d ON f.document_id = d.document_id
		LEFT JOIN project_ms p ON f.project_id = p.project_id
		WHERE f.deleted_at IS NULL 
		AND SPLIT_PART(f.form_number, '/', 2) = $1
		AND f.form_uuid NOT IN (SELECT form_uuid FROM RecentForms)
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3;
    `, divisionCode, limit, offset)
	if err != nil {
		log.Println("Error executing older query:", err)
		return nil, err
	}
	defer rows.Close()

	// Memindai hasil query dan mengembalikan data
	return scanTimelineHistory(rows)
}

func scanTimelineHistory(rows *sql.Rows) ([]models.TimelineHistory, error) {
	var historyList []models.TimelineHistory

	for rows.Next() {
		var history models.TimelineHistory
		var updatedBy sql.NullString
		var updatedAt sql.NullTime
		var formUUID string

		err := rows.Scan(
			&formUUID, &history.FormNumber, &history.FormTicket, &history.FormStatus,
			&history.DocumentUUID, &history.DocumentName,
			&history.ProjectUUID, &history.ProjectName,
			&history.CreatedBy, &history.CreatedAt, &updatedBy, &updatedAt,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		history.FormUUID = formUUID
		history.UpdatedBy = updatedBy
		history.UpdatedAt = updatedAt

		historyList = append(historyList, history)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in row iteration:", err)
		return nil, err
	}

	// Jika tidak ada data, return list kosong, bukan error
	if len(historyList) == 0 {
		log.Println("No other data available")
		return []models.TimelineHistory{}, nil
	}

	return historyList, nil
}

func GetDocumentCountPerMonthSuperAdmin(db *sql.DB, year int) ([]models.MonthlyDocumentCount, error) {
	query := `
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM') AS month, 
			COUNT(*) AS count
		FROM form_ms
		WHERE EXTRACT(YEAR FROM created_at) = $1
		AND deleted_at IS NULL
		GROUP BY month
		ORDER BY month;
	`

	rows, err := db.Query(query, year)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var counts []models.MonthlyDocumentCount
	for rows.Next() {
		var count models.MonthlyDocumentCount
		if err := rows.Scan(&count.Month, &count.Count); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		counts = append(counts, count)
	}

	if counts == nil {
		return []models.MonthlyDocumentCount{}, nil
	}

	return counts, nil
}

func GetDocumentCountPerMonthAdmin(db *sql.DB, year int, divisionCode string) ([]models.MonthlyDocumentCount, error) {
	query := `
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM') AS month, 
			COUNT(*) AS count
		FROM form_ms
		WHERE EXTRACT(YEAR FROM created_at) = $1
		AND deleted_at IS NULL
		AND SPLIT_PART(form_number, '/', 2) = $2
		GROUP BY month
		ORDER BY month;
	`

	rows, err := db.Query(query, year, divisionCode)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var counts []models.MonthlyDocumentCount
	for rows.Next() {
		var count models.MonthlyDocumentCount
		if err := rows.Scan(&count.Month, &count.Count); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		counts = append(counts, count)
	}

	if counts == nil {
		return []models.MonthlyDocumentCount{}, nil
	}

	return counts, nil
}

// GetDocumentStatusCountPerMonth menghitung jumlah dokumen berdasarkan status dalam bulan tertentu
func GetDocumentStatusCountPerMonthSuperAdmin(db *sql.DB, year, month int) ([]models.DocumentStatusCount, error) {
	query := `
		SELECT 
			form_status AS status, 
			COUNT(*) AS count
		FROM form_ms
		WHERE EXTRACT(YEAR FROM created_at) = $1 
		AND EXTRACT(MONTH FROM created_at) = $2
		AND deleted_at IS NULL
		GROUP BY form_status
		ORDER BY count DESC;
	`

	rows, err := db.Query(query, year, month)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var statusCounts []models.DocumentStatusCount
	for rows.Next() {
		var statusCount models.DocumentStatusCount
		if err := rows.Scan(&statusCount.Status, &statusCount.Count); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		statusCounts = append(statusCounts, statusCount)
	}

	if statusCounts == nil {
		return []models.DocumentStatusCount{}, nil
	}

	return statusCounts, nil
}

func GetDocumentStatusCountPerMonthAdmin(db *sql.DB, year, month int, divisionCode string) ([]models.DocumentStatusCount, error) {
	query := `
		SELECT 
			form_status AS status, 
			COUNT(*) AS count
		FROM form_ms
		WHERE EXTRACT(YEAR FROM created_at) = $1 
		AND EXTRACT(MONTH FROM created_at) = $2
		AND deleted_at IS NULL
		AND SPLIT_PART(form_number, '/', 2) = $3
		GROUP BY form_status
		ORDER BY count DESC;
	`

	rows, err := db.Query(query, year, month, divisionCode)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var statusCounts []models.DocumentStatusCount
	for rows.Next() {
		var statusCount models.DocumentStatusCount
		if err := rows.Scan(&statusCount.Status, &statusCount.Count); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		statusCounts = append(statusCounts, statusCount)
	}

	if statusCounts == nil {
		return []models.DocumentStatusCount{}, nil
	}

	return statusCounts, nil
}

func GetFormCountPerDocumentPerMonthSuperAdmin(db *sql.DB, year int, month int) ([]models.MonthlyFormCount, error) {
	query := `
		SELECT  
			d.document_name AS document_name, 
			EXTRACT(MONTH FROM f.created_at) AS month,
			COUNT(f.form_id) AS count
		FROM form_ms f
		JOIN document_ms d ON f.document_id = d.document_id
		WHERE EXTRACT(YEAR FROM f.created_at) = $1
		AND EXTRACT(MONTH FROM f.created_at) = $2
		AND f.deleted_at IS NULL
		GROUP BY d.document_name, EXTRACT(MONTH FROM f.created_at) 
		ORDER BY d.document_name;
    `

	rows, err := db.Query(query, year, month)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var counts []models.MonthlyFormCount
	for rows.Next() {
		var count models.MonthlyFormCount
		if err := rows.Scan(&count.DocumentName, &count.Month, &count.Count); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		counts = append(counts, count)
	}

	if counts == nil {
		return []models.MonthlyFormCount{}, nil
	}

	return counts, nil
}

func GetFormCountPerDocumentPerMonthAdmin(db *sql.DB, year int, month int, divisionCode string) ([]models.MonthlyFormCount, error) {
	query := `
		SELECT 
			d.document_name AS document_name, 
			EXTRACT(MONTH FROM f.created_at) AS month,
			COUNT(f.form_id) AS count
		FROM form_ms f
		JOIN document_ms d ON f.document_id = d.document_id
		WHERE EXTRACT(YEAR FROM f.created_at) = $1
		AND EXTRACT(MONTH FROM f.created_at) = $2
		AND f.deleted_at IS NULL
		AND SPLIT_PART(f.form_number, '/', 2) = $3
		GROUP BY d.document_name, EXTRACT(MONTH FROM f.created_at) 
		ORDER BY d.document_name;
	`

	rows, err := db.Query(query, year, month, divisionCode)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var counts []models.MonthlyFormCount
	for rows.Next() {
		var count models.MonthlyFormCount
		if err := rows.Scan(&count.DocumentName, &count.Month, &count.Count); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		counts = append(counts, count)
	}

	if counts == nil {
		return []models.MonthlyFormCount{}, nil
	}

	return counts, nil
}
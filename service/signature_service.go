package service

import (
	// "database/sql"
	"document/models"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func GetAllPersonalName() ([]models.Personal, error) {
	getUserAppRole := []models.Personal{}

	// Lakukan query ke database lain
	rows, err := db2.Queryx("SELECT u.user_id, pd.personal_name FROM user_ms u JOIN personal_data_ms pd ON u.user_id = pd.user_id WHERE u.deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		place := models.Personal{}
		err := rows.StructScan(&place)
		if err != nil {
			log.Println("Error scanning row to struct:", err)
			continue
		}
		getUserAppRole = append(getUserAppRole, place)
	}

	return getUserAppRole, nil
}

func GetSignatureForm(id string) ([]models.Signatories, error) {
	var signatories []models.Signatories

	err := db.Select(&signatories, `SELECT 
	sf.sign_uuid, 
	sf.name, 
	sf.position, 
	sf.role_sign, 
	sf.is_sign, 
	sf.created_by, 
	sf.created_at, 
	sf.updated_by, 
	sf.updated_at, 
	sf.deleted_by, 
	sf.deleted_at
FROM 
	sign_form sf 
	JOIN form_ms fm ON sf.form_id = fm.form_id 
WHERE 
	fm.form_uuid = $1 AND sf.deleted_at IS NULL`, id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return signatories, nil

}

func GetSpecSignatureByID(id string) (models.Signatorie, error) {
	var signatories models.Signatorie
	err := db.Get(&signatories, "SELECT sign_uuid, name, position, role_sign, is_sign, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM sign_form sf WHERE sign_uuid = $1 AND deleted_at IS NULL", id)
	if err != nil {
		log.Print(err)
		return models.Signatorie{}, err
	}

	return signatories, nil
}

func GetUserIDSign(id string) (models.UserIDSign, error) {
	var userID models.UserIDSign
	err := db.Get(&userID, "SELECT user_id, sign_uuid FROM sign_form WHERE sign_uuid = $1", id)
	if err != nil {
		log.Print(err)
		return models.UserIDSign{}, err
	}
	return userID, nil
}

func UpdateFormSignature(updateSign models.UpdateSign, id string, username string) error {

	currentTime := time.Now()

	_, err := db.NamedExec("UPDATE sign_form SET is_sign = :is_sign, updated_by = :updated_by, updated_at = :updated_at, sign_img = :image WHERE sign_uuid = :id", map[string]interface{}{
		"is_sign":    updateSign.IsSign,
		"updated_by": username,
		"updated_at": currentTime,
		"image":      updateSign.Image, // Menyimpan path gambar
		"id":         id,
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdateFormSignatureGuest(updateSign models.UpdateSignGuest, id string) error {

	currentTime := time.Now()

	_, err := db.NamedExec("UPDATE sign_form SET is_sign = :is_sign, updated_by = :updated_by, updated_at = :updated_at, sign_img = :image WHERE sign_uuid = :id", map[string]interface{}{
		"is_sign":    updateSign.IsSign,
		"updated_by": updateSign.Name,
		"updated_at": currentTime,
		"image":      updateSign.Image, // Menyimpan path gambar
		"id":         id,
	})
	if err != nil {
		return err
	}
	return nil
}

// ga kepake, tp ku edit bjir
// func GetUserRoleByFormID(userID int, formUUID string) (string, error) {
// 	var roleSign string
// 	query := `
// 		SELECT sf.role_sign
// 		FROM sign_form sf
// 		JOIN form_ms f ON sf.form_id = f.form_id
// 		WHERE sf.user_id = $1 AND f.form_uuid = $2
// 		LIMIT 1
// 	`
// 	err := db.Get(&roleSign, query, userID, formUUID)
// 	if err != nil {
// 		return "", err
// 	}
// 	return roleSign, nil
// }

func GetUserRolesByFormID(userID int, formUUID string) ([]string, error) {
	var roles []string
	query := `
        SELECT sf.role_sign
        FROM sign_form sf
        JOIN form_ms f ON sf.form_id = f.form_id
        WHERE sf.user_id = $1 AND f.form_uuid = $2
    `
	err := db.Select(&roles, query, userID, formUUID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func AddApproval(addApproval models.AddApproval, id string, username string, userID int) error {
	// Ambil role_sign dari database
	userRole, err := GetUserRolesByFormID(userID, id)
	if err != nil {
		return err
	}

	isAuthorized := false
	for _, role := range userRole {
		if role == "Atasan Penerima" {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return fmt.Errorf("hanya 'Atasan Penerima' yang dapat menyetujui atau menolak form")
	}

	// Lanjutkan dengan pembaruan form jika role_sign valid
	currentTime := time.Now()
	fmt.Printf("Debug: IsApprove=%v, Reason=%s, UpdatedBy=%s, ID=%s\n", addApproval.IsApproval, addApproval.Reason, username, id)
	_, err = db.Exec("UPDATE form_ms SET is_approve = $1, reason = $2, updated_by = $3, updated_at = $4 WHERE form_uuid = $5",
		addApproval.IsApproval, addApproval.Reason, username, currentTime, id)

	if err != nil {
		return err
	}
	return nil
}

// sekalian agar nomor da auto masuk ke form itcm
func AddApprovalDA(addApproval models.AddApproval, id string, username string, userID int) error {
	// Ambil role_sign dari database
	userRole, err := GetUserRolesByFormID(userID, id)
	if err != nil {
		return err
	}

	isAuthorized := false
	for _, role := range userRole {
		if role == "Atasan Penerima" {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return fmt.Errorf("hanya 'Atasan Penerima' yang dapat menyetujui atau menolak form")
	}

	// Query untuk mengambil form_id berdasarkan id yang diberikan
	var formDAID string
	err = db.QueryRow(`
        SELECT form_id
        FROM form_ms
        WHERE form_uuid = $1
    `, id).Scan(&formDAID)

	if err != nil {
		return err
	}

	// Ambil itcm_form_uuid dari database
	var itcmFormUUID string
	err = db.QueryRow(`
    SELECT 
        (f.form_data->>'itcm_form_uuid')::text AS itcm_form_uuid
    FROM form_ms f
    LEFT JOIN document_ms d ON f.document_id = d.document_id
    WHERE f.form_uuid = $1 AND d.document_code = 'DA'
`, id).Scan(&itcmFormUUID)
	if err != nil {
		return err
	}

	// ambil form_id milik itcm berdasarkan itcmFormUUID
	var formITCMID string
	err = db.QueryRow(`
			SELECT form_id
			FROM form_ms
			WHERE form_uuid = $1
	`, itcmFormUUID).Scan(&formITCMID)
	if err != nil {
		return err
	}

	var formDANumber string
	err = db.QueryRow(`
	SELECT form_number
	FROM form_ms
	WHERE form_id = $1
`, formDAID).Scan(&formDANumber)
	if err != nil {
		return err
	}

	fmt.Println("ID milik DA", formDAID)
	fmt.Println("UUID milik ITCM", itcmFormUUID)
	fmt.Println("ID milik ITCM", formITCMID)
	fmt.Println("number milik da", formDANumber)

	// Pastikan nilai `IsApproval` sudah didefinisikan dalam `addApproval`
	if addApproval.IsApproval {
		_, err := db.Exec(`
		UPDATE form_ms 
		SET form_data = JSONB_SET(form_data, '{no_da}', to_jsonb($1::text))
		WHERE form_id = $2;
	`, formDANumber, formITCMID)
		if err != nil {
			return err
		}
	}

	// Lanjutkan dengan pembaruan form jika role_sign valid
	currentTime := time.Now()
	_, err = db.NamedExec("UPDATE form_ms SET is_approve = :is_approve, reason = :reason, updated_by = :updated_by, updated_at = :updated_at WHERE form_uuid = :id", map[string]interface{}{
		"is_approve": addApproval.IsApproval,
		"reason":     addApproval.Reason,
		"updated_by": username,
		"updated_at": currentTime,
		"id":         id,
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdateSignInfo(signatory models.UpdateSignForm, id string, username string) (models.UpdateSignForm, error) {
	currentTime := time.Now()
	personalName := signatory.Name

	personalNames, err := GetAllPersonalName()
	if err != nil {
		log.Println("Error getting personal names:", err)
		return models.UpdateSignForm{}, err
	}
	var userID int
	for _, personal := range personalNames {
		if personal.PersonalName == personalName {
			userID, err = strconv.Atoi(personal.UserID)
			if err != nil {
				return models.UpdateSignForm{}, err
			}
			break
		}
	}

	if userID == 0 {
		log.Printf("User ID not found for personal name: %s\n", personalName)
		return models.UpdateSignForm{}, errors.New("User ID not found for personal name")
	}

	_, err = db.NamedExec("UPDATE sign_form SET user_id = :user_id, name = :name, position = :position, role_sign = :role_sign, updated_by = :updated_by, updated_at = :updated_at WHERE sign_uuid = :sign_uuid", map[string]interface{}{
		"user_id":    userID,
		"name":       personalName,
		"position":   signatory.Position,
		"role_sign":  signatory.Role,
		"updated_by": username,
		"updated_at": currentTime,
		"sign_uuid":  id,
	})
	if err != nil {
		return models.UpdateSignForm{}, err
	}

	return signatory, nil
}

func AddSignInfo(signatory models.AddSignInfo, username string) error {
	currentTime := time.Now()
	uuidObj := uuid.New()
	uuidString := uuidObj.String()

	var formID int
	err := db.Get(&formID, "SELECT form_id FROM form_ms WHERE form_uuid = $1", signatory.FormUUID)
	if err != nil {
		log.Println("Error getting form_id:", err)
		return err
	}
	personalName := signatory.Name
	personalNames, err := GetAllPersonalName()
	if err != nil {
		log.Println("Error getting personal names:", err)
		return err
	}
	var userID int
	for _, personal := range personalNames {
		if personal.PersonalName == personalName {
			userID, err = strconv.Atoi(personal.UserID)
			if err != nil {
				return err
			}
			break
		}
	}

	if userID == 0 {
		log.Printf("User ID not found for personal name: %s\n", personalName)
		return errors.New("User ID not found for personal name")
	}
	_, err = db.NamedExec("INSERT INTO sign_form (sign_uuid, form_id, user_id, name, position, role_sign, created_by) VALUES (:sign_uuid, :form_id, :user_id, :name, :position, :role_sign, :created_by)", map[string]interface{}{
		"sign_uuid":  uuidString,
		"user_id":    userID,
		"form_id":    formID,
		"name":       signatory.Name,
		"position":   signatory.Position,
		"role_sign":  signatory.Role,
		"created_by": username,
		"created_at": currentTime,
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteSignInfo(id, username string) error {
	currentTime := time.Now()
	result, err := db.NamedExec("UPDATE sign_form SET deleted_by = :deleted_by, deleted_at = :deleted_at WHERE sign_uuid = :id", map[string]interface{}{
		"deleted_by": username,
		"deleted_at": currentTime,
		"id":         id,
	})

	if err != nil {
		log.Print(err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound // Mengembalikan error jika tidak ada rekaman yang cocok
	}

	return nil
}

func SignatureNotif(userID int) ([]models.Notif, error) {
	rows, err := db.Query(`
SELECT 
		f.form_uuid, f.form_number, f.form_ticket, f.form_status,
		d.document_code, d.document_name, sf.role_sign, sf.is_sign, sf.created_at, sf.updated_at, sf.deleted_at
		FROM 
		form_ms f
	LEFT JOIN 
		document_ms d ON f.document_id = d.document_id
	LEFT JOIN 
		sign_form sf ON f.form_id = sf.form_id
		WHERE
		sf.user_id = $1 AND f.deleted_at IS NULL AND f.form_status = 'Published'
		ORDER BY sf.created_at DESC
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all form data
	var forms []models.Notif

	// Iterate through the rows
	for rows.Next() {
		// Scan the row into the Forms struct
		var form models.Notif
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentCode,
			&form.DocumentName,
			&form.RoleSign,
			&form.IsSign,
			&form.CreatedAt,
			&form.UpdatedAt,
			&form.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}
	// Return the forms as JSON response
	return forms, nil
}

func ApproveNotif(userID int) ([]models.NotifApproval, error) {
	rows, err := db.Query(`
SELECT 
		f.form_uuid, f.form_number, f.form_ticket, f.is_approve,
		d.document_code, d.document_name, sf.role_sign, sf.is_sign, sf.created_at, sf.updated_at, sf.deleted_at
		FROM 
		form_ms f
	LEFT JOIN 
		document_ms d ON f.document_id = d.document_id
	LEFT JOIN 
		sign_form sf ON f.form_id = sf.form_id
		WHERE
		sf.user_id = $1 AND f.deleted_at IS NULL AND f.is_approve = True
		ORDER BY sf.created_at DESC
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all form data
	var forms []models.NotifApproval

	// Iterate through the rows
	for rows.Next() {
		// Scan the row into the Forms struct
		var form models.NotifApproval
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.IsApprove,
			&form.DocumentCode,
			&form.DocumentName,
			&form.RoleSign,
			&form.IsSign,
			&form.CreatedAt,
			&form.UpdatedAt,
			&form.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}
	// Return the forms as JSON response
	return forms, nil
}

package service

import (
	"database/sql"
	"document/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func generateFormNumberHA(documentID int64, divisionCode string, recursionCount int) (string, error) {
	const maxRecursionCount = 1000

	// Check if the maximum recursion count is reached
	if recursionCount > maxRecursionCount {
		return "", errors.New("Maximum recursion count exceeded")
	}

	documentCode, err := GetDocumentCode(documentID)
	if err != nil {
		return "", fmt.Errorf("Failed to get document code: %v", err)
	}

	// Get the latest form number for the given document ID
	var latestFormNumber sql.NullString
	err = db.Get(&latestFormNumber, "SELECT MAX(form_number) FROM form_ms WHERE document_id = $1", documentID)
	if err != nil {
		return "", fmt.Errorf("Error getting latest form number: %v", err)
	}

	// Initialize formNumber to 1 if latestFormNumber is NULL
	formNumber := 1
	if latestFormNumber.Valid {
		// Parse the latest form number
		var latestFormNumberInt int
		_, err := fmt.Sscanf(latestFormNumber.String, "%d", &latestFormNumberInt)
		if err != nil {
			return "", fmt.Errorf("Error parsing latest form number: %v", err)
		}
		// Increment the latest form number
		formNumber = latestFormNumberInt + 1
	}

	// Get current year and month
	year := time.Now().Year()
	month := time.Now().Month()

	// Convert month to Roman numeral
	romanMonth, err := convertToRoman(int(month))
	if err != nil {
		return "", fmt.Errorf("Error converting month to Roman numeral: %v", err)
	}

	fmt.Println("latest", latestFormNumber)
	fmt.Println("document code", documentCode)

	// Format the form number according to the specified format
	formNumberString := fmt.Sprintf("%04d", formNumber)
	formNumberWithDivision := fmt.Sprintf("%s/%s/%s/%s/%d", formNumberString, divisionCode, documentCode, romanMonth, year)
	// formNumberWithDivision := fmt.Sprintf("%s/%s/%s/%s/%d", formNumberString, "PED", "F", romanMonth, year)

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM form_ms WHERE form_number = $1 and document_id = $2", formNumberString, documentID)
	if err != nil {
		return "", fmt.Errorf("Error checking existing form number: %v", err)
	}
	if count > 0 {
		// If the form number already exists, recursively call the function again
		return generateFormNumberHA(documentID, divisionCode, recursionCount+1)
	}

	fmt.Println(formNumberWithDivision)
	return formNumberWithDivision, nil
}

func AddHakAkses(addForm models.FormHA, infoHA []models.AddInfoHAReq, haReq models.HAReq, isPublished bool, userID int, divisionCode string, recrusionCount int, username string, signatories []models.Signatory) error {
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()
	appID := currentTimestamp + int64(uniqueID)
	uuidObj := uuid.New()
	uuidString := uuidObj.String()

	formStatus := "Draft"
	if isPublished {
		formStatus = "Published"
	}

	var documentID int64
	err := db.Get(&documentID, "SELECT document_id FROM document_ms WHERE document_uuid = $1", addForm.DocumentUUID)
	if err != nil {
		log.Println("Error getting document_id:", err)
		return err
	}

	formNumberHA, err := generateFormNumberHA(documentID, divisionCode, recrusionCount+1)
	if err != nil {
		// Handle error
		log.Println("Error generating form number:", err)
		return err
	}

	formData, err := json.Marshal(haReq)
	if err != nil {
		log.Println("Error marshaling ITCM struct:", err)
		return err
	}
	_, err = db.NamedExec("INSERT INTO form_ms (form_id, form_uuid, document_id, user_id, project_id, form_number, form_ticket, form_status, form_data, created_by) VALUES (:form_id, :form_uuid, :document_id, :user_id, :project_id, :form_number, :form_ticket, :form_status, :form_data, :created_by)", map[string]interface{}{
		"form_id":     appID,
		"form_uuid":   uuidString,
		"document_id": documentID,
		"user_id":     userID,
		"project_id":  nil,
		"form_number": formNumberHA,
		"form_ticket": addForm.FormTicket,
		"form_status": formStatus,
		"form_data":   formData, // Convert JSON to string
		"created_by":  username,
	})

	// fmt.Println(formNumberHA)

	if err != nil {
		return err
	}
	personalNames, err := GetAllPersonalName() // Mengambil daftar semua personal name
	if err != nil {
		log.Println("Error getting personal names:", err)
		return err
	}

	fmt.Println("bejir", infoHA)
	for _, info := range infoHA {
		uuidString := uuid.New().String()

		_, err := db.NamedExec("INSERT INTO hak_akses (ha_uuid, form_id, nama_pengguna, ruang_lingkup, jangka_waktu, created_by) VALUES (:ha_uuid, :form_id, :nama_pengguna, :ruang_lingkup, :jangka_waktu, :created_by)", map[string]interface{}{
			"ha_uuid":       uuidString,
			"form_id":       appID,
			"nama_pengguna": info.NamaPengguna,
			"ruang_lingkup": info.RuangLingkup,
			"jangka_waktu":  info.JangkaWaktu,
			"created_by":    username,
		})
		if err != nil {
			return err
		}
	}

	for _, signatory := range signatories {
		uuidString := uuid.New().String()

		// Mencari user_id yang sesuai dengan personal_name yang dipilih
		var userID string
		for _, personal := range personalNames {
			if personal.PersonalName == signatory.Name {
				userID = personal.UserID
				break
			}
		}

		// Memastikan user_id ditemukan untuk personal_name yang dipilih
		if userID == "" {
			log.Printf("User ID not found for personal name: %s\n", signatory.Name)
			continue
		}

		_, err := db.NamedExec("INSERT INTO sign_form (sign_uuid, form_id, user_id, name, position, role_sign, created_by) VALUES (:sign_uuid, :form_id, :user_id, :name, :position, :role_sign, :created_by)", map[string]interface{}{
			"sign_uuid":  uuidString,
			"user_id":    userID,
			"form_id":    appID,
			"name":       signatory.Name,
			"position":   signatory.Position,
			"role_sign":  signatory.Role,
			"created_by": username,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func AddHakAksesReview(addForm models.FormHA, infoHA []models.AddInfoHA, ha models.HA, isPublished bool, userID int, divisionCode string, recrusionCount int, username string, signatories []models.Signatory) error {
    tx, err := db.Beginx()
    if err != nil {
        log.Println("Error starting transaction:", err)
        return err
    }

    // Defer rollback jika ada error
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            log.Println("Transaction rolled back due to panic:", r)
        }
    }()

    currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
    uniqueID := uuid.New().ID()
    appID := currentTimestamp + int64(uniqueID)
    uuidObj := uuid.New()
    uuidString := uuidObj.String()

    formStatus := "Draft"
    if isPublished {
        formStatus = "Published"
    }

    var documentID int64
    err = tx.Get(&documentID, "SELECT document_id FROM document_ms WHERE document_uuid = $1", addForm.DocumentUUID)
    if err != nil {
        log.Println("Error getting document_id:", err)
        tx.Rollback()
        return err
    }

    formNumberHA, err := generateFormNumberHA(documentID, divisionCode, recrusionCount+1)
    if err != nil {
        log.Println("Error generating form number:", err)
        tx.Rollback()
        return err
    }

    formData, err := json.Marshal(ha)
    if err != nil {
        log.Println("Error marshaling HA struct:", err)
        tx.Rollback()
        return err
    }

    _, err = tx.NamedExec("INSERT INTO form_ms (form_id, form_uuid, document_id, user_id, project_id, form_number, form_ticket, form_status, form_data, created_by) VALUES (:form_id, :form_uuid, :document_id, :user_id, :project_id, :form_number, :form_ticket, :form_status, :form_data, :created_by)", map[string]interface{}{
        "form_id":     appID,
        "form_uuid":   uuidString,
        "document_id": documentID,
        "user_id":     userID,
        "project_id":  nil,
        "form_number": formNumberHA,
        "form_ticket": addForm.FormTicket,
        "form_status": formStatus,
        "form_data":   formData,
        "created_by":  username,
    })
    if err != nil {
        log.Println("Error inserting into form_ms:", err)
        tx.Rollback()
        return err
    }

    personalNames, err := GetAllPersonalName()
    if err != nil {
        log.Println("Error getting personal names:", err)
        tx.Rollback()
        return err
    }

    for _, info := range infoHA {
        uuidString := uuid.New().String()

        _, err := tx.NamedExec("INSERT INTO hak_akses_info (info_uuid, form_id, name, instansi, position, username, password, scope, created_by) VALUES (:info_uuid, :form_id, :name, :instansi, :position, :username, :password, :scope, :created_by)", map[string]interface{}{
            "info_uuid":  uuidString,
            "form_id":    appID,
            "name":       info.Name,
            "instansi":   info.Instansi,
            "position":   info.Position,
            "username":   info.Username,
            "password":   info.Password,
            "scope":      info.Scope,
            "created_by": username,
        })
        if err != nil {
            log.Println("Error inserting into hak_akses_info:", err)
            tx.Rollback()
            return err
        }
    }

    for _, signatory := range signatories {
        uuidString := uuid.New().String()
        var signatoryUserID string

        for _, personal := range personalNames {
            if personal.PersonalName == signatory.Name {
                signatoryUserID = personal.UserID
                break
            }
        }

        if signatoryUserID == "" {
            log.Printf("User ID not found for personal name: %s\n", signatory.Name)
            tx.Rollback()
            return fmt.Errorf("user ID not found for personal name: %s", signatory.Name)
        }

        _, err := tx.NamedExec("INSERT INTO sign_form (sign_uuid, form_id, user_id, name, position, role_sign, created_by) VALUES (:sign_uuid, :form_id, :user_id, :name, :position, :role_sign, :created_by)", map[string]interface{}{
            "sign_uuid":  uuidString,
            "user_id":    signatoryUserID,
            "form_id":    appID,
            "name":       signatory.Name,
            "position":   signatory.Position,
            "role_sign":  signatory.Role,
            "created_by": username,
        })
        if err != nil {
            log.Println("Error inserting into sign_form:", err)
            tx.Rollback()
            return err
        }
    }

    err = tx.Commit()
    if err != nil {
        log.Println("Error committing transaction:", err)
        tx.Rollback()
        return err
    }

    return nil
}


func GetAllHakAkses() ([]models.FormsHAReq, error) {
	rows, err := db.Query(`SELECT
    f.form_uuid,
    f.form_number,
    f.form_ticket,
    f.form_status,
    d.document_name,
    f.created_by,
    f.created_at,
    f.updated_by,
    f.updated_at,	
    f.deleted_by,
    f.deleted_at,
    (f.form_data->>'form_type')::text AS form_type,
    (f.form_data->>'form_name')::text AS form_name,
    (f.form_data->>'nama_tim')::text AS nama_tim,
    (f.form_data->>'product_manager')::text AS product_manager,
    (f.form_data->>'nama_pengusul')::text AS nama_pengusul,
    (f.form_data->>'tanggal_usul')::text AS tanggal_usul
FROM
    form_ms f
LEFT JOIN
    document_ms d ON f.document_id = d.document_id
WHERE
    d.document_code = 'HA' AND f.deleted_at IS NULL
    AND ((f.form_data->>'form_type') = 'Permintaan' OR (f.form_data->>'form_type') = 'Penghapusan')
ORDER BY f.form_number DESC;

	`)
	var forms []models.FormsHAReq
	//rows, err := db.Query(&forms, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var form models.FormsHAReq
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&form.FormType, // Urutan form_type dipindahkan ke depan
			&form.FormName, // Urutan form_type dipindahkan ke depan
			&form.NamaTim,
			&form.ProductManager,
			&form.NamaPengusul,
			&form.TanggalUsul,
		)

		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}

	return forms, nil
}

func GetAllHakAksesReview() ([]models.FormsHA, error) {
	rows, err := db.Query(`SELECT
		f.form_uuid,
		f.form_number,
		f.form_ticket,
		f.form_status,
		d.document_name,
		f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
		(f.form_data->>'form_name')::text AS form_name
	FROM
		form_ms f
	LEFT JOIN
		document_ms d ON f.document_id = d.document_id
	WHERE
		d.document_code = 'HA' AND f.deleted_at IS NULL
    AND (f.form_data->>'form_type') = 'Review'
		ORDER BY f.form_number DESC;
	`)
	var forms []models.FormsHA
	//rows, err := db.Query(&forms, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var form models.FormsHA
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&form.FormName,
		)
		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}

	return forms, nil
}

type FormwoilahWithSignatories struct {
	Form        models.FormsHAReq        `json:"form"`
	InfoHA      []models.HakAksesRequest `json:"hak_akses_info"`
	Signatories []models.SignatoryHA     `json:"signatories"`
}

func GetSpecAllHA(id string) (*FormwoilahWithSignatories, error) {
	var FormWithSignatories FormwoilahWithSignatories

	// Ambil data form
	err := db.Get(&FormWithSignatories.Form, `
			SELECT 
					f.form_uuid,
					f.form_number,
					f.form_status,
					d.document_name,
					CASE
						WHEN f.is_approve IS NULL THEN 'Menunggu Disetujui'
						WHEN f.is_approve = false THEN 'Tidak Disetujui'
						WHEN f.is_approve = true THEN 'Disetujui'
					END AS approval_status,
					COALESCE(f.reason, '') AS reason,
					f.created_by,
					f.created_at,
					f.updated_by,
					f.updated_at,
					f.deleted_by,
					f.deleted_at,
					(f.form_data->>'form_type')::text AS form_type,
					(f.form_data->>'form_name')::text AS form_name,
					(f.form_data->>'nama_tim')::text AS nama_tim,
					(f.form_data->>'product_manager')::text AS product_manager,
					(f.form_data->>'nama_pengusul')::text AS nama_pengusul,
					(f.form_data->>'tanggal_usul')::text AS tanggal_usul
			FROM
					form_ms f
			LEFT JOIN 
					document_ms d ON f.document_id = d.document_id
			WHERE
					f.form_uuid = $1 AND d.document_code = 'HA' AND f.deleted_at IS NULL
    			AND ((f.form_data->>'form_type') = 'Permintaan' OR (f.form_data->>'form_type') = '	')
	`, id)
	if err != nil {
		return nil, err
	}

	// Ambil data hak akses info
	err = db.Select(&FormWithSignatories.InfoHA, `
			SELECT 
					ha_uuid,
					nama_pengguna,
					ruang_lingkup,
					jangka_waktu
			FROM
					hak_akses
			WHERE
					form_id IN (
							SELECT form_id FROM form_ms WHERE form_uuid = $1 AND deleted_at IS NULL
					)
	`, id)
	if err != nil {
		return nil, err
	}

	// Ambil data signatories
	err = db.Select(&FormWithSignatories.Signatories, `
			 SELECT 
            sign_uuid,
            name AS signatory_name,
            position AS signatory_position,
            role_sign,
            is_sign,
						CASE
							WHEN sign_img IS NOT NULL AND sign_img != '' THEN CONCAT('/assets/images/signatures/', sign_img)
							ELSE ''
						END AS sign_img,
						updated_at
        FROM
					sign_form
			WHERE
					form_id IN (
							SELECT form_id FROM form_ms WHERE form_uuid = $1 AND deleted_at IS NULL
					)
	`, id)
	if err != nil {
		return nil, err
	}

	return &FormWithSignatories, nil
}

type FormWithSignatories struct {
	Form        models.FormsHA        `json:"form"`
	InfoHA      []models.HakAksesInfo `json:"hak_akses_info"`
	Signatories []models.SignatoryHA  `json:"signatories"`
}

func GetSpecAllHAReview(id string) (*FormWithSignatories, error) {
	var formWithSignatories FormWithSignatories

	// Ambil data form
	err := db.Get(&formWithSignatories.Form, `
			SELECT 
					f.form_uuid,
					f.form_number,
					f.form_status,
					d.document_name,
					f.created_by,
					f.created_at,
					f.updated_by,
					f.updated_at,
					f.deleted_by,
					f.deleted_at,
					(f.form_data->>'form_name')::text AS form_name
			FROM
					form_ms f
			LEFT JOIN 
					document_ms d ON f.document_id = d.document_id
			WHERE
					f.form_uuid = $1 AND d.document_code = 'HA' AND f.deleted_at IS NULL
    			AND (f.form_data->>'form_type') = 'Review'
	`, id)
	if err != nil {
		return nil, err
	}

	// Ambil data hak akses info
	err = db.Select(&formWithSignatories.InfoHA, `
			SELECT 
					info_uuid,
					name AS info_name,
					instansi,
					position,
					username,
					password,
					scope
			FROM
					hak_akses_info
			WHERE
					form_id IN (
							SELECT form_id FROM form_ms WHERE form_uuid = $1 AND deleted_at IS NULL
					)
	`, id)
	if err != nil {
		return nil, err
	}

	// Ambil data signatories
	err = db.Select(&formWithSignatories.Signatories, `
	SELECT 
		sign_uuid,
		name AS signatory_name,
		position AS signatory_position,
		role_sign,
		is_sign,
		CASE 
			WHEN sign_img IS NOT NULL AND sign_img != '' 
			THEN '/assets/images/signatures/' || sign_img 
			ELSE ''
		END AS sign_img,
		updated_at
	FROM
		sign_form
	WHERE
		form_id IN (
			SELECT form_id FROM form_ms WHERE form_uuid = $1 AND deleted_at IS NULL
		)
	`, id)
	if err != nil {
		return nil, err
	}

	return &formWithSignatories, nil
}

func GetSpecHakAksesReq(id string) (models.FormsHAReq, error) {
	var specHA models.FormsHAReq

	err := db.Get(&specHA, `SELECT 
	f.form_uuid,
	f.form_status,
	d.document_name,
	f.created_by,
	f.created_at,
	f.updated_by,
	f.updated_at,
	f.deleted_by,
	f.deleted_at,
		(f.form_data->>'form_type')::text AS form_type,
		(f.form_data->>'form_name')::text AS form_name,
    (f.form_data->>'nama_tim')::text AS nama_tim,
    (f.form_data->>'product_manager')::text AS product_manager,
    (f.form_data->>'nama_pengusul')::text AS nama_pengusul,
    (f.form_data->>'tanggal_usul')::text AS tanggal_usul
FROM
	form_ms f
LEFT JOIN 
	document_ms d ON f.document_id = d.document_id
WHERE
	f.form_uuid = $1 AND d.document_code = 'HA' AND f.deleted_at IS NULL
	`, id)
	if err != nil {
		return models.FormsHAReq{}, err
	}

	return specHA, nil

}

func GetSpecHakAkses(id string) (models.FormsHA, error) {
	var specHA models.FormsHA

	err := db.Get(&specHA, `SELECT 
	f.form_uuid,
	f.form_status,
	d.document_name,
	f.created_by,
	f.created_at,
	f.updated_by,
	f.updated_at,
	f.deleted_by,
	f.deleted_at,
	(f.form_data->>'form_name')::text AS form_name
FROM
	form_ms f
LEFT JOIN 
	document_ms d ON f.document_id = d.document_id
WHERE
	f.form_uuid = $1 AND d.document_code = 'HA' AND f.deleted_at IS NULL
    AND (f.form_data->>'form_type') = 'Review'
	`, id)
	if err != nil {
		return models.FormsHA{}, err
	}

	return specHA, nil

}

func UpdateHakAkses(id string, username string, req models.FormRequest, db *sql.DB) error {
	log.Printf("Updating HA with ID: %s, User: %s, Data: %+v", id, username, req)
	currentTime := time.Now()

	// Mulai transaksi database
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}
	defer func() {
		if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			tx.Rollback()
		}
	}()

	// **1. Update Form Data**
	formData, err := json.Marshal(req.FormData)
	if err != nil {
		log.Println("Error marshaling form struct:", err)
		return err
	}

	formStatus := "Draft"
	_, err = tx.Exec(`
		UPDATE form_ms 
		SET form_data = $1, form_status = $2, updated_at = $3, updated_by = $4 
		WHERE form_uuid = $5`,
		formData, formStatus, currentTime, username, id,
	)
	if err != nil {
		log.Println("Error updating form_ms:", err)
		return err
	}

	// **2. Ambil form_id dari form_uuid**
	var formID string
	err = tx.QueryRow("SELECT form_id FROM form_ms WHERE form_uuid = $1", id).Scan(&formID)
	if err == sql.ErrNoRows {
		log.Println("Form ID tidak ditemukan!")
		return fmt.Errorf("Form ID tidak ditemukan!")
	} else if err != nil {
		log.Println("Error getting form_id:", err)
		return err
	}

	// **3. Hapus Hak Akses (Deleted)**
	for _, ha := range req.InfoHA.Deleted {
		_, err = tx.Exec("DELETE FROM hak_akses WHERE ha_uuid = $1", ha.HAUUID)
		if err != nil {
			log.Println("Error deleting hak_akses:", err)
			return err
		}
	}

	// **4. Update Hak Akses (Updated)**
	for _, ha := range req.InfoHA.Updated {
		_, err = tx.Exec(`
			UPDATE hak_akses 
			SET nama_pengguna = $1, ruang_lingkup = $2, jangka_waktu = $3, updated_by = $4, updated_at = $5
			WHERE ha_uuid = $6`,
			ha.NamaPengguna, ha.RuangLingkup, ha.JangkaWaktu, username, currentTime, ha.HAUUID,
		)
		if err != nil {
			log.Println("Error updating hak_akses:", err)
			return err
		}
	}

	// **5. Tambah Hak Akses (Added)**
	for _, ha := range req.InfoHA.Added {
		uuidString := uuid.New().String()
		_, err = tx.Exec(`
			INSERT INTO hak_akses (ha_uuid, form_id, nama_pengguna, ruang_lingkup, jangka_waktu, created_by, created_at, updated_by, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			uuidString, formID, ha.NamaPengguna, ha.RuangLingkup, ha.JangkaWaktu, username, currentTime, username, currentTime,
		)
		if err != nil {
			log.Println("Error inserting hak_akses:", err)
			return err
		}
	}

	// **6. Hapus & Tambah Signatories**
	_, err = tx.Exec("DELETE FROM sign_form WHERE form_id = $1", formID)
	if err != nil {
		log.Println("Error deleting sign_form:", err)
		return err
	}

	// Ambil daftar user yang valid
	personalNames, err := GetAllPersonalName()
	if err != nil {
		log.Println("Error getting personal names:", err)
		return err
	}

	// Insert signatories
	for _, signatory := range req.Signatory {
		uuidString := uuid.New().String()
		var userID string
		for _, personal := range personalNames {
			if personal.PersonalName == signatory.Name {
				userID = personal.UserID
				break
			}
		}

		if userID == "" {
			log.Printf("User ID not found for personal name: %s\n", signatory.Name)
			continue
		}

		_, err = tx.Exec(`
			INSERT INTO sign_form (sign_uuid, form_id, user_id, name, position, role_sign, created_by) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`, 
			uuidString, formID, userID, signatory.Name, signatory.Position, signatory.Role, username,
		)
		if err != nil {
			log.Println("Error inserting sign_form:", err)
			return err
		}
	}

	// **7. Commit transaksi jika semua berhasil**
	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}

	log.Println("UpdateHakAkses successfully completed")
	return nil
}

func UpdateHakAksesReview(id string, username string, req models.FormRequestReview, db *sql.DB) error {
    log.Printf("Updating HA Review with ID: %s, User: %s", id, username)
    currentTime := time.Now()

    // Mulai transaksi
    tx, err := db.Begin()
    if err != nil {
        log.Println("Error starting transaction:", err)
        return err
    }
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            log.Println("Transaction rolled back due to panic:", p)
            err = fmt.Errorf("unexpected error: %v", p)
        } else if err != nil {
            log.Println("Rolling back transaction due to error:", err)
            tx.Rollback()
        }
    }()

    // **1. Update Form Data**
    formData, err := json.Marshal(req.FormData)
    if err != nil {
        log.Println("Error marshaling form struct:", err)
        return err
    }

	formStatus := "Draft"
    result, err := tx.Exec(`UPDATE form_ms SET form_data = $1, form_status = $2, updated_at = $3, updated_by = $4 WHERE form_uuid = $5`,
        formData, formStatus, currentTime, username, id)
    if err != nil {
        log.Println("Error updating form_ms:", err)
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        log.Println("No rows affected, form_uuid not found:", id)
        return fmt.Errorf("form_uuid tidak ditemukan")
    }

    var formID string
    err = tx.QueryRow("SELECT form_id FROM form_ms WHERE form_uuid = $1", id).Scan(&formID)
    if err == sql.ErrNoRows {
        log.Println("Form ID tidak ditemukan!")
        return fmt.Errorf("form ID tidak ditemukan!")
    } else if err != nil {
        log.Println("Error getting form_id:", err)
        return err
    }

    // **2. Hapus & Tambah Signatories**
    _, err = tx.Exec("DELETE FROM sign_form WHERE form_id = $1", formID)
    if err != nil {
        log.Println("Error deleting sign_form records:", err)
        return err
    }

    personalNames, err := GetAllPersonalName()
    if err != nil {
        log.Println("Error getting personal names:", err)
        return err
    }

    for _, signatory := range req.Signatory {
        uuidString := uuid.New().String()
        var userID string
        for _, personal := range personalNames {
            if personal.PersonalName == signatory.Name {
                userID = personal.UserID
                break
            }
        }

        if userID == "" {
            log.Printf("User ID not found for personal name: %s\n", signatory.Name)
            continue
        }

        _, err = tx.Exec(
            `INSERT INTO sign_form (sign_uuid, form_id, user_id, name, position, role_sign, created_by) 
             VALUES ($1, $2, $3, $4, $5, $6, $7)`,
            uuidString, formID, userID, signatory.Name, signatory.Position, signatory.Role, username)
        if err != nil {
            log.Println("Error inserting sign_form:", err)
            return err
        }
    }

    // **3. Hapus, Update, dan Tambah Hak Akses Info**
    for _, info := range req.InfoHAReview.Deleted {
        _, err = tx.Exec("DELETE FROM hak_akses_info WHERE form_id = $1 AND name = $2", formID, info.InfoName)
        if err != nil {
            log.Println("Error deleting hak_akses_info:", err)
            return err
        }
    }

    for _, info := range req.InfoHAReview.Added {
		uuidString := uuid.New().String()
		_, err = tx.Exec(
			`INSERT INTO hak_akses_info (info_uuid, form_id, name, instansi, position, username, password, scope, created_by, created_at, updated_by, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			uuidString, formID, info.InfoName, info.Instansi, info.Position, info.Username, info.Password, info.Scope,
			username, currentTime, username, currentTime)
		if err != nil {
			log.Println("Error inserting hak_akses_info:", err)
			return err
		}
	}
	
	for _, info := range req.InfoHAReview.Updated {
		_, err = tx.Exec(
			`UPDATE hak_akses_info 
			 SET instansi = $1, position = $2, username = $3, password = $4, scope = $5, updated_by = $6, updated_at = $7
			 WHERE form_id = $8 AND name = $9`,
			info.Instansi, info.Position, info.Username, info.Password, info.Scope, username, currentTime, formID, info.InfoName)
		if err != nil {
			log.Println("Error updating hak_akses_info:", err)
			return err
		}
	}

    err = tx.Commit()
    if err != nil {
        log.Println("Error committing transaction:", err)
        return err
    }

    log.Println("UpdateHakAksesReview successfully completed")
    return nil
}

func GetInfoHA(id string) ([]models.HakAksesInfo, error) {
	var infoHA []models.HakAksesInfo
	err := db.Select(&infoHA, `SELECT 
	info_uuid,
	name AS info_name,
	instansi,
	position,
	username,
	password,
	scope
FROM
	hak_akses_info
WHERE
	form_id IN (
		SELECT form_id FROM form_ms WHERE form_uuid = $1 AND deleted_at IS NULL
    AND (f.form_data->>'form_type') = 'Review'
	)
`, id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return infoHA, nil
}

func MyFormHA(userID int) ([]models.FormsHAReq, error) {
	rows, err := db.Query(`SELECT
		f.form_uuid,
		f.form_number,
		f.form_ticket,
		f.form_status,
		d.document_name,
		f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,    
		(f.form_data->>'form_type')::text AS form_type,
    (f.form_data->>'nama_tim')::text AS nama_tim,
    (f.form_data->>'form_name')::text AS form_name,
    (f.form_data->>'product_manager')::text AS product_manager,
    (f.form_data->>'nama_pengusul')::text AS nama_pengusul,
    (f.form_data->>'tanggal_usul')::text AS tanggal_usul
	FROM
		form_ms f
	LEFT JOIN
		document_ms d ON f.document_id = d.document_id
	WHERE
	f.user_id = $1 AND d.document_code = 'HA' AND  f.deleted_at IS NULL
    AND ((f.form_data->>'form_type') = 'Permintaan' OR (f.form_data->>'form_type') = 'Penghapusan')
		ORDER BY f.form_number DESC;
	`, userID)
	var forms []models.FormsHAReq
	//rows, err := db.Query(&forms, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var form models.FormsHAReq
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&form.FormType,
			&form.NamaTim,
			&form.FormName,
			&form.ProductManager,
			&form.NamaPengusul,
			&form.TanggalUsul,
		)

		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}

	return forms, nil
}

func MyFormHAReview(userID int) ([]models.FormsHA, error) {
	rows, err := db.Query(`SELECT
		f.form_uuid,
		f.form_number,
		f.form_ticket,
		f.form_status,
		d.document_name,
		f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
		(f.form_data->>'form_name')::text AS form_name
	FROM
		form_ms f
	LEFT JOIN
		document_ms d ON f.document_id = d.document_id
	WHERE
	f.user_id = $1 AND d.document_code = 'HA' AND  f.deleted_at IS NULL
    AND (f.form_data->>'form_type') = 'Review'
		ORDER BY f.form_number DESC;
	`, userID)
	var forms []models.FormsHA
	//rows, err := db.Query(&forms, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var form models.FormsHA
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&form.FormName,
		)
		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}

	return forms, nil
}

func GetFormsByAdmin() ([]models.FormsHA, error) {
	var forms []models.FormsHA
	query := `SELECT
		f.form_uuid, f.form_status,                                                                                               
		d.document_name,
		f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
		(f.form_data->>'form_name')::text AS form_name
	FROM
		form_ms f
	LEFT JOIN
		document_ms d ON f.document_id = d.document_id
	WHERE
		d.document_code = 'HA' AND f.deleted_at IS NULL
    AND (f.form_data->>'form_type') = 'Review'
		ORDER BY f.form_number DESC;
	`

	// Assuming 'db' is an *sqlx.DB instance
	err := db.Select(&forms, query)
	if err != nil {
		return nil, err
	}

	return forms, nil
}

// menampilkan form berdasar user/ milik dia sendiri
func SignatureUserHA(userID int) ([]models.FormsHA, error) {
	rows, err := db.Query(`SELECT
		f.form_uuid,
		f.form_number,
		f.form_status,
		d.document_name,
		CASE
			WHEN f.is_approve IS NULL THEN 'Menunggu Disetujui'
			WHEN f.is_approve = false THEN 'Tidak Disetujui'
			WHEN f.is_approve = true THEN 'Disetujui'
		END AS ApprovalStatus,
		COALESCE(f.reason, '') AS reason, 
		f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
		(f.form_data->>'form_name')::text AS form_name
		FROM 
		form_ms f
	LEFT JOIN 
		document_ms d ON f.document_id = d.document_id
	LEFT JOIN 
		sign_form sf ON f.form_id = sf.form_id
	WHERE
		sf.user_id = $1 AND d.document_code = 'HA' AND (f.form_data->>'form_type')::text = 'Review' AND f.deleted_at IS NULL
	ORDER BY f.form_number DESC;
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []models.FormsHA

	for rows.Next() {
		var form models.FormsHA
		var formName sql.NullString // Gunakan sql.NullString untuk menangani NULL

		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormStatus,
			&form.DocumentName,
			&form.ApprovalStatus,
			&form.Reason,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&formName, // Simpan ke sql.NullString dulu
		)

		if err != nil {
			return nil, err
		}

		// Konversi NullString ke string biasa
		if formName.Valid {
			form.FormName = formName.String
		} else {
			form.FormName = "" // Beri nilai default jika NULL
		}

		forms = append(forms, form)
	}

	return forms, nil
}

func GetHACode() (models.DocCodeName, error) {
	var documentCode models.DocCodeName

	err := db.Get(&documentCode, "SELECT document_uuid FROM document_ms WHERE document_code = 'HA'")

	if err != nil {
		return models.DocCodeName{}, err
	}
	return documentCode, nil
}

func FormHAByDivision(divisionCode string) ([]models.FormsHAReq, error) {
	var form []models.FormsHAReq

	// Now use the retrieved documentID in the query
	errSelect := db.Select(&form, `
			SELECT 
			f.form_uuid,
			f.form_number,
			f.form_ticket,
			f.form_status,
			f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			d.document_name,
			  (f.form_data->>'form_type')::text AS form_type,
			  (f.form_data->>'form_name')::text AS form_name,
				(f.form_data->>'nama_tim')::text AS nama_tim,
				(f.form_data->>'product_manager')::text AS product_manager,
				(f.form_data->>'nama_pengusul')::text AS nama_pengusul,
				(f.form_data->>'tanggal_usul')::text AS tanggal_usul,
			CASE
				WHEN f.is_approve IS NULL THEN 'Menunggu Disetujui'
				WHEN f.is_approve = false THEN 'Tidak Disetujui'
				WHEN f.is_approve = true THEN 'Disetujui'
			END AS approval_status -- Alias the CASE expression as ApprovalStatus
			FROM 
			form_ms f
		LEFT JOIN 
			document_ms d ON f.document_id = d.document_id
			WHERE
			d.document_code = 'HA' AND f.deleted_at IS NULL AND SPLIT_PART(f.form_number, '/', 2) = $1 
    AND ((f.form_data->>'form_type') = 'Permintaan' OR (f.form_data->>'form_type') = 'Penghapusan')
		ORDER BY f.form_number DESC;
	`, divisionCode)

	if errSelect != nil {
		log.Print(errSelect)
		return nil, errSelect
	}

	if len(form) == 0 {
		return nil, sql.ErrNoRows
	}

	return form, nil
}

func FormHAByDivisionReview(divisionCode string) ([]models.FormsHA, error) {
	var form []models.FormsHA

	// Now use the retrieved documentID in the query
	errSelect := db.Select(&form, `
			SELECT 
			f.form_uuid,
			f.form_number,
			f.form_ticket,
			f.form_status,
			f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			d.document_name,
			(f.form_data->>'form_name')::text AS form_name,
			CASE
				WHEN f.is_approve IS NULL THEN 'Menunggu Disetujui'
				WHEN f.is_approve = false THEN 'Tidak Disetujui'
				WHEN f.is_approve = true THEN 'Disetujui'
			END AS approval_status -- Alias the CASE expression as ApprovalStatus
			FROM 
			form_ms f
		LEFT JOIN 
			document_ms d ON f.document_id = d.document_id
			WHERE
			d.document_code = 'HA' AND f.deleted_at IS NULL AND SPLIT_PART(f.form_number, '/', 2) = $1 
    AND (f.form_data->>'form_type') = 'Review'
		ORDER BY f.form_number DESC;
	`, divisionCode)

	if errSelect != nil {
		log.Print(errSelect)
		return nil, errSelect
	}

	if len(form) == 0 {
		return nil, sql.ErrNoRows
	}

	return form, nil
}

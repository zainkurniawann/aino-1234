package service

import (
	"database/sql"
	"document/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func AddBA(addForm models.Form, ba models.BA, isPublished bool, userID int, username string, divisionCode string, recursionCount int, signatories []models.Signatory) error {
	// Mulai transaksi
	tx, err := db.Beginx()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
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
		return err
	}

	var projectID int64
	err = tx.Get(&projectID, "SELECT project_id FROM project_ms WHERE project_uuid = $1", addForm.ProjectUUID)
	if err != nil {
		log.Println("Error getting project_id:", err)
		return err
	}

	var documentCode string
	err = tx.Get(&documentCode, "SELECT document_code FROM document_ms WHERE document_uuid = $1", addForm.DocumentUUID)
	if err != nil {
		log.Println("Error getting document code:", err)
		return err
	}

	formNumber, err := generateFormNumber(documentID, divisionCode, recursionCount+1)
	if err != nil {
		log.Println("Error generating project form number:", err)
		return err
	}

	// Marshal ITCM struct to JSON
	baJSON, err := json.Marshal(ba)
	if err != nil {
		log.Println("Error marshaling ITCM struct:", err)
		return err
	}

	_, err = tx.NamedExec("INSERT INTO form_ms (form_id, form_uuid, document_id, user_id, project_id, form_number, form_ticket, form_status, form_data, created_by) VALUES (:form_id, :form_uuid, :document_id, :user_id, :project_id, :form_number, :form_ticket, :form_status, :form_data, :created_by)", map[string]interface{}{
		"form_id":     appID,
		"form_uuid":   uuidString,
		"document_id": documentID,
		"user_id":     userID,
		"project_id":  projectID,
		"form_number": formNumber,
		"form_ticket": addForm.FormTicket,
		"form_status": formStatus,
		"form_data":   baJSON, // Convert JSON to string
		"created_by":  username,
	})
	if err != nil {
		return err
	}

	personalNames, err := GetAllPersonalName()
	if err != nil {
		log.Println("Error getting personal names:", err)
		return err
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

		_, err = tx.NamedExec("INSERT INTO sign_form (sign_uuid, form_id, user_id, name, position, role_sign, created_by) VALUES (:sign_uuid, :form_id, :user_id, :name, :position, :role_sign, :created_by)", map[string]interface{}{
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

	// Commit transaksi jika semua query sukses
	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}

	return nil
}

func AddBeritaAcara(addForm models.Form, beritaAcara models.BeritaAcara, isPublished bool, userID int, username string, divisionCode string, recursionCount int, signatories []models.Signatory) error {
	currentTime := time.Now()
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()
	appID := currentTimestamp + int64(uniqueID)
	uuidObj := uuid.New()
	uuidString := uuidObj.String()

	picUUID := uuid.New().String() // Membuat UUID baru untuk PIC

	if beritaAcara.Jenis == "Peminjaman" {
		beritaAcara.PicUUID = picUUID // Menyimpan ke struktur berita acara
	} else {
		beritaAcara.PicUUID = ""
		fmt.Println("Bukan peminjaman")
	}

	// fmt.Println("pic uuid form", beritaAcara.PicUUID)
	// fmt.Println("pic uuid pic", picUUID)

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

	var documentCode string
	err = db.Get(&documentCode, "SELECT document_code FROM document_ms WHERE document_uuid = $1", addForm.DocumentUUID)
	if err != nil {
		log.Println("Error getting document code:", err)
		return err
	}

	formNumber, err := generateFormNumber(documentID, divisionCode, recursionCount+1)
	if err != nil {
		log.Println("Error generating project form number:", err)
		return err
	}

	// Marshal ITCM struct to JSON
	baJSON, err := json.Marshal(beritaAcara)
	if err != nil {
		log.Println("Error marshaling ITCM struct:", err)
		return err
	}

	_, err = db.NamedExec("INSERT INTO form_ms (form_id, form_uuid, form_ticket, document_id, user_id, form_number, form_status, form_data, created_by) VALUES (:form_id, :form_uuid, :form_ticket, :document_id, :user_id, :form_number, :form_status, :form_data, :created_by)", map[string]interface{}{
		"form_id":     appID,
		"form_uuid":   uuidString,
		"form_ticket": "",
		"document_id": documentID,
		"user_id":     userID,
		"form_number": formNumber,
		"form_status": formStatus,
		"form_data":   baJSON, // Convert JSON to string
		"created_by":  username,
	})

	if err != nil {
		return err
	}

	fmt.Println("peliss", beritaAcara.AssetUUID)
	var assetID int64
	err = db.Get(&assetID, "SELECT asset_id FROM assets_ms WHERE asset_uuid = $1", beritaAcara.AssetUUID)
	if err != nil {
		log.Println("Error getting asset_uuid:", err)
		return err
	}

	var assetStatus string
	if beritaAcara.Jenis == "Peminjaman" {
		assetStatus = "Dipinjam"
	} else if beritaAcara.Jenis == "Pengembalian" {
		assetStatus = "Tersedia"
	}
	// Update status asset di database
	_, err = db.NamedExec(`UPDATE assets_ms 
                       SET asset_status = :asset_status, updated_by = :updated_by, updated_at = :updated_at 
                       WHERE asset_id = :asset_id`, map[string]interface{}{
		"asset_status": assetStatus,
		"updated_by":   username,
		"updated_at":   currentTime,
		"asset_id":     assetID,
	})

	if err != nil {
		return err
	}

	if beritaAcara.Jenis == "Peminjaman" {
		var latestPicNumber sql.NullString
		err = db.Get(&latestPicNumber, "SELECT MAX(CAST(REGEXP_REPLACE(pic_description, '[^0-9]', '', 'g') AS INTEGER)) FROM pic_ms WHERE asset_id = $1", assetID)
		if err != nil {
			return fmt.Errorf("Error getting latest form number: %v", err)
		}

		// Initialize picNumber to 1 if latestPicNumber is NULL
		picNumber := 1
		if latestPicNumber.Valid {
			re := regexp.MustCompile(`\d+`) // Regular expression to capture numbers
			match := re.FindString(latestPicNumber.String)

			if match != "" {
				var latestPicNumberInt int
				_, err := fmt.Sscanf(match, "%d", &latestPicNumberInt) // Parsing number
				if err != nil {
					return fmt.Errorf("Error parsing latest pic number: %v", err)
				}
				// Increment the latest pic number
				picNumber = latestPicNumberInt + 1
			} else {
				log.Println("No number found in latest pic description, starting with 1")
			}
		}

		fmt.Println("Latest PIC Number from DB:", latestPicNumber.String)
		fmt.Println("New PIC Number:", picNumber)

		// Ensure PIC description is formatted correctly
		PICWithNumber := fmt.Sprintf("PIC %d", picNumber)
		fmt.Println("New PIC Description for UUID:", PICWithNumber)
		var endedAtValue *string
		if beritaAcara.Ended != "" { // Jika tidak kosong, simpan sebagai pointer
			endedAtValue = &beritaAcara.Ended
		}

		_, err = db.NamedExec("INSERT INTO pic_ms (pic_uuid, asset_id, pic_name, pic_description, created_by, start_at, ended_at) VALUES (:pic_uuid, :asset_id, :pic_name, :pic_description, :created_by, :start_at, :ended_at)", map[string]interface{}{
			"pic_uuid":        picUUID,
			"asset_id":        assetID,
			"pic_name":        beritaAcara.NamaPIC,
			"pic_description": PICWithNumber,
			"created_by":      username,
			"start_at":        beritaAcara.Start,
			"ended_at":        endedAtValue,
		})
		if err != nil {
			return err
		}
	} else {
		log.Println("Jenis Berita Acara bukan Peminjaman, update ended_at berdasarkan PIC yang mengembalikan aset")

		var endedAtValue interface{} // Gunakan interface agar bisa nil atau time.Time

		if beritaAcara.Ended != "" { // Pastikan ada nilai
			parsedTime, err := time.Parse("2006-01-02", beritaAcara.Ended) // Parsing tanggal tanpa jam
			if err != nil {
				return fmt.Errorf("Error parsing ended_at: %v", err)
			}
			endedAtValue = parsedTime // Assign sebagai time.Time
			log.Println("Parsed ended_at:", parsedTime)
		} else {
			endedAtValue = nil // Jika tidak ada nilai, masukkan NULL
		}
		log.Printf("Type of endedAtValue: %T, Value: %v\n", endedAtValue, endedAtValue)

		_, err = db.Exec(`
				UPDATE pic_ms 
				SET ended_at = $1 
				WHERE pic_uuid = (
					SELECT pic_uuid FROM pic_ms 
					WHERE asset_id = $2 
					ORDER BY CAST(REGEXP_REPLACE(pic_description, '[^0-9]', '', 'g') AS INTEGER) DESC 
					LIMIT 1
				)
			`, endedAtValue, assetID)

		if err != nil {
			return fmt.Errorf("Error updating ended_at: %v", err)
		}
	}

	// Ambil daftar semua personal names
	personalNames, err := GetAllPersonalName()
	if err != nil {
		log.Println("Error getting personal names:", err)
		return err
	}

	// Buat map untuk pencarian cepat
	personalNameMap := make(map[string]string)
	for _, personal := range personalNames {
		personalNameMap[personal.PersonalName] = personal.UserID
	}

	for _, signatory := range signatories {
		uuidString := uuid.New().String()

		// Cek apakah name ada di dropdown (personalNameMap)
		userID, exists := personalNameMap[signatory.Name]

		// Jika tidak ada di dropdown, generate UUID baru untuk user_id
		if !exists {
			log.Printf("User ID not found for personal name: %s. Generating new ID...\n", signatory.Name)

			// Generate ID unik berbasis timestamp + UUID
			newID := time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6)
			userID = strconv.FormatInt(newID, 10) // Konversi int64 ke string
		}

		// Insert ke sign_form dengan user_id yang valid
		_, err := db.NamedExec(`
		INSERT INTO sign_form (sign_uuid, form_id, user_id, name, position, role_sign, is_guest, created_by) 
		VALUES (:sign_uuid, :form_id, :user_id, :name, :position, :role_sign, :is_guest, :created_by)`,
			map[string]interface{}{
				"sign_uuid":  uuidString,
				"user_id":    userID,
				"form_id":    appID,
				"name":       signatory.Name,
				"position":   signatory.Position,
				"role_sign":  signatory.Role,
				"is_guest":   !exists, // Jika user tidak ada di DB, is_guest = true
				"created_by": username,
			})
		if err != nil {
			log.Println("Error inserting into sign_form:", err)
			return err
		}
	}
	return nil
}

func AddAsset(addAsset models.Asset, pic []models.Pic, userID int, username string, divisionCode string, recursionCount int, assetImgJSON string) error {
	const maxRecursionCount = 1000

	// Check if the maximum recursion count is reached
	if recursionCount > maxRecursionCount {
		return errors.New("Maximum recursion count exceeded")

	}
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()
	appID := currentTimestamp + int64(uniqueID)
	uuidObj := uuid.New()
	uuidString := uuidObj.String()

	var err error

	var latestCode sql.NullString
	query := fmt.Sprintf(`
		SELECT asset_code 
		FROM assets_ms 
		WHERE asset_type = '%s'
		ORDER BY CAST(REGEXP_REPLACE(asset_code, '[^0-9]', '', 'g') AS INTEGER) DESC 
		LIMIT 1
	`, addAsset.AssetType)

	err = db.Get(&latestCode, query)
	formNumber := 1

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("🔍 Tidak ditemukan asset_code sebelumnya, mulai dari 1.")
		} else {
			return fmt.Errorf("Error getting latest form number: %v", err)
		}
	} else if latestCode.Valid {
		re := regexp.MustCompile(`\D*(\d+)$`) // Cari angka terakhir di akhir string
		matches := re.FindStringSubmatch(latestCode.String)

		if len(matches) == 2 {
			latestFormNumberInt, err := strconv.Atoi(matches[1])
			if err == nil {
				formNumber = latestFormNumberInt + 1
			} else {
				log.Println("⚠️ Gagal parsing angka dari asset_code:", latestCode.String)
			}
		} else {
			log.Println("⚠️ Tidak ditemukan angka valid dalam asset_code:", latestCode.String)
		}
	}

	year := time.Now().Year()
	month := int(time.Now().Month())

	// Mapping kode lokasi
	locationMap := map[string]string{
		"Jakarta":    "JKT",
		"Yogyakarta": "YOG",
		"Surabaya":   "SBY",
		"Bandung":    "BDG",
		"Semarang":   "SMG",
		"Medan":      "MDN",
	}

	// Mapping kode tipe asset
	assetTypeMap := map[string]string{
		"Laptop":  "LP",
		"PC":      "PC",
		"Printer": "PRN",
		"Scanner": "SCN",
	}

	// Dapatkan kode lokasi
	locCode, exists := locationMap[addAsset.Lokasi]
	if !exists {
		if len(addAsset.Lokasi) > 3 {
			locCode = strings.ToUpper(addAsset.Lokasi[:3])
		} else {
			locCode = strings.ToUpper(addAsset.Lokasi)
		}
	}

	// Dapatkan kode tipe asset
	typeCode, exists := assetTypeMap[addAsset.AssetType]
	if !exists {
		if len(addAsset.AssetType) > 2 {
			typeCode = strings.ToUpper(addAsset.AssetType[:2])
		} else {
			typeCode = strings.ToUpper(addAsset.AssetType)
		}
	}

	// Buat kode asset baru
	newAssetCode := fmt.Sprintf("%s%d", typeCode, formNumber)

	// Format final asset code
	assetCode := fmt.Sprintf("HD/%d/%02d/%s/%s", year, month, locCode, newAssetCode)

	fmt.Println("Generated Asset Code:", assetCode)

	if assetImgJSON == "" {
		assetImgJSON = "[]"
	}

	_, err = db.NamedExec("INSERT INTO assets_ms (asset_id, asset_uuid, asset_code, asset_name, serial_number, asset_specification, procurement_date, price, asset_description, system_classification, asset_location, asset_status, asset_img, asset_type, created_by) VALUES (:asset_id, :asset_uuid, :asset_code, :asset_name, :serial_number, :asset_specification, :procurement_date, :price, :asset_description, :system_classification, :asset_location, :asset_status, :asset_img, :asset_type, :created_by)", map[string]interface{}{
		"asset_id":              appID,
		"asset_uuid":            uuidString,
		"asset_code":            assetCode,
		"asset_name":            addAsset.NamaAsset,
		"serial_number":         addAsset.SerialNumber,
		"asset_specification":   addAsset.Spesifikasi,
		"procurement_date":      addAsset.TglPengadaan,
		"price":                 addAsset.Harga,
		"asset_description":     addAsset.Deskripsi,
		"system_classification": addAsset.Klasifikasi,
		"asset_location":        addAsset.Lokasi,
		"asset_status":          addAsset.Status,
		"asset_img":             assetImgJSON,
		"asset_type":            addAsset.AssetType,
		"created_by":            username,
		"created_at":            currentTimestamp,
	})

	if err != nil {
		return err
	}

	for _, pic := range pic {
		uuidString := uuid.New().String()
		var endedAtValue *string
		if pic.Ended != "" {
			endedAtValue = &pic.Ended
		}

		log.Printf("Inserting PIC: %+v\n", map[string]interface{}{
			"pic_uuid":        uuidString,
			"asset_id":        appID,
			"pic_name":        pic.NamaPic,
			"pic_description": pic.Keterangan,
			"created_by":      username,
		})

		_, err := db.NamedExec("INSERT INTO pic_ms (pic_uuid, asset_id, pic_name, pic_description, start_at, ended_at, created_by) VALUES (:pic_uuid, :asset_id, :pic_name, :pic_description, :start_at, :ended_at, :created_by)", map[string]interface{}{
			"pic_uuid":        uuidString,
			"asset_id":        appID,
			"pic_name":        pic.NamaPic,
			"pic_description": pic.Keterangan,
			"start_at":        pic.Start,
			"ended_at":        endedAtValue,
			"created_by":      username,
		})
		if err != nil {
			log.Println("Error inserting PIC:", err)
			return err
		}
	}
	return nil
}

func GetAllFormBA() ([]models.FormsBA, error) {
	rows, err := db.Query(`
		SELECT 
			f.form_uuid,  f.form_number, f.form_ticket, f.form_status,
			d.document_name,
			p.project_name,
			f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			(f.form_data->>'judul')::text AS judul,
			(f.form_data->>'tanggal')::text AS tanggal,
			(f.form_data->>'nama_aplikasi')::text AS nama_aplikasi,
			(f.form_data->>'no_da')::text AS no_da,
			(f.form_data->>'no_itcm')::text AS no_itcm,
			(f.form_data->>'dilakukan_oleh')::text AS dilakukan_oleh,
			(f.form_data->>'didampingi_oleh')::text AS didampingi_oleh
			FROM 
			form_ms f
		LEFT JOIN 
			document_ms d ON f.document_id = d.document_id
		LEFT JOIN 
			project_ms p ON f.project_id = p.project_id
			WHERE
			d.document_code = 'BA' AND f.deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all form data
	var forms []models.FormsBA

	// Iterate through the rows
	for rows.Next() {
		// Scan the row into the Forms struct
		var form models.FormsBA
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.ProjectName,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&form.Judul,
			&form.Tanggal,
			&form.AppName,
			&form.NoDA,
			&form.NoITCM,
			&form.DilakukanOleh,
			&form.DidampingiOleh,
		)
		if err != nil {
			return nil, err
		}

		// Append the form data to the slice
		forms = append(forms, form)
	}
	// Return the forms as JSON response
	return forms, nil
}

// type formBA struct {
// 	Asset        models.Asset    `json:"asset"`
// 	Pic []models.Pic `json:"pic"`
// }

func GetAllFormBAAssets() ([]models.FormsBeritaAcara, error) {
	rows, err := db.Query(`
				SELECT 
			f.form_uuid,  f.form_number, f.form_status,
			d.document_name,
			f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			(f.form_data->>'jenis')::text AS jenis,
			(f.form_data->>'nama_pic')::text AS nama_pic,
			(f.form_data->>'asset_uuid')::text AS asset_uuid,
			(f.form_data->>'kode_asset')::text AS kode_asset,
			(f.form_data->>'jabatan_pic')::text AS jabatan_pic,
			(f.form_data->>'pihak_pertama')::text AS pihak_pertama
			FROM 
			form_ms f
		LEFT JOIN 
			document_ms d ON f.document_id = d.document_id
			WHERE
			d.document_code = 'BA' AND f.deleted_at IS NULL 
		AND  ((f.form_data->>'jenis')::text = 'Peminjaman' OR (f.form_data->>'jenis')::text = 'Pengembalian')
		ORDER BY form_number DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all form data
	var forms []models.FormsBeritaAcara

	// Iterate through the rows
	for rows.Next() {
		// Scan the row into the Forms struct
		var form models.FormsBeritaAcara
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormStatus,
			&form.DocumentName,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&form.Jenis,
			&form.NamaPIC,
			&form.AssetUUID,
			&form.KodeAsset,
			&form.JabatanPIC,
			&form.PihakPertama,
		)
		if err != nil {
			return nil, err
		}

		// Append the form data to the slice
		forms = append(forms, form)
	}
	// Return the forms as JSON response
	return forms, nil
}

func GetAllAssets() ([]models.Asset, error) {
	// 	var formBAWithSignatories formBA

	// 	err := db.Get(&formBAWithSignatories.Asset, `
	//         SELECT
	// 				asset_uuid,
	// 			asset_code
	// 	FROM
	// 	 	assets_ms
	//     `)

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	err = db.Select(&formBAWithSignatories.Pic, `
	// 	SELECT
	// 	 	pic_name, pic_description
	// 	 FROM
	// 	 	pic_ms ORDER BY pic_description ASC

	// `)
	// if err != nil {
	// return nil, err
	// }

	// return &formBAWithSignatories, nil

	// Query untuk mendapatkan semua aset
	rows, err := db.Query(`SELECT 
		asset_id,
		asset_uuid, 
		asset_code, 
		asset_name,
		serial_number,
		asset_specification,
		procurement_date,
		price,
		asset_description,
		system_classification,
		asset_location,
		asset_status,
		asset_type
	FROM 
		assets_ms
WHERE deleted_by IS NULL
ORDER BY CAST(SUBSTRING(asset_code FROM 'LP([0-9]+)$') AS INTEGER) ASC
;
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset

	// Mengambil data aset
	for rows.Next() {
		var asset models.Asset
		var assetID string
		err := rows.Scan(
			&assetID,
			&asset.AssetUUID,
			&asset.Kode,
			&asset.NamaAsset,
			&asset.SerialNumber,
			&asset.Spesifikasi,
			&asset.TglPengadaan,
			&asset.Harga,
			&asset.Deskripsi,
			&asset.Klasifikasi,
			&asset.Lokasi,
			&asset.Status,
			&asset.AssetType,
		)
		if err != nil {
			return nil, err
		}
		asset.AssetID = assetID //kalo gaada code ini nnti eror

		assets = append(assets, asset)
	}

	// Query untuk mendapatkan data PIC
	picRows, err := db.Query(`SELECT 
		asset_id, pic_uuid, pic_name, pic_description 
	FROM 
		pic_ms ORDER BY
    CAST(NULLIF(SUBSTRING(pic_description FROM 5), '') AS INTEGER) ASC;
		`)
	if err != nil {
		return nil, err
	}
	defer picRows.Close()

	// Mengambil data PIC
	picMap := make(map[string][]models.Pic)
	for picRows.Next() {
		var pic models.Pic
		var assetID string
		err := picRows.Scan(
			&assetID,
			&pic.PicUUID,
			&pic.NamaPic,
			&pic.Keterangan,
		)
		if err != nil {
			return nil, err
		}
		// Menambahkan ke picMap berdasarkan asset_id
		picMap[assetID] = append(picMap[assetID], pic)
	}

	// Menggabungkan PIC ke dalam aset
	for i := range assets {
		if pics, exists := picMap[assets[i].AssetID]; exists { // Gunakan asset_id atau yang sesuai
			assets[i].Pic = pics
		}
	}

	return assets, nil
}

func GetAllBAbyUserID(userID int) ([]models.FormsBA, error) {
	rows, err := db.Query(`
		SELECT 
			f.form_uuid,  f.form_number, f.form_ticket, f.form_status,
			d.document_name,
			p.project_name,
			f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			(f.form_data->>'judul')::text AS judul,
			(f.form_data->>'tanggal')::text AS tanggal,
			(f.form_data->>'nama_aplikasi')::text AS nama_aplikasi,
			(f.form_data->>'no_da')::text AS no_da,
			(f.form_data->>'no_itcm')::text AS no_itcm,
			(f.form_data->>'dilakukan_oleh')::text AS dilakukan_oleh,
			(f.form_data->>'didampingi_oleh')::text AS didampingi_oleh
			FROM 
			form_ms f
		LEFT JOIN 
			document_ms d ON f.document_id = d.document_id
		LEFT JOIN 
			project_ms p ON f.project_id = p.project_id
			WHERE
			f.user_id = $1 AND d.document_code = 'BA' AND f.project_id IS NOT NULL AND f.deleted_at IS NULL
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all form data
	var forms []models.FormsBA

	fmt.Println("scan woi")
	// Iterate through the rows
	for rows.Next() {
		// Scan the row into the Forms struct
		var form models.FormsBA
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.ProjectName,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&form.Judul,
			&form.Tanggal,
			&form.AppName,
			&form.NoDA,
			&form.NoITCM,
			&form.DilakukanOleh,
			&form.DidampingiOleh,
		)
		if err != nil {
			return nil, err
		}

		// Append the form data to the slice
		forms = append(forms, form)
	}
	// Return the forms as JSON response
	return forms, nil
}

func GetAllBAbyAdmin() ([]models.FormsBA, error) {
	rows, err := db.Query(`
		SELECT 
			f.form_uuid, f.form_number, f.form_ticket, f.form_status,
			d.document_name,
			p.project_name,
			f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			(f.form_data->>'judul')::text AS judul,
			(f.form_data->>'tanggal')::text AS tanggal,
			(f.form_data->>'nama_aplikasi')::text AS nama_aplikasi,
			(f.form_data->>'no_da')::text AS no_da,
			(f.form_data->>'no_itcm')::text AS no_itcm,
			(f.form_data->>'dilakukan_oleh')::text AS dilakukan_oleh,
			(f.form_data->>'didampingi_oleh')::text AS didampingi_oleh
			FROM 
			form_ms f
		LEFT JOIN 
			document_ms d ON f.document_id = d.document_id
		LEFT JOIN 
			project_ms p ON f.project_id = p.project_id
			WHERE
			d.document_code = 'BA' AND f.deleted_at IS NULL AND f.project_id IS NOT NULL ORDER BY f.form_number DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all form data
	var forms []models.FormsBA

	// Iterate through the rows
	for rows.Next() {
		// Scan the row into the Forms struct
		var form models.FormsBA
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.ProjectName,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&form.Judul,
			&form.Tanggal,
			&form.AppName,
			&form.NoDA,
			&form.NoITCM,
			&form.DilakukanOleh,
			&form.DidampingiOleh,
		)
		if err != nil {
			return nil, err
		}

		// Append the form data to the slice
		forms = append(forms, form)
	}
	// Return the forms as JSON response
	return forms, nil
}

func GetSpecBA(id string) (models.FormsBA, error) {
	var specBA models.FormsBA
	err := db.Get(&specBA, `SELECT 
	f.form_uuid,f.form_number, f.form_ticket, f.form_status,
	d.document_name,
	p.project_name,
	f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
	(f.form_data->>'judul')::text AS judul,
	(f.form_data->>'tanggal')::text AS tanggal,
	(f.form_data->>'nama_aplikasi')::text AS nama_aplikasi,
	(f.form_data->>'no_da')::text AS no_da,
	(f.form_data->>'no_itcm')::text AS no_itcm,
	(f.form_data->>'dilakukan_oleh')::text AS dilakukan_oleh,
	(f.form_data->>'didampingi_oleh')::text AS didampingi_oleh
	FROM 
	form_ms f
LEFT JOIN 
	document_ms d ON f.document_id = d.document_id
LEFT JOIN 
	project_ms p ON f.project_id = p.project_id
	WHERE
	f.form_uuid = $1 AND d.document_code = 'BA'  AND f.deleted_at IS NULL
	`, id)

	if err != nil {
		return models.FormsBA{}, err
	}

	return specBA, nil
}

func GetSpecAsset(id string) (models.Asset, error) {
	var specBA models.Asset
	err := db.QueryRow(`
		SELECT 
			asset_uuid, 
			asset_code, 
			asset_name,
			serial_number,
			asset_specification,
			TO_CHAR(procurement_date, 'YYYY-MM-DD') AS procurement_date,
			price,
			asset_description,
			system_classification,
			asset_location,
			asset_status,
			asset_type
		FROM assets_ms
		WHERE asset_uuid = $1
		AND deleted_at IS NULL
	`, id).Scan(
		&specBA.AssetUUID,
		&specBA.Kode,
		&specBA.NamaAsset,
		&specBA.SerialNumber,
		&specBA.Spesifikasi,
		&specBA.TglPengadaan,
		&specBA.Harga,
		&specBA.Deskripsi,
		&specBA.Klasifikasi,
		&specBA.Lokasi,
		&specBA.Status,
		&specBA.AssetType,
	)
	log.Println("hasil query scan", err)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("🔍 Asset tidak ditemukan dengan UUID:", id)
			return models.Asset{}, nil // Mengembalikan objek kosong tanpa error
		}
		log.Println("❌ Error query GetSpecAsset:", err)
		return models.Asset{}, err
	}

	return specBA, nil
}

type FormBAWithSignatories struct {
	Form        models.FormsBAAll    `json:"form"`
	Signatories []models.SignatoryHA `json:"signatories"`
}

func GetSpecAllBA(id string) (*FormBAWithSignatories, error) {
	var formBAWithSignatories FormBAWithSignatories

	err := db.Get(&formBAWithSignatories.Form, `
        SELECT 
            f.form_uuid, 
            f.form_number, 
            f.form_ticket, 
            f.form_status,
            d.document_name,
            p.project_name,
            f.created_by, 
            f.created_at, 
            f.updated_by, 
            f.updated_at, 
            f.deleted_by, 
            f.deleted_at,
            (f.form_data->>'judul')::text AS judul,
            (f.form_data->>'tanggal')::text AS tanggal,
            (f.form_data->>'nama_aplikasi')::text AS nama_aplikasi,
            (f.form_data->>'no_da')::text AS no_da,
            (f.form_data->>'no_itcm')::text AS no_itcm,
            (f.form_data->>'dilakukan_oleh')::text AS dilakukan_oleh,
            (f.form_data->>'didampingi_oleh')::text AS didampingi_oleh,
        COALESCE((f.form_data->>'keterangan')::text, '') AS keterangan
        FROM
            form_ms f
        LEFT JOIN 
            document_ms d ON f.document_id = d.document_id
        LEFT JOIN 
            project_ms p ON f.project_id = p.project_id
        WHERE
            f.form_uuid = $1 
            AND d.document_code = 'BA'  
            AND f.deleted_at IS NULL
    `, id)

	if err != nil {
		return nil, err
	}

	err = db.Select(&formBAWithSignatories.Signatories, `
        SELECT 
			sign_uuid,
			name AS signatory_name,
			position AS signatory_position,
			role_sign,
			is_guest,
			is_sign,
			CASE
				WHEN sign_img IS NOT NULL AND sign_img != '' THEN CONCAT('/assets/images/signatures/', sign_img)
				ELSE ''
			END AS sign_img,
			updated_at
		FROM sign_form
		WHERE form_id IN (
			SELECT form_id 
			FROM form_ms 
			WHERE form_uuid = $1 
			AND deleted_at IS NULL
            )
    `, id)
	if err != nil {
		return nil, err
	}

	return &formBAWithSignatories, nil
}

type FormBAAssetWithSignatories struct {
	Form        models.FormsBeritaAcaraAsset `json:"form"`
	Signatories []models.SignatoryHA         `json:"signatories"`
}

func GetSpecAllBAAssets(id string) (*FormBAAssetWithSignatories, error) {
	var formBAAssetWithSignatories FormBAAssetWithSignatories
	var BAEvidenceImg []byte

	query := `
    WITH image_data AS (
        SELECT 
            f.form_uuid,
            f.form_data,
            CASE 
                WHEN jsonb_typeof(f.form_data->'image') = 'array' 
                THEN f.form_data->'image' 
                ELSE jsonb_build_array(f.form_data->'image') 
            END AS image_array
        FROM form_ms f
        WHERE f.form_uuid = $1
    )
    SELECT 
        f.form_uuid, f.form_number, f.form_status,
        d.document_name,
        a.asset_name, a.asset_type, a.serial_number, a.asset_specification,
		SPLIT_PART(a.asset_name, ' ', 1) AS merk,
    	SUBSTRING(a.asset_name FROM POSITION(' ' IN a.asset_name) + 1) AS model,
        COALESCE(
            jsonb_agg(DISTINCT CONCAT('/assets/images/pp/', img)) FILTER (WHERE img IS NOT NULL AND img <> ''), 
            '[]'::jsonb
        ) AS image,
        f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
        (f.form_data->>'jenis')::text AS jenis,
        (f.form_data->>'nama_pic')::text AS nama_pic,
        (f.form_data->>'asset_uuid')::text AS asset_uuid,
        (f.form_data->>'kode_asset')::text AS kode_asset,
        COALESCE((f.form_data->>'reason')::text, '') AS reason,
        COALESCE((f.form_data->>'aksesoris')::text, '') AS aksesoris,
        COALESCE((f.form_data->>'kondisi')::text, '') AS kondisi,
        (f.form_data->>'jabatan_pic')::text AS jabatan_pic,
        (f.form_data->>'pihak_pertama')::text AS pihak_pertama,
        (f.form_data->>'jabatan_pihak_pertama')::text AS jabatan_pihak_pertama
    FROM form_ms f
    LEFT JOIN assets_ms a ON (f.form_data->>'asset_uuid') = a.asset_uuid
    LEFT JOIN document_ms d ON f.document_id = d.document_id
    LEFT JOIN image_data id ON id.form_uuid = f.form_uuid
    LEFT JOIN LATERAL jsonb_array_elements_text(id.image_array) AS img ON TRUE
    WHERE f.form_uuid = $1
    AND d.document_code = 'BA' 
    AND f.deleted_at IS NULL
    GROUP BY 
        f.form_uuid, f.form_number, f.form_status, d.document_name, 
        a.asset_name, a.asset_type, a.serial_number, a.asset_specification, 
        f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at, 
        f.form_data;`

	err := db.QueryRow(query, id).Scan(
		&formBAAssetWithSignatories.Form.FormUUID,
		&formBAAssetWithSignatories.Form.FormNumber,
		&formBAAssetWithSignatories.Form.FormStatus,
		&formBAAssetWithSignatories.Form.DocumentName,
		&formBAAssetWithSignatories.Form.NamaAsset,
		&formBAAssetWithSignatories.Form.AssetType,
		&formBAAssetWithSignatories.Form.SerialNumber,
		&formBAAssetWithSignatories.Form.Spesifikasi,
		&formBAAssetWithSignatories.Form.Merk,
		&formBAAssetWithSignatories.Form.Model,
		&BAEvidenceImg,
		&formBAAssetWithSignatories.Form.CreatedBy,
		&formBAAssetWithSignatories.Form.CreatedAt,
		&formBAAssetWithSignatories.Form.UpdatedBy,
		&formBAAssetWithSignatories.Form.UpdatedAt,
		&formBAAssetWithSignatories.Form.DeletedBy,
		&formBAAssetWithSignatories.Form.DeletedAt,
		&formBAAssetWithSignatories.Form.Jenis,
		&formBAAssetWithSignatories.Form.NamaPIC,
		&formBAAssetWithSignatories.Form.AssetUUID,
		&formBAAssetWithSignatories.Form.KodeAsset,
		&formBAAssetWithSignatories.Form.Reason,
		&formBAAssetWithSignatories.Form.Aksesoris,
		&formBAAssetWithSignatories.Form.Kondisi,
		&formBAAssetWithSignatories.Form.JabatanPIC,
		&formBAAssetWithSignatories.Form.PihakPertama,
		&formBAAssetWithSignatories.Form.JabatanPihakPertama,
	)

	if err != nil {
		return nil, err
	}

	// **Convert JSON image array to Go Slice**
	if len(BAEvidenceImg) == 0 || string(BAEvidenceImg) == "null" {
		formBAAssetWithSignatories.Form.Image = []string{}
	} else if err := json.Unmarshal(BAEvidenceImg, &formBAAssetWithSignatories.Form.Image); err != nil {
		log.Println("Error parsing asset_img:", err)
		formBAAssetWithSignatories.Form.Image = []string{}
	}

	// Ambil data Signatories
	err = db.Select(&formBAAssetWithSignatories.Signatories, `
		SELECT 
			sign_uuid,
			name AS signatory_name,
			position AS signatory_position,
			role_sign,
			is_guest,
			is_sign,
			CASE
				WHEN sign_img IS NOT NULL AND sign_img != '' THEN CONCAT('/assets/images/signatures/', sign_img)
				ELSE ''
			END AS sign_img,
			updated_at
		FROM sign_form
		WHERE form_id IN (
			SELECT form_id 
			FROM form_ms 
			WHERE form_uuid = $1 
			AND deleted_at IS NULL
		)
	`, id)

	if err != nil {
		return nil, err
	}

	return &formBAAssetWithSignatories, nil
}

type AssetWithPIC struct {
	Asset models.Asset `json:"asset"`
	PIC   []models.Pic `json:"pic"`
}

func GetSpecAllAsset(id string) (*AssetWithPIC, error) {
	var assetWithPIC AssetWithPIC
	var assetImgJSON []byte

	err := db.QueryRow(`
        SELECT 
			a.asset_uuid, 
			a.asset_code, 
			a.asset_name,
			SPLIT_PART(a.asset_name, ' ', 1) AS merk,
			SUBSTRING(a.asset_name FROM POSITION(' ' IN a.asset_name) + 1) AS model,
			a.serial_number,
			a.asset_specification,
			TO_CHAR(a.procurement_date, 'YYYY-MM-DD') AS procurement_date,
			a.price,
			a.asset_description,
			a.system_classification,
			a.asset_location,
			a.asset_status,
			a.asset_type,
			a.created_at,
			a.created_by,
			COALESCE(
				jsonb_agg(DISTINCT 
					CASE 
						WHEN img IS NOT NULL AND img <> '' 
						THEN CONCAT('/assets/images/asset_img/', img::text) 
					END
				) FILTER (WHERE img IS NOT NULL AND img <> ''), 
				'[]'::jsonb
			) AS asset_img
		FROM 
			assets_ms a
		LEFT JOIN LATERAL jsonb_array_elements_text(a.asset_img::jsonb) AS img ON TRUE
		WHERE
			a.asset_uuid = $1
			AND a.deleted_at IS NULL
		GROUP BY 
			a.asset_uuid, a.asset_code, a.asset_name, merk, model, a.serial_number, 
			a.asset_specification, a.procurement_date, a.price, a.asset_description, 
			a.system_classification, a.asset_location, a.asset_status, a.asset_type, 
			a.created_at, a.created_by;
`, id).Scan(
		&assetWithPIC.Asset.AssetUUID,
		&assetWithPIC.Asset.Kode,
		&assetWithPIC.Asset.NamaAsset,
		&assetWithPIC.Asset.Merk,
		&assetWithPIC.Asset.Model,
		&assetWithPIC.Asset.SerialNumber,
		&assetWithPIC.Asset.Spesifikasi,
		&assetWithPIC.Asset.TglPengadaan,
		&assetWithPIC.Asset.Harga,
		&assetWithPIC.Asset.Deskripsi,
		&assetWithPIC.Asset.Klasifikasi,
		&assetWithPIC.Asset.Lokasi,
		&assetWithPIC.Asset.Status,
		&assetWithPIC.Asset.AssetType,
		&assetWithPIC.Asset.CreatedAt,
		&assetWithPIC.Asset.CreatedBy,
		&assetImgJSON,
	)

	if err != nil {
		return nil, err
	}

	// Parsing JSON ke array string
	if len(assetImgJSON) == 0 || string(assetImgJSON) == "null" {
		assetWithPIC.Asset.AssetImg = []string{}
	} else if err := json.Unmarshal(assetImgJSON, &assetWithPIC.Asset.AssetImg); err != nil {
		log.Println("Error parsing asset_img:", err)
		assetWithPIC.Asset.AssetImg = []string{}
	}

	// **Cek jika hanya berisi "/assets/images/asset_img/" tanpa nama file**
	if len(assetWithPIC.Asset.AssetImg) == 1 && assetWithPIC.Asset.AssetImg[0] == "/assets/images/asset_img/" {
		assetWithPIC.Asset.AssetImg = []string{} // Kosongkan array
	}

	// Ambil data PIC
	rows, err := db.Query(`
    SELECT pic_uuid, pic_name, pic_description, start_at, ended_at 
    FROM pic_ms 
    WHERE asset_id IN (
        SELECT asset_id 
        FROM assets_ms 
        WHERE asset_uuid = $1
    )`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var pic models.Pic
		var startAt, endedAt sql.NullString

		err := rows.Scan(&pic.PicUUID, &pic.NamaPic, &pic.Keterangan, &startAt, &endedAt)
		if err != nil {
			return nil, err
		}

		// Jika NULL, berikan nilai default
		if startAt.Valid {
			pic.Start = startAt.String
		} else {
			pic.Start = ""
		}

		if endedAt.Valid {
			pic.Ended = endedAt.String
		} else {
			pic.Ended = ""
		}

		assetWithPIC.PIC = append(assetWithPIC.PIC, pic)
	}

	return &assetWithPIC, nil
}

func UpdateBA(updateBA models.Form, data models.BA, username string, userID int, isPublished bool, id string, signatories []models.Signatory) (models.Form, error) {
	currentTime := time.Now()
	formStatus := "Draft"
	if isPublished {
		formStatus = "Published"
	}

	var projectID int64
	err := db.Get(&projectID, "SELECT project_id FROM project_ms WHERE project_uuid = $1", updateBA.ProjectUUID)
	if err != nil {
		log.Println("Error getting project_id:", err)
		return models.Form{}, err
	}

	daJSON, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling DampakAnalisa struct:", err)
		return models.Form{}, err
	}
	log.Println("DampakAnalisa JSON:", string(daJSON)) // Periksa hasil marshaling

	_, err = db.NamedExec(`
	UPDATE form_ms 
	SET user_id = :user_id, form_ticket = :form_ticket, project_id = :project_id, 
	    form_status = :form_status, form_data = :form_data, 
	    updated_by = :updated_by, updated_at = :updated_at 
	WHERE form_uuid = :id`,
	map[string]interface{}{
		"user_id":     userID,
		"form_ticket": updateBA.FormTicket,
		"project_id":  projectID, 
		"form_status": formStatus,
		"form_data":   daJSON,
		"updated_by":  username,
		"updated_at":  currentTime,
		"id":          id,
	})
	if err != nil {
		return models.Form{}, err
	}

	var formID string
	err = db.Get(&formID, "SELECT form_id FROM form_ms WHERE form_uuid = $1", id)
	if err != nil {
		log.Println("Error getting form_id:", err)
		return models.Form{}, err
	}

	_, err = db.Exec("DELETE FROM sign_form WHERE form_id = $1", formID)
	if err != nil {
		log.Println("Error deleting sign_form records:", err)
		return models.Form{}, err
	}

	personalNames, err := GetAllPersonalName()
	if err != nil {
		log.Println("Error getting personal names:", err)
		return models.Form{}, err
	}

	for _, signatory := range signatories {
		uuidString := uuid.New().String()

		log.Printf("Processing signatory: %+v\n", signatory)
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

		_, err := db.NamedExec("INSERT INTO sign_form (sign_uuid, form_id, user_id, name, position, role_sign, created_by) VALUES (:sign_uuid, :form_id, :user_id, :name, :position, :role_sign, :created_by)", map[string]interface{}{
			"sign_uuid":  uuidString,
			"user_id":    userID,
			"form_id":    formID, // Adjusted to use documentID
			"name":       signatory.Name,
			"position":   signatory.Position,
			"role_sign":  signatory.Role,
			"created_by": username,
		})
		if err != nil {
			return models.Form{}, err
		}
	}

	return updateBA, nil
}

func UpdateAsset(id string, username string, req models.UpdateImageRequest, db *sql.DB) error {
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

	// **1. Update Asset**
	_, err = tx.Exec(`
		UPDATE assets_ms SET 
			asset_name = $1, serial_number = $2, asset_specification = $3, 
			procurement_date = $4, price = $5, asset_description = $6, 
			system_classification = $7, asset_location = $8, asset_status = $9, 
			asset_type = $10, updated_by = $11, updated_at = $12 
		WHERE asset_uuid = $13`,
		req.Asset.NamaAsset, req.Asset.SerialNumber, req.Asset.Spesifikasi,
		req.Asset.TglPengadaan, req.Asset.Harga, req.Asset.Deskripsi,
		req.Asset.Klasifikasi, req.Asset.Lokasi, req.Asset.Status,
		req.Asset.AssetType, username, time.Now(), id,
	)
	if err != nil {
		log.Println("Error updating asset:", err)
		return err
	}

	// **2. Update gambar (added & deleted)**
	for _, imgBase64 := range req.Asset.Image.Added {
		_, err = tx.Exec(`INSERT INTO assets_ms (asset_img) VALUES ($1)`,
			imgBase64)
		if err != nil {
			log.Println("Error adding asset image:", err)
			return err
		}
	}

	for _, assetimg := range req.Asset.Image.Deleted {
		_, err = tx.Exec(`DELETE FROM assets_ms WHERE asset_id = $1 AND asset_img = $2`, id, assetimg)
		if err != nil {
			log.Println("Error deleting asset image:", err)
			return err
		}
	}

	// **3. Ambil asset_id dari asset_uuid**
	var assetID string
	err = tx.QueryRow("SELECT asset_id FROM assets_ms WHERE asset_uuid = $1", id).Scan(&assetID)
	if err != nil {
		log.Println("Error getting asset_id:", err)
		return err
	}

	// **4. Hapus PIC yang dihapus**
	for _, pic := range req.Pic.Deleted {
		_, err = tx.Exec("DELETE FROM pic_ms WHERE pic_uuid = $1 AND asset_id = $2", pic.PicUUID, assetID)
		if err != nil {
			log.Println("Error deleting PIC:", err)
			return err
		}
	}

	// **5. Update PIC yang diubah**
	for _, pic := range req.Pic.Updated {
		_, err = tx.Exec(`
			UPDATE pic_ms 
			SET pic_name = $1, pic_description = $2, start_at = $3, ended_at = $4, updated_by = $5, updated_at = $6
			WHERE pic_uuid = $7 AND asset_id = $8`,
			pic.NamaPic, pic.Keterangan, pic.Start, pic.Ended, username, time.Now(), pic.PicUUID, assetID,
		)
		if err != nil {
			log.Println("Error updating PIC:", err)
			return err
		}
	}

	// **6. Tambah PIC baru**
	for _, pic := range req.Pic.Added {
		_, err = tx.Exec(`
			INSERT INTO pic_ms (pic_uuid, asset_id, pic_name, pic_description, start_at, ended_at, created_by, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			uuid.New().String(), assetID, pic.NamaPic, pic.Keterangan, pic.Start, pic.Ended, username, time.Now(),
		)
		if err != nil {
			log.Println("Error inserting PIC:", err)
			return err
		}
	}

	// Commit transaksi jika semua berhasil
	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}

	log.Println("UpdateAsset successfully completed")
	return nil
}

func FormBAByDivision(divisionCode string) ([]models.FormsBA, error) {
	var form []models.FormsBA

	// Now use the retrieved documentID in the query
	errSelect := db.Select(&form, `
			SELECT 
			f.form_uuid, f.form_number, f.form_ticket, f.form_status,
			d.document_name,
			p.project_name,
			f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			(f.form_data->>'judul')::text AS judul,
			(f.form_data->>'tanggal')::text AS tanggal,
			(f.form_data->>'nama_aplikasi')::text AS nama_aplikasi,
			(f.form_data->>'no_da')::text AS no_da,
			(f.form_data->>'no_itcm')::text AS no_itcm,
			(f.form_data->>'dilakukan_oleh')::text AS dilakukan_oleh,
			(f.form_data->>'didampingi_oleh')::text AS didampingi_oleh
			FROM 
			form_ms f
		LEFT JOIN 
			document_ms d ON f.document_id = d.document_id
		LEFT JOIN 
			project_ms p ON f.project_id = p.project_id
			WHERE
			d.document_code = 'BA' AND f.deleted_at IS NULL AND f.project_id IS NOT NULL AND SPLIT_PART(f.form_number, '/', 2) = $1
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

func DeleteBeritaAcara(id string, username string, jenis string) error {
	currentTime := time.Now()

	// Mendapatkan assetUUID dari form
	var assetUUID string
	var picUUID string
	err := db.QueryRow("SELECT form_data->>'asset_uuid' AS asset_uuid, form_data->>'pic_uuid' AS pic_uuid FROM form_ms WHERE form_uuid = $1", id).Scan(&assetUUID, &picUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no asset_uuid found for form_uuid %s", id)
		}
		return err
	}

	// Menentukan status baru berdasarkan jenis
	var newStatus string
	if jenis == "Peminjaman" {
		// Jika jenis penghapusan adalah Peminjaman, set status menjadi Dipinjam
		newStatus = "Tersedia"
	} else if jenis == "Pengembalian" {
		// Jika jenis penghapusan adalah Pengembalian, set status menjadi Tersedia
		newStatus = "Dipinjam"
	} else {
		return fmt.Errorf("unexpected jenis: %s", jenis)
	}

	// Menghapus (soft delete) asset di form_ms
	assetResult, err := db.Exec("UPDATE form_ms SET deleted_by = $1, deleted_at = $2 WHERE form_uuid = $3", username, currentTime, id)
	if err != nil {
		return err
	}

	// Cek apakah ada baris yang terpengaruh
	assetRowsAffected, err := assetResult.RowsAffected()
	if err != nil {
		return err
	}
	if assetRowsAffected == 0 {
		return ErrNotFound
	}

	// Mendapatkan assetID dari asset_uuid
	var assetID int
	err = db.QueryRow("SELECT asset_id FROM assets_ms WHERE asset_uuid = $1", assetUUID).Scan(&assetID)
	if err != nil {
		return fmt.Errorf("error retrieving asset_id for asset_uuid %s: %v", assetUUID, err)
	}

	fmt.Println("new", newStatus)

	// Memperbarui status aset
	fmt.Println("asset id", assetID)
	_, err = db.Exec("UPDATE assets_ms SET asset_status = $1 WHERE asset_id = $2", newStatus, assetID)
	if err != nil {
		return err
	}

	fmt.Println("pic uuid", picUUID)
	// Menghapus entri terkait di pic_ms secara permanen
	if jenis == "Peminjaman" {
		picResult, err := db.Exec("DELETE FROM pic_ms WHERE pic_uuid = $1", picUUID)
		if err != nil {
			return err
		}

		// Cek apakah ada baris yang terpengaruh untuk PIC
		formRowsAffected, err := picResult.RowsAffected()
		if err != nil {
			return err
		}
		if formRowsAffected == 0 {
			return ErrNotFound
		}
	}

	return nil
}

func DeleteAsset(id string, username string) error {
	currentTime := time.Now()
	log.Println("DeleteAsset called with id:", id, "and username:", username)

	assetResult, err := db.Exec("UPDATE assets_ms SET deleted_by = $1, deleted_at = $2 WHERE asset_uuid = $3", username, currentTime, id)
	if err != nil {
		log.Println("Error deleting asset:", err)
		return err
	}

	assetRowsAffected, err := assetResult.RowsAffected()
	if err != nil {
		log.Println("Error getting affected rows for asset:", err)
		return err
	}
	log.Println("Asset rows affected:", assetRowsAffected)

	if assetRowsAffected == 0 {
		log.Println("No asset found with id:", id)
		return ErrNotFound
	}

	// formResult, err := db.Exec("UPDATE form_ms SET deleted_by = $1, deleted_at = $2 WHERE form_data->>'asset_uuid' = $3", username, currentTime, id)
	// if err != nil {
	// 	log.Println("Error deleting related form:", err)
	// 	return err
	// }

	// formRowsAffected, err := formResult.RowsAffected()
	// if err != nil {
	// 	log.Println("Error getting affected rows for form:", err)
	// 	return err
	// }
	// log.Println("Form rows affected:", formRowsAffected)

	// if formRowsAffected == 0 {
	// 	log.Println("No form found related to asset id:", id)
	// 	return ErrNotFound
	// }

	return nil
}

// menampilkan formulir sesuai dengan nama signature user tersebut. required signature
func SignatureUserBA(userID int) ([]models.FormsBA, error) {
	rows, err := db.Query(`
		SELECT 
			f.form_uuid, f.form_number, f.form_ticket, f.form_status,
			d.document_name,
			p.project_name,
			f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			(f.form_data->>'judul')::text AS judul,
			(f.form_data->>'tanggal')::text AS tanggal,
			(f.form_data->>'nama_aplikasi')::text AS nama_aplikasi,
			(f.form_data->>'no_da')::text AS no_da,
			(f.form_data->>'no_itcm')::text AS no_itcm,
			(f.form_data->>'dilakukan_oleh')::text AS dilakukan_oleh,
			(f.form_data->>'didampingi_oleh')::text AS didampingi_oleh
		FROM 
			form_ms f
		LEFT JOIN 
			document_ms d ON f.document_id = d.document_id
		LEFT JOIN 
			project_ms p ON f.project_id = p.project_id
		LEFT JOIN 
			sign_form sf ON f.form_id = sf.form_id
		WHERE
			sf.user_id = $1 AND d.document_code = 'BA' AND f.project_id IS NOT NULL AND f.deleted_at IS NULL
		ORDER BY f.form_number DESC;
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []models.FormsBA

	for rows.Next() {
		var form models.FormsBA
		var projectName, judul, tanggal, appName, noDA, noITCM, dilakukanOleh, didampingiOleh sql.NullString // Gunakan sql.NullString

		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&projectName,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&judul, &tanggal, &appName, &noDA, &noITCM, &dilakukanOleh, &didampingiOleh,
		)
		if err != nil {
			return nil, err
		}

		// Handle NULL values
		form.ProjectName = getStringOrDefault(projectName)
		form.Judul = getStringOrDefault(judul)
		form.Tanggal = getStringOrDefault(tanggal)
		form.AppName = getStringOrDefault(appName)
		form.NoDA = getStringOrDefault(noDA)
		form.NoITCM = getStringOrDefault(noITCM)
		form.DilakukanOleh = getStringOrDefault(dilakukanOleh)
		form.DidampingiOleh = getStringOrDefault(didampingiOleh)

		forms = append(forms, form)
	}

	return forms, nil
}

// Fungsi bantu untuk menangani sql.NullString
func getStringOrDefault(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func GetBACode() (models.DocCodeName, error) {
	var documentCode models.DocCodeName

	err := db.Get(&documentCode, "SELECT document_uuid FROM document_ms WHERE document_code = 'BA'")

	if err != nil {
		return models.DocCodeName{}, err
	}
	return documentCode, nil
}

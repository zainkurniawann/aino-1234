package service

// func UpdateAsset(assetData models.Asset, addedPIC, deletedPIC, subjectsPIC []models.Pic, username string, userID int) (models.Asset, error) {
// 	currentTime := time.Now()
// 	var err error

// 	_, err = db.NamedExec(`
// 		UPDATE assets_ms 
// 		SET asset_name = :nama_asset, serial_number = :serial_number, asset_specification = :spesifikasi, 
// 		    procurement_date = :tgl_pengadaan, price = :harga, asset_description = :deskripsi, 
// 		    system_classification = :klasifikasi, asset_location = :lokasi, asset_status = :status, 
// 		    updated_by = :updated_by, updated_at = :updated_at 
// 		WHERE asset_uuid = :id`, map[string]interface{}{
// 		"nama_asset":    assetData.NamaAsset,
// 		"serial_number": assetData.SerialNumber,
// 		"spesifikasi":   assetData.Spesifikasi,
// 		"tgl_pengadaan": assetData.TglPengadaan,
// 		"harga":         assetData.Harga,
// 		"deskripsi":     assetData.Deskripsi,
// 		"klasifikasi":   assetData.Klasifikasi,
// 		"lokasi":        assetData.Lokasi,
// 		"status":        assetData.Status,
// 		"updated_by":    username,
// 		"updated_at":    currentTime,
// 		"id":            assetData.AssetUUID,
// 	})
// 	if err != nil {
// 		return models.Asset{}, err
// 	}

// 	var assetID string
// 	err = db.Get(&assetID, "SELECT asset_id FROM assets_ms WHERE asset_uuid = $1", assetData.AssetUUID)
// 	if err != nil {
// 		return models.Asset{}, err
// 	}

// 	// Hapus PIC yang ada di deletedPIC
// 	for _, pic := range deletedPIC {
// 		_, err := db.Exec("DELETE FROM pic_ms WHERE pic_uuid = $1", pic.PicUUID)
// 		if err != nil {
// 			return models.Asset{}, err
// 		}
// 	}

// 	// Tambah PIC yang ada di addedPIC
// 	for _, pic := range addedPIC {
// 		uuidString := uuid.New().String()
// 		_, err := db.NamedExec("INSERT INTO pic_ms (pic_uuid, asset_id, pic_name, pic_description, created_by) VALUES (:pic_uuid, :asset_id, :pic_name, :pic_description, :created_by)", map[string]interface{}{
// 			"pic_uuid":        uuidString,
// 			"asset_id":        assetID,
// 			"pic_name":        pic.NamaPic,
// 			"pic_description": pic.Keterangan,
// 			"created_by":      username,
// 		})
// 		if err != nil {
// 			return models.Asset{}, err
// 		}
// 	}

// 	return assetData, nil
// }

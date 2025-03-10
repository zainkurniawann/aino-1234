package controller

import (
	"database/sql"
	"document/models"
	"document/service"
	"encoding/base64"
	"encoding/json"
	"fmt"

	// "io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetSignatureForm(c echo.Context) error {
	id := c.Param("id")

	var getAppRole []models.Signatories

	getAppRole, err := service.GetSignatureForm(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Signatory tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			log.Print(err)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, getAppRole)
}

func GetSpecSignatureByID(c echo.Context) error {
	id := c.Param("id")

	var getAppRole models.Signatorie

	getAppRole, err := service.GetSpecSignatureByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Signatory tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			log.Print(err)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, getAppRole)
}

// UpdateSignature updates the signature of a user
func UpdateSignatureGuest(c echo.Context) error {
	id := c.Param("id")

	// Pastikan signature ID ditemukan di database
	previousContent, errGet := service.GetSpecSignatureByID(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal mengupdate signature. Signature tidak ditemukan!",
			Status:  false,
		})
	}

	var editSign models.UpdateSignGuest
	if err := c.Bind(&editSign); err != nil {
		log.Print(err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data invalid!",
			Status:  false,
		})
	}

	// Handle Base64 image
	signImg := editSign.Image
	if signImg == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Gambar tidak ditemukan!",
			Status:  false,
		})
	}

	// Pastikan format Base64 benar
	parts := strings.SplitN(signImg, ",", 2)
	if len(parts) < 2 {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Format gambar tidak valid!",
			Status:  false,
		})
	}

	// Decode the Base64 image
	imgData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Gagal mendekode gambar!",
			Status:  false,
		})
	}

	// Buat folder jika belum ada
	imageFolder := "assets/images/signatures"
	if err := os.MkdirAll(imageFolder, os.ModePerm); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Gagal membuat folder penyimpanan!",
			Status:  false,
		})
	}

	// Generate unique filename
	filename := fmt.Sprintf("signature_%s.png", uuid.New().String())
	dst := filepath.Join(imageFolder, filename)

	// Simpan gambar ke file
	err = os.WriteFile(dst, imgData, 0644)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Gagal menyimpan file!",
			Status:  false,
		})
	}

	// Simpan nama file ke database (relatif path)
	editSign.Image = filename

	// Validasi data sebelum update
	if err = c.Validate(&editSign); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	// Update signature di database
	errService := service.UpdateFormSignatureGuest(editSign, id)
	if errService != nil {
		log.Println("Kesalahan selama pembaruan:", errService)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}

	log.Println(previousContent)
	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Berhasil menambahkan tanda tangan!",
		Status:  true,
	})
}

func UpdateSignature(c echo.Context) error {
	id := c.Param("id")

	// Retrieve previous signature data
	previousContent, errGet := service.GetSpecSignatureByID(id)
	if errGet != nil {
			return c.JSON(http.StatusNotFound, &models.Response{
					Code:    404,
					Message: "Gagal mengupdate signature. Signature tidak ditemukan!",
					Status:  false,
			})
	}

	// Handle JWT token
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    401,
					"message": "Token tidak valid!",
					"status":  false,
			})
	}
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    401,
					"message": "Token tidak valid!",
					"status":  false,
			})
	}

	// Decode JWT claims
	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    401,
					"message": "Token tidak valid!",
					"status":  false,
			})
	}
	userName := c.Get("user_name").(string)

	// Verify user
	userIDFromToken := claims.UserID
	signatory, err := service.GetUserIDSign(id)
	if err != nil || signatory.UserID != userIDFromToken {
			return c.JSON(http.StatusUnauthorized, &models.Response{
					Code:    401,
					Message: "Anda tidak memiliki izin untuk mengupdate tanda tangan ini!",
					Status:  false,
			})
	}

	var editSign models.UpdateSign
	if err := c.Bind(&editSign); err != nil {
			log.Print(err)
			return c.JSON(http.StatusBadRequest, &models.Response{
					Code:    400,
					Message: "Data invalid!",
					Status:  false,
			})
	}
	fmt.Println(editSign)

	// Handle Base64 image
	signImg := editSign.Image
	if signImg == "" {
			return c.JSON(http.StatusBadRequest, &models.Response{
					Code:    400,
					Message: "Gambar tidak ditemukan!",
					Status:  false,
			})
	}

	// Extract the Base64 data
	parts := strings.Split(signImg, ",")
	if len(parts) != 2 {
			return c.JSON(http.StatusBadRequest, &models.Response{
					Code:    400,
					Message: "Data gambar tidak valid!",
					Status:  false,
			})
	}

	// Decode the Base64 image
	imgData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
			return c.JSON(http.StatusBadRequest, &models.Response{
					Code:    400,
					Message: "Gagal mendekode gambar!",
					Status:  false,
			})
	}

	// Generate a unique filename using UUID
	uniqueID := uuid.New().String() // Generates a new UUID
	filename := fmt.Sprintf("signature_%s.png", uniqueID) // Format filename

	// Set folder for saving image
	dst := filepath.Join("assets/images/signatures", filename) // Use the unique filename

	// Save the image to file
	err = os.WriteFile(dst, imgData, 0644) // Using WriteFile to create the file
	if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
					Code:    500,
					Message: "Gagal menyimpan file!",
					Status:  false,
			})
	}

	// Prepare data to update signature
	editSign.Image = filename // Simpan nama file yang unik ke database

	err = c.Validate(&editSign)
	if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, &models.Response{
					Code:    422,
					Message: "Data tidak boleh kosong!",
					Status:  false,
			})
	}

	// Update signature in the service
	errService := service.UpdateFormSignature(editSign, id, userName)
	if errService != nil {
			log.Println("Kesalahan selama pembaruan:", errService)
			return c.JSON(http.StatusInternalServerError, &models.Response{
					Code:    500,
					Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
					Status:  false,
			})
	}

	log.Println(previousContent)
	return c.JSON(http.StatusOK, &models.Response{
			Code:    200,
			Message: "Berhasil menambahkan tanda tangan!",
			Status:  true,
	})
}

func AddApproval(c echo.Context) error {
	id := c.Param("id")
	perviousContent, errGet := service.ShowFormById(id)
	if errGet != nil {
		log.Print(errGet)
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal menambahkan approval. Formulir tidak ditemukan!",
			Status:  false,
		})
	}
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hapus "Bearer " dari tokenString
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	//dekripsi token JWE
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	userName := c.Get("user_name").(string) // Mengambil userUUID dari konteks

	userID := c.Get("user_id").(int)

	fmt.Println("User ID :", userID)
	// Token yang sudah dideskripsi
	fmt.Println("Token yang sudah dideskripsi:", decrypted)

	// User UUID
	fmt.Println("User name:", userName)

	// Lakukan validasi token
	if userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	var editSign models.AddApproval
	if err := c.Bind(&editSign); err != nil {
		log.Print(err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data invalid!",
			Status:  false,
		})
	}

	// Sebelum pemanggilan fungsi AddApproval
	log.Printf("Nilai IsApproval sebelum pemanggilan AddApproval: %v", editSign.IsApproval)
	if !editSign.IsApproval && editSign.Reason == "" {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Alasan harus diisi jika tidak menyetujui.",
			Status:  false,
		})
	}

	err = c.Validate(&editSign)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong",
			Status:  false,
		})
	}

	if err == nil {
		errService := service.AddApproval(editSign, id, userName, userID)
		if errService != nil {
			log.Println("Kesalahan selama pembaruan:", errService)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: errService.Error(),
				Status:  false,
			})
		}

		log.Println(perviousContent)
		return c.JSON(http.StatusOK, &models.Response{
			Code:    200,
			Message: "Berhasil menambahkan approval!",
			Status:  true,
		})
	} else {
		log.Println("Kesalahan sebelum pembaruan:", err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}

}

func AddApprovalDA(c echo.Context) error {
	id := c.Param("id")
	perviousContent, errGet := service.ShowFormById(id)
	if errGet != nil {
		log.Print(errGet)
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal menambahkan approval. Formulir tidak ditemukan!",
			Status:  false,
		})
	}
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hapus "Bearer " dari tokenString
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	//dekripsi token JWE
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	userName := c.Get("user_name").(string) // Mengambil userUUID dari konteks

	userID := c.Get("user_id").(int)

	fmt.Println("User ID :", userID)
	// Token yang sudah dideskripsi
	fmt.Println("Token yang sudah dideskripsi:", decrypted)

	// User UUID
	fmt.Println("User name:", userName)

	// Lakukan validasi token
	if userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	var editSign models.AddApproval
	if err := c.Bind(&editSign); err != nil {
		log.Print(err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data invalid!",
			Status:  false,
		})
	}

	// Sebelum pemanggilan fungsi AddApproval
	log.Printf("Nilai IsApproval sebelum pemanggilan AddApproval: %v", editSign.IsApproval)
	if !editSign.IsApproval && editSign.Reason == "" {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Alasan harus diisi jika tidak menyetujui.",
			Status:  false,
		})
	}

	err = c.Validate(&editSign)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong",
			Status:  false,
		})
	}

	if err == nil {
		errService := service.AddApprovalDA(editSign, id, userName, userID)
		if errService != nil {
			log.Println("Kesalahan selama pembaruan:", errService)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: errService.Error(),
				Status:  false,
			})
		}

		log.Println(perviousContent)
		return c.JSON(http.StatusOK, &models.Response{
			Code:    200,
			Message: "Berhasil menambahkan approval!",
			Status:  true,
		})
	} else {
		log.Println("Kesalahan sebelum pembaruan:", err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}

}

func UpdateSignInfo(c echo.Context) error {
	id := c.Param("id")
	perviousContent, errGet := service.GetSpecSignatureByID(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal mengupdate info signature. Signature tidak ditemukan!",
			Status:  false,
		})
	}
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	userName := c.Get("user_name").(string)
	fmt.Println("Token yang sudah dideskripsi:", decrypted)

	fmt.Println("User name:", userName)

	if userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	var editSign models.UpdateSignForm
	if err := c.Bind(&editSign); err != nil {
		log.Print(err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data invalid!",
			Status:  false,
		})
	}

	err = c.Validate(&editSign)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	// Lakukan pembaruan tanda tangan
	_, err = service.UpdateSignInfo(editSign, id, userName)
	if err != nil {
		log.Println("Kesalahan selama pembaruan:", err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}

	log.Println(perviousContent)
	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Berhasil menambahkan info signature!",
		Status:  true,
	})
}

func AddSignInfo(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hapus "Bearer " dari tokenString
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	//dekripsi token JWE
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	userName := c.Get("user_name").(string) // Mengambil userUUID dari konteks

	// Token yang sudah dideskripsi
	fmt.Println("Token yang sudah dideskripsi:", decrypted)

	// User UUID
	fmt.Println("User name:", userName)

	// Lakukan validasi token
	if userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	var addSignInfo models.AddSignInfo
	if err := c.Bind(&addSignInfo); err != nil {
		log.Println("Gagal melakukan binding data:", err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	if err := c.Validate(&addSignInfo); err != nil {
		log.Println("Gagal melakukan binding data:", err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	// Panggil service untuk menambahkan informasi tanda tangan
	if err := service.AddSignInfo(addSignInfo, userName); err != nil {
		log.Println("Gagal menambahkan informasi tanda tangan:", err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Coba beberapa saat lagi",
			Status:  false,
		})
	}

	return c.JSON(http.StatusCreated, &models.Response{
		Code:    201,
		Message: "Berhasil menambahkan informasi tanda tangan!",
		Status:  true,
	})
}

func DeleteSignInfo(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken" // Ganti dengan kunci yang benar

	// Periksa apakah tokenString tidak kosong
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hapus "Bearer " dari tokenString
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	// Langkah 1: Mendekripsi token JWE
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	userName := c.Get("user_name").(string)
	id := c.Param("id")

	_, errGet := service.GetSpecSignatureByID(id)
	if errGet != nil {
		log.Println("Kesalahan saat penghapusan:", errGet)
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal menghapus signature. Signature tidak ditemukan!",
			Status:  false,
		})
	}

	errService := service.DeleteSignInfo(id, userName)
	if errService != nil {
		log.Println("Kesalahan saat penghapusan:", errService)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})

	}

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Signature berhasil dihapus!",
		Status:  true,
	})

}


func SignatureNotif(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hapus "Bearer " dari tokenString
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	//dekripsi token JWE
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	userID := c.Get("user_id").(int)
	roleCode := c.Get("role_code").(string)

	fmt.Println("User ID :", userID)
	fmt.Println("Role code", roleCode)
	form, err := service.SignatureNotif(userID)
	if err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, form)
}

func ApproveNotif(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hapus "Bearer " dari tokenString
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	//dekripsi token JWE
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	userID := c.Get("user_id").(int)
	roleCode := c.Get("role_code").(string)

	fmt.Println("User ID :", userID)
	fmt.Println("Role code", roleCode)
	form, err := service.ApproveNotif(userID)
	if err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, form)
}
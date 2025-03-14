package controller

import (
	"database/sql"
	"document/database"
	"document/models"
	"document/service"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func AddBA(c echo.Context) error {
	const maxRecursionCount = 1000
	recursionCount := 0 
	var addFormRequest struct {
		IsPublished bool               `json:"isPublished"`
		FormData    models.Form        `json:"formData"`
		BA          models.BA          `json:"data_ba"` 
		Signatory   []models.Signatory `json:"signatories"`
	}

	if err := c.Bind(&addFormRequest); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	if len(addFormRequest.Signatory) == 0 || addFormRequest.BA == (models.BA{}) {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data boleh kosong!",
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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	divisionCode := c.Get("division_code").(string)
	userID := c.Get("user_id").(int) 
	userName := c.Get("user_name").(string)
	addFormRequest.FormData.UserID = userID
	addFormRequest.FormData.Created_by = userName

	if userID == 0 && userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	whitespace := regexp.MustCompile(`^\s`)
	if whitespace.MatchString(addFormRequest.FormData.FormTicket) || whitespace.MatchString(addFormRequest.FormData.FormNumber) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Ticket atau Nomor tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	errVal := c.Validate(&addFormRequest.FormData)
	if errVal == nil {
		addroleErr := service.AddBA(addFormRequest.FormData, addFormRequest.BA, addFormRequest.IsPublished, userID, userName, divisionCode, recursionCount, addFormRequest.Signatory)

		if addroleErr != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Coba beberapa saat lagi",
				Status:  false,
			})
		}

		return c.JSON(http.StatusCreated, &models.Response{
			Code:    201,
			Message: "Berhasil menambahkan formulir berita acara!",
			Status:  true,
		})

	} else {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}
}

func isBeritaAcaraEmpty(ba models.BeritaAcara) bool {
	return ba.AssetUUID == "" || ba.PihakPertama == "" || ba.JabatanPihakPertama == "" || ba.NamaPIC == "" || ba.JabatanPIC == ""
}

func AddBAAsset(c echo.Context) error {
	const maxRecursionCount = 1000
	recursionCount := 0 // Set nilai awal untuk recursionCount
	var addFormRequest struct {
		IsPublished bool               `json:"isPublished"`
		FormData    models.Form        `json:"formData"`
		BeritaAcara models.BeritaAcara `json:"beritaAcara"` // Tambahkan BA ke dalam struct request
		Signatory   []models.Signatory `json:"signatories"`
	}

	if err := c.Bind(&addFormRequest); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	if len(addFormRequest.Signatory) == 0 || isBeritaAcaraEmpty(addFormRequest.BeritaAcara) {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak boleh kosong!",
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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	divisionCode := c.Get("division_code").(string)
	userID := c.Get("user_id").(int) // Mengambil userUUID dari konteks
	userName := c.Get("user_name").(string)
	addFormRequest.FormData.UserID = userID
	addFormRequest.FormData.Created_by = userName

	if userID == 0 && userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	assetImg := addFormRequest.BeritaAcara.Image

	if assetImg == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Gambar tidak ditemukan!",
			Status:  false,
		})
	}

	parts := strings.Split(assetImg, ",")
	if len(parts) != 2 {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data gambar tidak valid!",
			Status:  false,
		})
	}

	imgData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Gagal mendekode gambar!",
			Status:  false,
		})
	}

	uniqueID := uuid.New().String()                      // Generates a new UUID
	filename := fmt.Sprintf("evidence_%s.png", uniqueID) // Format filename

	dst := filepath.Join("assets/images/pp", filename) // Use the unique filename

	err = os.WriteFile(dst, imgData, 0644) // Using WriteFile to create the file
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Gagal menyimpan file!",
			Status:  false,
		})
	}

	addFormRequest.BeritaAcara.Image = filename

	errVal := c.Validate(&addFormRequest.FormData)
	if errVal == nil {
		addroleErr := service.AddBeritaAcara(addFormRequest.FormData, addFormRequest.BeritaAcara, addFormRequest.IsPublished, userID, userName, divisionCode, recursionCount, addFormRequest.Signatory)

		if addroleErr != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Coba beberapa saat lagi",
				Status:  false,
			})
		}

		return c.JSON(http.StatusCreated, &models.Response{
			Code:    201,
			Message: "Berhasil menambahkan formulir berita acara!",
			Status:  true,
		})

	} else {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}
}

func GetBACode(c echo.Context) error {
	documentCode, err := service.GetBACode()
	if err != nil {
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	return c.JSON(http.StatusOK, documentCode)
}

func AddAsset(c echo.Context) error {
	const maxRecursionCount = 1000
	recursionCount := 0

	var addAssetRequest struct {
		AssetData models.Asset `json:"assetData"`
		DataPIC   []models.Pic `json:"data_pic"`
	}

	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan atau tidak valid!",
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

	var claims JwtCustomClaims
	if err := json.Unmarshal([]byte(decrypted), &claims); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	divisionCode, _ := c.Get("division_code").(string)
	userID, _ := c.Get("user_id").(int)
	userName, _ := c.Get("user_name").(string)

	if divisionCode != "GA" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	if err := c.Bind(&addAssetRequest); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	if len(addAssetRequest.AssetData.NamaAsset) == 0 {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Nama asset tidak boleh kosong!",
			Status:  false,
		})
	}

	if len(addAssetRequest.DataPIC) == 0 {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data PIC tidak boleh kosong!",
			Status:  false,
		})
	}

	if addAssetRequest.AssetData.Lokasi == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Lokasi asset tidak boleh kosong!",
			Status:  false,
		})
	}

	parsedDate, err := time.Parse("2006-01-02", addAssetRequest.AssetData.TglPengadaan)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Format tanggal pengadaan tidak valid! Gunakan format YYYY-MM-DD.",
			Status:  false,
		})
	}
	log.Println("parsedDate", parsedDate)

	if addAssetRequest.AssetData.AssetType == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Tipe asset tidak boleh kosong!",
			Status:  false,
		})
	}

	var savedImages []string
	for _, imgBase64 := range addAssetRequest.AssetData.AssetImg {
		imgBase64 = cleanBase64(strings.TrimSpace(imgBase64)) // Bersihkan Base64

		if !isValidBase64(imgBase64) {
			return c.JSON(http.StatusBadRequest, &models.Response{
				Code:    400,
				Message: "Format gambar tidak valid! Pastikan dalam format Base64.",
				Status:  false,
			})
		}

		imgData, err := base64.StdEncoding.DecodeString(imgBase64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &models.Response{
				Code:    400,
				Message: "Gagal mendekode gambar! Periksa format Base64.",
				Status:  false,
			})
		}

		uniqueID := uuid.New().String()
		filename := fmt.Sprintf("asset_%s.png", uniqueID)
		folderPath := "assets/images/asset_img"
		dst := filepath.Join(folderPath, filename)

		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			err = os.MkdirAll(folderPath, os.ModePerm)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, &models.Response{
					Code:    500,
					Message: "Gagal membuat folder penyimpanan!",
					Status:  false,
				})
			}
		}

		err = os.WriteFile(dst, imgData, 0644)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Gagal menyimpan file!",
				Status:  false,
			})
		}

		savedImages = append(savedImages, filename)
	}

	assetImgJSON, err := json.Marshal(savedImages)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Gagal mengolah gambar!",
			Status:  false,
		})
	}

	assetToSave := addAssetRequest.AssetData
	assetToSave.AssetImg = nil

	err = service.AddAsset(assetToSave, addAssetRequest.DataPIC, userID, userName, divisionCode, recursionCount, string(assetImgJSON))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Coba beberapa saat lagi",
			Status:  false,
		})
	}

	return c.JSON(http.StatusCreated, &models.Response{
		Code:    201,
		Message: "Berhasil menambahkan aset!",
		Status:  true,
	})
}









func isValidBase64(str string) bool {
	str = strings.TrimSpace(str)
	if len(str)%4 != 0 {
		return false
	}
	base64Regex := regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)
	if !base64Regex.MatchString(str) {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(str)
	return err == nil
}

func cleanBase64(img string) string {
	parts := strings.Split(img, ",")
	if len(parts) == 2 {
		return parts[1] // Ambil hanya bagian Base64-nya
	}
	return img // Jika tidak ada prefix, langsung kembalikan
}

func GetAllFormBA(c echo.Context) error {
	form, err := service.GetAllFormBA()
	if err != nil {
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, form)
}

func GetAllFormBAAssets(c echo.Context) error {
	form, err := service.GetAllFormBAAssets()
	if err != nil {
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, form)
}

func GetAllAssets(c echo.Context) error {
	form, err := service.GetAllAssets()
	if err != nil {
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, form)
}

func GetSpecBA(c echo.Context) error {
	id := c.Param("id")

	var getDoc models.FormsBA

	getDoc, err := service.GetSpecBA(id)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Code:    404,
				Message: "Formulir berita acara tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, getDoc)
}

func GetSpecAllBA(c echo.Context) error {
	id := c.Param("id")

	formBAWithSignatories, err := service.GetSpecAllBA(id)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Code:    404,
				Message: "Formulir Berita Acara tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	responseData := map[string]interface{}{
		"form":        formBAWithSignatories.Form,
		"signatories": formBAWithSignatories.Signatories,
	}
	return c.JSON(http.StatusOK, responseData)
}

func GetSpecAllBAAssets(c echo.Context) error {
	id := c.Param("id")

	formBAWithSignatories, err := service.GetSpecAllBAAssets(id)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Code:    404,
				Message: "Formulir Berita Acara tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	responseData := map[string]interface{}{
		"form":        formBAWithSignatories.Form,
		"signatories": formBAWithSignatories.Signatories,
	}
	return c.JSON(http.StatusOK, responseData)
}

func GetSpecAllAsset(c echo.Context) error {
	id := c.Param("id")

	assetWithPIC, err := service.GetSpecAllAsset(id)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Code:    404,
				Message: "Asset tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	responseData := map[string]interface{}{
		"asset": assetWithPIC.Asset,
		"pic":   assetWithPIC.PIC,
	}
	return c.JSON(http.StatusOK, responseData)
}

func GetAllFormBAbyUserID(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	userID := c.Get("user_id").(int)

	form, err := service.GetAllBAbyUserID(userID)
	if err != nil {
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, form)

}

func GetAllFormBAAdmin(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	form, err := service.GetAllBAbyAdmin()
	if err != nil {
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, form)

}

func UpdateFormBA(c echo.Context) error {

	id := c.Param("id")

	var updateFormRequest struct {
		IsPublished bool               `json:"isPublished"`
		FormData    models.Form        `json:"formData"`
		Signatory   []models.Signatory `json:"signatories"`
		BA          models.BA          `json:"data_ba"`
	}

	if err := c.Bind(&updateFormRequest); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var userID int
	var userName string

	if claims, ok := c.Get("user_id").(int); ok {
		userID = claims
	} else {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	if name, ok := c.Get("user_name").(string); ok {
		userName = name
	} else {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	updateFormRequest.FormData.UserID = userID

	var updatedBy sql.NullString
	if userName != "" {
		updatedBy.String = userName
		updatedBy.Valid = true
	} else {
		updatedBy.Valid = false
	}

	updateFormRequest.FormData.Updated_by = updatedBy


	if userID == 0 && userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	whitespace := regexp.MustCompile(`^\s`)
	if whitespace.MatchString(updateFormRequest.FormData.FormTicket) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Ticket tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	if err := c.Validate(&updateFormRequest.FormData); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	previousContent, errGet := service.GetSpecBA(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal mengupdate formulir. Formulir tidak ditemukan!",
			Status:  false,
		})
	}
	log.Println("previousContent", previousContent)

	_, errService := service.UpdateBA(updateFormRequest.FormData, updateFormRequest.BA, userName, userID, updateFormRequest.IsPublished, id, updateFormRequest.Signatory)
	if errService != nil {
		if errService.Error() == "You are not authorized to update this form" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Anda tidak diizinkan untuk memperbarui formulir ini",
				"status":  false,
			})
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Formulir Berita Acara berhasil diperbarui!",
		Status:  true,
	})
}

func UpdateBeritaAcara(c echo.Context) error {
	id := c.Param("id")

	var updateFormRequest struct {
		IsPublished bool               `json:"isPublished"`
		FormData    models.Form        `json:"formData"`
		Signatory   []models.Signatory `json:"signatories"`
		BA          models.BA          `json:"data_ba"`
	}
	if err := c.Bind(&updateFormRequest); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	var userID int
	var userName string
	if claims, ok := c.Get("user_id").(int); ok {
		userID = claims
	} else {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	if name, ok := c.Get("user_name").(string); ok {
		userName = name
	} else {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	updateFormRequest.FormData.UserID = userID

	var updatedBy sql.NullString
	if userName != "" {
		updatedBy.String = userName
		updatedBy.Valid = true
	} else {
		updatedBy.Valid = false
	}

	updateFormRequest.FormData.Updated_by = updatedBy


	if userID == 0 && userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	whitespace := regexp.MustCompile(`^\s`)
	if whitespace.MatchString(updateFormRequest.FormData.FormTicket) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Ticket tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}
	if err := c.Validate(&updateFormRequest.FormData); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	previousContent, errGet := service.GetSpecBA(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal mengupdate formulir. Formulir tidak ditemukan!",
			Status:  false,
		})
	}
	log.Println("previousContent", previousContent)

	_, errService := service.UpdateBA(updateFormRequest.FormData, updateFormRequest.BA, userName, userID, updateFormRequest.IsPublished, id, updateFormRequest.Signatory)
	if errService != nil {
		if errService.Error() == "You are not authorized to update this form" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Anda tidak diizinkan untuk memperbarui formulir ini",
				"status":  false,
			})
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Formulir Berita Acara berhasil diperbarui!",
		Status:  true,
	})
}

func UpdateAsset(c echo.Context) error {
	id := c.Param("id")

	var updateRequest models.UpdateImageRequest

	if err := c.Bind(&updateRequest); err != nil {
		c.Logger().Error("Error binding request:", err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, &models.Response{
			Code:    401,
			Message: "Token tidak valid atau tidak ditemukan!",
			Status:  false,
		})
	}

	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")
	secretKey := "secretJwToken"

	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		c.Logger().Error("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, &models.Response{
			Code:    401,
			Message: "Token tidak valid!",
			Status:  false,
		})
	}

	var claims JwtCustomClaims
	if err := json.Unmarshal([]byte(decrypted), &claims); err != nil {
		c.Logger().Error("Gagal mengurai klaim token:", err)
		return c.JSON(http.StatusUnauthorized, &models.Response{
			Code:    401,
			Message: "Token tidak valid!",
			Status:  false,
		})
	}

	userName := claims.UserName
	if userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid atau username tidak ditemukan!",
			"status":  false,
		})
	}

	if err := c.Validate(&updateRequest.Asset); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data asset tidak boleh kosong!",
			Status:  false,
		})
	}

	_, errGet := service.GetSpecAsset(id)
	if errGet != nil {
		c.Logger().Error("Asset tidak ditemukan:", errGet)
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Asset tidak ditemukan!",
			Status:  false,
		})
	}

	for _, base64Img := range updateRequest.Asset.Image.Added {
		if !isValidBase64(base64Img) {
			return c.JSON(http.StatusBadRequest, models.Response{
				Code:    400,
				Message: "Format gambar tidak valid!",
				Status:  false,
			})
		}
	}

	errService := service.UpdateAsset(id, userName, updateRequest, database.DB.DB)
	if errService != nil {
		c.Logger().Error("Kesalahan saat update:", errService)

		if errService.Error() == "You are not authorized to update this asset" {
			return c.JSON(http.StatusUnauthorized, &models.Response{
				Code:    401,
				Message: "Anda tidak diizinkan untuk memperbarui asset ini",
				Status:  false,
			})
		}
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal, silakan coba lagi nanti!",
			Status:  false,
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Asset berhasil diperbarui!",
		Status:  true,
	})
}

func FormBAByDivision(c echo.Context) error {

	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "User ID tidak ditemukan!",
			"status":  false,
		})
	}
	log.Println("userID", userID)

	c.Set("division_code", claims.DivisionCode)
	divisionCode, ok := c.Get("division_code").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Division Code tidak ditemukan!",
			"status":  false,
		})
	}


	myform, err := service.FormBAByDivision(divisionCode)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Code:    404,
				Message: "Form tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}
	return c.JSON(http.StatusOK, myform)
}

func DeleteBeritaAcara(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var userName string

	if name, ok := c.Get("user_name").(string); ok {
		userName = name
	} else {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}
	id := c.Param("id")
	perviousContent, errGet := service.GetSpecAllBAAssets(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal menghapus BA. BA tidak ditemukan!",
			Status:  false,
		})
	}

	jenis := perviousContent.Form.BeritaAcara.Jenis

	errService := service.DeleteBeritaAcara(id, userName, jenis)
	if errService != nil {
		if errService.Error() == "You are not authorized to update this form" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Anda tidak diizinkan untuk menghapus asset ini",
				"status":  false,
			})
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Formulir berhasil dihapus!",
		Status:  true,
	})
}

func DeleteAsset(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var userName string

	if name, ok := c.Get("user_name").(string); ok {
		userName = name
	} else {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}
	id := c.Param("id")
	perviousContent, errGet := service.GetSpecAsset(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal menghapus asset. asset tidak ditemukan!",
			Status:  false,
		})
	}
	log.Println("previousContent", perviousContent)

	errService := service.DeleteAsset(id, userName)
	if errService != nil {
		if errService.Error() == "You are not authorized to update this form" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Anda tidak diizinkan untuk menghapus asset ini",
				"status":  false,
			})
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Formulir berhasil dihapus!",
		Status:  true,
	})
}

func SignatureUserBA(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	userID := c.Get("user_id").(int)

	form, err := service.SignatureUserBA(userID)
	if err != nil {
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, form)

}

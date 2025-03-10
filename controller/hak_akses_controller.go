package controller

import (
	"database/sql"
	"document/database"
	"document/models"
	"document/service"
	"encoding/json"
	"fmt"

	// "io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

func AddHA(c echo.Context) error {
	const maxRecursionCount = 1000
	recursionCount := 0 // Set nilai awal untuk recursionCount
	var addFormRequest struct {
		IsPublished bool                  `json:"isPublished"`
		FormData    models.FormHA         `json:"formData"`
		HakAksesReq models.HAReq          `json:"haReq"`
		HARequest   []models.AddInfoHAReq `json:"data_info_ha_req"`
		Signatory   []models.Signatory    `json:"signatories"`
	}

	if err := c.Bind(&addFormRequest); err != nil {
		log.Print(err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	fmt.Println("pelis", addFormRequest.HARequest)

	fmt.Println("c adalah", c)
	fmt.Printf("Received request data: %+v\n", addFormRequest)

	fmt.Println("ini hareq", addFormRequest.HARequest)
	fmt.Println("ini sign", addFormRequest.Signatory)

	if len(addFormRequest.Signatory) == 0 || len(addFormRequest.HARequest) == 0 {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	fmt.Println("Nilai isPublished yang diterima di backend:", addFormRequest.IsPublished)

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
	divisionCode := c.Get("division_code").(string)
	userID := c.Get("user_id").(int) // Mengambil userUUID dari konteks
	userName := c.Get("user_name").(string)
	addFormRequest.FormData.UserID = userID
	addFormRequest.FormData.Created_by = userName
	// addFormRequest.FormData.isProject = false
	// addFormRequest.FormData.projectCode =
	// Token yang sudah dideskripsi
	fmt.Println("Woiiiiiiiiiiiiiii")
	fmt.Println("Token yang sudah dideskripsi:", decrypted)
	fmt.Println("User ID:", userID)
	fmt.Println("User Name:", userName)
	fmt.Println("Division Code:", divisionCode)
	// Lakukan validasi token
	if userID == 0 && userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	// Validasi spasi untuk Code, Name, dan NumberFormat
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
		addroleErr := service.AddHakAkses(addFormRequest.FormData, addFormRequest.HARequest, addFormRequest.HakAksesReq, addFormRequest.IsPublished, userID, divisionCode, recursionCount, userName, addFormRequest.Signatory)
		if addroleErr != nil {
			log.Print(addroleErr)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Coba beberapa saat lagi",
				Status:  false,
			})
		}

		return c.JSON(http.StatusCreated, &models.Response{
			Code:    201,
			Message: "Berhasil menambahkan formulir hak akses!",
			Status:  true,
		})

	} else {
		fmt.Println(errVal)
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}
}

func AddHAReview(c echo.Context) error {
	const maxRecursionCount = 1000
	recursionCount := 0 // Set nilai awal untuk recursionCount
	var addFormRequest struct {
		IsPublished bool               `json:"isPublished"`
		FormData    models.FormHA      `json:"formData"`
		HakAkses    models.HA          `json:"ha"`
		HA          []models.AddInfoHA `json:"data_info_ha"`
		Signatory   []models.Signatory `json:"signatories"`
	}

	if err := c.Bind(&addFormRequest); err != nil {
		log.Print(err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	fmt.Println("woilaa", addFormRequest.HakAkses)

	fmt.Println("c adalah", c)
	fmt.Printf("Received request data: %+v\n", addFormRequest)

	if len(addFormRequest.Signatory) == 0 || len(addFormRequest.HA) == 0 {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	fmt.Println("Nilai isPublished yang diterima di backend:", addFormRequest.IsPublished)

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
	divisionCode := c.Get("division_code").(string)
	userID := c.Get("user_id").(int) // Mengambil userUUID dari konteks
	userName := c.Get("user_name").(string)
	addFormRequest.FormData.UserID = userID
	addFormRequest.FormData.Created_by = userName
	// addFormRequest.FormData.isProject = false
	// addFormRequest.FormData.projectCode =
	// Token yang sudah dideskripsi
	fmt.Println("Woiiiiiiiiiiiiiii")
	fmt.Println("Token yang sudah dideskripsi:", decrypted)
	fmt.Println("User ID:", userID)
	fmt.Println("User Name:", userName)
	fmt.Println("Division Code:", divisionCode)
	// Lakukan validasi token
	if userID == 0 && userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	// Validasi spasi untuk Code, Name, dan NumberFormat
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
		addroleErr := service.AddHakAksesReview(addFormRequest.FormData, addFormRequest.HA, addFormRequest.HakAkses, addFormRequest.IsPublished, userID, divisionCode, recursionCount, userName, addFormRequest.Signatory)
		if addroleErr != nil {
			log.Print(addroleErr)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Coba beberapa saat lagi",
				Status:  false,
			})
		}

		return c.JSON(http.StatusCreated, &models.Response{
			Code:    201,
			Message: "Berhasil menambahkan formulir hak akses!",
			Status:  true,
		})

	} else {
		fmt.Println(errVal)
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}
}

// menampilkan documen code milik ha
func GetHACode(c echo.Context) error {
	documentCode, err := service.GetHACode()
	if err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	return c.JSON(http.StatusOK, documentCode)
}

// menampilkan form tanpa token
func GetAllFormHA(c echo.Context) error {
	form, err := service.GetAllHakAkses()
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

// menampilkan form tanpa token
func GetAllFormHAReview(c echo.Context) error {
	form, err := service.GetAllHakAksesReview()
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

// func GetSpecHA(c echo.Context) error {
// 	id := c.Param("id")

// 	var getDoc models.FormsBA

// 	getDoc, err := service.GetSpecHA(id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Print(err)
// 			response := models.Response{
// 				Code:    404,
// 				Message: "Formulir berita acara tidak ditemukan!",
// 				Status:  false,
// 			}
// 			return c.JSON(http.StatusNotFound, response)
// 		} else {
// 			log.Print(err)
// 			return c.JSON(http.StatusInternalServerError, &models.Response{
// 				Code:    500,
// 				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
// 				Status:  false,
// 			})
// 		}
// 	}

// 	return c.JSON(http.StatusOK, getDoc)
// }

func GetSpecAllHA(c echo.Context) error {
	id := c.Param("id")

	// Ambil data formulir dan signatories
	formWithSignatories, err := service.GetSpecAllHA(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Formulir Hak Akses tidak ditemukan!",
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

	// Siapkan data respons
	responseData := map[string]interface{}{
		"form":         formWithSignatories.Form,
		"data_info_ha": formWithSignatories.InfoHA,
		"signatories":  formWithSignatories.Signatories,
	}

	return c.JSON(http.StatusOK, responseData)
}

func GetSpecAllHAReview(c echo.Context) error {
	id := c.Param("id")

	// Ambil data formulir dan signatories
	formWithSignatories, err := service.GetSpecAllHAReview(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Formulir Hak Akses tidak ditemukan!",
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

	// Siapkan data respons
	responseData := map[string]interface{}{
		"form":         formWithSignatories.Form,
		"data_info_ha": formWithSignatories.InfoHA,
		"signatories":  formWithSignatories.Signatories,
	}

	return c.JSON(http.StatusOK, responseData)
}

func GetSpecHakAkses(c echo.Context) error {
	id := c.Param("id")

	var getDoc models.FormsHA

	getDoc, err := service.GetSpecHakAkses(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Formulir hak akses tidak ditemukan!",
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

	return c.JSON(http.StatusOK, getDoc)
}

func UpdateHakAkses(c echo.Context) error {
	var updateFormRequest models.FormRequest

	// **1. Bind JSON ke struct**
	if err := c.Bind(&updateFormRequest); err != nil {
		log.Println("Error binding data:", err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data invalid atau format tidak sesuai!",
			Status:  false,
		})
	}

	// **2. Ambil ID dari parameter URL**
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "ID Hak Akses tidak ditemukan dalam parameter!",
			Status:  false,
		})
	}

	// **3. Cek apakah Hak Akses ada di database**
	previousContent, errGet := service.GetSpecHakAksesReq(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal mengupdate Hak Akses. Data tidak ditemukan!",
			Status:  false,
		})
	}

	// **4. Ambil token dari header**
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// **5. Cek format token**
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// **6. Dekripsi token**
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		log.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// **7. Parse klaim token**
	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		log.Println("Gagal mengurai klaim token:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// **8. Ambil username dari token**
	userName := claims.UserName
	if userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid atau username tidak ditemukan!",
			"status":  false,
		})
	}

	if updateFormRequest.FormData.NamaTim == "" || updateFormRequest.FormData.NamaPengusul == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Nama Tim dan Nama Pengusul tidak boleh kosong!",
			"status":  false,
		})
	}
	
	for _, signatory := range updateFormRequest.Signatory {
		if signatory.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": "UUID dan Name pada signatories tidak boleh kosong!",
				"status":  false,
			})
		}
	}
	
	// **9. Panggil service untuk update Hak Akses**
	errService := service.UpdateHakAkses(id, userName, updateFormRequest, database.DB.DB)
	if errService != nil {
		log.Println("Error during update:", errService)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan pada server. Silakan coba lagi!",
			Status:  false,
		})
	}

	log.Println("Previous Content:", previousContent)
	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Berhasil mengupdate Hak Akses!",
		Status:  true,
	})
}

func UpdateHakAksesReview(c echo.Context) error {
	var updateFormRequest models.FormRequestReview

	// Proses binding data dari request body
	if err := c.Bind(&updateFormRequest); err != nil {
		log.Println("Error binding data:", err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	// Ambil ID dari parameter URL
	id := c.Param("id")

	// Cek apakah Hak Akses sebelumnya ada
	previousContent, errGet := service.GetSpecHakAkses(id)
	if errGet != nil {
		log.Println("Error getting previous Hak Akses:", errGet)
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Hak akses tidak ditemukan!",
			Status:  false,
		})
	}

	// Verifikasi token
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
		log.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Parsing token JWT
	var claims JwtCustomClaims
	if err := json.Unmarshal([]byte(decrypted), &claims); err != nil {
		log.Println("Gagal mengurai klaim JWT:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Ambil username dari context Echo
	userName, ok := c.Get("user_name").(string)
	if !ok || userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	// Validasi `FormName`
	if updateFormRequest.FormData.FormName == "" {
		log.Println("Validation error: FormName kosong")
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    422,
			Message: "Form name tidak boleh kosong!",
			Status:  false,
		})
	}

	// Validasi `InfoHA` tidak boleh kosong
	if len(updateFormRequest.InfoHAReview.Added) == 0 && len(updateFormRequest.InfoHAReview.Updated) == 0 && len(updateFormRequest.InfoHAReview.InfoHAReview) == 0 && len(updateFormRequest.InfoHAReview.Deleted) == 0 {
		log.Println("Validation error: Hak Akses Info kosong")
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    422,
			Message: "Data Hak Akses Info tidak boleh kosong!",
			Status:  false,
		})
	}

	// Panggil service untuk update
	errService := service.UpdateHakAksesReview(id, userName, updateFormRequest, database.DB.DB)
	if errService != nil {
		log.Println("Error during update:", errService)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}

	log.Println("Previous Content:", previousContent)
	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Berhasil mengupdate formulir Hak Akses!",
		Status:  true,
	})
}

// menampilkan form dari user/ milik dia sendiri
func MyFormsHA(c echo.Context) error {
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
	form, err := service.MyFormHA(userID)
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

// menampilkan form dari user/ milik dia sendiri
func MyFormsHAReview(c echo.Context) error {
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
	form, err := service.MyFormHAReview(userID)
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

// menampilkan form itcm admin
func GetAllFormHAReviewAdmin(c echo.Context) error {
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
	form, err := service.GetFormsByAdmin()
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

func FormHAByDivision(c echo.Context) error {

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

	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "User ID tidak ditemukan!",
			"status":  false,
		})
	}
	fmt.Println("User ID :", userID)

	c.Set("division_code", claims.DivisionCode)
	divisionCode, ok := c.Get("division_code").(string)
	if !ok {
		// fmt.Println("Division Code is not set or invalid type")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Division Code tidak ditemukan!",
			"status":  false,
		})
	}

	fmt.Println("Division Code :", divisionCode)

	myform, err := service.FormHAByDivision(divisionCode)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Form tidak ditemukan!",
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
	return c.JSON(http.StatusOK, myform)
}

func FormHAByDivisionReview(c echo.Context) error {

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

	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "User ID tidak ditemukan!",
			"status":  false,
		})
	}
	fmt.Println("User ID :", userID)

	c.Set("division_code", claims.DivisionCode)
	divisionCode, ok := c.Get("division_code").(string)
	if !ok {
		// fmt.Println("Division Code is not set or invalid type")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Division Code tidak ditemukan!",
			"status":  false,
		})
	}

	fmt.Println("Division Code :", divisionCode)

	myform, err := service.FormHAByDivisionReview(divisionCode)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Form tidak ditemukan!",
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
	return c.JSON(http.StatusOK, myform)
}

// func PublishHA(c echo.Context) error {

// 	const maxRecursionCount = 1000
// 	recursionCount := 0 // Set nilai awal untuk recursionCount
// 	id := c.Param("id")
// 	var updateFormRequest struct {
// 		IsPublished bool               `json:"isPublished"`
// 		FormData    models.FormPublish `json:"formData"`
// 	}

// 	if err := c.Bind(&updateFormRequest); err != nil {
// 		log.Print("error saat binding:", err)
// 		return c.JSON(http.StatusBadRequest, &models.Response{
// 			Code:    400,
// 			Message: "Data tidak valid!",
// 			Status:  false,
// 		})
// 	}

// 	tokenString := c.Request().Header.Get("Authorization")
// 	secretKey := "secretJwToken"

// 	if tokenString == "" {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Token tidak ditemukan!",
// 			"status":  false,
// 		})
// 	}

// 	// Periksa apakah tokenString mengandung "Bearer "
// 	if !strings.HasPrefix(tokenString, "Bearer ") {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Token tidak valid!",
// 			"status":  false,
// 		})
// 	}

// 	// Hapus "Bearer " dari tokenString
// 	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

// 	//dekripsi token JWE
// 	decrypted, err := DecryptJWE(tokenOnly, secretKey)
// 	if err != nil {
// 		fmt.Println("Gagal mendekripsi token:", err)
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Token tidak valid!",
// 			"status":  false,
// 		})
// 	}

// 	var claims JwtCustomClaims
// 	errJ := json.Unmarshal([]byte(decrypted), &claims)
// 	if errJ != nil {
// 		fmt.Println("Gagal mengurai klaim:", errJ)
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Token tidak valid!",
// 			"status":  false,
// 		})
// 	}
// 	var userID int
// 	var userName string
// 	if claims, ok := c.Get("user_id").(int); ok {
// 		userID = claims
// 	} else {
// 		// Jika gagal mengonversi ke int, tangani kesalahan di sini
// 		log.Println("Tidak dapat mengonversi user_id ke int")
// 		return c.JSON(http.StatusBadRequest, &models.Response{
// 			Code:    400,
// 			Message: "Data tidak valid!",
// 			Status:  false,
// 		})
// 	}

// 	if name, ok := c.Get("user_name").(string); ok {
// 		userName = name
// 	} else {
// 		// Jika gagal mendapatkan nama pengguna, tangani kesalahan di sini
// 		log.Println("Tidak dapat mengonversi user_name ke string")
// 		return c.JSON(http.StatusBadRequest, &models.Response{
// 			Code:    400,
// 			Message: "Data tidak valid!",
// 			Status:  false,
// 		})
// 	}

// 	//updateFormRequest.FormData.UserID = userID

// 	divisionCode := c.Get("division_code").(string)
// 	updateFormRequest.FormData.UserID = userID

// 	var updatedBy sql.NullString
// 	if userName != "" {
// 		updatedBy.String = userName
// 		updatedBy.Valid = true
// 	} else {
// 		updatedBy.Valid = false
// 	}

// 	updateFormRequest.FormData.Updated_by = updatedBy

// 	// Token yang sudah dideskripsi
// 	fmt.Println("Token yang sudah dideskripsi:", decrypted)
// 	fmt.Println("User ID:", userID)
// 	fmt.Println("user name: ", userName)
// 	fmt.Println("division code:", divisionCode)

// 	// Lakukan validasi token
// 	if userID == 0 && userName == "" {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Invalid token atau token tidak ditemukan!",
// 			"status":  false,
// 		})
// 	}

// 	// if userID != updateFormRequest.FormData.UserID {
// 	// 	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 	// 		"code":    401,
// 	// 		"message": "Anda tidak diizinkan untuk memperbarui formulir ini",
// 	// 		"status":  false,
// 	// 	})
// 	// }
// 	whitespace := regexp.MustCompile(`^\s`)
// 	if whitespace.MatchString(updateFormRequest.FormData.FormTicket) {
// 		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
// 			Code:    422,
// 			Message: "Ticket tidak boleh dimulai dengan spasi!",
// 			Status:  false,
// 		})
// 	}

// 	if whitespace.MatchString(updateFormRequest.FormData.FormNumber) {
// 		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
// 			Code:    422,
// 			Message: "Name tidak boleh dimulai dengan spasi!",
// 			Status:  false,
// 		})
// 	}

// 	if err := c.Validate(&updateFormRequest.FormData); err != nil {
// 		return c.JSON(http.StatusInternalServerError, &models.Response{
// 			Code:    422,
// 			Message: "Data tidak boleh kosong!",
// 			Status:  false,
// 		})
// 	}

// 	previousContent, errGet := service.ShowFormById(id)
// 	if errGet != nil {
// 		log.Print(errGet)
// 		return c.JSON(http.StatusNotFound, &models.Response{
// 			Code:    404,
// 			Message: "Gagal mengupdate formulir. Formulir tidak ditemukan!",
// 			Status:  false,
// 		})
// 	}
// 	if previousContent.FormStatus == "Published" {
// 		return c.JSON(http.StatusBadRequest, &models.Response{
// 			Code:    400,
// 			Message: "Tidak dapat memperbarui dokumen yang sudah dipublish",
// 			Status:  false,
// 		})
// 	}

// 	_, errService := service.UpdateForm(updateFormRequest.FormData, id, updateFormRequest.IsPublished, userName, userID, divisionCode, recursionCount)
// 	if errService != nil {
// 		log.Println("Kesalahan selama pembaruan:", errService)
// 		if errService.Error() == "You are not authorized to update this form" {
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Anda tidak diizinkan untuk memperbarui formulir ini",
// 				"status":  false,
// 			})
// 		} else {
// 			return c.JSON(http.StatusInternalServerError, &models.Response{
// 				Code:    500,
// 				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
// 				Status:  false,
// 			})
// 		}
// 	}

// 	log.Println(previousContent)
// 	return c.JSON(http.StatusOK, &models.Response{
// 		Code:    200,
// 		Message: "Formulir berhasil diperbarui!",
// 		Status:  true,
// 	})
// }

func SignatureUserHA(c echo.Context) error {
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
	form, err := service.SignatureUserHA(userID)
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

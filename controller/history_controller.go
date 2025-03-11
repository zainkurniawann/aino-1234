package controller

import (
	// "document/models"
	"document/models"
	"document/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	// "time"

	"github.com/labstack/echo/v4"
)

// func GetTimelineHistory(c echo.Context) error {
// 	// Ambil token dan periksa role
// 	tokenString := c.Request().Header.Get("Authorization")
// 	if tokenString == "" {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Token tidak ditemukan!",
// 			"status":  false,
// 		})
// 	}

// 	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Token tidak valid!",
// 			"status":  false,
// 		})
// 	}

// 	var claims JwtCustomClaims
// 	err = json.Unmarshal([]byte(decrypted), &claims)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Token tidak valid!",
// 			"status":  false,
// 		})
// 	}

// 	// Hanya superadmin (SA) dan admin (A) yang boleh mengakses
// 	if claims.RoleCode != "SA" && claims.RoleCode != "A" {
// 		return c.JSON(http.StatusForbidden, map[string]interface{}{
// 			"code":    403,
// 			"message": "Akses ditolak!",
// 			"status":  false,
// 		})
// 	}

// 	// Ambil data history dari service
// 	history, err := service.GetTimelineHistory(db.DB)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"code":    500,
// 			"message": "Gagal mengambil data history",
// 			"status":  false,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, history)
// }

func GetRecentTimelineHistorySuperAdmin(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Ambil data recent timeline dari service
	history, err := service.GetRecentTimelineHistorySuperAdmin(db.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data history terbaru",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, history)
}

func GetRecentTimelineHistoryAdmin(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
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

	// Ambil data recent timeline dari service
	history, err := service.GetRecentTimelineHistoryAdmin(db.DB, divisionCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data history terbaru",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, history)
}

func GetOlderTimelineHistorySuperAdmin(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Ambil query parameter untuk pagination
	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 3 // Default limit
		// Kirim pesan error bahwa limit tidak valid, tetapi tetap lanjutkan eksekusi
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid limit parameter, using default value",
			"status":  false,
			"data": map[string]interface{}{
				"limit": limit, // Nilai default yang digunakan
			},
		})
	}

	offsetStr := c.QueryParam("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
		// Kirim pesan error bahwa offset tidak valid, tetapi tetap lanjutkan eksekusi
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid offset parameter, using default value",
			"status":  false,
			"data": map[string]interface{}{
				"offset": offset, // Nilai default yang digunakan
			},
		})
	}
	
	log.Println("Received limit:", limitStr, "offset:", offsetStr)

	// Lanjutkan dengan logika aplikasi Anda
	// Contoh: Query database dengan limit dan offset yang sudah diatur
	results, err := service.GetOlderTimelineHistorySuperAdmin(db.DB, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data history lama",
			"status":  false,
		})
	}

	// Jika results kosong, return 200 dengan pesan "no other data"
	if len(results) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    200,
			"message": "no other data",
			"status":  true,
			"data": map[string]interface{}{
				"result": []models.TimelineHistory{}, // Data kosong
			},
		})
	}

	// Kembalikan hasil query
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Success",
		"status":  true,
		"data": map[string]interface{}{
			"limit":  limit,  // Nilai limit yang digunakan
			"offset": offset, // Nilai offset yang digunakan
			"result": results,
		},
	})
}

func GetOlderTimelineHistoryAdmin(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Division Code tidak ditemukan!",
			"status":  false,
		})
	}
	fmt.Println("Division Code :", divisionCode)

	// Ambil query parameter untuk pagination
	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 3 // Default limit
		// Kirim pesan error bahwa limit tidak valid, tetapi tetap lanjutkan eksekusi
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid limit parameter, using default value",
			"status":  false,
			"data": map[string]interface{}{
				"limit": limit, // Nilai default yang digunakan
			},
		})
	}

	offsetStr := c.QueryParam("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
		// Kirim pesan error bahwa offset tidak valid, tetapi tetap lanjutkan eksekusi
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid offset parameter, using default value",
			"status":  false,
			"data": map[string]interface{}{
				"offset": offset, // Nilai default yang digunakan
			},
		})
	}

	// Lanjutkan dengan logika aplikasi Anda
	// Contoh: Query database dengan limit dan offset yang sudah diatur
	results, err := service.GetOlderTimelineHistoryAdmin(db.DB, divisionCode, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data history lama",
			"status":  false,
		})
	}

	// Jika results kosong, return 200 dengan pesan "no other data"
	if len(results) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    200,
			"message": "no other data",
			"status":  true,
			"data": map[string]interface{}{
				"result": []models.TimelineHistory{}, // Data kosong
			},
		})
	}

	// Kembalikan hasil query
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Success",
		"status":  true,
		"data": map[string]interface{}{
			"limit":  limit,  // Nilai limit yang digunakan
			"offset": offset, // Nilai offset yang digunakan
			"result": results,
		},
	})
}

func GetDocumentCountPerMonthSuperAdmin(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Ambil tahun dari query parameter, default ke tahun sekarang
	yearParam := c.QueryParam("year")
	year := time.Now().Year()

	if yearParam != "" {
		parsedYear, err := strconv.Atoi(yearParam)
		if err == nil && parsedYear > 0 {
			year = parsedYear
		}
	}

	// Ambil data jumlah dokumen per bulan dari service
	counts, err := service.GetDocumentCountPerMonthSuperAdmin(db.DB, year)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data jumlah dokumen per bulan",
			"status":  false,
		})
	}

	// Kembalikan data dalam response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   counts,
		"status": true,
	})
}

func GetDocumentCountPerMonthAdmin(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
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

	// Ambil tahun dari query parameter, default ke tahun sekarang
	yearParam := c.QueryParam("year")
	year := time.Now().Year()

	if yearParam != "" {
		parsedYear, err := strconv.Atoi(yearParam)
		if err == nil && parsedYear > 0 {
			year = parsedYear
		}
	}

	// Ambil data jumlah dokumen per bulan dari service
	counts, err := service.GetDocumentCountPerMonthAdmin(db.DB, year, divisionCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data jumlah dokumen per bulan",
			"status":  false,
		})
	}

	// Kembalikan data dalam response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   counts,
		"status": true,
	})
}

func GetDocumentStatusCountPerMonthHandlerSuperAdmin(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Ambil bulan dan tahun dari query parameter, default ke bulan dan tahun sekarang
	monthStr := c.QueryParam("month")
	yearStr := c.QueryParam("year")

	now := time.Now()
	month := now.Month()
	year := now.Year()

	if monthStr != "" {
		parsedMonth, err := strconv.Atoi(monthStr)
		if err == nil && parsedMonth >= 1 && parsedMonth <= 12 {
			month = time.Month(parsedMonth)
		} else {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": "Parameter 'month' harus antara 1-12",
				"status":  false,
			})
		}
	}

	if yearStr != "" {
		parsedYear, err := strconv.Atoi(yearStr)
		if err == nil && parsedYear >= 2000 {
			year = parsedYear
		} else {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": "Parameter 'year' tidak valid",
				"status":  false,
			})
		}
	}

	// Ambil data jumlah dokumen berdasarkan status per bulan dari service
	statusCounts, err := service.GetDocumentStatusCountPerMonthSuperAdmin(db.DB, year, int(month))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data jumlah dokumen berdasarkan status per bulan",
			"status":  false,
		})
	}

	// Kembalikan data dalam response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   statusCounts,
		"status": true,
	})
}

func GetDocumentStatusCountPerMonthHandlerAdmin(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
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

	// Ambil bulan dan tahun dari query parameter, default ke bulan dan tahun sekarang
	monthStr := c.QueryParam("month")
	yearStr := c.QueryParam("year")

	now := time.Now()
	month := now.Month()
	year := now.Year()

	if monthStr != "" {
		parsedMonth, err := strconv.Atoi(monthStr)
		if err == nil && parsedMonth >= 1 && parsedMonth <= 12 {
			month = time.Month(parsedMonth)
		} else {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": "Parameter 'month' harus antara 1-12",
				"status":  false,
			})
		}
	}

	if yearStr != "" {
		parsedYear, err := strconv.Atoi(yearStr)
		if err == nil && parsedYear >= 2000 {
			year = parsedYear
		} else {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": "Parameter 'year' tidak valid",
				"status":  false,
			})
		}
	}

	// Ambil data jumlah dokumen berdasarkan status per bulan dari service
	statusCounts, err := service.GetDocumentStatusCountPerMonthAdmin(db.DB, year, int(month), divisionCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data jumlah dokumen berdasarkan status per bulan",
			"status":  false,
		})
	}

	// Kembalikan data dalam response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   statusCounts,
		"status": true,
	})
}

func GetFormCountPerDocumentPerMonthSuperAdmin(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Ambil tahun dari query parameter, default ke tahun sekarang
	yearParam := c.QueryParam("year")
	now := time.Now()
	year := time.Now().Year()
	month := int(now.Month())

	if yearParam != "" {
		parsedYear, err := strconv.Atoi(yearParam)
		if err == nil && parsedYear > 0 {
			year = parsedYear
		}
	}

	// Ambil data jumlah formulir per dokumen per bulan dari service
	counts, err := service.GetFormCountPerDocumentPerMonthSuperAdmin(db.DB, year, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data jumlah formulir per dokumen per bulan",
			"status":  false,
		})
	}

	// Kembalikan data dalam response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   counts,
		"status": true,
	})
}

func GetFormCountPerDocumentPerMonthAdmin(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	c.Set("division_code", claims.DivisionCode)
	divisionCode, ok := c.Get("division_code").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Division Code tidak ditemukan!",
			"status":  false,
		})
	}

	// Ambil tahun dari query parameter, default ke tahun sekarang
	yearParam := c.QueryParam("year")
	now := time.Now()
	year := time.Now().Year()
	month := int(now.Month())

	if yearParam != "" {
		parsedYear, err := strconv.Atoi(yearParam)
		if err == nil && parsedYear > 0 {
			year = parsedYear
		}
	}

	// Ambil data jumlah formulir per dokumen per bulan dari service
	counts, err := service.GetFormCountPerDocumentPerMonthAdmin(db.DB, year, month, divisionCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data jumlah formulir per dokumen per bulan",
			"status":  false,
		})
	}

	// Kembalikan data dalam response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   counts,
		"status": true,
	})
}
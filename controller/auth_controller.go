package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ticketing-system-backend/auth-service/model"
	"github.com/ticketing-system-backend/auth-service/service"
	"github.com/ticketing-system-backend/auth-service/utils"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Token   string      `json:"token,omitempty"`
	User    interface{} `json:"user,omitempty"`
}

func loginHandler(c *gin.Context, dashboardOnly bool) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "format request tidak valid", nil)
		return
	}

	if validationErrors := utils.ValidateAndFormat(&req); validationErrors != nil {
		utils.Respond(c, http.StatusBadRequest, false, "validasi gagal", validationErrors)
		return
	}

	user, err := service.Login(req.Email, req.Password, dashboardOnly)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "gagal membuat token", nil)
		return
	}

	utils.Respond(c, http.StatusOK, true, "login berhasil", gin.H{
		"token": token,
		"user":  mapUserToResponse(user),
	})
}

// DashboardLogin godoc
//
//	@Summary		Login dashboard
//	@Description	Login khusus dashboard untuk semua role KECUALI customer
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			login	body		controller.LoginRequest	true	"Login data"
//	@Success		200		{object}	controller.LoginResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/login/dashboard [post]
func DashboardLogin(c *gin.Context) {
	loginHandler(c, true)
}

// MobileLogin godoc
//
//	@Summary		Login mobile
//	@Description	Login untuk semua role (termasuk customer)
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			login	body		controller.LoginRequest	true	"Login data"
//	@Success		200		{object}	controller.LoginResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/login/mobile [post]
func MobileLogin(c *gin.Context) {
	loginHandler(c, false)
}

// ====================
// Helper & Mapper
// ====================

func mapUserToResponse(u *model.User) gin.H {
	roles := make([]gin.H, 0)
	for _, r := range u.Roles {
		roles = append(roles, gin.H{
			"id":        r.ID,
			"nama":      r.Nama,
			"deskripsi": r.Deskripsi,
			"level":     r.Level,
		})
	}
	return gin.H{
		"id":           u.ID,
		"nama_lengkap": u.NamaLengkap,
		"email":        u.Email,
		"roles":        roles,
	}
}

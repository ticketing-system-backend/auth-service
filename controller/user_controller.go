package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ticketing-system-backend/auth-service/config"
	"github.com/ticketing-system-backend/auth-service/model"
	"github.com/ticketing-system-backend/auth-service/repository"
	"github.com/ticketing-system-backend/auth-service/utils"
)

type CreateUserRequest struct {
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	RoleIDs     []uint `json:"role_ids" binding:"required"`
}

type UpdateUserRequest struct {
	NamaLengkap string  `json:"nama_lengkap" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Password    *string `json:"password"`
	RoleIDs     []uint  `json:"role_ids" binding:"required"`
}

type UserResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data model.User
}

// GetAllUser		godoc
// @Summary     Get all User
// @Tags        User
// @Accept      json
// @Produce     json
// @Success     200  {object}  UserResponse
// @Failure     400  {object}  map[string]interface{}
// @Failure     404  {object}  map[string]interface{}
// @Router      /users [get]
// @Security    BearerAuth
func GetAllUsers(c *gin.Context) {
	users, err := repository.GetAllUser()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "gagal mengambil user", nil)
		return
	}
	utils.Respond(c, http.StatusOK, true, "daftar user berhasil diambil", users)
}

// GetUserById godoc
// @Summary     Get user by ID
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id   path      int  true  "User ID"
// @Success     200  {object}  UserResponse
// @Failure     400  {object}  map[string]interface{}
// @Failure     404  {object}  map[string]interface{}
// @Router      /users/{id} [get]
// @Security    BearerAuth
func GetUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := repository.FindUserById(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "user tidak ditemukan", nil)
		return
	}
	utils.Respond(c, http.StatusOK, true, "data user berhasil diambil", user)
}

//	CreateUser 	godoc
//	@Summary		Create user
//	@Tags				User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		CreateUserRequest	true	"Create user data"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/users [post]
//	@Security		BearerAuth
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "data user tidak valid", nil)
		return
	}

	if validationErrors := utils.ValidateAndFormat(&req); validationErrors != nil {
		utils.Respond(c, http.StatusBadRequest, false, "validasi gagal", validationErrors)
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "gagal hash password", nil)
		return
	}

	// Ambil roles langsung dari DB
	var roles []model.Role
	if err := config.DB.Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "role tidak valid", nil)
		return
	}

	// Buat user
	user := model.User{
		NamaLengkap: req.NamaLengkap,
		Email:       req.Email,
		Password:    hashed,
		Roles:       roles,
	}

	if err := repository.CreateUser(&user); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "gagal menyimpan user", nil)
		return
	}

	// Ambil data user lengkap dari repository
	userData, err := repository.FindUserById(user.ID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "gagal ambil user", nil)
		return
	}

	utils.Respond(c, http.StatusCreated, true, "user berhasil dibuat", userData)
}

//	UpdateUser	godoc
//	@Summary		Update user
//	@Tags				User
//	@Accept			json
//	@Produce		json
// 	@Param      id    path    	int                        		true  "User ID"
//	@Param			user	body			UpdateUserRequest	true	"Update user data"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/users/{id}	[put]
//	@Security		BearerAuth
func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID tidak valid", nil)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "data tidak valid", nil)
		return
	}

	if validationErrors := utils.ValidateAndFormat(&req); validationErrors != nil {
		utils.Respond(c, http.StatusBadRequest, false, "validasi gagal", validationErrors)
	}

	// Ambil roles dari DB
	var roles []model.Role
	if err := config.DB.Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "role tidak valid", nil)
		return
	}

	// Ambil data user yang akan diupdate
	user, err := repository.FindUserById(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "user tidak ditemukan", nil)
		return
	}

	// Update field dari req
	user.NamaLengkap = req.NamaLengkap
	user.Email = req.Email
	user.Roles = roles

	if req.Password != nil && *req.Password != "" {
		hashed, err := utils.HashPassword(*req.Password)
		if err != nil {
			utils.Respond(c, http.StatusInternalServerError, false, "gagal hash password", nil)
			return
		}
		user.Password = hashed
	}

	if err := repository.UpdateUser(user); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "gagal update user", nil)
		return
	}

	utils.Respond(c, http.StatusOK, true, "user berhasil diperbarui", user)
}
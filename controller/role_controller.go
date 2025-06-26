package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ticketing-system-backend/auth-service/model"
	"github.com/ticketing-system-backend/auth-service/repository"
	"github.com/ticketing-system-backend/auth-service/utils"
)

type CreateRoleRequest struct {
	Nama      string `json:"nama" validate:"required"`
	Deskripsi string `json:"deskripsi" validate:"required"`
	Level     string `json:"level" validate:"required"`
}

type UpdateRoleRequest struct {
	Nama      string `json:"nama" validate:"required"`
	Deskripsi string `json:"deskripsi" validate:"required"`
	Level     string `json:"level" validate:"required"`
}

type RoleResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data model.Role
}

// GetAllRole		godoc
// @Summary     Get all Role
// @Tags        Role
// @Success     200  {object}  RoleResponse
// @Router      /roles [get]
// @Security    BearerAuth
func GetAllRoles(c *gin.Context) {
	roles, err := repository.GetAllRole()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "gagal mengambil role", nil)
		return
	}
	utils.Respond(c, http.StatusOK, true, "daftar role berhasil diambil", roles)
}

// GetRoleById	godoc
// @Summary     Get role by ID
// @Tags        Role
// @Param				id	path	int	true	"Role ID"
// @Success     200  {object}  RoleResponse
// @Router      /roles/{id} [get]
// @Security    BearerAuth
func GetRoleById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	role, err := repository.FindRoleById(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "role tidak ditemukan", nil)
		return
	}
	utils.Respond(c, http.StatusOK, true, "data role berhasil diambil", role)
}

// CreateRole		godoc
// @Summary     Create Role
// @Tags        Role
// @Param				role	body	CreateRoleRequest	true	"Create role data"
// @Success     200  {object}  RoleResponse
// @Router      /roles [post]
// @Security    BearerAuth
func CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "data tidak valid", nil)
		return
	}

	if validationErrors := utils.ValidateAndFormat(&req); validationErrors != nil {
		utils.Respond(c, http.StatusBadRequest, false, "validasi gagal", validationErrors)
	}

role := model.Role{
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
		Level:     req.Level,
	}


	if err := repository.CreateRole(&role); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "gagal menyimpan role", nil)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "role berhasil dibuat", role)
}

// UpdateRole		godoc
// @Summary     Update Role
// @Tags        Role
// @Param				id	path	int	true	"Role ID"
// @Param				role	body	UpdateRoleRequest	true	"Update role data"
// @Success     200  {object}  RoleResponse
// @Router      /roles/{id} [put]
// @Security    BearerAuth
func UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID tidak valid", nil)
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "data tidak valid", nil)
		return
	}

	if validationErrors := utils.ValidateAndFormat(&req); validationErrors != nil {
		utils.Respond(c, http.StatusBadRequest, false, "validasi gagal", validationErrors)
		return
	}

	role := model.Role{
		ID: uint(id),
		Nama: req.Nama,
		Deskripsi: req.Deskripsi,
		Level: req.Level,
	}

	if err := repository.UpdateRole(&role); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "gagal update role", nil)
		return
	}
	utils.Respond(c, http.StatusOK, true, "role berhasil diperbarui", role)
}

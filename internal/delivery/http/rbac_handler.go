package http

import (
	"net/http"
	"strconv"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"
	"github.com/gin-gonic/gin"
)

type RBACHandler struct {
	rbacUsecase *usecase.RBACUsecase
}

func NewRBACHandler(rbacUsecase *usecase.RBACUsecase) *RBACHandler {
	return &RBACHandler{rbacUsecase: rbacUsecase}
}

// ListPermissions godoc
// @Summary      Daftar seluruh granular permissions
// @Description  Mengembalikan seluruh kode izin granular yang terdaftar dalam sistem dikelompokkan per modul
// @Tags         Admin - RBAC
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /admin/rbac/permissions [get]
func (h *RBACHandler) ListPermissions(c *gin.Context) {
	perms, err := h.rbacUsecase.ListPermissions(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal memuat daftar permission: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Daftar permission berhasil dimuat", perms, nil)
}

// ListRoles godoc
// @Summary      Daftar roles dan perizinan matriks
// @Description  Mengembalikan daftar seluruh role beserta array permissions yang terhubung
// @Tags         Admin - RBAC
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /admin/rbac/roles [get]
func (h *RBACHandler) ListRoles(c *gin.Context) {
	roles, err := h.rbacUsecase.ListRolesWithPermissions(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal memuat daftar role: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Daftar role berhasil dimuat", roles, nil)
}

// UpdateRolePermissions godoc
// @Summary      Perbarui matriks permission untuk role tertentu
// @Description  Mengatur daftar permission codes yang diberikan kepada role tertentu
// @Tags         Admin - RBAC
// @Accept       json
// @Produce      json
// @Param        id     path      int                                 true  "Role ID"
// @Param        input  body      usecase.UpdateRolePermissionsInput  true  "Daftar kode permissions baru"
// @Security     BearerAuth
// @Success      200    {object}  response.Response
// @Failure      400    {object}  response.Response
// @Failure      500    {object}  response.Response
// @Router       /admin/rbac/roles/{id}/permissions [put]
func (h *RBACHandler) UpdateRolePermissions(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		response.BadRequest(c, "Role ID tidak valid.", nil)
		return
	}

	var input usecase.UpdateRolePermissionsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload permission codes tidak valid: "+err.Error(), nil)
		return
	}

	if err := h.rbacUsecase.UpdateRolePermissions(c.Request.Context(), roleID, input); err != nil {
		response.InternalServerError(c, "Gagal memperbarui perizinan role: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Perizinan role berhasil diperbarui", nil, nil)
}

// ListUsers godoc
// @Summary      Daftar pengguna dan penugasan role
// @Description  Mengembalikan daftar akun pengguna dengan role dan perizinan efektif mereka
// @Tags         Admin - RBAC
// @Produce      json
// @Param        q         query     string  false  "Pencarian nama/email/username/NIP"
// @Param        page      query     int     false  "Nomor halaman (default: 1)"
// @Param        per_page  query     int     false  "Jumlah item per halaman (default: 15)"
// @Security     BearerAuth
// @Success      200       {object}  response.Response
// @Failure      500       {object}  response.Response
// @Router       /admin/rbac/users [get]
func (h *RBACHandler) ListUsers(c *gin.Context) {
	q := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

	users, total, err := h.rbacUsecase.ListUsers(c.Request.Context(), q, page, perPage)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat daftar pengguna: "+err.Error())
		return
	}

	meta := &response.Meta{
		Page:       page,
		PerPage:    perPage,
		TotalItems: int64(total),
		TotalPages: (total + perPage - 1) / perPage,
	}

	response.Success(c, http.StatusOK, "Daftar pengguna berhasil dimuat", users, meta)
}

// AssignUserRoles godoc
// @Summary      Tetapkan roles untuk pengguna
// @Description  Mengatur penugasan role (multi-role) kepada user tertentu
// @Tags         Admin - RBAC
// @Accept       json
// @Produce      json
// @Param        id     path      int                           true  "User ID"
// @Param        input  body      usecase.AssignUserRolesInput  true  "Daftar role IDs yang diberikan"
// @Security     BearerAuth
// @Success      200    {object}  response.Response
// @Failure      400    {object}  response.Response
// @Failure      500    {object}  response.Response
// @Router       /admin/rbac/users/{id}/roles [put]
func (h *RBACHandler) AssignUserRoles(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.BadRequest(c, "User ID tidak valid.", nil)
		return
	}

	var input usecase.AssignUserRolesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Payload role IDs tidak valid: "+err.Error(), nil)
		return
	}

	if err := h.rbacUsecase.AssignUserRoles(c.Request.Context(), userID, input); err != nil {
		response.InternalServerError(c, "Gagal menetapkan role pengguna: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Role pengguna berhasil diperbarui", nil, nil)
}

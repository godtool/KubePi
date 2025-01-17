package role

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/kubepi/service/api/v1/commons"
	"github.com/KubeOperator/kubepi/service/api/v1/session"
	v1Role "github.com/KubeOperator/kubepi/service/model/v1/role"
	"github.com/KubeOperator/kubepi/service/server"
	"github.com/KubeOperator/kubepi/service/service/v1/common"
	"github.com/KubeOperator/kubepi/service/service/v1/role"
	"github.com/KubeOperator/kubepi/service/service/v1/rolebinding"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	roleService        role.Service
	roleBindingService rolebinding.Service
}

func NewHandler() *Handler {
	return &Handler{
		roleService:        role.NewService(),
		roleBindingService: rolebinding.NewService(),
	}
}

func (h *Handler) SearchRoles() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		var conditions commons.SearchConditions
		if err := ctx.ReadJSON(&conditions); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		groups, total, err := h.roleService.Search(pageNum, pageSize, conditions.Conditions, common.DBOptions{})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", pkgV1.Page{Items: groups, Total: total})
	}
}

// List Roles
// @Tags roles
// @Summary List all roles
// @Description List all roles
// @Accept  json
// @Produce  json
// @Success 200 {object} []v1Role.Role
// @Security ApiKeyAuth
// @Router /roles [get]
func (h *Handler) ListRoles() iris.Handler {
	return func(ctx *context.Context) {
		roles, err := h.roleService.List(common.DBOptions{})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", roles)
	}
}

// Create Role
// @Tags roles
// @Summary Create role
// @Description Create role
// @Accept  json
// @Produce  json
// @Param request body v1Role.Role true "request"
// @Success 200 {object} v1Role.Role
// @Security ApiKeyAuth
// @Router /roles [post]
func (h *Handler) CreateRole() iris.Handler {
	return func(ctx *context.Context) {
		var req v1Role.Role
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		req.CreatedBy = profile.Name
		if err := h.roleService.Create(&req, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

// Delete Role
// @Tags roles
// @Summary Delete role by name
// @Description Delete role by name
// @Accept  json
// @Produce  json
// @Param name path string true "角色名称"
// @Success 200 {object} v1Role.Role
// @Security ApiKeyAuth
// @Router /roles/{name} [delete]
func (h *Handler) DeleteRole() iris.Handler {
	return func(ctx *context.Context) {
		roleName := ctx.Params().GetString("name")
		tx, _ := server.DB().Begin(true)
		txOptions := common.DBOptions{DB: tx}
		rbs, err := h.roleBindingService.GetRoleBindingsByRoleName(roleName, txOptions)
		if err != nil && !errors.As(err, &storm.ErrNotFound) {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		for i := range rbs {
			if err := h.roleBindingService.Delete(rbs[i].Name, txOptions); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		if err := h.roleService.Delete(roleName, txOptions); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		_ = tx.Commit()
	}
}

// Update Role
// @Tags roles
// @Summary Update role by name
// @Description Update role by name
// @Accept  json
// @Produce  json
// @Param name path string true "角色名称"
// @Success 200 {object} v1Role.Role
// @Security ApiKeyAuth
// @Router /roles/{name} [put]
func (h *Handler) UpdateRole() iris.Handler {
	return func(ctx *context.Context) {
		roleName := ctx.Params().GetString("name")
		if roleName == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("invalid resource name %s", roleName))
			return
		}
		var req v1Role.Role
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.roleService.Update(roleName, &req, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

// Get Role
// @Tags roles
// @Summary Get role by name
// @Description Get role by name
// @Accept  json
// @Produce  json
// @Param name path string true "权限名称"
// @Success 200 {object} v1Role.Role
// @Security ApiKeyAuth
// @Router /roles/{name} [get]
func (h *Handler) GetRole() iris.Handler {
	return func(ctx *context.Context) {
		roleName := ctx.Params().GetString("name")
		if roleName == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("invalid resource name %s", roleName))
			return
		}
		r, err := h.roleService.Get(roleName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", r)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/roles")
	sp.Post("/search", handler.SearchRoles())
	sp.Get("/", handler.ListRoles())
	sp.Get("/:name", handler.GetRole())
	sp.Post("/", handler.CreateRole())
	sp.Delete("/:name", handler.DeleteRole())
	sp.Put("/:name", handler.UpdateRole())
}

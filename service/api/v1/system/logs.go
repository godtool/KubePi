package system

import (
	"errors"

	"github.com/KubeOperator/kubepi/service/api/v1/commons"
	"github.com/KubeOperator/kubepi/service/service/v1/common"
	"github.com/KubeOperator/kubepi/service/service/v1/system"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	systemService system.Service
}

func NewHandler() *Handler {
	return &Handler{
		systemService: system.NewService(),
	}
}

func (h *Handler) OperationLogsSearch() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)

		var conditions commons.SearchConditions
		if err := ctx.ReadJSON(&conditions); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		systems, total, err := h.systemService.SearchOperationLogs(pageNum, pageSize, conditions.Conditions, common.DBOptions{})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", pkgV1.Page{Items: systems, Total: total})
	}
}

func (h *Handler) LoginLogsSearch() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)

		var conditions commons.SearchConditions
		if err := ctx.ReadJSON(&conditions); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		systems, total, err := h.systemService.SearchLoginLogs(pageNum, pageSize, conditions.Conditions, common.DBOptions{})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", pkgV1.Page{Items: systems, Total: total})
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/systems")
	sp.Post("/login/logs/search", handler.LoginLogsSearch())
	sp.Post("/operation/logs/search", handler.OperationLogsSearch())
}

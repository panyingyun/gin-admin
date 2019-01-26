package ctl

import (
	"strings"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// User 用户管理
// @Name User
// @Description 用户管理接口
type User struct {
	UserBll *bll.User `inject:""`
}

// Query 查询数据
func (a *User) Query(ctx *context.Context) {
	switch ctx.Query("q") {
	case "page":
		a.QueryPage(ctx)
	default:
		ctx.ResError(errors.NewBadRequestError("未知的查询类型"))
	}
}

// QueryPage 查询分页数据
// @Summary 查询分页数据
// @Param Access-Token header string false "访问令牌"
// @Param current query int true "分页索引" 1
// @Param pageSize query int true "分页大小" 10
// @Param user_name query string false "用户名(模糊查询)"
// @Param real_name query string false "真实姓名(模糊查询)"
// @Param role_id query string false "角色ID"
// @Param status query int false "状态(1:启用 2:停用)"
// @Success 200 []schema.UserPageShow "分页查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 400 option.Interface "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/users?q=page
func (a *User) QueryPage(ctx *context.Context) {
	var params schema.UserQueryParam
	params.UserName = ctx.Query("user_name")
	params.RealName = ctx.Query("real_name")
	params.RoleID = ctx.Query("role_id")
	params.Status = util.S(ctx.Query("status")).Int()

	items, pr, err := a.UserBll.QueryPage(ctx.CContext(), params, ctx.GetPaginationParam())
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResPage(items, pr)
}

// Get 查询指定数据
// Get 查询指定数据
// @Summary 查询指定数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.User
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 404 string "{error:{code:0,message:资源不存在}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/users/{id}
func (a *User) Get(ctx *context.Context) {
	item, err := a.UserBll.Get(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResSuccess(item)
}

// Create 创建数据
// @Summary 创建数据
// @Param Access-Token header string false "访问令牌"
// @Param body body schema.User true
// @Success 200 option.Interface "{record_id:记录ID}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/users
func (a *User) Create(ctx *context.Context) {
	var item schema.User
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	newItem, err := a.UserBll.Create(ctx.CContext(), item)
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResSuccess(context.HTTPNewItem{RecordID: newItem.RecordID})
}

// Update 更新数据
// @Summary 更新数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Param body body schema.User true
// @Success 200 option.Interface "{status:OK}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/users/{id}
func (a *User) Update(ctx *context.Context) {
	var item schema.User
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	err := a.UserBll.Update(ctx.CContext(), ctx.Param("id"), item)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除数据
// @Summary 删除数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Success 200 option.Interface "{status:OK}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router DELETE /api/v1/users/{id}
func (a *User) Delete(ctx *context.Context) {
	err := a.UserBll.Delete(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// DeleteMany 删除多条数据
// @Summary 删除多条数据
// @Param Access-Token header string false "访问令牌"
// @Param batch query string true "记录ID（多个以,分隔）"
// @Success 200 option.Interface "{status:OK}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router DELETE /api/v1/users
func (a *User) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")
	if len(ids) == 0 {
		ctx.ResError(errors.NewBadRequestError("无效的请求参数"))
		return
	}

	for _, id := range ids {
		err := a.UserBll.Delete(ctx.CContext(), id)
		if err != nil {
			ctx.ResError(err)
			return
		}
	}

	ctx.ResOK()
}

// Enable 启用数据
// @Summary 启用数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Success 200 option.Interface "{status:OK}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router PATCH /api/v1/users/{id}/enable
func (a *User) Enable(ctx *context.Context) {
	err := a.UserBll.UpdateStatus(ctx.CContext(), ctx.Param("id"), 1)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// Disable 禁用数据
// @Summary 禁用数据
// @Param Access-Token header string false "访问令牌"
// @Param id path string true "记录ID"
// @Success 200 option.Interface "{status:OK}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router PATCH /api/v1/users/{id}/disable
func (a *User) Disable(ctx *context.Context) {
	err := a.UserBll.UpdateStatus(ctx.CContext(), ctx.Param("id"), 2)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

package system

import (
	"github.com/flipped-aurora/gf-vue-admin/app/api/system/internal"
	"github.com/flipped-aurora/gf-vue-admin/app/model/system/request"
	"github.com/flipped-aurora/gf-vue-admin/app/service/system"
	"github.com/flipped-aurora/gf-vue-admin/library/auth"
	"github.com/flipped-aurora/gf-vue-admin/library/common"
	"github.com/flipped-aurora/gf-vue-admin/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"strconv"
)

var User = new(user)

type user struct{}

// Register
// @Tags SysUser
// @Summary 用户注册账号
// @Produce  application/json
// @Param data body systemReq.Register true "用户名, 昵称, 密码, 角色ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"注册成功"}"
// @Router /user/register [post]
func (b *user) Register(r *ghttp.Request) *response.Response {
	var info request.UserRegister
	if err := r.Parse(&info); err != nil {
		return &response.Response{Error: err, MessageCode: response.ErrorCreated}
	}
	data, err := system.User.Register(&info)
	if err != nil {
		return &response.Response{Data: g.Map{"user": data}, Error: err, Message: "注册失败!"}
	}
	return &response.Response{Data: g.Map{"user": data}, Message: "注册成功!"}
}

// GetUserInfo
// @Tags SysUser
// @Summary 获取用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUserInfo [get]
func (b *user) GetUserInfo(r *ghttp.Request) *response.Response {
	var info request.UserFind
	if err := r.Parse(&info); err != nil {
		return &response.Response{Error: err, MessageCode: response.ErrorFirst}
	}
	claims := internal.NewClaims(r)
	if info.Uuid = claims.GetUserUuid(); info.Uuid == "" || claims.Error() != nil {
		return &response.Response{Error: claims.Error(), Message: "获取用户信息失败!"}
	}
	data, err := system.User.Find(&info)
	if err != nil {
		return &response.Response{Error: err, Message: "获取用户信息失败!"}
	}
	return &response.Response{Data: g.Map{"userInfo": data}, Message: "获取用户信息成功!"}
}

// SetUserInfo
// @Tags SysUser
// @Summary 设置用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysUser true "ID, 用户名, 昵称, 头像链接"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /user/setUserInfo [put]
func (b *user) SetUserInfo(r *ghttp.Request) *response.Response {
	var info request.UserUpdate
	if err := r.Parse(&info); err != nil {
		return &response.Response{Error: err, MessageCode: response.ErrorUpdated}
	}
	data, err := system.User.Update(&info)
	if err != nil {
		return &response.Response{Error: err, Message: "设置失败!"}
	}
	return &response.Response{Data: g.Map{"userInfo": data}, Message: "设置成功!"}
}

// ChangePassword
// @Tags SysUser
// @Summary 用户修改密码
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.UserChangePassword true "用户名, 原密码, 新密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/changePassword [post]
func (b *user) ChangePassword(r *ghttp.Request) *response.Response {
	var info request.UserChangePassword
	if err := r.Parse(&info); err != nil {
		return &response.Response{Error: err, MessageCode: response.ErrorUpdated}
	}
	if err := system.User.ChangePassword(&info); err != nil {
		return &response.Response{Error: err, Message: "修改失败!"}
	}
	return &response.Response{Message: "修改成功!"}
}

// Delete
// @Tags SystemUser
// @Summary 删除用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "用户ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /user/deleteUser [delete]
func (b *user) Delete(r *ghttp.Request) *response.Response {
	var info common.GetByID
	if err := r.Parse(&info); err != nil {
		return &response.Response{Error: err, MessageCode: response.ErrorDeleted}
	}
	claims := internal.NewClaims(r)
	if id := claims.GetUserID(); id == 0 || claims.Error() != nil {
		if id == info.ToUint() {
			return &response.Response{Error: claims.Error(), Message: "自我删除失败!"}
		}
	}
	if err := system.User.Delete(&info); err != nil {
		return &response.Response{Error: err, MessageCode: response.ErrorDeleted}
	}
	return &response.Response{MessageCode: response.SuccessDeleted}
}

// GetList
// @Tags SysUser
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUserList [post]
func (b *user) GetList(r *ghttp.Request) *response.Response {
	var info common.PageInfo
	if err := r.Parse(&info); err != nil {
		return &response.Response{Error: err, MessageCode: response.ErrorGetList}
	}
	list, total, err := system.User.GetList(&info)
	if err != nil {
		return &response.Response{Error: err, MessageCode: response.ErrorGetList}
	}
	return &response.Response{Data: common.NewPageResult(list, total, info), MessageCode: response.SuccessGetList}
}

// @Tags SysUser
// @Summary 更改用户权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.SetUserAuth true "用户UUID, 角色ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/setUserAuthority [post]
func (b *user) SetUserAuthority(r *ghttp.Request) *response.Response {
	var info request.UserSetAuthority
	if err := r.Parse(&info); err != nil {
		return &response.Response{Error: err, MessageCode: response.ErrorUpdated}
	}
	claims := internal.NewClaims(r)
	if info.ID = claims.GetUserID(); info.ID == 0 || claims.Error() != nil {
		err := claims.Error()
		return &response.Response{Error: err, Message: err.Error()}
	}
	if info.Uuid = claims.GetUserUuid(); info.Uuid == "" || claims.Error() != nil {
		err := claims.Error()
		return &response.Response{Error: err, Message: err.Error()}
	}
	if err := system.User.SetAuthority(&info); err != nil {
		return &response.Response{Error: err, Message: "修改失败!"}
	}
	_claims := claims.GetUserClaims()
	_claims.AuthorityId = info.AuthorityId
	if token, err := auth.NewJWT().CreateToken(_claims); err != nil {
		return &response.Response{Error: err, Message: "修改失败!"}
	} else {
		r.Response.Header().Set("new-token", token)
		r.Response.Header().Set("new-expires-at", strconv.FormatInt(_claims.ExpiresAt, 10))
		return &response.Response{Message: "修改成功!"}
	}
}

// @Tags SysUser
// @Summary 设置用户权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.SetUserAuthorities true "用户UUID, 角色ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/setUserAuthorities [post]
func (b *user) SetUserAuthorities(r *ghttp.Request) *response.Response {
	var info request.UserSetAuthorities
	if err := r.Parse(&info); err != nil {
		return &response.Response{Error: err, MessageCode: response.ErrorUpdated}
	}
	if err := system.User.SetUserAuthorities(&info); err != nil {
		return &response.Response{Error: err, Message: "修改失败!"}
	}
	return &response.Response{Message: "修改成功!"}
}
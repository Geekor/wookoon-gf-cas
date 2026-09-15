package wookooncas

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

/**
 * 从 context 中读取用户 UID
 * ---------------------------------------------
 * 前提： 使用了 CasAuthed 认证中间件的接口中方可提取
 */
func GetUserUid(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	return r.GetCtxVar(CtxUserId).String()
}

/**
 * 从 context 中读取用户 账号
 * ---------------------------------------------
 * 前提： 使用了 CasAuthed 认证中间件的接口中方可提取
 */
func GetUserPassport(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	return r.GetCtxVar(CtxUserPassport).String()
}

/**
 * 从 context 中读取用户昵称
 * ---------------------------------------------
 * 前提： 使用了 CasAuthed 认证中间件的接口中方可提取
 */
func GetUserName(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	return r.GetCtxVar(CtxUserName).String()
}

/**
 * 从 context 中读取用户邮箱
 * ---------------------------------------------
 * 前提： 使用了 CasAuthed 认证中间件的接口中方可提取
 */
func GetUserEmail(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	return r.GetCtxVar(CtxUserEmail).String()
}

/**
 * 从 context 中读取用户角色
 * ---------------------------------------------
 * 前提： 使用了 CasAuthed 认证中间件的接口中方可提取
 */
func GetUserRoles(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	return r.GetCtxVar(CtxUserRoles).String()
}

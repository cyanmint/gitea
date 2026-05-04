// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"errors"
	"net/http"

	"code.gitea.io/gitea/modules/auth/password"
	"code.gitea.io/gitea/modules/optional"
	"code.gitea.io/gitea/modules/setting"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
	user_service "code.gitea.io/gitea/services/user"
)

// ChangePassword changes the currently signed-in user's password.
func ChangePassword(ctx *context.APIContext) {
	// swagger:operation POST /user/change_password user userChangePassword
	// ---
	// summary: Change the authenticated user's password
	// consumes:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/ChangePasswordOption"
	// responses:
	//   "204":
	//     description: "No Content – password changed successfully"
	//   "400":
	//     "$ref": "#/responses/error"
	//   "403":
	//     "$ref": "#/responses/forbidden"

	form := web.GetForm(ctx).(*api.ChangePasswordOption)

	if ctx.Doer.IsPasswordSet() && !ctx.Doer.ValidatePassword(form.OldPassword) {
		ctx.APIError(http.StatusBadRequest, errors.New("current password is incorrect"))
		return
	}

	opts := &user_service.UpdateAuthOptions{
		Password:           optional.Some(form.NewPassword),
		MustChangePassword: optional.Some(false),
	}
	if err := user_service.UpdateAuth(ctx, ctx.Doer, opts); err != nil {
		switch {
		case errors.Is(err, password.ErrMinLength):
			ctx.APIError(http.StatusBadRequest, errors.New("password too short"))
		case errors.Is(err, password.ErrComplexity):
			ctx.APIError(http.StatusBadRequest, errors.New("password not complex enough"))
		default:
			ctx.APIErrorInternal(err)
		}
		return
	}
	_ = setting.MinPasswordLength // used for error message context only
	ctx.Status(http.StatusNoContent)
}

// DeleteSelf deletes the currently authenticated user's account.
func DeleteSelf(ctx *context.APIContext) {
	// swagger:operation DELETE /user user userDeleteSelf
	// ---
	// summary: Delete the authenticated user
	// responses:
	//   "204":
	//     description: "No Content – account deleted"
	//   "403":
	//     "$ref": "#/responses/forbidden"

	if ctx.Doer.IsAdmin {
		ctx.APIError(http.StatusForbidden, errors.New("admins cannot delete their own account via the API"))
		return
	}
	if err := user_service.DeleteUser(ctx, ctx.Doer, false); err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
"errors"
"net/http"

"code.gitea.io/gitea/models/db"
user_model "code.gitea.io/gitea/models/user"
"code.gitea.io/gitea/modules/auth/password"
"code.gitea.io/gitea/modules/log"
"code.gitea.io/gitea/modules/setting"
api "code.gitea.io/gitea/modules/structs"
"code.gitea.io/gitea/modules/web"
"code.gitea.io/gitea/services/context"
"code.gitea.io/gitea/services/convert"
)

// Register creates a new user account via open self-registration.
func Register(ctx *context.APIContext) {
// swagger:operation POST /user/register user userRegister
// ---
// summary: Register a new user (open registration)
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - name: body
//   in: body
//   required: true
//   schema:
//     "$ref": "#/definitions/RegisterUserOption"
// responses:
//   "201":
//     "$ref": "#/responses/User"
//   "400":
//     "$ref": "#/responses/error"
//   "403":
//     "$ref": "#/responses/forbidden"
//   "422":
//     "$ref": "#/responses/validationError"

if setting.Service.DisableRegistration || setting.Service.AllowOnlyExternalRegistration {
ctx.APIError(http.StatusForbidden, "registration is disabled on this instance")
return
}

form := web.GetForm(ctx).(*api.RegisterUserOption)

if len(form.Password) < setting.MinPasswordLength {
ctx.APIError(http.StatusBadRequest, errors.New("password too short"))
return
}
if !password.IsComplexEnough(form.Password) {
ctx.APIError(http.StatusBadRequest, errors.New("password does not meet complexity requirements"))
return
}
if err := password.IsPwned(ctx, form.Password); err != nil {
if password.IsErrIsPwnedRequest(err) {
log.Error("IsPwned request error: %v", err)
}
ctx.APIError(http.StatusBadRequest, errors.New("password has been found in a data breach; choose a different password"))
return
}
if !user_model.IsEmailDomainAllowed(form.Email) {
ctx.APIError(http.StatusUnprocessableEntity, errors.New("email domain is not allowed"))
return
}

u := &user_model.User{
Name:   form.Username,
Email:  form.Email,
Passwd: form.Password,
}

meta := &user_model.Meta{
InitialIP:        ctx.RemoteAddr(),
InitialUserAgent: ctx.Req.UserAgent(),
}
if err := user_model.CreateUser(ctx, u, meta, nil); err != nil {
switch {
case user_model.IsErrUserAlreadyExist(err),
db.IsErrNameReserved(err),
db.IsErrNamePatternNotAllowed(err),
db.IsErrNameCharsNotAllowed(err),
user_model.IsErrEmailAlreadyUsed(err),
user_model.IsErrEmailCharIsNotSupported(err),
user_model.IsErrEmailInvalid(err):
ctx.APIError(http.StatusUnprocessableEntity, err)
default:
ctx.APIErrorInternal(err)
}
return
}

ctx.JSON(http.StatusCreated, convert.ToUser(ctx, u, ctx.Doer))
}

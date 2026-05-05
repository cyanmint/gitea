// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"net/http"

	actions_model "code.gitea.io/gitea/models/actions"
	repo_model "code.gitea.io/gitea/models/repo"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
)

// GetActionsPermissions returns the owner-level Actions permissions for the authenticated user.
func GetActionsPermissions(ctx *context.APIContext) {
	// swagger:operation GET /user/actions/permissions user userGetActionsPermissions
	// ---
	// summary: Get the Actions permissions for the authenticated user
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/UserActionsPermissions"

	cfg, err := actions_model.GetOwnerActionsConfig(ctx, ctx.Doer.ID)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	mode := string(cfg.TokenPermissionMode)
	if mode == "" {
		mode = string(repo_model.ActionsTokenPermissionModePermissive)
	}
	ids := cfg.AllowedCrossRepoIDs
	if ids == nil {
		ids = []int64{}
	}
	ctx.JSON(http.StatusOK, api.UserActionsPermissions{
		TokenPermissionMode: mode,
		AllowedCrossRepoIDs: ids,
	})
}

// UpdateActionsPermissions updates the owner-level Actions permissions for the authenticated user.
func UpdateActionsPermissions(ctx *context.APIContext) {
	// swagger:operation PUT /user/actions/permissions user userUpdateActionsPermissions
	// ---
	// summary: Update the Actions permissions for the authenticated user
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UserActionsPermissions"
	// responses:
	//   "200":
	//     "$ref": "#/responses/UserActionsPermissions"
	//   "400":
	//     "$ref": "#/responses/error"

	form := web.GetForm(ctx).(*api.UserActionsPermissions)

	// Validate permission mode
	mode := repo_model.ActionsTokenPermissionMode(form.TokenPermissionMode)
	validModes := repo_model.ActionsTokenPermissionMode("").EnumValues()
	isValid := false
	for _, v := range validModes {
		if mode == v {
			isValid = true
			break
		}
	}
	if form.TokenPermissionMode != "" && !isValid {
		ctx.APIError(http.StatusBadRequest, "invalid token_permission_mode")
		return
	}

	cfg, err := actions_model.GetOwnerActionsConfig(ctx, ctx.Doer.ID)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	if form.TokenPermissionMode != "" {
		cfg.TokenPermissionMode = mode
	}
	if form.AllowedCrossRepoIDs != nil {
		cfg.AllowedCrossRepoIDs = form.AllowedCrossRepoIDs
	}
	if err := actions_model.SetOwnerActionsConfig(ctx, ctx.Doer.ID, cfg); err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	modeStr := string(cfg.TokenPermissionMode)
	if modeStr == "" {
		modeStr = string(repo_model.ActionsTokenPermissionModePermissive)
	}
	ids := cfg.AllowedCrossRepoIDs
	if ids == nil {
		ids = []int64{}
	}
	ctx.JSON(http.StatusOK, api.UserActionsPermissions{
		TokenPermissionMode: modeStr,
		AllowedCrossRepoIDs: ids,
	})
}

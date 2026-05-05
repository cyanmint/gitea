// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"net/http"

	auth_model "code.gitea.io/gitea/models/auth"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
)

// ListOAuth2Grants lists all OAuth2 grants for the authenticated user.
func ListOAuth2Grants(ctx *context.APIContext) {
	// swagger:operation GET /user/applications/grants user userListOAuth2Grants
	// ---
	// summary: List all OAuth2 grants for the authenticated user
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/OAuth2GrantList"
	//   "401":
	//     "$ref": "#/responses/unauthorized"

	grants, err := auth_model.GetOAuth2GrantsByUserID(ctx, ctx.Doer.ID)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	result := make([]*api.OAuth2Grant, 0, len(grants))
	for _, g := range grants {
		appName := ""
		if g.Application != nil {
			appName = g.Application.Name
		}
		result = append(result, &api.OAuth2Grant{
			ID:            g.ID,
			UserID:        g.UserID,
			ApplicationID: g.ApplicationID,
			ApplicationName: appName,
			Scope:         g.Scope,
			Created:       g.CreatedUnix.AsTime(),
		})
	}
	ctx.JSON(http.StatusOK, result)
}

// RevokeOAuth2Grant revokes an OAuth2 grant for the authenticated user.
func RevokeOAuth2Grant(ctx *context.APIContext) {
	// swagger:operation DELETE /user/applications/grants/{id} user userRevokeOAuth2Grant
	// ---
	// summary: Revoke an OAuth2 grant for the authenticated user
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: Grant ID
	//   required: true
	//   type: integer
	//   format: int64
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "401":
	//     "$ref": "#/responses/unauthorized"
	//   "404":
	//     "$ref": "#/responses/notFound"

	grantID := ctx.PathParamInt64("id")
	if err := auth_model.RevokeOAuth2Grant(ctx, grantID, ctx.Doer.ID); err != nil {
		ctx.APIError(http.StatusNotFound, "grant not found or not owned by user")
		return
	}
	ctx.Status(http.StatusNoContent)
}

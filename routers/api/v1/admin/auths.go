// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package admin

import (
	"net/http"

	"code.gitea.io/gitea/models/auth"
	"code.gitea.io/gitea/models/db"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
)

// ListAuthSources returns all configured authentication sources
func ListAuthSources(ctx *context.APIContext) {
	// swagger:operation GET /admin/auths admin adminListAuthSources
	// ---
	// summary: List all authentication sources
	// produces:
	// - application/json
	// parameters:
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/AuthSourceList"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	sources, total, err := db.FindAndCount[auth.Source](ctx, auth.FindSourcesOptions{})
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	res := make([]*api.AuthSource, 0, len(sources))
	for _, s := range sources {
		res = append(res, &api.AuthSource{
			ID:            s.ID,
			Name:          s.Name,
			Type:          int(s.Type),
			TypeName:      s.TypeName(),
			IsActive:      s.IsActive,
			IsSyncEnabled: s.IsSyncEnabled,
			Created:       s.CreatedUnix.AsTime(),
			Updated:       s.UpdatedUnix.AsTime(),
		})
	}

	ctx.SetTotalCountHeader(total)
	ctx.JSON(http.StatusOK, res)
}

// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package web

import (
	"net/http"

	"code.gitea.io/gitea/modules/templates"
	"code.gitea.io/gitea/services/context"
)

const tplSPA templates.TplName = "spa"

// SPA renders the SPA shell page. This catch-all handler serves index.html for
// any route not matched by more specific handlers, allowing the Vue Router
// running in the browser to take over client-side navigation.
func SPA(ctx *context.Context) {
	ctx.HTML(http.StatusOK, tplSPA)
}

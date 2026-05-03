// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package web

import (
	"net/http"
	"strings"

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

// SPAOrFeed serves the SPA shell for normal browser requests but delegates to
// feedHandler for RSS/Atom feed requests (URL path ending in ".rss" or ".atom").
// This is used for the "/{username}" route which serves both the user/org profile
// page (SPA) and user activity feeds (Go handler).
func SPAOrFeed(feedHandler func(*context.Context)) func(*context.Context) {
	return func(ctx *context.Context) {
		username := ctx.PathParam("username")
		if strings.HasSuffix(username, ".rss") || strings.HasSuffix(username, ".atom") {
			feedHandler(ctx)
			return
		}
		SPA(ctx)
	}
}

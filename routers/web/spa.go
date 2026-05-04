// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package web

import (
"net/http"
"strings"

"code.gitea.io/gitea/services/context"
)

// SPA returns HTTP 404 for requests that previously served the SPA shell.
// The frontend is now deployed independently (e.g. GitHub Pages) and is no
// longer served by the Go binary.
func SPA(ctx *context.Context) {
ctx.Status(http.StatusNotFound)
}

// SPAOrFeed serves RSS/Atom feeds for requests whose username path-param ends
// in ".rss" or ".atom", and returns 404 for all other requests.
func SPAOrFeed(feedHandler func(*context.Context)) func(*context.Context) {
return func(ctx *context.Context) {
username := ctx.PathParam("username")
if strings.HasSuffix(username, ".rss") || strings.HasSuffix(username, ".atom") {
feedHandler(ctx)
return
}
ctx.Status(http.StatusNotFound)
}
}

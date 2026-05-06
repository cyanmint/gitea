// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package admin

import (
	"net/http"
	"sort"

	"code.gitea.io/gitea/models/db"
	system_model "code.gitea.io/gitea/models/system"
	"code.gitea.io/gitea/modules/cache"
	"code.gitea.io/gitea/modules/json"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/setting/config"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
)

// SelfCheck runs admin self-check diagnostics and returns the results
func SelfCheck(ctx *context.APIContext) {
	// swagger:operation GET /admin/self_check admin adminSelfCheck
	// ---
	// summary: Run admin self-check diagnostics
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/SelfCheckResult"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	result := &api.SelfCheckResult{
		StartupProblems:              setting.StartupProblems,
		InconsistentCollationColumns: []string{},
	}
	if result.StartupProblems == nil {
		result.StartupProblems = []string{}
	}

	r, err := db.CheckCollationsDefaultEngine()
	if err == nil && r != nil {
		result.DatabaseCollationMismatch = !r.CollationEquals(r.DatabaseCollation, r.ExpectedCollation)
		result.DatabaseCollationCaseInsensitive = !r.IsCollationCaseSensitive(r.DatabaseCollation)
		if r.InconsistentCollationColumns != nil {
			result.InconsistentCollationColumns = r.InconsistentCollationColumns
		}
	}

	elapsed, cacheErr := cache.Test()
	if cacheErr != nil {
		result.CacheError = cacheErr.Error()
	} else {
		result.CacheElapsedMs = elapsed.Milliseconds()
		result.CacheSlow = elapsed > cache.SlowCacheThreshold
	}

	ctx.JSON(http.StatusOK, result)
}

// GetConfigSettings returns all dynamic configuration settings
func GetConfigSettings(ctx *context.APIContext) {
	// swagger:operation GET /admin/config/settings admin adminGetConfigSettings
	// ---
	// summary: List all dynamic configuration settings
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/ConfigSettingList"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	_, dbSettings, err := system_model.GetAllSettings(ctx)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	keys := config.GetAllConfigOptionKeys()
	sort.Strings(keys)

	res := make([]*api.ConfigSetting, 0, len(keys))
	for _, key := range keys {
		val := dbSettings[key]
		res = append(res, &api.ConfigSetting{Key: key, Value: val})
	}

	ctx.JSON(http.StatusOK, res)
}

// SetConfigSettings updates one or more dynamic configuration settings
func SetConfigSettings(ctx *context.APIContext) {
	// swagger:operation PATCH /admin/config/settings admin adminSetConfigSettings
	// ---
	// summary: Update dynamic configuration settings
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   required: true
	//   schema:
	//     type: array
	//     items:
	//       "$ref": "#/definitions/ConfigSetting"
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	//   "400":
	//     "$ref": "#/responses/error"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	var settings []*api.ConfigSetting
	if err := json.NewDecoder(ctx.Req.Body).Decode(&settings); err != nil {
		ctx.APIError(http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}

	toSave := make(map[string]string, len(settings))
	for _, s := range settings {
		if config.GetConfigOption(s.Key) == nil {
			ctx.APIError(http.StatusBadRequest, "unknown config key: "+s.Key)
			return
		}
		toSave[s.Key] = s.Value
	}

	if err := system_model.SetSettings(ctx, toSave); err != nil {
		ctx.APIErrorInternal(err)
		return
	}
	config.GetDynGetter().InvalidateCache()
	ctx.Status(http.StatusNoContent)
}

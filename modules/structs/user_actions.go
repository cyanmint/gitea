// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package structs

// UserActionsPermissions holds the Actions permissions for a user or organisation owner.
// swagger:model
type UserActionsPermissions struct {
	// TokenPermissionMode is the default Actions token permission mode.
	// Possible values: "permissive", "restricted"
	// example: permissive
	TokenPermissionMode string `json:"token_permission_mode"`
	// AllowedCrossRepoIDs is the list of repository IDs allowed for cross-repo access.
	AllowedCrossRepoIDs []int64 `json:"allowed_cross_repo_ids"`
}

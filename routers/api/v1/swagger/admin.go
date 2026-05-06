// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package swagger

import (
	api "code.gitea.io/gitea/modules/structs"
)

// AuthSourceList
// swagger:response AuthSourceList
type swaggerResponseAuthSourceList struct {
	// in:body
	Body []api.AuthSource `json:"body"`
}

// QueueStatList
// swagger:response QueueStatList
type swaggerResponseQueueStatList struct {
	// in:body
	Body []api.QueueStat `json:"body"`
}

// StacktraceResponse
// swagger:response StacktraceResponse
type swaggerResponseStacktrace struct {
	// in:body
	Body api.StacktraceResult `json:"body"`
}

// SelfCheckResult
// swagger:response SelfCheckResult
type swaggerResponseSelfCheckResult struct {
	// in:body
	Body api.SelfCheckResult `json:"body"`
}

// ConfigSettingList
// swagger:response ConfigSettingList
type swaggerResponseConfigSettingList struct {
	// in:body
	Body []api.ConfigSetting `json:"body"`
}

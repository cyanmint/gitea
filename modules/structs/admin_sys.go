// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package structs

import "time"

// AuthSource represents an authentication source (e.g. LDAP, OAuth2, SMTP)
// swagger:model
type AuthSource struct {
	// The unique identifier of the authentication source
	ID int64 `json:"id"`
	// The display name of the authentication source
	Name string `json:"name"`
	// The numeric type code (LDAP=2, DLDAP=5, SMTP=3, PAM=4, OAuth2=6, SSPI=7)
	Type int `json:"type"`
	// The human-readable name of the type (e.g. "LDAP (via BindDN)")
	TypeName string `json:"type_name"`
	// Whether this source is enabled
	IsActive bool `json:"is_active"`
	// Whether periodic sync is enabled for this source
	IsSyncEnabled bool `json:"is_sync_enabled"`
	// When the source was created
	Created time.Time `json:"created"`
	// When the source was last updated
	Updated time.Time `json:"updated"`
}

// QueueStat represents runtime statistics for a single managed queue
// swagger:model
type QueueStat struct {
	// Internal queue identifier
	ID int64 `json:"id"`
	// Queue name
	Name string `json:"name"`
	// Queue backend type (e.g. "level", "channel")
	Type string `json:"type"`
	// Item type processed by this queue
	ItemTypeName string `json:"item_type_name"`
	// Number of active worker goroutines
	WorkerNumber int `json:"worker_number"`
	// Number of workers currently processing items
	WorkerActiveNumber int `json:"worker_active_number"`
	// Maximum allowed workers (-1 = unlimited)
	WorkerMaxNumber int `json:"worker_max_number"`
	// Number of items waiting in the queue
	QueueItemNumber int `json:"queue_item_number"`
}

// StackEntry is a single frame in a goroutine stack trace
// swagger:model
type StackEntry struct {
	// Function name
	Function string `json:"function"`
	// Source file path
	File string `json:"file"`
	// Line number in the source file
	Line int `json:"line"`
}

// StackLabel is a pprof label attached to a goroutine stack
// swagger:model
type StackLabel struct {
	// Label name
	Name string `json:"name"`
	// Label value
	Value string `json:"value"`
}

// GoroutineStack describes a set of goroutines sharing the same stack trace
// swagger:model
type GoroutineStack struct {
	// Number of goroutines with this stack trace
	Count int64 `json:"count"`
	// Human-readable description
	Description string `json:"description"`
	// pprof labels
	Labels []*StackLabel `json:"labels,omitempty"`
	// Stack frames (bottom of call stack first)
	Entry []*StackEntry `json:"entry,omitempty"`
}

// ProcessInfo represents a Gitea-managed process and its goroutine stacks
// swagger:model
type ProcessInfo struct {
	// Process ID
	PID string `json:"pid"`
	// Parent process ID
	ParentPID string `json:"parent_pid"`
	// Human-readable description
	Description string `json:"description"`
	// When the process was started
	Start time.Time `json:"start"`
	// Process type ("system" or "request" etc.)
	Type string `json:"type"`
	// Child processes
	Children []*ProcessInfo `json:"children,omitempty"`
	// Goroutine stacks associated with this process
	Stacks []*GoroutineStack `json:"stacks,omitempty"`
}

// SelfCheckResult contains the results of the admin self-check diagnostics
// swagger:model
type SelfCheckResult struct {
	// List of startup configuration problems detected
	StartupProblems []string `json:"startup_problems"`
	// Whether the database collation check found a mismatch
	DatabaseCollationMismatch bool `json:"database_collation_mismatch"`
	// Whether the database collation is case-insensitive (may cause issues)
	DatabaseCollationCaseInsensitive bool `json:"database_collation_case_insensitive"`
	// List of columns with inconsistent collations
	InconsistentCollationColumns []string `json:"inconsistent_collation_columns"`
	// Error message from the cache test, empty if cache is healthy
	CacheError string `json:"cache_error"`
	// Whether the cache responded slower than the threshold
	CacheSlow bool `json:"cache_slow"`
	// Cache response duration in milliseconds
	CacheElapsedMs int64 `json:"cache_elapsed_ms"`
}

// ConfigSetting represents one editable dynamic configuration key/value pair
// swagger:model
type ConfigSetting struct {
	// Configuration key (e.g. "repository.default_branch")
	Key string `json:"key"`
	// Current value as JSON (number, bool, or quoted string)
	Value string `json:"value"`
}

// StacktraceResult is the response body for the stacktrace endpoint
// swagger:model
type StacktraceResult struct {
	// Total number of goroutines reported by the Go runtime
	NumGoroutine int `json:"num_goroutine"`
	// Number of Gitea-managed processes
	ProcessCount int `json:"process_count"`
	// Number of goroutines in the pprof profile
	GoroutineCount int64 `json:"goroutine_count"`
	// Process tree with goroutine stacks
	Processes []*ProcessInfo `json:"processes"`
}

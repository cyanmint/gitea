// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package admin

import (
	"net/http"
	"runtime"
	"sort"

	"code.gitea.io/gitea/modules/process"
	"code.gitea.io/gitea/modules/queue"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
)

// ListQueues returns runtime statistics for all managed queues
func ListQueues(ctx *context.APIContext) {
	// swagger:operation GET /admin/monitor/queues admin adminListQueues
	// ---
	// summary: List all managed queues and their runtime statistics
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/QueueStatList"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	managed := queue.GetManager().ManagedQueues()

	ids := make([]int64, 0, len(managed))
	for id := range managed {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	res := make([]*api.QueueStat, 0, len(managed))
	for _, id := range ids {
		mq := managed[id]
		res = append(res, &api.QueueStat{
			ID:                 id,
			Name:               mq.GetName(),
			Type:               mq.GetType(),
			ItemTypeName:       mq.GetItemTypeName(),
			WorkerNumber:       mq.GetWorkerNumber(),
			WorkerActiveNumber: mq.GetWorkerActiveNumber(),
			WorkerMaxNumber:    mq.GetWorkerMaxNumber(),
			QueueItemNumber:    mq.GetQueueItemNumber(),
		})
	}

	ctx.JSON(http.StatusOK, res)
}

// GetStacktrace returns goroutine stacktraces for all Gitea-managed processes
func GetStacktrace(ctx *context.APIContext) {
	// swagger:operation GET /admin/monitor/stacktrace admin adminGetStacktrace
	// ---
	// summary: Get goroutine stacktrace for all managed processes
	// produces:
	// - application/json
	// parameters:
	// - name: show
	//   in: query
	//   description: which goroutines to include ("all" or "process")
	//   type: string
	// responses:
	//   "200":
	//     "$ref": "#/responses/StacktraceResponse"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	show := ctx.FormString("show")
	if show == "" {
		show = "all"
	}
	noSystem := show == "process"

	processes, processCount, goroutineCount, err := process.GetManager().ProcessStacktraces(false, noSystem)
	if err != nil {
		ctx.APIErrorInternal(err)
		return
	}

	var convertProcesses func([]*process.Process) []*api.ProcessInfo
	convertProcesses = func(procs []*process.Process) []*api.ProcessInfo {
		out := make([]*api.ProcessInfo, 0, len(procs))
		for _, p := range procs {
			info := &api.ProcessInfo{
				PID:         string(p.PID),
				ParentPID:   string(p.ParentPID),
				Description: p.Description,
				Start:       p.Start,
				Type:        p.Type,
			}
			for _, s := range p.Stacks {
				gs := &api.GoroutineStack{
					Count:       s.Count,
					Description: s.Description,
				}
				for _, lbl := range s.Labels {
					gs.Labels = append(gs.Labels, &api.StackLabel{Name: lbl.Name, Value: lbl.Value})
				}
				for _, e := range s.Entry {
					gs.Entry = append(gs.Entry, &api.StackEntry{Function: e.Function, File: e.File, Line: e.Line})
				}
				info.Stacks = append(info.Stacks, gs)
			}
			info.Children = convertProcesses(p.Children)
			out = append(out, info)
		}
		return out
	}

	ctx.JSON(http.StatusOK, &api.StacktraceResult{
		NumGoroutine:   runtime.NumGoroutine(),
		ProcessCount:   processCount,
		GoroutineCount: goroutineCount,
		Processes:      convertProcesses(processes),
	})
}

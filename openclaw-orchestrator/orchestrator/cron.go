package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// CronSchedule mirrors the stable schedule fields returned by OpenClaw cron.
// The fixed OpenClaw version stores one of atMs, everyMs or expr depending on kind.
type CronSchedule struct {
	Kind    string `json:"kind"`
	AtMs    int64  `json:"atMs,omitempty"`
	EveryMs int64  `json:"everyMs,omitempty"`
	Expr    string `json:"expr,omitempty"`
	TZ      string `json:"tz,omitempty"`
}

type CronPayload struct {
	Kind       string   `json:"kind"`
	Message    string   `json:"message,omitempty"`
	Text       string   `json:"text,omitempty"`
	Model      string   `json:"model,omitempty"`
	ToolsAllow []string `json:"toolsAllow,omitempty"`
}

type CronState struct {
	NextRunAtMs       int64  `json:"nextRunAtMs,omitempty"`
	LastRunAtMs       int64  `json:"lastRunAtMs,omitempty"`
	LastRunStatus     string `json:"lastRunStatus,omitempty"`
	LastStatus        string `json:"lastStatus,omitempty"`
	LastError         string `json:"lastError,omitempty"`
	LastRunError      string `json:"lastRunError,omitempty"`
	LastErrorReason   string `json:"lastErrorReason,omitempty"`
	LastRunSummary    string `json:"lastRunSummary,omitempty"`
	RunningAtMs       int64  `json:"runningAtMs,omitempty"`
	ConsecutiveErrors int    `json:"consecutiveErrors,omitempty"`
}

type CronJob struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	DisplayName    string       `json:"displayName,omitempty"`
	DeclarationKey string       `json:"declarationKey,omitempty"`
	Description    string       `json:"description,omitempty"`
	Enabled        bool         `json:"enabled"`
	AgentID        string       `json:"agentId,omitempty"`
	SessionTarget  string       `json:"sessionTarget,omitempty"`
	SessionKey     string       `json:"sessionKey,omitempty"`
	WakeMode       string       `json:"wakeMode,omitempty"`
	Schedule       CronSchedule `json:"schedule"`
	Payload        CronPayload  `json:"payload"`
	State          CronState    `json:"state"`
	Status         string       `json:"status,omitempty"`
}

type CronListResult struct {
	Jobs       []CronJob `json:"jobs"`
	Total      int       `json:"total,omitempty"`
	HasMore    bool      `json:"hasMore,omitempty"`
	NextOffset int       `json:"nextOffset,omitempty"`
}

type CronRun struct {
	RunID          string `json:"runId,omitempty"`
	JobID          string `json:"jobId,omitempty"`
	Status         string `json:"status,omitempty"`
	DeliveryStatus string `json:"deliveryStatus,omitempty"`
	Summary        string `json:"summary,omitempty"`
	Error          string `json:"error,omitempty"`
	ErrorReason    string `json:"errorReason,omitempty"`
	Cause          string `json:"cause,omitempty"`
	Provider       string `json:"provider,omitempty"`
	Model          string `json:"model,omitempty"`
	StartedAtMs    int64  `json:"startedAtMs,omitempty"`
	EndedAtMs      int64  `json:"endedAtMs,omitempty"`
	RunAtMs        int64  `json:"runAtMs,omitempty"`
}

type CronRunsResult struct {
	Entries    []CronRun `json:"entries"`
	Total      int       `json:"total,omitempty"`
	HasMore    bool      `json:"hasMore,omitempty"`
	NextOffset int       `json:"nextOffset,omitempty"`
}

type CronJobInput struct {
	ID             string
	Name           string
	Description    string
	DeclarationKey string
	Enabled        bool
	AgentID        string
	SessionTarget  string
	SessionKey     string
	ScheduleKind   string
	ScheduleValue  string
	Timezone       string
	Message        string
	Model          string
}

func (s *Service) CronList(ctx context.Context) (CronListResult, error) {
	out, err := s.run(ctx, nil, "cron", "list", "--all", "--json")
	if err != nil {
		return CronListResult{}, err
	}
	var result CronListResult
	if err := json.Unmarshal(out, &result); err != nil {
		return CronListResult{}, fmt.Errorf("%w: cron list: %v", ErrInvalidResponse, err)
	}
	for index := range result.Jobs {
		normalizeCronJob(&result.Jobs[index])
	}
	if result.Total == 0 {
		result.Total = len(result.Jobs)
	}
	return result, nil
}

func (s *Service) CronStatus(ctx context.Context) (map[string]any, error) {
	out, err := s.run(ctx, nil, "cron", "status", "--json")
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("%w: cron status: %v", ErrInvalidResponse, err)
	}
	return result, nil
}

func (s *Service) CronAdd(ctx context.Context, input CronJobInput) (CronJob, error) {
	args, err := cronMutationArgs("add", input)
	if err != nil {
		return CronJob{}, err
	}
	out, err := s.run(ctx, nil, args...)
	if err != nil {
		return CronJob{}, err
	}
	job, err := decodeCronJob(out)
	if err != nil {
		return CronJob{}, fmt.Errorf("cron add response: %w", err)
	}
	return job, nil
}

func (s *Service) CronEdit(ctx context.Context, input CronJobInput) (CronJob, error) {
	input.ID = strings.TrimSpace(input.ID)
	if input.ID == "" {
		return CronJob{}, fmt.Errorf("cron job id is required")
	}
	args, err := cronMutationArgs("edit", input)
	if err != nil {
		return CronJob{}, err
	}
	out, err := s.run(ctx, nil, args...)
	if err != nil {
		return CronJob{}, err
	}
	job, err := decodeCronJob(out)
	if err == nil {
		return job, nil
	}
	job, found, listErr := s.cronJobByID(ctx, input.ID)
	if listErr != nil {
		return CronJob{}, fmt.Errorf("cron edit completed but status refresh failed: %w", listErr)
	}
	if !found {
		return CronJob{}, fmt.Errorf("cron edit response: %w", err)
	}
	return job, nil
}

func (s *Service) CronRemove(ctx context.Context, id string) (map[string]any, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("cron job id is required")
	}
	out, err := s.run(ctx, nil, "cron", "rm", id, "--json")
	if err != nil {
		return nil, err
	}
	return decodeJSONObject(out)
}

func (s *Service) CronEnable(ctx context.Context, id string, enabled bool) (CronJob, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return CronJob{}, fmt.Errorf("cron job id is required")
	}
	command := "disable"
	if enabled {
		command = "enable"
	}
	out, err := s.run(ctx, nil, "cron", command, id)
	if err != nil {
		return CronJob{}, err
	}
	job, err := decodeCronJob(out)
	if err == nil {
		return job, nil
	}
	job, found, listErr := s.cronJobByID(ctx, id)
	if listErr != nil {
		return CronJob{}, fmt.Errorf("cron %s completed but status refresh failed: %w", command, listErr)
	}
	if !found {
		return CronJob{ID: id, Enabled: enabled, Status: map[bool]string{true: "idle", false: "disabled"}[enabled]}, nil
	}
	return job, nil
}

func (s *Service) CronRun(ctx context.Context, id string) (map[string]any, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("cron job id is required")
	}
	out, err := s.run(ctx, nil, "cron", "run", id)
	if err != nil {
		return nil, err
	}
	result, decodeErr := decodeJSONObject(out)
	if decodeErr == nil {
		return result, nil
	}
	message := strings.TrimSpace(string(out))
	if message == "" {
		message = "OpenClaw 已接收立即执行请求"
	}
	return map[string]any{"submitted": true, "message": message}, nil
}

func (s *Service) CronRuns(ctx context.Context, id string, limit int) (CronRunsResult, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return CronRunsResult{}, fmt.Errorf("cron job id is required")
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}
	out, err := s.run(ctx, nil, "cron", "runs", "--id", id, "--limit", strconv.Itoa(limit))
	if err != nil {
		return CronRunsResult{}, err
	}
	var result CronRunsResult
	if err := json.Unmarshal(out, &result); err != nil {
		return CronRunsResult{}, fmt.Errorf("%w: cron runs: %v", ErrInvalidResponse, err)
	}
	return result, nil
}

func (s *Service) cronJobByID(ctx context.Context, id string) (CronJob, bool, error) {
	list, err := s.CronList(ctx)
	if err != nil {
		return CronJob{}, false, err
	}
	for _, job := range list.Jobs {
		if job.ID == id {
			return job, true, nil
		}
	}
	return CronJob{}, false, nil
}

func cronMutationArgs(command string, input CronJobInput) ([]string, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.AgentID = strings.TrimSpace(input.AgentID)
	input.SessionTarget = strings.TrimSpace(input.SessionTarget)
	input.SessionKey = strings.TrimSpace(input.SessionKey)
	input.ScheduleKind = strings.ToLower(strings.TrimSpace(input.ScheduleKind))
	input.ScheduleValue = strings.TrimSpace(input.ScheduleValue)
	input.Timezone = strings.TrimSpace(input.Timezone)
	input.Message = strings.TrimSpace(input.Message)
	input.Model = strings.TrimSpace(input.Model)
	if input.Name == "" {
		return nil, fmt.Errorf("cron job name is required")
	}
	if input.Description == "" {
		input.Description = "STA-100 定时任务"
	}
	if input.Message == "" {
		return nil, fmt.Errorf("cron job message is required")
	}
	if input.ScheduleKind == "" {
		input.ScheduleKind = "cron"
	}
	if input.ScheduleValue == "" {
		return nil, fmt.Errorf("cron schedule is required")
	}
	if input.SessionTarget == "" {
		input.SessionTarget = "isolated"
	}

	args := []string{"cron", command}
	if command == "add" {
		args = append(args, "--json")
	}
	if command == "edit" {
		if input.ID == "" {
			return nil, fmt.Errorf("cron job id is required")
		}
		args = append(args, input.ID)
	}
	args = append(args, "--name", input.Name, "--description", input.Description)
	if input.DeclarationKey != "" && command == "add" {
		args = append(args, "--declaration-key", strings.TrimSpace(input.DeclarationKey))
	}
	if input.AgentID != "" {
		args = append(args, "--agent", input.AgentID)
	}
	if input.SessionKey != "" {
		args = append(args, "--session-key", input.SessionKey)
	}
	args = append(args, "--session", input.SessionTarget, "--no-deliver", "--message", input.Message)
	switch input.ScheduleKind {
	case "every":
		args = append(args, "--every", input.ScheduleValue)
	case "at":
		args = append(args, "--at", input.ScheduleValue)
	case "cron":
		args = append(args, "--cron", input.ScheduleValue)
		if input.Timezone != "" {
			args = append(args, "--tz", input.Timezone)
		}
	default:
		return nil, fmt.Errorf("unsupported cron schedule kind %q", input.ScheduleKind)
	}
	if input.Model != "" {
		args = append(args, "--model", input.Model)
	}
	if input.Enabled && command == "edit" {
		args = append(args, "--enable")
	} else if !input.Enabled && command == "add" {
		args = append(args, "--disabled")
	} else if !input.Enabled {
		args = append(args, "--disable")
	}
	return args, nil
}

func decodeCronJob(out []byte) (CronJob, error) {
	var envelope struct {
		Job *CronJob `json:"job"`
	}
	if err := json.Unmarshal(out, &envelope); err != nil {
		return CronJob{}, fmt.Errorf("%w: %v", ErrInvalidResponse, err)
	}
	if envelope.Job != nil {
		normalizeCronJob(envelope.Job)
		return *envelope.Job, nil
	}
	var job CronJob
	if err := json.Unmarshal(out, &job); err != nil || job.ID == "" {
		return CronJob{}, fmt.Errorf("%w: expected cron job object", ErrInvalidResponse)
	}
	normalizeCronJob(&job)
	return job, nil
}

func decodeJSONObject(out []byte) (map[string]any, error) {
	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("%w: expected JSON object: %v", ErrInvalidResponse, err)
	}
	return result, nil
}

func normalizeCronJob(job *CronJob) {
	if job == nil {
		return
	}
	if job.Status == "" {
		switch {
		case !job.Enabled:
			job.Status = "disabled"
		case job.State.RunningAtMs > 0:
			job.Status = "running"
		case job.State.LastRunStatus != "":
			job.Status = job.State.LastRunStatus
		case job.State.LastStatus != "":
			job.Status = job.State.LastStatus
		default:
			job.Status = "idle"
		}
	}
}

package runtime

import (
	"errors"
	"fmt"
)
type ExecutionStatus string

type ExecutionState struct {
	AgentID string 
	ExecutionID string 
	RetryCount int 
	CurrentAgentLoopTurn int //Gioi han 20 lan lap 
	ErrorMessage string 
	Status ExecutionStatus
	SessionID string //nen luu them de co gi co the ho tro them thong tin cho lan agent loop nay
}
const (
	ExecutionStatusInitializing      ExecutionStatus = "INITIALIZING"
	ExecutionStatusBuildingContext   ExecutionStatus = "BUILDING_CONTEXT"
	ExecutionStatusCallingProvider   ExecutionStatus = "CALLING_PROVIDER"
	ExecutionStatusExecutingTool     ExecutionStatus = "EXECUTING_TOOL"
	ExecutionStatusStreamingResponse ExecutionStatus = "STREAMING_RESPONSE"
	ExecutionStatusFinalizing        ExecutionStatus = "FINALIZING"

	ExecutionStatusCompleted    ExecutionStatus = "COMPLETED"
	ExecutionStatusFailed       ExecutionStatus = "FAILED"
	ExecutionStatusCancelled    ExecutionStatus = "CANCELLED"
	ExecutionStatusLimitReached ExecutionStatus = "LIMIT_REACHED"
)

// Day chinh la cac status ma agent duoc phep chuyen trang thia.
var ExecutionMap = map[ExecutionStatus]map[ExecutionStatus]bool{
	ExecutionStatusInitializing: {
		ExecutionStatusBuildingContext: true,
		ExecutionStatusFailed:          true,
		ExecutionStatusCancelled:       true,
	},
	ExecutionStatusBuildingContext: {
		ExecutionStatusCallingProvider: true,
		ExecutionStatusFailed:          true,
		ExecutionStatusCancelled:       true,
	},
	ExecutionStatusCallingProvider: {
		ExecutionStatusExecutingTool:     true,
		ExecutionStatusStreamingResponse: true,
		ExecutionStatusFailed:            true,
		ExecutionStatusLimitReached:      true,
		ExecutionStatusCancelled : true, 
	},
	ExecutionStatusExecutingTool: {
		ExecutionStatusCallingProvider: true,
		ExecutionStatusFailed:          true,
		ExecutionStatusLimitReached:    true,
		ExecutionStatusCancelled : true, 
	},
	ExecutionStatusStreamingResponse: {
		ExecutionStatusFinalizing: true,
		ExecutionStatusFailed:     true,
		ExecutionStatusCancelled:  true,
	},
	ExecutionStatusFinalizing: {
		ExecutionStatusCompleted: true,
		ExecutionStatusFailed:    true,
	},
}

func InvalidTransitionError(msg string) error {
	return errors.New(msg)
}

func (s *ExecutionState) Transform(next ExecutionStatus) error {
	currentStatus := s.Status
	transition, ok := ExecutionMap[currentStatus]
	if !ok {
		return InvalidTransitionError(
			fmt.Sprintf("Cannot change from %s to %s", s.Status, next),
		)
	}
	_, ok = transition[next] //Bien ok duoc khai bao lai nen khong the dung duoc :=, phai su dung dau =
	if !ok {
		return InvalidTransitionError(
			fmt.Sprintf("Cannot change from %s to %s", s.Status, next),
		)
	}
	return nil
}

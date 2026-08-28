package runtime

import (
	"errors"
	"fmt"
)

//Yeu cau cho moi lan chay agent
type ExecutionRequest struct {
	ExecutionID string 
	AgentID string 
	SessionID string 
	UserID string 

	Prompt string 
	Model string 
	Provider string 

	Steam bool 
	Metadata map[string]string 
} 
type ExecutionStatus string 
const (
	ExecutionStatusInitializing ExecutionStatus = "INITIALIZING" 
	ExecutionStatusBuildingContext ExecutionStatus = "BUILDING_CONTEXT" 
	ExecutionStatusCallingProvider ExecutionStatus = "CALLING_PROVIDER" 
	ExecutionStatusExecutingTool ExecutionStatus = "EXECUTING_TOOL" 
	ExecutionStatusStreamingResponse ExecutionStatus = "STREAMING_RESPONSE"
	ExecutionStatusFinalizing ExecutionStatus = "FINALIZING"

	ExecutionStatusCompleted ExecutionStatus = "COMPLETED" 
	ExecutionStatusFailed ExecutionStatus = "FAILED" 
	ExecutionStatusCancelled ExecutionStatus = "CANCELLED" 
	ExecutionStatusLimitReached ExecutionStatus = "LIMIT_REACHED" 
) 

type ExecutionState struct {
	ExecutionID string 
	AgentID string 
	SessionID string 

	CurrentModelTurn int 
	TotalCall int 
	TotalRetry int 

	Status ExecutionStatus 
	ErrorMessage string 
}

//Khoi tao cac luong ma agent se duoc phep goi. Agent chi co the thuc hien viec chuyen qua lai giua cac trang thai nay 

var ExecutionMap = map[ExecutionStatus]map[ExecutionStatus]bool{
	ExecutionStatusInitializing: {
		ExecutionStatusBuildingContext: true, 
		ExecutionStatusFailed: true, 
		ExecutionStatusCancelled: true, 
	}, 
	ExecutionStatusBuildingContext: {
		ExecutionStatusCallingProvider: true, 
		ExecutionStatusFailed: true, 
		ExecutionStatusCancelled: true, 
	}, 
	ExecutionStatusCallingProvider: {
		ExecutionStatusExecutingTool: true, 
		ExecutionStatusStreamingResponse: true, 
		ExecutionStatusFailed: true, 
		ExecutionStatusLimitReached: true, 
	},  
	ExecutionStatusExecutingTool: {
		ExecutionStatusCallingProvider: true, 
		ExecutionStatusFailed: true, 
		ExecutionStatusLimitReached: true, 
	}, 
	ExecutionStatusStreamingResponse: {
		ExecutionStatusFinalizing: true, 
		ExecutionStatusFailed: true, 
		ExecutionStatusCancelled: true, 
	}, 
	ExecutionStatusFinalizing: {
		ExecutionStatusCompleted : true,
		ExecutionStatusFailed: true,  
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
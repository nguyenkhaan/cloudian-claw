package runtime

//Dai dien cho 1 yeu cau client gui den Agent cua chung ta -> SessionID, AgentID, UserID, Prompt
type ExecutionRequest struct {
	ExecutionID string 
	UserID string 
	AgentID string 
	SessionID string 
	Prompt string 
	Model string 
	Provider string 
	Stream bool //Agent se chay dang stream hay la dang response hoan chinh 
}

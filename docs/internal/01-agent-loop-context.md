# AGENT LOOP 
## 1. Luồng hoạt động 
Hệ thống Agent của chúng ta sẽ có một luồng chạy như sau: 

Gateway -> Agent Router -> Agent Loop -> (Provider, Tool Registry, ...) 

- Gateway: Đóng vai trò là cửa vào của hệ thống 
    - Xác nhận gateway token 
    - Kiểm tra request cơ bản 
    - Giới hạn tốc độ hoặc quota 
    - Từ chối request không hợp lệ 

- Agent Router: Bộ định tuyến tới đúng Agent / Session 
    - Tìm Agent theo ID hoặc tên 
    - Kiểm tra Agent có đang tồn tại và đang active 
    - Kiểm tra session 
    - Tạo `ExecutionRequest` dựa vào các thông tin trên để truyền vào Agent Loop 
- Agent Loop: runtime trung tâm 
    - Xây dựng `RuntimeContext` 
    - Nạp các skill, history, rule 
    - Gọi Provider tới các LLM 
    - Thực hiện vòng lặp: loop (<= 20times): THINK -> ACT -> OBSERVE 
    - Gọi `ToolCall` và xử lý result `ToolError` và `ToolResult` 

Các interface làm nhiệm vụ như các contract hoặc các đối tượng / sự kiện phải xử lý trong Agent Loop. 

- ExecutionRequest: | ExecutionID, AgentID, SessionID, UserID, Prompt, Model, Provider, Stream, Metadata | Yêu cầu mỗi khi bắt đầu một lần chạy Agent 
- RunTimeContext: | AgentID, SessionID, UserID, Workspace, RestrictToWorkspace, SystemPrompt, History, Skills, Rules, Tools, Model, Provider, MaxToolCalls, MaxProviderRetries, MaxContextTokens, MaxContinuationDepth | Gom toàn bộ dữ liệu Agent cần trong một lần chạy. Hạn chế việc để Agent đọc database trong bước này 
- ExecutionState | | Thể hiện chi tiết từng bước trạng thái của quá trình Execution. Tại State này thì các thông số của Execution đang là gì? 


### Execution State 
Với ExecutionState. Chúng ta cần khai báo thêm các state trung gian. Đóng vai trò như định danh để biết bước hiện tại của Execution 
```go 
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
```

Các bước chuyển state phải được định nghĩa trước và phải được cân nhắc tính hợp lệ, mỗi khi Agent muốn thực hiện việc chuyển state. Chúng ta sẽ có 1 bản đồ chuyển state cung với 1 hàm để Agent có thể xem xét việc chuyển trạng thái có hợp lệ không 
```go
ExecutionStateMap := map[ExecutionStatus]map[ExecutionStatus]bool{}  
```
```go
func (s *ExecutionState) Transition(next ExecutionStatus) error /  
```

### ToolCall, ToolResult, ToolError 
Đây là các runtime contract chứa input / output, sẽ dùng trong quá trình thực thi. Cần phân biệt được với các entity model khác 
```go
type ToolCall struct {
    ID        string
    ToolID    string
    ToolName  string
    Arguments map[string]any
}
```

```go
type ToolResult struct {
    ToolCallID string
    Content    string
    Output     any
    IsError    bool
}
```
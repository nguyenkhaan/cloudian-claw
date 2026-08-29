# AGENT LOOP 
## 1. Luồng hoạt động của hệ thống 
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
## 2. Agent Loop 
### Khái niệm 
When a model receives a request, it cannot fully address the problem, it needs to reach out the outside: read files, query databases, call APIs, execute codes... These things help the agent verify or judges its work quality.

A LLM model or Agent cannot do this (because Agent can only do things once that the LLM request). The Agent Loop is what makes different 

Agent Loop is an orchestration layer enable this. It manages the cycle or reasoning and action, allows model to tackle problem requiring multiple steps. 

### How to the Agent Loop Work 
The Agent Loop operates on a simple rule: invoked the model -> check if wants to use a tools -> execute the tool, then invoked the model again with the result... until we have final answer

**For example:** We have to analyze codebase for security vulnerabilities: 

1. The model receives the request. Firstly, it needs to know the folder structure. It requests a file listing tool 
2. The model now sees the directory structure in its context. It identifies the main application entry point and requests the file reader tool to examine it.
3. The model sees the application code. It notices database queries and decides to examine the database module for potential SQL injection. It requests the file reader again.
4. The model sees the database module and identifies a vulnerability: user input concatenated directly into SQL queries. To assess the scope, it requests a code search tool to find all call sites of the vulnerable function.
5. The model sees 12 call sites in the search results. It now has everything it needs. Rather than requesting another tool, it produces a terminal response: a report detailing the vulnerability, affected locations, and remediation steps.

Each phase: from Invoked Model -> Seed data to context, we call a `depth`. We have to define `MAXIMUM_CONTINUATION_DEPTH` to limit the depth in agent loop, prevent agent from running infinitely. 

In each phase, Agent will run 3 mini step: Think (Reasoning with LLM) -> Act (Execute Tool) -> Observe (Update Context)

For example: 
```
Depth N: 
    -> Think: invoked model and receive answer 
    -> Act: If model require tool call: 
        + Execute tool 
    -> Observe: Give tool result to context 
    Continue Depth N + 1 
Depth (N + 1): 
    -> 
    -> 
    -> 
```
### System Workflows 

Our system acts following flows: 
    Gateway -> Execution Request -> Agent Requirement Context (Tools, Rules, Skills...) Preparation -> Agent Loop (execute ToolCall, update state per iteration) -> Finialize answer (Execution Result) 

We have a runtime layer, contains multiple structs to serve agent loop pipeline: `ExecutionRequest`, `ExecutionState`, `ExeuctionStatus`, `TooCall`, `ExecutionResult`, `ModelResponse`

Explain structs table: 


### 
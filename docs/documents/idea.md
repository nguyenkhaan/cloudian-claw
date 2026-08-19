# AI Agent Gateway System — Kiến trúc lõi (Local, rút gọn từ Goclaw)

> Hệ thống AI Agent chạy local, tập trung vào các thành phần cốt lõi: Agent, Skill, Memory, Tool, Provider, Session, Rules.
> **Đã loại bỏ khỏi phạm vi:** Channel, MCP, Subagent, Delegation, Teams, TTS, Webhook, Cron — các tính năng mở rộng/phối hợp đa tác nhân không thuộc lõi một Agent system.
> Kiến trúc kỹ thuật chi tiết (tech stack, sơ đồ component, data flow BE) được tách sang tài liệu System Design riêng.

---

## Blackbox Exploring

**Mục tiêu**: Kiểm kê từ ngoài vào trong — chức năng nào của Goclaw được giữ lại cho bản lõi này, chia theo BE (service/dữ liệu) và FE (màn hình/tương tác).

### Backend (BE)

| Nhóm | Chức năng giữ lại | Ghi chú |
|---|---|---|
| Auth & Gateway | Xác thực bằng Gateway API Token | Không có khái niệm Team/Org — 1 instance local, 1 chủ sở hữu. Token được cấu hình tĩnh trong file .env và lưu tại LocalStorage phía Client, xác thực qua header Authorization: Bearer. |
| Provider | Lưu & xác thực API key provider (OpenRouter, Vercel...), khám phá model theo provider | Nhiều provider cùng lúc, không giới hạn 1 model cố định |
| Agent | CRUD agent (SOUL.md, IDENTITY.md), clone agent làm template | Không có Subagent — mỗi agent độc lập, không phân cấp cha/con |
| Skill | Load, validate, version, gán skill cho agent | Skill lấy từ filesystem local, không có kho skill từ xa/marketplace. Skill thuần túy là tài liệu hướng dẫn bằng ngôn ngữ tự nhiên (Markdown + YAML frontmatter), không chứa mã thực thi; nhiệm vụ thực thi thuộc về Tool Registry. |
| Tool | Đăng ký tool, kiểm soát quyền truy cập tool theo agent, log thực thi | Không có Delegation (agent giao việc cho agent khác qua tool) |
| Rules | CRUD rule, gán rule cho agent, versioning, sandbox test | Gán theo từng agent, không có khái niệm "nhóm agent"/Team |
| Session | Tạo, lưu, khôi phục session; lịch sử hội thoại | — |
| Memory | Cấu hình embedding (provider, model, chunk size, overlap), nén hội thoại (compaction) | Lưu trữ vector dạng mảng số thực (double precision[] hoặc jsonb) trong PostgreSQL và tính độ tương đồng (Cosine Similarity) trực tiếp trong Go để đảm bảo tính di động cục bộ. |
| Monitoring | Thu thập trace, latency, token usage, error rate | — |
| Realtime Events | Phát sự kiện real-time (chat message, reasoning step, tool call, skill activation) | Chỉ phục vụ hiển thị nội bộ, không phải Webhook ra ngoài |
| API Key | Tạo/thu hồi API key có scope, xem usage | Dùng để agent hoặc client ngoài gọi vào hệ thống — không phải Channel tích hợp nền tảng thứ ba |
| Config | Server, AI Defaults, Quota, Tools (cấu hình toàn hệ thống) | Quota là per-agent/per-user cá nhân, không phải per-team |

**Loại bỏ khỏi BE**: Channel adapter (Slack/Discord/Telegram...), MCP protocol handler, Subagent orchestration, Delegation (agent-to-agent task handoff), Team/Org multi-tenancy, TTS engine, Webhook dispatcher, Cron scheduler.

### Frontend (FE)

Giao diện FE sẽ gồm các màn hình chủ yếu sau: 

| Nhóm | Chức năng giữ lại | Ghi chú |
|---|---|---|
| Onboarding Setup | Nhập Gateway token → cấu hình provider → validate model → tạo agent đầu tiên | Thực hiện setup một số thông tin quan trọng. Luồng bắt buộc trước khi sử dụng hệ thống |
| Overview Dashboard | Trạng thái agent, sức khỏe provider, hoạt động gần đây, resource usage | — |
| Chat Interface | Chat real-time với agent, streaming, xem tool call, đính kèm file/ảnh | — |
| Agent Management | Tạo/sửa/xóa/clone agent, xem chi tiết (model, skill, tool, memory stats) | — |
| Session Management | List, rename, archive/xóa, chuyển session | — |
| Rules Management | CRUD rule, gán rule, sandbox test | — |
| Skill Management | Browse thư viện skill, thêm skill custom, bật/tắt theo agent, version | — |
| Tool Management | Bật/tắt tool, cấu hình tham số, policy truy cập, xem log | — |
| Provider Management | Thêm/xóa provider, nhập & validate key, xem model, đặt provider mặc định | — |
| API Key Management | Tạo key có scope, đặt hạn dùng, thu hồi, xem usage | — |
| Memory Management | Xem memory usage/agent, cấu hình embedding, ngưỡng compaction | — |
| Monitoring | Trace theo agent, latency, token, error rate | — |
| Realtime Event Feed | Thẻ sự kiện real-time (chat, reasoning, tool call, skill activation) | — |
| Terminal CLI | Chat mode tương tác, lệnh `/agent /skill /tool`, tự động phát hiện skill/rule, coding assistant | Không có lệnh liên quan Channel/MCP/Team |
| Config Page | Tab Server / AI Defaults / Quota / Tools | — |

**Loại bỏ khỏi FE**: màn hình Channel integration, MCP server config, Subagent tree view, Delegation flow, Team/Member management, TTS voice settings, Webhook config, Cron job scheduler.

#### Mô tả chi tiết cho từng màn hình: 
- Overview Setup: Gồm 3 bước thiết lập: Provider -> Model -> Agent.
-  
---

## User Usecase

Mỗi usecase trình bày dạng luồng thao tác đánh số + mô tả ngắn.

### UC-01: Onboarding lần đầu
**Mô tả**: Người dùng thiết lập hệ thống trước khi sử dụng.
1. Nhập Gateway API Token
2. Chọn một provider và nhập API Token của Provider đó (Ví dụ OpenRouter)
3. Hệ thống validate key → hiển thị danh sách model khả dụng
4. Người dùng chọn model mặc định
5. Tạo Agent đầu tiên: điền SOUL.md (tính cách, chuyên môn), IDENTITY.md dựa trên LLM Model đó. 
6. Hệ thống khởi tạo agent → sao khi khởi tạo thành công, chuyển người dùng vào dashboard. 

### UC-02: Tạo Agent mới
**Mô tả**: Tạo một agent với cấu hình riêng.
1. Mở Agent Management → Create Agent
2. Điền SOUL.md (identity, tone, chuyên môn)
3. Chọn Provider + Model
4. Bật/tắt các Skill muốn gán
5. Bật/tắt các Tool được phép dùng
6. Gán Rule hành vi (nếu có)
7. Submit → Agent được khởi tạo, sẵn sàng chat

### UC-03: Chat với Agent
**Mô tả**: Tương tác hội thoại với agent đã tạo.
1. Chọn agent từ danh sách hoặc mở session cũ
2. Nhập tin nhắn (kèm file/ảnh nếu cần). File đính kèm sẽ được upload lên thư mục cục bộ của backend, lưu filepath trong DB và hiển thị qua static file endpoint, thiết kế trừu tượng để sau này dễ chuyển đổi sang Cloudinary.
3. Hệ thống stream phản hồi theo thời gian thực
4. Nếu agent gọi tool → hiển thị tool call + kết quả trong panel
5. Session được lưu tự động sau mỗi lượt

### UC-04: Quản lý Skill cho Agent
**Mô tả**: Bật/tắt hoặc cập nhật skill cho một agent cụ thể.
1. Mở Skill Management
2. Browse thư viện skill (local)
3. Chọn agent cần gán
4. Bật/tắt skill tương ứng
5. Hệ thống validate dependency giữa các skill
6. Lưu cấu hình → agent áp dụng ngay ở lượt chat kế tiếp

### UC-05: Quản lý Tool & quyền truy cập
**Mô tả**: Kiểm soát agent được dùng tool nào.
1. Mở Tool Management
2. Chọn agent
3. Bật/tắt từng tool
4. Cấu hình tham số tool (nếu có)
5. Lưu → Tool Registry cập nhật policy áp dụng cho agent đó

### UC-06: Thiết lập Rules cho Agent
**Mô tả**: Định nghĩa ràng buộc hành vi.
1. Mở Rules Management → Create Rule
2. Viết nội dung rule
3. Test rule trong sandbox (Thực hiện cuộc gọi LLM thực tế qua mock session với tin nhắn mẫu và rule thử nghiệm để đánh giá phản hồi, không lưu vào DB Session thực tế)
4. Gán rule cho agent
5. Lưu → rule có hiệu lực ở lượt chat kế tiếp

### UC-07: Quản lý Session
**Mô tả**: Theo dõi và thao tác trên các phiên hội thoại.
1. Mở Session Management
2. Xem danh sách session (agent dùng, thời lượng, số tin nhắn)
3. Chọn một session để đổi tên / archive / xóa / mở lại
4. Nếu mở lại → quay về Chat Interface với ngữ cảnh cũ

### UC-08: Cấu hình Memory / Embedding
**Mô tả**: Điều chỉnh cách agent lưu và truy xuất ký ức dài hạn.
1. Mở Agent Memory Management
2. Chọn agent
3. Chọn provider/model embedding
4. Chỉnh chunk size, overlap
5. Chỉnh ngưỡng compaction (threshold, keep recent, max token)
6. Lưu → áp dụng cho các hội thoại tiếp theo của agent

### UC-09: Theo dõi Monitoring
**Mô tả**: Giám sát hiệu năng và lỗi của agent.
1. Mở Monitoring
2. Chọn agent hoặc xem toàn hệ thống
3. Xem trace, latency, token usage, error rate
4. Drill-down vào một trace cụ thể để debug

### UC-10: Dùng Terminal CLI
**Mô tả**: Tương tác agent qua dòng lệnh cho tác vụ lập trình.
1. Chạy CLI, vào chế độ chat tương tác
2. Gõ lệnh `/agent <tên>` để chọn agent
3. CLI tự động quét và áp dụng skill/rule liên quan
4. Gửi yêu cầu (coding, troubleshooting...)
5. Agent đọc/ghi file cục bộ theo tool được cấp quyền
6. Nếu cần cấu hình phức tạp → CLI gợi ý chuyển sang Dashboard

### UC-11: Tạo API Key cho tích hợp ngoài
**Mô tả**: Cho phép client/agent ngoài gọi vào hệ thống.
1. Mở API Key Management → Create Key
2. Chọn scope quyền
3. Đặt ngày hết hạn (tùy chọn)
4. Hệ thống sinh key → hiển thị một lần
5. Theo dõi usage của key trong danh sách

### UC-12: Cấu hình hệ thống
**Mô tả**: Điều chỉnh thiết lập toàn cục.
1. Mở Config Page
2. Chọn tab: Server / AI Defaults / Quota / Tools
3. Chỉnh giá trị tương ứng (host, port, model mặc định, temperature, quota, tool toggle)
4. Lưu → áp dụng toàn hệ thống

---

## System Usecase

Chức năng hệ thống phải thực hiện để phục vụ các User Usecase ở trên.

| ID | System Usecase | Phục vụ UC | Mô tả |
|---|---|---|---|
| SU-01 | Xác thực & định tuyến request | UC-01, UC-03, UC-10 | Gateway xác thực token, validate request, định tuyến tới đúng service (không qua Channel adapter) |
| SU-02 | Đăng ký & khám phá Provider | UC-01, UC-02, UC-12 | Lưu key, gọi API provider để validate, liệt kê model khả dụng, cache kết quả |
| SU-03 | Vòng đời Agent | UC-02, UC-04, UC-05, UC-06 | Tạo/sửa/xóa/clone agent; ràng buộc agent với model, skill, tool, rule đã chọn |
| SU-04 | Nạp & xác thực Skill | UC-04 | Đọc skill từ filesystem, kiểm tra cú pháp/dependency, giữ bản hoạt động cuối nếu bản mới lỗi |
| SU-05 | Đăng ký & thực thi Tool có kiểm soát quyền | UC-05, UC-10 | Kiểm tra policy trước khi cho agent gọi tool; ghi log mọi lần gọi |
| SU-06 | Áp dụng Rules trong pipeline sinh phản hồi | UC-06 | Chèn rule đã gán vào ngữ cảnh trước khi gọi model; hỗ trợ test rule độc lập trong sandbox |
| SU-07 | Quản lý & khôi phục Session | UC-03, UC-07 | Lưu message pair vào DB, khôi phục session từ checkpoint gần nhất nếu lỗi giữa chừng |
| SU-08 | Pipeline Embedding & Compaction | UC-08 | Sinh embedding theo cấu hình, truy xuất context liên quan, nén hội thoại khi vượt ngưỡng |
| SU-09 | Thu thập Telemetry & Trace | UC-09 | Ghi metric, trace span, error event theo từng agent; buffer khi observability suy giảm |
| SU-10 | Phát sự kiện Realtime | UC-03, UC-09 | Emit event (message, reasoning step, tool call, skill activation) qua WebSocket/SSE tới FE |
| SU-11 | Cấp phát & kiểm soát API Key | UC-11 | Sinh key có scope, kiểm tra hạn dùng khi request tới, ghi nhận usage |
| SU-12 | Quản lý cấu hình toàn cục | UC-12 | Lưu & áp dụng Server / AI Defaults / Quota / Tools; các service đọc config này khi khởi động và khi thay đổi |
| SU-13 | Streaming phản hồi model + tool loop có giới hạn | UC-03, UC-10 | Gọi LLM provider, xử lý tool call, đưa kết quả tool trở lại model; giới hạn độ sâu đệ quy (không vòng lặp agent vô hạn, không Subagent/Delegation) |

---

## Ghi chú phạm vi

- **System Design** (tech stack, sơ đồ kiến trúc BE, đặc tả component, data flow diagram) nằm ở tài liệu riêng, sẽ được cung cấp sau — không lặp lại ở đây.
- Danh sách loại bỏ áp dụng xuyên suốt cả 3 section: **Channel, MCP, Subagent, Delegation, Teams, TTS** — cùng các tính năng phụ thuộc trực tiếp vào chúng (Webhook, Cron, quota theo team...).
- Mọi chức năng còn lại trong file gốc (SOUL.md/IDENTITY.md, Skill, Tool, Memory, Provider, Session, Rules, Monitoring, Realtime Events, API Key, Config, CLI cơ bản) được giữ nguyên vì thuộc nhóm lõi Agent system.
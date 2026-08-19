# PRD: Cloudclaw Product Overview

## 1. Giới thiệu tổng quan (Introduction / Overview)
**Cloudclaw** là hệ thống AI Agent Gateway cục bộ (local-first, single-owner) hoạt động như một bệ phóng (Agent Harness) để chạy và kiểm thử các Mô hình ngôn ngữ lớn (LLM) dưới vai trò các AI Agent độc lập. Hệ thống hỗ trợ đắc lực cho các tác vụ lập trình (coding assistant) và xử lý sự cố. 

Hệ thống được thiết kế xoay quanh trung tâm là **Agent Loop** (vòng điều phối thực thi), kết hợp với các dịch vụ chuyên biệt khác như Quản lý Provider, Đăng ký Tool, Lưu trữ Memory, Session và quản lý Skill/Rule để định hướng hành vi của Agent.

---

## 2. Mục tiêu (Goals)
* **Khởi chạy Cục bộ dễ dàng:** Không yêu cầu các thiết lập môi trường phức tạp (không bắt buộc cài đặt extension `pgvector` trên PostgreSQL, tối giản hóa việc cài đặt offline).
* **Quản lý Agent linh hoạt:** Cấu hình Agent chi tiết bằng `SOUL.md` / `IDENTITY.md` cùng các ràng buộc về mặt hành vi (Rules) và khả năng (Skills).
* **Vòng lặp thực thi an toàn (Bounded Agent Loop):** Đảm bảo LLM có thể gọi các Tool hệ thống một cách bảo mật và có giới hạn chống lặp vô hạn.
* **Tương tác thời gian thực:** Stream kết quả và theo dõi chi tiết quá trình suy nghĩ, kích hoạt skill hay gọi tool của agent thông qua Web Dashboard hoặc Terminal CLI.

---

## 3. User Stories (Các câu chuyện người dùng)

### US-001: Onboarding thiết lập hệ thống lần đầu
**Description:** Là người dùng mới, tôi muốn thực hiện thiết lập hệ thống từng bước (nhập Gateway token, cấu hình provider và tạo agent đầu tiên) để tôi có thể nhanh chóng bắt đầu sử dụng hệ thống.
**Acceptance Criteria:**
- [ ] Giao diện Onboarding yêu cầu nhập Gateway API Token (xác thực trực tiếp với token cấu hình trong `.env`).
- [ ] Người dùng cấu hình Provider (nhập API key) -> hệ thống validate key thành công và lấy danh sách model khả dụng.
- [ ] Người dùng điền file cấu hình `SOUL.md` và `IDENTITY.md` để khởi tạo Agent đầu tiên.
- [ ] Lưu Gateway Token thành công vào LocalStorage của client để xác thực các request sau này.
- [ ] Sau khi thiết lập xong, chuyển hướng thành công vào Dashboard chính.
- [ ] Verify in browser using dev-browser skill.

### US-002: Tạo và cấu hình Agent mới
**Description:** Là người dùng, tôi muốn tạo hoặc clone một agent mới với các cấu hình SOUL, Model, Skill, Rule và Tool riêng biệt để đáp ứng các nghiệp vụ chuyên biệt.
**Acceptance Criteria:**
- [ ] Có giao diện form nhập đầy đủ thông tin: tên Agent, `SOUL.md`, `IDENTITY.md`, chọn Provider/Model, bật/tắt Skill, Rule và phân quyền Tool.
- [ ] Tính năng clone Agent cho phép sao chép nhanh cấu hình từ agent có sẵn sang agent mới.
- [ ] Mọi cập nhật cấu hình của Agent chỉ có hiệu lực ở lượt thực thi kế tiếp (không ảnh hưởng tới Runtime Context đang chạy).
- [ ] Verify in browser using dev-browser skill.

### US-003: Chat thời gian thực với Agent và xem Tool Call
**Description:** Là người dùng, tôi muốn nhắn tin với Agent, gửi file đính kèm, nhận phản hồi dạng streaming và xem các bước gọi tool thời gian thực để theo dõi và tương tác hiệu quả.
**Acceptance Criteria:**
- [ ] Hiển thị khung chat hỗ trợ stream tin nhắn (SSE / WebSocket) ổn định.
- [ ] Hỗ trợ tải tệp đính kèm (ảnh, tài liệu). Tệp được lưu tạm thời tại thư mục cục bộ của backend (`data/attachments/`), lưu relative path trong DB và phục vụ qua static file endpoint.
- [ ] Giao diện hiển thị các sự kiện gọi tool (Tool Call) trực quan: Đang gọi, Tham số, Trạng thái (Thành công/Lỗi).
- [ ] Bất kỳ lỗi thực thi tool nào đều được hiển thị lỗi và trả ngược lại cho model xử lý tiếp thay vì crash hệ thống.
- [ ] Verify in browser using dev-browser skill.

### US-004: Quản lý thư viện Skill cục bộ
**Description:** Là người dùng, tôi muốn quét và duyệt danh sách các Skill trên ổ đĩa local để liên kết chúng với các Agent mong muốn.
**Acceptance Criteria:**
- [ ] Hiển thị danh sách các Skill quét được từ filesystem.
- [ ] Skill được nạp thành công từ các tệp Markdown kèm YAML frontmatter để lấy metadata (không chứa code thực thi).
- [ ] Nội dung của Skill được inject trực tiếp vào Agent System Prompt khi build Runtime Context.
- [ ] Cho phép bật/tắt và gán Skill cho từng Agent từ Dashboard.
- [ ] Verify in browser using dev-browser skill.

### US-005: Quản lý Tool và cấu hình phân quyền
**Description:** Là người dùng, tôi muốn bật/tắt các Tool ở mức toàn cục hoặc phân quyền sử dụng Tool cho từng Agent để đảm bảo an toàn hệ thống.
**Acceptance Criteria:**
- [ ] Giao diện quản lý danh sách các Tool hệ thống (ví dụ: đọc file, chạy lệnh terminal...).
- [ ] Cho phép bật/tắt Tool toàn cục (Global Tool Enabled).
- [ ] Cho phép tick chọn danh sách Tool được phép sử dụng cho mỗi Agent (Agent Tool Allowed).
- [ ] Ghi lại lịch sử (log) thực thi của từng Tool.
- [ ] Verify in browser using dev-browser skill.

### US-006: Tạo Rule và chạy thử trong Sandbox
**Description:** Là người dùng, tôi muốn tạo các quy tắc ứng xử (Rule), gán cho Agent và chạy thử trong môi trường Sandbox để xem LLM phản ứng như thế nào trước khi lưu chính thức.
**Acceptance Criteria:**
- [ ] Giao diện Sandbox cho phép nhập tin nhắn mẫu và Rule thử nghiệm.
- [ ] Khi chạy Sandbox, hệ thống tạo mock session, gọi LLM thực tế bằng provider đang chọn và trả ra kết quả để người dùng đánh giá.
- [ ] Cuộc gọi thử nghiệm không ghi lịch sử vào database Session chính thức.
- [ ] Verify in browser using dev-browser skill.

### US-007: Quản lý phiên hội thoại (Session)
**Description:** Là người dùng, tôi muốn xem danh sách các cuộc hội thoại cũ, đổi tên, lưu trữ (archive), xóa hoặc tiếp tục trò chuyện ở phiên đó.
**Acceptance Criteria:**
- [ ] Hiển thị danh sách Session phân chia theo Agent.
- [ ] Hỗ trợ các nút thao tác: đổi tên, lưu trữ, xóa.
- [ ] Khi mở lại một Session cũ, toàn bộ lịch sử trò chuyện được phục hồi chính xác (bao gồm cả trình tự gọi tool trước đó).
- [ ] Verify in browser using dev-browser skill.

### US-008: Cấu hình Memory và Compaction
**Description:** Là người dùng, tôi muốn thiết lập thông số Embedding cho Memory và kiểm soát ngưỡng tự động nén hội thoại để tối ưu ngữ cảnh.
**Acceptance Criteria:**
- [ ] Cấu hình được: Embedding provider, Model, chunk size, overlap.
- [ ] Thiết lập ngưỡng compaction: keep recent messages, max token budget.
- [ ] Hỗ trợ nút kích hoạt nén thủ công (Manual Compaction) trực tiếp trên Session.
- [ ] Verify in browser using dev-browser skill.

### US-010: Tương tác qua Terminal CLI
**Description:** Là nhà phát triển thích dùng dòng lệnh, tôi muốn sử dụng CLI để chat với Agent và thao tác nhanh các tính năng mà không cần mở trình duyệt.
**Acceptance Criteria:**
- [ ] CLI cung cấp lệnh `/agent <tên>` để chọn agent, `/session` để quản lý phiên và chế độ chat trực tiếp.
- [ ] Tích hợp cơ chế stream phản hồi trực tiếp trên terminal.
- [ ] Hiển thị rõ ràng các sự kiện gọi Tool đang diễn ra.

### US-011: Quản lý API Key cho ứng dụng ngoài
**Description:** Là người dùng, tôi muốn tạo các API Key có phân quyền (scope) và thời hạn riêng để các ứng dụng bên ngoài có thể gọi Agent thông qua Gateway.
**Acceptance Criteria:**
- [ ] Cho phép tạo API Key kèm cấu hình scope và ngày hết hạn.
- [ ] Chỉ hiển thị Plain API Key duy nhất 1 lần khi khởi tạo.
- [ ] Hỗ trợ thu hồi (revoke) key và xem lịch sử usage (số lần gọi).
- [ ] Verify in browser using dev-browser skill.

### US-012: Cấu hình toàn cục (Global Configuration)
**Description:** Là người quản trị, tôi muốn cấu hình các tham số mặc định của Server và AI để áp dụng chung cho toàn hệ thống.
**Acceptance Criteria:**
- [ ] Cấu hình thông số Server (Host, Port).
- [ ] Cấu hình AI Defaults (Default Provider, Default Model, Temperature, Max Tokens).
- [ ] Thiết lập hạn mức mặc định (Quota per User/Agent).
- [ ] Verify in browser using dev-browser skill.

---

## 4. Yêu cầu chức năng (Functional Requirements)
* **FR-1 [Xác thực]:** Hệ thống Gateway bắt buộc kiểm tra Gateway API Token qua header `Authorization: Bearer <token>`. Mọi request không khớp với giá trị cấu hình tĩnh trong file `.env` đều bị reject với mã lỗi 401.
* **FR-2 [Provider Discovery]:** Hệ thống phải cung cấp cổng trừu tượng gọi API của các LLM Provider tương thích với định dạng OpenAI. Có API kiểm tra (validate) key và cache danh sách các model của provider đó.
* **FR-3 [Vòng đời Agent]:** Quản lý cấu hình Agent trong DB và thư mục cục bộ. Khi thay đổi cấu hình Agent, hệ thống phải đảm bảo phiên thực thi đang chạy (Execution) sử dụng bản snapshot ngữ cảnh bất biến và không bị ảnh hưởng.
* **FR-4 [Nạp Skill]:** Parser của hệ thống Go sẽ đọc và phân tích cấu trúc YAML frontmatter trong tệp Markdown của các Skill. Trích xuất chính xác thông tin hướng dẫn và inject vào Runtime Context.
* **FR-5 [Tool Registry]:** Đóng vai trò là chốt chặn an toàn duy nhất để kiểm tra quyền và thực thi các Tool hệ thống (chạy terminal, thao tác file). Đảm bảo kiểm tra đúng giới hạn workspace và timeout.
* **FR-6 [Rule Engine]:** Khi dựng Runtime Context, hệ thống sẽ gộp tất cả Rule được gán cho Agent thành các chỉ thị ràng buộc bổ sung đưa vào System Prompt của Provider.
* **FR-7 [Session Persistence]:** Lưu trữ tuần tự toàn bộ lịch sử tin nhắn của phiên trò chuyện, bao gồm cả thứ tự gọi tool và kết quả của tool đó. Lỗi ghi DB không được rollback luồng stream tin nhắn đang trả về cho người dùng.
* **FR-8 [Memory Vector Search]:** Thực hiện tính toán Cosine Similarity trực tiếp bằng Go trên tập dữ liệu mảng số thực float8/JSONB được lưu trong PostgreSQL để tìm kiếm ngữ cảnh dài hạn phù hợp.
* **FR-9 [Observability]:** Ghi log chi tiết và trace span cho mỗi Execution. Các spans con phải liên kết chặt chẽ với root trace của Execution.
* **FR-10 [Realtime Publisher]:** Phát hành các sự kiện trạng thái (execution events) dạng JSON qua WebSocket/SSE. Lỗi đường truyền không làm gián đoạn Agent Loop.
* **FR-11 [API Key Auth]:** Gateway hỗ trợ xác thực API Key từ client ngoài, đối chiếu scope của key và ghi nhận lịch sử usage.
* **FR-12 [Global Configuration]:** Lưu trữ và áp dụng cài đặt toàn cục vào DB, cho phép Agent cấu hình ghi đè (override) các mặc định này.
* **FR-13 [Bounded Execution]:** Giới hạn Agent Loop tối đa số lần gọi tool liên tục (max tool calls), độ sâu đệ quy (max continuation depth) và thời gian thực thi (timeout) để ngăn chặn vòng lặp vô hạn.

---

## 5. Non-Goals (Các mục tiêu ngoài phạm vi)
* Không xây dựng hệ thống tài khoản nhiều người dùng/tổ chức (Multi-tenant/Multi-org).
* Không có cơ chế phối hợp đa tác nhân (Multi-agent orchestration), các Agent hoạt động hoàn toàn độc lập (không có Subagent hay Delegation).
* Không xây dựng kho lưu trữ skill trực tuyến (Skill Marketplace).
* Không tích hợp kênh thứ ba (Slack, Discord, Telegram...).
* Không hỗ trợ giao thức MCP (Model Context Protocol).
* Không chuyển đổi văn bản thành giọng nói (TTS) hay Webhook bên ngoài.

---

## 6. Thiết kế kỹ thuật và Kiến trúc (Technical Considerations)
* **Ngôn ngữ & Framework:** Backend viết bằng **Go (Golang)**; Frontend viết bằng **React + TypeScript + TailwindCSS / Vanilla CSS**.
* **Database:** **PostgreSQL** là nguồn lưu trữ bền vững. Thay vì dùng `pgvector`, tầng Go sẽ chịu trách nhiệm tính toán Cosine Similarity trên các vector dạng mảng số thực để tăng tính di động.
* **Cơ chế tải file:** File đính kèm lưu tạm thời trên thư mục cục bộ của backend. Tầng lưu trữ được trừu tượng hóa qua interface để sau này dễ dàng chuyển đổi sang sử dụng **Cloudinary** hoặc dịch vụ lưu trữ đám mây tương đương.

---

## 7. Chỉ số thành công (Success Metrics)
* **Khởi động nhanh:** Cài đặt cục bộ đơn giản, tự động chạy migration và khởi động server Go thành công dưới 3 giây.
* **Stream mượt mà:** Thời gian từ lúc nhận request đến khi bắt đầu stream ký tự đầu tiên (Time To First Token) dưới 800ms trên đường truyền ổn định.
* **Độ trễ thấp cho tìm kiếm vector:** Thuật toán Cosine Similarity viết bằng Go xử lý tìm kiếm trong database dưới 15ms cho quy mô 1000 chunks ký ức.

---

## 8. Câu hỏi mở (Open Questions)
1. Khi chuyển sang Cloudinary, chúng ta có giữ cơ chế fallback lưu trữ cục bộ nếu cloud service lỗi không?
2. Có cần cấu hình giới hạn dung lượng tối đa cho mỗi tệp tải lên ở màn hình quản trị hệ thống không?

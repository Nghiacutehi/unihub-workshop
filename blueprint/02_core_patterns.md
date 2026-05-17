# 2. Core Technical Patterns

UniHub Workshop áp dụng hàng loạt các Technical Pattern cao cấp để xử lý bài toán **High Concurrency** (tải cao) và **Data Consistency** (Tính nhất quán dữ liệu).

## 2.1 Virtual Waiting Room (Phòng chờ ảo)
**Vấn đề:** 12,000 sinh viên cùng đăng ký, nếu cho phép hit DB liên tục sẽ sập DB (Connection Pool Exhaustion).
**Giải pháp:** 
- Khởi tạo hàng đợi ZSET trên Redis.
- Khi tải tăng cao, request bị chặn ở Middleware/Service và đẩy vào Queue. User nhận mã `202 Queued` và tự động Polling (hỏi lại server liên tục).
- Lệnh Lua Script chạy định kỳ sẽ "promote" những người ở lâu nhất sang trạng thái "Active" để họ được vào đăng ký.

## 2.2 Event-Driven Architecture (Thỏ RabbitMQ)
**Vấn đề:** Các xử lý liên quan khóa DB và gửi Email mất thời gian.
**Giải pháp:**
- Áp dụng mô hình Publish - Subscribe.
- API chỉ làm nhiệm vụ kiểm tra tính hợp lệ cơ bản, sau đó đẩy JSON Payload (Event) vào RabbitMQ và phản hồi `202 Accepted` ngay trong 10ms.
- `Registration Worker` chạy ngầm, từ từ consume message, trừ vé, và gửi thông báo.

## 2.3 Pessimistic Locking (Khóa Bi Quan)
**Vấn đề:** 10,000 sinh viên tranh nhau 1,000 chỗ ngồi. Dễ dính Race Condition (Bán lố - Oversell).
**Giải pháp:**
- Sử dụng lệnh SQL: `SELECT available_seats FROM workshops WHERE id = $1 FOR UPDATE`.
- Lệnh `FOR UPDATE` khóa (Lock) row đó ở mức độ Database.
- Khi Worker 1 đang check chỗ và trừ vé, Worker 2 (hoặc thread khác) phải "Wait" cho đến khi Worker 1 `COMMIT`. Đảm bảo chỗ ngồi không bao giờ âm.

## 2.4 Idempotency (Tính Lũy Đẳng)
**Vấn đề:** Người dùng click nút "Đăng ký" 2 lần liên tiếp hoặc hệ thống mạng bị chập chờn tự động Retry. Dẫn tới trừ chỗ 2 lần.
**Giải pháp:**
- **Lớp 1 (Redis):** Tạo một key `reg:lock:{workshop_id}:{student_id}` bằng thao tác nguyên tử `SETNX` (Set if Not Exists) với TTL 5 phút. Lần gọi thứ 2 sẽ bị từ chối ngay ở Cache.
- **Lớp 2 (Postgres):** Đặt Constraint Unique (`uq_user_workshop`) trên bảng registrations.

## 2.5 Circuit Breaker (Cầu Dao Tự Động)
**Vấn đề:** Hệ thống phụ thuộc API AI của Google (Gemini). Nếu Google sập hoặc khóa Quota, luồng tạo Workshop bị treo cứng chờ timeout.
**Giải pháp:**
- Wrap (bọc) lệnh gọi Gemini bằng Circuit Breaker.
- Track số lỗi. Nếu lỗi 50% trong thời gian ngắn -> Cầu dao đóng (OPEN). Mọi lệnh gọi tới AI sau đó sẽ bị từ chối ngay lập tức (Fail Fast) trong 30 giây để bảo vệ resource, sau 30 giây mới mở lại (HALF-OPEN) để test thử.

## 2.6 Pipe-and-Filter AI Processing
**Vấn đề:** File PDF upload lên chứa nhiều rác, ký tự lạ, cần xử lý nhiều bước.
**Giải pháp:**
- Dữ liệu chảy qua một đường ống (Pipe).
- Bộ lọc 1 (Filter): `extractText` (Lọc Text thuần từ PDF).
- Bộ lọc 2 (Filter): `cleanText` (Xóa khoảng trắng thừa, xóa control characters).
- Bộ lọc 3 (Filter): Gọi Gemini AI.
- Dễ dàng thay thế hoặc cắm thêm một bộ lọc mới mà không ảnh hưởng code cũ.

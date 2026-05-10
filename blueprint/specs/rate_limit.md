# Đặc tả: Kiểm soát tải đột biến (rate_limit.md)

## 1. Mô tả
Tài liệu này quy định cơ chế bảo vệ hệ thống trước tình trạng quá tải lưu lượng (Traffic Saturation). Hệ thống phải đảm bảo khả năng tiếp nhận và xử lý ổn định cho 12.000 request trong thời gian ngắn, ngăn chặn việc cạn kiệt tài nguyên của Backend API và Database.

## 2. Giải thuật xử lý
Hệ thống triển khai thuật toán Token Bucket tại tầng API Gateway:
- Mỗi định danh (MSSV hoặc IP) được cấp một xô (Bucket) chứa Token.
- Request chỉ được phép đi tiếp nếu xô còn Token.
- Token được nạp lại vào xô theo một tốc độ cố định (Refill Rate).

## 3. Các bước xử lý
1. API Gateway tiếp nhận Request và trích xuất định danh (JWT sub hoặc Client IP).
2. Gateway gửi truy vấn tới Redis để kiểm tra số lượng Token còn lại trong xô của định danh đó.
3. Nếu Token > 0:
   - Trừ 1 Token trong Redis.
   - Chuyển tiếp Request vào Backend Service.
4. Nếu Token = 0:
   - Từ chối Request ngay lập tức.
   - Phản hồi mã lỗi HTTP 429.

## 4. Kịch bản lỗi
- Redis Connection Timeout: Trong trường hợp không thể kết nối tới Redis, Gateway sẽ chuyển sang chế độ Fail-open (cho phép request qua) hoặc áp dụng Local Rate Limit (giới hạn tại RAM của Gateway) để bảo vệ hệ thống ở mức tối thiểu.
- Vượt ngưỡng (Rate Limit Exceeded): Trả về mã lỗi HTTP 429 kèm Header Retry-After để thông báo cho Client thời gian cần chờ.

## 5. Ràng buộc
- Thời gian xử lý kiểm tra tại Gateway không được vượt quá 50ms.
- Dữ liệu bộ đếm trong Redis phải có thời gian hết hạn (TTL) để giải phóng bộ nhớ.

## 6. Tiêu chí chấp nhận
- Hệ thống duy trì trạng thái hoạt động (Availability) khi có 12.000 sinh viên truy cập đồng thời.
- Không có hiện tượng sập dịch vụ (Cascading Failure) do quá tải request.



# Đặc tả: rate_limit.md

## 1. Mô tả
Tài liệu quy định cơ chế bảo vệ hệ thống trước `Spike Traffic` với lưu lượng dự kiến 12.000 `Request` trong 10 phút. Cơ chế kiểm soát tải được thực thi trực tiếp tại `API Gateway` nhằm đảm bảo `Availability` cho `Business Logic` và ngăn chặn tình trạng `Resource Exhaustion` tại `Database`.

## 2. Thuật toán Token Bucket
Hệ thống sử dụng thuật toán `Token Bucket` với cấu hình tham số như sau:

* `Capacity`: 50 `Tokens`. Cho phép hệ thống tiếp nhận `Burst Traffic` ban đầu từ mỗi `Client`.
* `Refill Rate`: 5 `Tokens/second`. Giới hạn tốc độ xử lý trung bình.
* `Identifier`: Xác định thông qua `user_id` trích xuất từ `JWT Payload`
* `Atomicity`: Sử dụng `Lua Scripting` thực thi trực tiếp trên `Redis Engine`. Đảm bảo các bước kiểm tra (Read) và trừ (Decrement) `Token` diễn ra trong một giao dịch duy nhất (`Atomic Operation`), ngăn chặn triệt để lỗi `Race Condition` khi `Client` gửi yêu cầu đồng thời.

## 3. Luồng xử lý (Workflow)
Quy trình đánh giá và định tuyến `Request` tại tầng `API Gateway`:

1. `API Gateway` tiếp nhận `HTTP Request` và trích xuất `Identifier`.
2. Gọi `Redis` thực thi `Lua Script` để kiểm tra trạng thái `Token Bucket`.
3. Nhánh Hợp lệ (`Token > 0`): `Redis` tự động trừ 1 `Token`, `API Gateway` đính kèm `Header` trạng thái và `Forward Request` vào `Backend Services`.
4. Nhánh Từ chối (`Token = 0`): `API Gateway` hủy `Request`, phản hồi mã `429 Too Many Requests` ngay lập tức.


## 4. Quản lý bộ nhớ (Redis Memory Management)
Hệ thống áp dụng chính sách giải phóng tài nguyên tự động:

* `TTL`: 300s (5 phút).
* Ràng buộc: `TTL` bắt buộc phải lớn hơn thời gian làm đầy `Bucket` (`Capacity` / `Refill Rate`) để tránh lỗ hổng `Bypass`.Mọi `Key` khởi tạo trong `Redis` bắt buộc phải đính kèm `TTL`, Khi thời gian `Idle` vượt 300 giây, `Redis` tự động dọn dẹp `Key` để giải phóng `RAM`.

## 5. Fallback Strategy (In-memory Rate Limiting)
Phương án dự phòng bảo vệ `Database` khi cụm `Redis` gặp sự cố:

* Điều kiện kích hoạt: `Redis Connection Timeout` vượt quá ngưỡng ngắt mạch `50ms`.
* Luồng thực thi: `API Gateway` kích hoạt `Circuit Breaker`, tạm ngắt kết nối với `Redis`.
* Cơ chế RAM: Khởi tạo `In-memory Cache` (Sử dụng `Caffeine Cache` hoặc `ConcurrentHashMap` hỗ trợ `Atomic Operations`) trực tiếp trên `Heap Memory` của `Gateway`.
* Ràng buộc Fallback: Cấu hình `Size-based Eviction` (Tối đa 10.000 `Keys`) để ngăn chặn tình trạng `Out of Memory` (OOM) cho chính `API Gateway`.
## 6. Client-side Handshake
Yêu cầu bắt buộc đối với `Frontend` nhằm triệt tiêu `Spam Request` tại tầng `Presentation`:

* Nhận diện: Bắt `HTTP Status Code 429`.
* Trích xuất: Đọc thông số thời gian từ `Header` `Retry-After`.
* Thực thi: Cập nhật `State` của `UI`, thiết lập thuộc tính `Disabled` cho `Button` đăng ký trong khoảng thời gian tương ứng.

## 7. Xử lý ngoại lệ

| Trạng thái | HTTP Status Code | Phản hồi hệ thống |
| :--- | :--- | :--- |
| `Rate Limit Exceeded` | `429 Too Many Requests` | Từ chối `Request`, đính kèm `Header` `Retry-After`. |
| `Redis Timeout / Failure` | N/A | Kích hoạt `Local Rate Limit`, ghi nhận `Warning Log`. |
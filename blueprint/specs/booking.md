# Đặc tả: Tranh chấp chỗ ngồi và Đăng ký (booking.md)

## 1. Mô tả
Tài liệu quy định luồng xử lý đăng ký tham dự Workshop. Trọng tâm của module này là giải quyết bài toán `Concurrency` (Tranh chấp đồng thời) khi số lượng `Request` vượt quá số lượng chỗ ngồi có sẵn (ví dụ: hàng trăm sinh viên tranh 60 chỗ). Hệ thống kết hợp kiến trúc `Event-Driven` để điều hòa tải và `Pessimistic Locking` để đảm bảo tính toàn vẹn dữ liệu.

## 2. Kiến trúc giải pháp
Hệ thống không cho phép `API` ghi trực tiếp vào `Database` khi có `Spike Traffic`, thay vào đó sử dụng mô hình xử lý bất đồng bộ:

* `Message Broker`: Sử dụng `RabbitMQ` hoặc `Kafka` làm vùng đệm hứng chịu toàn bộ lưu lượng `Request` hợp lệ từ `API Gateway`. Đóng vai trò `Load Leveling` (San phẳng tải).
* `Background Workers`: Các tiến trình chạy ngầm sẽ kéo (`Consume`) từng `Message` từ `Queue` ra để xử lý tuần tự theo tốc độ mà `Database` có thể chịu đựng.
* `Concurrency Control`: Sử dụng cơ chế khóa bi quan (`Pessimistic Locking`) tại tầng `PostgreSQL` để đảm bảo không xảy ra hiện tượng `Overbooking`.

## 3. Cơ chế Khóa dữ liệu (Database Locking)
Hệ thống từ chối sử dụng `Optimistic Locking` do tỷ lệ xung đột (`Contention Rate`) trong 3 phút đầu là cực kỳ cao, dễ dẫn đến hiện tượng `Rollback` hàng loạt làm cạn kiệt `CPU`. Thay vào đó, hệ thống áp dụng `Pessimistic Locking`:

* Cú pháp thực thi: Sử dụng lệnh `SELECT ... FOR UPDATE` trong một `Database Transaction`.
* Hành vi: Khi `Worker A` đang đọc số lượng chỗ ngồi của Workshop X, dòng dữ liệu (`Row`) của Workshop X sẽ bị khóa. `Worker B` muốn đọc dòng này bắt buộc phải đưa vào trạng thái chờ (`Block/Wait`) cho đến khi `Worker A` hoàn tất `Transaction` (đăng ký xong và trừ đi 1 chỗ).
* Trade-off: Chấp nhận giảm `Throughput` của `Database` để đổi lấy `Strong Consistency` (Sự nhất quán tuyệt đối).

## 4. Luồng xử lý (Workflow)
Quy trình thực thi luồng đăng ký từ `Client` đến `Database`:

1. `Client` gửi `Request` đăng ký chứa `workshop_id`.
2. `Backend API` tiếp nhận, khởi tạo một `Correlation ID` (Mã theo dõi) và đẩy `Message` vào `Registration Queue`.
3. `API` phản hồi ngay lập tức cho `Client` trạng thái `HTTP 202 Accepted` (Đang xử lý) kèm `Correlation ID`.
4. `Client` thiết lập cơ chế `Polling` (hoặc `Server-Sent Events`) để theo dõi trạng thái của `Correlation ID`.
5. `Background Worker` lấy `Message` từ `Queue` và mở một `Transaction` trên `PostgreSQL`.
6. Thực thi `SELECT available_seats FROM workshops WHERE id = ? FOR UPDATE`.
7. Nhánh Đủ chỗ (`available_seats > 0`):
   - Giảm `available_seats` đi 1.
   - Insert bản ghi vào bảng `registrations` với trạng thái `SUCCESS` (hoặc `PENDING_PAYMENT` nếu có phí).
   - `Commit Transaction` (Giải phóng khóa).
8. Nhánh Hết chỗ (`available_seats = 0`):
   - `Rollback Transaction` (Giải phóng khóa).
   - Đánh dấu trạng thái của `Correlation ID` thành `FAILED_FULL`.

## 5. Kịch bản lỗi
* `Queue Overflow`: `Message Broker` đạt giới hạn hàng đợi. Hệ thống kích hoạt `Circuit Breaker`, từ chối nhận thêm `Request` và trả về `HTTP 503 Service Unavailable`.
* `Database Deadlock`: Nếu xảy ra khóa chéo, `PostgreSQL` sẽ tự động hủy một `Transaction`. `Worker` bắt được lỗi `Deadlock Exception` sẽ tự động đưa `Message` trở lại `Queue` để `Retry` tối đa 3 lần.
* `Worker Crash`: Nếu `Worker` chết giữa chừng khi đang giữ `Message` (chưa gửi lệnh `ACK` cho `Broker`), `Message Broker` sẽ đưa `Message` đó sang một `Worker` khác xử lý để không mất yêu cầu của sinh viên.

## 6. Xử lý ngoại lệ

| Lớp thành phần | Mã lỗi / Trạng thái | Hành động xử lý |
| :--- | :--- | :--- |
| `API Layer` | `202 Accepted` | Tiếp nhận thành công, sinh viên chờ kết quả. |
| `Worker Layer` | `available_seats = 0` | Cập nhật `Status = REJECTED`, gửi thông báo hết chỗ. |
| `Broker Layer` | `Connection Refused` | Bật chế độ `Graceful Degradation`, báo lỗi bảo trì. |
| `Database Layer` | `Lock Timeout` | `Rollback`, đẩy `Message` vào `Dead Letter Queue` (DLQ). |
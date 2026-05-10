# Đặc tả: checkin.md

## 1. Mô tả
Tài liệu quy định luồng nghiệp vụ `Attendance Tracking` tại hiện trường sự kiện. Chức năng được thiết kế tuân thủ hệ tư tưởng `AP System` trong `CAP Theorem`, ưu tiên `Availability` và `Partition Tolerance` để nhân sự duy trì hoạt động xác thực khi `Network Connection` bị gián đoạn.

## 2. Kiến trúc giải pháp
Hệ thống áp dụng kiến trúc `Offline-first` kết hợp `Eventual Consistency`:

* `Local Persistence`: Tích hợp `SQLite` trực tiếp tại `Client App` làm `Cache Layer` bền vững.
* `Data Synchronization`: Thiết lập `Background Worker` tại `Client` thực thi đồng bộ dữ liệu.
* `Idempotency`: Áp dụng `UUIDv4` làm `Primary Key` tại `Local DB` cho mỗi lượt `Scan` nhằm loại trừ rủi ro `Data Duplication` khi thực thi `Retry`.
* `Security`: Áp dụng chuẩn mã hóa `AES-256` cho toàn bộ file `Database` của `SQLite`.

## 3. Workflow

### 3.1. Giai đoạn Local Check-in
1. `Client` khởi chạy `Camera Module` thực thi `QR Decode`.
2. `Validation`: Thực thi `Signature Verification` bằng `Public Key` lưu tại `Client`.
3. Kiểm tra `Network Status`:
   - Phân nhánh `Connected`: Gửi `HTTP POST /api/v1/checkin/live`.
   - Phân nhánh `Disconnected`: Thực thi `INSERT` vào `SQLite`.
4. Khởi tạo `Record` trong `SQLite` với `Schema`:
   - `id`: `UUIDv4`
   - `student_id`: `String`
   - `workshop_id`: `String`
   - `scanned_at`: `UNIX Timestamp`
   - `sync_status`: `PENDING`

### 3.2. Giai đoạn Background Synchronization
1. `Network Listener` tại `Client` phát hiện `Connection` chuyển trạng thái `UP`.
2. Truy vấn `SQLite` lấy tập hợp `Records` có `sync_status = 'PENDING'`.
3. Đóng gói `Records` thành `JSON Array`.
4. Gọi `HTTP POST /api/v1/checkin/sync` thực thi phương thức `Bulk Insert`.
5. `API Gateway` định tuyến `Request` tới `Checkin Service`.
6. `Checkin Service` thực thi `UPDATE` vào `PostgreSQL`.
7. `Backend` phản hồi `HTTP 200 OK` kèm danh sách `id` đã xử lý.
8. `Client` tiếp nhận `Response` và thực thi `DELETE` các `Records` tương ứng trong `SQLite`.

## 4. Xử lý ngoại lệ

| Lỗi (Exception) | Nguyên nhân | Hành vi hệ thống |
| :--- | :--- | :--- |
| `Data Conflict` | Quét cùng một `QR Code` trên nhiều `Devices` khác nhau trong trạng thái `Offline`. | `Backend` đối chiếu `scanned_at`, giữ bản ghi có `Timestamp` nhỏ nhất, loại bỏ bản ghi còn lại và xuất `Warning Log`. |
| `Invalid Signature` | `Payload` của `QR Code` bị chỉnh sửa hoặc không thuộc `System`. | `Client Validation` fail, từ chối lưu vào `SQLite`, ném ngoại lệ `InvalidTokenException`. |
| `Network Timeout` | Mất kết nối khi đang truyền `Bulk Request`. | `Client` giữ nguyên `sync_status = 'PENDING'`, tái thực thi `Sync` sau khoảng thời gian `Retry Interval`. |

## 5. Ràng buộc kỹ thuật
* `Sync Interval`: Tiến trình đồng bộ chạy định kỳ 30s/lần thông qua `Cron Job` tại `Client`.
* `Payload Size Limit`: Thuật toán `Pagination` tại `Client` giới hạn tối đa 500 `Records` cho mỗi `Bulk Request` để ngăn chặn `Payload Too Large`.
* `Performance`: Tốc độ xử lý `Local Check-in` tại `Client` (Bao gồm `Validation` và `SQLite INSERT`) không vượt quá 100ms.
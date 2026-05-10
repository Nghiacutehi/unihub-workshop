# Đặc tả: payment.md

## 1. Mô tả
Tài liệu quy định luồng xử lý đăng ký tham dự đối với các `Workshop` có thu phí. Trọng tâm của module là tích hợp an toàn với `Payment Gateway` (Cổng thanh toán bên thứ ba), đảm bảo nguyên tắc `ACID` cho giao dịch tài chính, ngăn chặn lỗi `Double Charge` (Trừ tiền hai lần) và duy trì `Availability` cho toàn hệ thống khi đối tác gặp sự cố.

## 2. Kiến trúc giải pháp
Hệ thống kết hợp nhiều `Design Patterns` và cơ chế chịu lỗi để quản lý vòng đời giao dịch:

* `Strategy Pattern`: Trừu tượng hóa các phương thức thanh toán (ví dụ: `MoMoStrategy`, `VNPayStrategy`). Hỗ trợ mở rộng cổng thanh toán mới tuân thủ nguyên lý `OCP` (Open/Closed Principle).
* `Circuit Breaker`: Áp dụng mẫu thiết kế ngắt mạch để bảo vệ `Thread Pool` của hệ thống khi `Payment Gateway` phản hồi chậm hoặc `Timeout`.
* `Idempotency`: Quản lý tính lũy đẳng của giao dịch thông qua `Redis`, đảm bảo một `Webhook` hoặc `Retry Request` chỉ được xử lý cập nhật trạng thái đúng một lần.

## 3. Workflow
Quy trình thực thi từ lúc khởi tạo đến khi xác nhận thanh toán:

1. `Client` gửi `POST Request` yêu cầu đăng ký `Workshop` có phí.
2. `Backend` thực thi `Pessimistic Locking` kiểm tra `available_seats`. Nếu hợp lệ, tạm giữ chỗ và tạo bản ghi `Registration` với trạng thái `PENDING_PAYMENT`.
3. Khởi tạo `Transaction_ID` duy nhất và gọi `API` của `Payment Gateway` để tạo `Checkout URL`.
4. Trả `Checkout URL` về `Client`. Hệ thống bắt đầu đếm ngược `Payment TTL` (thời gian giữ chỗ, ví dụ: 15 phút).
5. `Client` thực hiện thanh toán trên giao diện của `Payment Gateway`.
6. `Payment Gateway` gửi `Webhook Callback` báo trạng thái giao dịch về `API Gateway`.
7. `Backend` thực thi `Signature Verification` để xác thực nguồn gốc `Webhook`.
8. Kiểm tra `Idempotency Key` trong `Redis`. Nếu đây là `Request` lần đầu, thực thi nghiệp vụ cập nhật:
   - Chuyển trạng thái `Registration` thành `SUCCESS`.
   - Sinh `QR Code`.
   - Đẩy sự kiện vào `Notification Queue` để gửi `Email`.
9. `Backend` trả về `HTTP 200 OK` cho `Payment Gateway`.

## 4. Cơ chế chống trùng lặp (Idempotency)
Để giải quyết lỗi gọi `Webhook` nhiều lần từ đối tác hoặc do mạng chập chờn:

* Khởi tạo `Key`: Cấu trúc `payment:idempotency:{transaction_id}`.
* Thực thi kiểm tra: Khi nhận `Webhook`, hệ thống sử dụng lệnh `SETNX` (Set if Not eXists) của `Redis` với `TTL = 24h`.
* Phân nhánh:
   - Nếu `SETNX` trả về `true`: Giao dịch lần đầu, hệ thống tiến hành xử lý `Business Logic`.
   - Nếu `SETNX` trả về `false`: Khóa đã tồn tại, hệ thống xác định đây là `Duplicate Request`. Lập tức bỏ qua `Business Logic` và trả về `HTTP 200 OK` để đối tác ngừng gửi lại `Webhook`.

## 5. Cơ chế chịu lỗi (Circuit Breaker)
Cấu hình ngắt mạch bảo vệ hệ thống trước sự cố từ `Payment Gateway`:

* `State CLOSED`: Trạng thái bình thường, `Request` được phép gọi sang `Payment Gateway`.
* `State OPEN`: Kích hoạt khi tỷ lệ `Error Rate` (Lỗi 5xx hoặc Timeout) vượt quá 50% trong khung thời gian 10 giây. Toàn bộ `Request` gọi sang `Payment Gateway` bị hệ thống chặn lập tức (`Fail-fast`).
* `Fallback`: Khi mạch `OPEN`, hệ thống trả về thông báo lỗi, `Frontend` tự động thực thi `Graceful Degradation` (Ví dụ: Tạm ẩn nút "Thanh toán MoMo", chỉ cho phép đăng ký Workshop miễn phí).
* `State HALF-OPEN`: Sau khoảng thời gian `Sleep Window` (ví dụ: 30 giây), mạch cho phép một lượng nhỏ `Request` đi qua để dò đường. Nếu thành công, mạch chuyển về `CLOSED`.

## 6. Xử lý ngoại lệ

| Trường hợp | Giai đoạn | Xử lý hệ thống |
| :--- | :--- | :--- |
| `Payment Timeout` | Trạng thái chờ `Webhook` quá 15 phút | Background Worker quét định kỳ, tự động chuyển trạng thái giao dịch thành `CANCELLED`, hoàn trả chỗ ngồi (`available_seats + 1`). |
| `Invalid Signature` | Kiểm tra `Webhook` | Hủy `Request`, ném ngoại lệ `SecurityException`, ghi nhận log tấn công. |
| `Payment Failed` | `Webhook` trả về mã lỗi tài khoản | Cập nhật trạng thái `FAILED`, giữ nguyên `available_seats`, thông báo `Client`. |

## 7. Ràng buộc kỹ thuật
* `Data Type`: Các biến lưu trữ số tiền (`Amount`) phải sử dụng kiểu dữ liệu `Decimal` hoặc `Numeric`, tuyệt đối không sử dụng `Float` hoặc `Double` để tránh sai số dấu phẩy động.
* `Timeout Config`: Lệnh gọi `API` ra bên ngoài (`Outbound Call`) bắt buộc phải cấu hình `Connection Timeout` (tối đa 3s) và `Read Timeout` (tối đa 5s).
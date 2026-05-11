# Đặc tả: notification.md

## 1. Mô tả
Tài liệu quy định cơ chế vận hành của hệ thống thông báo đa kênh (`Multi-channel Notification System`). Hệ thống chịu trách nhiệm gửi thông tin xác nhận đăng ký, mã `QR Code` và các cập nhật thay đổi trạng thái Workshop tới người dùng thông qua nhiều phương thức truyền thông khác nhau.

## 2. Kiến trúc giải pháp
Hệ thống áp dụng kết hợp các mẫu thiết kế để đảm bảo tính mở rộng và tách biệt trách nhiệm:

* `Observer Pattern`: `Notification Service` đóng vai trò là một `Observer` lắng nghe các sự kiện (`Events`) từ `Registration Service`.
* `Strategy Pattern`: Mỗi kênh thông báo (Email, Web Notification, Telegram) được triển khai như một `Strategy` độc lập kế thừa từ `NotificationStrategy` `Interface`.
* `Open/Closed Principle (OCP)`: Việc bổ sung kênh thông báo mới chỉ yêu cầu tạo thêm một `Class` `Strategy` mới mà không cần chỉnh sửa logic cốt lõi.
* `Asynchronous Processing`: Sử dụng `Message Queue` để xử lý việc gửi thông báo nhằm tránh gây nghẽn cho luồng đăng ký chính.

## 3. Workflow
Quy trình thực thi luồng thông báo:

1. `Registration Service` hoàn tất cập nhật CSDL, phát một sự kiện (kèm `event_id` duy nhất) vào `Notification Queue`.
2. `Notification Dispatcher` (thuộc `Background Worker`) tiêu thụ `Message` từ hàng đợi.
3. `Dispatcher` khởi tạo các bản ghi vào bảng `notifications` trên `PostgreSQL` với trạng thái `PENDING`.
4. Đối với mỗi kênh, `Dispatcher` gọi `Execute()` trên `Strategy` tương ứng:
   - `EmailStrategy`: Khởi tạo nội dung `HTML`, gọi `SMTP Server`.
   - - `WebNotificationStrategy`: Lưu thông báo vào bảng `notifications` trong CSDL để hiển thị tại biểu tượng "Chuông thông báo" trên giao diện Web khi sinh viên đăng nhập.
5. Hệ thống thực thi `UPDATE` bảng `notifications` thành `SENT` (kèm `sent_at`) hoặc `FAILED` (kèm `error_message`) dựa trên phản hồi từ bên thứ ba.

## 4. Quản lý dữ liệu lịch sử (Persistence)
Hệ thống lưu trữ mọi luồng thông báo tại bảng `notifications` nhằm phục vụ tính năng "Lịch sử thông báo" (In-app Inbox) cho sinh viên và hỗ trợ truy vết lỗi. Cấu trúc lưu vết bao gồm:

* `user_id`: Định danh người nhận.
* `channel`: Phương thức gửi (`EMAIL`, `WEB`).
* `status`: Trạng thái quá trình gửi (`PENDING`, `SENT`, `FAILED`).
* `event_id`: Mã sự kiện gốc (Ví dụ: `REG_SUCCESS_33333333`).

## 5. Xử lý ngoại lệ

| Lỗi (Exception) | Nguyên nhân | Hành vi hệ thống |
| :--- | :--- | :--- |
| `Provider Down` | Dịch vụ `Email` hoặc `Message Broker` gặp sự cố. | Kích hoạt `Retry Policy` với `Exponential Backoff`. Cập nhật trạng thái `FAILED`. Nếu quá số lần thử lại, đẩy bản ghi vào `Dead Letter Queue` (`DLQ`). |
| `Invalid Recipient` | Địa chỉ Email sai định dạng hoặc `User_ID`  không tồn tại. | Ngừng xử lý cho kênh đó, ghi log `Error Message` vào CSDL và tiếp tục gửi qua các kênh khác. |
| `Rate Limit` | Vượt quá giới hạn gửi của `Provider`. | Tạm dừng gửi (`Pause`) và lập lịch lại sau một khoảng thời gian chờ. |

## 6. Ràng buộc kỹ thuật
* `Concurrency`: Cho phép xử lý song song các `Strategy` khác nhau cho cùng một sự kiện để tối ưu tốc độ.
* `Idempotency`: Thiết lập `UNIQUE (event_id, channel)` tại tầng `Database` để chặn tuyệt đối việc gửi trùng lặp thông báo cho cùng một sự kiện qua cùng một kênh khi luồng xử lý bị `Retry`.
* `Template Engine`: Nội dung thông báo phải được quản lý qua các `Templates` riêng biệt, hỗ trợ đa ngôn ngữ.

## 7. Tiêu chí chấp nhận
* Sinh viên nhận được thông báo trong vòng dưới 30 giây kể từ khi đăng ký thành công.
* Hệ thống vẫn vận hành bình thường nếu một trong các kênh thông báo (ví dụ Email) gặp sự cố.
* Có khả năng bật/tắt các kênh thông báo thông qua tệp cấu hình hệ thống mà không cần biên dịch lại mã nguồn.
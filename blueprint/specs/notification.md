
# Đặc tả: notification.md

## 1. Mô tả
Tài liệu quy định cơ chế vận hành của hệ thống thông báo đa kênh (`Multi-channel Notification System`). [cite_start]Hệ thống chịu trách nhiệm gửi thông tin xác nhận đăng ký, mã `QR Code` và các cập nhật thay đổi trạng thái Workshop tới người dùng thông qua nhiều phương thức truyền thông khác nhau.

## 2. Kiến trúc giải pháp
Hệ thống áp dụng kết hợp các mẫu thiết kế để đảm bảo tính mở rộng và tách biệt trách nhiệm:

* `Observer Pattern`: `Notification Service` đóng vai trò là một `Observer` lắng nghe các sự kiện (`Events`) từ `Registration Service`. Khi một đơn đăng ký được xác nhận, một sự kiện `RegistrationSuccessful` sẽ được phát đi.
* `Strategy Pattern`: Mỗi kênh thông báo (Email, In-app Push, Telegram) được triển khai như một `Strategy` độc lập kế thừa từ `NotificationStrategy` `Interface`.
* [cite_start]`Open/Closed Principle (OCP)`: Việc bổ sung kênh thông báo mới chỉ yêu cầu tạo thêm một `Class` `Strategy` mới mà không cần chỉnh sửa logic cốt lõi của `Notification Dispatcher`.
* `Asynchronous Processing`: Sử dụng `Message Queue` để xử lý việc gửi thông báo nhằm tránh gây nghẽn cho luồng đăng ký chính.

## 3. Workflow
Quy trình thực thi luồng thông báo:

1. `Registration Service` hoàn tất cập nhật `Database`, phát một sự kiện vào `Notification Queue`.
2. `Notification Dispatcher` (thuộc `Background Worker`) tiêu thụ `Message` từ hàng đợi.
3. `Dispatcher` xác định danh sách các kênh cần gửi dựa trên cấu hình hệ thống hoặc tùy chọn của người dùng.
4. Đối với mỗi kênh, `Dispatcher` gọi `Execute()` trên `Strategy` tương ứng:
   - `EmailStrategy`: Khởi tạo nội dung `HTML`, gọi `SMTP Server` hoặc `Mail Service`.
   - `InAppPushStrategy`: Gọi `Firebase Cloud Messaging` (`FCM`) để đẩy thông báo tới `Mobile App`.
5. Hệ thống ghi nhận trạng thái gửi (`Sent`, `Failed`) vào nhật ký hệ thống.

## 4. Cấu trúc dữ liệu Strategy

| Thành phần | Kiểu dữ liệu | Mô tả |
| :--- | :--- | :--- |
| `NotificationID` | `UUID` | Định danh duy nhất cho mỗi yêu cầu thông báo. |
| `Recipient` | `String` | Địa chỉ nhận (Email, Device Token, hoặc Chat ID). |
| `Payload` | `JSON` | Chứa nội dung thông báo, tiêu đề và các thông số tùy biến. |
| `ProviderConfig` | `Object` | Chứa thông tin kết nối (API Key, Endpoint) của dịch vụ bên thứ ba. |

## 5. Xử lý ngoại lệ

| Lỗi (Exception) | Nguyên nhân | Hành vi hệ thống |
| :--- | :--- | :--- |
| `Provider Down` | Dịch vụ `Email` hoặc `Push` gặp sự cố. | Kích hoạt `Retry Policy` với `Exponential Backoff`. Nếu quá số lần thử lại, đẩy bản ghi vào `Dead Letter Queue` (`DLQ`). |
| `Invalid Recipient` | Địa chỉ Email sai định dạng hoặc `Token` thiết bị hết hạn. | Ngừng xử lý cho kênh đó, ghi log `Error` và tiếp tục gửi qua các kênh khác. |
| `Rate Limit` | Vượt quá giới hạn gửi của `Provider`. | Tạm dừng gửi (`Pause`) và lập lịch lại sau một khoảng thời gian chờ. |

## 6. Ràng buộc kỹ thuật
* `Concurrency`: Cho phép xử lý song song các `Strategy` khác nhau cho cùng một sự kiện để tối ưu tốc độ.
* `Idempotency`: Kiểm tra `Event_ID` trước khi gửi để tránh gửi trùng lặp thông báo cho cùng một sự kiện đăng ký.
* `Template Engine`: Nội dung thông báo phải được quản lý qua các `Templates` riêng biệt, hỗ trợ đa ngôn ngữ.

## 7. Tiêu chí chấp nhận
* Sinh viên nhận được thông báo trong vòng dưới 30 giây kể từ khi đăng ký thành công.
* Hệ thống vẫn vận hành bình thường nếu một trong các kênh thông báo (ví dụ Email) gặp sự cố.
* Có khả năng bật/tắt các kênh thông báo thông qua tệp cấu hình hệ thống mà không cần biên dịch lại mã nguồn.
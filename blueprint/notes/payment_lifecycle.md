# Vòng đời Giao dịch & Cơ chế Tự phục hồi (Payment Lifecycle)

Tài liệu này mô tả cách hệ thống UniHub quản lý các đơn đăng ký chờ thanh toán, đảm bảo tính toàn vẹn dữ liệu và hiệu suất cao.

## 1. Cơ chế TTL (Time-To-Live) cho Đơn hàng
Mọi đơn đăng ký Workshop có phí đều có thời hạn thanh toán nhất định để tránh việc giữ chỗ ảo quá lâu.

*   **Thời gian chờ:** 15 phút (Mặc định).
*   **Trạng thái kích hoạt:** `PENDING_PAYMENT`.
*   **Hành động khi hết hạn:**
    1.  Xóa bản ghi `payments` liên quan.
    2.  Xóa bản ghi `registrations`.
    3.  **Hoàn trả ghế (Release Seat)**: Cộng lại số ghế trống trong DB.
    4.  **Đồng bộ Redis**: Cập nhật lại Seat Counter trên Redis Cache để người khác có thể đăng ký ngay lập tức.

## 2. Worker Dọn dẹp (Cleanup Worker)
Hệ thống chạy một Background Worker định kỳ (mỗi 5 phút) để quét và xử lý các đơn hàng hết hạn.
*   **File xử lý:** `payment_service.go` -> `CleanupExpiredPayments`.
*   **Đảm bảo:** Sử dụng Transaction để đảm bảo việc xóa dữ liệu và hoàn trả ghế diễn ra nguyên tử (Atomic).

## 3. Tính Lũy đẳng (Idempotency)
Để chống lại việc đăng ký trùng lặp hoặc xử lý Webhook nhiều lần:

### Tại Tầng Đăng ký:
*   **Database Constraint:** Ràng buộc `uq_user_workshop` ngăn chặn một người dùng có 2 bản ghi cho cùng 1 Workshop.
*   **Worker Check:** Nếu phát hiện lỗi trùng lặp, Worker sẽ đánh dấu `FAILED` và ACK tin nhắn thay vì thử lại (retry), giúp giải phóng hàng đợi.

### Tại Tầng Webhook (Thanh toán):
*   **Status Lock**: Trước khi xử lý Webhook, hệ thống kiểm tra trạng thái hiện tại của Payment. Nếu đã là `SUCCESS` hoặc `FAILED`, Webhook sẽ bị bỏ qua (Ignore) để tránh xử lý thừa.
*   **Xác thực Chữ ký**: Mọi Webhook từ Sandbox/Ngân hàng đều phải đi kèm chữ ký bảo mật để ngăn chặn các yêu cầu giả mạo trạng thái thanh toán.

## 4. Đồng bộ hóa Email Thông báo
Khi trạng thái thanh toán chuyển sang `SUCCESS`, hệ thống thực hiện đồng thời:
1.  Ký số RSA cho vé (Ticket Signature).
2.  Gửi Email xác nhận đăng ký với giao diện đồng nhất (Premium Template), kèm dòng xác nhận thanh toán thành công.
3.  Cập nhật thông báo In-app (Web Notification).

# Độ tin cậy & Hiệu năng Hệ thống (Reliability & Performance)

Tài liệu này ghi nhận các cơ chế tự phục hồi và tối ưu hóa hạ tầng của UniHub.

## 1. Cơ chế Tự phục hồi RabbitMQ (Self-Healing)

Backend tích hợp cơ chế tự động kết nối lại khi Publisher bị mất liên kết với RabbitMQ Server.

*   **Logic:**
    -   Sử dụng `NotifyClose` để lắng nghe sự kiện ngắt kết nối.
    -   Thực hiện vòng lặp `handleReconnect` với chiến lược thử lại (Retry).
    -   Đảm bảo các yêu cầu đăng ký không bị mất (Lossless) khi Broker tạm thời ngoại tuyến.

## 2. Xử lý đồng thời (Concurrency Control)

Hệ thống sử dụng mô hình **Pessimistic Locking** để ngăn chặn Overbooking.

*   **Tại Database:** Khi xử lý đăng ký từ Queue, Worker thực hiện `SELECT ... FOR UPDATE` trên hàng dữ liệu của Workshop.
*   **Tại Cache:** Sử dụng Redis Rate Limiter và Lua Script để lọc bớt yêu cầu rác trước khi đẩy vào Queue, giảm tải cho Database.

## 3. Tiêu chuẩn Dữ liệu & Timezone

*   Toàn bộ hệ thống thống nhất múi giờ **Asia/Ho_Chi_Minh** (ICT).
*   Mọi bản ghi `created_at` và `scanned_at` đều được lưu dưới dạng Timestamp có múi giờ để đảm bảo tính nhất quán khi truy vấn báo cáo.
*   Backend chịu trách nhiệm format thời gian sang định dạng ISO-8601 trước khi trả về cho Mobile/Web.

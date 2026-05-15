# Đặc tả Kỹ thuật: Circuit Breaker (Ngắt mạch tự động)

Tài liệu mô tả cơ chế bảo vệ hệ thống trước sự cố từ các dịch vụ bên thứ ba (AI Model, Payment Gateway).

## 1. Cơ chế hoạt động

Hệ thống triển khai mẫu thiết kế **Circuit Breaker** với 3 trạng thái cơ bản:

-   **CLOSED (Đóng):** Trạng thái bình thường. Mọi yêu cầu được gửi đi. Nếu tỷ lệ lỗi vượt quá ngưỡng (`threshold`), mạch sẽ chuyển sang **OPEN**.
-   **OPEN (Mở):** Trạng thái ngắt mạch. Mọi yêu cầu bị từ chối ngay lập tức (`Fail-fast`) để tránh treo tài nguyên và cho đối tác thời gian hồi phục. Sau một khoảng thời gian chờ (`sleepWindow`), mạch chuyển sang **HALF-OPEN**.
-   **HALF-OPEN (Nửa mở):** Trạng thái thăm dò. Hệ thống cho phép một số lượng giới hạn các yêu cầu đi qua.
    -   Nếu các yêu cầu này thành công liên tiếp, mạch chuyển về **CLOSED**.
    -   Nếu có bất kỳ yêu cầu nào thất bại, mạch quay lại **OPEN**.

## 2. Các điểm áp dụng thực tế

### A. AI Summary Service
-   **Ngưỡng lỗi:** 50% trong cửa sổ 60 giây.
-   **Thời gian ngủ:** 60 giây.
-   **Mục đích:** Tránh treo tiến trình worker khi Gemini API bị giới hạn lưu lượng (Rate limit) hoặc gặp sự cố mạng.

### B. Payment Gateway
-   **Ngưỡng lỗi:** 50% trong cửa sổ 30 giây.
-   **Thời gian ngủ:** 30 giây.
-   **Mục đích:** Đảm bảo hệ thống không tạo thêm registration mới cho các workshop có phí khi cổng thanh toán đang bảo trì, giúp giảm thiểu rủi ro dữ liệu rác.

## 3. Theo dõi trạng thái (Monitoring)

Hệ thống cung cấp API nội bộ cho Admin để kiểm tra trạng thái hiện tại của các bộ ngắt mạch:
-   `GET /api/v1/admin/payment/circuit-breaker`
-   `GET /api/v1/admin/payment/gateway-status` (Tích hợp cả Mock Toggle)

## 4. Lợi ích kiến trúc
-   **Resilience (Tính bền bỉ):** Hệ thống vẫn hoạt động ổn định dù một phần dịch vụ bên ngoài gặp sự cố.
-   **Resource Protection:** Ngăn chặn việc chiếm dụng kết nối và bộ nhớ vào các yêu cầu vô vọng.

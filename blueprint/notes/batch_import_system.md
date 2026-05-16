# Hệ thống Điều phối Dữ liệu hàng loạt (Batch Import System)

UniHub tích hợp một hệ thống Import dữ liệu sinh viên mạnh mẽ, tuân thủ mô hình **Batch Sequential Pipeline** và được tối ưu hóa cho hiệu suất cao.

## 1. Kiến trúc Hệ thống (Backend)

Hệ thống được xây dựng trên 4 giai đoạn chính:

*   **Extract**: Đọc dữ liệu từ tệp CSV được tải lên.
*   **Transform**: Kiểm tra tính hợp lệ của dữ liệu (MSSV, Email, trùng lặp). Bước này đã được tối ưu bằng cách bỏ qua việc băm mật khẩu (giả định mật khẩu trong file đã được băm sẵn).
*   **Load (High Performance)**: Sử dụng cơ chế **Bulk Upsert** (Insert on Conflict Update) trong một Database Transaction duy nhất để nạp hàng ngàn bản ghi chỉ trong vài giây.
*   **Cleanup**: Lưu trữ (Archive) tệp tin sau khi xử lý thành công để giải phóng không gian lưu trữ tạm thời.

## 2. Trung tâm Điều phối Dữ liệu (Frontend)

Giao diện quản trị cao cấp cho phép Admin kiểm soát toàn bộ luồng dữ liệu:

*   **Lập lịch thông minh**: Mặc định hệ thống tự động quét và xử lý vào lúc **02:00 sáng** hàng ngày để tránh giờ cao điểm.
*   **Điều khiển thủ công**: Nút **"Chạy ngay"** cho phép kích hoạt tiến trình xử lý ngay lập tức với phản hồi trạng thái thời gian thực.
*   **Báo cáo trạng thái**: Sử dụng các Badge màu sắc đồng bộ (`Pending`, `Processing`, `Completed`, `Failed`).

## 3. Cơ chế Xử lý và Hiển thị Lỗi

Một trong những điểm nhấn "Premium" của hệ thống là khả năng chẩn đoán lỗi chi tiết:

*   **Error Logging**: Mọi dòng dữ liệu lỗi đều được ghi lại vào bảng `import_errors` kèm theo lý do cụ thể (thiếu MSSV, sai định dạng email, trùng ID nội bộ).
*   **Chi tiết lỗi (Sheet UI)**: Khi một Job có lỗi, Admin có thể nhấn vào nút **"Lỗi"** để mở ngăn kéo chi tiết. Giao diện hiển thị số dòng, lý do lỗi và **dữ liệu thô (raw data)** giúp việc sửa lỗi trong file CSV trở nên cực kỳ dễ dàng.

## 4. Công cụ Hỗ trợ Dữ liệu (Tools)

Để phục vụ việc kiểm thử và vận hành, các script sau đã được cung cấp tại `src/data/scripts/`:

*   **`generate_sample.go`**: Sinh 5.000 bản ghi sinh viên mẫu với mật khẩu băm sẵn.
*   **`generate_error_sample.go`**: Sinh tệp tin mẫu chứa các dòng lỗi cố ý (thiếu ID, sai email, trùng lặp) để kiểm tra giao diện lỗi.
*   **`cleanup_test_users.go`**: (Chạy qua Backend) Giúp dọn dẹp toàn bộ dữ liệu thử nghiệm dựa trên thời gian tạo và tiền tố MSSV `2026`.

## 5. Lưu ý vận hành

*   **Đồng bộ giao diện**: Tất cả nút bấm và trạng thái (`Processing`) đều được đồng bộ màu `primary` của UniHub.
*   **An toàn dữ liệu**: Luôn ưu tiên cơ chế **Upsert** để không làm mất dữ liệu của sinh viên hiện có khi cập nhật danh sách mới.

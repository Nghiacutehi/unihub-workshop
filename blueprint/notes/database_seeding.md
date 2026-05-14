# Database Seeding & Testing Guide

Hướng dẫn thiết lập môi trường dữ liệu mẫu để kiểm thử toàn diện hệ thống UniHub.

## 1. Script Seed Dữ liệu
Toàn bộ logic tạo dữ liệu nằm tại file: `src/data/seed_data.sql`.

### Cách sử dụng:
1. Copy nội dung file `seed_data.sql`.
2. Dán vào SQL Editor của Supabase Dashboard.
3. Nhấn **Run**.

## 2. Các dữ liệu được tạo ra
Script sẽ tạo ra một môi trường giả lập với đầy đủ kịch bản:
- **100 Vé Đăng ký thành công:** Đã được gán ngẫu nhiên cho sinh viên và workshop.
- **Thanh toán:** Tự động tạo bản ghi `payments` cho các workshop có phí.
- **Thông báo:** Tạo lịch sử thông báo đã gửi cho sinh viên.
- **Nhật ký Import:** Tạo sẵn 1 Job thành công và 1 Job thất bại kèm danh sách lỗi mẫu.

## 3. Lưu ý về Ràng buộc (Constraints)
- **uq_user_workshop:** Một sinh viên không thể đăng ký 2 lần cho cùng 1 workshop.
- **uq_event_channel:** Một sự kiện chỉ có duy nhất 1 thông báo cho mỗi kênh (EMAIL). Script đã xử lý bằng lệnh `ON CONFLICT DO NOTHING`.
- **available_seats:** Luôn được cập nhật tự động sau khi insert registrations mẫu.

## 4. Tài khoản Test mặc định
- **Admin:** `admin` / pass: `123456`
- **Staff:** `staff01` / pass: `123456`
- **Student:** `21127001` - `21127017` / pass: `123456`

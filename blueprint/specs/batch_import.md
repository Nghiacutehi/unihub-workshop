# Đặc tả: batch_import.md

## 1. Mô tả
[cite_start]Tài liệu quy định luồng nghiệp vụ đồng bộ dữ liệu sinh viên từ hệ thống cũ vào `UniHub Workshop` thông qua tập tin `CSV` . [cite_start]Quy trình thực thi định kỳ vào ban đêm (Off-peak hours) để cập nhật danh sách sinh viên hợp lệ, phục vụ công tác xác thực và đăng ký workshop[cite: 40].

## 2. Kiến trúc giải pháp
Module áp dụng kiến trúc `Batch Sequential`. Dữ liệu được xử lý qua các giai đoạn nối tiếp, mỗi giai đoạn phải hoàn thành 100% trước khi chuyển sang giai đoạn kế tiếp:

* `Job Scheduler`: Sử dụng `Cron Job` hoặc `Quartz` để kích hoạt tiến trình tự động.
* `Staging Area`: Sử dụng bảng tạm (`Temporary Table`) trong `PostgreSQL` để lưu trữ dữ liệu thô trước khi làm sạch.
* `Transaction Management`: Áp dụng `Chunk-oriented Processing`. Dữ liệu được chia nhỏ thành từng lô (`Batch`) để thực thi nhằm tránh treo `Database Connection` và hỗ trợ `Rollback` từng phần khi gặp lỗi.

## 3. Workflow
Quy trình thực thi gồm 4 giai đoạn cốt lõi:

### 3.1. Giai đoạn Extract (Trích xuất)
1. [cite_start]`Background Worker` kiểm tra sự tồn tại của file `CSV` tại thư mục chỉ định[cite: 39].
2. Thực thi `File Validation`: Kiểm tra định dạng (Encoding), cấu trúc cột và dung lượng file.
3. Đọc dữ liệu và `Bulk Insert` vào `Staging Table`.

### 3.2. Giai đoạn Transform (Biến đổi và Làm sạch)
1. Kiểm tra tính hợp lệ của từng bản ghi (Đúng định dạng MSSV, Email, số điện thoại).
2. Xử lý trùng lặp: Loại bỏ các dòng trùng MSSV trong cùng một file.
3. Ánh xạ (`Mapping`) dữ liệu từ cấu trúc CSV vào `Domain Model` của hệ thống UniHub.

### 3.3. Giai đoạn Load (Nạp dữ liệu)
1. Thực hiện `Upsert` (Update or Insert) vào bảng `users` chính thức.
2. Sử dụng câu lệnh `ON CONFLICT (student_id) DO UPDATE` để cập nhật thông tin nếu sinh viên đã tồn tại, hoặc thêm mới nếu chưa có.
3. Thực thi theo từng `Chunk` (ví dụ: 1.000 bản ghi/giao dịch).

### 3.4. Giai đoạn Cleanup (Dọn dẹp)
1. Di chuyển file `CSV` đã xử lý vào thư mục `Archive` (Lưu trữ).
2. Xóa dữ liệu trong `Staging Table`.
3. Ghi log tổng kết: Số bản ghi thành công, số bản ghi lỗi, thời gian xử lý.

## 4. Xử lý ngoại lệ

| Loại lỗi | Giai đoạn | Hành vi hệ thống |
| :--- | :--- | :--- |
| `File Corrupted` | `Extract` | [cite_start]Ngừng tiến trình, bắn thông báo `Critical Error` tới `Organizer`, không tác động tới `Database`[cite: 55]. |
| `Data Type Mismatch` | `Transform` | Bỏ qua dòng lỗi (`Skip`), ghi nhận vào `Error Log` kèm số dòng cụ thể, tiếp tục xử lý các dòng còn lại. |
| `Database Deadlock` | `Load` | Thực hiện `Retry` giao dịch tối đa 3 lần. Nếu vẫn thất bại, thực hiện `Rollback Chunk` hiện tại và tiếp tục `Chunk` tiếp theo. |

## 5. Ràng buộc kỹ thuật
* `Schedule`: Mặc định thực thi vào lúc 02:00 AM hàng ngày.
* `Performance`: Phải xử lý tối thiểu 1.000 bản ghi/giây để đảm bảo hoàn thành trước giờ cao điểm sáng.
* `Logging`: Bắt buộc ghi lại `Tracking ID` cho mỗi đợt `Import` để có thể truy vết nguồn gốc dữ liệu khi xảy ra khiếu nại về thông tin sinh viên.

## 6. Tiêu chí chấp nhận
* [cite_start]Toàn bộ sinh viên mới có trong file `CSV` phải đăng nhập được vào hệ thống sau khi quá trình `Import` kết thúc[cite: 40].
* [cite_start]Không làm ảnh hưởng đến hiệu năng của các `API` xem workshop đang chạy đồng thời (Graceful Execution)[cite: 55].
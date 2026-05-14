# Chiến lược Offline-first & Đồng bộ dữ liệu

Tài liệu hướng dẫn về cơ chế lưu trữ bền vững tại máy khách (Client Persistence) và quy trình đồng bộ hóa (Synchronization).

## 1. Lưu trữ Local (SQLite)

Hệ thống sử dụng SQLite (`expo-sqlite`) để lưu trữ các lượt điểm danh khi thiết bị mất kết nối mạng (Offline).

### Bảo mật dữ liệu (AES-256)
- Toàn bộ dữ liệu nhạy cảm trong SQLite được mã hóa bằng thuật toán **AES-256** (thư viện `crypto-js`).
- Dữ liệu được mã hóa trước khi thực hiện lệnh `INSERT` và chỉ được giải mã khi cần hiển thị hoặc đồng bộ.
- **Mục tiêu:** Chống xem trộm dữ liệu nếu thiết bị bị chiếm quyền điều khiển (Root/Jailbreak).

### Cấu trúc bảng `local_checkins`
| Trường | Kiểu dữ liệu | Mô tả |
| :--- | :--- | :--- |
| `id` | TEXT (PK) | Mã UUIDv4 duy nhất cho mỗi lượt quét. |
| `data` | TEXT | JSON đã mã hóa AES (chứa `student_id`, `workshop_id`, ...). |
| `sync_status` | TEXT | Trạng thái: `PENDING` (Chờ) hoặc `SYNCED` (Đã đồng bộ). |
| `created_at` | INTEGER | Timestamp lúc quét mã. |

## 2. Quy trình Đồng bộ (Sync Workflow)

Hệ thống áp dụng cơ chế **Eventual Consistency** thông qua `useSync` hook.

1. **Chu kỳ:** Tiến trình chạy ngầm định kỳ mỗi **30 giây**.
2. **Bulk Upload:** 
   - App gom tối đa 500 bản ghi `PENDING`.
   - Gửi một Request duy nhất lên API `/api/checkin/sync` để tối ưu băng thông.
3. **Idempotency:** Backend sử dụng `id` (UUIDv4) để đảm bảo nếu một bản ghi bị gửi trùng 2 lần, hệ thống PostgreSQL sẽ không tạo dữ liệu rác.
4. **Dọn dẹp:** Chỉ sau khi nhận phản hồi thành công (HTTP 200) từ Server, App mới thực hiện xóa bản ghi khỏi SQLite.

## 3. Xử lý ngoại lệ

- **Mất mạng kéo dài:** App vẫn cho phép nhân viên quét mã bình thường, dung lượng lưu trữ tối đa phụ thuộc vào bộ nhớ điện thoại (khuyến nghị tối đa 10,000 bản ghi).
- **Lỗi Server (500):** App giữ nguyên trạng thái `PENDING` và thử lại ở chu kỳ tiếp theo.
- **Lỗi Dữ liệu (400):** Ghi Log và thông báo cho nhân viên nếu mã QR không hợp lệ sau khi đã lưu local.

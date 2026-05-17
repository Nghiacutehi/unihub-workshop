# 3. Database Design & Synchronization

## 3.1 Cấu Trúc PostgreSQL (Source of Truth)
Cơ sở dữ liệu được chuẩn hóa, quản lý dữ liệu gốc.

1. **Bảng Users:** `user_id (PK)`, `password_hash`, `role`, `email`, v.v.
2. **Bảng Workshops:** `id (PK)`, `title`, `available_seats`, `ai_summary`, v.v.
3. **Bảng Registrations:** Lịch sử đăng ký.
   - Primary Key: `id` (UUID).
   - Foreign Keys: `student_id`, `workshop_id`.
   - **Constraint quan trọng:** `UNIQUE(student_id, workshop_id)` chống sinh viên đăng ký 1 sự kiện nhiều lần.

## 3.2 Tối Ưu Nhập Liệu Lớn (Bulk Upsert)
Khi admin upload file CSV chứa 12,000 học sinh:
- API không dùng vòng lặp `for` để gọi `INSERT` 12,000 lần (Sẽ gây nghẽn Connection).
- Batch Import Service chia danh sách thành các mảng (chunks) 1000 records.
- Sử dụng tính năng `UNNEST` của PostgreSQL để truyền vào toàn bộ mảng và chạy `INSERT ... ON CONFLICT (user_id) DO UPDATE SET ...` trong 1 query duy nhất.
- Đảm bảo tốc độ nhập liệu mili-giây và an toàn lặp lại (Upsert Idempotent).

## 3.3 Offline-First Mobile Database (SQLite)
App di động cho Staff check-in phải hoạt động được dưới tầng hầm (nơi không có wifi/3G).

**Cấu trúc bảng SQLite (App):**
- Bảng `Checkins`: `id`, `student_id`, `workshop_id`, `status` (PENDING / SYNCED), `timestamp`.

**Luồng dữ liệu:**
1. Khi có Internet: App kéo dữ liệu danh sách Workshop về lưu tạm.
2. Staff mang máy xuống hầm. Quét QR Code. App verify chữ ký số RSA hoàn toàn Offline. 
3. Nếu vé hợp lệ, lưu vào SQLite với `status = 'PENDING'`.
4. Staff mang máy lên mặt đất. Background Task `useSync.ts` quét định kỳ 30s/lần.
5. Thấy có record PENDING -> Bắn API Bulk Sync `/api/v1/checkin/sync` lên Backend.
6. Backend trả về OK -> SQLite chuyển record thành `SYNCED`.

# Quy chuẩn kỹ thuật Backend & Database

Tài liệu này ghi lại các quyết định kỹ thuật về cách lưu trữ và xử lý dữ liệu để đảm bảo tính đồng nhất giữa Web, Mobile và Server.

## 1. Xử lý Ngày tháng & Múi giờ (DateTime)

Hệ thống sử dụng múi giờ Việt Nam (**UTC+7**) làm chuẩn hiển thị, nhưng lưu trữ theo chuẩn quốc tế.

*   **Kiểu dữ liệu:** Luôn sử dụng `TIMESTAMPTZ` (Timestamp with time zone) trong PostgreSQL.
*   **Quy tắc Insert/Update:** 
    *   Luôn đính kèm múi giờ khi thao tác SQL: `'2026-05-14 08:00:00+07'`.
    *   Tránh sử dụng giờ hệ thống (`now()`) nếu không chắc chắn Server đang đặt ở đâu.
*   **Quy tắc hiển thị (Client):** 
    *   Mobile và Web sẽ nhận dữ liệu dạng chuỗi ISO từ Supabase.
    *   Sử dụng `new Date(string)` để tự động chuyển đổi về giờ địa phương của thiết bị người dùng.

## 2. Trạng thái Workshop (Status logic)

Trạng thái của một Workshop được xác định qua 2 tầng:

1.  **Tầng Database (`status` column):**
    *   `DRAFT`: Nháp (Không hiển thị cho sinh viên/nhân sự).
    *   `PUBLISHED`: Đã xuất bản (Hiển thị trên hệ thống).
2.  **Tầng Ứng dụng (Dynamic Status):**
    *   Logic này được tính toán "on-the-fly" dựa trên thời gian thực:
        *   `In Progress`: `start_time <= now <= end_time`.
        *   `Upcoming`: `now < start_time`.
        *   `Completed`: `now > end_time`.

## 3. Xác thực & Phân quyền (Auth & RBAC)

*   **User ID:** Cột `id` trong bảng `users` phải khớp hoàn toàn với `UID` do Supabase Auth cấp phát.
*   **Roles:** Hệ thống chỉ chấp nhận 3 Role chuẩn: `STUDENT`, `STAFF`, `ADMIN`.
*   **Mobile Access:** Chỉ có `STAFF` và `ADMIN` mới có quyền sử dụng App Mobile Check-in.

## 4. Bảo mật QR & Xác thực (Security)

Hệ thống sử dụng cơ chế chữ ký số để chống gian lận:
*   **Thuật toán:** Ed25519 (Asymmetric Signing).
*   **Private Key:** Lưu tại biến môi trường Backend, dùng để ký vé.
*   **Public Key:** Cung cấp qua Endpoint `GET /api/auth/keys` cho Client xác thực.

## 5. Đồng bộ hóa & Tính toàn vẹn (Sync & Idempotency)

Khi tiếp nhận dữ liệu từ Mobile (đặc biệt là dữ liệu Offline):
*   **UUID v4:** Sử dụng `UUID v4` cho toàn bộ `id` bản ghi điểm danh để chống trùng lặp dữ liệu (Idempotency) khi Client thực hiện Retry Sync.
*   **Conflict Resolution:** Nếu phát hiện 1 sinh viên điểm danh cùng 1 workshop tại nhiều thiết bị khác nhau, Server ưu tiên bản ghi có `scanned_at` nhỏ nhất.
*   **Bulk API:** Endpoint `/api/checkin/sync` phải hỗ trợ xử lý hàng loạt bản ghi (tối đa 500) trong một Transaction để đảm bảo hiệu năng.

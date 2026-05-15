# Quản lý Phiên & Bộ nhớ tạm (Auth & Caching)

Tài liệu này quy định cách hệ thống xử lý trạng thái người dùng và tối ưu dữ liệu trên trình duyệt.

## 1. Cơ chế Xác thực (JWT)

Hệ thống sử dụng JWT (JSON Web Token) để xác thực người dùng.

*   **Lưu trữ:** Token được lưu trong `Http-Only Cookie` (tên: `unihub_token`).
*   **Truyền tải:** Gửi qua Header `Authorization: Bearer <token>` trong mọi request.
*   **Thời hạn:** 24 giờ kể từ lúc đăng nhập thành công.

## 2. Quản lý Session trên Frontend

Để tăng tốc độ hiển thị UI, thông tin User cơ bản được lưu trong cookie `unihub_session`.

*   **Dữ liệu:** `{ id, name, student_id }`.
*   **Xử lý 401 (Unauthorized):**
    -   Khi API trả về lỗi 401, hệ thống tự động gọi `auth.clearSession()`.
    -   Xóa sạch các cookie `unihub_token` và `unihub_session`.
    -   Chuyển hướng người dùng về trang `/login`.

## 3. Chiến lược Caching (User-Specific Cache)

Để tránh rò rỉ dữ liệu khi nhiều người dùng chung một máy tính, hệ thống áp dụng Cache theo User ID.

*   **Key định danh:** `unihub_workshops_cache_{userId}`.
*   **Vị trí:** `localStorage`.
*   **Quy trình làm mới:**
    1.  Khi Hook `useWorkshops` khởi chạy, nó kiểm tra `userId` hiện tại.
    2.  Nếu `userId` thay đổi (đăng nhập tài khoản khác), toàn bộ dữ liệu trong State và `localStorage` cũ sẽ bị xóa bỏ.
    3.  Hệ thống thực hiện Fetch mới và ghi đè vào Cache của User mới.

## 4. Đồng bộ thời gian (Timezone)

*   **Múi giờ hệ thống:** `Asia/Ho_Chi_Minh` (UTC+7).
*   **Quy tắc:**
    -   Backend (Go) chuyển đổi mọi mốc thời gian sang ICT trước khi trả về JSON.
    -   Frontend (Next.js) định dạng hiển thị `DD/MM/YYYY HH:mm` theo chuẩn Việt Nam.

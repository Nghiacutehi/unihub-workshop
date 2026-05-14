# UniHub Authentication & Authorization

Tài liệu đặc tả cơ chế bảo mật và phân quyền của hệ thống UniHub Workshop.

## 1. Cơ chế Xác thực (Authentication)

Hệ thống sử dụng cơ chế xác thực dựa trên dữ liệu người dùng được lưu trữ trực tiếp trong bảng `public.users` của PostgreSQL.

- **Định danh:** Người dùng có thể đăng nhập bằng **Email** hoặc **Mã sinh viên (user_id)**.
- **Mật khẩu:** Hiện tại đang lưu trữ dưới dạng bản rõ (plaintext) cho mục đích phát triển. 
  - *Lưu ý:* Cần triển khai `bcrypt` hoặc `argon2` để băm mật khẩu trước khi demo thực tế.
- **Session:** Sau khi đăng nhập thành công, thông tin người dùng được lưu vào **Cookie** tên là `unihub_session` dưới dạng JSON đã mã hóa URI.

## 2. Phân quyền (Authorization)

Hệ thống phân chia làm 3 vai trò chính dựa trên cột `role`:

| Vai trò | Quyền hạn | Route truy cập |
| :--- | :--- | :--- |
| **ADMIN** | Quản lý toàn bộ hệ thống, workshops, người dùng. | `/admin/*`, `/` |
| **STAFF** | Quản lý workshop cụ thể, quét mã QR điểm danh. | `/admin/*` (giới hạn), Mobile App |
| **STUDENT** | Xem workshop, đăng ký, nhận vé QR. | `/` (Home), Mobile App (User mode) |

### Bảo vệ Route (Middleware)
- Toàn bộ route `/admin/*` được bảo vệ bởi `src/web/middleware.ts`.
- Nếu người dùng chưa đăng nhập hoặc có `role = 'STUDENT'`, hệ thống sẽ tự động điều hướng về trang chủ hoặc trang đăng nhập.

## 3. Bảo mật trên Mobile
- App Mobile sử dụng `AsyncStorage` để lưu trữ Session.
- **Staff App** chỉ cho phép người dùng có role `ADMIN` hoặc `STAFF` đăng nhập. Nếu sinh viên cố tình đăng nhập vào app của nhân viên, hệ thống sẽ báo lỗi và từ chối cấp quyền.

## 4. Kế hoạch bảo mật tiếp theo
- Triển khai JWT (JSON Web Token) để thay thế việc lưu JSON thuần trong Cookie.
- Kết nối với Backend Go để thực hiện xác thực tập trung.
- Triển khai mã hóa RSA cho mã QR để chống làm giả vé.

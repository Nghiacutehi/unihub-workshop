# Đặc tả: Phân quyền và Bảo mật (auth.md)

## 1. Mô tả
Tài liệu này quy định cơ chế xác thực (Authentication) và ủy quyền (Authorization). Hệ thống triển khai Stateless Authentication sử dụng JWT và mô hình RBAC. Logic xác thực và ủy quyền được tách biệt thành các module độc lập tại tầng Middleware. Còn tài khoản và mật khẩu để xác thực sinh viên đã được lưu trữ qua file CSV



## 2. Ma trận Phân quyền (Role-Permission Matrix)
Dưới đây là định nghĩa chi tiết các vai trò và quyền hạn tương ứng trên các tài nguyên hệ thống để đảm bảo tính nhất quán trong lập trình:

| Tài nguyên (Resource) | Student | Staff (Check-in) | Organizer (Admin) |
| :--- | :--- | :--- | :--- |
| Workshop Information | Read | Read | CRUD |
| Registration | Create | Read | Read / Update |
| QR Check-in | Deny | Create (Scan) | Read |
| AI Summary | Read | Deny | Create |
| CSV Data Import | Deny | Deny | Create (Execute) |
| User Statistics | Deny | Deny | Read |

## 3. Quy trình Xác thực và Ủy quyền
### 3.1. Giai đoạn Xác thực (Authentication)
1. Người dùng gửi thông tin định danh qua API `POST /auth/login`.
2. Hệ thống kiểm tra mật khẩu đã băm (Hashed password) trong PostgreSQL.
3. Nếu khớp, sinh JWT chứa:
    - `sub`: Định danh người dùng (MSSV/Username).
    - `role`: Vai trò của người dùng (STUDENT/STAFF/ORGANIZER).
    - `exp`: Thời gian hết hạn (1 giờ kể từ lúc tạo).

### 3.2. Giai đoạn Ủy quyền (Authorization)
1. API Gateway trích xuất Bearer Token từ Header của request.
2. Thực hiện giải mã và xác minh chữ ký bằng `JWT_SECRET_KEY`.
3. Kiểm tra trường `exp`. Nếu hết hạn, hủy bỏ yêu cầu.
4. Trích xuất trường `role` và đối chiếu với yêu cầu của API Endpoint:
    - Nếu Role hợp lệ: Chuyển tiếp request đến dịch vụ xử lý.
    - Nếu Role không hợp lệ: Trả về lỗi `403 Forbidden`.

## 4. Kịch bản lỗi
* **Sai thông tin đăng nhập:** Trả về `401 Unauthorized`.
* **Token bị chỉnh sửa:** Chữ ký không khớp, trả về `401 Unauthorized`.
* **Truy cập trái quyền:** Ví dụ Sinh viên truy cập API xóa workshop, trả về `403 Forbidden`.
* **Token hết hạn:** Trả về `401 Unauthorized` kèm thông báo cần đăng nhập lại.

## 5. Ràng buộc kỹ thuật
* **Algorithm:** Sử dụng HS256 để ký Token.
* **Storage:** Client lưu JWT tại HttpOnly Cookie.Khóa bí mật của server được lưu trong .env
* **Env Management:** Khóa bí mật phải được lưu trong biến môi trường `AUTH_SECRET`.
* **State:** Tuyệt đối không sử dụng `HttpSession` hoặc lưu trữ trạng thái tại RAM máy chủ.
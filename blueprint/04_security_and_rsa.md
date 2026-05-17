# 4. Security: RSA-2048 Digital Signatures

Vấn đề sống còn: Sinh viên có thể fake (làm giả) mã QR bằng cách tạo một chuỗi JSON có chứa MSSV của mình để qua mặt Staff. Làm thế nào để phân biệt QR chuẩn và QR tự chế offline (không cần gọi API check)?

## Cách giải quyết: Chữ Ký Số Phi Đối Xứng (Asymmetric Cryptography)

Hệ thống sử dụng bộ khóa RSA-2048.
- **Private Key (Bí mật):** Chỉ nằm duy nhất ở Backend. Dùng để "Ký".
- **Public Key (Công khai):** Nằm ở App Mobile của Staff. Dùng để "Xác minh".

### Luồng Hoạt Động

1. **Sinh mã (Tại Backend khi đăng ký thành công):**
   - Payload gốc: `{"sid":"21127100", "wid":"workshop-uuid-1234"}`
   - Backend tính mã băm SHA-256 của payload này.
   - Backend sử dụng **Private Key** để mã hóa chuỗi băm đó tạo thành một chuỗi lòng vòng gọi là Chữ ký (Signature) dạng Base64.
   - Trả về cho Web hiển thị mã QR chứa: `{ Payload Gốc + Signature }`.

2. **Quét mã & Kiểm tra (Tại App Mobile Offline):**
   - Staff quét QR lấy được chuỗi JSON.
   - Mobile bóc tách Payload Gốc và Signature.
   - Mobile tự tính lại mã băm SHA-256 của Payload Gốc.
   - Mobile sử dụng **Public Key** giải mã cái Signature kia. Nó sẽ ói ra một chuỗi mã băm.
   - So sánh: Nếu mã băm Mobile tự tính **TRÙNG KHỚP** với mã băm vừa giải mã được -> QR này chắc chắn được tạo ra từ Backend. VÉ THẬT!
   - Kẻ gian có thể fake Payload, nhưng vì không có Private Key, hắn không thể tạo ra Signature khớp với Payload.

## Bảo Mật Authentication
Toàn bộ API được bảo vệ bằng JSON Web Token (JWT) và Role-Based Access Control (RBAC). 
- Route `/api/v1/admin/*` yêu cầu token của Role `ADMIN`.
- Route `/api/v1/registrations/*` yêu cầu token của Role `STUDENT`.

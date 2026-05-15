# Đặc tả Kỹ thuật: Bảo mật RSA & Offline Verification

Tài liệu chi tiết về cơ chế ký số (Signing) và xác thực vé (Verification) không cần mạng trên ứng dụng Mobile.

## 1. Cơ sở hạ tầng khóa (Key Infrastructure)

Hệ thống sử dụng cặp khóa RSA (2048-bit):
-   **Private Key:** Lưu trữ an toàn tại Server (Biến môi trường `RSA_PRIVATE_KEY`). Dùng để ký số dữ liệu vé sau khi đăng ký/thanh toán thành công.
-   **Public Key:** Cung cấp công khai qua API. Mobile App dùng để xác thực tính toàn vẹn của mã QR.

## 2. Quy trình Đồng bộ Public Key

Mobile App thực hiện lấy Public Key trong các trường hợp:
1.  **Lần đầu khởi chạy:** Khi `AsyncStorage` chưa có khóa.
2.  **Verify thất bại:** Nếu quét QR bị báo lỗi chữ ký, App sẽ tự động thử đồng bộ lại khóa mới trước khi kết luận vé giả.
3.  **Thao tác thủ công:** Người dùng nhấn "Đồng bộ" trong trang Lịch sử/Cài đặt.

**Endpoint:** `GET /api/v1/auth/public-key`

## 3. Quy trình Xoay vòng khóa (Key Rotation)

Để đảm bảo an toàn tối đa, quy trình rotate key được thực hiện theo các bước sau:

1.  **Server side:**
    -   Tạo cặp khóa RSA mới.
    -   Cập nhật `RSA_PRIVATE_KEY` trong file cấu hình/môi trường của Backend.
    -   Khởi động lại Backend. API `/public-key` sẽ trả về khóa mới ngay lập tức.
2.  **Data side:**
    -   Các vé cũ đã ký bằng khóa cũ vẫn tồn tại trong DB của sinh viên.
    -   *Lưu ý:* Khi nhân viên Staff mở App Mobile, họ **phải** nhấn đồng bộ hoặc khởi động lại App để nạp Public Key mới nhất.
3.  **Compatibility (Tính tương thích):**
    -   Hệ thống hiện tại chỉ hỗ trợ **1 Key Active**.
    -   *Khuyến nghị:* Thực hiện rotate key vào khung giờ thấp điểm (đêm khuya) để tránh gián đoạn việc quét vé tại sự kiện.
    -   *Nâng cao (Future):** Triển khai mảng khóa (Key Array) để hỗ trợ Verify cả khóa cũ và khóa mới trong giai đoạn chuyển đổi.

## 4. Cấu trúc dữ liệu ký (Signature Format)

Dữ liệu gốc (`Payload`) để ký:
`student_id | user_uuid | workshop_id` (Nối bằng dấu gạch đứng `|`)

Mã QR chứa JSON:
```json
{
  "sid": "MSSV",
  "uid": "User UUID",
  "wid": "Workshop UUID",
  "sig": "Base64 Signature"
}
```
Mobile App thực hiện băm SHA256 chuỗi `sid|uid|wid` và dùng Public Key để Verify với trường `sig`.

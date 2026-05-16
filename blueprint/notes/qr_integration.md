# Hệ thống Tích hợp QR & Điểm danh Offline-First (UniHub)

Tài liệu này quy định chuẩn định dạng mã QR và cơ chế ký số bảo mật cho hệ thống vé điện tử UniHub.

## 1. Cấu trúc mã QR (JSON Payload)

Mã QR chứa chính xác 4 trường dữ liệu để tối ưu tốc độ quét và độ bảo mật.

| Trường | Mô tả | Ví dụ |
| :--- | :--- | :--- |
| `sid` | Student ID (MSSV) | `HE150123` |
| `uid` | User UUID | `550e8400-e29b-41d4-a716-446655440000` |
| `wid` | Workshop ID | `workshop-001` |
| `sig` | RSA Signature (Base64) | `A2b5... (256+ bytes)` |

**Ví dụ JSON:**
```json
{
  "sid": "HE150123",
  "uid": "550e8400-e29b-41d4-a716-446655440000",
  "wid": "workshop-001",
  "sig": "iPZCXeWuDsEvJ0vcUL8I874dsEtRkT+nkukiQLFI..."
}
```

## 2. Cơ chế Ký số (Backend)

*   **Thuật toán:** RSA-2048 với SHA-256 (RSASSA-PKCS1-v1_5).
*   **Dữ liệu thô (Raw Data):** Chữ ký được tạo ra từ chuỗi kết hợp theo định dạng: `sid|uid|wid`.
*   **Quy trình:**
    1.  Tạo chuỗi `rawData = sid + "|" + uid + "|" + wid`.
    2.  Băm `rawData` bằng SHA-256.
    3.  Ký bản băm bằng **RSA Private Key**.
    4.  Encode kết quả sang Base64 để đưa vào trường `sig`.

## 3. Cơ chế Xác thực (Mobile Staff App)

Xác thực được thực hiện 100% Offline trên thiết bị của nhân viên để đảm bảo tốc độ.

1.  **Lấy Public Key:** Mobile tải Public Key từ `/api/v1/auth/public-key` và lưu vào `AsyncStorage`.
2.  **Verify:**
    -   Parse mã QR để lấy 4 trường.
    -   Tạo lại chuỗi `rawData` từ `sid`, `uid`, `wid`.
    -   Sử dụng Public Key để giải mã `sig` và đối soát với bản băm của `rawData`.
3.  **Kết quả:** Nếu khớp, vé hợp lệ. Nếu không, hiển thị cảnh báo "VÉ GIẢ MẠO".

## 4. Bảo mật Offline

*   Hệ thống ngăn chặn việc sử dụng lại vé (Double Spending) bằng cách lưu lịch sử quét vào **SQLite** nội bộ trên điện thoại.
## 5. Mã QR Thanh toán (VietQR - Napas 247)

Khác với mã QR vé (Ticket QR) dùng để điểm danh, hệ thống sử dụng mã QR thanh toán động để sinh viên hoàn tất đăng ký Workshop có phí.

*   **Chuẩn:** VietQR (Napas 247).
*   **Loại:** QR Động (Dynamic QR) - Tự động điền số tiền và nội dung.
*   **Tham số tích hợp:**
    *   **Beneficiary:** Ngân hàng đối tác & Số tài khoản của UniHub.
    *   **Amount:** Lấy trực tiếp từ thuộc tính `price` của Workshop.
    *   **AddInfo (Nội dung):** Tuân thủ cú pháp `UNIPAY [Transaction_ID]` để hệ thống tự động đối soát Webhook.
*   **Tích hợp Frontend:** Sử dụng API `vietqr.io` để sinh ảnh QR thời gian thực, đảm bảo tính tiện dụng (quét bằng App ngân hàng bất kỳ).

# Giao thức Check-in Offline (Mobile Staff App)

Tài liệu này hướng dẫn cách ứng dụng Mobile xử lý việc quét vé và xác thực khi không có mạng.

## 1. Đồng bộ Public Key

Để Verify Offline, Mobile App cần nạp Public Key của Server.

*   **API:** `GET /api/v1/auth/public-key`.
*   **Thời điểm gọi:** 
    -   Khi App khởi động (Splash Screen).
    -   Khi người dùng nhấn nút "Đồng bộ" trong Sync Screen.
*   **Lưu trữ:** Lưu vào `AsyncStorage` (Key: `@unihub_public_key`).

## 2. Luồng xử lý quét vé (Scanning Flow)

Khi nhân viên quét một mã QR:

1.  **Parse JSON:** Lấy `sid`, `uid`, `wid`, `sig`.
2.  **Verify Signature:** 
    -   Sử dụng thư viện `jsrsasign` để Verify chuỗi `sid|uid|wid` với `sig` bằng Public Key đã lưu.
    -   Nếu thất bại -> Hiển thị lỗi "CHỮ KÝ SAI/GIẢ MẠO".
3.  **Check Workshop ID:** So sánh `wid` trong QR với ID của Workshop hiện tại mà nhân viên đang check-in.
    -   Nếu không khớp -> Hiển thị lỗi "SAI WORKSHOP".
4.  **Check Duplicate (Offline):** Truy vấn trong SQLite nội bộ xem `sid` này đã được quét cho `wid` này chưa.
    -   Nếu đã có -> Hiển thị lỗi "VÉ ĐÃ SỬ DỤNG".
5.  **Save Local:** Nếu mọi thứ hợp lệ, lưu bản ghi vào SQLite với trạng thái `PENDING`.

## 3. Đồng bộ dữ liệu về Server (Back-sync)

Khi có mạng trở lại, ứng dụng tự động đẩy các bản ghi `PENDING` về Server thông qua Queue hoặc API đồng bộ hàng loạt (Bulk Sync).

*   **Lưu ý:** Chỉ sau khi đồng bộ thành công, bản ghi trong SQLite mới chuyển sang trạng thái `SYNCED`.

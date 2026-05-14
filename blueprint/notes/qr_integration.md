# Hệ thống Tích hợp QR & Điểm danh Offline-First (UniHub)

Tài liệu này ghi lại kiến trúc, luồng xử lý và các cơ chế vận hành thực tế của hệ thống điểm danh QR sau khi đã tối ưu hóa.

## 1. Kiến trúc Tổng quan (Architecture)
Hệ thống tuân thủ nguyên tắc **Offline-First**, đảm bảo tốc độ phản hồi < 50ms ngay cả khi không có mạng.

*   **Client (Mobile):** Quét mã -> Kiểm tra ID & Phòng (Local) -> Kiểm tra trùng lập (Local + Global Cache) -> Lưu SQLite -> Đồng bộ ngầm.
*   **Backend (Next.js API):** Nhận gói đồng bộ -> Cập nhật Database theo UUID -> Trả về kết quả xác nhận.
*   **Database:** SQLite (Mobile) & PostgreSQL (Server) đồng bộ trạng thái qua `sync_status`.

## 2. Cấu trúc dữ liệu mã QR
Mã QR sử dụng định dạng JSON thu gọn để tối ưu tốc độ nhận diện. Lớp chữ ký số đã được gỡ bỏ ở Client để ưu tiên hiệu năng:
```json
{
  "sid": "21127003",                    // MSSV (Hiển thị cho nhân viên đối chiếu)
  "uid": "ae94d6d6-5d64-4b6e-b81b...", // User UUID (Khóa chính hệ thống - Xử lý siêu tốc)
  "wid": "22222222-0000-0000-0000...", // Workshop ID (Dùng để kiểm tra đúng phòng/đúng buổi)
  "sig": "..."                         // (Legacy) Chỉ dùng cho xác thực phía Server nếu cần
}
```

## 3. Luồng xử lý quét mã (Scanning Logic)
Quy trình kiểm tra 3 bước cực nhẹ tại Client:
1.  **Đúng phòng (Room Check):** So sánh `wid` trong QR với `workshopId` đang mở. Chặn tuyệt đối việc quét nhầm vé của Workshop khác.
2.  **Chống quét trùng (Duplicate Check):** Tra cứu MSSV trong bảng `already_checked_in` (Toàn cục) và `local_checkins` (Tại thiết bị).
3.  **Ghi nhận (Storage):** Lưu vào SQLite với trạng thái `PENDING` và phản hồi tức thì cho nhân viên.

## 4. Cơ chế Đồng bộ (Sync Mechanism)
*   **Trạng thái SYNCED:** Bản ghi chuyển từ `PENDING` sang `SYNCED` sau khi Server xác nhận. Dữ liệu vẫn được giữ lại trong ca làm việc để tra cứu.
*   **Global Cache:** Mỗi khi mở Workshop, App tải danh sách những người đã điểm danh trước đó từ Server về thiết bị (Bảng `already_checked_in`). Điều này giúp chặn gian lận giữa các cổng soát vé khác nhau.
*   **Auto-Cleanup:** Tự động dọn dẹp các bản ghi của Workshop đã kết thúc > 24 giờ để tối ưu bộ nhớ.

## 5. Giao diện & Trải nghiệm (UI/UX)
Hệ thống sử dụng ngôn ngữ thiết kế **Midnight Indigo** đồng bộ toàn hệ thống:
*   **Màu chủ đạo:** `#312E81` (Dùng cho Nút bấm chính, Header và Logo).
*   **Phản hồi trạng thái:**
    *   **Xanh lá (#10B981):** Thành công (Vé hợp lệ).
    *   **Vàng/Cam (#F59E0B):** Cảnh báo (Vé đã quét rồi).
    *   **Đỏ (#EF4444):** Lỗi (Sai phòng hoặc mã QR không hợp lệ).

## 6. Lưu ý Vận hành
*   **Kết nối:** Mobile App kết nối với Backend qua địa chỉ IP Local (Cấu hình tại `services/crypto.ts`).
*   **Hiệu năng:** Đã loại bỏ polyfill Crypto giúp giảm dung lượng App và tăng tính ổn định trên các dòng máy cũ.

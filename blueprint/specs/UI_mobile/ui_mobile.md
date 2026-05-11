# Đặc tả: Giao diện Mobile App Check-in (ui_mobile.md)

## 1. Mô tả
Tài liệu quy định cấu trúc màn hình, luồng tương tác và trải nghiệm người dùng (`UX`) cho ứng dụng di động dành riêng cho nhân sự (`Staff`). Ứng dụng tập trung vào hiệu năng quét mã `QR Code` và khả năng vận hành trong điều kiện mất kết nối mạng (`Offline-first`).

## 2. Design System & UI/UX (Đồng nhất với hệ thống Web)
Để đảm bảo tính chuyên nghiệp và nhận diện thương hiệu UniHub, ứng dụng di động tuân thủ các quy tắc thiết kế từ phân hệ Web.

### 2.1. Bảng màu (Color Palette)
* **Màu chủ đạo (Primary):** `Cobalt Blue` (#1E40AF). Sử dụng cho Header và các nút bấm hành động chính.
* **Màu nền (Background):** `Slate 50` (#F8FAFC) để tạo cảm giác sạch sẽ, hiện đại.
* **Màu trạng thái (Semantic):**
    * `Success`: Emerald Green (#059669) - Hiển thị khi quét QR thành công.
    * `Error`: Rose Red (#E11D48) - Hiển thị khi vé không hợp lệ hoặc đã quét rồi.
    * `Warning`: Amber Orange (#D97706) - Hiển thị trạng thái đang chờ đồng bộ (Offline).

### 2.2. Quy tắc Component
* **Nút bấm:** Bo góc `rounded-md` (8px). Sử dụng hiệu ứng phản hồi khi chạm (`Touch Feedback`).
* **Khoảng trắng:** Tuân thủ hệ thống lưới **8-point**, đảm bảo khoảng cách giữa các phần tử tối thiểu 16px để dễ dàng thao tác bằng ngón tay.
* **Typography:** Sử dụng font `Inter` hoặc hệ thống font mặc định của thiết bị (San Francisco/Roboto) để tối ưu tốc độ đọc.

---

## 3. Cấu trúc màn hình (Screens)

### 3.1. Màn hình Đăng nhập (Login)
* **Giao diện:** Đơn giản, tập trung vào form nhập mã nhân sự và mật khẩu.
* **Bảo mật:** Lưu trữ `JWT` an toàn trên thiết bị. Chỉ tài khoản có `role: STAFF` mới được phép vào ứng dụng.

### 3.2. Màn hình Chọn Workshop (Workshop Selection)
* **Mô tả:** Danh sách các Workshop đang hoặc sắp diễn ra trong ngày.
* **Tính năng:** Nhân sự chọn đúng Workshop mình đang trực để bắt đầu quét mã. Mỗi Card hiển thị: Tên Workshop, Phòng, và số lượng sinh viên đã check-in thực tế.

### 3.3. Màn hình Quét QR (QR Scanner - Core Feature)
* **Giao diện:**
    * Vùng quét (Viewfinder) nằm chính giữa màn hình.
    * Đèn Flash: Nút bật/tắt để hỗ trợ quét trong điều kiện thiếu sáng.
    * Thông số: Hiển thị nhanh số lượng `Checked-in / Total` ở góc trên.
* **Phản hồi tức thì (Instant Feedback):**
    * **Thành công:** Màn hình nháy xanh nhẹ, rung nhẹ (`Haptic Feedback`) và phát tiếng "Beep" ngắn.
    * **Thất bại:** Màn hình nháy đỏ, rung mạnh hơn và phát tiếng cảnh báo. Hiển thị lý do lỗi (Vé sai, Đã quét, Sai Workshop).

### 3.4. Màn hình Lịch sử & Đồng bộ (Sync Status)
* **Mô tả:** Danh sách các lượt quét vừa thực hiện.
* **Trạng thái đồng bộ:**
    * Biểu tượng **Đám mây xanh**: Đã lưu lên server.
    * Biểu tượng **Đồng hồ vàng**: Đang lưu tại máy (`Offline Mode`), chờ có mạng để đồng bộ.
* **Nút bấm:** `[Đồng bộ ngay]` (Chỉ hiện khi có kết nối mạng trở lại).

---

## 4. Xử lý Trải nghiệm Offline (AP Architecture)
Để đảm bảo nhân sự không bị gián đoạn khi mạng yếu, giao diện phải thể hiện rõ các trạng thái sau:

| Trạng thái mạng | Phản hồi giao diện (UI Response) |
| :--- | :--- |
| **Mất kết nối** | Hiển thị Banner cảnh báo: "Đang hoạt động Ngoại tuyến". Mọi dữ liệu quét được ghi trực tiếp vào `SQLite`. |
| **Quét Offline** | Sau mỗi lượt quét, hiển thị thông báo: "Đã lưu cục bộ". |
| **Có mạng lại** | Tự động thực hiện `Background Sync` và hiển thị thông báo: "Đã đồng bộ thành công {x} dữ liệu lên hệ thống". |

---

## 5. Tiêu chí chấp nhận & Hiệu năng
* **Tốc độ:** Thời gian từ lúc Camera nhận diện mã QR đến khi lưu vào `SQLite` và phản hồi UI phải dưới **100ms**.
* **Độ ổn định:** App không được văng (Crash) khi quét liên tục hoặc khi chuyển đổi trạng thái mạng.
* **Pin:** Tối ưu hóa việc sử dụng Camera để không gây nóng máy hoặc hao pin quá nhanh trong ca trực 4 tiếng.
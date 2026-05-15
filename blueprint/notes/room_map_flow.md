# Room Map Flow & Performance Strategy

Tài liệu này mô tả chi tiết luồng xử lý sơ đồ phòng (Room Layout Map) trong hệ thống UniHub Workshop, bao gồm cơ chế hiển thị, xử lý lỗi và bảo vệ hệ thống khi quá tải.

## 1. Tổng quan (Overview)
Tính năng sơ đồ phòng cho phép sinh viên xem sơ đồ chi tiết của địa điểm tổ chức workshop ngay trong ứng dụng. Vì ảnh sơ đồ có dung lượng lớn, hệ thống áp dụng cơ chế "Smart Loading" để đảm bảo hiệu năng tối ưu.

## 2. Luồng xử lý chi tiết (System Flow)

### A. Phía Server (Go Backend)
Server chịu trách nhiệm quyết định việc có cung cấp dữ liệu ảnh hay không dựa trên tình trạng tải thực tế.
- **Vị trí xử lý:** `internal/service/workshop_service.go`
- **Cơ chế:** Kiểm tra lưu lượng truy cập hiện tại (`currentRequestRate`).
- **Ngưỡng giới hạn (Threshold):** **12,000 requests/giây**.
- **Logic:**
  - Nếu `Tải < 12,000`: Trả về đầy đủ thông tin bao gồm `room_layout_url`.
  - Nếu `Tải > 12,000`: Gán `room_layout_url = null` trước khi trả về JSON.
- **Mục tiêu:** Tiết kiệm băng thông (bandwidth) và giảm tải cho Cloud Storage/Static Server khi hệ thống đang ở đỉnh điểm truy cập.

### B. Phía Client (Student Web)
Frontend tập trung vào trải nghiệm người dùng và bắt lỗi hiển thị.

#### 1. Kiểm tra dữ liệu (Data Validation)
- Khi người dùng nhấn nút **"Mở sơ đồ"**:
  - Kiểm tra nếu `roomLayoutUrl` tồn tại: Hiển thị Overlay sơ đồ.
  - Nếu không tồn tại (do Server chặn hoặc DB chưa có): Hiển thị thông báo lỗi qua `toast.error`: *"Sơ đồ phòng cho workshop [Tên] chưa sẵn sàng."*

#### 2. Xử lý lỗi tải ảnh (Image Loading Error)
- Sử dụng sự kiện `onError` trong thẻ `img`.
- Nếu URL hợp lệ nhưng không thể tải ảnh (404, lỗi mạng...):
  - Hiển thị thông báo `toast`: *"Không thể tải sơ đồ phòng. Vui lòng kiểm tra kết nối mạng."*
  - Tự động thay thế bằng ảnh dự phòng (Placeholder) để tránh làm vỡ giao diện.

#### 3. Quản lý trạng thái (State Management)
- **Reset State:** Sử dụng `useEffect` để theo dõi trạng thái đóng/mở của Dialog.
- **Logic:** Khi Dialog đóng (`open === false`), trạng thái `showMap` sẽ tự động được reset về `false`.
- **Mục tiêu:** Đảm bảo khi người dùng mở Workshop tiếp theo, sơ đồ không bị tự động hiện lên từ lần xem trước.

## 3. Cấu trúc dữ liệu (Data Structure)

| Trường (API) | Kiểu dữ liệu | Mô tả |
| :--- | :--- | :--- |
| `room_layout_url` | `string (nullable)` | Đường dẫn ảnh trong thư mục `/public/maps/` |

## 4. Bảo trì & Debug (Maintenance)
- Để kiểm tra luồng quá tải: Chỉnh sửa giá trị `currentRequestRate` trong `workshop_service.go` vượt ngưỡng 12,000.
- Log chi tiết quá trình nhận URL được in ra trong Console trình duyệt để hỗ trợ kỹ thuật.

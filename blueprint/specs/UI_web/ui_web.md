# Đặc tả: Giao diện Web Application (ui_web.md)

## 1. Mô tả
Tài liệu quy định cấu trúc màn hình (`Wireframe Text`), luồng tương tác người dùng (`User Flow`) và các ràng buộc về trải nghiệm người dùng (`UX Constraints`) cho hai phân hệ Web chính: Phân hệ Sinh viên (`Student Portal`) và Phân hệ Quản trị (`Admin Dashboard`). Ứng dụng được xây dựng theo mô hình `Single Page Application (SPA)`.

## 2. Phân hệ Sinh viên (Student Portal)
Phân hệ này ưu tiên thiết kế `Mobile-first` (tương thích tốt trên trình duyệt điện thoại) và tích hợp các cơ chế phản hồi thời gian thực khi hệ thống tải cao.

### 2.1. Cấu trúc Màn hình (Screens)

* **Trang chủ (Home / Workshop Catalog):**
    * **Thanh điều hướng (Navbar):** Logo UniHub, Hộp thư thông báo (Biểu tượng chuông), Avatar người dùng.
    * **Khu vực Tìm kiếm & Lọc:** Hỗ trợ lọc theo ngày diễn ra, chủ đề, và loại vé (Miễn phí / Có phí).
    * **Danh sách Workshop (Grid View):** Hiển thị các `Card` chứa thông tin cơ bản: Tên sự kiện, Diễn giả, Thời gian, và **Thanh tiến trình (Progress Bar)** thể hiện tỷ lệ lấp đầy chỗ ngồi (`available_seats` / `capacity`).

* **Trang Chi tiết (Workshop Detail):**
    * **Nội dung AI Summary:** Khối văn bản làm nổi bật 5 gạch đầu dòng tóm tắt mục tiêu Workshop được trích xuất tự động từ luồng `Pipe-and-Filter`.
    * **Thông tin Địa điểm:** Tên phòng và nút xem Sơ đồ phòng (`Room Layout`).
    * **Khu vực Đăng ký (Call to Action):** * Hiển thị giá vé (hoặc chữ "Miễn phí").
        * Nút bấm `[Đăng ký ngay]`. Nút này sẽ chuyển sang trạng thái `Disabled` nếu hệ thống trả về mã `429 Too Many Requests`.

* **Popup Đăng ký & Thanh toán (Checkout Modal):**
    * **Workshop Miễn phí:** Hiển thị xác nhận thông tin và nút `[Xác nhận đăng ký]`. Khi bấm, giao diện hiển thị `Loading Spinner` chờ hệ thống xử lý `Pessimistic Locking`.
    * **Workshop Có phí:** Hiển thị danh sách Phương thức thanh toán (MoMo, VNPay). Áp dụng `Graceful Degradation`: Nếu API Backend báo cổng thanh toán đang bị ngắt mạch (`Circuit Breaker OPEN`), tự động ẩn các nút thanh toán này và hiển thị thông báo bảo trì.

* **Trang Lịch sử & Vé (My Tickets):**
    * Danh sách các sự kiện đã đăng ký thành công (`Status = SUCCESS`).
    * Mỗi vé hiển thị một mã `QR Code` sắc nét để nhân sự quét tại cửa, kèm trạng thái `Check-in` (Chưa tham dự / Đã tham dự).


### Luồng truy cập không cần đăng nhập (Public Access)
Sinh viên có thể thực hiện các hành động sau mà không cần tài khoản:
* **Xem Danh sách Workshop:** Duyệt các sự kiện, lọc theo chủ đề và thời gian.
* **Xem Chi tiết Workshop:** Đọc thông tin diễn giả, phòng tổ chức và **nội dung AI Summary**.
* **Xem trạng thái chỗ ngồi:** Theo dõi thanh tiến trình lấp đầy theo thời gian thực.

### Luồng bắt buộc đăng nhập (Authenticated Access)
Hệ thống sẽ yêu cầu đăng nhập bằng `student_id` để lấy `JWT` chứa `user_id` khi sinh viên thực hiện:
* **Đăng ký tham gia:** Bấm nút `[Đăng ký ngay]`. Hệ thống cần `user_id` để thực hiện khóa bi quan và lưu bản ghi đăng ký.
* **Thanh toán:** Thực hiện thanh toán cho các workshop có phí thông qua cổng thanh toán tích hợp.
* **Vé của tôi:** Xem lịch sử các workshop đã đăng ký thành công và lấy mã `QR Code` để check-in.

## 3. Phân hệ Ban tổ chức (Admin Dashboard)
Phân hệ này ưu tiên thiết kế cho màn hình Desktop, sử dụng các cấu trúc dữ liệu dạng Bảng (`DataGrid`) để quản lý khối lượng lớn thông tin.

### 3.1. Cấu trúc Màn hình (Screens)

* **Trang Đăng nhập (Login):** Form nhập `student_id` (hoặc `username`) và `password`. Hệ thống lưu trữ `JWT` vào bộ nhớ an toàn sau khi xác thực thành công.
* **Bảng điều khiển (Overview Dashboard):**
    * Các thẻ thông kê (Thống kê số lượng Workshop đang mở, tổng số Sinh viên đã đăng ký, tỷ lệ lấp đầy trung bình).
    * Biểu đồ lưu lượng đăng ký theo thời gian thực.
* **Quản lý Workshop (Workshop Management):**
    * Bảng danh sách Workshop hỗ trợ phân trang (`Pagination`), sắp xếp và tìm kiếm.
    * Các hành động tại mỗi dòng: Cập nhật thông tin, Đổi giờ, Đổi phòng, hoặc Hủy sự kiện.
* **Form Tạo Mới / Cập Nhật Workshop:**
    * Các trường nhập liệu chuẩn: Tiêu đề, Diễn giả, Sức chứa, Giá vé, Thời gian.
    * **Khu vực Upload PDF:** Nơi Ban tổ chức tải lên tài liệu giới thiệu. Sau khi tải lên, hệ thống hiển thị trạng thái *"Đang xử lý AI Summary..."* và cập nhật kết quả vào text box bên dưới để Ban tổ chức duyệt trước khi Xuất bản (`Publish`).
* **Quản lý Hệ thống & Đồng bộ (System & Integration):**
    * Bảng theo dõi lịch sử `Batch Import` từ luồng CSV.
    * Khu vực hiển thị `Import Errors`: Danh sách các dòng dữ liệu bị lỗi (sai định dạng email, trùng lặp) từ đêm hôm trước để Ban tổ chức rà soát.

## 4. Xử lý Trải nghiệm Người dùng (UX Constraints)

Nhằm đối phó với tải trọng đột biến và đảm bảo tính nhất quán của dữ liệu, `Frontend` phải bắt buộc tuân thủ các quy tắc sau:

| Tình huống / API Status | Hành vi của Giao diện (UI Behavior) |
| :--- | :--- |
| **HTTP 429 Rate Limit** | Bắt Header `Retry-After`. Thay đổi nội dung nút bấm thành *"Vui lòng đợi {x} giây"*, chuyển nút sang màu xám (`Disabled`) và khóa mọi tương tác click cho đến khi đếm ngược kết thúc. |
| **HTTP 202 Accepted** | API xử lý bất đồng bộ. Giao diện không chuyển trang ngay mà hiển thị màn hình chờ *"Đang xếp hàng giữ chỗ..."*. Khởi chạy cơ chế `Polling` (gọi API 2 giây/lần) để hỏi kết quả từ Backend. |
| **HTTP 409 Conflict** | Xảy ra khi sinh viên vừa bấm đăng ký thì vừa đúng lúc có người khác lấy chỗ cuối cùng. Hiển thị thông báo Toast màu đỏ: *"Rất tiếc, Workshop đã hết chỗ ngay trước mắt bạn."* |
| **Mất kết nối mạng** | Cảnh báo trạng thái Offline. Hủy bỏ các lệnh gọi API ghi dữ liệu (như Đăng ký) nhưng vẫn cho phép sinh viên xem các trang đã được lưu Cache nội bộ (`Service Worker`). |

## 5. Tiêu chí chấp nhận
* Toàn bộ luồng đăng ký của sinh viên có thể thực hiện tối đa trong 3 cú click chuột.
* Giao diện hiển thị phản hồi trạng thái rõ ràng, không để người dùng ở trạng thái không biết hệ thống có đang xử lý yêu cầu hay không (`Silent Failure`).




## 6. Design System & UI/UX Guidelines
Toàn bộ Frontend bắt buộc phải tuân thủ nghiêm ngặt Hệ thống thiết kế (Design System) dưới đây để đảm bảo tính chuyên nghiệp, tránh hiện tượng lạm dụng màu sắc và sai lệch tỷ lệ. Có thể tích hợp thông qua cấu hình `Theme` của Ant Design và `tailwind.config.js`.

### 6.1. Bảng màu (Color Palette)
Tuyệt đối không sử dụng các màu nguyên bản (Pure Colors). Hệ thống áp dụng quy tắc phối màu 60-30-10:

* **Màu chủ đạo (Primary Brand):** Sử dụng các tone màu mang tính học thuật và tin cậy. Gợi ý: `Cobalt Blue` (Xanh xám) hoặc `Indigo` (Chàm). Tránh dùng màu xanh dương chói.
* **Màu nền & Cấu trúc (Neutral/Base):** * Nền trang web (Background): Xám cực nhạt (`#F8FAFC` - Slate 50) để làm nổi bật các Card trắng.
    * Chữ viết (Text): Không dùng màu Đen `#000000`. Phải dùng Xám than (`#1E293B` - Slate 800) cho tiêu đề và Xám nhạt hơn (`#475569` - Slate 600) cho văn bản phụ.
* **Màu ngữ nghĩa (Semantic Colors):** Phải dùng các tone màu Pastel hoặc được giảm độ bão hòa (Desaturated).
    * `Success` (Thành công / Check-in): Xanh ngọc (`Emerald`).
    * `Warning` (Sắp hết chỗ): Vàng cam (`Amber`).
    * `Error/Destructive` (Lỗi / Hủy): Đỏ mận (`Rose`), tuyệt đối không dùng Đỏ tươi.

### 6.2. Nghệ thuật chữ (Typography)
* **Font Family:** Bắt buộc sử dụng các Font chữ Sans-serif hiện đại, chuyên trị hiển thị màn hình như `Inter`, `Roboto` hoặc `SF Pro`. Không dùng nhiều hơn 1 họ Font.
* **Hệ thống phân cấp (Hierarchy):** Áp dụng kích thước chữ theo thang tỷ lệ chuẩn (12px, 14px, 16px, 20px, 24px, 32px). Cột dữ liệu trong bảng Admin mặc định dùng size 14px để tiết kiệm không gian.

### 6.3. Tính nhất quán của Component (Component Consistency)
* **Nút bấm (Buttons):**
    * *Kích thước:* Bắt buộc cố định chiều cao. Size nhỏ (32px), Size chuẩn (40px), Size lớn cho Call-to-action (48px).
    * *Bo góc (Border-radius):* Áp dụng bo góc vừa phải (`rounded-md` tương đương 6px hoặc 8px) cho toàn bộ Input và Button. Tránh bo tròn hoàn toàn (Pill-shape) vì làm giảm tính nghiêm túc của hệ thống giáo dục.
    * *Trạng thái (States):* Mọi nút bấm bắt buộc phải có hiệu ứng chuyển màu nhẹ khi lướt chuột qua (`:hover`), khi nhấn (`:active`) và màu xám mờ khi bị khóa (`:disabled`).
* **Khoảng trắng (Spacing & Grid):**
    * Tuân thủ tuyệt đối **Hệ thống lưới 8-point (8pt Grid System)**. 
    * Mọi khoảng cách (Padding, Margin, Gap) giữa các thành phần bắt buộc phải là bội số của 4 hoặc 8 (VD: 4px, 8px, 16px, 24px, 32px). Tuyệt đối cấm các con số lẻ tẻ như 5px, 11px, 17px.

### 6.4. Trạng thái tải và Phản hồi (Loading & Feedback)
* Tránh sử dụng màn hình trắng hoặc cục xoay (Spinner) giữa màn hình một cách vô nghĩa.
* Khi tải danh sách Workshop hoặc bảng dữ liệu Admin, bắt buộc sử dụng hiệu ứng **Skeleton Loading** (Khối xám nhấp nháy giữ đúng khung layout) để giảm cảm giác chờ đợi của người dùng.
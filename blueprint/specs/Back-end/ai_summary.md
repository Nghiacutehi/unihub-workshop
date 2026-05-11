# Đặc tả: ai_summary.md

## 1. Mô tả
Tài liệu quy định luồng xử lý tự động tóm tắt nội dung từ các tập tin tài liệu Workshop dưới định dạng `PDF`. Hệ thống sử dụng mô hình ngôn ngữ lớn (`LLM`) để trích xuất các thông tin cốt lõi, giúp Sinh viên nhanh chóng nắm bắt mục tiêu và nội dung chính của sự kiện.

## 2. Kiến trúc giải pháp
Module được thiết kế theo kiến trúc `Pipe-and-Filter`. Dữ liệu chảy qua một chuỗi các bộ lọc (`Filters`) độc lập, trong đó đầu ra của bộ lọc này là đầu vào của bộ lọc kế tiếp:

* `Pipe`: Cơ chế truyền dẫn dữ liệu giữa các giai đoạn xử lý.
* `Filter 1 (Extraction)`: Trích xuất văn bản thô từ tập tin `PDF`.
* `Filter 2 (Cleaning)`: Làm sạch dữ liệu, loại bỏ các ký tự đặc biệt, khoảng trắng thừa và `Stopwords`.
* `Filter 3 (Prompt Construction)`: Đóng gói văn bản đã làm sạch vào một cấu trúc `Prompt` chuyên dụng cho tác vụ tóm tắt.
* `Filter 4 (AI Integration)`: Gửi `Payload` tới dịch vụ `AI Model` và nhận kết quả.
* `Sink (Persistence)`: Lưu trữ kết quả tóm tắt cuối cùng vào `PostgreSQL`.

## 3. Workflow
Quy trình thực thi chi tiết luồng xử lý:

1. `Organizer` tải tập tin `PDF` lên thông qua `Admin Dashboard`.
2. Hệ thống lưu trữ tập tin vào `Object Storage` và đẩy một `Job` vào hàng đợi xử lý.
3. `Filter 1` thực thi đọc nội dung `Binary` của `PDF` và chuyển đổi sang dạng `Plain Text`.
4. `Filter 2` thực hiện chuẩn hóa văn bản, giới hạn độ dài ký tự (`Token Limit`) để phù hợp với ngữ cảnh của `AI Model`.
5. `Filter 3` kết hợp văn bản với các chỉ dẫn (`Instruction`) như: "Tóm tắt trong 5 gạch đầu dòng", "Ngôn ngữ: Tiếng Việt".
6. `Filter 4` thực hiện `External API Call` tới `AI Model`. Sử dụng cơ chế `Retry` nếu gặp lỗi mạng.
7. Kết quả trả về được `Sink` kiểm tra định dạng cuối cùng trước khi cập nhật vào trường `summary` của bảng `workshops` trong `Database`.

## 4. Xử lý ngoại lệ

| Lỗi (Exception) | Nguyên nhân | Hành vi hệ thống |
| :--- | :--- | :--- |
| `Unsupported Format` | File tải lên không phải `PDF` hoặc bị mã hóa mật khẩu. | Ngừng xử lý, phản hồi mã lỗi `400 Bad Request`, yêu cầu kiểm tra lại định dạng file. |
| `Content Too Large` | Số lượng `Token` vượt quá giới hạn của `AI Model`. | Thực hiện thuật toán `Truncation` (cắt bớt phần cuối) hoặc chia nhỏ văn bản để tóm tắt từng phần. |
| `AI Service Timeout` | Dịch vụ `AI` (như OpenAI, Gemini) không phản hồi hoặc quá tải. | Kích hoạt `Circuit Breaker`, trả về trạng thái `Summary Pending` và thực hiện `Retry` vào thời điểm khác. |

## 5. Ràng buộc kỹ thuật
* `Reliability`: Sử dụng `Circuit Breaker` cho bộ lọc `AI Integration` để tránh treo tiến trình xử lý khi đối tác gặp sự cố.
* `Performance`: Giai đoạn trích xuất văn bản (`Extraction`) không được chiếm dụng quá 512MB `RAM` để đảm bảo an toàn cho máy chủ.
* `Data Consistency`: Bản tóm tắt chỉ được hiển thị trên giao diện người dùng sau khi trạng thái của `Filter` cuối cùng là `Success`.

## 6. Tiêu chí chấp nhận
* Hệ thống trích xuất đúng nội dung văn bản từ các `PDF` có cấu trúc phức tạp (cột, bảng).
* Bản tóm tắt đảm bảo tính ngắn gọn, súc tích và hiển thị đúng tại trang chi tiết Workshop.
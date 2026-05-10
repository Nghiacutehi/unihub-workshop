# UniHub Workshop - Project Proposal

## Bối cảnh và Vấn đề
Trường Đại học A hiện đang tổ chức "Tuần lễ kỹ năng và nghề nghiệp" thường niên kéo dài trong 5 ngày, với tần suất 8 đến 12 workshop diễn ra song song mỗi ngày. Quy trình hiện tại phụ thuộc hoàn toàn vào việc sử dụng Google Form để ghi nhận đăng ký và thao tác gửi email thông báo thủ công. Khi quy mô sự kiện mở rộng, kiến trúc nghiệp vụ hiện tại bộc lộ rõ điểm yếu chí mạng (Single Point of Failure), dẫn đến tình trạng hệ thống quá tải, quản lý phân tán và hoàn toàn không đáp ứng được nhu cầu thực tế của nhà trường.

## Mục tiêu
Mục tiêu cốt lõi của dự án là xây dựng và triển khai hệ thống UniHub Workshop nhằm số hóa toàn diện quy trình vận hành, từ giai đoạn sinh viên xem lịch, đăng ký cho đến bước check-in thực tế tại sự kiện. Hệ thống được thiết kế với mục tiêu định lượng khắt khe: phải duy trì tính khả dụng (Availability) và xử lý thành công lưu lượng tải đột biến của 12.000 sinh viên truy cập trong 10 phút đầu tiên mở cổng đăng ký, đặc biệt là kiểm soát 60% lưu lượng dồn dập tập trung trong 3 phút đầu.

## Người dùng và Nhu cầu
Hệ thống phục vụ ba nhóm người dùng chính với các phân quyền và nhu cầu chuyên biệt:
* **Sinh viên:** Cần hệ thống để xem danh sách, sơ đồ phòng, đăng ký workshop (miễn phí hoặc có phí) và nhận mã QR để check-in. Nhu cầu nhận thông báo xác nhận tự động qua đa kênh (Email, App và kiến trúc phải hỗ trợ mở rộng sang Telegram trong tương lai).
* **Ban tổ chức:** Cần hệ thống quản trị (Admin Web) để tạo mới, điều chỉnh (đổi phòng, đổi giờ, hủy) workshop và theo dõi dữ liệu thống kê theo thời gian thực. Có nhu cầu tải file PDF để hệ thống tự động trích xuất, làm sạch văn bản và gọi API AI tóm tắt nội dung giới thiệu.
* **Nhân sự check-in:** Cần ứng dụng di động để quét mã QR tại cửa phòng, yêu cầu tính khả dụng cao để đảm bảo không ùn tắc tại cửa sự kiện.

## Phạm vi
* **Trong phạm vi (In-scope):** Phát triển hoàn chỉnh các thành phần bao gồm Web Admin cho Ban tổ chức, Mobile App Check-in cho nhân sự, hệ thống Ứng dụng/Web cho Sinh viên và Backend API lõi. Xây dựng các luồng xử lý ngầm (background jobs) phục vụ việc import file CSV định kỳ và giao tiếp với mô hình AI. Triển khai thực tế các cơ chế kiến trúc chịu lỗi như Rate Limiting, Circuit Breaker và Idempotency.
* **Ngoài phạm vi (Out-of-scope):** Không can thiệp hoặc thay đổi mã nguồn của hệ thống quản lý sinh viên hiện tại của trường (chỉ tiếp nhận dữ liệu một chiều). Không đấu nối với Payment Gateway của ngân hàng thật; các lỗi giao dịch, timeout sẽ được giả lập ở mức độ hệ thống để kiểm thử luồng thiết kế.

## Rủi ro và Ràng buộc
Dự án đối mặt với 5 thách thức kỹ thuật cốt lõi buộc phải được giải quyết triệt để trong giai đoạn thiết kế kiến trúc:
1. **Tranh chấp chỗ ngồi:** Các workshop có giới hạn (ví dụ: 60 chỗ) nhận hàng trăm request cùng lúc. Bắt buộc phải duy trì tính Nhất quán mạnh (Strong Consistency) để không xảy ra tình trạng bán trùng vé.
2. **Tải trọng đột biến:** Cần cơ chế kiểm soát linh hoạt để bảo vệ Backend API khỏi sự cố sập hệ thống do hàng vạn sinh viên truy cập đồng thời, đảm bảo tính công bằng trong việc điều hướng luồng request.
3. **Thanh toán không ổn định:** Cổng thanh toán có rủi ro sập hoặc timeout. Yêu cầu luồng xem thông tin và đăng ký workshop miễn phí không bị ảnh hưởng (Graceful Degradation), đồng thời đảm bảo không trừ tiền hai lần khi client tự động retry.
4. **Check-in offline:** Hạ tầng mạng tại các phòng sự kiện không ổn định. Mobile App bắt buộc phải ghi nhận check-in cục bộ (AP System) và tự động đồng bộ lên máy chủ trung tâm khi có kết nối mạng, đảm bảo không mất mát hay xung đột dữ liệu.
5. **Tích hợp một chiều CSV:** Quá trình import file dữ liệu sinh viên vào ban đêm phải xử lý an toàn các dòng dữ liệu lỗi hoặc trùng lặp, đảm bảo không tạo ra Single Point of Failure làm gián đoạn toàn bộ dịch vụ đang vận hành.
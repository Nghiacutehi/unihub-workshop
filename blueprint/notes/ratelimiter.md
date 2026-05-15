# Đặc tả Kỹ thuật: Rate Limiting & Virtual Waiting Room

Tài liệu mô tả cơ chế kiểm soát lưu lượng và phòng chờ ảo để bảo vệ hệ thống trước tải cao (12,000+ sinh viên).

## 1. Kiến trúc Đa tầng (Multi-layer Defense)

Hệ thống sử dụng 3 lớp bảo vệ để đảm bảo tính ổn định:

1.  **Local Token Bucket:** Giới hạn tải ngay tại bộ nhớ RAM của từng server instance để chặn các cuộc tấn công DDoS cơ bản hoặc lỗi lặp vô tận từ client.
2.  **Distributed Token Bucket (Redis + Lua):** Giới hạn tải dựa trên User ID hoặc IP trên toàn bộ cluster server, sử dụng script Lua để đảm bảo tính nguyên tử.
3.  **Virtual Waiting Room:** Chỉ dành riêng cho luồng Đăng ký Workshop.

## 2. Phòng chờ ảo (Virtual Waiting Room)

Khi số lượng yêu cầu đăng ký vượt quá ngưỡng xử lý song song của Database (mặc định 100), hệ thống sẽ chuyển người dùng vào hàng chờ.

### Luồng xử lý:
1.  **Check Status:** Hệ thống kiểm tra xem người dùng đã có "Token cho phép đăng ký" (Active Token) chưa.
2.  **Enqueue:** Nếu chưa có và hệ thống đang đầy, người dùng được đưa vào Sorted Set trong Redis (`waiting_room:queue:{workshop_id}`).
3.  **Position Calculation:** Trả về vị trí hiện tại dựa trên rank trong Sorted Set.
4.  **Promotion:** Một worker hoặc script tự động sẽ đẩy người dùng từ hàng chờ (`Queue`) sang trạng thái `Granted` khi có chỗ trống.
5.  **TTL:** Token cho phép đăng ký có thời hạn 5 phút. Nếu sinh viên không đăng ký trong thời gian này, chỗ sẽ được nhường cho người tiếp theo.

## 3. Atomic Seat Limiter (Double-Check)

Để tránh "Race Condition" và giảm tải cho PostgreSQL:
- **Redis Check:** Trước khi đẩy vào RabbitMQ, hệ thống thực hiện `DECR` số lượng ghế trong Redis.
- **Atomic Operation:** Sử dụng Lua script để kiểm tra `if seats > 0 then seats -= 1`.
- **Consistency:** Nếu Redis báo hết chỗ, yêu cầu bị từ chối ngay lập tức (HTTP 409). Nếu DB sau đó báo lỗi (sai lệch dữ liệu), hệ thống sẽ thực hiện `INCR` lại Redis để hoàn trả ghế.

## 4. Tham số cấu hình (Default)

- **Capacity:** 100 tokens.
- **Refill Rate:** 10 tokens/second.
- **Waiting Room TTL:** 300s (5 phút).
- **Queue Max Size:** Không giới hạn (dựa trên dung lượng Redis).

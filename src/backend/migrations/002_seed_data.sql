-- Seed data for UniHub Workshop
-- Passwords are bcrypt hashed version of "password123"

-- Users (password: password123)
INSERT INTO users (student_id, password_hash, full_name, email, phone, role) VALUES
('21127001', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Nguyễn Văn An', 'an.nguyen@student.edu.vn', '0901234001', 'STUDENT'),
('21127002', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Trần Thị Bình', 'binh.tran@student.edu.vn', '0901234002', 'STUDENT'),
('21127003', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Lê Hoàng Cường', 'cuong.le@student.edu.vn', '0901234003', 'STUDENT'),
('21127004', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Phạm Minh Duy', 'duy.pham@student.edu.vn', '0901234004', 'STUDENT'),
('21127005', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Hoàng Thị Em', 'em.hoang@student.edu.vn', '0901234005', 'STUDENT'),
('ADMIN001', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Nguyễn Admin', 'admin@unihub.edu.vn', '0909999001', 'ORGANIZER'),
('STAFF001', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Trần Staff', 'staff@unihub.edu.vn', '0909999002', 'STAFF')
ON CONFLICT (student_id) DO NOTHING;

-- Workshops (5-day event, Monday to Friday)
INSERT INTO workshops (title, description, speaker, room, start_time, end_time, capacity, available_seats, price, status) VALUES
-- Day 1 (Monday)
('Kỹ năng CV và Phỏng vấn', 'Workshop hướng dẫn viết CV chuyên nghiệp và kỹ năng trả lời phỏng vấn', 'ThS. Nguyễn Minh Tuấn', 'A101', '2025-06-02 08:00:00+07', '2025-06-02 10:00:00+07', 60, 58, 0.00, 'PUBLISHED'),
('Lập trình Python cho người mới', 'Khóa học cơ bản về Python dành cho sinh viên không chuyên IT', 'TS. Trần Văn Hùng', 'B201', '2025-06-02 08:00:00+07', '2025-06-02 11:00:00+07', 40, 40, 50000.00, 'PUBLISHED'),
('Design Thinking Workshop', 'Phương pháp tư duy thiết kế trong giải quyết vấn đề', 'MBA. Lê Thu Hà', 'C301', '2025-06-02 13:00:00+07', '2025-06-02 16:00:00+07', 50, 50, 0.00, 'PUBLISHED'),
('Khởi nghiệp từ A-Z', 'Chia sẻ kinh nghiệm khởi nghiệp từ các founder thành công', 'CEO Phạm Đức Minh', 'A102', '2025-06-02 14:00:00+07', '2025-06-02 17:00:00+07', 80, 80, 100000.00, 'PUBLISHED'),

-- Day 2 (Tuesday)
('Data Science Fundamentals', 'Giới thiệu về khoa học dữ liệu và ứng dụng thực tế', 'PGS.TS. Ngô Thanh Long', 'B201', '2025-06-03 08:00:00+07', '2025-06-03 11:00:00+07', 45, 45, 75000.00, 'PUBLISHED'),
('Public Speaking Mastery', 'Rèn luyện kỹ năng thuyết trình trước đám đông', 'ThS. Vũ Hoàng Anh', 'A101', '2025-06-03 08:00:00+07', '2025-06-03 10:00:00+07', 60, 60, 0.00, 'PUBLISHED'),
('Digital Marketing 101', 'Marketing số cho người mới bắt đầu', 'Nguyễn Thị Lan - CMO TechCorp', 'C301', '2025-06-03 13:00:00+07', '2025-06-03 16:00:00+07', 55, 55, 50000.00, 'PUBLISHED'),

-- Day 3 (Wednesday)
('Cloud Computing with AWS', 'Thực hành triển khai ứng dụng trên Amazon Web Services', 'AWS Solutions Architect Trần Minh', 'B201', '2025-06-04 08:00:00+07', '2025-06-04 12:00:00+07', 35, 35, 150000.00, 'PUBLISHED'),
('Quản lý tài chính cá nhân', 'Kiến thức cơ bản về đầu tư và quản lý chi tiêu', 'CFA Lê Bảo Ngọc', 'A101', '2025-06-04 13:00:00+07', '2025-06-04 15:00:00+07', 70, 70, 0.00, 'PUBLISHED'),
('UI/UX Design Workshop', 'Thiết kế giao diện và trải nghiệm người dùng với Figma', 'Senior Designer Phạm Anh Khoa', 'C301', '2025-06-04 13:00:00+07', '2025-06-04 17:00:00+07', 40, 40, 80000.00, 'PUBLISHED'),

-- Day 4 (Thursday)
('Machine Learning in Practice', 'Ứng dụng ML vào các bài toán thực tế', 'TS. Hoàng Văn Đức', 'B201', '2025-06-05 08:00:00+07', '2025-06-05 12:00:00+07', 30, 30, 200000.00, 'PUBLISHED'),
('Networking & Personal Branding', 'Xây dựng thương hiệu cá nhân trong thời đại số', 'LinkedIn Expert Mai Phương', 'A101', '2025-06-05 14:00:00+07', '2025-06-05 16:00:00+07', 65, 65, 0.00, 'PUBLISHED'),

-- Day 5 (Friday)
('Hackathon Kick-off', 'Cuộc thi lập trình 24 giờ - Giải thưởng 50 triệu đồng', 'BTC Tuần lễ Kỹ năng', 'Hall A', '2025-06-06 08:00:00+07', '2025-06-07 08:00:00+07', 100, 100, 0.00, 'PUBLISHED'),
('Career Fair & Closing', 'Hội chợ việc làm và lễ bế mạc Tuần lễ Kỹ năng', 'Ban tổ chức', 'Hall B', '2025-06-06 13:00:00+07', '2025-06-06 17:00:00+07', 200, 200, 0.00, 'PUBLISHED');

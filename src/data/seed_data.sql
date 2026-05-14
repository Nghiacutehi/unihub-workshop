-- ==========================================
-- 1. XÓA SẠCH DỮ LIỆU CŨ (RESET)
-- ==========================================
TRUNCATE public.notifications, public.payments, public.registrations, 
         public.workshops, public.users, public.import_errors, public.import_jobs 
CASCADE;

-- ==========================================
-- 2. TẠO NGƯỜI DÙNG (ADMIN, STAFF, STUDENTS)
-- ==========================================
INSERT INTO public.users (user_id, password_hash, full_name, email, role) VALUES 
('admin', '123456', 'Quản trị viên UniHub', 'admin@unihub.edu.vn', 'ADMIN'),
('staff01', '123456', 'Lý Gia Thành (Staff)', 'staff01@unihub.edu.vn', 'STAFF'),
('staff02', '123456', 'Vương Tiểu Nhị (Staff)', 'staff02@unihub.edu.vn', 'STAFF');

INSERT INTO public.users (user_id, password_hash, full_name, email, role)
SELECT 
    '211270' || LPAD(i::text, 2, '0'), '123456', 
    'Sinh viên ' || i, 'sv' || i || '@student.hcmus.edu.vn', 'STUDENT'
FROM generate_series(1, 17) AS i;

-- ==========================================
-- 3. TẠO WORKSHOPS (10 BUỔI)
-- ==========================================
INSERT INTO public.workshops (id, title, speaker, room, start_time, end_time, capacity, available_seats, price, status)
VALUES (
    '22222222-0000-0000-0000-000000000004',
    'Thực chiến Cloud AWS cho người mới',
    'Anh Phạm C', 'Phòng I.53',
    CURRENT_TIMESTAMP - interval '1 hour', 
    CURRENT_TIMESTAMP + interval '2 hours',
    60, 60, 0.00, 'PUBLISHED'
);

INSERT INTO public.workshops (title, speaker, room, start_time, end_time, capacity, available_seats, price, status)
SELECT 
    'Workshop ' || themes.name, themes.speaker, 'Phòng ' || themes.room,
    CURRENT_TIMESTAMP + (i - 4) * interval '1 day',
    CURRENT_TIMESTAMP + (i - 4) * interval '1 day' + interval '2 hours',
    50, 50, themes.price, 'PUBLISHED'
FROM (VALUES 
    (1, 'Kỹ năng Viết CV IT', 'Chị Lê B', 'E.302', 0),
    (2, 'Giải thuật & Leetcode', 'Anh Nguyễn A', 'C.21', 50000),
    (3, 'System Design', 'Anh Trần E', 'I.41', 100000),
    (5, 'UI/UX Trendy', 'Chị Đặng H', 'B.11', 0),
    (6, 'Blockchain Web3', 'Anh Hoàng F', 'D.05', 150000),
    (7, 'Softskills', 'Ms. Emma', 'E.10', 0),
    (8, 'Mobile App', 'Anh Lý G', 'I.53', 50000),
    (9, 'Data Science', 'Dr. Ngô I', 'Hội trường', 0),
    (10, 'Cyber Security', 'Anh Bùi K', 'Lab 01', 80000)
) AS themes(i, name, speaker, room, price);

-- ==========================================
-- 4. TẠO 100 ĐĂNG KÝ, THANH TOÁN & THÔNG BÁO
-- ==========================================
DO $$
DECLARE
    v_user_id UUID;
    v_workshop_id UUID;
    v_reg_id UUID;
    v_price NUMERIC;
    v_i INTEGER;
BEGIN
    FOR v_i IN 1..100 LOOP
        -- Chọn ngẫu nhiên User và Workshop
        SELECT id INTO v_user_id FROM public.users WHERE role = 'STUDENT' ORDER BY random() LIMIT 1;
        SELECT id, price INTO v_workshop_id, v_price FROM public.workshops ORDER BY random() LIMIT 1;
        
        -- 4.1 Chèn Đăng ký
        INSERT INTO public.registrations (user_id, workshop_id, status, ticket_signature)
        VALUES (v_user_id, v_workshop_id, 'SUCCESS', 'SIG_' || md5(v_i::text))
        ON CONFLICT (user_id, workshop_id) DO NOTHING
        RETURNING id INTO v_reg_id;

        IF v_reg_id IS NOT NULL THEN
            -- 4.2 Chèn Thanh toán (Nếu workshop có phí)
            INSERT INTO public.payments (registration_id, transaction_id, amount, provider, status)
            VALUES (v_reg_id, 'TXN_' || md5(v_reg_id::text), v_price, 'VNPAY', 'COMPLETED');

            -- 4.3 Chèn Thông báo mẫu (Thêm ON CONFLICT để tránh lỗi uq_event_channel)
            INSERT INTO public.notifications (user_id, registration_id, channel, title, content, status, event_id)
            VALUES (v_user_id, v_reg_id, 'EMAIL', 'Đăng ký thành công', 'Bạn đã đăng ký thành công Workshop!', 'SENT', v_workshop_id::text)
            ON CONFLICT DO NOTHING;
        END IF;
    END LOOP;
END $$;

-- ==========================================
-- 5. TẠO DỮ LIỆU IMPORT (Lịch sử hệ thống)
-- ==========================================
INSERT INTO public.import_jobs (file_name, status, total_records, success_count, error_count)
VALUES 
('danh_sach_sv_dot_1.csv', 'COMPLETED', 50, 48, 2),
('danh_sach_sv_dot_2.xlsx', 'FAILED', 100, 0, 100);

-- Chèn lỗi mẫu cho các Job trên
INSERT INTO public.import_errors (job_id, row_number, raw_data, error_reason)
SELECT 
    id, 5, '2112700x, Nguyen Van An, an@mail.com', 'Định dạng MSSV không hợp lệ'
FROM public.import_jobs WHERE file_name = 'danh_sach_sv_dot_1.csv';

-- ==========================================
-- 6. CẬP NHẬT GHẾ TRỐNG
-- ==========================================
UPDATE public.workshops w
SET available_seats = capacity - (SELECT count(*) FROM public.registrations r WHERE r.workshop_id = w.id AND r.status = 'SUCCESS');

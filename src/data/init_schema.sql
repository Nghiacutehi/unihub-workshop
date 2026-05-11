-- Kích hoạt extension để tự động sinh UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================
-- BẢNG 1: USERS (Người dùng hệ thống)
-- ==========================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id VARCHAR(20) UNIQUE NOT NULL, 
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'STUDENT', 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index tăng tốc truy vấn đăng nhập và import CSV
CREATE INDEX idx_users_student_id ON users(student_id);

-- ==========================================
-- BẢNG 2: WORKSHOPS (Thông tin sự kiện)
-- ==========================================
CREATE TABLE workshops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    speaker TEXT NOT NULL, 
    room VARCHAR(50) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    capacity INT NOT NULL CHECK (capacity > 0),
    available_seats INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00, 
    summary TEXT, 
    status VARCHAR(20) NOT NULL DEFAULT 'PUBLISHED',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Chống cháy: Không bao giờ cho phép ghế rớt xuống số âm
    CONSTRAINT chk_positive_seats CHECK (available_seats >= 0 AND available_seats <= capacity)
);

-- ==========================================
-- BẢNG 3: REGISTRATIONS (Vé đăng ký)
-- ==========================================
CREATE TABLE registrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workshop_id UUID NOT NULL REFERENCES workshops(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING_PAYMENT', 
    qr_code TEXT, 
    is_checked_in BOOLEAN NOT NULL DEFAULT FALSE,
    scanned_at TIMESTAMP WITH TIME ZONE, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Chống spam: 1 user chỉ đăng ký 1 workshop 1 lần
    CONSTRAINT uq_user_workshop UNIQUE (user_id, workshop_id)
);

-- Index tối ưu truy xuất lịch sử đăng ký
CREATE INDEX idx_registrations_user ON registrations(user_id);
CREATE INDEX idx_registrations_workshop ON registrations(workshop_id);

-- ==========================================
-- BẢNG 4: PAYMENTS (Giao dịch thanh toán)
-- ==========================================
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    registration_id UUID NOT NULL REFERENCES registrations(id) ON DELETE CASCADE,
    transaction_id VARCHAR(100) UNIQUE NOT NULL, 
    amount DECIMAL(10, 2) NOT NULL, 
    provider VARCHAR(50) NOT NULL, 
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- BẢNG 5: NOTIFICATIONS (Lịch sử thông báo)
-- ==========================================
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    registration_id UUID REFERENCES registrations(id) ON DELETE SET NULL,
    channel VARCHAR(20) NOT NULL, 
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', 
    event_id VARCHAR(100) NOT NULL, 
    error_message TEXT, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP WITH TIME ZONE,
    
    -- Chống gửi trùng lặp thông báo
    CONSTRAINT uq_event_channel UNIQUE (event_id, channel)
);

-- Index tối ưu truy vấn cho icon cái chuông
CREATE INDEX idx_notifications_user ON notifications(user_id);

-- ==========================================
-- BẢNG 6: IMPORT_JOBS (Lịch sử các đợt chạy Batch Import)
-- ==========================================
CREATE TABLE import_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PROCESSING', 
    total_records INT DEFAULT 0,
    success_count INT DEFAULT 0,
    error_count INT DEFAULT 0,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- ==========================================
-- BẢNG 7: IMPORT_ERRORS (Chi tiết các dòng CSV bị lỗi)
-- ==========================================
CREATE TABLE import_errors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id UUID NOT NULL REFERENCES import_jobs(id) ON DELETE CASCADE,
    row_number INT NOT NULL, 
    raw_data TEXT, 
    error_reason TEXT NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
# Admin Portal Architecture

Tài liệu hướng dẫn cấu trúc code và các thành phần dùng chung trong trang quản trị Admin.

## 1. Thành phần dùng chung (Shared Components)

Để đảm bảo tính đồng nhất (DRY - Don't Repeat Yourself), các trang Admin phải sử dụng:

### PageHeader (`src/web/components/admin/page-header.tsx`)
Mọi trang (Dashboard, Workshops, Logs) đều phải bắt đầu bằng component này để hiển thị Tiêu đề và Nút hành động chính.
```tsx
<PageHeader 
  title="Tiêu đề trang"
  description="Mô tả ngắn gọn"
  actionLabel="Nút bấm (nếu có)"
  onAction={() => handleAction()}
/>
```

### WorkshopTable (`src/web/components/admin/workshop-table.tsx`)
Bảng hiển thị danh sách Workshop có tích hợp sẵn phân trang và các trạng thái màu sắc.
- Prop `showActions={false}`: Dùng cho Dashboard (ẩn nút Sửa/Xóa).
- Prop `showActions={true}`: Dùng cho trang Quản lý chuyên sâu.

## 2. Quản lý dữ liệu (Data Management)

### Custom Hooks
Sử dụng Custom Hooks để tách biệt logic gọi API khỏi UI Component.
- **`useWorkshops`**: Tự động fetch dữ liệu, quản lý trạng thái `loading` và `error`.

## 3. Phân quyền (Role-based Access)
- Trang Admin nằm trong nhóm Route `/admin/*`.
- Chỉ người dùng có `role = 'ADMIN'` hoặc `'STAFF'` mới có quyền truy cập thông qua Middleware của Supabase.

# UniHub Design System (Midnight Indigo)

Tài liệu này quy định về ngôn ngữ thiết kế đồng bộ giữa Web (Admin) và Mobile (Staff App).

## 1. Bảng màu chủ đạo (Color Palette)

Hệ thống sử dụng tông màu **Midnight Indigo** làm chủ đạo để tạo cảm giác chuyên nghiệp, tin cậy.

| Thành phần | Mã màu (HEX) | Sử dụng |
| :--- | :--- | :--- |
| **Primary** | `#3730A3` | Nút bấm chính, Sidebar active, Icon chính. |
| **Primary Hover** | `#312E81` | Trạng thái hover của nút bấm chính. |
| **Secondary (Light)** | `#EEF2FF` | Background của Tag, Nút bấm phụ, Hover nền. |
| **Background** | `#F8FAFC` | Nền chính của toàn bộ ứng dụng. |
| **Text Primary** | `#0F172A` | Tiêu đề, văn bản quan trọng. |
| **Text Muted** | `#64748B` | Mô tả phụ, chú thích. |

## 2. Quy tắc giao diện (UI Rules)

### Web (Next.js + Tailwind)
- **Shadow:** Sử dụng đổ bóng có màu để tạo chiều sâu: `shadow-lg shadow-indigo-500/20`.
- **Bo góc (Radius):** Ưu tiên `rounded-2xl` cho các Card và `rounded-xl` cho Nút bấm.
- **Hover:** Loại bỏ hoàn toàn màu tím (`purple`), thay thế bằng `indigo-50` hoặc `primary/90`.

### Mobile (React Native)
- Toàn bộ màu sắc phải được lấy từ file `src/mobile/src/constants/theme.ts`. Không fix cứng mã hex trong các component.
- Trạng thái Workshop:
    - **In Progress:** Border Indigo, Nền xanh nhạt.
    - **Completed:** Nền đỏ nhạt (`rose-50`), Chữ đỏ.

## 3. Typography
- Font chữ: **Inter** hoặc **System Sans-serif**.
- Cỡ chữ tiêu đề: `28px - 32px` (Bold).
- Cỡ chữ nội dung: `14px - 16px`.

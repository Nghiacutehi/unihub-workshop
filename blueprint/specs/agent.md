# Đặc tả: Quy tắc viết Code cho AI Agent (agent.md)

## 1. Mục đích
Tài liệu này là **bộ luật tối thượng** (Source of Truth) quy định cách AI Agent phải tổ chức mã nguồn, đặt tên, phân chia thư mục, viết logic và xử lý lỗi khi tham gia phát triển dự án **UniHub Workshop**. Mọi dòng code được sinh ra bởi AI phải tuân thủ nghiêm ngặt các quy tắc bên dưới. Nếu có xung đột giữa tài liệu này và yêu cầu tức thời của người dùng, **ưu tiên theo yêu cầu người dùng** nhưng phải cảnh báo về sự vi phạm.

---

## 2. Nguyên tắc kiến trúc cốt lõi (Architectural Principles)

### 2.1. SOLID

| Nguyên tắc | Ý nghĩa | Quy tắc áp dụng |
| :--- | :--- | :--- |
| **S** — Single Responsibility | Mỗi file / module / component chỉ chịu trách nhiệm cho **một việc duy nhất**. | Một component không được vừa fetch data, vừa render UI, vừa xử lý form validation. Tách thành `useXxxQuery` (hook), `XxxForm` (UI), `xxxSchema` (validation). |
| **O** — Open/Closed | Mở cho việc mở rộng, đóng cho việc sửa đổi. | Khi thêm phương thức thanh toán mới (ví dụ: ZaloPay), **không** được sửa code cũ của MoMo/VNPay. Thay vào đó, tạo một module mới triển khai cùng interface `PaymentGateway`. |
| **L** — Liskov Substitution | Các lớp con phải thay thế được lớp cha mà không làm hỏng logic. | Mọi adapter thanh toán (`MomoAdapter`, `VnpayAdapter`) phải triển khai đầy đủ contract của interface `PaymentGateway` mà không thêm tiền điều kiện hoặc giảm hậu điều kiện. |
| **I** — Interface Segregation | Không ép buộc một module phụ thuộc vào các hàm nó không sử dụng. | Component `WorkshopCard` chỉ cần nhận `{ title, speaker, date, seatRatio }`, **không** được truyền cả object `Workshop` nặng 20 trường vào. Định nghĩa `type WorkshopCardProps` riêng. |
| **D** — Dependency Inversion | Module cấp cao không phụ thuộc vào module cấp thấp; cả hai phụ thuộc vào abstraction. | Logic đăng ký (`BookingService`) không được import trực tiếp `PostgresClient`. Thay vào đó, nhận một `DatabasePort` (interface/type) qua constructor hoặc dependency injection. |

### 2.2. DRY — Don't Repeat Yourself
* **Cấm** copy-paste code giữa `Student-Front-end` và `Admin-Front-end`. Nếu hai phân hệ dùng chung logic (ví dụ: format ngày giờ, xử lý lỗi API, type definitions), logic đó **phải** được đặt trong thư mục `shared/` hoặc `packages/` dùng chung.
* **Cấm** viết lại hàm tiện ích (utility) đã tồn tại. Trước khi tạo một hàm mới, AI **bắt buộc** phải tìm kiếm trong codebase xem hàm tương tự đã có chưa.
* **Ngoại lệ hợp lệ:** Cho phép lặp lại nếu việc trừu tượng hóa tạo ra coupling không cần thiết (ví dụ: hai component tình cờ giống nhau nhưng phục vụ domain hoàn toàn khác nhau).

### 2.3. KISS — Keep It Simple, Stupid
* Ưu tiên giải pháp đơn giản nhất mà vẫn đúng yêu cầu. **Cấm** over-engineering.
* Không tạo abstraction layer nếu chỉ có **đúng 1 implementation**. Ví dụ: Không tạo `ILogger` → `ConsoleLogger` nếu chỉ dùng `console.log`.
* Không sử dụng Design Pattern (Factory, Strategy, Observer...) nếu bài toán không đủ phức tạp để biện minh cho nó. Rule of Three: Chỉ áp dụng pattern khi logic lặp lại **ít nhất 3 lần**.
* Code phải đọc được như văn bản tự nhiên. Nếu cần comment dài dòng để giải thích một đoạn code, hãy viết lại đoạn code đó cho rõ ràng hơn thay vì comment.

### 2.4. YAGNI — You Ain't Gonna Need It
* **Tuyệt đối cấm** viết code "phòng xa" cho tính năng chưa có trong đặc tả (`specs/`). Không tạo sẵn endpoint, table, component, hay service cho tính năng tương lai.
* Không thêm cấu hình, tham số, hoặc abstraction layer "cho linh hoạt" nếu hiện tại chỉ có 1 trường hợp sử dụng.
* **Ví dụ vi phạm:** Tạo sẵn module `EmailNotification`, `SMSNotification`, `TelegramNotification` trong khi spec `notification.md` chỉ yêu cầu `In-App Notification`.

### 2.5. Loose Coupling & High Cohesion
* **Loose Coupling (Liên kết lỏng):** Các module giao tiếp với nhau thông qua **interface / type / event**, không import trực tiếp implementation cụ thể. Khi thay đổi nội bộ một module, các module khác không bị ảnh hưởng.
* **High Cohesion (Gắn kết cao):** Mọi file, hàm, biến bên trong một module phải phục vụ cùng một mục đích nghiệp vụ. Không nhồi nhét logic không liên quan vào cùng một module.
* **Ví dụ đúng:** Thư mục `features/booking/` chứa tất cả code liên quan đến đăng ký (hook, component, type, service). Thư mục `features/payment/` chứa tất cả code liên quan đến thanh toán. Hai thư mục giao tiếp qua shared types, không import lẫn nhau.

---

## 3. Quy tắc tổ chức thư mục (Directory Structure)

### 3.1. Cấu trúc tổng thể dự án

```text
/unihub-workshop (Root)
│
├── /blueprint                          # Tài liệu thiết kế (Specs, ERD, Architecture)
│   └── /specs
│       ├── /Back-end                   # Đặc tả các module Backend
│       ├── /UI_web                     # Đặc tả giao diện Web
│       ├── /UI_mobile                  # Đặc tả giao diện Mobile
│       └── agent.md                    # File này
│
├── /src                                # Toàn bộ mã nguồn
│   ├── /shared                         # ★ Code dùng chung giữa các app
│   │   ├── /types                      # TypeScript interfaces & DTOs
│   │   ├── /constants                  # Hằng số, enum, config chung
│   │   ├── /utils                      # Hàm tiện ích (formatDate, cn, v.v.)
│   │   └── /validators                 # Zod schemas dùng chung
│   │
│   ├── /web
│   │   ├── /student-portal             # Next.js App — Phân hệ Sinh viên
│   │   └── /admin-dashboard            # Next.js App — Phân hệ Quản trị
│   │
│   ├── /mobile
│   │   └── /staff-app                  # React Native / Flutter — App Check-in
│   │
│   ├── /backend                        # API Server (NestJS / Go / v.v.)
│   │
│   └── /data                           # Seed data, SQL scripts, CSV mẫu
│
├── .gitignore
└── README.md
```

### 3.2. Cấu trúc bên trong mỗi Web App (Next.js)

```text
/student-portal (hoặc /admin-dashboard)
│
├── /app                                # Next.js App Router
│   ├── layout.tsx                      # Root layout (font, theme, metadata)
│   ├── page.tsx                        # Trang chủ
│   ├── /workshops
│   │   ├── page.tsx                    # Danh sách workshop
│   │   └── /[id]
│   │       └── page.tsx                # Chi tiết workshop
│   └── /my-tickets
│       └── page.tsx                    # Vé của tôi (Authenticated)
│
├── /components
│   ├── /ui                             # Atomic components (Button, Card, Input...)
│   │   └── button.tsx                  # Shadcn/UI primitives — KHÔNG sửa đổi
│   ├── /layout                         # Components bố cục (Navbar, Footer, Sidebar)
│   │   └── navbar.tsx
│   └── /features                       # Components nghiệp vụ theo tính năng
│       ├── /workshop                   # Components liên quan đến Workshop
│       │   ├── workshop-card.tsx
│       │   ├── workshop-grid.tsx
│       │   └── workshop-filters.tsx
│       ├── /booking                    # Components liên quan đến Đăng ký
│       │   ├── checkout-modal.tsx
│       │   └── booking-status.tsx
│       └── /ticket                     # Components liên quan đến Vé
│           └── ticket-card.tsx
│
├── /hooks                              # Custom React hooks
│   ├── use-workshops.ts                # Hook fetch danh sách workshop
│   ├── use-booking.ts                  # Hook xử lý đăng ký
│   └── use-auth.ts                     # Hook xác thực người dùng
│
├── /lib                                # Thư viện nội bộ
│   ├── api-client.ts                   # Cấu hình Axios/Fetch wrapper
│   ├── utils.ts                        # Re-export từ shared/utils
│   └── constants.ts                    # Re-export từ shared/constants
│
├── /styles                             # CSS chuyên biệt cho app này
│   └── globals.css
│
└── /public                             # Static assets (images, icons)
```

### 3.3. Quy tắc đặt tên thư mục & file

| Đối tượng | Quy tắc | Ví dụ đúng | Ví dụ sai |
| :--- | :--- | :--- | :--- |
| **Thư mục** | `kebab-case`, số ít, viết thường | `workshop-card/`, `my-tickets/` | `WorkshopCard/`, `My_Tickets/` |
| **File component (.tsx)** | `kebab-case` | `workshop-card.tsx` | `WorkshopCard.tsx`, `workshopCard.tsx` |
| **File hook (.ts)** | Bắt đầu bằng `use-`, `kebab-case` | `use-workshops.ts` | `useWorkshops.ts`, `workshops-hook.ts` |
| **File type/interface (.ts)** | `kebab-case` | `workshop.types.ts` | `IWorkshop.ts`, `WorkshopTypes.ts` |
| **File utility (.ts)** | `kebab-case`, mô tả chức năng | `format-date.ts`, `cn.ts` | `helpers.ts`, `misc.ts` |
| **Biến & Hàm (JS/TS)** | `camelCase` | `availableSeats`, `formatDate()` | `available_seats`, `format_date()` |
| **Component React** | `PascalCase` (chỉ trong code, không phải tên file) | `export function WorkshopCard()` | `export function workshopCard()` |
| **Hằng số** | `SCREAMING_SNAKE_CASE` | `MAX_RETRY_COUNT`, `API_BASE_URL` | `maxRetryCount`, `apiBaseUrl` |
| **Type / Interface** | `PascalCase`, không prefix `I` | `type Workshop = {...}` | `interface IWorkshop`, `type TWorkshop` |

---

## 4. Quy tắc viết Code (Coding Standards)

### 4.1. TypeScript
* **Strict Mode:** Bắt buộc bật `"strict": true` trong `tsconfig.json`. Không được dùng `any`. Nếu type chưa rõ, dùng `unknown` và type-guard.
* **Type vs Interface:** Dùng `type` cho Union/Intersection và data shapes. Dùng `interface` chỉ khi cần `extends` hoặc declaration merging. Ưu tiên `type` cho nhất quán.
* **Enum:** Tránh dùng TypeScript `enum`. Thay thế bằng `as const` objects:
  ```typescript
  // ✅ Đúng
  export const BookingStatus = {
    PENDING: "PENDING",
    SUCCESS: "SUCCESS",
    FAILED: "FAILED",
  } as const;
  type BookingStatus = (typeof BookingStatus)[keyof typeof BookingStatus];

  // ❌ Sai
  enum BookingStatus { PENDING, SUCCESS, FAILED }
  ```

### 4.2. React & Next.js
* **Server Components mặc định:** Mọi component là Server Component trừ khi cần interactivity (state, event handler, browser API). Chỉ thêm `"use client"` khi thực sự cần thiết.
* **Không dùng `useEffect` để fetch data:** Sử dụng Server Components hoặc React Server Actions. Nếu cần client-side fetching, dùng thư viện chuyên dụng (SWR, TanStack Query) thay vì `useEffect` + `useState` thủ công.
* **Component Size:** Một file component không được vượt quá **150 dòng code** (không tính import và type). Nếu vượt quá, bắt buộc tách thành các sub-component.
* **Props Drilling:** Nếu prop phải truyền qua **hơn 2 cấp** component, phải sử dụng React Context hoặc State Management. Không được truyền tay qua 3-4 tầng.

### 4.3. Styling (CSS / Tailwind)
* Tuân thủ hoàn toàn **Design System** đã quy định trong `ui_web.md` mục 6 và `ui_mobile.md` mục 2.
* **Tailwind class order:** Tuân theo thứ tự: Layout → Spacing → Sizing → Typography → Colors → Effects → States.
  ```tsx
  // ✅ Đúng thứ tự
  <div className="flex items-center gap-4 p-4 w-full text-sm text-slate-800 bg-white rounded-md shadow-sm hover:shadow-md">

  // ❌ Lộn xộn
  <div className="shadow-sm text-sm p-4 hover:shadow-md rounded-md flex bg-white w-full gap-4 text-slate-800 items-center">
  ```
* **Cấm Magic Number trong CSS:** Mọi giá trị spacing phải tuân thủ hệ thống 8-point (bội số của 4 hoặc 8). Cấm các giá trị lẻ tẻ như `5px`, `13px`, `17px`.

### 4.4. Xử lý lỗi (Error Handling)
* **Không bao giờ nuốt lỗi (Swallow Errors):**
  ```typescript
  // ❌ Cấm tuyệt đối
  try { await fetchData(); } catch (e) { /* ignore */ }

  // ✅ Bắt buộc xử lý
  try {
    await fetchData();
  } catch (error) {
    logger.error("Failed to fetch workshops", { error, workshopId });
    throw new AppError("FETCH_FAILED", "Không thể tải dữ liệu workshop");
  }
  ```
* **Xử lý HTTP Status:** Tuân thủ bảng UX Constraints trong `ui_web.md` mục 4 (HTTP 429, 202, 409, Offline).
* **Error Boundary:** Mỗi feature section phải được bọc bởi React Error Boundary. Lỗi ở component Workshop không được crash toàn bộ trang.

### 4.5. Bảo mật (Security)
* **Không hardcode secrets:** Mọi API key, JWT secret, database URL phải nằm trong file `.env` và được khai báo trong `.env.example` (không chứa giá trị thật).
* **Sanitize Input:** Mọi dữ liệu từ người dùng phải được validate bằng Zod schema trước khi gửi lên API.
* **XSS Prevention:** Không dùng `dangerouslySetInnerHTML`. Nếu bắt buộc phải render HTML (ví dụ: AI Summary), sử dụng thư viện sanitize chuyên dụng.

---

## 5. Quy tắc Git & Version Control

### 5.1. Commit Message
Tuân thủ **Conventional Commits**:
```
<type>(<scope>): <mô tả ngắn gọn>

Ví dụ:
feat(booking): thêm luồng đăng ký workshop miễn phí
fix(auth): sửa lỗi token hết hạn không redirect về login
refactor(workshop): tách WorkshopCard thành sub-components
docs(specs): cập nhật đặc tả payment module
chore(deps): nâng cấp next.js lên 16.3
style(ui): chỉnh spacing theo 8-point grid
```

| Type | Mô tả |
| :--- | :--- |
| `feat` | Tính năng mới |
| `fix` | Sửa lỗi |
| `refactor` | Tái cấu trúc code (không thay đổi hành vi) |
| `docs` | Thay đổi tài liệu |
| `style` | Thay đổi giao diện / format code |
| `chore` | Cập nhật build, dependencies, CI/CD |
| `test` | Thêm hoặc sửa test |

### 5.2. Branching
* **Không commit trực tiếp lên `main`** (trừ hotfix khẩn cấp).
* Tạo nhánh theo format: `<type>/<mô-tả-ngắn>`. Ví dụ: `feat/checkout-modal`, `fix/qr-scanner-crash`.

### 5.3. Gitignore
* **Tuyệt đối không push** các thư mục/file sau lên Git:
  - `node_modules/`, `.next/`, `dist/`, `build/`
  - `.env`, `.env.local` (chỉ push `.env.example`)
  - File tạm của IDE: `.vscode/settings.json`, `.idea/`
  - File tạm của OS: `.DS_Store`, `Thumbs.db`, `~$*`, `*.tmp`
  - Lock file không phù hợp: Nếu dự án dùng `pnpm`, chỉ giữ `pnpm-lock.yaml`, xóa `package-lock.json` (và ngược lại).

---

## 6. Quy tắc tối ưu hiệu năng (Performance)

* **Lazy Loading:** Các trang con (route) phải được tải lười (dynamic import) để giảm bundle size ban đầu.
* **Image Optimization:** Sử dụng `next/image` thay cho thẻ `<img>` thông thường. Bắt buộc khai báo `width`, `height` và `alt`.
* **Memoization:** Chỉ dùng `React.memo`, `useMemo`, `useCallback` khi có bằng chứng rõ ràng về vấn đề hiệu năng (đo bằng React DevTools Profiler). Không memo phòng xa — vi phạm YAGNI.
* **Bundle Analysis:** Trước mỗi lần release, kiểm tra bundle size. Một component UI đơn lẻ không được import thư viện nặng hơn 50KB gzipped.

---

## 7. Quy tắc Documentation & Comments

* **Self-documenting Code:** Code phải tự giải thích qua cách đặt tên biến và hàm. Comment chỉ dùng để giải thích **tại sao** (why), không giải thích **cái gì** (what).
  ```typescript
  // ❌ Comment vô nghĩa
  // Kiểm tra xem available_seats có lớn hơn 0 không
  if (availableSeats > 0) { ... }

  // ✅ Comment giải thích lý do
  // Sử dụng Pessimistic Lock để tránh Overbooking khi 200+ sinh viên
  // tranh chấp 60 chỗ trong 3 phút đầu (xem booking.md mục 3)
  await db.query("SELECT ... FOR UPDATE");
  ```
* **JSDoc cho Public API:** Mọi hàm/type được export ra ngoài module phải có JSDoc mô tả tham số, giá trị trả về và ví dụ sử dụng.
* **TODO Format:** `// TODO(tên-người): Mô tả việc cần làm — Deadline: YYYY-MM-DD`

---

## 8. Checklist trước khi AI sinh code

Trước khi output bất kỳ đoạn code nào, AI Agent **bắt buộc** phải tự kiểm tra:

- [ ] Code có vi phạm Single Responsibility không? (Mỗi file/hàm chỉ làm 1 việc)
- [ ] Code có lặp lại logic đã tồn tại trong codebase không? (DRY)
- [ ] Giải pháp có đơn giản nhất có thể không? (KISS)
- [ ] Có đang viết code cho tính năng chưa có trong specs không? (YAGNI)
- [ ] File mới có được đặt đúng thư mục theo cấu trúc mục 3 không?
- [ ] Tên file, biến, hàm có tuân thủ Naming Convention mục 3.3 không?
- [ ] Component có vượt quá 150 dòng không?
- [ ] Có import trực tiếp implementation thay vì abstraction không? (Loose Coupling)
- [ ] Styling có tuân thủ Design System (8-point grid, bảng màu chuẩn) không?
- [ ] Lỗi có được xử lý đầy đủ, không nuốt lỗi không?
- [ ] Có hardcode secret hay magic number nào không?
- [ ] Commit message có đúng format Conventional Commits không?

---

## 9. Tiêu chí chấp nhận
* Mọi code được AI sinh ra phải **biên dịch thành công** (`tsc --noEmit` không lỗi).
* Mọi component phải render được mà không crash (không có runtime error khi mount).
* Cấu trúc thư mục phải **khớp chính xác** với sơ đồ trong mục 3 của tài liệu này.
* Code phải đọc hiểu được bởi một lập trình viên trung bình trong vòng **30 giây** cho mỗi hàm.

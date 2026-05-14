"use client"

import { WorkshopCard, type Workshop } from "@/components/student/workshop-card"

// Mock data for workshops
const mockWorkshops: Workshop[] = [
  {
    id: "1",
    title: "Giới thiệu về Trí tuệ Nhân tạo và Ứng dụng trong Thực tế",
    speaker: "TS. Nguyễn Văn Minh",
    speakerTitle: "Giảng viên CNTT",
    date: "15/05/2026",
    time: "09:00 - 11:30",
    location: "Hội trường A - Tòa nhà B2",
    capacity: 150,
    availableSeats: 23,
    ticketType: "free",
    category: "Công nghệ",
  },
  {
    id: "2",
    title: "UI/UX Design: Từ Ý tưởng đến Sản phẩm Hoàn chỉnh",
    speaker: "Trần Thị Hương",
    speakerTitle: "Senior Designer tại FPT",
    date: "17/05/2026",
    time: "14:00 - 17:00",
    location: "Phòng Lab D401",
    capacity: 50,
    availableSeats: 8,
    ticketType: "paid",
    price: 50000,
    category: "Thiết kế",
  },
  {
    id: "3",
    title: "Khởi nghiệp Sinh viên: Cơ hội và Thách thức",
    speaker: "Lê Quang Đức",
    speakerTitle: "CEO Startup XYZ",
    date: "20/05/2026",
    time: "08:30 - 11:00",
    location: "Hội trường C - Tòa nhà A1",
    capacity: 200,
    availableSeats: 0,
    ticketType: "free",
    category: "Kinh doanh",
  },
  {
    id: "4",
    title: "Digital Marketing: Chiến lược Thu hút Khách hàng Gen Z",
    speaker: "Phạm Minh Tú",
    speakerTitle: "Marketing Manager",
    date: "22/05/2026",
    time: "13:30 - 16:00",
    location: "Phòng Seminar E201",
    capacity: 80,
    availableSeats: 45,
    ticketType: "paid",
    price: 75000,
    category: "Marketing",
  },
  {
    id: "5",
    title: "Kỹ năng Thuyết trình Chuyên nghiệp",
    speaker: "ThS. Hoàng Thu Hà",
    speakerTitle: "Chuyên gia đào tạo",
    date: "25/05/2026",
    time: "09:00 - 12:00",
    location: "Phòng họp F102",
    capacity: 40,
    availableSeats: 12,
    ticketType: "free",
    category: "Kỹ năng mềm",
  },
  {
    id: "6",
    title: "React & Next.js: Xây dựng Web App Hiện đại",
    speaker: "Võ Thanh Sơn",
    speakerTitle: "Fullstack Developer",
    date: "28/05/2026",
    time: "14:00 - 17:30",
    location: "Phòng Lab D402",
    capacity: 35,
    availableSeats: 35,
    ticketType: "paid",
    price: 100000,
    category: "Công nghệ",
  },
]

export function WorkshopGrid() {
  return (
    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {mockWorkshops.map((workshop) => (
        <WorkshopCard key={workshop.id} workshop={workshop} />
      ))}
    </div>
  )
}

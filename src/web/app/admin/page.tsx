"use client"

import { Plus, Search, Filter } from "lucide-react"
import { AdminLayout } from "@/components/admin/admin-layout"
import { WorkshopTable, Workshop } from "@/components/admin/workshop-table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

// Mock data for workshops
const mockWorkshops: Workshop[] = [
  {
    id: "1",
    title: "AI trong Giáo dục: Xu hướng và Ứng dụng",
    speaker: "TS. Nguyễn Văn An",
    datetime: "2024-03-15T09:00:00",
    capacity: 100,
    registered: 95,
    price: 0,
    status: "open",
  },
  {
    id: "2",
    title: "Design Thinking cho Sinh viên",
    speaker: "ThS. Trần Thị Bình",
    datetime: "2024-03-18T14:00:00",
    capacity: 50,
    registered: 48,
    price: 50000,
    status: "open",
  },
  {
    id: "3",
    title: "Khởi nghiệp từ Ý tưởng đến Thực tế",
    speaker: "CEO Lê Minh Cường",
    datetime: "2024-03-10T08:30:00",
    capacity: 200,
    registered: 200,
    price: 100000,
    status: "completed",
  },
  {
    id: "4",
    title: "Soft Skills trong Môi trường Doanh nghiệp",
    speaker: "HR Director Phạm Thị Dung",
    datetime: "2024-03-20T10:00:00",
    capacity: 80,
    registered: 45,
    price: 0,
    status: "open",
  },
  {
    id: "5",
    title: "Python cho Phân tích Dữ liệu",
    speaker: "TS. Hoàng Văn Em",
    datetime: "2024-03-12T13:30:00",
    capacity: 60,
    registered: 60,
    price: 75000,
    status: "cancelled",
  },
  {
    id: "6",
    title: "Kỹ năng Thuyết trình Chuyên nghiệp",
    speaker: "MC Ngô Thanh Phong",
    datetime: "2024-03-22T09:00:00",
    capacity: 40,
    registered: 28,
    price: 0,
    status: "open",
  },
  {
    id: "7",
    title: "Cloud Computing Fundamentals",
    speaker: "AWS Expert Vũ Đức Giang",
    datetime: "2024-03-25T14:00:00",
    capacity: 70,
    registered: 55,
    price: 150000,
    status: "open",
  },
  {
    id: "8",
    title: "Marketing số trong Kỷ nguyên 4.0",
    speaker: "CMO Đặng Thị Hoa",
    datetime: "2024-03-08T10:00:00",
    capacity: 90,
    registered: 90,
    price: 80000,
    status: "completed",
  },
  {
    id: "9",
    title: "Quản lý Tài chính Cá nhân",
    speaker: "Financial Advisor Bùi Văn Ích",
    datetime: "2024-03-28T08:30:00",
    capacity: 120,
    registered: 67,
    price: 0,
    status: "open",
  },
  {
    id: "10",
    title: "UX/UI Design Workshop",
    speaker: "Design Lead Cao Thị Kim",
    datetime: "2024-03-30T13:00:00",
    capacity: 35,
    registered: 32,
    price: 120000,
    status: "open",
  },
  {
    id: "11",
    title: "Blockchain và Ứng dụng Thực tiễn",
    speaker: "CTO Đinh Văn Long",
    datetime: "2024-04-02T09:00:00",
    capacity: 55,
    registered: 20,
    price: 200000,
    status: "open",
  },
  {
    id: "12",
    title: "Tiếng Anh Giao tiếp Doanh nghiệp",
    speaker: "IELTS Trainer Emma Wilson",
    datetime: "2024-04-05T14:30:00",
    capacity: 30,
    registered: 28,
    price: 0,
    status: "open",
  },
]

import { useWorkshops } from "@/hooks/use-workshops"
import { Loader2, AlertCircle } from "lucide-react"
import { PageHeader } from "@/components/admin/page-header"

export default function DashboardPage() {
  const { workshops, loading, error } = useWorkshops()

  const handleCreate = () => {
    console.log("Create new workshop")
  }

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-[80vh] gap-4">
        <Loader2 className="h-10 w-10 animate-spin text-primary" />
        <p className="text-muted-foreground font-medium">Đang tải dữ liệu hệ thống...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-[80vh] gap-4 p-8">
        <div className="bg-rose-50 text-rose-600 p-4 rounded-2xl flex items-center gap-3 border border-rose-100">
          <AlertCircle className="h-6 w-6" />
          <p className="font-bold">Lỗi kết nối: {error}</p>
        </div>
        <Button onClick={() => window.location.reload()} variant="outline">Thử lại</Button>
      </div>
    )
  }

  return (
    <div className="p-8 space-y-8 max-w-7xl mx-auto">
      <PageHeader 
        title="Tổng quan hệ thống"
        description="Thống kê hoạt động và quản lý các buổi Workshop UniHub"
        actionLabel="Tạo Workshop mới"
        onAction={handleCreate}
      />

      {/* Stats Cards */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm hover:shadow-md transition-shadow">
          <p className="text-sm font-medium text-slate-500 uppercase tracking-wider">Tổng Workshop</p>
          <p className="mt-2 text-4xl font-bold text-slate-900 leading-tight">
            {workshops.length}
          </p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm hover:shadow-md transition-shadow">
          <p className="text-sm font-medium text-slate-500 uppercase tracking-wider">Đang mở đăng ký</p>
          <p className="mt-2 text-4xl font-bold text-emerald-600 leading-tight">
            {workshops.filter((w) => w.status === "open").length}
          </p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm hover:shadow-md transition-shadow">
          <p className="text-sm font-medium text-slate-500 uppercase tracking-wider">Tổng vé thành công</p>
          <p className="mt-2 text-4xl font-bold text-primary leading-tight">
            {workshops.reduce((sum, w) => sum + w.registered, 0)}
          </p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm hover:shadow-md transition-shadow">
          <p className="text-sm font-medium text-slate-500 uppercase tracking-wider">
            Tỷ lệ lấp đầy
          </p>
          <p className="mt-2 text-4xl font-bold text-amber-600 leading-tight">
            {workshops.length > 0 
              ? Math.round(
                  (workshops.reduce((sum, w) => sum + w.registered, 0) /
                    workshops.reduce((sum, w) => sum + w.capacity, 0)) *
                    100
                )
              : 0}
            %
          </p>
        </div>
      </div>

      {/* Workshop Table Section */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-bold text-slate-900">Hoạt động gần đây</h2>
          <Button variant="ghost" className="text-primary hover:text-primary/80 hover:bg-indigo-50 font-bold" asChild>
            <a href="/admin/workshops">Xem chi tiết</a>
          </Button>
        </div>
        
        <WorkshopTable
          workshops={workshops.slice(0, 5)}
          showActions={false}
        />
      </div>
    </div>
  )
}

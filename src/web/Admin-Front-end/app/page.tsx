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

export default function AdminWorkshopsPage() {
  const handleEdit = (workshop: Workshop) => {
    console.log("Edit workshop:", workshop)
  }

  const handleDelete = (workshop: Workshop) => {
    console.log("Delete workshop:", workshop)
  }

  return (
    <AdminLayout>
      <div className="space-y-6">
        {/* Page Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-semibold text-slate-800">
              Quản lý Workshop
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Quản lý tất cả các workshop trong hệ thống UniHub
            </p>
          </div>
          <Button className="bg-indigo-700 hover:bg-indigo-800">
            <Plus className="mr-2 h-4 w-4" />
            Tạo Workshop mới
          </Button>
        </div>

        {/* Filters */}
        <div className="flex flex-wrap items-center gap-4">
          <div className="relative flex-1 min-w-[240px] max-w-md">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
            <Input
              type="search"
              placeholder="Tìm kiếm theo tiêu đề, diễn giả..."
              className="h-10 pl-10 text-sm"
            />
          </div>
          <Select defaultValue="all">
            <SelectTrigger className="w-[160px] h-10">
              <Filter className="mr-2 h-4 w-4 text-slate-400" />
              <SelectValue placeholder="Trạng thái" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Tất cả</SelectItem>
              <SelectItem value="open">Đang mở</SelectItem>
              <SelectItem value="completed">Đã xong</SelectItem>
              <SelectItem value="cancelled">Đã hủy</SelectItem>
            </SelectContent>
          </Select>
          <Select defaultValue="all">
            <SelectTrigger className="w-[160px] h-10">
              <SelectValue placeholder="Loại vé" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Tất cả</SelectItem>
              <SelectItem value="free">Miễn phí</SelectItem>
              <SelectItem value="paid">Có phí</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div className="rounded-lg border border-slate-200 bg-white p-4">
            <p className="text-sm font-medium text-slate-500">Tổng Workshop</p>
            <p className="mt-1 text-2xl font-semibold text-slate-800">
              {mockWorkshops.length}
            </p>
          </div>
          <div className="rounded-lg border border-slate-200 bg-white p-4">
            <p className="text-sm font-medium text-slate-500">Đang mở</p>
            <p className="mt-1 text-2xl font-semibold text-emerald-600">
              {mockWorkshops.filter((w) => w.status === "open").length}
            </p>
          </div>
          <div className="rounded-lg border border-slate-200 bg-white p-4">
            <p className="text-sm font-medium text-slate-500">Tổng đăng ký</p>
            <p className="mt-1 text-2xl font-semibold text-indigo-600">
              {mockWorkshops.reduce((sum, w) => sum + w.registered, 0)}
            </p>
          </div>
          <div className="rounded-lg border border-slate-200 bg-white p-4">
            <p className="text-sm font-medium text-slate-500">
              Tỷ lệ lấp đầy TB
            </p>
            <p className="mt-1 text-2xl font-semibold text-amber-600">
              {Math.round(
                (mockWorkshops.reduce((sum, w) => sum + w.registered, 0) /
                  mockWorkshops.reduce((sum, w) => sum + w.capacity, 0)) *
                  100
              )}
              %
            </p>
          </div>
        </div>

        {/* Workshop Table */}
        <WorkshopTable
          workshops={mockWorkshops}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      </div>
    </AdminLayout>
  )
}

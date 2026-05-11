"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import {
  LayoutDashboard,
  CalendarDays,
  FileSpreadsheet,
  Settings,
  GraduationCap,
} from "lucide-react"
import { cn } from "@/lib/utils"

const menuItems = [
  {
    label: "Dashboard",
    href: "/admin",
    icon: LayoutDashboard,
  },
  {
    label: "Quản lý Workshop",
    href: "/admin/workshops",
    icon: CalendarDays,
  },
  {
    label: "Nhật ký CSV Import",
    href: "/admin/csv-logs",
    icon: FileSpreadsheet,
  },
  {
    label: "Cài đặt",
    href: "/admin/settings",
    icon: Settings,
  },
]

export function AdminSidebar() {
  const pathname = usePathname()

  return (
    <aside className="fixed left-0 top-0 z-40 h-screen w-64 border-r border-slate-200 bg-white">
      {/* Logo */}
      <div className="flex h-16 items-center gap-2 border-b border-slate-200 px-6">
        <div className="flex h-9 w-9 items-center justify-center rounded-md bg-indigo-700">
          <GraduationCap className="h-5 w-5 text-white" />
        </div>
        <span className="text-xl font-semibold text-slate-800">UniHub</span>
      </div>

      {/* Navigation */}
      <nav className="flex flex-col gap-1 p-4">
        {menuItems.map((item) => {
          const isActive = pathname === item.href || 
            (item.href !== "/admin" && pathname.startsWith(item.href))
          
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center gap-3 rounded-md px-3 py-2.5 text-sm font-medium transition-colors",
                isActive
                  ? "bg-indigo-50 text-indigo-700"
                  : "text-slate-600 hover:bg-slate-100 hover:text-slate-800"
              )}
            >
              <item.icon className={cn(
                "h-5 w-5",
                isActive ? "text-indigo-700" : "text-slate-500"
              )} />
              {item.label}
            </Link>
          )
        })}
      </nav>
    </aside>
  )
}

"use client"

import { useState, useEffect, useCallback } from "react"
import { Bell, Ticket, QrCode, Clock, MapPin, X } from "lucide-react"
import { QRCodeSVG } from "qrcode.react"
import { 
  Popover, 
  PopoverContent, 
  PopoverTrigger 
} from "@/components/ui/popover"
import { 
  Dialog, 
  DialogContent, 
  DialogHeader, 
  DialogTitle,
  DialogDescription
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Separator } from "@/components/ui/separator"
import { api, auth } from "@/lib/api-client"
import { format } from "date-fns"
import { vi } from "date-fns/locale"

interface Registration {
  id: string
  workshopId: string
  status: string
  ticketSignature?: string
  workshopTitle: string
  workshopRoom: string
  startTime: string
}

export function NotificationBell() {
  const [registrations, setRegistrations] = useState<Registration[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedTicket, setSelectedTicket] = useState<Registration | null>(null)
  const [isQRModalOpen, setIsQRModalOpen] = useState(false)

  const fetchRegistrations = useCallback(async () => {
    if (!auth.isAuthenticated()) return
    setLoading(true)
    try {
      const response = await api.get<Registration[]>('/api/v1/registrations/my')
      if (response.success && response.data) {
        // Chỉ lấy các vé đã thành công và có chữ ký
        const validTickets = response.data.filter(r => r.status === 'SUCCESS' && r.ticketSignature)
        setRegistrations(validTickets)
      }
    } catch (error) {
      console.error("Lỗi khi tải thông báo:", error)
    } finally {
      setLoading(false)
    }
  }, [])

  const user = auth.getUser() as any
  const userId = user?.id || 'guest'

  useEffect(() => {
    setRegistrations([]) // Reset dữ liệu cũ
    fetchRegistrations()
  }, [fetchRegistrations, userId])

  const openQR = (reg: Registration) => {
    setSelectedTicket(reg)
    setIsQRModalOpen(true)
  }

  const downloadQR = () => {
    if (!selectedTicket) return
    
    const svg = document.querySelector(".qr-container svg") as SVGElement
    if (!svg) return

    const canvas = document.createElement("canvas")
    const svgData = new XMLSerializer().serializeToString(svg)
    const img = new Image()
    const svgBlob = new Blob([svgData], { type: "image/svg+xml;charset=utf-8" })
    const url = URL.createObjectURL(svgBlob)

    img.onload = () => {
      canvas.width = img.width * 2 // Tăng độ phân giải
      canvas.height = img.height * 2
      const ctx = canvas.getContext("2d")
      if (ctx) {
        ctx.fillStyle = "white"
        ctx.fillRect(0, 0, canvas.width, canvas.height)
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height)
        const pngUrl = canvas.toDataURL("image/png")
        const downloadLink = document.createElement("a")
        downloadLink.href = pngUrl
        downloadLink.download = `UniHub-QR-${selectedTicket.workshopTitle.replace(/\s+/g, '-')}.png`
        document.body.appendChild(downloadLink)
        downloadLink.click()
        document.body.removeChild(downloadLink)
      }
      URL.revokeObjectURL(url)
    }
    img.src = url
  }

  const safeFormat = (dateStr: string | undefined, formatStr: string) => {
    if (!dateStr) return "--:--"
    const date = new Date(dateStr)
    if (isNaN(date.getTime()) || date.getFullYear() <= 1) return "--:--"
    return format(date, formatStr, { locale: vi })
  }

  return (
    <>
      <Popover onOpenChange={(open) => open && fetchRegistrations()}>
        <PopoverTrigger asChild>
          <Button variant="ghost" size="icon" className="relative hover:bg-slate-100 transition-colors">
            <Bell className="h-5 w-5 text-slate-600" />
            {registrations.length > 0 && (
              <span className="absolute right-2 top-2 h-2.5 w-2.5 rounded-full bg-destructive border-2 border-white animate-pulse" />
            )}
            <span className="sr-only">Thông báo</span>
          </Button>
        </PopoverTrigger>
        <PopoverContent align="end" className="w-80 p-0 shadow-2xl border-slate-200 overflow-hidden rounded-xl">
          <div className="bg-slate-50/50 p-4 border-b">
            <h3 className="font-bold text-slate-900 flex items-center gap-2">
              <Ticket className="h-4 w-4 text-primary" />
              Vé của tôi ({registrations.length})
            </h3>
          </div>
          
          <ScrollArea className="h-[350px]">
            {loading && (
              <div className="p-8 text-center text-sm text-slate-400">Đang tải vé...</div>
            )}
            {!loading && registrations.length === 0 && (
              <div className="p-12 text-center">
                <div className="mx-auto w-12 h-12 bg-slate-100 rounded-full flex items-center justify-center mb-3">
                  <Bell className="h-6 w-6 text-slate-300" />
                </div>
                <p className="text-sm text-slate-500">Bạn chưa có vé nào thành công.</p>
              </div>
            )}
            <div className="flex flex-col">
              {registrations.map((reg) => (
                <div 
                  key={reg.id} 
                  className="p-4 hover:bg-slate-50 transition-colors group border-b last:border-0"
                >
                  <p className="font-semibold text-sm text-slate-900 line-clamp-1 group-hover:text-primary transition-colors">
                    {reg.workshopTitle}
                  </p>
                  <div className="mt-2 space-y-1">
                    <div className="flex items-center gap-2 text-xs text-slate-500">
                      <Clock className="h-3 w-3" />
                      {safeFormat(reg.startTime, "HH:mm, dd/MM")}
                    </div>
                    <div className="flex items-center gap-2 text-xs text-slate-500">
                      <MapPin className="h-3 w-3" />
                      {reg.workshopRoom}
                    </div>
                  </div>
                  <Button 
                    variant="outline" 
                    size="sm" 
                    className="w-full mt-3 h-8 text-xs font-bold border-primary/20 hover:bg-primary hover:text-white transition-all gap-2"
                    onClick={() => openQR(reg)}
                  >
                    <QrCode className="h-3.5 w-3.5" />
                    XEM MÃ QR
                  </Button>
                </div>
              ))}
            </div>
          </ScrollArea>
          <div className="p-3 bg-slate-50 text-center border-t">
            <button className="text-[10px] font-bold text-slate-400 hover:text-primary uppercase tracking-widest">
              Xem tất cả lịch sử
            </button>
          </div>
        </PopoverContent>
      </Popover>

      {/* QR Code Modal */}
      <Dialog open={isQRModalOpen} onOpenChange={setIsQRModalOpen}>
        <DialogContent className="sm:max-w-[380px] p-0 overflow-hidden border-none shadow-2xl rounded-3xl">
          <div className="bg-primary p-6 text-center relative">
            <div className="mx-auto w-12 h-12 bg-white/10 rounded-xl flex items-center justify-center mb-3 backdrop-blur-sm">
              <QrCode className="h-6 w-6 text-white" />
            </div>
            <DialogTitle className="text-lg font-bold text-white mb-1 leading-tight">
              Mã Check-in Workshop
            </DialogTitle>
            <DialogDescription className="text-white/70 text-xs">
              Đưa mã này cho nhân viên tại cửa.
            </DialogDescription>
          </div>

          <div className="p-6 bg-white text-center">
            <div className="bg-slate-50 p-4 rounded-2xl inline-block shadow-inner border border-slate-100 qr-container">
              {selectedTicket?.ticketSignature && (
                <QRCodeSVG 
                  value={selectedTicket.ticketSignature}
                  size={180}
                  level="H"
                  includeMargin={false}
                  imageSettings={{
                    src: "/icon.svg", 
                    x: undefined,
                    y: undefined,
                    height: 36,
                    width: 36,
                    excavate: true,
                  }}
                />
              )}
            </div>

            <div className="mt-4">
              <Button 
                variant="default" 
                size="default" 
                className="w-full font-bold gap-2 rounded-xl shadow-lg shadow-primary/20 hover:shadow-primary/40 transition-all h-11"
                onClick={downloadQR}
              >
                <QrCode className="h-4 w-4" />
                TẢI VỀ MÁY (PNG)
              </Button>
            </div>

            <div className="mt-6 space-y-1">
              <p className="font-bold text-base text-slate-900 tracking-tight line-clamp-1">
                {selectedTicket?.workshopTitle}
              </p>
              <div className="flex items-center justify-center gap-3 text-xs font-medium text-slate-400">
                <span className="flex items-center gap-1">
                  <Clock className="h-3 w-3 text-primary" />
                  {safeFormat(selectedTicket?.startTime, "HH:mm")}
                </span>
                <span className="w-1 h-1 rounded-full bg-slate-200" />
                <span className="flex items-center gap-1">
                  <MapPin className="h-3 w-3 text-primary" />
                  {selectedTicket?.workshopRoom}
                </span>
              </div>
            </div>

            <p className="mt-6 text-[9px] text-slate-400 uppercase tracking-widest font-medium">
              UniHub Ticket System • ID: {selectedTicket?.id.slice(0, 8)}
            </p>
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}

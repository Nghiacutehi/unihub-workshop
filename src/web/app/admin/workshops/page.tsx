'use client'

import { useState, useEffect } from 'react'
import { supabase } from '@/lib/supabase'
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { QRCodeSVG } from 'qrcode.react'
import { Calendar, MapPin, Users, QrCode, Edit, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { PageHeader } from '@/components/admin/page-header'
import { signTicket } from '@/lib/crypto'

export default function AdminWorkshopsPage() {
  const [workshops, setWorkshops] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedWorkshop, setSelectedWorkshop] = useState<any>(null)
  const [registrations, setRegistrations] = useState<any[]>([])
  const [loadingRegs, setLoadingRegs] = useState(false)

  // 1. Fetch danh sách Workshop
  useEffect(() => {
    fetchWorkshops()
  }, [])

  const fetchWorkshops = async () => {
    setLoading(true)
    const { data, error } = await supabase
      .from('workshops')
      .select('*')
      .order('start_time', { ascending: false })

    if (error) toast.error('Lỗi tải danh sách Workshop')
    else setWorkshops(data || [])
    setLoading(false)
  }

  const handleCreate = () => {
    toast.info('Tính năng Tạo Workshop đang được phát triển')
  }

  // 2. Fetch danh sách sinh viên đăng ký khi nhấn xem QR
  const fetchRegistrations = async (workshopId: string) => {
    setLoadingRegs(true)
    const { data, error } = await supabase
      .from('registrations')
      .select(`
        id,
        user_id,
        ticket_signature,
        users (
          user_id,
          full_name,
          email
        )
      `)
      .eq('workshop_id', workshopId)
      .eq('status', 'SUCCESS') // Chỉ lấy những vé đã thanh toán thành công

    if (error) toast.error('Lỗi tải danh sách đăng ký')
    else setRegistrations(data || [])
    setLoadingRegs(false)
  }

  const formatDateTime = (str: string) => {
    return new Date(str).toLocaleString('vi-VN', {
      day: '2-digit', month: '2-digit', year: 'numeric',
      hour: '2-digit', minute: '2-digit'
    })
  }

  return (
    <div className="p-8 space-y-8 max-w-7xl mx-auto">
      <PageHeader
        title="Quản lý Workshop"
        description="Chỉnh sửa nội dung và quản lý mã QR điểm danh."
        actionLabel="Tạo Workshop mới"
        onAction={handleCreate}
      />

      {loading ? (
        <div className="flex justify-center py-20">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : (
        <div className="rounded-xl border bg-card shadow-sm overflow-hidden">
          <Table>
            <TableHeader className="bg-muted/50">
              <TableRow>
                <TableHead className="w-[300px]">Workshop</TableHead>
                <TableHead>Địa điểm & Thời gian</TableHead>
                <TableHead>Trạng thái</TableHead>
                <TableHead>Chỗ ngồi</TableHead>
                <TableHead className="text-right">Thao tác</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {workshops.map((ws) => (
                <TableRow key={ws.id} className="hover:bg-muted/30 transition-colors">
                  <TableCell>
                    <div className="font-bold text-base">{ws.title}</div>
                    <div className="text-xs text-muted-foreground mt-1 line-clamp-1">{ws.speaker}</div>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2 text-sm">
                      <Calendar className="h-3.5 w-3.5 text-primary" />
                      {formatDateTime(ws.start_time)}
                    </div>
                    <div className="flex items-center gap-2 text-sm text-muted-foreground mt-1">
                      <MapPin className="h-3.5 w-3.5" />
                      {ws.room}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant={ws.status === 'PUBLISHED' ? 'default' : 'secondary'}>
                      {ws.status}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <Users className="h-4 w-4 text-muted-foreground" />
                      <span className="text-sm font-medium">
                        {ws.capacity - ws.available_seats} / {ws.capacity}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button variant="outline" size="sm" onClick={() => toast.info('Tính năng Edit đang phát triển')}>
                        <Edit className="h-4 w-4 mr-1" /> Sửa
                      </Button>

                      <Dialog>
                        <DialogTrigger asChild>
                          <Button variant="secondary" size="sm" onClick={() => fetchRegistrations(ws.id)}>
                            <QrCode className="h-4 w-4 mr-1" /> Xem QR
                          </Button>
                        </DialogTrigger>
                        <DialogContent className="!max-w-[90vw] !w-[90vw] max-h-[90vh] overflow-y-auto">
                          <DialogHeader>
                            <DialogTitle>Mã QR Sinh viên - {ws.title}</DialogTitle>
                          </DialogHeader>

                          {loadingRegs ? (
                            <div className="flex justify-center py-10"><Loader2 className="animate-spin" /></div>
                          ) : (
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 p-4">
                              {registrations.length === 0 ? (
                                <div className="col-span-full text-center py-10 text-muted-foreground">
                                  Chưa có sinh viên nào đăng ký thành công.
                                </div>
                              ) : (
                                registrations.map((reg) => (
                                  <div key={reg.id} className="flex flex-col items-center p-4 border rounded-2xl bg-white shadow-sm">
                                    <div className="bg-slate-50 p-3 rounded-xl mb-4">
                                      <QRCodeSVG
                                        value={JSON.stringify({
                                          sid: reg.users?.user_id,
                                          uid: reg.user_id, // Thêm UUID để Backend xử lý siêu tốc
                                          wid: ws.id,
                                          sig: signTicket(reg.users?.user_id || '', ws.id)
                                        })}
                                        size={160}
                                        level="H"
                                        includeMargin={true}
                                      />
                                    </div>
                                    <div className="text-center space-y-1">
                                      <div className="font-bold text-slate-900">{reg.users?.full_name}</div>
                                      <div className="text-xs font-mono text-primary bg-primary/5 px-2 py-0.5 rounded-full inline-block">
                                        {reg.users?.user_id}
                                      </div>
                                    </div>
                                  </div>
                                ))
                              )}
                            </div>
                          )}
                        </DialogContent>
                      </Dialog>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  )
}

import { useState, useCallback } from 'react'
import { api, APIError } from '@/lib/api-client'
import { toast } from 'sonner'

export function useRegistration(workshopId: string) {
  const [isRegistering, setIsRegistering] = useState(false)
  const [regStatus, setRegStatus] = useState<string | null>(null)
  const [waitingPosition, setWaitingPosition] = useState<number | null>(null)

  const pollWaitingRoom = useCallback(async () => {
    const interval = setInterval(async () => {
      try {
        const response = await api.get<any>(`/api/v1/registrations/waiting-room/${workshopId}`)
        
        if (response.data?.status === 1) { // GRANTED
          clearInterval(interval)
          handleRegister() // Thử đăng ký lại khi đã được cấp quyền
        } else if (response.data?.position) {
          setWaitingPosition(response.data.position)
        }
      } catch (error) {
        console.error('Lỗi Polling Waiting Room:', error)
      }
    }, 5000)
    return () => clearInterval(interval)
  }, [workshopId])

  const pollRegistrationStatus = useCallback(async (correlationId: string) => {
    const interval = setInterval(async () => {
      try {
        const response = await api.get<any>(`/api/v1/registrations/status/${correlationId}`)
        
        if (response.data?.status !== 'PROCESSING') {
          clearInterval(interval)
          setIsRegistering(false)
          setRegStatus(null)
          
          if (response.data?.status === 'SUCCESS' || response.data?.status === 'PENDING_PAYMENT') {
            toast.success('Đăng ký thành công! Vui lòng kiểm tra email hoặc mục cá nhân.')
            // Tự động reload trang hoặc cập nhật UI nếu cần
            window.location.reload() 
          } else {
            toast.error(response.data?.message || 'Yêu cầu đăng ký bị từ chối.')
          }
        }
      } catch (error) {
        console.error('Lỗi Polling Status:', error)
      }
    }, 2000)
    return () => clearInterval(interval)
  }, [])

  const handleRegister = useCallback(async () => {
    setIsRegistering(true)
    setRegStatus('Đang gửi yêu cầu...')
    
    try {
      const response = await api.post<any>('/api/v1/registrations', { 
        workshop_id: workshopId 
      })

      // 202 Accepted — Đã vào hàng đợi xử lý
      if (response.data?.correlationId) {
        setRegStatus('Hệ thống đang xử lý...')
        pollRegistrationStatus(response.data.correlationId)
      }
    } catch (error) {
      if (error instanceof APIError) {
        if (error.status === 401) {
          toast.error('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.')
        } else if (error.status === 429) {
          setRegStatus('Đang trong phòng chờ...')
          pollWaitingRoom()
        } else {
          toast.error(error.message || 'Đăng ký thất bại')
        }
      } else {
        toast.error(error instanceof Error ? error.message : 'Lỗi kết nối hệ thống')
      }
      setIsRegistering(false)
      setRegStatus(null)
    }
  }, [workshopId, pollWaitingRoom, pollRegistrationStatus])

  return {
    isRegistering,
    regStatus,
    waitingPosition,
    handleRegister
  }
}

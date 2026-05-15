import { useState, useCallback, useRef, useEffect } from 'react'
import { api, APIError } from '@/lib/api-client'
import { toast } from 'sonner'

export function useRegistration(workshopId: string) {
  const [isRegistering, setIsRegistering] = useState(false)
  const [regStatus, setRegStatus] = useState<string | null>(null)
  const [waitingPosition, setWaitingPosition] = useState<number | null>(null)
  const [showPaymentDialog, setShowPaymentDialog] = useState(false)
  const [paymentInfo, setPaymentInfo] = useState<{ url: string, amount: number, title: string } | null>(null)

  // Dùng Ref để phá vỡ vòng lặp phụ thuộc giữa handleRegister và pollWaitingRoom
  const handleRegisterRef = useRef<() => Promise<void>>()

  // 1. Polling Status
  const pollRegistrationStatus = useCallback(async (correlationId: string) => {
    const interval = setInterval(async () => {
      try {
        const response = await api.get<any>(`/api/v1/registrations/status/${correlationId}`)

        if (response.data?.status !== 'PROCESSING') {
          clearInterval(interval)
          setIsRegistering(false)
          setRegStatus(null)

          if (response.data?.status === 'SUCCESS') {
            toast.success('Đăng ký thành công! Vui lòng kiểm tra email.')
            window.dispatchEvent(new CustomEvent('registration-success', { detail: { workshopId } }))
            window.dispatchEvent(new CustomEvent(`workshop-reg-success-${workshopId}`))
          } else if (response.data?.status === 'PENDING_PAYMENT') {
            setPaymentInfo({
              url: response.data.paymentUrl,
              amount: response.data.paymentAmount,
              title: response.data.message || 'Thanh toán đăng ký Workshop'
            })
            setShowPaymentDialog(true)
            toast.info('Vui lòng hoàn tất thanh toán để nhận vé.')
            window.dispatchEvent(new CustomEvent('registration-success', { detail: { workshopId } }))
            window.dispatchEvent(new CustomEvent(`workshop-reg-success-${workshopId}`))
          } else {
            toast.error(response.data?.message || 'Yêu cầu đăng ký bị từ chối.')
          }
        }
      } catch (error) {
        console.error('Lỗi Polling Status:', error)
      }
    }, 2000)
    return () => clearInterval(interval)
  }, [workshopId])

  // 2. Polling Waiting Room
  const pollWaitingRoom = useCallback(async () => {
    const interval = setInterval(async () => {
      try {
        const response = await api.get<any>(`/api/v1/registrations/waiting-room/${workshopId}`)

        if (response.data?.status === 1) { // GRANTED
          clearInterval(interval)
          if (handleRegisterRef.current) {
            handleRegisterRef.current()
          }
        } else if (response.data?.position) {
          setWaitingPosition(response.data.position)
        }
      } catch (error) {
        console.error('Lỗi Polling Waiting Room:', error)
      }
    }, 5000)
    return () => clearInterval(interval)
  }, [workshopId])

  // 3. Handle Register
  const handleRegister = useCallback(async () => {
    setIsRegistering(true)
    setRegStatus('Đang gửi yêu cầu...')

    try {
      const response = await api.post<any>('/api/v1/registrations', {
        workshop_id: workshopId
      })

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
  }, [workshopId, pollRegistrationStatus, pollWaitingRoom])

  // Cập nhật ref mỗi khi handleRegister thay đổi
  useEffect(() => {
    handleRegisterRef.current = handleRegister
  }, [handleRegister])

  return {
    isRegistering,
    regStatus,
    waitingPosition,
    showPaymentDialog,
    setShowPaymentDialog,
    paymentInfo,
    handleRegister
  }
}

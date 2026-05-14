'use client'

import { useState, useEffect, useCallback } from 'react'
import { supabase } from '@/lib/supabase'
import { Workshop } from '@/components/admin/workshop-table'

export function useWorkshops() {
  const [workshops, setWorkshops] = useState<Workshop[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchWorkshops = useCallback(async () => {
    try {
      setLoading(true)
      const { data, error } = await supabase
        .from('workshops')
        .select('*')
        .order('start_time', { ascending: false })

      if (error) throw error

      // Map dữ liệu từ DB sang format của Component Workshop
      const formattedData: Workshop[] = (data || []).map(ws => ({
        id: ws.id,
        title: ws.title,
        speaker: ws.speaker || 'Chưa xác định',
        datetime: ws.start_time,
        capacity: ws.capacity,
        registered: ws.capacity - ws.available_seats,
        price: Number(ws.price),
        status: ws.status.toLowerCase() as any
      }))

      setWorkshops(formattedData)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchWorkshops()
  }, [fetchWorkshops])

  return {
    workshops,
    loading,
    error,
    refresh: fetchWorkshops
  }
}

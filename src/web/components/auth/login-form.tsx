'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Eye, EyeOff, Mail, Lock } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { toast } from 'sonner'
import { supabase } from '@/lib/supabase'

export function LoginForm() {
  const router = useRouter()
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    rememberMe: false,
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    
    try {
      // 1. Truy vấn bảng users bằng email HOẶC user_id (MSSV)
      const { data: user, error } = await supabase
        .from('users')
        .select('id, user_id, password_hash, full_name, email, role')
        .or(`email.eq.${formData.email},user_id.eq.${formData.email}`)
        .single()

      if (error || !user) {
        throw new Error('Tài khoản không tồn tại trong hệ thống')
      }

      // 2. So sánh mật khẩu (Plaintext cho dev, sau này dùng bcrypt)
      if (user.password_hash !== formData.password) {
        throw new Error('Mật khẩu không chính xác')
      }

      // 3. Lưu session vào Cookie để Middleware và Client đọc được
      const sessionData = {
        id: user.id,
        user_id: user.user_id,
        full_name: user.full_name,
        email: user.email,
        role: user.role,
      }
      
      const cookieValue = encodeURIComponent(JSON.stringify(sessionData))
      document.cookie = `unihub_session=${cookieValue}; path=/; max-age=3600; SameSite=Lax`

      toast.success(`Chào mừng ${user.full_name}!`)

      // 4. Chuyển hướng chuẩn xác theo Role từ Database
      if (user.role === 'ADMIN' || user.role === 'STAFF') {
        router.push('/admin')
      } else {
        router.push('/')
      }
    } catch (error: any) {
      toast.error(error.message || 'Đăng nhập thất bại')
    } finally {
      setIsLoading(false)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }))
  }

  return (
    <form onSubmit={handleSubmit} className="w-full space-y-6">
      {/* Email Field */}
      <div className="space-y-2">
        <label htmlFor="email" className="text-sm font-medium text-foreground">
          Email hoặc Mã sinh viên
        </label>
        <div className="relative">
          <Mail className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
          <Input
            id="email"
            name="email"
            type="email"
            placeholder="Nhập email hoặc mã sinh viên"
            value={formData.email}
            onChange={handleChange}
            className="pl-10 h-11 bg-background border-input"
            required
          />
        </div>
      </div>

      {/* Password Field */}
      <div className="space-y-2">
        <label htmlFor="password" className="text-sm font-medium text-foreground">
          Mật khẩu
        </label>
        <div className="relative">
          <Lock className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
          <Input
            id="password"
            name="password"
            type={showPassword ? 'text' : 'password'}
            placeholder="Nhập mật khẩu"
            value={formData.password}
            onChange={handleChange}
            className="pl-10 pr-10 h-11 bg-background border-input"
            required
          />
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
            aria-label={showPassword ? 'Ẩn mật khẩu' : 'Hiện mật khẩu'}
          >
            {showPassword ? (
              <EyeOff className="h-5 w-5" />
            ) : (
              <Eye className="h-5 w-5" />
            )}
          </button>
        </div>
      </div>

      {/* Remember Me & Forgot Password */}
      <div className="flex items-center justify-between text-sm">
        <label className="flex items-center gap-2 cursor-pointer">
          <Checkbox
            name="rememberMe"
            checked={formData.rememberMe}
            onCheckedChange={(checked) =>
              setFormData(prev => ({
                ...prev,
                rememberMe: checked === true,
              }))
            }
          />
          <span className="text-foreground">Ghi nhớ tôi</span>
        </label>
        <a
          href="#"
          className="text-primary hover:underline font-medium"
        >
          Quên mật khẩu?
        </a>
      </div>

      {/* Submit Button */}
      <Button
        type="submit"
        disabled={isLoading}
        className="w-full h-11 text-base font-semibold rounded-md"
      >
        {isLoading ? 'Đang đăng nhập...' : 'Đăng nhập'}
      </Button>

      {/* Sign Up Link */}
      <p className="text-center text-sm text-muted-foreground">
        Chưa có tài khoản?{' '}
        <a
          href="#"
          className="text-primary hover:underline font-medium"
        >
          Đăng ký ngay
        </a>
      </p>
    </form>
  )
}

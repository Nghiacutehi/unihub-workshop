import { NextResponse, type NextRequest } from 'next/server'

/**
 * Middleware bảo vệ route — đọc session từ Cookie (không phụ thuộc Supabase Auth).
 * Logic phân quyền tuân thủ auth.md mục 3.2.
 */
export function middleware(request: NextRequest) {
  const sessionCookie = request.cookies.get('unihub_session')?.value
  const pathname = request.nextUrl.pathname

  // 1. Chưa đăng nhập + không phải trang login + không phải API -> chuyển về login
  if (!sessionCookie && !pathname.startsWith('/login') && !pathname.startsWith('/api')) {
    return NextResponse.redirect(new URL('/login', request.url))
  }

  // 2. Đã đăng nhập + đang ở trang login -> chuyển về trang chủ theo role
  if (sessionCookie && pathname.startsWith('/login')) {
    try {
      const user = JSON.parse(decodeURIComponent(sessionCookie))
      const isAdmin = user.role === 'ADMIN' || user.role === 'STAFF'
      const dest = isAdmin ? '/admin' : '/'
      return NextResponse.redirect(new URL(dest, request.url))
    } catch {
      // Cookie hỏng -> xóa và cho vào login
      const res = NextResponse.redirect(new URL('/login', request.url))
      res.cookies.delete('unihub_session')
      return res
    }
  }

  // 3. Bảo vệ /admin — chỉ ADMIN hoặc STAFF mới vào được
  if (sessionCookie && pathname.startsWith('/admin')) {
    try {
      const user = JSON.parse(decodeURIComponent(sessionCookie))
      const isAdmin = user.role === 'ADMIN' || user.role === 'STAFF'
      if (!isAdmin) {
        return NextResponse.redirect(new URL('/', request.url))
      }
    } catch {
      return NextResponse.redirect(new URL('/login', request.url))
    }
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    '/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)',
  ],
}

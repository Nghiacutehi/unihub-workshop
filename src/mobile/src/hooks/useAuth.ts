import { useState } from 'react';
import { Alert } from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { supabase } from '../services/supabase';

const SESSION_KEY = '@unihub_session';

/** Thông tin phiên đăng nhập lưu trên máy */
export type UserSession = {
  id: string;
  user_id: string;
  full_name: string;
  email: string;
  role: string;
};

/**
 * Hook xử lý xác thực người dùng.
 * Truy vấn trực tiếp bảng `users` trong PostgreSQL (đúng auth.md mục 3.1).
 */
export function useAuth() {
  const [loading, setLoading] = useState(false);

  /** Đăng nhập bằng email + password, so sánh với password_hash trong DB */
  const login = async (email: string, password: string): Promise<UserSession | null> => {
    if (!email || !password) {
      Alert.alert('Thiếu thông tin', 'Vui lòng nhập Email và Mật khẩu');
      return null;
    }

    setLoading(true);
    try {
      // 1. Truy vấn bảng users bằng email HOẶC user_id (MSSV)
      const { data: user, error } = await supabase
        .from('users')
        .select('id, user_id, password_hash, full_name, email, role')
        .or(`email.eq.${email},user_id.eq.${email}`)
        .single();

      if (error || !user) {
        throw new Error('Tài khoản không tồn tại trên hệ thống');
      }

      // 2. So sánh mật khẩu (Plaintext cho dev)
      if (user.password_hash !== password) {
        throw new Error('Mật khẩu không chính xác');
      }

      // 3. Kiểm tra quyền — Chỉ Staff/Admin được vào Mobile
      if (user.role !== 'STAFF' && user.role !== 'ADMIN') {
        throw new Error(`Tài khoản sinh viên (${user.user_id}) không có quyền vào ứng dụng Staff`);
      }

      // 4. Tạo session và lưu vào AsyncStorage
      const session: UserSession = {
        id: user.id,
        user_id: user.user_id,
        full_name: user.full_name,
        email: user.email,
        role: user.role,
      };
      await AsyncStorage.setItem(SESSION_KEY, JSON.stringify(session));

      return session;
    } catch (error: any) {
      Alert.alert('Đăng nhập thất bại', error.message);
      return null;
    } finally {
      setLoading(false);
    }
  };

  /** Đăng xuất — xóa session khỏi bộ nhớ máy */
  const logout = async () => {
    await AsyncStorage.removeItem(SESSION_KEY);
  };

  /** Đọc session đã lưu (dùng khi mở lại app) */
  const getSession = async (): Promise<UserSession | null> => {
    const raw = await AsyncStorage.getItem(SESSION_KEY);
    return raw ? JSON.parse(raw) : null;
  };

  return { login, logout, getSession, loading };
}

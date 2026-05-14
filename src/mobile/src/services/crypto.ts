/**
 * File này hiện tại chỉ chứa các cấu hình cơ bản.
 * Logic mã hóa AES đã được loại bỏ để tăng tính ổn định cho hệ thống Offline.
 */

export const API_BASE_URL = 'http://192.168.1.8:3000';
const PUBLIC_KEY = 'LIenPZCXeWuDsEvJ0vcUL8I874dsEtRkT+nkukiQLFI=';

export async function fetchPublicKey(): Promise<string | null> {
  return PUBLIC_KEY;
}

// Giữ lại hàm trống để tránh lỗi ở các file khác nếu có gọi (sẽ sớm dọn dẹp triệt để)
export function verifyTicket(): boolean { return true; }
export function encryptData(data: string): string { return data; }
export function decryptData(data: string): string { return data; }

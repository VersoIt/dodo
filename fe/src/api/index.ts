import axios, { AxiosInstance, AxiosResponse } from 'axios'
import type { ApiResponse, Product, Order, PromoCode, Analytics, User } from '../types'

const client: AxiosInstance = axios.create({
  baseURL: '',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Auth Interceptor
client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Global Error Interceptor
client.interceptors.response.use(
  (response: AxiosResponse) => response.data,
  (error) => {
    const message = error.response?.data?.error || 'Произошла непредвиденная ошибка'
    // Here we can trigger global toast or redirect to login on 401
    return Promise.reject(new Error(message))
  }
)

// --- Catalog API ---
export const catalogApi = {
  getProducts: (): Promise<ApiResponse<Product[]>> => client.get('/api/v1/catalog/products'),
  getProduct: (id: string): Promise<ApiResponse<Product>> => client.get(`/api/v1/catalog/products/${id}`),
  createProduct: (payload: Partial<Product>): Promise<ApiResponse<Product>> => client.post('/api/v1/catalog/products', payload),
  updateProduct: (id: string, payload: Partial<Product>): Promise<ApiResponse<Product>> => client.put(`/api/v1/catalog/products/${id}`, payload),
}

// --- Orders API ---
export const ordersApi = {
  createOrder: (data: any): Promise<ApiResponse<Order>> => client.post('/api/v1/orders', data),
  getMyOrders: (): Promise<ApiResponse<Order[]>> => client.get('/api/v1/orders/my'),
  getAllOrders: (): Promise<ApiResponse<Order[]>> => client.get('/api/v1/orders/all'),
  getOrder: (id: string): Promise<ApiResponse<Order>> => client.get(`/api/v1/orders/${id}`),
  updateStatus: (id: string, status: string): Promise<ApiResponse<Order>> => client.patch(`/api/v1/orders/${id}/status`, { status }),
  payOrder: (id: string): Promise<ApiResponse<any>> => client.post(`/api/v1/orders/${id}/pay`),
  listPromos: (): Promise<ApiResponse<PromoCode[]>> => client.get('/api/v1/promos'),
  createPromoCode: (data: Partial<PromoCode>): Promise<ApiResponse<PromoCode>> => client.post('/api/v1/promos', data),
  deletePromo: (id: string): Promise<ApiResponse<any>> => client.delete(`/api/v1/promos/${id}`),
  checkPromoCode: (code: string): Promise<ApiResponse<PromoCode>> => client.get(`/api/v1/promos/check/${code}`),
  getAnalytics: (): Promise<ApiResponse<Analytics>> => client.get('/api/v1/orders/analytics'),
  exportReport: (start?: string, end?: string): Promise<ApiResponse<Blob>> => 
    client.get('/api/v1/orders/report', { 
      params: { start, end },
      responseType: 'blob' 
    }).then(res => ({ success: true, data: res.data })),
}

// --- Kitchen API ---
export const kitchenApi = {
  listTickets: (): Promise<ApiResponse<{ tickets: any[] }>> => client.get('/api/v1/kitchen/tickets'),
  updateStatus: (id: string, status: string): Promise<ApiResponse<any>> => client.patch(`/api/v1/kitchen/tickets/${id}/status`, { status }),
}

// --- Logistics API ---
export const logisticsApi = {
  listDeliveries: (): Promise<ApiResponse<{ deliveries: any[] }>> => client.get('/api/v1/logistics/deliveries'),
  assignCourier: (id: string, courierId: string): Promise<ApiResponse<any>> => client.post(`/api/v1/logistics/orders/${id}/assign`, { courier_id: courierId }),
  updateStatus: (id: string, status: string): Promise<ApiResponse<any>> => client.patch(`/api/v1/logistics/orders/${id}/status`, { status }),
  updateLocation: (id: string, lat: number, lng: number): Promise<ApiResponse<any>> => client.post(`/api/v1/logistics/orders/${id}/location`, { lat, lng }),
}

// --- Auth API ---
export const authApi = {
  login: (credentials: any): Promise<ApiResponse<{ token: string, user: User }>> => client.post('/api/v1/auth/login', credentials),
  register: (data: any): Promise<ApiResponse<User>> => client.post('/api/v1/auth/register', data),
  getMe: (): Promise<ApiResponse<User>> => client.get('/api/v1/auth/me'),
  updateProfile: (data: Partial<User>): Promise<ApiResponse<User>> => client.patch('/api/v1/auth/me', data),
}

// --- Chat API ---
export const chatApi = {
  getHistory: (orderId: string | number, limit: number = 50): Promise<ApiResponse<any[]>> => 
    client.get(`/api/v1/chat/history`, { params: { order_id: orderId, limit } }),
  syncMessages: (orderId: string | number, afterId: number): Promise<ApiResponse<any[]>> => 
    client.get(`/api/v1/chat/sync`, { params: { order_id: orderId, after_id: afterId } }),
}

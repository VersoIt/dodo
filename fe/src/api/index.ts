import axios from 'axios'

// --- Catalog API ---
export const catalogApi = {
  getProducts: () => axios.get('/api/v1/catalog/products').then(r => r.data),
  getProduct: (id: string) => axios.get(`/api/v1/catalog/products/${id}`).then(r => r.data),
  createProduct: (payload: any) => axios.post('/api/v1/catalog/products', payload).then(r => r.data),
  updateProduct: (id: string, payload: any) => axios.put(`/api/v1/catalog/products/${id}`, payload).then(r => r.data),
}

// --- Orders API ---
export const ordersApi = {
  createOrder: (data: any) => axios.post('/api/v1/orders', data).then(r => r.data),
  getMyOrders: () => axios.get('/api/v1/orders/my').then(r => r.data),
  getAllOrders: () => axios.get('/api/v1/orders/all').then(r => r.data),
  getOrder: (id: string) => axios.get(`/api/v1/orders/${id}`).then(r => r.data),
  updateStatus: (id: string, status: string) => axios.patch(`/api/v1/orders/${id}/status`, { status }).then(r => r.data),
  payOrder: (id: string) => axios.post(`/api/v1/orders/${id}/pay`).then(r => r.data),
}

// --- Auth API ---
export const authApi = {
  login: (credentials: any) => axios.post('/api/v1/auth/login', credentials).then(r => r.data),
  register: (data: any) => axios.post('/api/v1/auth/register', data).then(r => r.data),
  getMe: () => axios.get('/api/v1/auth/me').then(r => r.data),
  updateProfile: (data: any) => axios.patch('/api/v1/auth/me', data).then(r => r.data),
}

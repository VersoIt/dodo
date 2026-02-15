// --- Generic API Response ---
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: string
}

// --- Domain Models ---
export interface Product {
  id: string
  name: string
  description: string
  base_price: number
  category_id: number
  image_url: string
  is_available: boolean
}

export interface Order {
  order_id: string
  order_number: string
  status: string
  final_price: number
  items: OrderItem[]
  address: Address
  created_at: string
}

export interface OrderItem {
  product_id: string
  product_name: string
  quantity: number
}

export interface Address {
  city: string
  street: string
  house: string
  apartment?: string
  floor?: string
  entrance?: string
  comment?: string
}

export interface PromoCode {
  id: string
  code: string
  type: 'percent' | 'fixed'
  amount: number
  active: boolean
}

export interface Analytics {
  total_revenue: number
  orders_count: number
  avg_check: number
  top_products: ProductStat[]
}

export interface ProductStat {
  name: string
  count: number
  revenue: number
}

export interface User {
  id: string
  email: string
  full_name: string
  role: string
  phone: string
}

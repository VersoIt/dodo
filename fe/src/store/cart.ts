import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface CartItem {
  id: string
  name: string
  price: number
  quantity: number
  imageUrl: string
}

export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>([])

  const totalItems = computed(() => items.value.reduce((acc, item) => acc + item.quantity, 0))
  const totalPrice = computed(() => items.value.reduce((acc, item) => acc + item.price * item.quantity, 0))

  function addToCart(product: any) {
    const existingItem = items.value.find(item => item.id === product.id)
    if (existingItem) {
      existingItem.quantity++
    } else {
      items.value.push({
        id: product.id,
        name: product.name,
        price: product.price,
        imageUrl: product.imageUrl,
        quantity: 1
      })
    }
  }

  function removeFromCart(productId: string) {
    items.value = items.value.filter(item => item.id !== productId)
  }

  function updateQuantity(productId: string, quantity: number) {
    const item = items.value.find(item => item.id === productId)
    if (item) {
      item.quantity = Math.max(1, quantity)
    }
  }

  function clearCart() {
    items.value = []
  }

  return { items, totalItems, totalPrice, addToCart, removeFromCart, updateQuantity, clearCart }
})

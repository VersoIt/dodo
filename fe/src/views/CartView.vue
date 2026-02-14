<script setup lang="ts">
import { useCartStore } from '../store/cart'
import { useRouter } from 'vue-router'
import { Trash2, Minus, Plus, ShoppingBag } from 'lucide-vue-next'
import axios from 'axios'
import { ref } from 'vue'

const cartStore = useCartStore()
const router = useRouter()
const isPlacingOrder = ref(false)

const handleCheckout = async () => {
  if (cartStore.items.length === 0) return
  
  try {
    isPlacingOrder.value = true
    // In a real app, we'd check if user is logged in
    const orderData = {
      items: cartStore.items.map(item => ({
        productId: item.id,
        quantity: item.quantity
      })),
      address: {
        street: 'Sample Street 123',
        city: 'Pizza City'
      }
    }
    
    const response = await axios.post('/api/v1/orders', orderData)
    const orderId = response.data.id
    
    cartStore.clearCart()
    router.push(`/order/${orderId}`)
  } catch (error) {
    console.error('Checkout failed:', error)
    alert('Failed to place order. Please try again.')
  } finally {
    isPlacingOrder.value = false
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto">
    <h1 class="text-3xl font-bold mb-8">Your Cart</h1>

    <div v-if="cartStore.items.length === 0" class="card bg-base-100 shadow-xl p-12 text-center">
      <div class="flex flex-col items-center gap-4">
        <ShoppingBag class="w-16 h-16 text-base-content/20" />
        <h2 class="text-2xl font-semibold">Your cart is empty</h2>
        <p class="text-base-content/60">Looks like you haven't added any pizzas yet.</p>
        <router-link to="/" class="btn btn-primary mt-4">Go to Menu</router-link>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Items List -->
      <div class="lg:col-span-2 space-y-4">
        <div v-for="item in cartStore.items" :key="item.id" class="card card-side bg-base-100 shadow-sm border border-base-200 p-4">
          <figure class="w-24 h-24 rounded-xl overflow-hidden flex-shrink-0">
            <img :src="item.imageUrl" :alt="item.name" class="w-full h-full object-cover" />
          </figure>
          <div class="card-body py-0 px-4 justify-between">
            <div class="flex justify-between items-start">
              <div>
                <h3 class="font-bold text-lg">{{ item.name }}</h3>
                <p class="text-primary font-semibold">${{ item.price }}</p>
              </div>
              <button @click="cartStore.removeFromCart(item.id)" class="btn btn-ghost btn-sm btn-circle text-error">
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
            
            <div class="flex justify-between items-center">
              <div class="join border border-base-300">
                <button @click="cartStore.updateQuantity(item.id, item.quantity - 1)" class="btn btn-ghost btn-xs join-item">
                  <Minus class="w-3 h-3" />
                </button>
                <span class="btn btn-ghost btn-xs join-item pointer-events-none">{{ item.quantity }}</span>
                <button @click="cartStore.updateQuantity(item.id, item.quantity + 1)" class="btn btn-ghost btn-xs join-item">
                  <Plus class="w-3 h-3" />
                </button>
              </div>
              <p class="font-bold">${{ (item.price * item.quantity).toFixed(2) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Summary -->
      <div class="lg:col-span-1">
        <div class="card bg-base-100 shadow-xl border border-base-200">
          <div class="card-body">
            <h2 class="card-title mb-4">Order Summary</h2>
            <div class="space-y-2">
              <div class="flex justify-between">
                <span>Subtotal</span>
                <span>${{ cartStore.totalPrice.toFixed(2) }}</span>
              </div>
              <div class="flex justify-between">
                <span>Delivery</span>
                <span class="text-success">FREE</span>
              </div>
              <div class="divider"></div>
              <div class="flex justify-between font-bold text-xl">
                <span>Total</span>
                <span>${{ cartStore.totalPrice.toFixed(2) }}</span>
              </div>
            </div>
            <div class="card-actions mt-6">
              <button 
                @click="handleCheckout" 
                class="btn btn-primary btn-block"
                :disabled="isPlacingOrder"
              >
                <span v-if="isPlacingOrder" class="loading loading-spinner"></span>
                Place Order
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

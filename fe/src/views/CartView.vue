<script setup lang="ts">
import { useCartStore } from '../store/cart'
import { useRouter } from 'vue-router'
import { Trash2, Minus, Plus, ShoppingBag, CreditCard, QrCode, ShieldCheck, X } from 'lucide-vue-next'
import axios from 'axios'
import { inject, ref } from 'vue'

const cartStore = useCartStore()
const router = useRouter()
const isPlacingOrder = ref(false)
const isPaying = ref(false)
const showPaymentModal = ref(false)
const createdOrderId = ref<string | null>(null)
const addToast = inject('addToast') as (msg: string, type?: any) => void

const handleCheckout = async () => {
  if (cartStore.items.length === 0) return
  
  try {
    isPlacingOrder.value = true
    const orderData = {
      items: cartStore.items.map(item => ({
        product_id: item.id,
        quantity: item.quantity
      })),
      address: {
        street: 'Main Street 1',
        city: 'Pizza Town'
      }
    }
    
    const response = await axios.post('/api/v1/orders', orderData)
    if (response.data.success) {
      createdOrderId.value = response.data.data.order_id
      showPaymentModal.value = true
    } else {
      addToast(response.data.error || 'Failed to place order', 'error')
    }
  } catch (error: any) {
    console.error('Checkout failed:', error)
    addToast(error.response?.data?.error || 'Failed to place order. Are you logged in?', 'error')
  } finally {
    isPlacingOrder.value = false
  }
}

const confirmPayment = async () => {
  if (!createdOrderId.value) return
  
  try {
    isPaying.value = true
    // Simulate API call to pay for order
    const response = await axios.post(`/api/v1/orders/${createdOrderId.value}/pay`)
    
    if (response.data.success) {
      cartStore.clearCart()
      addToast('Payment successful! Your order is being prepared.', 'success')
      showPaymentModal.value = false
      router.push(`/order/${createdOrderId.value}/success`)
    } else {
      addToast(response.data.error || 'Payment failed', 'error')
    }
  } catch (error: any) {
    addToast('Payment failed. Please try again.', 'error')
  } finally {
    isPaying.value = false
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto pb-20">
    <h1 class="text-4xl font-black mb-8 tracking-tight">Your Cart</h1>

    <div v-if="cartStore.items.length === 0" class="card bg-base-100 shadow-xl p-12 text-center border border-base-200">
      <div class="flex flex-col items-center gap-4">
        <div class="bg-base-200 p-6 rounded-full">
          <ShoppingBag class="w-16 h-16 text-base-content/20" />
        </div>
        <h2 class="text-2xl font-bold">Your cart is empty</h2>
        <p class="text-base-content/60">Looks like you haven't added any pizzas yet.</p>
        <router-link to="/" class="btn btn-primary btn-wide mt-4 rounded-xl">Go to Menu</router-link>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Items List -->
      <div class="lg:col-span-2 space-y-4">
        <div v-for="item in cartStore.items" :key="item.id" class="card card-side bg-base-100 shadow-sm border border-base-200 p-4 transition-all hover:shadow-md">
          <figure class="w-24 h-24 rounded-2xl overflow-hidden flex-shrink-0 shadow-inner">
            <img :src="item.imageUrl" :alt="item.name" class="w-full h-full object-cover" />
          </figure>
          <div class="card-body py-0 px-4 justify-between">
            <div class="flex justify-between items-start">
              <div>
                <h3 class="font-bold text-lg leading-tight">{{ item.name }}</h3>
                <p class="text-primary font-black mt-1">${{ item.price.toFixed(2) }}</p>
              </div>
              <button @click="cartStore.removeFromCart(item.id)" class="btn btn-ghost btn-sm btn-circle text-error/40 hover:text-error hover:bg-error/10">
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
            
            <div class="flex justify-between items-center mt-2">
              <div class="join bg-base-200/50 p-0.5 rounded-xl border border-base-300">
                <button @click="cartStore.updateQuantity(item.id, item.quantity - 1)" class="btn btn-ghost btn-xs join-item px-2">
                  <Minus class="w-3 h-3" />
                </button>
                <span class="btn btn-ghost btn-xs join-item pointer-events-none font-bold text-xs">{{ item.quantity }}</span>
                <button @click="cartStore.updateQuantity(item.id, item.quantity + 1)" class="btn btn-ghost btn-xs join-item px-2">
                  <Plus class="w-3 h-3" />
                </button>
              </div>
              <p class="font-black text-lg">${{ (item.price * item.quantity).toFixed(2) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Summary -->
      <div class="lg:col-span-1">
        <div class="card bg-base-100 shadow-2xl border border-primary/10 sticky top-24 overflow-hidden">
          <div class="bg-primary/5 px-6 py-4 border-b border-primary/10">
            <h2 class="font-bold text-lg flex items-center gap-2"><CreditCard class="w-5 h-5 text-primary" /> Order Summary</h2>
          </div>
          <div class="card-body p-6">
            <div class="space-y-3">
              <div class="flex justify-between text-sm">
                <span class="opacity-60">Subtotal</span>
                <span class="font-bold">${{ cartStore.totalPrice.toFixed(2) }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="opacity-60">Delivery Fee</span>
                <span class="text-success font-bold uppercase text-[10px] tracking-widest">Free</span>
              </div>
              <div class="divider my-1 opacity-50"></div>
              <div class="flex justify-between items-end pt-2">
                <span class="font-bold text-lg">Total</span>
                <div class="text-right">
                  <p class="text-3xl font-black text-primary leading-none">${{ cartStore.totalPrice.toFixed(2) }}</p>
                  <p class="text-[10px] uppercase font-bold opacity-40 mt-1">VAT Included</p>
                </div>
              </div>
            </div>
            <div class="card-actions mt-8">
              <button 
                @click="handleCheckout" 
                class="btn btn-primary btn-block btn-lg rounded-2xl shadow-lg shadow-primary/30 gap-3"
                :disabled="isPlacingOrder"
              >
                <span v-if="isPlacingOrder" class="loading loading-spinner"></span>
                Place Order
              </button>
            </div>
            
            <div class="mt-6 flex items-center justify-center gap-2 opacity-40">
              <ShieldCheck class="w-4 h-4" />
              <span class="text-[10px] font-bold uppercase tracking-widest">Secure Checkout</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- PAYMENT MODAL -->
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="showPaymentModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
        <div class="fixed inset-0 bg-black/80 backdrop-blur-md" @click="!isPaying && (showPaymentModal = false)"></div>
        
        <Transition
          appear
          enter-active-class="transition duration-500 delay-100 ease-out"
          enter-from-class="opacity-0 translate-y-12 scale-95"
          enter-to-class="opacity-100 translate-y-0 scale-100"
        >
          <div class="relative bg-base-100 w-full max-w-md rounded-[2.5rem] shadow-2xl border border-white/10 overflow-hidden">
            <!-- Close Button -->
            <button 
              v-if="!isPaying"
              @click="showPaymentModal = false" 
              class="absolute top-6 right-6 btn btn-ghost btn-circle btn-sm bg-base-200/50"
            >
              <X class="w-4 h-4" />
            </button>

            <div class="p-10 flex flex-col items-center">
              <div class="bg-primary/10 p-4 rounded-3xl mb-6">
                <QrCode class="w-12 h-12 text-primary" />
              </div>
              <h3 class="font-black text-3xl mb-2 text-center leading-tight">Scan to Pay</h3>
              <p class="text-base-content/60 text-center mb-8 px-4">
                Total amount: <span class="font-black text-primary">${{ cartStore.totalPrice.toFixed(2) }}</span>
              </p>

              <!-- QR MOCKUP -->
              <div class="relative bg-white p-6 rounded-[2rem] shadow-inner mb-8 group overflow-hidden">
                <div class="absolute inset-0 bg-primary/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                <svg viewBox="0 0 100 100" class="w-48 h-48 text-black fill-current relative z-10">
                  <path d="M0 0h35v10H10v25H0V0zM65 0h35v35h-10V10H65V0zM0 65h10v25h25v10H0V65zM90 65h10v35H65v-10h25V65z" />
                  <path d="M20 20h15v15H20V20zM65 20h15v15H65V20zM20 65h15v15H20V65z" />
                  <path d="M42 20h16v16h-16V20zM20 42h16v16h-16V42zM42 42h16v16h-16V42zM64 42h16v16h-16V42zM42 64h16v16h-16V64zM64 64h8v8h-8V64zM72 72h8v8h-8V72z" />
                </svg>
              </div>

              <div class="flex flex-col w-full gap-3">
                <button 
                  @click="confirmPayment" 
                  class="btn btn-primary btn-xl btn-block rounded-2xl h-16 text-lg font-black shadow-xl shadow-primary/30"
                  :disabled="isPaying"
                >
                  <span v-if="isPaying" class="loading loading-spinner"></span>
                  <Check v-else class="w-6 h-6" />
                  Confirm Payment
                </button>
                <div class="flex items-center justify-center gap-2 mt-4 opacity-40">
                  <ShieldCheck class="w-4 h-4" />
                  <span class="text-[10px] font-bold uppercase tracking-widest">End-to-End Encrypted</span>
                </div>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.btn-xl {
  @apply min-h-[4rem] text-lg;
}
</style>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import { CheckCircle2, Clock, Truck, PackageCheck, CreditCard } from 'lucide-vue-next'

const route = useRoute()
const orderId = route.params.id as string
const order = ref<any>(null)
const loading = ref(true)
let pollInterval: any = null

const fetchOrder = async () => {
  try {
    const response = await axios.get(`/api/v1/orders/${orderId}`)
    order.value = response.data
  } catch (error) {
    console.error('Failed to fetch order:', error)
    // Mock for demo
    if (!order.value) {
      order.value = {
        id: orderId,
        status: 0, // Received
        totalPrice: 25.98,
        items: [
          { name: 'Pepperoni', quantity: 1, price: 14.99 },
          { name: 'Margherita', quantity: 1, price: 12.99 }
        ]
      }
    }
  } finally {
    loading.value = false
  }
}

const getStatusIndex = (status: number) => {
  // 0: Received, 1: Kitchen, 2: Ready, 3: InDelivery, 4: Delivered, 5: Cancelled
  return status
}

const steps = [
  { label: 'Order Received', icon: PackageCheck },
  { label: 'In the Kitchen', icon: Clock },
  { label: 'Ready for Pickup', icon: CheckCircle2 },
  { label: 'Out for Delivery', icon: Truck },
  { label: 'Enjoy your Pizza!', icon: CheckCircle2 }
]

const handlePayment = async () => {
  try {
    await axios.post(`/api/v1/orders/${orderId}/pay`)
    fetchOrder()
  } catch (error) {
    alert('Payment failed')
  }
}

onMounted(() => {
  fetchOrder()
  pollInterval = setInterval(fetchOrder, 5000)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>

<template>
  <div class="max-w-3xl mx-auto py-8">
    <div v-if="loading" class="flex justify-center py-20">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>

    <div v-else-if="order" class="space-y-8">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-3xl font-bold">Order #{{ orderId.slice(0, 8) }}</h1>
          <p class="text-base-content/60">Thank you for your order!</p>
        </div>
        <div class="badge badge-lg badge-primary font-bold">
          {{ steps[getStatusIndex(order.status)]?.label || 'Processing' }}
        </div>
      </div>

      <!-- Payment Section if needed -->
      <div v-if="order.status === 0" class="card bg-secondary text-secondary-content shadow-xl">
        <div class="card-body flex-row items-center gap-6">
          <CreditCard class="w-12 h-12" />
          <div class="flex-grow">
            <h2 class="card-title text-2xl">Payment Required</h2>
            <p>Please complete your payment to start the preparation.</p>
          </div>
          <button @click="handlePayment" class="btn btn-neutral btn-lg">Pay Now</button>
        </div>
      </div>

      <!-- Tracking Steps -->
      <div class="card bg-base-100 shadow-xl border border-base-200">
        <div class="card-body">
          <h2 class="card-title mb-6 text-xl">Order Status</h2>
          <ul class="steps steps-vertical md:steps-horizontal w-full">
            <li 
              v-for="(step, index) in steps" 
              :key="index"
              class="step"
              :class="{ 'step-primary': index <= getStatusIndex(order.status) }"
            >
              <div class="flex flex-col items-center gap-1">
                <component :is="step.icon" class="w-6 h-6 mb-1" />
                <span class="text-xs font-semibold">{{ step.label }}</span>
              </div>
            </li>
          </ul>
        </div>
      </div>

      <!-- Order Details -->
      <div class="card bg-base-100 shadow-xl border border-base-200">
        <div class="card-body">
          <h2 class="card-title mb-4">Order Summary</h2>
          <div class="divide-y divide-base-200">
            <div v-for="item in order.items" :key="item.name" class="py-3 flex justify-between">
              <div>
                <span class="font-bold">{{ item.quantity }}x</span> {{ item.name }}
              </div>
              <span class="font-semibold">${{ (item.price * item.quantity).toFixed(2) }}</span>
            </div>
            <div class="pt-4 mt-4 flex justify-between text-xl font-bold">
              <span>Total</span>
              <span>${{ order.totalPrice.toFixed(2) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.step::after {
  @apply bg-base-300;
}
.step-primary::after {
  @apply bg-primary text-primary-content;
}
</style>

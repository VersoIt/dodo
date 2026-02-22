<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { PartyPopper, CheckCircle2, Home, Package } from 'lucide-vue-next'
import { ordersApi } from '../api'

const route = useRoute()
const router = useRouter()
const orderId = route.params.id as string
const orderNumber = ref('')
const loading = ref(true)

const fetchOrderDetails = async () => {
  try {
    const response = await ordersApi.getOrder(orderId)
    if (response.success && response.data) {
      orderNumber.value = response.data.order_number
    }
  } catch (err) {
    console.error('Failed to load order info')
  } finally {
    loading.value = false
  }
}

onMounted(fetchOrderDetails)
</script>

<template>
  <div class="min-h-[70vh] flex items-center justify-center px-4 py-12">
    <div class="max-w-md w-full text-center space-y-8 animate-in fade-in zoom-in-95 duration-700">
      <div class="relative inline-flex mb-4">
        <div class="absolute inset-0 bg-success/20 rounded-full animate-ping duration-1000"></div>
        <div class="relative bg-success text-success-content p-6 rounded-full shadow-2xl">
          <PartyPopper class="w-16 h-16" />
        </div>
      </div>

      <div class="space-y-3">
        <h1 class="text-4xl font-black tracking-tight text-base-content uppercase">Спасибо!</h1>
        <p class="text-xl font-bold text-primary uppercase tracking-widest">Заказ успешно оформлен</p>
        <p class="text-base-content/60 leading-relaxed pt-2">
          Ваша оплата подтверждена. Наши мастера уже начали готовить вашу вкуснейшую пиццу!
        </p>
      </div>

      <div class="bg-base-100 border border-base-300 rounded-[2.5rem] p-8 shadow-xl space-y-6">
        <div class="flex justify-between items-center text-sm uppercase font-bold tracking-widest opacity-40">
          <span>Номер заказа</span>
          <span>Статус</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-2xl font-black font-mono">#{{ orderNumber || '...' }}</span>
          <div class="badge badge-success gap-2 py-4 px-4 font-bold text-xs uppercase text-white">
            <CheckCircle2 class="w-4 h-4" /> Оплачено
          </div>
        </div>
        
        <div class="divider opacity-50"></div>
        
        <div class="space-y-3">
          <button @click="router.push(`/order/${orderId}`)" class="btn btn-primary btn-block btn-lg rounded-2xl gap-3 shadow-lg shadow-primary/20">
            <Package class="w-5 h-5" /> Отслеживать статус
          </button>
          <button @click="router.push('/')" class="btn btn-ghost btn-block gap-2 opacity-60 hover:opacity-100">
            <Home class="w-4 h-4" /> На главную
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

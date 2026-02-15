<script setup lang="ts">
import { ref, onMounted, inject, computed, watch } from 'vue'
import { ChefHat, Clock, CheckCircle2, Play, User, AlertCircle } from 'lucide-vue-next'
import { useAuthStore } from '../store/auth'
import { ordersApi } from '../api'
import { ORDER_STATUS } from '../constants'
import { shortId } from '../utils/format'
import AppModal from '../components/shared/AppModal.vue'

const orders = ref<any[]>([])
const loading = ref(true)
const authStore = useAuthStore()
const addToast = inject('addToast') as (msg: string, type?: any) => void

const showConfirmModal = ref(false)
const pendingAction = ref<{ orderId: string, status: string, title: string, description: string } | null>(null)

const fetchKitchenOrders = async () => {
  try {
    loading.value = true
    const res = await ordersApi.getAllOrders()
    if (res.success) orders.value = res.data || []
  } catch (err) {
    console.error('Failed to fetch orders:', err)
  } finally {
    loading.value = false
  }
}

const filteredOrders = computed(() => {
  const myId = authStore.user?.id || authStore.user?.user_id
  if (!myId) return []
  
  return orders.value.filter((o: any) => {
    const chefId = o.chef_id || o.chefId
    const isUnassignedPaid = o.status === ORDER_STATUS.PAID && (!chefId || chefId === "")
    const isMyActiveCooking = o.status === ORDER_STATUS.COOKING && chefId === myId
    return isUnassignedPaid || isMyActiveCooking
  })
})

const openConfirm = (order: any, nextStatus: string) => {
  const isStart = nextStatus === ORDER_STATUS.COOKING
  pendingAction.value = {
    orderId: order.order_id,
    status: nextStatus,
    title: isStart ? 'Принять заказ?' : 'Завершить готовку?',
    description: isStart 
      ? `Взять заказ #${order.order_number.split('-').pop()} в работу?` 
      : 'Подтвердите, что пицца готова к выдаче курьеру.'
  }
  showConfirmModal.value = true
}

const handleConfirm = async () => {
  if (!pendingAction.value) return
  try {
    const { orderId, status } = pendingAction.value
    await ordersApi.updateStatus(orderId, status)
    addToast(status === ORDER_STATUS.COOKING ? 'Вы начали готовить' : 'Заказ готов!', 'success')
    showConfirmModal.value = false
    pendingAction.value = null
    await fetchKitchenOrders()
  } catch (err: any) {
    addToast(err.response?.data?.error || 'Ошибка обновления', 'error')
    showConfirmModal.value = false
  }
}

watch(() => authStore.user, (u) => { if (u) fetchKitchenOrders() }, { immediate: true })
onMounted(fetchKitchenOrders)
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-12 gap-6">
      <div class="flex items-center gap-5">
        <div class="bg-primary p-5 rounded-[1.75rem] shadow-2xl shadow-primary/20 animate-in zoom-in duration-500">
          <ChefHat class="w-10 h-10 text-primary-content" />
        </div>
        <div>
          <h1 class="text-5xl font-black tracking-tighter uppercase italic">Кухня</h1>
          <p class="text-base-content/40 font-black uppercase text-[10px] tracking-[0.3em] mt-1 ml-1">Live Order Queue</p>
        </div>
      </div>
      <div class="flex items-center gap-4 bg-base-100 border border-base-200 p-2.5 rounded-3xl shadow-sm">
        <div class="w-12 h-12 rounded-2xl bg-primary/5 flex items-center justify-center text-primary font-black text-xl shadow-inner">
          {{ authStore.user?.name?.charAt(0).toUpperCase() }}
        </div>
        <div class="pr-6">
          <p class="text-[9px] font-black uppercase opacity-30 leading-none mb-1.5 tracking-widest">Chef on duty</p>
          <p class="text-sm font-black">{{ authStore.user?.name || 'Loading...' }}</p>
        </div>
      </div>
    </div>

    <div v-if="loading && orders.length === 0" class="flex justify-center py-40"><span class="loading loading-spinner loading-lg text-primary"></span></div>

    <div v-else-if="filteredOrders.length === 0" class="text-center py-48 bg-base-100 rounded-[4rem] border-2 border-dashed border-base-300 animate-in fade-in duration-700">
      <div class="bg-base-200 p-10 rounded-full inline-block mb-8 shadow-inner"><Clock class="w-16 h-16 opacity-10" /></div>
      <h2 class="text-3xl font-black opacity-20 uppercase tracking-tighter italic">Все заказы готовы</h2>
      <p class="text-base-content/30 max-w-xs mx-auto mt-3 font-bold uppercase text-[10px] tracking-widest">Отличная работа! Отдохните немного.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-10">
      <div v-for="order in filteredOrders" :key="order.order_id" 
        class="card bg-base-100 shadow-xl border border-base-200 overflow-hidden transition-all duration-500 hover:shadow-2xl group"
        :class="{ 'ring-4 ring-primary ring-inset shadow-primary/10 scale-[1.01]': order.status === ORDER_STATUS.COOKING }"
      >
        <div class="p-10">
          <div class="flex justify-between items-start mb-8">
            <div class="space-y-2">
              <div class="flex items-center gap-3">
                <span class="text-4xl font-black tracking-tighter italic">#{{ order.order_number.split('-').pop() }}</span>
                <span v-if="order.status === ORDER_STATUS.COOKING" class="badge badge-primary font-black text-[9px] h-5 tracking-widest px-3">ACTIVE</span>
              </div>
              <p class="text-[10px] font-bold opacity-30 uppercase tracking-[0.2em] font-mono">{{ shortId(order.order_id) }}</p>
            </div>
            <div class="badge font-black uppercase text-[10px] py-4 px-5 rounded-xl border-none shadow-sm" 
              :class="order.status === ORDER_STATUS.COOKING ? 'badge-warning' : 'badge-info'">
              {{ order.status === ORDER_STATUS.COOKING ? 'Готовится' : 'Оплачен' }}
            </div>
          </div>

          <div class="divider opacity-10 my-0"></div>
          
          <div class="py-10 space-y-4">
            <div v-for="item in order.items" :key="item.product_id" class="flex justify-between items-center bg-base-200/40 p-6 rounded-[1.5rem] border border-base-300/10 group-hover:bg-base-200/60 transition-colors">
              <span class="font-black text-base tracking-tight"><span class="text-primary font-black mr-4 text-2xl italic">x{{ item.quantity }}</span> {{ item.product_name }}</span>
            </div>
          </div>

          <div class="card-actions mt-4">
            <button 
              v-if="order.status === ORDER_STATUS.PAID"
              @click="openConfirm(order, ORDER_STATUS.COOKING)"
              class="btn btn-primary btn-block h-20 rounded-3xl gap-4 font-black uppercase shadow-2xl shadow-primary/20 text-xl italic transition-all hover:scale-[1.02]"
            >
              <Play class="w-7 h-7" /> Начать готовку
            </button>
            <button 
              v-if="order.status === ORDER_STATUS.COOKING"
              @click="openConfirm(order, ORDER_STATUS.READY)"
              class="btn btn-success btn-block h-20 rounded-3xl gap-4 font-black uppercase shadow-2xl shadow-success/20 text-xl italic transition-all hover:scale-[1.02] text-white"
            >
              <CheckCircle2 class="w-7 h-7" /> Завершить
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- ACTION CONFIRMATION MODAL -->
    <AppModal :show="showConfirmModal" @close="showConfirmModal = false">
      <div class="p-12 text-center">
          <div class="bg-primary/10 w-28 h-24 rounded-[2.5rem] flex items-center justify-center mx-auto mb-10 shadow-inner">
            <AlertCircle class="w-14 h-14 text-primary" />
          </div>
          <h3 class="font-black text-4xl uppercase tracking-tighter mb-4 leading-none italic">{{ pendingAction?.title }}</h3>
          <p class="text-base-content/50 text-base mb-12 leading-relaxed font-bold px-6">
            {{ pendingAction?.description }}
          </p>
          <div class="flex flex-col gap-4">
            <button @click="handleConfirm" class="btn btn-primary h-20 rounded-3xl text-xl font-black uppercase shadow-2xl shadow-primary/20 tracking-tight">Подтвердить</button>
            <button @click="showConfirmModal = false" class="btn btn-ghost h-16 rounded-2xl font-black uppercase text-[10px] tracking-widest opacity-40">Отмена</button>
          </div>
      </div>
    </AppModal>
  </div>
</template>

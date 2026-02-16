<script setup lang="ts">
import { ref, onMounted, inject, computed, watch } from 'vue'
import { Truck, MapPin, CheckCircle2, Package, User, AlertCircle } from 'lucide-vue-next'
import { useAuthStore } from '../store/auth'
import { logisticsApi } from '../api'
import { ORDER_STATUS } from '../constants'
import { shortId } from '../utils/format'
import AppModal from '../components/shared/AppModal.vue'

const orders = ref<any[]>([])
const loading = ref(true)
const authStore = useAuthStore()
const addToast = inject('addToast') as (msg: string, type?: any) => void

const showConfirmModal = ref(false)
const pendingAction = ref<{ orderId: string, status: string, title: string, description: string } | null>(null)

const fetchLogisticsOrders = async () => {
  try {
    loading.value = true
    const res = await logisticsApi.listDeliveries()
    if (res.success && res.data) {
      orders.value = (res.data.deliveries || []).map((d: any) => ({
        ...d,
        // Map logistics status to ORDER_STATUS
        // pending -> READY
        // on_way -> DELIVERING
        // delivered -> COMPLETED
        status: d.status === 'pending' || d.status === 'assigned' ? ORDER_STATUS.READY : 
                d.status === 'on_way' ? ORDER_STATUS.DELIVERING : 
                d.status === 'delivered' ? ORDER_STATUS.COMPLETED : d.status,
        address: {
          city: d.city,
          street: d.street,
          house: d.house,
          apartment: d.apartment
        }
      }))
    }
  } catch (err) { console.error('Failed to fetch orders:', err) } finally { loading.value = false }
}

const filteredOrders = computed(() => {
  const myId = authStore.user?.id || authStore.user?.user_id
  if (!myId) return []
  return orders.value.filter((o: any) => {
    const courierId = o.courier_id || o.courierId
    const isUnassignedReady = o.status === ORDER_STATUS.READY && (!courierId || courierId === "")
    const isMyActiveDelivery = o.status === ORDER_STATUS.DELIVERING && courierId === myId
    return isUnassignedReady || isMyActiveDelivery
  })
})

const openConfirm = (order: any, nextStatus: string) => {
  const isStart = nextStatus === ORDER_STATUS.DELIVERING
  pendingAction.value = {
    orderId: order.order_id, status: nextStatus, title: isStart ? 'Забрать заказ?' : 'Заказ доставлен?',
    description: isStart ? `Начать доставку заказа #${order.order_number.split('-').pop()}?` : 'Подтвердите, что вы успешно передали пиццу клиенту.'
  }
  showConfirmModal.value = true
}

const handleConfirm = async () => {
  if (!pendingAction.value) return
  const myId = authStore.user?.id || authStore.user?.user_id
  if (!myId) return

  try {
    const { orderId, status } = pendingAction.value
    
    if (status === ORDER_STATUS.DELIVERING) {
      // First assign, then start delivery (on_way)
      await logisticsApi.assignCourier(orderId, myId)
      await logisticsApi.updateStatus(orderId, 'delivering')
    } else if (status === ORDER_STATUS.COMPLETED) {
      await logisticsApi.updateStatus(orderId, 'completed')
    }
    
    addToast(status === ORDER_STATUS.DELIVERING ? 'Удачной дороги!' : 'Заказ доставлен!', 'success')
    showConfirmModal.value = false; pendingAction.value = null; await fetchLogisticsOrders()
  } catch (err: any) { addToast(err.response?.data?.error || 'Ошибка обновления', 'error'); showConfirmModal.value = false }
}

watch(() => authStore.user, (u) => { if (u) fetchLogisticsOrders() }, { immediate: true })
onMounted(fetchLogisticsOrders)
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-12 gap-6">
      <div class="flex items-center gap-5">
        <div class="bg-secondary p-5 rounded-[1.75rem] shadow-2xl shadow-secondary/20 animate-in zoom-in duration-500"><Truck class="w-10 h-10 text-secondary-content" /></div>
        <div>
          <h1 class="text-5xl font-black tracking-tighter uppercase text-secondary">Доставка</h1>
          <p class="text-base-content/40 font-black uppercase text-[10px] tracking-[0.3em] mt-1 ml-1">Live Logistics Fleet</p>
        </div>
      </div>
      <div class="flex items-center gap-4 bg-base-100 border border-base-200 p-2.5 rounded-3xl shadow-sm">
        <div class="w-12 h-12 rounded-2xl bg-secondary/5 flex items-center justify-center text-secondary font-black text-xl shadow-inner">{{ authStore.user?.name?.charAt(0).toUpperCase() }}</div>
        <div class="pr-6"><p class="text-[9px] font-black uppercase opacity-30 leading-none mb-1.5 tracking-widest">Courier on route</p><p class="text-sm font-black">{{ authStore.user?.name || 'Loading...' }}</p></div>
      </div>
    </div>

    <div v-if="loading && orders.length === 0" class="flex justify-center py-40"><span class="loading loading-spinner loading-lg text-secondary"></span></div>
    <div v-else-if="filteredOrders.length === 0" class="text-center py-48 bg-base-100 rounded-[4rem] border-2 border-dashed border-base-300 animate-in fade-in duration-700">
      <div class="bg-base-200 p-10 rounded-full inline-block mb-8 shadow-inner"><Package class="w-16 h-16 opacity-10" /></div>
      <h2 class="text-3xl font-black opacity-20 uppercase tracking-tighter">Нет активных заказов</h2>
      <p class="text-base-content/30 max-w-xs mx-auto mt-3 font-bold uppercase text-[10px] tracking-widest">Ожидайте готовности новых пицц.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-10">
      <div v-for="order in filteredOrders" :key="order.order_id" class="card bg-base-100 shadow-xl border border-base-200 overflow-hidden transition-all duration-500 hover:shadow-2xl group" :class="{ 'ring-4 ring-secondary ring-inset shadow-secondary/10 scale-[1.01]': order.status === ORDER_STATUS.DELIVERING }">
        <div class="bg-secondary/5 px-10 py-5 flex justify-between items-center border-b border-secondary/10">
           <div><span class="text-2xl font-black text-secondary tracking-tighter">#{{ order.order_number.split('-').pop() }}</span><span v-if="order.status === ORDER_STATUS.DELIVERING" class="badge badge-secondary font-black text-[9px] h-5 tracking-widest px-3 ml-3">MY ROUTE</span></div>
           <div class="badge font-black uppercase text-[10px] py-4 px-5 rounded-xl border-none shadow-sm" :class="order.status === ORDER_STATUS.READY ? 'badge-success text-white' : 'badge-info'">{{ order.status === ORDER_STATUS.READY ? 'Готов к выдаче' : 'В пути' }}</div>
        </div>
        <div class="card-body p-10">
          <div class="flex items-start gap-8 mb-10">
            <div class="bg-error/10 p-6 rounded-[2rem] shadow-sm group-hover:scale-110 transition-transform duration-500"><MapPin class="w-10 h-10 text-error" /></div>
            <div class="flex-1">
              <p class="text-[10px] font-black uppercase tracking-[0.3em] text-base-content/30 mb-2">Target Destination</p>
              <p class="font-black text-3xl leading-tight mb-2 tracking-tight">{{ order.address.street }}</p>
              <p class="text-base-content/60 font-bold text-lg">{{ order.address.city }}</p>
              <div v-if="order.address.house" class="flex flex-wrap gap-4 mt-6"><div class="bg-base-200/50 border border-base-300 px-5 py-3 rounded-2xl text-base font-black">Дом {{ order.address.house }}</div><div v-if="order.address.apartment" class="bg-base-200/50 border border-base-300 px-5 py-3 rounded-2xl text-base font-black">Кв. {{ order.address.apartment }}</div></div>
            </div>
          </div>
          <div class="divider opacity-10 my-0"></div>
          <div class="py-10">
            <p class="text-[10px] font-black uppercase tracking-[0.3em] text-base-content/30 mb-5">Order Content</p>
            <div class="flex flex-wrap gap-3"><span v-for="item in order.items" :key="item.product_id" class="badge badge-lg bg-base-100 border-2 border-base-200 font-black text-sm py-6 px-6 rounded-[1.25rem] group-hover:border-secondary/20 transition-colors"><span class="text-secondary mr-3 font-black text-xl">{{ item.quantity }}x</span> {{ item.product_name }}</span></div>
          </div>
          <div class="card-actions mt-auto pt-4"><button v-if="order.status === ORDER_STATUS.READY" @click="openConfirm(order, ORDER_STATUS.DELIVERING)" class="btn btn-secondary btn-block h-20 rounded-3xl gap-4 font-black uppercase shadow-2xl shadow-secondary/20 text-xl transition-all hover:scale-[1.02]"><Truck class="w-8 h-8" /> Принять заказ</button><button v-if="order.status === ORDER_STATUS.DELIVERING" @click="openConfirm(order, ORDER_STATUS.COMPLETED)" class="btn btn-success btn-block h-20 rounded-3xl gap-4 font-black uppercase shadow-2xl shadow-success/20 text-xl transition-all hover:scale-[1.02] text-white"><CheckCircle2 class="w-8 h-8" /> Доставлено</button></div>
        </div>
      </div>
    </div>

    <AppModal :show="showConfirmModal" @close="showConfirmModal = false">
      <div class="p-12 text-center">
          <div class="bg-secondary/10 w-28 h-24 rounded-[2.5rem] flex items-center justify-center mx-auto mb-10 shadow-inner"><AlertCircle class="w-14 h-14 text-secondary" /></div>
          <h3 class="font-black text-4xl uppercase tracking-tighter mb-4 leading-none text-secondary">{{ pendingAction?.title }}</h3>
          <p class="text-base-content/50 text-base mb-12 leading-relaxed font-bold px-6">{{ pendingAction?.description }}</p>
          <div class="flex flex-col gap-4"><button @click="handleConfirm" class="btn btn-secondary h-20 rounded-3xl text-xl font-black uppercase shadow-2xl shadow-secondary/20 tracking-tight">Подтвердить</button><button @click="showConfirmModal = false" class="btn btn-ghost h-16 rounded-2xl font-black uppercase text-[10px] tracking-widest opacity-40">Отмена</button></div>
      </div>
    </AppModal>
  </div>
</template>

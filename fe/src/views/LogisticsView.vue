<script setup lang="ts">
import { ref, onMounted, inject, computed, watch } from 'vue'
import axios from 'axios'
import { Truck, MapPin, CheckCircle2, Package, User, AlertCircle, Play } from 'lucide-vue-next'
import { useAuthStore } from '../store/auth'

const orders = ref<any[]>([])
const loading = ref(true)
const authStore = useAuthStore()
const addToast = inject('addToast') as (msg: string, type?: any) => void

const showConfirmModal = ref(false)
const pendingAction = ref<{ orderId: string, status: string, title: string, description: string } | null>(null)

const STATUS_READY = 'ready'
const STATUS_DELIVERING = 'delivering'
const STATUS_COMPLETED = 'completed'

const fetchLogisticsOrders = async () => {
  try {
    loading.value = true
    const response = await axios.get('/api/v1/orders/all') 
    if (response.data.success) {
      orders.value = response.data.data || []
    }
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
    const courierId = o.courier_id || o.courierId
    const isUnassignedReady = o.status === STATUS_READY && (!courierId || courierId === "")
    const isMyActiveDelivery = o.status === STATUS_DELIVERING && courierId === myId
    return isUnassignedReady || isMyActiveDelivery
  })
})

const openConfirm = (order: any, nextStatus: string) => {
  const isStart = nextStatus === STATUS_DELIVERING
  pendingAction.value = {
    orderId: order.order_id,
    status: nextStatus,
    title: isStart ? 'Забрать заказ?' : 'Доставлено?',
    description: isStart 
      ? `Взять заказ #${order.order_number.split('-').pop()} на доставку?` 
      : 'Подтвердите, что вы успешно передали заказ клиенту.'
  }
  showConfirmModal.value = true
}

const handleConfirm = async () => {
  if (!pendingAction.value) return
  
  try {
    const { orderId, status } = pendingAction.value
    await axios.patch(`/api/v1/orders/${orderId}/status`, { status })
    addToast(status === STATUS_DELIVERING ? 'Заказ принят в доставку' : 'Заказ доставлен!', 'success')
    showConfirmModal.value = false
    pendingAction.value = null
    await fetchLogisticsOrders()
  } catch (err: any) {
    addToast(err.response?.data?.error || 'Ошибка обновления', 'error')
    showConfirmModal.value = false
  }
}

watch(() => authStore.user, (u) => { if (u) fetchLogisticsOrders() }, { immediate: true })
onMounted(fetchLogisticsOrders)
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-10 gap-4">
      <div class="flex items-center gap-4">
        <div class="bg-secondary p-4 rounded-[1.5rem] shadow-xl shadow-secondary/20">
          <Truck class="w-10 h-10 text-secondary-content" />
        </div>
        <div>
          <h1 class="text-4xl font-black tracking-tighter uppercase italic text-secondary">Логистика</h1>
          <p class="text-base-content/50 font-bold uppercase text-[10px] tracking-[0.2em]">Доставка заказов</p>
        </div>
      </div>
      <div class="flex items-center gap-4 bg-base-100 border border-base-200 p-2 rounded-2xl shadow-sm">
        <div class="w-10 h-10 rounded-xl bg-base-200 flex items-center justify-center">
          <User class="w-5 h-5 opacity-40" />
        </div>
        <div class="pr-4">
          <p class="text-[10px] font-black uppercase opacity-30 leading-none mb-1">Курьер</p>
          <p class="text-sm font-black">{{ authStore.user?.name || 'Загрузка...' }}</p>
        </div>
      </div>
    </div>

    <div v-if="loading && orders.length === 0" class="flex justify-center py-32">
      <span class="loading loading-spinner loading-lg text-secondary"></span>
    </div>

    <div v-else-if="filteredOrders.length === 0" class="text-center py-40 bg-base-100 rounded-[3.5rem] border-2 border-dashed border-base-300">
      <div class="bg-base-200 p-8 rounded-full inline-block mb-6"><Package class="w-12 h-12 opacity-10" /></div>
      <h2 class="text-2xl font-black opacity-20 uppercase tracking-tighter">Нет заказов</h2>
      <p class="text-base-content/30 max-w-xs mx-auto mt-2 font-bold uppercase text-[10px] tracking-widest">Ждем, когда кухня закончит готовку</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <div v-for="order in filteredOrders" :key="order.order_id" 
        class="card bg-base-100 shadow-xl border border-base-200 overflow-hidden transition-all duration-300"
        :class="{ 'ring-4 ring-secondary ring-inset': order.status === STATUS_DELIVERING }"
      >
        <div class="bg-secondary/5 px-8 py-4 flex justify-between items-center border-b border-secondary/10">
           <span class="font-black text-secondary uppercase tracking-tighter">#{{ order.order_number.split('-').pop() }}</span>
           <div class="badge font-black uppercase text-[10px] py-3.5 px-4 rounded-lg border-none" 
            :class="order.status === STATUS_READY ? 'badge-success text-white' : 'badge-info'">
            {{ order.status === STATUS_READY ? 'Готов к выдаче' : 'В пути' }}
           </div>
        </div>
        
        <div class="card-body p-8">
          <div class="flex items-start gap-6 mb-8">
            <div class="bg-error/10 p-5 rounded-[1.5rem] shadow-sm"><MapPin class="w-8 h-8 text-error" /></div>
            <div class="flex-1">
              <p class="text-[10px] font-black uppercase tracking-widest text-base-content/30 mb-2">Адрес доставки</p>
              <p class="font-black text-2xl leading-tight mb-1">{{ order.address.street }}</p>
              <p class="text-base-content/60 font-bold tracking-tight">{{ order.address.city }}</p>
              <div v-if="order.address.house" class="flex gap-3 mt-5">
                <div class="bg-base-200 px-4 py-2 rounded-xl text-sm font-black">Дом: {{ order.address.house }}</div>
                <div v-if="order.address.apartment" class="bg-base-200 px-4 py-2 rounded-xl text-sm font-black">Кв: {{ order.address.apartment }}</div>
              </div>
            </div>
          </div>

          <div class="divider opacity-20 my-0"></div>

          <div class="py-8">
            <p class="text-[10px] font-black uppercase tracking-widest text-base-content/30 mb-4">Состав заказа</p>
            <div class="flex flex-wrap gap-2">
              <span v-for="item in order.items" :key="item.product_id" class="badge badge-lg bg-base-100 border-2 border-base-200 font-black text-xs py-5 px-5 rounded-2xl">
                <span class="text-secondary mr-2 font-black text-base">x{{ item.quantity }}</span> {{ item.product_name }}
              </span>
            </div>
          </div>

          <div class="card-actions">
            <button 
              v-if="order.status === STATUS_READY"
              @click="openConfirm(order, STATUS_DELIVERING)"
              class="btn btn-secondary btn-block rounded-2xl gap-3 font-black uppercase shadow-lg shadow-secondary/10 h-16 text-lg"
            >
              <Truck class="w-6 h-6" /> Забрать заказ
            </button>
            <button 
              v-if="order.status === STATUS_DELIVERING"
              @click="openConfirm(order, STATUS_COMPLETED)"
              class="btn btn-success btn-block rounded-2xl gap-3 font-black uppercase shadow-lg shadow-success/10 h-16 text-lg text-white"
            >
              <CheckCircle2 class="w-6 h-6" /> Доставлено
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Confirmation Modal -->
    <Transition name="modal-fade">
      <div v-if="showConfirmModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
        <div class="fixed inset-0 bg-black/80" @click="showConfirmModal = false"></div>
        <Transition name="modal-zoom" appear>
          <div class="relative bg-base-100 w-full max-w-sm rounded-[3rem] shadow-2xl border border-base-200 p-10 text-center">
              <div class="bg-secondary/10 w-24 h-24 rounded-full flex items-center justify-center mx-auto mb-8 animate-in zoom-in duration-300">
                <AlertCircle class="w-12 h-12 text-secondary" />
              </div>
              <h3 class="font-black text-3xl uppercase tracking-tighter mb-3 leading-none">{{ pendingAction?.title }}</h3>
              <p class="text-base-content/50 text-sm mb-10 leading-relaxed font-medium px-4">
                {{ pendingAction?.description }}
              </p>
              <div class="flex flex-col gap-3">
                <button @click="handleConfirm" class="btn btn-secondary btn-lg rounded-2xl h-16 font-black uppercase shadow-xl shadow-secondary/20 tracking-tight">
                  Подтвердить
                </button>
                <button @click="showConfirmModal = false" class="btn btn-ghost btn-lg rounded-2xl font-bold opacity-40">
                  Отмена
                </button>
              </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.15s ease; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
.modal-zoom-enter-active { transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1); }
.modal-zoom-leave-active { transition: all 0.15s ease-in; }
.modal-zoom-enter-from, .modal-zoom-leave-to { opacity: 0; transform: scale(0.97) translateY(8px); }
</style>

<script setup lang="ts">
import { ref, onMounted, inject, computed, watch } from 'vue'
import axios from 'axios'
import { ChefHat, Clock, CheckCircle2, Play, User, AlertCircle } from 'lucide-vue-next'
import { useAuthStore } from '../store/auth'

const orders = ref<any[]>([])
const loading = ref(true)
const authStore = useAuthStore()
const addToast = inject('addToast') as (msg: string, type?: any) => void

const showConfirmModal = ref(false)
const pendingAction = ref<{ orderId: string, status: string, title: string, description: string } | null>(null)

const STATUS_PAID = 'paid'
const STATUS_COOKING = 'cooking'
const STATUS_READY = 'ready'

const fetchKitchenOrders = async () => {
  try {
    loading.value = true
    const response = await axios.get('/api/v1/orders/all') 
    if (response.data.success) {
      orders.value = response.data.data || []
      // Debug log to see exactly what backend returns
      console.log('[Kitchen] Data from backend:', JSON.stringify(orders.value, null, 2))
    }
  } catch (err) {
    console.error('Failed to fetch orders:', err)
  } finally {
    loading.value = false
  }
}

const filteredOrders = computed(() => {
  const myId = authStore.user?.id || authStore.user?.user_id
  if (!myId) {
    console.warn('[Kitchen] No user ID found in store')
    return []
  }
  
  return orders.value.filter((o: any) => {
    const chefId = o.chef_id || o.chefId
    
    const isUnassignedPaid = o.status === STATUS_PAID && (!chefId || chefId === "")
    const isMyActiveCooking = o.status === STATUS_COOKING && chefId === myId
    
    if (o.status === STATUS_COOKING) {
      console.log(`[Kitchen] Order ${o.order_number}: status=cooking, chef_id=${chefId}, my_id=${myId}, MATCH=${chefId === myId}`)
    }
    
    return isUnassignedPaid || isMyActiveCooking
  })
})

const openConfirm = (order: any, nextStatus: string) => {
  const isStart = nextStatus === STATUS_COOKING
  pendingAction.value = {
    orderId: order.order_id,
    status: nextStatus,
    title: isStart ? 'Accept Order?' : 'Complete Cooking?',
    description: isStart 
      ? `Accept Order #${order.order_number.split('-').pop()}? This will assign it to you.` 
      : 'Finish preparing this pizza and move it to ready state.'
  }
  showConfirmModal.value = true
}

const handleConfirm = async () => {
  if (!pendingAction.value) return
  
  try {
    const { orderId, status } = pendingAction.value
    console.log(`[Kitchen] Patching order ${orderId} to status ${status}`)
    const resp = await axios.patch(`/api/v1/orders/${orderId}/status`, { status })
    console.log('[Kitchen] Patch response:', resp.data)
    
    addToast(status === STATUS_COOKING ? 'Started cooking!' : 'Pizza ready!', 'success')
    showConfirmModal.value = false
    pendingAction.value = null
    await fetchKitchenOrders()
  } catch (err: any) {
    console.error('[Kitchen] Status update failed:', err)
    addToast(err.response?.data?.error || 'Update failed', 'error')
    showConfirmModal.value = false
  }
}

watch(() => authStore.user, (u) => { if (u) fetchKitchenOrders() }, { immediate: true })
onMounted(fetchKitchenOrders)
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-10 gap-4">
      <div class="flex items-center gap-4">
        <div class="bg-primary p-4 rounded-[1.5rem] shadow-xl shadow-primary/20">
          <ChefHat class="w-10 h-10 text-primary-content" />
        </div>
        <div>
          <h1 class="text-4xl font-black tracking-tighter uppercase italic">Kitchen</h1>
          <p class="text-base-content/50 font-bold uppercase text-[10px] tracking-[0.2em]">Live Orders</p>
        </div>
      </div>
      <div class="flex items-center gap-4 bg-base-100 border border-base-200 p-2 rounded-2xl shadow-sm">
        <div class="w-10 h-10 rounded-xl bg-base-200 flex items-center justify-center">
          <User class="w-5 h-5 opacity-40" />
        </div>
        <div class="pr-4">
          <p class="text-[10px] font-black uppercase opacity-30 leading-none mb-1">Active Chef</p>
          <p class="text-sm font-black">{{ authStore.user?.name || 'Loading...' }}</p>
        </div>
      </div>
    </div>

    <div v-if="loading && orders.length === 0" class="flex justify-center py-32">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>

    <div v-else-if="filteredOrders.length === 0" class="text-center py-40 bg-base-100 rounded-[3.5rem] border-2 border-dashed border-base-300">
      <div class="bg-base-200 p-8 rounded-full inline-block mb-6"><Clock class="w-12 h-12 opacity-10" /></div>
      <h2 class="text-2xl font-black opacity-20 uppercase tracking-tighter">No orders found</h2>
      <p class="text-base-content/30 max-w-xs mx-auto mt-2 font-bold uppercase text-[10px] tracking-widest">Everything is under control</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      <div v-for="order in filteredOrders" :key="order.order_id" 
        class="card bg-base-100 shadow-xl border border-base-200 overflow-hidden transition-all duration-300"
        :class="{ 'ring-4 ring-primary ring-inset': order.status === STATUS_COOKING }"
      >
        <div class="p-8">
          <div class="flex justify-between items-start mb-6">
            <div class="space-y-1">
              <div class="flex items-center gap-2">
                <span class="text-3xl font-black tracking-tighter">#{{ order.order_number.split('-').pop() }}</span>
                <span v-if="order.status === STATUS_COOKING" class="badge badge-primary font-black text-[8px] h-4">IN WORK</span>
              </div>
              <p class="text-[9px] font-bold opacity-30 uppercase tracking-[0.2em] font-mono">{{ order.order_id.slice(0,12) }}</p>
            </div>
            <div class="badge font-black uppercase text-[10px] py-3.5 px-4 rounded-lg border-none" 
              :class="order.status === STATUS_COOKING ? 'badge-warning' : 'badge-info'">
              {{ order.status }}
            </div>
          </div>

          <div class="divider opacity-20 my-0"></div>
          
          <div class="py-8 space-y-3">
            <div v-for="item in order.items" :key="item.product_id" class="flex justify-between items-center bg-base-200/40 p-4 rounded-2xl border border-base-300/10">
              <span class="font-bold text-sm tracking-tight"><span class="text-primary font-black mr-3 text-lg">x{{ item.quantity }}</span> {{ item.product_name }}</span>
            </div>
          </div>

          <div class="card-actions">
            <button 
              v-if="order.status === STATUS_PAID"
              @click="openConfirm(order, STATUS_COOKING)"
              class="btn btn-primary btn-block rounded-2xl gap-3 font-black uppercase shadow-lg shadow-primary/10 h-16 text-lg"
            >
              <Play class="w-6 h-6" /> Start Cooking
            </button>
            <button 
              v-if="order.status === STATUS_COOKING"
              @click="openConfirm(order, STATUS_READY)"
              class="btn btn-success btn-block rounded-2xl gap-3 font-black uppercase shadow-lg shadow-success/10 h-16 text-lg text-white"
            >
              <CheckCircle2 class="w-6 h-6" /> Ready to Go
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
              <div class="bg-primary/10 w-24 h-24 rounded-full flex items-center justify-center mx-auto mb-8 animate-in zoom-in duration-300">
                <AlertCircle class="w-12 h-12 text-primary" />
              </div>
              <h3 class="font-black text-3xl uppercase tracking-tighter mb-3 leading-none">{{ pendingAction?.title }}</h3>
              <p class="text-base-content/50 text-sm mb-10 leading-relaxed font-medium px-4">
                {{ pendingAction?.description }}
              </p>
              <div class="flex flex-col gap-3">
                <button @click="handleConfirm" class="btn btn-primary btn-lg rounded-2xl h-16 font-black uppercase shadow-xl shadow-primary/20 tracking-tight">
                  Confirm
                </button>
                <button @click="showConfirmModal = false" class="btn btn-ghost btn-lg rounded-2xl font-bold opacity-40">
                  Cancel
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

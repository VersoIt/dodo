<script setup lang="ts">
import { ref, onMounted, inject, computed } from 'vue'
import { Tag, Plus, Trash2, TrendingUp, ShoppingBag, DollarSign, ChefHat, Truck, MessageSquare, Download } from 'lucide-vue-next'
import { ordersApi } from '../api'
import type { PromoCode, Analytics } from '../types'
import AppModal from '../components/shared/AppModal.vue'
import ChatComponent from '../components/ChatComponent.vue'

const addToast = inject('addToast') as (msg: string, type?: any) => void

const loading = ref(true)
const exporting = ref(false)
const promos = ref<PromoCode[]>([])
const analytics = ref<Analytics | null>(null)
const activeOrders = ref<any[]>([])

const showAddPromoModal = ref(false)
const showChatModal = ref(false)
const activeOrderId = ref<string | null>(null)
const newPromo = ref<{ code: string, amount: number, type: 'percent' | 'fixed' }>({ code: '', amount: 0, type: 'percent' })

const exportDates = ref({
  start: '',
  end: ''
})

const openChat = (orderId: string) => {
  activeOrderId.value = orderId
  showChatModal.value = true
}

const isPromoValid = computed(() => {
  return newPromo.value.code.trim().length > 0 && newPromo.value.amount > 0
})

const handleExport = async () => {
  try {
    exporting.value = true
    const res = await ordersApi.exportReport(exportDates.value.start, exportDates.value.end)
    if (res.success && res.data) {
      const url = window.URL.createObjectURL(new Blob([res.data]))
      const link = document.createElement('a')
      link.href = url
      const filename = `report_${new Date().toISOString().split('T')[0]}.xlsx`
      link.setAttribute('download', filename)
      document.body.appendChild(link)
      link.click()
      link.remove()
      addToast('Отчет загружен', 'success')
    }
  } catch (err) {
    addToast('Ошибка экспорта', 'error')
  } finally {
    exporting.value = false
  }
}

const fetchData = async () => {
  try {
    loading.value = true
    const [promosRes, analyticsRes, ordersRes] = await Promise.all([
      ordersApi.listPromos(),
      ordersApi.getAnalytics(),
      ordersApi.getAllOrders()
    ])
    if (promosRes.success && promosRes.data) promos.value = promosRes.data
    if (analyticsRes.success && analyticsRes.data) analytics.value = analyticsRes.data
    if (ordersRes.success && ordersRes.data) {
       activeOrders.value = ordersRes.data.filter((o: any) => o.status !== 'completed' && o.status !== 'canceled')
    }
  } catch (err) { console.error('Failed to load manager data') } finally { loading.value = false }
}

const handleAddPromo = async () => {
  if (!isPromoValid.value) {
    addToast('Заполните все поля корректно', 'error')
    return
  }
  try {
    const res = await ordersApi.createPromoCode(newPromo.value)
    if (res.success) {
      addToast('Промокод создан!', 'success'); showAddPromoModal.value = false
      newPromo.value = { code: '', amount: 0, type: 'percent' }; await fetchData()
    }
  } catch (err) { addToast('Ошибка создания', 'error') }
}

const handleDeletePromo = async (id: string) => {
  if (!confirm('Удалить промокод?')) return
  try { await ordersApi.deletePromo(id); await fetchData(); addToast('Удалено', 'success') } catch (err) { addToast('Ошибка', 'error') }
}

onMounted(fetchData)
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <div class="flex flex-col md:flex-row justify-between items-center mb-10 gap-6">
      <div>
        <h1 class="text-4xl font-black tracking-tighter uppercase">Панель Управления</h1>
        <p class="text-base-content/50 font-bold uppercase text-[10px] tracking-[0.2em] mt-1">Управление промокодами и аналитика</p>
      </div>

      <div class="flex flex-col sm:flex-row items-center gap-4 bg-base-100 p-4 rounded-[2rem] border border-base-200 shadow-sm">
        <div class="flex items-center gap-2">
          <input v-model="exportDates.start" type="date" class="input input-sm input-bordered rounded-xl font-bold text-xs" />
          <span class="opacity-30">—</span>
          <input v-model="exportDates.end" type="date" class="input input-sm input-bordered rounded-xl font-bold text-xs" />
        </div>
        <button @click="handleExport" :disabled="exporting" class="btn border-none text-primary-content bg-gradient-to-tr from-primary to-accent shadow-lg shadow-primary/30 hover:shadow-primary/50 hover:-translate-y-1 transition-all duration-300 ease-in-out rounded-xl font-black px-6 gap-2">
          <span v-if="exporting" class="loading loading-spinner loading-xs"></span>
          <Download v-else class="w-4 h-4" />
          Экспорт
        </button>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-32">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>

    <div v-else class="space-y-12 animate-in fade-in slide-in-from-right-4 duration-500">
      <!-- Analytics Dashboard -->
      <section v-if="analytics" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        <div class="card bg-primary text-primary-content shadow-xl rounded-[2rem]">
          <div class="card-body p-6">
            <div class="flex items-center gap-4">
              <div class="p-3 bg-white/20 rounded-2xl"><DollarSign class="w-6 h-6" /></div>
              <div>
                <p class="text-[10px] uppercase font-black tracking-widest opacity-70">Выручка</p>
                <h3 class="text-2xl font-black">{{ analytics.total_revenue.toLocaleString() }} ₽</h3>
              </div>
            </div>
          </div>
        </div>
        <div class="card bg-base-100 border border-base-200 shadow-xl rounded-[2rem]">
          <div class="card-body p-6">
            <div class="flex items-center gap-4">
              <div class="p-3 bg-secondary/10 text-secondary rounded-2xl"><ShoppingBag class="w-6 h-6" /></div>
              <div>
                <p class="text-[10px] uppercase font-black tracking-widest opacity-50 text-base-content">Заказов</p>
                <h3 class="text-2xl font-black text-base-content">{{ analytics.orders_count }}</h3>
              </div>
            </div>
          </div>
        </div>
        <div class="card bg-base-100 border border-base-200 shadow-xl rounded-[2rem]">
          <div class="card-body p-6">
            <div class="flex items-center gap-4">
              <div class="p-3 bg-accent/10 text-accent rounded-2xl"><TrendingUp class="w-6 h-6" /></div>
              <div>
                <p class="text-[10px] uppercase font-black tracking-widest opacity-50 text-base-content">Средний чек</p>
                <h3 class="text-2xl font-black text-base-content">{{ Math.round(analytics.avg_check).toLocaleString() }} ₽</h3>
              </div>
            </div>
          </div>
        </div>
        <div class="card bg-base-100 border border-base-200 shadow-xl rounded-[2rem]">
          <div class="card-body p-6">
            <div class="flex items-center gap-4">
              <div class="p-3 bg-warning/10 text-warning rounded-2xl"><ChefHat class="w-6 h-6" /></div>
              <div>
                <p class="text-[10px] uppercase font-black tracking-widest opacity-50 text-base-content">На кухне</p>
                <h3 class="text-2xl font-black text-base-content">{{ analytics.cooking_count }}</h3>
              </div>
            </div>
          </div>
        </div>
        <div class="card bg-base-100 border border-base-200 shadow-xl rounded-[2rem]">
          <div class="card-body p-6">
            <div class="flex items-center gap-4">
              <div class="p-3 bg-info/10 text-info rounded-2xl"><Truck class="w-6 h-6" /></div>
              <div>
                <p class="text-[10px] uppercase font-black tracking-widest opacity-50 text-base-content">В пути</p>
                <h3 class="text-2xl font-black text-base-content">{{ analytics.delivering_count }}</h3>
              </div>
            </div>
          </div>
        </div>
      </section>

      <div v-if="analytics && analytics.top_products.length > 0" class="card bg-base-100 border border-base-200 shadow-xl rounded-[2rem] overflow-hidden">
        <div class="p-8 pb-0">
          <h2 class="text-2xl font-black uppercase tracking-tight">Популярные товары</h2>
        </div>
        <div class="p-8 overflow-x-auto">
          <table class="table w-full">
            <thead>
              <tr class="text-[10px] uppercase font-black opacity-40 border-b-2 border-base-200">
                <th>Название</th>
                <th>Продано</th>
                <th>Выручка</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in analytics.top_products.sort((a,b) => b.revenue - a.revenue).slice(0, 5)" :key="p.name" class="border-base-200 font-bold">
                <td class="py-4">{{ p.name }}</td>
                <td class="py-4"><span class="badge badge-ghost font-black">{{ p.count }}</span></td>
                <td class="py-4">{{ p.revenue.toLocaleString() }} ₽</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Active Orders Section -->
      <div class="space-y-8">
        <h2 class="text-2xl font-black uppercase tracking-tight">Активные заказы (Чат)</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div v-for="order in activeOrders" :key="order.order_id" class="card bg-base-100 border border-base-200 shadow-sm hover:shadow-md transition-all rounded-[1.5rem] overflow-hidden">
            <div class="card-body p-5">
              <div class="flex justify-between items-start mb-2">
                <span class="font-black text-sm">#{{ order.order_number.split('-').pop() }}</span>
                <span class="badge badge-xs font-bold uppercase tracking-tighter">{{ order.status }}</span>
              </div>
              <p class="text-xs opacity-50 mb-4 line-clamp-1">{{ order.address?.street }}</p>
              <button @click="openChat(order.order_id)" class="btn btn-primary btn-sm btn-block rounded-xl gap-2 font-black">
                <MessageSquare class="w-4 h-4" /> Чат
              </button>
            </div>
          </div>
        </div>
        <div v-if="activeOrders.length === 0" class="text-center py-10 bg-base-100 rounded-[2rem] border-2 border-dashed border-base-300 opacity-40 font-black uppercase text-xs">Нет активных заказов</div>
      </div>

      <div class="space-y-8">
        <div class="flex justify-between items-center">
          <h2 class="text-2xl font-black uppercase tracking-tight">Активные промокоды</h2>
          <button @click="showAddPromoModal = true" class="btn btn-secondary rounded-2xl gap-2 font-black uppercase text-xs">
            <Plus class="w-4 h-4" /> Создать промокод
          </button>
        </div>
        
        <div v-if="promos.length === 0" class="text-center py-20 bg-base-100 rounded-[2rem] border-2 border-dashed border-base-300 opacity-40 font-black uppercase">Промокоды не созданы</div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div v-for="promo in promos" :key="promo.id" class="card bg-base-100 shadow-xl border border-base-200 rounded-[2rem] overflow-hidden group">
            <div class="card-body p-8">
              <div class="flex justify-between items-start">
                <div><span class="badge badge-secondary badge-lg font-black font-mono tracking-widest px-4 py-4 mb-3">{{ promo.code }}</span><p class="font-bold text-xs uppercase opacity-50">Скидка: <span class="text-base-content">{{ promo.amount }}{{ promo.type === 'percent' ? '%' : ' ₽' }}</span></p></div>
                <button @click="handleDeletePromo(promo.id)" class="btn btn-ghost btn-circle btn-sm text-error/30 hover:text-error hover:bg-error/10"><Trash2 class="w-4 h-4" /></button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <AppModal :show="showAddPromoModal" title="Новый промокод" @close="showAddPromoModal = false">
      <div class="p-10 space-y-6">
        <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-60 text-secondary">Код купона <span class="text-error">*</span></span></label><input v-model="newPromo.code" type="text" placeholder="SUMMER2026" class="input input-bordered border-secondary/30 bg-secondary/5 h-14 w-full rounded-2xl font-mono font-black uppercase focus:border-secondary" /></div>
        <div class="grid grid-cols-2 gap-6">
          <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Тип</span></label><select v-model="newPromo.type" class="select select-bordered w-full rounded-2xl h-14 font-bold"><option value="percent">Процент (%)</option><option value="fixed">Сумма (₽)</option></select></div>
          <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-60 text-secondary">Значение <span class="text-error">*</span></span></label><input v-model="newPromo.amount" type="number" class="input input-bordered border-secondary/30 bg-secondary/5 w-full rounded-2xl h-14 font-bold focus:border-secondary" /></div>
        </div>
        <button @click="handleAddPromo" :disabled="!isPromoValid" class="btn btn-secondary btn-block h-16 rounded-2xl font-black uppercase shadow-xl shadow-secondary/20 disabled:shadow-none mt-4">Создать промокод</button>
      </div>
    </AppModal>

    <AppModal :show="showChatModal" @close="showChatModal = false">
      <div class="p-4 h-[600px]">
        <ChatComponent v-if="activeOrderId" :orderId="activeOrderId" />
      </div>
    </AppModal>
  </div>
</template>

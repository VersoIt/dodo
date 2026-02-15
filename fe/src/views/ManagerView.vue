<script setup lang="ts">
import { ref, onMounted, inject } from 'vue'
import { BarChart3, Tag, TrendingUp, Users, DollarSign, Plus, Trash2 } from 'lucide-vue-next'
import { formatPrice } from '../utils/format'
import { ordersApi } from '../api'
import type { PromoCode, ProductStat, Analytics } from '../types'

interface Stats extends Analytics {}

const addToast = inject('addToast') as (msg: string, type?: any) => void

const activeTab = ref('kpi')
const loading = ref(true)
const stats = ref<Stats>({ total_revenue: 0, orders_count: 0, avg_check: 0, top_products: [] })
const promos = ref<PromoCode[]>([])

const showAddPromoModal = ref(false)
const newPromo = ref<{ code: string, amount: number, type: 'percent' | 'fixed' }>({ code: '', amount: 0, type: 'percent' })

const fetchData = async () => {
  try {
    loading.value = true
    const [statsRes, promosRes] = await Promise.all([
      ordersApi.getAnalytics(),
      ordersApi.listPromos()
    ])
    if (statsRes.success && statsRes.data) stats.value = statsRes.data
    if (promosRes.success && promosRes.data) promos.value = promosRes.data
  } catch (err) { console.error('Failed to load manager data') } finally { loading.value = false }
}

const handleAddPromo = async () => {
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
        <p class="text-base-content/50 font-bold uppercase text-[10px] tracking-[0.2em] mt-1">Реальные данные из системы</p>
      </div>
      <div class="tabs tabs-boxed bg-base-200 p-1.5 rounded-2xl">
        <a class="tab rounded-xl transition-all font-bold uppercase text-[10px] tracking-widest px-6 gap-2" :class="{ 'tab-active bg-primary text-primary-content shadow-lg': activeTab === 'kpi' }" @click="activeTab = 'kpi'"><BarChart3 class="w-4 h-4" /> Аналитика</a>
        <a class="tab rounded-xl transition-all font-bold uppercase text-[10px] tracking-widest px-6 gap-2" :class="{ 'tab-active bg-secondary text-secondary-content shadow-lg': activeTab === 'promos' }" @click="activeTab = 'promos'"><Tag class="w-4 h-4" /> Промокоды</a>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-32"><span class="loading loading-spinner loading-lg text-primary"></span></div>

    <div v-else-if="activeTab === 'kpi'" class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div class="stat bg-base-100 shadow-xl rounded-[2rem] border border-base-200">
          <div class="stat-figure text-primary bg-primary/10 p-3 rounded-xl"><DollarSign class="w-6 h-6" /></div>
          <div class="stat-title font-bold uppercase text-[10px] tracking-widest">Общая выручка</div>
          <div class="stat-value text-primary text-3xl font-black">{{ formatPrice(stats.total_revenue) }}</div>
          <div class="stat-desc font-medium text-success flex items-center gap-1"><TrendingUp class="w-3 h-3" /> В режиме реального времени</div>
        </div>
        <div class="stat bg-base-100 shadow-xl rounded-[2rem] border border-base-200">
          <div class="stat-figure text-secondary bg-secondary/10 p-3 rounded-xl"><Tag class="w-6 h-6" /></div>
          <div class="stat-title font-bold uppercase text-[10px] tracking-widest">Всего заказов</div>
          <div class="stat-value text-secondary text-3xl font-black">{{ stats.orders_count }}</div>
          <div class="stat-desc font-medium">Без учета отмененных</div>
        </div>
        <div class="stat bg-base-100 shadow-xl rounded-[2rem] border border-base-200">
          <div class="stat-figure text-accent bg-accent/10 p-3 rounded-xl"><BarChart3 class="w-6 h-6" /></div>
          <div class="stat-title font-bold uppercase text-[10px] tracking-widest">Средний чек</div>
          <div class="stat-value text-accent text-3xl font-black">{{ formatPrice(stats.avg_check) }}</div>
          <div class="stat-desc font-medium">Рассчитано по всем заказам</div>
        </div>
      </div>
      <div class="card bg-base-100 shadow-xl border border-base-200 rounded-[2.5rem]">
        <div class="card-body p-10">
          <h2 class="card-title font-black uppercase tracking-tight mb-8">Топ-5 товаров по продажам</h2>
          <div v-if="stats.top_products && stats.top_products.length > 0" class="overflow-x-auto">
            <table class="table">
              <thead><tr class="uppercase text-[10px] tracking-widest text-base-content/40 border-b-2 border-base-200"><th>Название</th><th>Продано шт.</th><th class="text-right">Выручка</th></tr></thead>
              <tbody><tr v-for="(p, i) in stats.top_products" :key="i" class="hover:bg-base-200/50 transition-colors font-bold"><td class="text-lg">{{ p.name }}</td><td class="text-secondary font-black text-xl">{{ p.count }}</td><td class="text-right">{{ formatPrice(p.revenue) }}</td></tr></tbody>
            </table>
          </div>
          <div v-else class="text-center py-10 opacity-30 font-bold uppercase">Нет данных о продажах</div>
        </div>
      </div>
    </div>

    <div v-else class="space-y-8 animate-in fade-in slide-in-from-right-4 duration-500">
      <div class="flex justify-between items-center"><h2 class="text-2xl font-black uppercase tracking-tight">Активные промокоды</h2><button @click="showAddPromoModal = true" class="btn btn-secondary rounded-2xl gap-2 font-black uppercase text-xs"><Plus class="w-4 h-4" /> Создать промокод</button></div>
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

    <AppModal :show="showAddPromoModal" title="Новый промокод" @close="showAddPromoModal = false">
      <div class="p-10 space-y-6">
        <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Код купона</span></label><input v-model="newPromo.code" type="text" placeholder="SUMMER2026" class="input input-bordered w-full rounded-2xl font-mono font-black h-14 uppercase" /></div>
        <div class="grid grid-cols-2 gap-6">
          <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Тип</span></label><select v-model="newPromo.type" class="select select-bordered w-full rounded-2xl h-14 font-bold"><option value="percent">Процент (%)</option><option value="fixed">Сумма (₽)</option></select></div>
          <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Значение</span></label><input v-model="newPromo.amount" type="number" class="input input-bordered w-full rounded-2xl h-14 font-bold" /></div>
        </div>
        <button @click="handleAddPromo" class="btn btn-secondary btn-block h-16 rounded-2xl font-black uppercase shadow-xl shadow-secondary/20 mt-4">Создать промокод</button>
      </div>
    </AppModal>
  </div>
</template>

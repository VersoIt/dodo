<script setup lang="ts">
import { ref, computed } from 'vue'
import { BarChart3, Tag, TrendingUp, Users, DollarSign, Plus, Trash2, Calendar } from 'lucide-vue-next'
import { formatPrice } from '../utils/format'
import AppModal from '../components/shared/AppModal.vue'

// --- TABS STATE ---
const activeTab = ref('kpi') // 'kpi' | 'promos'

// --- MOCK DATA FOR KPI ---
const stats = ref({
  totalRevenue: 1250500,
  ordersCount: 1450,
  avgCheck: 860,
  activeUsers: 320
})

const topProducts = ref([
  { name: 'Пепперони Экстрим', count: 450, revenue: 269550 },
  { name: 'Четыре Сыра', count: 320, revenue: 191680 },
  { name: 'Маргарита', count: 210, revenue: 94290 },
])

// --- PROMO CODES STATE ---
const promos = ref([
  { id: 1, code: 'WELCOME2026', discount: 15, type: 'percent', active: true },
  { id: 2, code: 'PIZZADAY', discount: 200, type: 'fixed', active: true },
  { id: 3, code: 'FREEDELIVERY', discount: 0, type: 'free_delivery', active: false },
])

const showAddPromoModal = ref(false)
const newPromo = ref({ code: '', discount: 0, type: 'percent' })

const addPromo = () => {
  promos.value.push({
    id: Date.now(),
    code: newPromo.value.code.toUpperCase(),
    discount: newPromo.value.discount,
    type: newPromo.value.type,
    active: true
  })
  showAddPromoModal.value = false
  newPromo.value = { code: '', discount: 0, type: 'percent' }
}

const removePromo = (id: number) => {
  promos.value = promos.value.filter(p => p.id !== id)
}

const togglePromo = (id: number) => {
  const p = promos.value.find(p => p.id === id)
  if (p) p.active = !p.active
}
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <!-- Header -->
    <div class="flex flex-col md:flex-row justify-between items-center mb-10 gap-6">
      <div>
        <h1 class="text-4xl font-black tracking-tighter uppercase italic text-base-content">Панель Управления</h1>
        <p class="text-base-content/50 font-bold uppercase text-[10px] tracking-[0.2em] mt-1">Analytics & Marketing</p>
      </div>
      
      <!-- Tabs -->
      <div class="tabs tabs-boxed bg-base-200 p-1.5 rounded-2xl">
        <a 
          class="tab rounded-xl transition-all font-bold uppercase text-[10px] tracking-widest px-6 gap-2"
          :class="{ 'tab-active bg-primary text-primary-content shadow-lg': activeTab === 'kpi' }"
          @click="activeTab = 'kpi'"
        >
          <BarChart3 class="w-4 h-4" /> Аналитика (KPI)
        </a>
        <a 
          class="tab rounded-xl transition-all font-bold uppercase text-[10px] tracking-widest px-6 gap-2"
          :class="{ 'tab-active bg-secondary text-secondary-content shadow-lg': activeTab === 'promos' }"
          @click="activeTab = 'promos'"
        >
          <Tag class="w-4 h-4" /> Промокоды
        </a>
      </div>
    </div>

    <!-- KPI TAB -->
    <div v-if="activeTab === 'kpi'" class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <!-- Stats Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="stat bg-base-100 shadow-xl rounded-[2rem] border border-base-200">
          <div class="stat-figure text-primary bg-primary/10 p-3 rounded-xl"><DollarSign class="w-6 h-6" /></div>
          <div class="stat-title font-bold uppercase text-[10px] tracking-widest">Выручка (Месяц)</div>
          <div class="stat-value text-primary text-3xl font-black">{{ formatPrice(stats.totalRevenue) }}</div>
          <div class="stat-desc font-medium text-success flex items-center gap-1"><TrendingUp class="w-3 h-3" /> +12% к прошлому</div>
        </div>
        
        <div class="stat bg-base-100 shadow-xl rounded-[2rem] border border-base-200">
          <div class="stat-figure text-secondary bg-secondary/10 p-3 rounded-xl"><Tag class="w-6 h-6" /></div>
          <div class="stat-title font-bold uppercase text-[10px] tracking-widest">Заказов</div>
          <div class="stat-value text-secondary text-3xl font-black">{{ stats.ordersCount }}</div>
          <div class="stat-desc font-medium">~48 заказов в день</div>
        </div>

        <div class="stat bg-base-100 shadow-xl rounded-[2rem] border border-base-200">
          <div class="stat-figure text-accent bg-accent/10 p-3 rounded-xl"><BarChart3 class="w-6 h-6" /></div>
          <div class="stat-title font-bold uppercase text-[10px] tracking-widest">Средний чек</div>
          <div class="stat-value text-accent text-3xl font-black">{{ formatPrice(stats.avgCheck) }}</div>
          <div class="stat-desc font-medium">Стабильный рост</div>
        </div>

        <div class="stat bg-base-100 shadow-xl rounded-[2rem] border border-base-200">
          <div class="stat-figure text-neutral bg-neutral/10 p-3 rounded-xl"><Users class="w-6 h-6" /></div>
          <div class="stat-title font-bold uppercase text-[10px] tracking-widest">Клиентов</div>
          <div class="stat-value text-3xl font-black">{{ stats.activeUsers }}</div>
          <div class="stat-desc font-medium">+24 новых за неделю</div>
        </div>
      </div>

      <!-- Top Products -->
      <div class="card bg-base-100 shadow-xl border border-base-200 rounded-[2.5rem]">
        <div class="card-body p-8">
          <h2 class="card-title font-black uppercase tracking-tight mb-6">Топ продаж</h2>
          <div class="overflow-x-auto">
            <table class="table">
              <thead>
                <tr class="uppercase text-[10px] tracking-widest text-base-content/40 border-b-2 border-base-200">
                  <th>Название</th>
                  <th>Продано шт.</th>
                  <th class="text-right">Выручка</th>
                  <th class="text-right">Доля</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(product, i) in topProducts" :key="i" class="hover:bg-base-200/50 transition-colors font-bold">
                  <td class="text-lg">{{ product.name }}</td>
                  <td class="text-secondary font-black">{{ product.count }}</td>
                  <td class="text-right">{{ formatPrice(product.revenue) }}</td>
                  <td class="text-right">
                    <progress class="progress progress-primary w-20" :value="product.count" max="500"></progress>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- PROMOS TAB -->
    <div v-else class="space-y-8 animate-in fade-in slide-in-from-right-4 duration-500">
      <div class="flex justify-between items-center">
        <h2 class="text-2xl font-black uppercase tracking-tight">Активные акции</h2>
        <button @click="showAddPromoModal = true" class="btn btn-secondary rounded-2xl gap-2 font-black uppercase text-xs shadow-lg shadow-secondary/20">
          <Plus class="w-4 h-4" /> Создать промокод
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div v-for="promo in promos" :key="promo.id" class="card bg-base-100 shadow-xl border border-base-200 rounded-[2rem] group hover:border-secondary/30 transition-all">
          <div class="card-body flex flex-row items-center justify-between p-8">
            <div>
              <div class="flex items-center gap-3 mb-2">
                <span class="badge badge-lg font-black font-mono tracking-widest px-4 py-4" :class="promo.active ? 'badge-secondary' : 'badge-ghost opacity-50'">
                  {{ promo.code }}
                </span>
                <span v-if="!promo.active" class="text-[10px] font-bold uppercase text-error">Отключен</span>
              </div>
              <p class="text-base-content/60 font-bold text-xs uppercase tracking-wider">
                Скидка: 
                <span class="text-base-content font-black">
                  {{ promo.type === 'percent' ? `-${promo.discount}%` : `-${promo.discount} ₽` }}
                </span>
              </p>
            </div>
            
            <div class="flex items-center gap-2">
              <input type="checkbox" class="toggle toggle-secondary" :checked="promo.active" @change="togglePromo(promo.id)" />
              <button @click="removePromo(promo.id)" class="btn btn-ghost btn-circle btn-sm hover:bg-error/10 hover:text-error">
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ADD PROMO MODAL -->
    <AppModal :show="showAddPromoModal" title="Новый промокод" @close="showAddPromoModal = false">
      <div class="p-8 space-y-6">
        <div class="form-control">
          <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Код купона</span></label>
          <input v-model="newPromo.code" type="text" placeholder="SUMMER2026" class="input input-bordered w-full rounded-xl font-mono font-black uppercase tracking-widest" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="form-control">
            <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Тип скидки</span></label>
            <select v-model="newPromo.type" class="select select-bordered w-full rounded-xl font-bold">
              <option value="percent">Процент (%)</option>
              <option value="fixed">Сумма (₽)</option>
            </select>
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Значение</span></label>
            <input v-model="newPromo.discount" type="number" class="input input-bordered w-full rounded-xl font-bold" />
          </div>
        </div>
        <button @click="addPromo" class="btn btn-secondary btn-block h-14 rounded-2xl font-black uppercase shadow-xl mt-4">
          Создать
        </button>
      </div>
    </AppModal>
  </div>
</template>

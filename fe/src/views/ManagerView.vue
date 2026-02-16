<script setup lang="ts">
import { ref, onMounted, inject } from 'vue'
import { Tag, Plus, Trash2 } from 'lucide-vue-next'
import { ordersApi } from '../api'
import type { PromoCode } from '../types'

const addToast = inject('addToast') as (msg: string, type?: any) => void

const loading = ref(true)
const promos = ref<PromoCode[]>([])

const showAddPromoModal = ref(false)
const newPromo = ref<{ code: string, amount: number, type: 'percent' | 'fixed' }>({ code: '', amount: 0, type: 'percent' })

const fetchData = async () => {
  try {
    loading.value = true
    const promosRes = await ordersApi.listPromos()
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
        <p class="text-base-content/50 font-bold uppercase text-[10px] tracking-[0.2em] mt-1">Управление промокодами</p>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-32"><span class="loading loading-spinner loading-lg text-primary"></span></div>

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
        <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Код купона</span></label><input v-model="newPromo.code" type="text" placeholder="SUMMER2026" class="input input-bordered h-14 w-full rounded-2xl font-mono font-black uppercase" /></div>
        <div class="grid grid-cols-2 gap-6">
          <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Тип</span></label><select v-model="newPromo.type" class="select select-bordered w-full rounded-2xl h-14 font-bold"><option value="percent">Процент (%)</option><option value="fixed">Сумма (₽)</option></select></div>
          <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Значение</span></label><input v-model="newPromo.amount" type="number" class="input input-bordered w-full rounded-2xl h-14 font-bold" /></div>
        </div>
        <button @click="handleAddPromo" class="btn btn-secondary btn-block h-16 rounded-2xl font-black uppercase shadow-xl shadow-secondary/20 mt-4">Создать промокод</button>
      </div>
    </AppModal>
  </div>
</template>

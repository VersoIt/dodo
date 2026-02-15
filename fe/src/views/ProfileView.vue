<script setup lang="ts">
import { ref, onMounted, inject } from 'vue'
import { useAuthStore } from '../store/auth'
import { Mail, Shield, Calendar, User, Package, Save, X, Phone } from 'lucide-vue-next'
import axios from 'axios'

const authStore = useAuthStore()
const addToast = inject('addToast') as (msg: string, type?: any) => void

const activeTab = ref('profile')
const isEditing = ref(false)
const orders = ref<any[]>([])
const loadingOrders = ref(false)

// Edit Form
const editName = ref('')
const editPhone = ref('')
const isSaving = ref(false)

const startEditing = () => {
  editName.value = authStore.user?.name || ''
  editPhone.value = authStore.user?.phone || ''
  isEditing.value = true
}

const saveProfile = async () => {
  try {
    isSaving.value = true
    const response = await axios.patch('/api/v1/auth/me', {
      name: editName.value,
      phone: editPhone.value
    })
    if (response.data.success) {
      authStore.user = response.data.data
      isEditing.value = false
      addToast('Профиль обновлен!', 'success')
    }
  } catch (err) {
    console.error('Failed to update profile:', err)
    addToast('Не удалось обновить профиль', 'error')
  } finally {
    isSaving.value = false
  }
}

const fetchOrders = async () => {
  try {
    loadingOrders.value = true
    const response = await axios.get('/api/v1/orders/my')
    if (response.data.success) {
      orders.value = response.data.data || []
    }
  } catch (err) {
    console.error('Failed to fetch orders:', err)
  } finally {
    loadingOrders.value = false
  }
}

onMounted(async () => {
  await authStore.fetchMe()
})
</script>

<template>
  <div class="max-w-4xl mx-auto py-8 px-4">
    <div class="flex flex-col md:flex-row gap-8">
      <!-- Sidebar -->
      <div class="w-full md:w-64 space-y-2">
        <button 
          @click="activeTab = 'profile'"
          class="btn btn-ghost w-full justify-start gap-3"
          :class="{ 'btn-active bg-primary/10 text-primary': activeTab === 'profile' }"
        >
          <User class="w-5 h-5" />
          Информация
        </button>
        <button 
          v-if="authStore.user?.role === 'client'"
          @click="activeTab = 'orders'; fetchOrders()"
          class="btn btn-ghost w-full justify-start gap-3"
          :class="{ 'btn-active bg-primary/10 text-primary': activeTab === 'orders' }"
        >
          <Package class="w-5 h-5" />
          Мои заказы
        </button>
        <div class="divider"></div>
        <button 
          @click="authStore.logout(); $router.push('/login')"
          class="btn btn-ghost w-full justify-start gap-3 text-error hover:bg-error/10"
        >
          Выйти
        </button>
      </div>

      <!-- Main Content -->
      <div class="flex-1">
        <div v-if="activeTab === 'profile'" class="card bg-base-100 shadow-xl border border-base-200 overflow-hidden transition-all duration-300">
          <div class="card-body">
            <div class="flex items-center justify-between mb-8">
              <div class="flex items-center gap-6">
                <div class="avatar placeholder">
                  <div class="bg-primary text-primary-content rounded-2xl w-20 shadow-lg">
                    <span class="text-3xl font-bold uppercase">{{ authStore.user?.name?.[0] || authStore.user?.email?.[0] || 'U' }}</span>
                  </div>
                </div>
                <div>
                  <h1 class="text-3xl font-bold tracking-tight">{{ authStore.user?.name || 'Любитель пиццы' }}</h1>
                  <p class="text-base-content/60">{{ authStore.user?.email }}</p>
                </div>
              </div>
              <button v-if="!isEditing" @click="startEditing" class="btn btn-outline btn-sm rounded-lg">Изменить</button>
            </div>

            <div v-if="!isEditing" class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="p-4 bg-base-200/50 border border-base-300/50 rounded-2xl space-y-1">
                <p class="text-[10px] font-black uppercase tracking-widest text-base-content/40 flex items-center gap-2">
                  <User class="w-3 h-3" /> Имя
                </p>
                <p class="font-bold text-lg">{{ authStore.user?.name || 'Не указано' }}</p>
              </div>
              <div class="p-4 bg-base-200/50 border border-base-300/50 rounded-2xl space-y-1">
                <p class="text-[10px] font-black uppercase tracking-widest text-base-content/40 flex items-center gap-2">
                  <Phone class="w-3 h-3" /> Телефон
                </p>
                <p class="font-bold text-lg">{{ authStore.user?.phone || 'Не указано' }}</p>
              </div>
              <div class="p-4 bg-base-200/50 border border-base-300/50 rounded-2xl space-y-1">
                <p class="text-[10px] font-black uppercase tracking-widest text-base-content/40 flex items-center gap-2">
                  <Mail class="w-3 h-3" /> Email
                </p>
                <p class="font-bold text-lg">{{ authStore.user?.email }}</p>
              </div>
              <div class="p-4 bg-base-200/50 border border-base-300/50 rounded-2xl space-y-1">
                <p class="text-[10px] font-black uppercase tracking-widest text-base-content/40 flex items-center gap-2">
                  <Shield class="w-3 h-3" /> Роль
                </p>
                <div class="badge badge-secondary font-bold uppercase text-[10px]">{{ authStore.user?.role || 'client' }}</div>
              </div>
            </div>

            <!-- Edit Mode -->
            <div v-else class="space-y-4 animate-in fade-in slide-in-from-bottom-2 duration-300">
              <div class="form-control w-full">
                <label class="label"><span class="label-text font-bold">Полное имя</span></label>
                <input v-model="editName" type="text" class="input input-bordered w-full rounded-xl" placeholder="Иван Иванов" />
              </div>
              <div class="form-control w-full">
                <label class="label"><span class="label-text font-bold">Номер телефона</span></label>
                <input v-model="editPhone" type="text" class="input input-bordered w-full rounded-xl" placeholder="+7 999 000 00 00" />
              </div>
              <div class="flex gap-2 justify-end mt-6">
                <button @click="isEditing = false" class="btn btn-ghost rounded-lg gap-2">
                  <X class="w-4 h-4" /> Отмена
                </button>
                <button @click="saveProfile" class="btn btn-primary rounded-lg gap-2" :disabled="isSaving">
                  <span v-if="isSaving" class="loading loading-spinner"></span>
                  <Save class="w-4 h-4" /> Сохранить
                </button>
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'orders'" class="space-y-4 animate-in fade-in slide-in-from-right-4 duration-300">
          <h2 class="text-2xl font-bold mb-4">История заказов</h2>
          
          <div v-if="loadingOrders" class="flex justify-center py-12">
            <span class="loading loading-spinner loading-lg text-primary"></span>
          </div>

          <div v-else-if="orders.length === 0" class="card bg-base-100 border border-dashed border-base-300 py-12">
            <div class="flex flex-col items-center text-center">
              <Package class="w-16 h-16 text-base-content/20 mb-4" />
              <h3 class="text-xl font-bold">Заказов пока нет</h3>
              <p class="text-base-content/60 mb-6">Проголодались? Пора заказать вкусную пиццу!</p>
              <router-link to="/" class="btn btn-primary">В меню</router-link>
            </div>
          </div>

          <div v-else v-for="order in orders" :key="order.order_id" class="card bg-base-100 shadow-md border border-base-200 hover:border-primary/30 transition-all hover:shadow-lg">
            <div class="card-body p-6 flex flex-row justify-between items-center">
              <div>
                <div class="flex items-center gap-3 mb-1">
                  <span class="font-mono font-bold text-lg tracking-tight">#{{ order.order_number }}</span>
                  <div class="badge badge-outline uppercase text-[10px] font-black tracking-widest" 
                    :class="{ 
                      'badge-success border-green-500 text-green-500': order.status === 'completed' || order.status === 'paid',
                      'badge-warning border-yellow-500 text-yellow-500': order.status === 'cooking' || order.status === 'created'
                    }">
                    {{ order.status }}
                  </div>
                </div>
                <p class="text-[10px] uppercase font-bold text-base-content/40">15 Февраля, 2026</p>
              </div>
              <div class="text-right">
                <p class="text-2xl font-black text-primary">{{ order.final_price?.toLocaleString() }} ₽</p>
                <router-link :to="'/order/' + order.order_id" class="btn btn-link btn-xs text-primary no-underline hover:underline p-0 h-auto min-h-0">Подробнее</router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

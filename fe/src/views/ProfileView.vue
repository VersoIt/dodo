<script setup lang="ts">
import { computed, inject, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { ordersApi } from '../api'
import { formatPrice } from '../utils/format'
import {
  BarChart3,
  BookOpen,
  Calendar,
  ChefHat,
  CircleDot,
  ClipboardList,
  IdCard,
  LayoutDashboard,
  LogOut,
  Mail,
  MessageSquare,
  Package,
  PencilLine,
  Phone,
  Save,
  Shield,
  Truck,
  User,
  X
} from 'lucide-vue-next'

type ToastType = 'success' | 'error' | 'info'
type StaffRole = 'manager' | 'chef' | 'courier'
type StaffStatus = 'Активен' | 'На смене' | 'Неактивен'
type ClientTab = 'profile' | 'orders'

interface StaffNavItem {
  id: string
  label: string
  icon: typeof User
  route?: string
  toast?: string
}

interface StaffProfileState {
  userId: string
  userIdHint: string
  staffId?: string
  staffIdHint?: string
  staffIdLabel?: string
  fullName: string
  email: string
  phone?: string
  roleLabel: string
  status?: StaffStatus
  statusLabel?: string
  registeredAt: string
}

interface StaffInfoCard {
  key: string
  label: string
  hint: string
  value: string
  icon: typeof User
  tone?: 'role' | 'status'
}

const router = useRouter()
const authStore = useAuthStore()
const addToast = inject<(message: string, type?: ToastType) => void>('addToast', () => undefined)

const clientTab = ref<ClientTab>('profile')
const clientEditing = ref(false)
const staffEditing = ref(false)
const orders = ref<any[]>([])
const loadingOrders = ref(false)
const isSaving = ref(false)

const clientDraft = ref({
  fullName: '',
  phone: ''
})

const staffDraft = ref({
  fullName: '',
  email: '',
  phone: '',
  status: 'Активен' as StaffStatus
})

const roleLabels = {
  client: 'Клиент',
  manager: 'Менеджер',
  chef: 'Повар',
  courier: 'Курьер'
} as const

const staffProfiles = ref<Record<StaffRole, StaffProfileState>>({
  manager: {
    userId: '1001',
    userIdHint: 'users.user_id',
    fullName: 'Руслан Лященко',
    email: 'manager@manager.com',
    roleLabel: 'Менеджер',
    registeredAt: '12.03.2024, 14:25'
  },
  chef: {
    userId: 'U-1024',
    userIdHint: 'users.user_id',
    staffId: 'CH-2048',
    staffIdHint: 'chefs.chef_id',
    staffIdLabel: 'ID повара',
    fullName: 'Алексей Петров',
    email: 'alexey.petrov@pizzagood.ru',
    phone: '+7 (999) 123-45-67',
    roleLabel: 'Повар',
    status: 'Активен',
    statusLabel: 'Статус повара',
    registeredAt: '15.03.2024'
  },
  courier: {
    userId: '1042',
    userIdHint: 'users.user_id',
    staffId: '238',
    staffIdHint: 'couriers.courier_id',
    staffIdLabel: 'ID курьера',
    fullName: 'Иван Иванов',
    email: 'ivan@pizzagood.ru',
    phone: '+7 (999) 999-99-99',
    roleLabel: 'Курьер',
    status: 'Активен',
    statusLabel: 'Статус курьера',
    registeredAt: '2024-09-18 14:32'
  }
})

const currentRole = computed(() => authStore.user?.role || 'client')
const isClient = computed(() => currentRole.value === 'client')
const currentStaffRole = computed<StaffRole | null>(() => (
  currentRole.value === 'manager' || currentRole.value === 'chef' || currentRole.value === 'courier'
    ? currentRole.value
    : null
))

const roleLabel = computed(() => roleLabels[currentRole.value as keyof typeof roleLabels] || 'Клиент')
const displayName = computed(() => authStore.user?.name || 'Пользователь')
const displayEmail = computed(() => authStore.user?.email || 'user@pizzagood.ru')
const displayPhone = computed(() => authStore.user?.phone || '')
const displayShortName = computed(() => displayName.value.trim().split(/\s+/)[0] || 'Пользователь')
const displayInitial = computed(() => displayName.value.trim().charAt(0).toUpperCase() || 'U')

const currentStaffProfile = computed<StaffProfileState | null>(() => {
  if (!currentStaffRole.value) return null

  const base = staffProfiles.value[currentStaffRole.value]

  return {
    ...base,
    fullName: authStore.user?.name || base.fullName,
    email: authStore.user?.email || base.email,
    phone: authStore.user?.phone || base.phone
  }
})

const staffNavItems = computed<StaffNavItem[]>(() => {
  switch (currentStaffRole.value) {
    case 'manager':
      return [
        { id: 'info', label: 'Информация', icon: User },
        { id: 'dashboard', label: 'Дашборд', icon: LayoutDashboard, route: '/manager' },
        { id: 'orders', label: 'Заказы', icon: ClipboardList, toast: 'Раздел заказов пока оставил как фронтовой мок.' },
        { id: 'directories', label: 'Справочники', icon: BookOpen, route: '/' },
        { id: 'reports', label: 'Отчеты', icon: BarChart3, toast: 'Отчеты сейчас собраны на панели менеджера.' }
      ]
    case 'chef':
      return [
        { id: 'info', label: 'Информация', icon: User },
        { id: 'kitchen', label: 'Заказы кухни', icon: ChefHat, route: '/kitchen' },
        { id: 'chat', label: 'Чат', icon: MessageSquare, toast: 'Служебный чат оставил локальным моковым сценарием.' }
      ]
    case 'courier':
      return [
        { id: 'info', label: 'Информация', icon: User },
        { id: 'deliveries', label: 'Доставки', icon: Truck, route: '/logistics' },
        { id: 'chat', label: 'Чат', icon: MessageSquare, toast: 'Служебный чат остается локальным фронтовым блоком.' }
      ]
    default:
      return []
  }
})

const clientInfoCards = computed(() => [
  { label: 'Имя', value: displayName.value || 'Не указано', icon: User },
  { label: 'Телефон', value: displayPhone.value || 'Не указано', icon: Phone },
  { label: 'Email', value: displayEmail.value || 'Не указано', icon: Mail },
  { label: 'Роль', value: roleLabel.value, icon: Shield }
])

const staffInfoCards = computed<StaffInfoCard[]>(() => {
  const profile = currentStaffProfile.value
  if (!profile) return []

  const cards: StaffInfoCard[] = [
    {
      key: 'user-id',
      label: 'ID пользователя',
      hint: profile.userIdHint,
      value: profile.userId,
      icon: IdCard
    }
  ]

  if (profile.staffId && profile.staffIdHint && profile.staffIdLabel) {
    cards.push({
      key: 'staff-id',
      label: profile.staffIdLabel,
      hint: profile.staffIdHint,
      value: profile.staffId,
      icon: IdCard
    })
  }

  cards.push(
    {
      key: 'full-name',
      label: 'ФИО',
      hint: 'users.name',
      value: profile.fullName,
      icon: User
    },
    {
      key: 'email',
      label: 'Email',
      hint: 'users.email',
      value: profile.email,
      icon: Mail
    }
  )

  if (profile.phone) {
    cards.push({
      key: 'phone',
      label: 'Телефон',
      hint: 'users.phone',
      value: profile.phone,
      icon: Phone
    })
  }

  cards.push({
    key: 'role',
    label: 'Роль',
    hint: 'users.role',
    value: profile.roleLabel,
    icon: Shield,
    tone: 'role'
  })

  if (profile.status && profile.statusLabel) {
    cards.push({
      key: 'status',
      label: profile.statusLabel,
      hint: 'users.status',
      value: profile.status,
      icon: CircleDot,
      tone: 'status'
    })
  }

  cards.push({
    key: 'registered-at',
    label: 'Дата регистрации',
    hint: 'users.created_at',
    value: profile.registeredAt,
    icon: Calendar
  })

  return cards
})

const roleBadgeClass = computed(() => {
  switch (currentStaffRole.value) {
    case 'manager':
      return 'bg-primary/10 text-primary border border-primary/15'
    case 'chef':
      return 'bg-primary/10 text-primary border border-primary/15'
    case 'courier':
      return 'bg-slate-900 text-white border border-slate-900/10'
    default:
      return 'bg-base-200 text-base-content/70 border border-base-300'
  }
})

const statusBadgeClass = computed(() => {
  const status = currentStaffProfile.value?.status

  switch (status) {
    case 'Активен':
      return 'bg-green-500/10 text-green-600 border border-green-500/15'
    case 'На смене':
      return 'bg-amber-500/10 text-amber-600 border border-amber-500/15'
    default:
      return 'bg-base-200 text-base-content/55 border border-base-300'
  }
})

const startClientEditing = () => {
  clientEditing.value = true
  clientDraft.value = {
    fullName: displayName.value,
    phone: displayPhone.value
  }
}

const startStaffEditing = () => {
  const profile = currentStaffProfile.value
  if (!profile) return

  staffDraft.value = {
    fullName: profile.fullName,
    email: profile.email,
    phone: profile.phone || '',
    status: profile.status || 'Активен'
  }
  staffEditing.value = true
}

const stopEditing = () => {
  staffEditing.value = false
}

const saveClientProfile = async () => {
  if (!clientDraft.value.fullName.trim()) {
    addToast('Введите имя для профиля.', 'error')
    return
  }

  try {
    isSaving.value = true
    authStore.setUser({
      ...authStore.user,
      name: clientDraft.value.fullName.trim(),
      phone: clientDraft.value.phone.trim()
    })
    clientEditing.value = false
    addToast('Профиль обновлен локально.', 'success')
  } finally {
    isSaving.value = false
  }
}

const saveStaffProfile = async () => {
  if (!currentStaffRole.value || !currentStaffProfile.value) return
  if (!staffDraft.value.fullName.trim() || !staffDraft.value.email.trim()) {
    addToast('Заполните имя и email.', 'error')
    return
  }

  try {
    isSaving.value = true

    const current = staffProfiles.value[currentStaffRole.value]

    staffProfiles.value[currentStaffRole.value] = {
      ...current,
      fullName: staffDraft.value.fullName.trim(),
      email: staffDraft.value.email.trim(),
      phone: current.phone !== undefined ? staffDraft.value.phone.trim() : current.phone,
      status: current.status ? staffDraft.value.status : current.status
    }

    authStore.setUser({
      ...authStore.user,
      name: staffDraft.value.fullName.trim(),
      email: staffDraft.value.email.trim(),
      phone: staffDraft.value.phone.trim()
    })

    staffEditing.value = false
    addToast('Карточка сотрудника обновлена локально.', 'success')
  } finally {
    isSaving.value = false
  }
}

const fetchOrders = async () => {
  try {
    loadingOrders.value = true
    const res = await ordersApi.getMyOrders()
    if (res.success) orders.value = res.data || []
  } catch (err) {
    console.error('Failed to fetch orders:', err)
  } finally {
    loadingOrders.value = false
  }
}

const handleStaffNav = (item: StaffNavItem) => {
  if (item.id === 'info') return
  if (item.route) {
    router.push(item.route)
    return
  }
  addToast(item.toast || 'Этот раздел пока оставил моковым.', 'info')
}

const logout = () => {
  authStore.logout()
  router.push('/login')
}

const formatOrderDate = (value?: string) => {
  if (!value) return 'Недавно'

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value

  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: 'long',
    year: 'numeric'
  }).format(date)
}

onMounted(async () => {
  await authStore.fetchMe()
})
</script>

<template>
  <div class="pb-12 pt-6">
    <div v-if="isClient" class="mx-auto max-w-4xl px-4">
      <div class="grid gap-8 md:grid-cols-[250px_minmax(0,1fr)]">
        <aside class="rounded-[2rem] border border-white/50 bg-base-100/85 p-4 shadow-[0_20px_60px_-40px_rgba(15,23,42,0.45)] backdrop-blur">
          <button
            class="sidebar-link"
            :class="{ 'sidebar-link-active': clientTab === 'profile' }"
            @click="clientTab = 'profile'"
          >
            <User class="h-5 w-5" />
            Информация
          </button>
          <button
            class="sidebar-link"
            :class="{ 'sidebar-link-active': clientTab === 'orders' }"
            @click="clientTab = 'orders'; fetchOrders()"
          >
            <Package class="h-5 w-5" />
            Мои заказы
          </button>
          <div class="my-4 h-px bg-base-200"></div>
          <button class="sidebar-link text-error hover:bg-error/10" @click="logout">
            <LogOut class="h-5 w-5" />
            Выход
          </button>
        </aside>

        <section class="rounded-[2.5rem] border border-white/50 bg-base-100/90 p-8 shadow-[0_28px_80px_-50px_rgba(15,23,42,0.5)] backdrop-blur">
          <div v-if="clientTab === 'profile'">
            <div class="mb-8 flex flex-col gap-6 sm:flex-row sm:items-start sm:justify-between">
              <div class="flex items-center gap-5">
                <div class="profile-avatar">
                  {{ displayInitial }}
                </div>
                <div>
                  <h1 class="text-4xl font-black tracking-[-0.06em] text-secondary">{{ displayName }}</h1>
                  <p class="mt-2 text-lg font-medium text-secondary/45">{{ displayEmail }}</p>
                </div>
              </div>

              <button class="ghost-action" @click="startClientEditing">
                <PencilLine class="h-4 w-4" />
                Изменить
              </button>
            </div>

            <div class="grid gap-4 md:grid-cols-2">
              <article v-for="card in clientInfoCards" :key="card.label" class="info-card">
                <div class="info-card-label">
                  <component :is="card.icon" class="h-4 w-4" />
                  {{ card.label }}
                </div>
                <div v-if="card.label === 'Роль'" class="mt-3 inline-flex rounded-full px-3 py-1 text-xs font-black uppercase tracking-[0.16em]" :class="roleBadgeClass">
                  {{ card.value }}
                </div>
                <p v-else class="mt-3 text-xl font-black tracking-tight text-secondary">{{ card.value }}</p>
              </article>
            </div>

            <form v-if="clientEditing" class="mt-8 grid gap-4 rounded-[2rem] border border-base-200 bg-base-100/85 p-5" @submit.prevent="saveClientProfile">
              <div class="grid gap-4 md:grid-cols-2">
                <label class="edit-field">
                  <span class="edit-field-label">Полное имя</span>
                  <input v-model="clientDraft.fullName" type="text" class="edit-input" placeholder="Иван Иванов" />
                </label>
                <label class="edit-field">
                  <span class="edit-field-label">Телефон</span>
                  <input v-model="clientDraft.phone" type="text" class="edit-input" placeholder="+7 999 000 00 00" />
                </label>
              </div>

              <div class="flex flex-wrap justify-end gap-3">
                <button type="button" class="ghost-action" @click="clientEditing = false">
                  <X class="h-4 w-4" />
                  Отмена
                </button>
                <button type="submit" class="primary-action" :disabled="isSaving">
                  <span v-if="isSaving" class="loading loading-spinner loading-sm"></span>
                  <Save v-else class="h-4 w-4" />
                  Сохранить
                </button>
              </div>
            </form>
          </div>

          <div v-else>
            <div class="mb-6">
              <p class="text-[11px] font-black uppercase tracking-[0.24em] text-secondary/35">История заказов</p>
              <h2 class="mt-2 text-3xl font-black tracking-[-0.06em] text-secondary">Мои заказы</h2>
            </div>

            <div v-if="loadingOrders" class="flex justify-center py-14">
              <span class="loading loading-spinner loading-lg text-primary"></span>
            </div>

            <div v-else-if="orders.length === 0" class="rounded-[2rem] border border-dashed border-base-300 bg-base-100/70 px-6 py-12 text-center">
              <Package class="mx-auto h-14 w-14 text-secondary/20" />
              <h3 class="mt-4 text-2xl font-black tracking-tight text-secondary">Заказов пока нет</h3>
              <p class="mt-2 text-sm font-medium text-secondary/45">Когда появятся заказы, они будут отображаться здесь.</p>
            </div>

            <div v-else class="space-y-4">
              <article v-for="order in orders" :key="order.order_id" class="rounded-[2rem] border border-base-200 bg-base-100/90 p-6 shadow-[0_14px_35px_-28px_rgba(15,23,42,0.45)]">
                <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                  <div>
                    <div class="flex flex-wrap items-center gap-3">
                      <span class="font-mono text-lg font-black tracking-tight text-secondary">#{{ order.order_number }}</span>
                      <span class="rounded-full bg-primary/10 px-3 py-1 text-[10px] font-black uppercase tracking-[0.18em] text-primary">
                        {{ order.status }}
                      </span>
                    </div>
                    <p class="mt-2 text-sm font-semibold text-secondary/45">{{ formatOrderDate(order.created_at) }}</p>
                  </div>

                  <div class="text-left sm:text-right">
                    <p class="text-2xl font-black tracking-tight text-primary">{{ formatPrice(order.final_price) }}</p>
                    <router-link :to="`/order/${order.order_id}`" class="mt-2 inline-flex text-sm font-black text-primary">
                      Подробнее
                    </router-link>
                  </div>
                </div>
              </article>
            </div>
          </div>
        </section>
      </div>
    </div>

    <div v-else-if="currentStaffProfile" class="mx-auto max-w-[1440px] px-4">
      <div class="grid gap-8 xl:grid-cols-[220px_minmax(0,1fr)]">
        <aside class="staff-sidebar">
          <div class="mb-8 px-2">
            <p class="text-4xl font-black tracking-[-0.08em] text-primary">Pizza Good</p>
          </div>

          <div class="space-y-2">
            <button
              v-for="item in staffNavItems"
              :key="item.id"
              class="sidebar-link"
              :class="{ 'sidebar-link-active': item.id === 'info' }"
              @click="handleStaffNav(item)"
            >
              <component :is="item.icon" class="h-5 w-5" />
              {{ item.label }}
            </button>
          </div>

          <div class="my-6 h-px bg-base-200"></div>

          <button class="sidebar-link text-error hover:bg-error/10" @click="logout">
            <LogOut class="h-5 w-5" />
            Выход
          </button>
        </aside>

        <section class="rounded-[2.6rem] border border-white/55 bg-base-100/92 p-7 shadow-[0_32px_90px_-60px_rgba(15,23,42,0.55)] backdrop-blur sm:p-9">
          <div class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
            <div class="flex flex-col gap-5 sm:flex-row sm:items-start">
              <div class="profile-avatar">
                {{ displayInitial }}
              </div>

              <div>
                <p v-if="currentStaffRole === 'manager'" class="text-sm font-semibold text-secondary/55">
                  Вы вошли как:
                  <span class="font-black text-primary">{{ roleLabel }}</span>
                </p>
                <h1 class="mt-1 text-4xl font-black tracking-[-0.07em] text-secondary sm:text-5xl">{{ currentStaffProfile.fullName }}</h1>
                <p class="mt-2 text-lg font-medium text-secondary/45">{{ currentStaffProfile.email }}</p>

                <div class="mt-4 flex flex-wrap gap-3">
                  <span class="inline-flex rounded-full px-3 py-1 text-[11px] font-black uppercase tracking-[0.18em]" :class="roleBadgeClass">
                    {{ currentStaffProfile.roleLabel }}
                  </span>
                  <span
                    v-if="currentStaffProfile.status"
                    class="inline-flex rounded-full px-3 py-1 text-[11px] font-black uppercase tracking-[0.18em]"
                    :class="statusBadgeClass"
                  >
                    {{ currentStaffProfile.status }}
                  </span>
                </div>
              </div>
            </div>

            <div class="flex flex-col gap-4 sm:flex-row sm:items-start">
              <div class="rounded-[1.5rem] border border-base-200 bg-base-100/80 px-5 py-4 text-sm shadow-[0_18px_40px_-34px_rgba(15,23,42,0.35)]">
                <div class="flex items-center gap-2 text-secondary/45">
                  <Shield class="h-4 w-4" />
                  <span>Вы вошли как: <strong class="font-black text-primary">{{ roleLabel }}</strong></span>
                </div>
                <div class="mt-3 flex items-center gap-2 text-secondary/45">
                  <User class="h-4 w-4" />
                  <span>Пользователь: <strong class="font-black text-secondary">{{ displayShortName }}</strong></span>
                </div>
              </div>

              <button v-if="!staffEditing" class="ghost-action self-start" @click="startStaffEditing">
                <PencilLine class="h-4 w-4" />
                Изменить
              </button>
            </div>
          </div>

          <div v-if="!staffEditing" class="mt-8 grid gap-4 md:grid-cols-2">
            <article v-for="card in staffInfoCards" :key="card.key" class="info-card">
              <div class="info-card-label">
                <component :is="card.icon" class="h-4 w-4" />
                {{ card.label }}
              </div>
              <p class="mt-1 text-xs font-semibold text-secondary/28">{{ card.hint }}</p>

              <div v-if="card.tone === 'role'" class="mt-4 inline-flex rounded-full px-3 py-1 text-xs font-black uppercase tracking-[0.16em]" :class="roleBadgeClass">
                {{ card.value }}
              </div>
              <div v-else-if="card.tone === 'status'" class="mt-4 inline-flex rounded-full px-3 py-1 text-xs font-black uppercase tracking-[0.16em]" :class="statusBadgeClass">
                {{ card.value }}
              </div>
              <p v-else class="mt-4 text-2xl font-black tracking-tight text-secondary">{{ card.value }}</p>
            </article>
          </div>

          <form v-else class="mt-8 space-y-5" @submit.prevent="saveStaffProfile">
            <div class="grid gap-4 md:grid-cols-2">
              <label class="edit-field">
                <span class="edit-field-label">ID пользователя</span>
                <span class="edit-field-hint">{{ currentStaffProfile.userIdHint }}</span>
                <input :value="currentStaffProfile.userId" type="text" class="edit-input read-only" readonly />
              </label>

              <label v-if="currentStaffProfile.staffId" class="edit-field">
                <span class="edit-field-label">{{ currentStaffProfile.staffIdLabel }}</span>
                <span class="edit-field-hint">{{ currentStaffProfile.staffIdHint }}</span>
                <input :value="currentStaffProfile.staffId" type="text" class="edit-input read-only" readonly />
              </label>

              <label class="edit-field">
                <span class="edit-field-label">ФИО</span>
                <span class="edit-field-hint">users.name</span>
                <input v-model="staffDraft.fullName" type="text" class="edit-input" />
              </label>

              <label class="edit-field">
                <span class="edit-field-label">Email</span>
                <span class="edit-field-hint">users.email</span>
                <input v-model="staffDraft.email" type="email" class="edit-input" />
              </label>

              <label v-if="currentStaffProfile.phone !== undefined" class="edit-field">
                <span class="edit-field-label">Телефон</span>
                <span class="edit-field-hint">users.phone</span>
                <input v-model="staffDraft.phone" type="text" class="edit-input" />
              </label>

              <label class="edit-field">
                <span class="edit-field-label">Роль</span>
                <span class="edit-field-hint">users.role</span>
                <input :value="currentStaffProfile.roleLabel" type="text" class="edit-input read-only" readonly />
              </label>

              <label v-if="currentStaffProfile.status" class="edit-field">
                <span class="edit-field-label">{{ currentStaffProfile.statusLabel }}</span>
                <span class="edit-field-hint">users.status</span>
                <select v-model="staffDraft.status" class="edit-input">
                  <option value="Активен">Активен</option>
                  <option value="На смене">На смене</option>
                  <option value="Неактивен">Неактивен</option>
                </select>
              </label>

              <label class="edit-field">
                <span class="edit-field-label">Дата регистрации</span>
                <span class="edit-field-hint">users.created_at</span>
                <input :value="currentStaffProfile.registeredAt" type="text" class="edit-input read-only" readonly />
              </label>
            </div>

            <div class="flex flex-wrap gap-3">
              <button type="submit" class="primary-action" :disabled="isSaving">
                <span v-if="isSaving" class="loading loading-spinner loading-sm"></span>
                <Save v-else class="h-4 w-4" />
                Сохранить изменения
              </button>

              <button type="button" class="ghost-action" @click="stopEditing">
                <X class="h-4 w-4" />
                Отмена
              </button>
            </div>
          </form>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.staff-sidebar {
  border-radius: 2rem;
  border: 1px solid rgba(255, 255, 255, 0.5);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.94), rgba(248, 250, 252, 0.88));
  padding: 1.4rem;
  box-shadow: 0 26px 70px -58px rgba(15, 23, 42, 0.55);
}

.sidebar-link {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.85rem;
  border-radius: 1.2rem;
  padding: 0.95rem 1rem;
  color: rgba(15, 23, 42, 0.72);
  font-size: 0.97rem;
  font-weight: 700;
  transition: all 0.2s ease;
}

.sidebar-link:hover {
  transform: translateX(2px);
  background: rgba(15, 23, 42, 0.04);
}

.sidebar-link-active {
  background: linear-gradient(135deg, rgba(255, 91, 102, 0.14), rgba(255, 91, 102, 0.06));
  color: rgb(255, 91, 102);
  box-shadow: inset 0 0 0 1px rgba(255, 91, 102, 0.08);
}

.profile-avatar {
  display: flex;
  height: 5.75rem;
  width: 5.75rem;
  align-items: center;
  justify-content: center;
  border-radius: 1.65rem;
  background: linear-gradient(160deg, #ff6d7a 0%, #ff3e54 100%);
  color: white;
  font-size: 2.4rem;
  font-weight: 900;
  box-shadow: 0 26px 55px -30px rgba(255, 62, 84, 0.6);
}

.info-card {
  border-radius: 1.45rem;
  border: 1px solid rgba(226, 232, 240, 0.95);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(248, 250, 252, 0.86));
  padding: 1.15rem 1.2rem 1.3rem;
  box-shadow: 0 18px 36px -32px rgba(15, 23, 42, 0.45);
}

.info-card-label {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  color: rgba(15, 23, 42, 0.42);
  font-size: 0.72rem;
  font-weight: 900;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.ghost-action,
.primary-action {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.55rem;
  border-radius: 1rem;
  padding: 0.9rem 1.35rem;
  font-size: 0.9rem;
  font-weight: 900;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  transition: all 0.2s ease;
}

.ghost-action {
  border: 1px solid rgba(148, 163, 184, 0.45);
  color: rgba(15, 23, 42, 0.8);
  background: white;
}

.ghost-action:hover {
  transform: translateY(-1px);
  border-color: rgba(255, 91, 102, 0.35);
  color: rgb(255, 91, 102);
}

.primary-action {
  border: 0;
  color: white;
  background: linear-gradient(135deg, #ff5b66 0%, #ff3f55 100%);
  box-shadow: 0 24px 50px -28px rgba(255, 63, 85, 0.72);
}

.primary-action:hover {
  transform: translateY(-1px);
  box-shadow: 0 26px 55px -26px rgba(255, 63, 85, 0.78);
}

.primary-action:disabled {
  transform: none;
  opacity: 0.75;
  box-shadow: none;
}

.edit-field {
  display: flex;
  flex-direction: column;
  border-radius: 1.45rem;
  border: 1px solid rgba(226, 232, 240, 0.95);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(248, 250, 252, 0.86));
  padding: 1rem 1.05rem 1.1rem;
  gap: 0.2rem;
}

.edit-field-label {
  color: rgba(15, 23, 42, 0.56);
  font-size: 0.86rem;
  font-weight: 800;
}

.edit-field-hint {
  color: rgba(15, 23, 42, 0.28);
  font-size: 0.74rem;
  font-weight: 700;
}

.edit-input {
  margin-top: 0.6rem;
  min-height: 3.35rem;
  border-radius: 1rem;
  border: 1px solid rgba(203, 213, 225, 0.95);
  background: white;
  padding: 0.9rem 1rem;
  color: rgb(30, 41, 59);
  font-size: 1rem;
  font-weight: 700;
  outline: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.edit-input:focus {
  border-color: rgba(255, 91, 102, 0.55);
  box-shadow: 0 0 0 4px rgba(255, 91, 102, 0.08);
}

.read-only {
  background: rgba(248, 250, 252, 0.9);
  color: rgba(15, 23, 42, 0.62);
}
</style>

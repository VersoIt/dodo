<script setup lang="ts">
import { computed, inject, ref } from 'vue'
import AppModal from '../components/shared/AppModal.vue'
import { useAuthStore } from '../store/auth'
import {
  AlertCircle,
  CalendarDays,
  CheckCircle2,
  ChefHat,
  Clock3,
  Flag,
  Play,
  ReceiptText,
  User,
  Wallet
} from 'lucide-vue-next'

type ToastType = 'success' | 'error' | 'info'
type KitchenTicketStatus = 'new' | 'cooking' | 'ready'
type ValueTone = 'success' | 'info' | 'warning'
type KitchenIcon = typeof ChefHat

interface KitchenItem {
  name: string
  quantity: number
  price: number
}

interface KitchenTicket {
  id: string
  publicOrderId: string
  kitchenTaskId: string
  chefName: string
  chefId: string
  createdAt: string
  acceptedAt: string | null
  completedAt: string | null
  status: KitchenTicketStatus
  paymentStatus: string
  items: KitchenItem[]
}

interface TicketRow {
  label: string
  value: string
  icon: KitchenIcon
  tone?: ValueTone
  dimmed?: boolean
  multiline?: boolean
}

interface PendingAction {
  ticketId: string
  nextStatus: Extract<KitchenTicketStatus, 'cooking' | 'ready'>
}

const authStore = useAuthStore()
const addToast = inject<(message: string, type?: ToastType) => void>('addToast', () => undefined)

const fallbackChefName = 'Алексей Смирнов'

const currentChefName = computed(() => authStore.user?.name || fallbackChefName)
const currentChefId = computed(() => String(authStore.user?.id || authStore.user?.user_id || '1042'))
const currentChefInitial = computed(() => currentChefName.value.charAt(0).toUpperCase())

const formatCurrency = (value: number) => `${new Intl.NumberFormat('ru-RU').format(value)} ₽`
const formatDateTime = (date: Date) =>
  new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
    .format(date)
    .replace(',', '')

const statusBadgeLabel = (status: KitchenTicketStatus) => {
  if (status === 'new') return 'Новое'
  if (status === 'cooking') return 'В работе'
  return 'Готово'
}

const statusFieldLabel = (status: KitchenTicketStatus) => {
  if (status === 'new') return 'Ожидает приготовления'
  if (status === 'cooking') return 'В работе'
  return 'Готов'
}

const statusTone = (status: KitchenTicketStatus): ValueTone => {
  if (status === 'new') return 'warning'
  if (status === 'cooking') return 'info'
  return 'success'
}

const statusBadgeClass = (status: KitchenTicketStatus) => {
  if (status === 'new') return 'bg-blue-500/10 text-blue-600'
  if (status === 'cooking') return 'bg-amber-500/12 text-amber-600'
  return 'bg-emerald-500/12 text-emerald-600'
}

const toneBadgeClass = (tone: ValueTone) => {
  if (tone === 'success') return 'bg-emerald-500/12 text-emerald-700'
  if (tone === 'info') return 'bg-sky-500/12 text-sky-700'
  return 'bg-amber-500/12 text-amber-700'
}

const ticketTotal = (ticket: KitchenTicket) =>
  ticket.items.reduce((sum, item) => sum + item.price * item.quantity, 0)

const ticketSummary = (ticket: KitchenTicket) =>
  ticket.items.map((item) => `${item.name} ×${item.quantity}`).join(', ')

const tickets = ref<KitchenTicket[]>([
  {
    id: 'KO-10458',
    publicOrderId: 'PG-2026.04.05-425e8a8b4d8d',
    kitchenTaskId: 'KCH-1048-01',
    chefName: fallbackChefName,
    chefId: 'CHF-023',
    createdAt: '05.04.2026 19:05',
    acceptedAt: '05.04.2026 19:08',
    completedAt: null,
    status: 'cooking',
    paymentStatus: 'Оплачен',
    items: [
      { name: 'Четыре сыра', quantity: 3, price: 690 },
      { name: 'Соус чесночный', quantity: 1, price: 90 },
      { name: 'Морс клюквенный', quantity: 2, price: 140 }
    ]
  },
  {
    id: 'KO-10459',
    publicOrderId: 'PG-2026.04.05-425e8a8b4d8e',
    kitchenTaskId: 'KCH-1048-02',
    chefName: fallbackChefName,
    chefId: 'CHF-023',
    createdAt: '05.04.2026 19:07',
    acceptedAt: null,
    completedAt: null,
    status: 'new',
    paymentStatus: 'Оплачен',
    items: [
      { name: 'Трюфельная с грибами', quantity: 1, price: 820 },
      { name: 'Четыре сыра', quantity: 1, price: 690 },
      { name: 'Пепперони Фреш', quantity: 1, price: 710 },
      { name: 'Лимонад', quantity: 2, price: 160 }
    ]
  },
  {
    id: 'KO-10460',
    publicOrderId: 'PG-2026.04.05-425e8a8b4d8f',
    kitchenTaskId: 'KCH-1048-03',
    chefName: fallbackChefName,
    chefId: 'CHF-023',
    createdAt: '05.04.2026 18:42',
    acceptedAt: '05.04.2026 18:42',
    completedAt: '05.04.2026 19:01',
    status: 'ready',
    paymentStatus: 'Оплачен',
    items: [
      { name: 'Маргарита Буррата', quantity: 2, price: 750 },
      { name: 'Соус чесночный', quantity: 2, price: 90 }
    ]
  }
])

const sortedTickets = computed(() => {
  const order: Record<KitchenTicketStatus, number> = { new: 0, cooking: 1, ready: 2 }

  return [...tickets.value].sort((left, right) => order[left.status] - order[right.status])
})

const buildTicketRows = (ticket: KitchenTicket): TicketRow[] => [
  {
    label: 'ID заказа',
    value: ticket.publicOrderId,
    icon: ReceiptText
  },
  {
    label: 'Дата создания задания',
    value: ticket.createdAt,
    icon: CalendarDays
  },
  {
    label: 'Статус заказа',
    value: statusFieldLabel(ticket.status),
    icon: Flag,
    tone: statusTone(ticket.status)
  },
  {
    label: 'Статус оплаты',
    value: ticket.paymentStatus,
    icon: Wallet,
    tone: 'success'
  },
  {
    label: 'Назначенный повар',
    value: `${ticket.chefId} / ${ticket.chefName}`,
    icon: User
  },
  {
    label: 'Время начала приготовления',
    value: ticket.acceptedAt || 'не начато',
    icon: Clock3,
    dimmed: !ticket.acceptedAt
  },
  {
    label: 'Время завершения приготовления',
    value: ticket.completedAt || 'не завершено',
    icon: CheckCircle2,
    dimmed: !ticket.completedAt
  }
]

const showConfirmModal = ref(false)
const pendingAction = ref<PendingAction | null>(null)

const selectedConfirmTicket = computed(() => {
  if (!pendingAction.value) return null
  return tickets.value.find((ticket) => ticket.id === pendingAction.value?.ticketId) || null
})

const isStartConfirm = computed(() => pendingAction.value?.nextStatus === 'cooking')

const confirmTitle = computed(() =>
  isStartConfirm.value ? 'Подтверждение принятия заказа' : 'Завершить приготовление?'
)

const confirmDescription = computed(() =>
  isStartConfirm.value
    ? 'После подтверждения система зафиксирует время начала приготовления и назначит заказ текущему повару.'
    : 'Подтвердите завершение производственного этапа по заказу на кухне Pizza Good.'
)

const confirmPrimaryLabel = computed(() =>
  isStartConfirm.value ? 'Подтвердить' : 'Подтвердить готовность'
)

const confirmRows = computed<TicketRow[]>(() => {
  const ticket = selectedConfirmTicket.value

  if (!ticket) return []

  if (isStartConfirm.value) {
    return [
      { label: 'ID заказа', value: ticket.publicOrderId, icon: ReceiptText },
      { label: 'ID кухонного задания', value: ticket.kitchenTaskId, icon: ReceiptText },
      { label: 'Текущий повар', value: currentChefName.value, icon: User },
      { label: 'ID повара', value: currentChefId.value, icon: User },
      { label: 'Текущий статус заказа', value: ticket.paymentStatus, icon: Wallet, tone: 'success' },
      { label: 'Новый статус', value: 'Готовится', icon: Flag, tone: 'info' },
      {
        label: 'Время начала приготовления',
        value: 'Будет установлено автоматически',
        icon: Clock3,
        dimmed: true
      },
      { label: 'Краткий состав заказа', value: ticketSummary(ticket), icon: ChefHat }
    ]
  }

  return [
    { label: 'ID заказа', value: ticket.publicOrderId, icon: ReceiptText },
    { label: 'ID кухонного задания', value: ticket.kitchenTaskId, icon: ReceiptText },
    { label: 'Повар', value: ticket.chefName, icon: User },
    { label: 'ID повара', value: ticket.chefId, icon: User },
    { label: 'Время начала приготовления', value: ticket.acceptedAt || 'не начато', icon: Clock3 },
    {
      label: 'Время завершения приготовления',
      value: 'Будет установлено автоматически после подтверждения',
      icon: CheckCircle2,
      dimmed: true,
      multiline: true
    },
    { label: 'Текущий статус', value: 'Готовится', icon: Flag, tone: 'info' },
    { label: 'Новый статус после подтверждения', value: 'Готов', icon: CheckCircle2, tone: 'success' },
    { label: 'Краткий состав заказа', value: ticketSummary(ticket), icon: ChefHat }
  ]
})

const openActionModal = (ticket: KitchenTicket) => {
  if (ticket.status === 'ready') return

  pendingAction.value = {
    ticketId: ticket.id,
    nextStatus: ticket.status === 'new' ? 'cooking' : 'ready'
  }
  showConfirmModal.value = true
}

const closeConfirmModal = () => {
  showConfirmModal.value = false
  pendingAction.value = null
}

const handleConfirm = () => {
  const action = pendingAction.value

  if (!action) return

  const now = formatDateTime(new Date())

  tickets.value = tickets.value.map((ticket) => {
    if (ticket.id !== action.ticketId) return ticket

    if (action.nextStatus === 'cooking') {
      return {
        ...ticket,
        status: 'cooking',
        chefName: currentChefName.value,
        chefId: currentChefId.value,
        acceptedAt: now
      }
    }

    return {
      ...ticket,
      status: 'ready',
      chefName: currentChefName.value,
      chefId: currentChefId.value,
      completedAt: now
    }
  })

  addToast(
    action.nextStatus === 'cooking'
      ? 'Заказ принят в работу'
      : 'Приготовление подтверждено, заказ готов к передаче',
    'success'
  )

  closeConfirmModal()
}

const actionButtonClass = (status: KitchenTicketStatus) => {
  if (status === 'new') {
    return 'bg-gradient-to-r from-primary to-red-500 text-white shadow-[0_18px_40px_-18px_rgba(255,82,94,0.85)] hover:-translate-y-0.5'
  }

  if (status === 'cooking') {
    return 'bg-gradient-to-r from-emerald-500 to-green-500 text-white shadow-[0_18px_40px_-18px_rgba(16,185,129,0.75)] hover:-translate-y-0.5'
  }

  return 'cursor-default bg-emerald-100 text-emerald-700 shadow-none'
}

const actionButtonLabel = (status: KitchenTicketStatus) => {
  if (status === 'new') return 'Принять заказ'
  if (status === 'cooking') return 'Завершить приготовление'
  return 'Передано в доставку'
}
</script>

<template>
  <div class="min-h-screen bg-[radial-gradient(circle_at_top,_rgba(255,255,255,0.98),_rgba(244,246,252,0.95)_46%,_rgba(236,240,247,0.92)_100%)]">
    <div class="mx-auto max-w-7xl px-4 py-8 md:px-6">
      <div class="mb-10 flex flex-col gap-6 xl:flex-row xl:items-center xl:justify-between">
        <div class="flex items-center gap-5">
          <div class="flex h-24 w-24 items-center justify-center rounded-[2rem] bg-gradient-to-br from-primary via-red-500 to-red-400 text-white shadow-[0_24px_55px_-24px_rgba(255,82,94,0.9)]">
            <ChefHat class="h-12 w-12" />
          </div>

          <div>
            <p class="text-xs font-black uppercase tracking-[0.42em] text-base-content/30">Live order queue</p>
            <h1 class="text-5xl font-black uppercase tracking-[-0.08em] text-base-content">Кухня</h1>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-4 rounded-[2rem] border border-base-200/80 bg-base-100/90 px-4 py-3 shadow-[0_18px_45px_-30px_rgba(15,23,42,0.35)] backdrop-blur">
          <div class="rounded-full bg-amber-500/12 px-4 py-2 text-xs font-black uppercase tracking-[0.28em] text-amber-700">
            Кухня
          </div>

          <div class="h-10 w-px bg-base-200"></div>

          <div class="flex items-center gap-3">
            <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-primary/8 text-lg font-black text-primary shadow-inner">
              {{ currentChefInitial }}
            </div>

            <div>
              <p class="text-[11px] font-black uppercase tracking-[0.24em] text-base-content/30">Вы вошли как</p>
              <p class="text-base font-black text-base-content">Повар {{ currentChefName }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="grid gap-8 xl:grid-cols-2">
        <article
          v-for="ticket in sortedTickets"
          :key="ticket.id"
          class="rounded-[2.5rem] border border-base-200/80 bg-base-100/95 p-7 shadow-[0_28px_70px_-38px_rgba(15,23,42,0.4)] transition duration-300 hover:-translate-y-1 hover:shadow-[0_34px_80px_-38px_rgba(15,23,42,0.45)]"
        >
          <div class="mb-7 flex items-start justify-between gap-4">
            <div class="flex items-start gap-4">
              <div class="flex h-14 w-14 items-center justify-center rounded-[1.25rem] border border-base-200 bg-base-100 text-base-content/55 shadow-inner">
                <ReceiptText class="h-7 w-7" />
              </div>

              <div>
                <h2 class="text-[2.15rem] font-black tracking-[-0.08em] text-base-content">{{ ticket.id }}</h2>
                <p class="text-[11px] font-bold uppercase tracking-[0.22em] text-base-content/32">ID кухонного задания</p>
              </div>
            </div>

            <span
              class="rounded-full px-4 py-2 text-[11px] font-black uppercase tracking-[0.28em]"
              :class="statusBadgeClass(ticket.status)"
            >
              {{ statusBadgeLabel(ticket.status) }}
            </span>
          </div>

          <div class="overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white/90">
            <div
              v-for="row in buildTicketRows(ticket)"
              :key="`${ticket.id}-${row.label}`"
              class="grid grid-cols-[minmax(0,1.2fr)_minmax(0,0.8fr)] items-center gap-4 border-b border-base-200/80 px-5 py-4 last:border-b-0"
            >
              <div class="flex items-center gap-3 text-base font-medium text-base-content/60">
                <component :is="row.icon" class="h-5 w-5 shrink-0 text-base-content/45" />
                <span>{{ row.label }}</span>
              </div>

              <div class="text-right text-base font-semibold text-base-content">
                <span
                  v-if="row.tone"
                  class="inline-flex rounded-full px-3 py-1 text-sm font-bold"
                  :class="toneBadgeClass(row.tone)"
                >
                  {{ row.value }}
                </span>
                <span
                  v-else
                  :class="[
                    row.dimmed ? 'text-base-content/38' : 'text-base-content',
                    row.multiline ? 'whitespace-pre-line' : ''
                  ]"
                >
                  {{ row.value }}
                </span>
              </div>
            </div>
          </div>

          <div class="mt-7">
            <div class="mb-4 flex items-center justify-between gap-4">
              <h3 class="text-2xl font-black tracking-[-0.05em] text-base-content">Состав заказа</h3>
              <p class="text-sm font-semibold text-base-content/45">{{ ticket.items.length }} позиции</p>
            </div>

            <div class="overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white">
              <div class="grid grid-cols-[1.35fr_0.7fr_0.8fr] bg-base-100/70 px-5 py-3 text-[11px] font-black uppercase tracking-[0.18em] text-base-content/38">
                <span>Товар</span>
                <span class="text-center">Количество</span>
                <span class="text-right">Цена на момент заказа</span>
              </div>

              <div
                v-for="item in ticket.items"
                :key="`${ticket.id}-${item.name}`"
                class="grid grid-cols-[1.35fr_0.7fr_0.8fr] items-center border-t border-base-200/80 px-5 py-4 text-base"
              >
                <span class="font-medium text-base-content/82">{{ item.name }}</span>
                <span class="text-center font-semibold text-base-content">{{ item.quantity }}</span>
                <span class="text-right font-semibold text-base-content">{{ formatCurrency(item.price) }}</span>
              </div>

              <div class="flex items-center justify-between border-t border-base-200/80 px-5 py-4">
                <span class="text-sm font-bold uppercase tracking-[0.18em] text-base-content/38">Сумма тикета</span>
                <span class="text-xl font-black tracking-[-0.04em] text-primary">{{ formatCurrency(ticketTotal(ticket)) }}</span>
              </div>
            </div>
          </div>

          <button
            class="mt-8 flex h-16 w-full items-center justify-center gap-3 rounded-[1.5rem] text-lg font-black uppercase tracking-tight transition"
            :class="actionButtonClass(ticket.status)"
            :disabled="ticket.status === 'ready'"
            @click="openActionModal(ticket)"
          >
            <Play v-if="ticket.status === 'new'" class="h-6 w-6" />
            <CheckCircle2 v-else class="h-6 w-6" />
            {{ actionButtonLabel(ticket.status) }}
          </button>
        </article>
      </div>
    </div>

    <AppModal :show="showConfirmModal" maxWidth="2xl" @close="closeConfirmModal">
      <div v-if="selectedConfirmTicket" class="px-7 pb-8 pt-7 sm:px-10 sm:pb-10">
        <div
          v-if="isStartConfirm"
          class="mx-auto mb-6 flex w-fit items-center rounded-full bg-base-200/70 px-5 py-3 text-sm font-semibold text-base-content/65 shadow-inner"
        >
          Вы вошли как:
          <span class="ml-1.5 font-black text-base-content">Повар</span>
        </div>

        <div class="mb-6 flex justify-center">
          <div class="flex h-24 w-24 items-center justify-center rounded-[2rem] bg-primary/10 text-primary shadow-inner">
            <AlertCircle class="h-12 w-12" />
          </div>
        </div>

        <h3 class="mx-auto max-w-xl text-center text-4xl font-black uppercase leading-none tracking-[-0.08em] text-base-content sm:text-5xl">
          {{ confirmTitle }}
        </h3>

        <p class="mx-auto mt-4 max-w-2xl text-center text-base font-medium leading-relaxed text-base-content/55">
          {{ confirmDescription }}
        </p>

        <div
          v-if="!isStartConfirm"
          class="mt-7 flex items-center gap-3 rounded-[1.35rem] bg-base-100 px-4 py-4 text-base font-semibold text-base-content/70 shadow-inner"
        >
          <User class="h-5 w-5 text-base-content/55" />
          <span>
            Вы вошли как:
            <span class="font-black text-base-content">Повар</span>
          </span>
        </div>

        <div class="mt-7 overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white shadow-[0_20px_48px_-34px_rgba(15,23,42,0.35)]">
          <div
            v-for="row in confirmRows"
            :key="`${selectedConfirmTicket.id}-${row.label}`"
            class="grid grid-cols-[minmax(0,1.12fr)_minmax(0,0.88fr)] items-center gap-4 border-b border-base-200/80 px-5 py-4 last:border-b-0"
          >
            <div class="flex items-center gap-3 text-base font-medium text-base-content/60">
              <component :is="row.icon" class="h-5 w-5 shrink-0 text-base-content/45" />
              <span>{{ row.label }}</span>
            </div>

            <div class="text-right text-base font-semibold text-base-content">
              <span
                v-if="row.tone"
                class="inline-flex rounded-full px-3 py-1 text-sm font-bold"
                :class="toneBadgeClass(row.tone)"
              >
                {{ row.value }}
              </span>
              <span
                v-else
                :class="[
                  row.dimmed ? 'text-base-content/45' : 'text-base-content',
                  row.multiline ? 'whitespace-pre-line' : ''
                ]"
              >
                {{ row.value }}
              </span>
            </div>
          </div>
        </div>

        <div
          v-if="!isStartConfirm"
          class="mt-6 flex items-start gap-4 rounded-[1.35rem] border border-base-200/90 bg-base-100/85 px-5 py-4"
        >
          <div class="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-base-200 text-base-content/65">
            <AlertCircle class="h-5 w-5" />
          </div>

          <p class="text-sm font-medium leading-relaxed text-base-content/58">
            После подтверждения система зафиксирует время завершения приготовления и передаст заказ в модуль доставки.
          </p>
        </div>

        <div class="mt-8 flex flex-col gap-3">
          <button
            class="h-16 rounded-[1.45rem] bg-gradient-to-r from-primary via-red-500 to-red-600 text-lg font-black uppercase tracking-tight text-white shadow-[0_24px_48px_-22px_rgba(255,82,94,0.9)] transition hover:-translate-y-0.5"
            @click="handleConfirm"
          >
            {{ confirmPrimaryLabel }}
          </button>

          <button
            class="h-12 text-sm font-black uppercase tracking-[0.24em] text-base-content/35 transition hover:text-base-content/55"
            @click="closeConfirmModal"
          >
            Отмена
          </button>
        </div>
      </div>
    </AppModal>
  </div>
</template>

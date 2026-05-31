<script setup lang="ts">
import { computed, inject, ref } from 'vue'
import AppModal from '../components/shared/AppModal.vue'
import { useAuthStore } from '../store/auth'
import {
  AlertCircle,
  CalendarDays,
  CheckCircle2,
  Clock3,
  CreditCard,
  MapPin,
  MessageSquare,
  Package,
  ReceiptText,
  Truck,
  User,
  Wallet
} from 'lucide-vue-next'

type ToastType = 'success' | 'error' | 'info'
type DeliveryTicketStatus = 'ready' | 'delivering' | 'delivered' | 'failed'
type PaymentState = 'paid' | 'awaiting'
type DeliveryTone = 'success' | 'info' | 'warning' | 'danger'
type DeliveryIcon = typeof Truck

interface DeliveryItem {
  name: string
  quantity: number
  price: number
}

interface DeliveryAddress {
  city: string
  street: string
  house: string
  apartment: string
  floor: string
  comment: string
}

interface DeliveryChatMessage {
  id: number
  author: 'Клиент' | 'Курьер' | 'Менеджер'
  text: string
  time: string
}

interface DeliveryTicket {
  id: string
  deliveryCode: string
  publicOrderId: string
  courierId: string
  courierName: string
  createdAt: string
  acceptedAt: string | null
  completedAt: string | null
  status: DeliveryTicketStatus
  paymentMethod: string
  paymentStatus: PaymentState
  total: number
  address: DeliveryAddress
  items: DeliveryItem[]
  courierComment: string
  chat: DeliveryChatMessage[]
}

interface DeliveryRow {
  label: string
  value: string
  tone?: DeliveryTone
  dimmed?: boolean
  multiline?: boolean
}

interface PendingAction {
  ticketId: string
  type: 'start' | 'finish'
}

const authStore = useAuthStore()
const addToast = inject<(message: string, type?: ToastType) => void>('addToast', () => undefined)

const fallbackCourierName = 'Иван Петров'
const rawCourierId = computed(() => authStore.user?.user_id ?? authStore.user?.id)

const currentCourierName = computed(() => authStore.user?.name || fallbackCourierName)
const currentCourierId = computed(() => {
  const value = rawCourierId.value

  if (typeof value === 'number') return String(value)

  if (typeof value === 'string' && /^\d+$/.test(value)) return value

  return '208'
})

const currentCourierInitial = computed(() => currentCourierName.value.charAt(0).toUpperCase())

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

const formatTime = (date: Date) =>
  new Intl.DateTimeFormat('ru-RU', {
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)

const paymentStatusLabel = (status: PaymentState) => (status === 'paid' ? 'Оплачено' : 'Ожидает оплаты')

const statusBadgeLabel = (status: DeliveryTicketStatus) => {
  if (status === 'ready') return 'Готов'
  if (status === 'delivering') return 'В пути'
  if (status === 'delivered') return 'Доставлен'
  return 'Не доставлен'
}

const statusFieldLabel = (status: DeliveryTicketStatus) => {
  if (status === 'ready') return 'Готов к доставке'
  if (status === 'delivering') return 'В пути'
  if (status === 'delivered') return 'Доставлен'
  return 'Не доставлен'
}

const statusTone = (status: DeliveryTicketStatus): DeliveryTone => {
  if (status === 'ready') return 'success'
  if (status === 'delivering') return 'info'
  if (status === 'delivered') return 'success'
  return 'danger'
}

const rowToneClass = (tone: DeliveryTone) => {
  if (tone === 'success') return 'bg-emerald-500/12 text-emerald-700'
  if (tone === 'info') return 'bg-sky-500/12 text-sky-700'
  if (tone === 'warning') return 'bg-amber-500/12 text-amber-700'
  return 'bg-rose-500/12 text-rose-700'
}

const statusBadgeClass = (status: DeliveryTicketStatus) => {
  if (status === 'ready') return 'bg-emerald-500/12 text-emerald-700'
  if (status === 'delivering') return 'bg-sky-500/12 text-sky-700'
  if (status === 'delivered') return 'bg-slate-800 text-white'
  return 'bg-rose-500/12 text-rose-700'
}

const paymentTone = (status: PaymentState): DeliveryTone => (status === 'paid' ? 'success' : 'warning')

const tickets = ref<DeliveryTicket[]>([
  {
    id: 'DLV-10458',
    deliveryCode: 'DL-2026.04.05-001284',
    publicOrderId: 'PG-2026.04.05-425e8a8b4d8d',
    courierId: '208',
    courierName: 'Иван Петров',
    createdAt: '05.04.2026 12:15',
    acceptedAt: null,
    completedAt: null,
    status: 'ready',
    paymentMethod: 'Банковская карта',
    paymentStatus: 'paid',
    total: 1290,
    address: {
      city: 'Москва',
      street: 'Проспект Ленина',
      house: '2',
      apartment: '15',
      floor: '4',
      comment: 'Домофон 15, позвонить за 5 минут'
    },
    items: [{ name: 'Четыре сыра', quantity: 3, price: 430 }],
    courierComment: '',
    chat: [
      { id: 1, author: 'Менеджер', text: 'Заказ готов к передаче, можно забирать.', time: '12:18' },
      { id: 2, author: 'Клиент', text: 'Буду на связи, домофон работает.', time: '12:20' }
    ]
  },
  {
    id: 'DLV-10459',
    deliveryCode: 'DL-2026.04.05-001285',
    publicOrderId: 'PG-2026.04.05-425e8a8b4d8e',
    courierId: '208',
    courierName: 'Иван Петров',
    createdAt: '05.04.2026 12:22',
    acceptedAt: '05.04.2026 12:30',
    completedAt: null,
    status: 'delivering',
    paymentMethod: 'Наличными',
    paymentStatus: 'awaiting',
    total: 1780,
    address: {
      city: 'Москва',
      street: 'Проспект Вернадского',
      house: '92',
      apartment: '48',
      floor: '11',
      comment: 'Оставить у двери'
    },
    items: [
      { name: 'Трюфельная с грибами', quantity: 1, price: 690 },
      { name: 'Четыре сыра', quantity: 1, price: 430 },
      { name: 'Пепперони Фреш', quantity: 1, price: 660 }
    ],
    courierComment: 'Передать лично в руки',
    chat: [
      { id: 1, author: 'Клиент', text: 'Если буду без связи, позвоните за 5 минут.', time: '12:24' },
      { id: 2, author: 'Курьер', text: 'Принял, подъеду примерно через 20 минут.', time: '12:33' }
    ]
  }
])

const sortedTickets = computed(() => {
  const order: Record<DeliveryTicketStatus, number> = {
    ready: 0,
    delivering: 1,
    delivered: 2,
    failed: 3
  }

  return [...tickets.value].sort((left, right) => order[left.status] - order[right.status])
})

const buildDeliveryRows = (ticket: DeliveryTicket): DeliveryRow[] => [
  { label: 'ID заказа', value: ticket.publicOrderId },
  { label: 'ID курьера', value: ticket.courierId },
  { label: 'Курьер', value: ticket.courierName },
  { label: 'Дата создания доставки', value: ticket.createdAt },
  {
    label: 'Дата принятия в доставку',
    value: ticket.acceptedAt || 'не принято',
    tone: ticket.acceptedAt ? undefined : 'danger',
    dimmed: !ticket.acceptedAt
  },
  {
    label: 'Дата завершения доставки',
    value: ticket.completedAt || 'не доставлено',
    tone: ticket.completedAt ? 'success' : ticket.status === 'failed' ? 'danger' : 'warning',
    dimmed: !ticket.completedAt
  },
  { label: 'Статус заказа', value: statusFieldLabel(ticket.status), tone: statusTone(ticket.status) }
]

const buildAddressRows = (ticket: DeliveryTicket): DeliveryRow[] => [
  { label: 'Город', value: ticket.address.city },
  { label: 'Улица', value: ticket.address.street },
  { label: 'Дом', value: ticket.address.house },
  { label: 'Квартира', value: ticket.address.apartment },
  { label: 'Этаж', value: ticket.address.floor },
  { label: 'Комментарий к доставке', value: ticket.address.comment, multiline: true }
]

const buildPaymentRows = (ticket: DeliveryTicket): DeliveryRow[] => [
  { label: 'Итоговая сумма', value: formatCurrency(ticket.total) },
  { label: 'Способ оплаты', value: ticket.paymentMethod },
  { label: 'Статус оплаты', value: paymentStatusLabel(ticket.paymentStatus), tone: paymentTone(ticket.paymentStatus) }
]

const showChatModal = ref(false)
const activeChatTicketId = ref<string | null>(null)
const chatDraft = ref('')

const selectedChatTicket = computed(() => {
  if (!activeChatTicketId.value) return null
  return tickets.value.find((ticket) => ticket.id === activeChatTicketId.value) || null
})

const selectedChatMessages = computed(() =>
  (selectedChatTicket.value?.chat ?? []).filter((message) => message.author !== '\u041a\u043b\u0438\u0435\u043d\u0442')
)

const openChatModal = (ticketId: string) => {
  activeChatTicketId.value = ticketId
  chatDraft.value = ''
  showChatModal.value = true
}

const closeChatModal = () => {
  showChatModal.value = false
  activeChatTicketId.value = null
  chatDraft.value = ''
}

const sendChatMessage = () => {
  const ticket = selectedChatTicket.value
  const text = chatDraft.value.trim()

  if (!ticket || !text) return

  tickets.value = tickets.value.map((item) =>
    item.id === ticket.id
      ? {
          ...item,
          chat: [
            ...item.chat,
            {
              id: item.chat.length + 1,
              author: 'Курьер',
              text,
              time: formatTime(new Date())
            }
          ]
        }
      : item
  )

  chatDraft.value = ''
}

const showConfirmModal = ref(false)
const pendingAction = ref<PendingAction | null>(null)
const finishPaymentReceived = ref(false)
const finishCourierComment = ref('')

const selectedConfirmTicket = computed(() => {
  if (!pendingAction.value) return null
  return tickets.value.find((ticket) => ticket.id === pendingAction.value?.ticketId) || null
})

const isStartConfirm = computed(() => pendingAction.value?.type === 'start')

const startConfirmInfoRows = computed<DeliveryRow[]>(() => {
  const ticket = selectedConfirmTicket.value

  if (!ticket) return []

  return [
    { label: 'ID доставки', value: ticket.deliveryCode },
    { label: 'ID заказа', value: ticket.publicOrderId },
    { label: 'Курьер', value: currentCourierName.value },
    { label: 'ID курьера', value: currentCourierId.value },
    { label: 'Дата создания доставки', value: ticket.createdAt },
    {
      label: 'Время принятия в доставку',
      value: 'Будет установлено автоматически после подтверждения',
      dimmed: true,
      multiline: true
    },
    { label: 'Текущий статус заказа', value: 'Готов', tone: 'success' },
    { label: 'Новый статус после подтверждения', value: 'В пути', tone: 'info' }
  ]
})

const finishConfirmInfoRows = computed<DeliveryRow[]>(() => {
  const ticket = selectedConfirmTicket.value

  if (!ticket) return []

  return [
    { label: 'ID доставки', value: ticket.deliveryCode },
    { label: 'ID заказа', value: ticket.publicOrderId },
    { label: 'Курьер', value: currentCourierName.value },
    { label: 'ID курьера', value: currentCourierId.value },
    { label: 'Время принятия в доставку', value: ticket.acceptedAt || 'не принято' },
    {
      label: 'Время завершения доставки',
      value: 'Будет установлено системой автоматически после подтверждения',
      dimmed: true,
      multiline: true
    },
    { label: 'Текущий статус заказа', value: 'В пути', tone: 'info' },
    { label: 'Новый статус', value: 'Доставлен', tone: 'success' },
    { label: 'Альтернативный статус', value: 'Не доставлен', tone: 'danger' }
  ]
})

const modalAddressRows = computed<DeliveryRow[]>(() => {
  const ticket = selectedConfirmTicket.value

  if (!ticket) return []

  return [
    { label: 'Город', value: ticket.address.city },
    { label: 'Улица', value: ticket.address.street },
    { label: 'Дом', value: ticket.address.house },
    { label: 'Квартира', value: ticket.address.apartment },
    { label: 'Этаж', value: ticket.address.floor },
    { label: 'Комментарий', value: ticket.address.comment, multiline: true }
  ]
})

const modalPaymentRows = computed<DeliveryRow[]>(() => {
  const ticket = selectedConfirmTicket.value

  if (!ticket) return []

  return [
    { label: 'Итоговая сумма', value: formatCurrency(ticket.total) },
    { label: 'Способ оплаты', value: ticket.paymentMethod },
    { label: 'Статус оплаты', value: paymentStatusLabel(ticket.paymentStatus), tone: paymentTone(ticket.paymentStatus) }
  ]
})

const openActionModal = (ticket: DeliveryTicket) => {
  if (ticket.status === 'ready') {
    pendingAction.value = { ticketId: ticket.id, type: 'start' }
  } else if (ticket.status === 'delivering') {
    pendingAction.value = { ticketId: ticket.id, type: 'finish' }
    finishPaymentReceived.value = ticket.paymentStatus === 'paid'
    finishCourierComment.value = ticket.courierComment
  } else {
    return
  }

  showConfirmModal.value = true
}

const closeConfirmModal = () => {
  showConfirmModal.value = false
  pendingAction.value = null
  finishPaymentReceived.value = false
  finishCourierComment.value = ''
}

const handleStartConfirm = () => {
  const action = pendingAction.value

  if (!action || action.type !== 'start') return

  const now = formatDateTime(new Date())

  tickets.value = tickets.value.map((ticket) =>
    ticket.id === action.ticketId
      ? {
          ...ticket,
          status: 'delivering',
          courierName: currentCourierName.value,
          courierId: currentCourierId.value,
          acceptedAt: now
        }
      : ticket
  )

  addToast('Доставка принята, заказ закреплён за текущим курьером', 'success')
  closeConfirmModal()
}

const handleFinishConfirm = (outcome: 'delivered' | 'failed') => {
  const action = pendingAction.value

  if (!action || action.type !== 'finish') return

  const now = formatDateTime(new Date())

  tickets.value = tickets.value.map((ticket) => {
    if (ticket.id !== action.ticketId) return ticket

    return {
      ...ticket,
      status: outcome,
      completedAt: now,
      courierName: currentCourierName.value,
      courierId: currentCourierId.value,
      paymentStatus: finishPaymentReceived.value ? 'paid' : ticket.paymentStatus,
      courierComment: finishCourierComment.value.trim()
    }
  })

  addToast(
    outcome === 'delivered'
      ? 'Доставка подтверждена'
      : 'Заказ отмечен как не доставленный',
    outcome === 'delivered' ? 'success' : 'info'
  )

  closeConfirmModal()
}

const actionButtonClass = (status: DeliveryTicketStatus) => {
  if (status === 'ready') {
    return 'bg-slate-800 text-white shadow-[0_18px_40px_-18px_rgba(15,23,42,0.8)] hover:-translate-y-0.5'
  }

  if (status === 'delivering') {
    return 'bg-slate-900 text-white shadow-[0_18px_40px_-18px_rgba(15,23,42,0.82)] hover:-translate-y-0.5'
  }

  if (status === 'delivered') {
    return 'cursor-default bg-emerald-100 text-emerald-700 shadow-none'
  }

  return 'cursor-default bg-rose-100 text-rose-700 shadow-none'
}

const actionButtonLabel = (status: DeliveryTicketStatus) => {
  if (status === 'ready') return 'Принять заказ'
  if (status === 'delivering') return 'Завершить доставку'
  if (status === 'delivered') return 'Доставка завершена'
  return 'Не доставлено'
}
</script>

<template>
  <div class="min-h-screen bg-[radial-gradient(circle_at_top,_rgba(255,255,255,0.98),_rgba(244,246,252,0.95)_46%,_rgba(236,240,247,0.92)_100%)]">
    <div class="mx-auto max-w-7xl px-4 py-8 md:px-6">
      <div class="mb-10 flex flex-col gap-6 xl:flex-row xl:items-center xl:justify-between">
        <div class="flex items-center gap-5">
          <div class="flex h-24 w-24 items-center justify-center rounded-[2rem] bg-gradient-to-br from-slate-800 via-slate-700 to-slate-900 text-white shadow-[0_24px_55px_-24px_rgba(15,23,42,0.9)]">
            <Truck class="h-12 w-12" />
          </div>

          <div>
            <p class="text-xs font-black uppercase tracking-[0.42em] text-base-content/30">Courier delivery module</p>
            <h1 class="text-5xl font-black uppercase tracking-[-0.08em] text-base-content">Доставка</h1>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-4 rounded-[2rem] border border-base-200/80 bg-base-100/90 px-4 py-3 shadow-[0_18px_45px_-30px_rgba(15,23,42,0.35)] backdrop-blur">
          <div class="rounded-full bg-sky-500/12 px-4 py-2 text-xs font-black uppercase tracking-[0.28em] text-sky-700">
            Доставка
          </div>

          <div class="h-10 w-px bg-base-200"></div>

          <div class="flex items-center gap-3">
            <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-slate-800/8 text-lg font-black text-slate-800 shadow-inner">
              {{ currentCourierInitial }}
            </div>

            <div>
              <p class="text-[11px] font-black uppercase tracking-[0.24em] text-base-content/30">Вы вошли как</p>
              <p class="text-base font-black text-base-content">Курьер {{ currentCourierName }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="grid gap-8 xl:grid-cols-2">
        <article
          v-for="ticket in sortedTickets"
          :key="ticket.id"
          class="rounded-[2.5rem] border border-base-200/80 bg-base-100/95 shadow-[0_28px_70px_-38px_rgba(15,23,42,0.4)] transition duration-300 hover:-translate-y-1 hover:shadow-[0_34px_80px_-38px_rgba(15,23,42,0.45)]"
        >
          <div class="flex items-center justify-between gap-4 border-b border-base-200/80 px-7 py-6">
            <h2 class="text-[2rem] font-black tracking-[-0.06em] text-base-content">ID доставки: {{ ticket.id }}</h2>

            <span
              class="rounded-full px-4 py-2 text-[11px] font-black uppercase tracking-[0.28em]"
              :class="statusBadgeClass(ticket.status)"
            >
              {{ statusBadgeLabel(ticket.status) }}
            </span>
          </div>

          <div class="space-y-5 px-7 py-6">
            <section class="overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white">
              <div class="flex items-center gap-4 border-b border-base-200/80 px-5 py-4">
                <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-sky-500/10 text-sky-600 shadow-inner">
                  <ReceiptText class="h-6 w-6" />
                </div>

                <div>
                  <p class="text-[10px] font-black uppercase tracking-[0.26em] text-base-content/32">Deliveries / couriers / order_statuses</p>
                </div>
              </div>

              <div class="grid gap-3 px-5 py-4">
                <div
                  v-for="row in buildDeliveryRows(ticket)"
                  :key="`${ticket.id}-${row.label}`"
                  class="grid grid-cols-[minmax(0,1fr)_minmax(0,0.95fr)] items-start gap-4 text-base"
                >
                  <span class="font-medium text-base-content/62">{{ row.label }}:</span>
                  <div class="text-right font-semibold text-base-content">
                    <span
                      v-if="row.tone"
                      class="inline-flex rounded-full px-3 py-1 text-sm font-bold"
                      :class="rowToneClass(row.tone)"
                    >
                      {{ row.value }}
                    </span>
                    <span
                      v-else
                      :class="[
                        row.dimmed ? 'text-base-content/40' : 'text-base-content',
                        row.multiline ? 'whitespace-pre-line' : ''
                      ]"
                    >
                      {{ row.value }}
                    </span>
                  </div>
                </div>
              </div>
            </section>

            <section class="overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white">
              <div class="flex items-center gap-4 border-b border-base-200/80 px-5 py-4">
                <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-rose-500/10 text-rose-600 shadow-inner">
                  <MapPin class="h-6 w-6" />
                </div>

                <p class="text-[10px] font-black uppercase tracking-[0.26em] text-base-content/32">Адрес доставки</p>
              </div>

              <div class="grid gap-3 px-5 py-4">
                <div
                  v-for="row in buildAddressRows(ticket)"
                  :key="`${ticket.id}-address-${row.label}`"
                  class="grid grid-cols-[minmax(0,1fr)_minmax(0,0.95fr)] items-start gap-4 text-base"
                >
                  <span class="font-medium text-base-content/62">{{ row.label }}:</span>
                  <span class="text-right font-semibold text-base-content" :class="row.multiline ? 'whitespace-pre-line' : ''">
                    {{ row.value }}
                  </span>
                </div>
              </div>
            </section>

            <section class="overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white">
              <div class="flex items-center gap-4 border-b border-base-200/80 px-5 py-4">
                <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-emerald-500/10 text-emerald-600 shadow-inner">
                  <Wallet class="h-6 w-6" />
                </div>

                <p class="text-[10px] font-black uppercase tracking-[0.26em] text-base-content/32">Оплата</p>
              </div>

              <div class="grid gap-3 px-5 py-4">
                <div
                  v-for="row in buildPaymentRows(ticket)"
                  :key="`${ticket.id}-payment-${row.label}`"
                  class="grid grid-cols-[minmax(0,1fr)_minmax(0,0.95fr)] items-start gap-4 text-base"
                >
                  <span class="font-medium text-base-content/62">{{ row.label }}:</span>
                  <div class="text-right font-semibold text-base-content">
                    <span
                      v-if="row.tone"
                      class="inline-flex rounded-full px-3 py-1 text-sm font-bold"
                      :class="rowToneClass(row.tone)"
                    >
                      {{ row.value }}
                    </span>
                    <span v-else>{{ row.value }}</span>
                  </div>
                </div>
              </div>
            </section>

            <section class="overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white">
              <div class="flex items-center gap-4 border-b border-base-200/80 px-5 py-4">
                <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-amber-500/10 text-amber-600 shadow-inner">
                  <Package class="h-6 w-6" />
                </div>

                <p class="text-[10px] font-black uppercase tracking-[0.26em] text-base-content/32">Состав заказа</p>
              </div>

              <div class="grid grid-cols-[1.4fr_0.6fr_0.6fr] px-5 py-3 text-[11px] font-black uppercase tracking-[0.18em] text-base-content/38">
                <span>Товар</span>
                <span class="text-center">Кол-во</span>
                <span class="text-right">Цена</span>
              </div>

              <div
                v-for="item in ticket.items"
                :key="`${ticket.id}-${item.name}`"
                class="grid grid-cols-[1.4fr_0.6fr_0.6fr] items-center border-t border-base-200/80 px-5 py-3 text-base"
              >
                <span class="font-medium text-base-content/84">{{ item.name }}</span>
                <span class="text-center font-semibold text-base-content">{{ item.quantity }}</span>
                <span class="text-right font-semibold text-base-content">{{ formatCurrency(item.price) }}</span>
              </div>
            </section>
          </div>

          <div class="flex items-center gap-4 px-7 pb-7">
            <button
              class="flex h-16 flex-1 items-center justify-center gap-3 rounded-[1.4rem] text-lg font-black uppercase tracking-tight transition"
              :class="actionButtonClass(ticket.status)"
              :disabled="ticket.status === 'delivered' || ticket.status === 'failed'"
              @click="openActionModal(ticket)"
            >
              <Truck class="h-6 w-6" />
              {{ actionButtonLabel(ticket.status) }}
            </button>

            <button
              class="flex h-16 w-16 items-center justify-center rounded-[1.4rem] border border-base-200 bg-white text-base-content/70 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md"
              @click="openChatModal(ticket.id)"
            >
              <MessageSquare class="h-7 w-7" />
            </button>
          </div>
        </article>
      </div>
    </div>

    <AppModal :show="showChatModal" maxWidth="lg" @close="closeChatModal">
      <div v-if="selectedChatTicket" class="max-h-[85vh] overflow-y-auto px-7 pb-7 pt-7 sm:px-8">
        <div class="mb-5 flex items-start justify-between gap-4">
          <div>
            <p class="text-[10px] font-black uppercase tracking-[0.24em] text-base-content/32">Чат по доставке</p>
            <h3 class="text-3xl font-black tracking-[-0.06em] text-base-content">{{ selectedChatTicket.id }}</h3>
          </div>

          <div class="rounded-full bg-slate-800/8 px-4 py-2 text-xs font-black uppercase tracking-[0.24em] text-slate-700">
            {{ selectedChatTicket.publicOrderId }}
          </div>
        </div>

        <div class="space-y-3 rounded-[1.8rem] border border-base-200/90 bg-base-100/80 p-4">
          <div
            v-for="message in selectedChatMessages"
            :key="`${selectedChatTicket.id}-${message.id}`"
            class="rounded-[1.35rem] px-4 py-3"
            :class="message.author === 'Курьер' ? 'bg-slate-800 text-white' : 'bg-white text-base-content shadow-sm'"
          >
            <div class="mb-2 flex items-center justify-between gap-3 text-xs font-black uppercase tracking-[0.2em]">
              <span :class="message.author === 'Курьер' ? 'text-white/70' : 'text-base-content/35'">{{ message.author }}</span>
              <span :class="message.author === 'Курьер' ? 'text-white/65' : 'text-base-content/35'">{{ message.time }}</span>
            </div>
            <p class="text-sm font-medium leading-relaxed">{{ message.text }}</p>
          </div>
        </div>

        <div class="mt-5 space-y-3">
            <textarea
              v-model="chatDraft"
              rows="3"
              class="textarea textarea-bordered w-full rounded-[1.35rem] border-base-200 bg-white px-4 py-3 text-base"
              placeholder="Напишите сообщение..."
            ></textarea>

          <div class="flex items-center justify-end gap-3">
            <button
              class="btn btn-ghost rounded-[1rem] px-5 text-sm font-black uppercase tracking-[0.18em] text-base-content/45"
              @click="closeChatModal"
            >
              Закрыть
            </button>

            <button
              class="btn rounded-[1rem] border-0 bg-slate-800 px-6 text-sm font-black uppercase tracking-[0.18em] text-white shadow-[0_18px_35px_-20px_rgba(15,23,42,0.8)]"
              @click="sendChatMessage"
            >
              Отправить
            </button>
          </div>
        </div>
      </div>
    </AppModal>

    <AppModal :show="showConfirmModal" maxWidth="2xl" @close="closeConfirmModal">
      <div v-if="selectedConfirmTicket" class="max-h-[88vh] overflow-y-auto px-7 pb-8 pt-7 sm:px-10 sm:pb-10">
        <div class="mx-auto mb-5 flex w-fit items-center rounded-full bg-base-200/70 px-5 py-3 text-sm font-semibold text-base-content/65 shadow-inner">
          Вы вошли как:
          <span class="ml-1.5 font-black text-base-content">Курьер</span>
        </div>

        <div class="mb-5 flex justify-center">
          <div
            class="flex h-24 w-24 items-center justify-center rounded-[2rem] shadow-inner"
            :class="isStartConfirm ? 'bg-slate-800/8 text-slate-800' : 'bg-base-200 text-slate-800'"
          >
            <Truck v-if="isStartConfirm" class="h-12 w-12" />
            <AlertCircle v-else class="h-12 w-12" />
          </div>
        </div>

        <h3 class="mx-auto max-w-2xl text-center text-4xl font-black uppercase leading-none tracking-[-0.08em] text-base-content sm:text-5xl">
          {{ isStartConfirm ? 'Подтвердить начало доставки' : 'Завершение доставки' }}
        </h3>

        <p class="mx-auto mt-4 max-w-2xl text-center text-base font-medium leading-relaxed text-base-content/55">
          {{
            isStartConfirm
              ? 'Проверьте данные перед фиксацией начала логистического этапа.'
              : 'Подтвердите результат передачи заказа клиенту.'
          }}
        </p>

        <section class="mt-7 overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white shadow-[0_20px_48px_-34px_rgba(15,23,42,0.35)]">
          <div class="flex items-center gap-3 border-b border-base-200/80 px-5 py-4">
            <ReceiptText class="h-5 w-5 text-base-content/55" />
            <h4 class="text-xl font-black tracking-[-0.04em] text-base-content">
              {{ isStartConfirm ? 'Данные доставки' : 'Информация о доставке' }}
            </h4>
          </div>

          <div class="grid gap-3 px-5 py-4">
            <div
              v-for="row in isStartConfirm ? startConfirmInfoRows : finishConfirmInfoRows"
              :key="`${selectedConfirmTicket.id}-confirm-${row.label}`"
              class="grid grid-cols-[minmax(0,1fr)_minmax(0,1fr)] items-start gap-4 text-base"
            >
              <span class="font-medium text-base-content/62">{{ row.label }}</span>
              <div class="text-right font-semibold text-base-content">
                <span
                  v-if="row.tone"
                  class="inline-flex rounded-full px-3 py-1 text-sm font-bold"
                  :class="rowToneClass(row.tone)"
                >
                  {{ row.value }}
                </span>
                <span
                  v-else
                  :class="[
                    row.dimmed ? 'text-base-content/42' : 'text-base-content',
                    row.multiline ? 'whitespace-pre-line' : ''
                  ]"
                >
                  {{ row.value }}
                </span>
              </div>
            </div>
          </div>
        </section>

        <section class="mt-5 overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white shadow-[0_20px_48px_-34px_rgba(15,23,42,0.35)]">
          <div class="flex items-center gap-3 border-b border-base-200/80 px-5 py-4">
            <MapPin class="h-5 w-5 text-base-content/55" />
            <h4 class="text-xl font-black tracking-[-0.04em] text-base-content">
              {{ isStartConfirm ? 'Краткий адрес' : 'Адрес доставки' }}
            </h4>
          </div>

          <div class="grid gap-3 px-5 py-4">
            <div
              v-for="row in modalAddressRows"
              :key="`${selectedConfirmTicket.id}-modal-address-${row.label}`"
              class="grid grid-cols-[minmax(0,1fr)_minmax(0,1fr)] items-start gap-4 text-base"
            >
              <span class="font-medium text-base-content/62">{{ row.label }}</span>
              <span class="text-right font-semibold text-base-content" :class="row.multiline ? 'whitespace-pre-line' : ''">
                {{ row.value }}
              </span>
            </div>
          </div>
        </section>

        <section class="mt-5 overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white shadow-[0_20px_48px_-34px_rgba(15,23,42,0.35)]">
          <div class="flex items-center gap-3 border-b border-base-200/80 px-5 py-4">
            <CreditCard class="h-5 w-5 text-base-content/55" />
            <h4 class="text-xl font-black tracking-[-0.04em] text-base-content">Оплата</h4>
          </div>

          <div class="grid gap-3 px-5 py-4">
            <div
              v-for="row in modalPaymentRows"
              :key="`${selectedConfirmTicket.id}-modal-payment-${row.label}`"
              class="grid grid-cols-[minmax(0,1fr)_minmax(0,1fr)] items-start gap-4 text-base"
            >
              <span class="font-medium text-base-content/62">{{ row.label }}</span>
              <div class="text-right font-semibold text-base-content">
                <span
                  v-if="row.tone"
                  class="inline-flex rounded-full px-3 py-1 text-sm font-bold"
                  :class="rowToneClass(row.tone)"
                >
                  {{ row.value }}
                </span>
                <span v-else>{{ row.value }}</span>
              </div>
            </div>

            <label
              v-if="!isStartConfirm"
              class="mt-2 flex cursor-pointer items-center gap-3 rounded-[1rem] border border-base-200 bg-base-100 px-4 py-3"
            >
              <input v-model="finishPaymentReceived" type="checkbox" class="checkbox checkbox-sm rounded-md border-base-300" />
              <span class="text-base font-semibold text-base-content/72">Оплата получена</span>
            </label>
          </div>
        </section>

        <div
          v-if="isStartConfirm"
          class="mt-5 flex items-start gap-4 rounded-[1.35rem] border border-base-200/90 bg-base-100/85 px-5 py-4"
        >
          <div class="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-base-200 text-base-content/65">
            <AlertCircle class="h-5 w-5" />
          </div>

          <p class="text-sm font-medium leading-relaxed text-base-content/58">
            После подтверждения система закрепит заказ за текущим курьером и зафиксирует время принятия в доставку.
          </p>
        </div>

        <section
          v-else
          class="mt-5 overflow-hidden rounded-[1.75rem] border border-base-200/90 bg-white shadow-[0_20px_48px_-34px_rgba(15,23,42,0.35)]"
        >
          <div class="flex items-center gap-3 border-b border-base-200/80 px-5 py-4">
            <MessageSquare class="h-5 w-5 text-base-content/55" />
            <h4 class="text-xl font-black tracking-[-0.04em] text-base-content">Комментарий курьера</h4>
          </div>

          <div class="px-5 py-4">
            <textarea
              v-model="finishCourierComment"
              rows="3"
              class="textarea textarea-bordered w-full rounded-[1.2rem] border-base-200 bg-white px-4 py-3 text-base"
              placeholder="Передано клиенту лично"
            ></textarea>
          </div>
        </section>

        <div class="mt-8 flex flex-col gap-3">
          <button
            v-if="isStartConfirm"
            class="h-16 rounded-[1.45rem] bg-slate-800 text-lg font-black uppercase tracking-tight text-white shadow-[0_24px_48px_-22px_rgba(15,23,42,0.9)] transition hover:-translate-y-0.5"
            @click="handleStartConfirm"
          >
            Подтвердить
          </button>

          <template v-else>
            <button
              class="h-16 rounded-[1.45rem] bg-slate-900 text-lg font-black uppercase tracking-tight text-white shadow-[0_24px_48px_-22px_rgba(15,23,42,0.9)] transition hover:-translate-y-0.5"
              @click="handleFinishConfirm('delivered')"
            >
              Подтвердить доставку
            </button>

            <button
              class="h-16 rounded-[1.45rem] bg-gradient-to-r from-red-500 to-red-600 text-lg font-black uppercase tracking-tight text-white shadow-[0_24px_48px_-22px_rgba(239,68,68,0.9)] transition hover:-translate-y-0.5"
              @click="handleFinishConfirm('failed')"
            >
              Заказ не доставлен
            </button>
          </template>

          <button
            class="h-12 rounded-[1.2rem] border border-base-300 text-sm font-black uppercase tracking-[0.24em] text-base-content/55 transition hover:border-base-content/20 hover:text-base-content/75"
            @click="closeConfirmModal"
          >
            Отмена
          </button>
        </div>
      </div>
    </AppModal>
  </div>
</template>

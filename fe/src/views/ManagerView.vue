<script setup lang="ts">
import { computed, inject, ref } from 'vue'
import AppModal from '../components/shared/AppModal.vue'
import {
  ArrowRight,
  BadgeCheck,
  BadgePercent,
  CalendarDays,
  CheckCircle2,
  ChefHat,
  CircleX,
  Clock3,
  CreditCard,
  DollarSign,
  Download,
  MapPin,
  MessageSquare,
  PackageCheck,
  Phone,
  ReceiptText,
  ShoppingBag,
  Ticket,
  Timer,
  Trash2,
  TrendingUp,
  Truck,
  User,
  Wallet,
  X
} from 'lucide-vue-next'

type ToastType = 'success' | 'error' | 'info'
type DashboardIcon = typeof DollarSign
type OrderTone = 'success' | 'warning' | 'info' | 'error'
type DetailTone = 'success' | 'warning' | 'default'

interface MetricCard {
  label: string
  value: string
  icon: DashboardIcon
  tone: 'primary' | 'amber' | 'green' | 'blue' | 'rose'
  featured?: boolean
}

interface ProductRow {
  name: string
  sold: number
  revenue: number
}

interface TransactionRow {
  id: string
  date: string
  client: string
  orderStatus: 'Доставлен' | 'Создан' | 'На кухне' | 'В пути' | 'Отменен'
  paymentStatus: 'PAID' | 'CREATED' | 'REFUND'
  paymentMethod: 'Карта' | 'Наличные' | 'Онлайн'
  courier: string
  promoCode: string | null
  amount: number
}

interface ActiveOrderCard {
  id: string
  status: 'Доставляется' | 'В пути' | 'На кухне'
  tone: 'green' | 'blue' | 'amber'
  address: string
  time: string
  eta: string
  note: string
}

interface PromoCard {
  id: string
  code: string
  type: 'percent' | 'fixed'
  amount: number
  expiresAt: string
  note: string
}

interface OrderItemLine {
  id: number
  name: string
  quantity: number
  unitPrice: number
  lineTotal: number
}

interface OrderAddress {
  city: string
  street: string
  house: string
  apartment: string
  floor: string
  comment: string
}

interface OrderConversationMessage {
  id: number
  role: 'manager' | 'client' | 'courier'
  text: string
  time: string
  status: string
}

interface OrderProfile {
  key: string
  aliases: string[]
  badge: string
  internalId: string
  createdAt: string
  statusStep: number
  statusLabel: string
  statusPill: string
  statusTone: OrderTone
  client: string
  phone: string
  courier: string
  promoCode: string | null
  paymentMethod: string
  paymentResult: string
  paymentState: 'PAID' | 'CREATED' | 'REFUND'
  transactionId: string
  address: OrderAddress
  subtotal: number
  discount: number
  total: number
  items: OrderItemLine[]
  chat: OrderConversationMessage[]
}

interface DetailRow {
  label: string
  value: string
  icon: DashboardIcon
  tone?: DetailTone
}

const addToast = inject<(message: string, type?: ToastType) => void>('addToast', () => undefined)

const trackingSteps = [
  { label: 'Принят', icon: PackageCheck },
  { label: 'Оплачен', icon: CreditCard },
  { label: 'Готовится', icon: ChefHat },
  { label: 'Готов', icon: CheckCircle2 },
  { label: 'В пути', icon: Truck },
  { label: 'Доставлен', icon: BadgeCheck }
]

const formatCurrency = (value: number) => `${new Intl.NumberFormat('ru-RU').format(value)} ₽`
const formatDate = (value: string) => value.split('-').reverse().join('.')
const formatPromoValue = (promo: PromoCard) => promo.type === 'percent' ? `${promo.amount}%` : formatCurrency(promo.amount)

const exportDates = ref({
  start: '2026-01-01',
  end: '2026-04-05'
})

const transactionFilters = ref({
  orderStatus: 'Все',
  paymentStatus: 'Все',
  paymentMethod: 'Все',
  courier: 'Все',
  promo: 'Все',
  minAmount: '',
  maxAmount: ''
})

const metrics = ref<MetricCard[]>([
  { label: 'Выручка', value: '6 490 ₽', icon: DollarSign, tone: 'primary', featured: true },
  { label: 'Заказов', value: '3', icon: ShoppingBag, tone: 'green' },
  { label: 'Средний чек', value: '2 163 ₽', icon: TrendingUp, tone: 'amber' },
  { label: 'На кухне', value: '1', icon: ChefHat, tone: 'amber' },
  { label: 'В пути', value: '1', icon: Truck, tone: 'blue' },
  { label: 'Доставлено вовремя', value: '87%', icon: BadgeCheck, tone: 'green' },
  { label: 'Среднее время приготовления', value: '18 мин', icon: Timer, tone: 'rose' },
  { label: 'Среднее время доставки', value: '31 мин', icon: Clock3, tone: 'blue' },
  { label: 'Отменено / не доставлено', value: '1', icon: CircleX, tone: 'rose' }
])

const popularProducts = ref<ProductRow[]>([
  { name: 'Четыре сыра', sold: 4, revenue: 3000 },
  { name: 'Пепперони Фреш', sold: 4, revenue: 2600 },
  { name: 'Трюфельная с грибами', sold: 1, revenue: 890 },
  { name: 'Маргарита Буратта', sold: 3, revenue: 2250 }
])

const transactions = ref<TransactionRow[]>([
  {
    id: '#425e8a',
    date: '2026-04-01',
    client: 'Алексей Иванов',
    orderStatus: 'Доставлен',
    paymentStatus: 'PAID',
    paymentMethod: 'Карта',
    courier: 'Руслан',
    promoCode: 'SUMMER2026',
    amount: 2490
  },
  {
    id: '#9f7254',
    date: '2026-04-02',
    client: 'Вероника Просек',
    orderStatus: 'Доставлен',
    paymentStatus: 'PAID',
    paymentMethod: 'Онлайн',
    courier: 'Иван',
    promoCode: null,
    amount: 2600
  },
  {
    id: '#e20235',
    date: '2026-04-04',
    client: 'Клиент сайта',
    orderStatus: 'Создан',
    paymentStatus: 'CREATED',
    paymentMethod: 'Карта',
    courier: '—',
    promoCode: null,
    amount: 1900
  },
  {
    id: '#c77140',
    date: '2026-03-29',
    client: 'Мария Степанова',
    orderStatus: 'На кухне',
    paymentStatus: 'PAID',
    paymentMethod: 'Карта',
    courier: 'Антон',
    promoCode: 'NIGHT10',
    amount: 1740
  },
  {
    id: '#a99018',
    date: '2026-03-21',
    client: 'Егор Власов',
    orderStatus: 'Отменен',
    paymentStatus: 'REFUND',
    paymentMethod: 'Онлайн',
    courier: '—',
    promoCode: 'WELCOME',
    amount: 980
  }
])

const activeOrders = ref<ActiveOrderCard[]>([
  {
    id: '#425e8ab84d8d',
    status: 'Доставляется',
    tone: 'green',
    address: 'Проспект Ленина',
    time: '12:35',
    eta: '13:05',
    note: 'Клиент просит позвонить за 5 минут'
  },
  {
    id: '#9f7254f682a0',
    status: 'В пути',
    tone: 'blue',
    address: 'Проспект Вернадского',
    time: '12:47',
    eta: '13:20',
    note: 'Без приборов, передать на ресепшн'
  },
  {
    id: '#e202355a7685',
    status: 'На кухне',
    tone: 'amber',
    address: 'Проспект Вернадского',
    time: '12:50',
    eta: '13:30',
    note: 'Добавить соус ранч отдельно'
  },
  {
    id: '#30a07c600639',
    status: 'В пути',
    tone: 'blue',
    address: '213 Sadovaya',
    time: '12:55',
    eta: '13:18',
    note: 'Домофон не работает, встречают у входа'
  }
])

const promos = ref<PromoCard[]>([
  {
    id: 'promo-1',
    code: 'SUMMER2026',
    type: 'percent',
    amount: 2,
    expiresAt: '30.06.2026',
    note: 'Для повторных клиентов и доставки'
  },
  {
    id: 'promo-2',
    code: 'NIGHT10',
    type: 'percent',
    amount: 10,
    expiresAt: '15.07.2026',
    note: 'После 21:00 на весь онлайн-заказ'
  }
])

const orderProfiles = ref<OrderProfile[]>([
  {
    key: '#425e8a',
    aliases: ['#425e8ab84d8d'],
    badge: '#ORD-2026-04-01-425E8AB84D8D',
    internalId: '10248',
    createdAt: '01.04.2026 14:05',
    statusStep: 1,
    statusLabel: 'Оплачен',
    statusPill: 'В работе',
    statusTone: 'warning',
    client: 'Алексей Иванов',
    phone: '+7 (999) 123-45-67',
    courier: 'Алексей Смирнов',
    promoCode: 'SUMMER2026',
    paymentMethod: 'Банковская карта',
    paymentResult: 'Успешно',
    paymentState: 'PAID',
    transactionId: 'TXN-9F45-AB12-77C3',
    address: {
      city: 'Москва',
      street: 'Проспект Ленина',
      house: '12',
      apartment: '45',
      floor: '7',
      comment: 'Позвонить за 5 минут, код домофона 258.'
    },
    subtotal: 2740,
    discount: 250,
    total: 2490,
    items: [
      { id: 1, name: 'Пицца Четыре сыра', quantity: 2, unitPrice: 750, lineTotal: 1500 },
      { id: 2, name: 'Кока-Кола 0.5', quantity: 1, unitPrice: 150, lineTotal: 150 },
      { id: 3, name: 'Соус чесночный', quantity: 3, unitPrice: 200, lineTotal: 600 },
      { id: 4, name: 'Картофель по-деревенски', quantity: 1, unitPrice: 490, lineTotal: 490 }
    ],
    chat: [
      { id: 1, role: 'client', text: 'Здравствуйте, можно без лука?', time: '01.04.2026 14:06', status: 'Прочитано' },
      { id: 2, role: 'manager', text: 'Добрый день. Да, отметил в заказе.', time: '01.04.2026 14:07', status: 'Прочитано' },
      { id: 3, role: 'courier', text: 'Заказ заберу через 10 минут.', time: '01.04.2026 14:22', status: 'Не прочитано' }
    ]
  },
  {
    key: '#9f7254',
    aliases: ['#9f7254f682a0'],
    badge: '#ORD-2026-04-02-9F7254F682A0',
    internalId: '10251',
    createdAt: '02.04.2026 16:40',
    statusStep: 5,
    statusLabel: 'Доставлен',
    statusPill: 'Доставлен',
    statusTone: 'success',
    client: 'Вероника Просек',
    phone: '+7 (905) 446-22-18',
    courier: 'Иван Петров',
    promoCode: null,
    paymentMethod: 'Онлайн',
    paymentResult: 'Успешно',
    paymentState: 'PAID',
    transactionId: 'TXN-9F72-54AA-1001',
    address: {
      city: 'Москва',
      street: 'Проспект Вернадского',
      house: '44',
      apartment: '18',
      floor: '3',
      comment: 'Оставить на стойке администратора.'
    },
    subtotal: 2600,
    discount: 0,
    total: 2600,
    items: [
      { id: 1, name: 'Пепперони Фреш', quantity: 2, unitPrice: 650, lineTotal: 1300 },
      { id: 2, name: 'Сырные бортики', quantity: 1, unitPrice: 450, lineTotal: 450 },
      { id: 3, name: 'Лимонад манго', quantity: 2, unitPrice: 425, lineTotal: 850 }
    ],
    chat: [
      { id: 1, role: 'client', text: 'Оставьте заказ на стойке администратора.', time: '02.04.2026 16:41', status: 'Прочитано' },
      { id: 2, role: 'manager', text: 'Принято, передаем курьеру.', time: '02.04.2026 16:43', status: 'Прочитано' },
      { id: 3, role: 'courier', text: 'Уже на месте, передал заказ администратору.', time: '02.04.2026 17:12', status: 'Прочитано' }
    ]
  },
  {
    key: '#e20235',
    aliases: ['#e202355a7685'],
    badge: '#ORD-2026-04-04-E202355A7685',
    internalId: '10255',
    createdAt: '04.04.2026 12:48',
    statusStep: 2,
    statusLabel: 'Готовится',
    statusPill: 'На кухне',
    statusTone: 'warning',
    client: 'Клиент сайта',
    phone: '+7 (901) 245-11-03',
    courier: '—',
    promoCode: null,
    paymentMethod: 'Карта',
    paymentResult: 'Создано',
    paymentState: 'CREATED',
    transactionId: 'TXN-E202-355A-7685',
    address: {
      city: 'Москва',
      street: 'Проспект Вернадского',
      house: '51',
      apartment: '14',
      floor: '11',
      comment: 'Добавить соус ранч отдельно.'
    },
    subtotal: 1900,
    discount: 0,
    total: 1900,
    items: [
      { id: 1, name: 'Мясная BBQ', quantity: 1, unitPrice: 990, lineTotal: 990 },
      { id: 2, name: 'Сырный соус', quantity: 2, unitPrice: 190, lineTotal: 380 },
      { id: 3, name: 'Картофель фри', quantity: 1, unitPrice: 530, lineTotal: 530 }
    ],
    chat: [
      { id: 1, role: 'client', text: 'Можно добавить еще салфетки?', time: '04.04.2026 12:49', status: 'Прочитано' },
      { id: 2, role: 'manager', text: 'Да, уже передал на кухню.', time: '04.04.2026 12:50', status: 'Прочитано' }
    ]
  },
  {
    key: '#c77140',
    aliases: [],
    badge: '#ORD-2026-03-29-C7714019AA2C',
    internalId: '10234',
    createdAt: '29.03.2026 21:12',
    statusStep: 2,
    statusLabel: 'На кухне',
    statusPill: 'В работе',
    statusTone: 'warning',
    client: 'Мария Степанова',
    phone: '+7 (903) 555-14-10',
    courier: 'Антон',
    promoCode: 'NIGHT10',
    paymentMethod: 'Банковская карта',
    paymentResult: 'Успешно',
    paymentState: 'PAID',
    transactionId: 'TXN-C771-4019-AA2C',
    address: {
      city: 'Москва',
      street: 'Ломоносовский проспект',
      house: '6',
      apartment: '21',
      floor: '9',
      comment: 'Приборы на двоих.'
    },
    subtotal: 1930,
    discount: 190,
    total: 1740,
    items: [
      { id: 1, name: 'Трюфельная с грибами', quantity: 1, unitPrice: 890, lineTotal: 890 },
      { id: 2, name: 'Цезарь ролл', quantity: 2, unitPrice: 425, lineTotal: 850 },
      { id: 3, name: 'Соус ранч', quantity: 1, unitPrice: 190, lineTotal: 190 }
    ],
    chat: [
      { id: 1, role: 'client', text: 'Можно добавить приборы на двоих?', time: '29.03.2026 21:15', status: 'Прочитано' },
      { id: 2, role: 'manager', text: 'Да, передали на кухню.', time: '29.03.2026 21:16', status: 'Прочитано' }
    ]
  },
  {
    key: '#a99018',
    aliases: [],
    badge: '#ORD-2026-03-21-A99018FF7A11',
    internalId: '10192',
    createdAt: '21.03.2026 18:01',
    statusStep: 0,
    statusLabel: 'Отменен',
    statusPill: 'Отменен',
    statusTone: 'error',
    client: 'Егор Власов',
    phone: '+7 (925) 456-00-77',
    courier: '—',
    promoCode: 'WELCOME',
    paymentMethod: 'Онлайн',
    paymentResult: 'Возврат',
    paymentState: 'REFUND',
    transactionId: 'TXN-A990-18FF-7A11',
    address: {
      city: 'Москва',
      street: 'Кутузовский проспект',
      house: '9',
      apartment: '41',
      floor: '5',
      comment: 'Клиент отменил заказ.'
    },
    subtotal: 1180,
    discount: 200,
    total: 980,
    items: [
      { id: 1, name: 'Маргарита', quantity: 1, unitPrice: 590, lineTotal: 590 },
      { id: 2, name: 'Лимонад цитрус', quantity: 1, unitPrice: 190, lineTotal: 190 },
      { id: 3, name: 'Наггетсы', quantity: 1, unitPrice: 400, lineTotal: 400 }
    ],
    chat: [
      { id: 1, role: 'client', text: 'Извините, заказ уже не нужен.', time: '21.03.2026 18:05', status: 'Прочитано' },
      { id: 2, role: 'manager', text: 'Хорошо, отменили и оформили возврат.', time: '21.03.2026 18:07', status: 'Прочитано' }
    ]
  },
  {
    key: '#30a07c',
    aliases: ['#30a07c600639'],
    badge: '#ORD-2026-04-04-30A07C600639',
    internalId: '10259',
    createdAt: '04.04.2026 12:55',
    statusStep: 4,
    statusLabel: 'В пути',
    statusPill: 'В пути',
    statusTone: 'info',
    client: 'Дарья Серова',
    phone: '+7 (916) 700-18-22',
    courier: 'Руслан',
    promoCode: 'SUMMER2026',
    paymentMethod: 'Наличные',
    paymentResult: 'Успешно',
    paymentState: 'PAID',
    transactionId: 'TXN-30A0-7C60-0639',
    address: {
      city: 'Москва',
      street: 'Sadovaya',
      house: '213',
      apartment: '1',
      floor: '1',
      comment: 'Домофон не работает, встречают у входа.'
    },
    subtotal: 2150,
    discount: 110,
    total: 2040,
    items: [
      { id: 1, name: 'Карбонара', quantity: 2, unitPrice: 790, lineTotal: 1580 },
      { id: 2, name: 'Чесночный соус', quantity: 1, unitPrice: 140, lineTotal: 140 },
      { id: 3, name: 'Вода без газа', quantity: 2, unitPrice: 160, lineTotal: 320 }
    ],
    chat: [
      { id: 1, role: 'client', text: 'Домофон не работает, пожалуйста, наберите по телефону.', time: '04.04.2026 12:54', status: 'Прочитано' },
      { id: 2, role: 'manager', text: 'Отметил в заказе, курьер в курсе.', time: '04.04.2026 12:56', status: 'Прочитано' },
      { id: 3, role: 'courier', text: 'Подъезжаю, буду через 6 минут.', time: '04.04.2026 13:08', status: 'Не прочитано' }
    ]
  }
])

const showAddPromoModal = ref(false)
const showOrderModal = ref(false)
const selectedOrderKey = ref<string | null>(null)
const newPromo = ref({
  code: '',
  amount: 10,
  type: 'percent' as PromoCard['type'],
  expiresAt: '31.07.2026',
  note: 'Локальный промокод менеджера'
})
const messageDraft = ref('')

const resolveOrderProfile = (orderId: string) => {
  return orderProfiles.value.find((profile) => profile.key === orderId || profile.aliases.includes(orderId)) ?? null
}

const openOrderDetails = (orderId: string) => {
  const profile = resolveOrderProfile(orderId)
  if (!profile) {
    addToast('Для этого заказа пока оставил только визуальную строку', 'info')
    return
  }

  selectedOrderKey.value = profile.key
  showOrderModal.value = true
}

const closeOrderDetails = () => {
  showOrderModal.value = false
}

const filteredTransactions = computed(() => {
  const minAmount = transactionFilters.value.minAmount ? Number(transactionFilters.value.minAmount) : null
  const maxAmount = transactionFilters.value.maxAmount ? Number(transactionFilters.value.maxAmount) : null

  return transactions.value.filter((item) => {
    if (transactionFilters.value.orderStatus !== 'Все' && item.orderStatus !== transactionFilters.value.orderStatus) return false
    if (transactionFilters.value.paymentStatus !== 'Все' && item.paymentStatus !== transactionFilters.value.paymentStatus) return false
    if (transactionFilters.value.paymentMethod !== 'Все' && item.paymentMethod !== transactionFilters.value.paymentMethod) return false
    if (transactionFilters.value.courier !== 'Все' && item.courier !== transactionFilters.value.courier) return false

    if (transactionFilters.value.promo === 'Без промокода' && item.promoCode) return false
    if (transactionFilters.value.promo !== 'Все' && transactionFilters.value.promo !== 'Без промокода' && item.promoCode !== transactionFilters.value.promo) return false

    if (minAmount !== null && item.amount < minAmount) return false
    if (maxAmount !== null && item.amount > maxAmount) return false

    return true
  })
})

const filteredRevenue = computed(() => filteredTransactions.value.reduce((sum, item) => sum + item.amount, 0))

const courierOptions = computed(() => ['Все', ...new Set(transactions.value.map((item) => item.courier))])
const promoOptions = computed(() => ['Все', 'Без промокода', ...new Set(transactions.value.map((item) => item.promoCode).filter(Boolean) as string[])])

const selectedOrderProfile = computed<OrderProfile | null>(() => {
  if (!selectedOrderKey.value) return null
  return resolveOrderProfile(selectedOrderKey.value)
})

const selectedTrackerWidth = computed(() => {
  if (!selectedOrderProfile.value) return '0%'
  return `${(selectedOrderProfile.value.statusStep / (trackingSteps.length - 1)) * 100}%`
})

const selectedOrderInfoGroups = computed(() => {
  const order = selectedOrderProfile.value
  const empty = { left: [] as DetailRow[], right: [] as DetailRow[] }
  if (!order) return empty

  const left: DetailRow[] = [
    { label: 'ID заказа', value: order.internalId, icon: PackageCheck },
    { label: 'Дата создания', value: order.createdAt, icon: Clock3 },
    { label: 'Текущий статус', value: order.statusLabel, icon: BadgeCheck, tone: order.statusTone === 'error' ? 'warning' : 'success' },
    { label: 'Клиент', value: order.client, icon: User },
    { label: 'Телефон клиента', value: order.phone, icon: Phone }
  ]

  const right: DetailRow[] = [
    { label: 'Курьер', value: order.courier, icon: Truck },
    { label: 'Промокод', value: order.promoCode ?? '—', icon: Ticket },
    { label: 'Способ оплаты', value: order.paymentMethod, icon: CreditCard },
    { label: 'Статус оплаты', value: order.paymentResult, icon: Wallet, tone: order.paymentState === 'REFUND' ? 'warning' : 'success' },
    { label: 'ID транзакции', value: order.transactionId, icon: ReceiptText }
  ]

  return { left, right }
})

const selectedOrderAddressRows = computed(() => {
  const order = selectedOrderProfile.value
  if (!order) return []

  return [
    { label: 'Город', value: order.address.city },
    { label: 'Улица', value: order.address.street },
    { label: 'Дом', value: order.address.house },
    { label: 'Квартира', value: order.address.apartment },
    { label: 'Этаж', value: order.address.floor },
    { label: 'Комментарий к доставке', value: order.address.comment }
  ]
})

const selectedOrderMessages = computed(() =>
  (selectedOrderProfile.value?.chat ?? []).filter((message) => message.role !== 'client')
)

const handleExport = () => {
  const rows = [
    ['ID заказа', 'Дата', 'Клиент', 'Статус заказа', 'Оплата', 'Курьер', 'Промокод', 'Сумма'],
    ...filteredTransactions.value.map((item) => [
      item.id,
      formatDate(item.date),
      item.client,
      item.orderStatus,
      `${item.paymentStatus} ${item.paymentMethod}`,
      item.courier,
      item.promoCode ?? '—',
      String(item.amount)
    ])
  ]

  const csv = rows.map((row) => row.map((cell) => `"${cell}"`).join(';')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `manager-dashboard-${exportDates.value.start}-${exportDates.value.end}.csv`
  document.body.appendChild(link)
  link.click()
  link.remove()
  window.URL.revokeObjectURL(url)
  addToast('Экспорт собран локально в CSV', 'success')
}

const handleDeletePromo = (promoId: string) => {
  promos.value = promos.value.filter((promo) => promo.id !== promoId)
  addToast('Промокод удален локально', 'success')
}

const handleAddPromo = () => {
  if (!newPromo.value.code.trim()) {
    addToast('Введите код промокода', 'error')
    return
  }

  promos.value.unshift({
    id: `promo-${Date.now()}`,
    code: newPromo.value.code.trim().toUpperCase(),
    amount: Number(newPromo.value.amount),
    type: newPromo.value.type,
    expiresAt: newPromo.value.expiresAt,
    note: newPromo.value.note.trim() || 'Локальный промокод менеджера'
  })

  showAddPromoModal.value = false
  newPromo.value = {
    code: '',
    amount: 10,
    type: 'percent',
    expiresAt: '31.07.2026',
    note: 'Локальный промокод менеджера'
  }

  addToast('Промокод добавлен локально', 'success')
}

const sendOrderMessage = () => {
  const text = messageDraft.value.trim()
  if (!text || !selectedOrderProfile.value) return

  selectedOrderProfile.value.chat.push({
    id: Date.now(),
    role: 'manager',
    text,
    time: 'Сейчас',
    status: 'Не прочитано'
  })

  messageDraft.value = ''
}

const metricToneClass = (tone: MetricCard['tone'], featured?: boolean) => {
  if (featured) return 'metric-card-featured'

  switch (tone) {
    case 'green':
      return 'metric-card-soft metric-card-green'
    case 'amber':
      return 'metric-card-soft metric-card-amber'
    case 'blue':
      return 'metric-card-soft metric-card-blue'
    case 'rose':
      return 'metric-card-soft metric-card-rose'
    default:
      return 'metric-card-soft metric-card-primary'
  }
}

const metricIconClass = (tone: MetricCard['tone'], featured?: boolean) => {
  if (featured) return 'bg-white/18 text-white'

  switch (tone) {
    case 'green':
      return 'bg-success/12 text-success'
    case 'amber':
      return 'bg-warning/16 text-warning'
    case 'blue':
      return 'bg-info/12 text-info'
    case 'rose':
      return 'bg-error/12 text-error'
    default:
      return 'bg-primary/12 text-primary'
  }
}

const orderStatusClass = (status: TransactionRow['orderStatus']) => {
  switch (status) {
    case 'Доставлен':
      return 'bg-success/15 text-success'
    case 'Создан':
      return 'bg-info/15 text-info'
    case 'На кухне':
      return 'bg-warning/18 text-warning'
    case 'В пути':
      return 'bg-primary/15 text-primary'
    default:
      return 'bg-error/15 text-error'
  }
}

const paymentStatusClass = (status: TransactionRow['paymentStatus']) => {
  switch (status) {
    case 'PAID':
      return 'bg-success/15 text-success'
    case 'CREATED':
      return 'bg-info/15 text-info'
    default:
      return 'bg-error/15 text-error'
  }
}

const activeOrderBadgeClass = (tone: ActiveOrderCard['tone']) => {
  switch (tone) {
    case 'green':
      return 'bg-success/15 text-success'
    case 'amber':
      return 'bg-warning/20 text-warning'
    default:
      return 'bg-info/15 text-info'
  }
}

const orderTonePillClass = (tone: OrderTone) => {
  switch (tone) {
    case 'success':
      return 'bg-success/15 text-success'
    case 'info':
      return 'bg-info/15 text-info'
    case 'error':
      return 'bg-error/15 text-error'
    default:
      return 'bg-warning/20 text-warning'
  }
}

const detailValueToneClass = (tone?: string) => {
  switch (tone) {
    case 'success':
      return 'text-success'
    case 'warning':
      return 'text-error'
    default:
      return 'text-secondary'
  }
}

const trackerCircleClass = (index: number) => {
  const currentStep = selectedOrderProfile.value?.statusStep ?? 0
  return index <= currentStep
    ? 'bg-gradient-to-r from-primary to-[#ff8a53] text-white shadow-[0_18px_32px_-18px_rgba(255,71,87,0.9)]'
    : 'bg-base-200 text-secondary/45'
}

const trackerIconClass = (index: number) => {
  const currentStep = selectedOrderProfile.value?.statusStep ?? 0
  return index <= currentStep ? 'text-primary' : 'text-secondary/25'
}

const trackerLabelClass = (index: number) => {
  const currentStep = selectedOrderProfile.value?.statusStep ?? 0
  return index <= currentStep ? 'text-secondary' : 'text-secondary/45'
}

const messageCardClass = (role: OrderConversationMessage['role']) => {
  switch (role) {
    case 'client':
      return 'border-base-300 bg-base-100'
    case 'courier':
      return 'border-warning/20 bg-warning/8'
    default:
      return 'border-primary/12 bg-primary/6'
  }
}

const messageBadgeClass = (role: OrderConversationMessage['role']) => {
  switch (role) {
    case 'client':
      return 'bg-base-200 text-secondary/60'
    case 'courier':
      return 'bg-warning/16 text-warning'
    default:
      return 'bg-primary/14 text-primary'
  }
}

const messageRoleLabel = (role: OrderConversationMessage['role']) => {
  switch (role) {
    case 'client':
      return 'Клиент'
    case 'courier':
      return 'Курьер'
    default:
      return 'Менеджер'
  }
}
</script>

<template>
  <div class="manager-shell min-h-screen px-4 py-6 sm:px-6 lg:px-8">
    <div class="mx-auto max-w-[1440px] space-y-6">
      <section class="manager-surface p-6 sm:p-8">
        <div class="relative z-10 flex flex-col gap-6 xl:flex-row xl:items-start xl:justify-between">
          <div>
            <p class="text-xs font-black uppercase tracking-[0.28em] text-secondary/35">Панель менеджера</p>
            <h1 class="mt-3 text-4xl font-black uppercase tracking-tight text-secondary sm:text-5xl">Аналитический дашборд</h1>
            <p class="mt-2 text-sm font-semibold text-secondary/45">Продажи, доставка и эффективность за период</p>
          </div>

          <div class="grid gap-4 sm:grid-cols-[1fr_auto_1fr_auto] xl:min-w-[640px]">
            <label class="space-y-2">
              <span class="field-label">Дата начала</span>
              <div class="field-shell">
                <CalendarDays class="h-4 w-4 text-secondary/35" />
                <input v-model="exportDates.start" type="date" class="h-full w-full bg-transparent text-sm font-bold text-secondary outline-none" />
              </div>
            </label>

            <div class="hidden items-end pb-3 text-2xl font-black text-secondary/20 sm:flex">–</div>

            <label class="space-y-2">
              <span class="field-label">Дата окончания</span>
              <div class="field-shell">
                <CalendarDays class="h-4 w-4 text-secondary/35" />
                <input v-model="exportDates.end" type="date" class="h-full w-full bg-transparent text-sm font-bold text-secondary outline-none" />
              </div>
            </label>

            <button type="button" class="export-button self-end" @click="handleExport">
              <Download class="h-4 w-4" />
              Экспорт
            </button>
          </div>
        </div>
      </section>

      <section class="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
        <article
          v-for="metric in metrics"
          :key="metric.label"
          class="metric-card"
          :class="metricToneClass(metric.tone, metric.featured)"
        >
          <div class="flex items-center gap-4">
            <div class="metric-icon" :class="metricIconClass(metric.tone, metric.featured)">
              <component :is="metric.icon" class="h-5 w-5" />
            </div>

            <div>
              <p class="text-[11px] font-black uppercase tracking-[0.18em]" :class="metric.featured ? 'text-white/70' : 'text-secondary/40'">
                {{ metric.label }}
              </p>
              <h2 class="mt-1 text-3xl font-black tracking-tight" :class="metric.featured ? 'text-white' : 'text-secondary'">
                {{ metric.value }}
              </h2>
            </div>
          </div>
        </article>
      </section>

      <section class="grid gap-6 xl:grid-cols-[0.8fr_1.5fr]">
        <article class="manager-surface p-6 sm:p-8">
          <div class="relative z-10">
            <div class="mb-6 flex items-center justify-between gap-4">
              <div>
                <h2 class="text-2xl font-black uppercase tracking-tight text-secondary">Популярные товары</h2>
                <p class="mt-1 text-sm font-semibold text-secondary/45">Топ по продажам и выручке</p>
              </div>
              <div class="metric-icon bg-warning/15 text-warning">
                <TrendingUp class="h-5 w-5" />
              </div>
            </div>

            <div class="overflow-x-auto">
              <table class="table">
                <thead>
                  <tr class="text-[11px] font-black uppercase tracking-[0.2em] text-secondary/35">
                    <th>Название</th>
                    <th>Продано</th>
                    <th class="text-right">Выручка</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="product in popularProducts" :key="product.name" class="font-semibold text-secondary">
                    <td class="py-4 font-black">{{ product.name }}</td>
                    <td class="py-4">
                      <span class="inline-flex rounded-full bg-base-200 px-3 py-1 text-sm font-black text-secondary">{{ product.sold }}</span>
                    </td>
                    <td class="py-4 text-right font-black">{{ formatCurrency(product.revenue) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <button type="button" class="section-link mt-6" @click="addToast('Оставил блок товаров только как локальный обзор', 'info')">
              Смотреть все товары
              <ArrowRight class="h-4 w-4" />
            </button>
          </div>
        </article>

        <article class="manager-surface p-6 sm:p-8">
          <div class="relative z-10">
            <div class="mb-6 flex items-center justify-between gap-4">
              <div>
                <h2 class="text-2xl font-black uppercase tracking-tight text-secondary">Все транзакции за период</h2>
                <div class="mt-1 flex flex-wrap items-center gap-3 text-sm font-semibold text-secondary/45">
                  <span>{{ filteredTransactions.length }} записей · {{ formatCurrency(filteredRevenue) }}</span>
                  <span class="rounded-full bg-primary/10 px-3 py-1 text-[11px] font-black uppercase tracking-[0.18em] text-primary">Только моки фронта</span>
                </div>
              </div>
              <div class="metric-icon bg-primary/10 text-primary">
                <ReceiptText class="h-5 w-5" />
              </div>
            </div>

            <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
              <div>
                <span class="field-label">Статус заказа</span>
                <select v-model="transactionFilters.orderStatus" class="select-shell">
                  <option>Все</option>
                  <option>Создан</option>
                  <option>На кухне</option>
                  <option>В пути</option>
                  <option>Доставлен</option>
                  <option>Отменен</option>
                </select>
              </div>

              <div>
                <span class="field-label">Статус оплаты</span>
                <select v-model="transactionFilters.paymentStatus" class="select-shell">
                  <option>Все</option>
                  <option>PAID</option>
                  <option>CREATED</option>
                  <option>REFUND</option>
                </select>
              </div>

              <div>
                <span class="field-label">Способ оплаты</span>
                <select v-model="transactionFilters.paymentMethod" class="select-shell">
                  <option>Все</option>
                  <option>Карта</option>
                  <option>Наличные</option>
                  <option>Онлайн</option>
                </select>
              </div>

              <div>
                <span class="field-label">Курьер</span>
                <select v-model="transactionFilters.courier" class="select-shell">
                  <option v-for="option in courierOptions" :key="option">{{ option }}</option>
                </select>
              </div>

              <div>
                <span class="field-label">Промокод</span>
                <select v-model="transactionFilters.promo" class="select-shell">
                  <option v-for="option in promoOptions" :key="option">{{ option }}</option>
                </select>
              </div>

              <div>
                <span class="field-label">Сумма от</span>
                <div class="field-shell">
                  <Wallet class="h-4 w-4 text-secondary/35" />
                  <input v-model="transactionFilters.minAmount" type="number" placeholder="0" class="h-full w-full bg-transparent text-sm font-bold text-secondary outline-none" />
                </div>
              </div>

              <div>
                <span class="field-label">Сумма до</span>
                <div class="field-shell">
                  <Wallet class="h-4 w-4 text-secondary/35" />
                  <input v-model="transactionFilters.maxAmount" type="number" placeholder="10000" class="h-full w-full bg-transparent text-sm font-bold text-secondary outline-none" />
                </div>
              </div>

              <div>
                <span class="field-label">Итог</span>
                <div class="field-shell justify-between">
                  <div class="flex items-center gap-2">
                    <Wallet class="h-4 w-4 text-secondary/35" />
                    <span class="text-sm font-bold text-secondary/50">Сумма</span>
                  </div>
                  <span class="text-sm font-black text-secondary">{{ formatCurrency(filteredRevenue) }}</span>
                </div>
              </div>
            </div>

            <div class="mt-6 overflow-x-auto">
              <table class="table">
                <thead>
                  <tr class="text-[11px] font-black uppercase tracking-[0.2em] text-secondary/35">
                    <th>ID заказа</th>
                    <th>Дата</th>
                    <th>Клиент</th>
                    <th>Статус заказа</th>
                    <th>Оплата</th>
                    <th>Курьер</th>
                    <th>Промокод</th>
                    <th class="text-right">Сумма</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="transaction in filteredTransactions"
                    :key="transaction.id"
                    class="order-row-button font-semibold text-secondary"
                    role="button"
                    tabindex="0"
                    @click="openOrderDetails(transaction.id)"
                    @keydown.enter="openOrderDetails(transaction.id)"
                    @keydown.space.prevent="openOrderDetails(transaction.id)"
                  >
                    <td class="py-4 font-black">{{ transaction.id }}</td>
                    <td class="py-4">{{ formatDate(transaction.date) }}</td>
                    <td class="py-4">{{ transaction.client }}</td>
                    <td class="py-4">
                      <span class="rounded-full px-3 py-1 text-xs font-black" :class="orderStatusClass(transaction.orderStatus)">
                        {{ transaction.orderStatus }}
                      </span>
                    </td>
                    <td class="py-4">
                      <div class="flex flex-col gap-2">
                        <span class="inline-flex w-fit items-center gap-1 rounded-full px-3 py-1 text-xs font-black" :class="paymentStatusClass(transaction.paymentStatus)">
                          <CreditCard class="h-3.5 w-3.5" />
                          {{ transaction.paymentStatus }}
                        </span>
                        <span class="text-xs font-bold text-secondary/45">{{ transaction.paymentMethod }}</span>
                      </div>
                    </td>
                    <td class="py-4">{{ transaction.courier }}</td>
                    <td class="py-4">
                      <span v-if="transaction.promoCode" class="rounded-full bg-secondary/10 px-3 py-1 text-xs font-black text-secondary">
                        {{ transaction.promoCode }}
                      </span>
                      <span v-else class="text-secondary/35">—</span>
                    </td>
                    <td class="py-4 text-right font-black">{{ formatCurrency(transaction.amount) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div v-if="filteredTransactions.length === 0" class="mt-4 rounded-[1.5rem] border border-dashed border-base-300 bg-base-200/40 px-6 py-10 text-center text-sm font-black uppercase tracking-[0.2em] text-secondary/35">
              По текущим фильтрам ничего не найдено
            </div>

            <button type="button" class="section-link mt-6" @click="addToast('Для задачи оставил только визуальный список транзакций', 'info')">
              Смотреть все транзакции
              <ArrowRight class="h-4 w-4" />
            </button>
          </div>
        </article>
      </section>

      <section class="grid gap-6 xl:grid-cols-[1.7fr_0.95fr]">
        <article class="manager-surface p-6 sm:p-8">
          <div class="relative z-10">
            <div class="mb-6 flex items-center justify-between gap-4">
              <div>
                <h2 class="text-2xl font-black uppercase tracking-tight text-secondary">Активные заказы</h2>
                <p class="mt-1 text-sm font-semibold text-secondary/45">Клик по любой карточке открывает форму отслеживания статуса</p>
              </div>
              <div class="metric-icon bg-info/10 text-info">
                <PackageCheck class="h-5 w-5" />
              </div>
            </div>

            <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
              <button
                v-for="order in activeOrders"
                :key="order.id"
                type="button"
                class="order-clickable-card rounded-[1.75rem] border border-base-300 bg-base-100 px-5 py-5 shadow-sm"
                @click="openOrderDetails(order.id)"
              >
                <div class="flex items-start justify-between gap-3">
                  <span class="text-sm font-black text-secondary">{{ order.id }}</span>
                  <span class="rounded-full px-3 py-1 text-[11px] font-black" :class="activeOrderBadgeClass(order.tone)">
                    {{ order.status }}
                  </span>
                </div>

                <p class="mt-4 text-left text-sm font-bold text-secondary/70">{{ order.address }}</p>
                <p class="mt-1 text-left text-xs font-semibold text-secondary/40">Создан в {{ order.time }} · ETA {{ order.eta }}</p>
                <p class="mt-4 min-h-[40px] text-left text-sm font-semibold text-secondary/55">{{ order.note }}</p>

                <span class="chat-button mt-4">
                  <PackageCheck class="h-4 w-4" />
                  Открыть заказ
                </span>
              </button>
            </div>

            <button type="button" class="section-link mt-6" @click="addToast('Здесь оставил моковый обзор активных заказов', 'info')">
              Смотреть все заказы
              <ArrowRight class="h-4 w-4" />
            </button>
          </div>
        </article>

        <article class="manager-surface p-6 sm:p-8">
          <div class="relative z-10">
            <div class="mb-6 flex items-center justify-between gap-4">
              <div>
                <h2 class="text-2xl font-black uppercase tracking-tight text-secondary">Активные промокоды</h2>
                <p class="mt-1 text-sm font-semibold text-secondary/45">Без API, но с локальным добавлением и удалением</p>
              </div>
              <button type="button" class="metric-icon bg-primary/10 text-primary" @click="showAddPromoModal = true">
                <BadgePercent class="h-5 w-5" />
              </button>
            </div>

            <div class="space-y-4">
              <article v-for="promo in promos" :key="promo.id" class="rounded-[1.75rem] border border-base-300 bg-base-100 px-5 py-5 shadow-sm">
                <div class="flex items-start justify-between gap-3">
                  <div>
                    <div class="inline-flex items-center gap-2 rounded-full bg-secondary px-4 py-2 font-mono text-sm font-black tracking-[0.2em] text-secondary-content">
                      <Ticket class="h-4 w-4" />
                      {{ promo.code }}
                    </div>
                    <p class="mt-4 text-sm font-semibold text-secondary/70">Скидка: <span class="font-black text-secondary">{{ formatPromoValue(promo) }}</span></p>
                    <p class="mt-1 text-sm font-semibold text-secondary/50">{{ promo.note }}</p>
                  </div>

                  <button type="button" class="rounded-full p-2 text-error/50 transition hover:bg-error/10 hover:text-error" @click="handleDeletePromo(promo.id)">
                    <Trash2 class="h-4 w-4" />
                  </button>
                </div>

                <div class="mt-4 flex items-center justify-between text-sm font-semibold text-secondary/45">
                  <span>Действует до</span>
                  <span class="font-black text-secondary">{{ promo.expiresAt }}</span>
                </div>
              </article>
            </div>

            <button type="button" class="section-link mt-6" @click="addToast('Блок промокодов работает только на фронте', 'info')">
              Смотреть все промокоды
              <ArrowRight class="h-4 w-4" />
            </button>
          </div>
        </article>
      </section>
    </div>

    <AppModal :show="showAddPromoModal" maxWidth="lg" title="Добавить промокод" @close="showAddPromoModal = false">
      <div class="space-y-5 p-8">
        <label class="space-y-2">
          <span class="field-label">Код</span>
          <div class="field-shell">
            <Ticket class="h-4 w-4 text-secondary/35" />
            <input v-model="newPromo.code" type="text" placeholder="SPRING20" class="h-full w-full bg-transparent text-sm font-bold uppercase text-secondary outline-none" />
          </div>
        </label>

        <div class="grid gap-4 sm:grid-cols-2">
          <label class="space-y-2">
            <span class="field-label">Размер скидки</span>
            <div class="field-shell">
              <BadgePercent class="h-4 w-4 text-secondary/35" />
              <input v-model.number="newPromo.amount" type="number" min="1" class="h-full w-full bg-transparent text-sm font-bold text-secondary outline-none" />
            </div>
          </label>

          <label class="space-y-2">
            <span class="field-label">Тип</span>
            <select v-model="newPromo.type" class="select-shell">
              <option value="percent">Процент</option>
              <option value="fixed">Фикс</option>
            </select>
          </label>
        </div>

        <label class="space-y-2">
          <span class="field-label">Действует до</span>
          <div class="field-shell">
            <CalendarDays class="h-4 w-4 text-secondary/35" />
            <input v-model="newPromo.expiresAt" type="text" placeholder="31.07.2026" class="h-full w-full bg-transparent text-sm font-bold text-secondary outline-none" />
          </div>
        </label>

        <label class="space-y-2">
          <span class="field-label">Описание</span>
          <textarea v-model="newPromo.note" rows="3" class="textarea-shell" placeholder="Короткое описание промокода"></textarea>
        </label>

        <div class="flex justify-end gap-3">
          <button type="button" class="btn btn-ghost rounded-2xl px-6 font-black uppercase tracking-[0.16em]" @click="showAddPromoModal = false">Отмена</button>
          <button type="button" class="export-button" @click="handleAddPromo">Сохранить</button>
        </div>
      </div>
    </AppModal>

    <AppModal :show="showOrderModal" maxWidth="7xl" @close="closeOrderDetails">
      <div v-if="selectedOrderProfile" class="max-h-[88vh] overflow-y-auto p-6 sm:p-8">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <div class="flex flex-wrap items-center gap-3">
              <h2 class="text-3xl font-black uppercase tracking-tight text-secondary sm:text-4xl">Заказ</h2>
              <span class="rounded-full bg-gradient-to-r from-primary to-[#ff8a53] px-4 py-2 text-[11px] font-black uppercase tracking-[0.22em] text-white">
                {{ selectedOrderProfile.badge }}
              </span>
            </div>
            <p class="mt-2 text-sm font-semibold text-secondary/45">Карточка заказа и отслеживание в реальном времени</p>
          </div>

          <div class="flex items-center gap-3">
            <span class="rounded-full px-5 py-3 text-xs font-black uppercase tracking-[0.18em]" :class="orderTonePillClass(selectedOrderProfile.statusTone)">
              {{ selectedOrderProfile.statusPill }}
            </span>
            <button type="button" class="order-detail-close" @click="closeOrderDetails">
              <X class="h-4 w-4" />
            </button>
          </div>
        </div>

        <section class="mt-6 grid gap-6 xl:grid-cols-[1fr_1.1fr]">
          <article class="manager-surface p-6 sm:p-8">
            <div class="relative z-10">
              <div class="mb-6 flex items-center gap-3">
                <ReceiptText class="h-5 w-5 text-primary" />
                <h3 class="text-2xl font-black tracking-tight text-secondary">Заказ</h3>
              </div>

              <div class="grid gap-6 md:grid-cols-2">
                <div class="rounded-[1.5rem] border border-base-300 bg-base-100/70">
                  <div v-for="row in selectedOrderInfoGroups.left" :key="row.label" class="order-detail-row">
                    <div class="flex items-center gap-3 text-secondary/55">
                      <component :is="row.icon" class="h-4 w-4" />
                      <span class="text-sm font-semibold">{{ row.label }}</span>
                    </div>
                    <span class="order-detail-value" :class="detailValueToneClass(row.tone)">{{ row.value }}</span>
                  </div>
                </div>

                <div class="rounded-[1.5rem] border border-base-300 bg-base-100/70">
                  <div v-for="row in selectedOrderInfoGroups.right" :key="row.label" class="order-detail-row">
                    <div class="flex items-center gap-3 text-secondary/55">
                      <component :is="row.icon" class="h-4 w-4" />
                      <span class="text-sm font-semibold">{{ row.label }}</span>
                    </div>
                    <span class="order-detail-value" :class="detailValueToneClass(row.tone)">{{ row.value }}</span>
                  </div>
                </div>
              </div>
            </div>
          </article>

          <article class="manager-surface p-6 sm:p-8">
            <div class="relative z-10">
              <div class="mb-6 flex items-center gap-3">
                <Clock3 class="h-5 w-5 text-primary" />
                <h3 class="text-2xl font-black tracking-tight text-secondary">Статус заказа</h3>
              </div>

              <div class="relative pt-4">
                <div class="absolute left-0 right-0 top-9 hidden h-1 rounded-full bg-base-300 md:block"></div>
                <div class="absolute left-0 top-9 hidden h-1 rounded-full bg-gradient-to-r from-primary to-[#ff8a53] md:block" :style="{ width: selectedTrackerWidth }"></div>

                <div class="grid gap-6 md:grid-cols-6">
                  <div v-for="(step, index) in trackingSteps" :key="step.label" class="flex flex-col items-center gap-3 text-center">
                    <div class="relative z-10 flex h-10 w-10 items-center justify-center rounded-full border-4 border-white text-sm font-black" :class="trackerCircleClass(index)">
                      {{ index + 1 }}
                    </div>
                    <component :is="step.icon" class="h-5 w-5" :class="trackerIconClass(index)" />
                    <span class="text-sm font-black" :class="trackerLabelClass(index)">{{ step.label }}</span>
                  </div>
                </div>
              </div>
            </div>
          </article>
        </section>

        <section class="mt-6 grid gap-6 xl:grid-cols-[0.95fr_0.7fr_1.15fr]">
          <article class="manager-surface p-6 sm:p-8">
            <div class="relative z-10">
              <div class="mb-6 flex items-center gap-3">
                <MapPin class="h-5 w-5 text-primary" />
                <h3 class="text-2xl font-black tracking-tight text-secondary">Адрес доставки</h3>
              </div>

              <div class="rounded-[1.5rem] border border-base-300 bg-base-100/70">
                <div v-for="row in selectedOrderAddressRows" :key="row.label" class="order-detail-row">
                  <span class="text-sm font-semibold text-secondary/55">{{ row.label }}</span>
                  <span class="order-detail-value max-w-[58%] text-right">{{ row.value }}</span>
                </div>
              </div>
            </div>
          </article>

          <article class="manager-surface p-6 sm:p-8">
            <div class="relative z-10">
              <div class="mb-6 flex items-center gap-3">
                <Wallet class="h-5 w-5 text-primary" />
                <h3 class="text-2xl font-black tracking-tight text-secondary">Финансы</h3>
              </div>

              <div class="space-y-4 rounded-[1.5rem] border border-base-300 bg-base-100/70 p-5">
                <div class="flex items-center justify-between text-sm font-semibold text-secondary/55">
                  <span>Сумма до скидки</span>
                  <span class="font-black text-secondary">{{ formatCurrency(selectedOrderProfile.subtotal) }}</span>
                </div>
                <div class="flex items-center justify-between text-sm font-semibold text-secondary/55">
                  <span>Скидка</span>
                  <span class="font-black text-secondary">{{ formatCurrency(selectedOrderProfile.discount) }}</span>
                </div>
                <div class="h-px bg-base-300"></div>
                <div class="flex items-center justify-between">
                  <span class="text-base font-black text-secondary">Итоговая сумма</span>
                  <span class="text-3xl font-black tracking-tight text-primary">{{ formatCurrency(selectedOrderProfile.total) }}</span>
                </div>
              </div>
            </div>
          </article>

          <article class="manager-surface p-6 sm:p-8">
            <div class="relative z-10">
              <div class="mb-6 flex items-center justify-between gap-3">
                <div class="flex items-center gap-3">
                  <MessageSquare class="h-5 w-5 text-primary" />
                  <h3 class="text-2xl font-black tracking-tight text-secondary">Чат по заказу</h3>
                </div>
                <span class="rounded-full bg-success/12 px-3 py-1 text-[11px] font-black uppercase tracking-[0.18em] text-success">Локальный мок</span>
              </div>

              <div class="space-y-4">
                <article v-for="message in selectedOrderMessages" :key="message.id" class="rounded-[1.5rem] border px-4 py-4 shadow-sm" :class="messageCardClass(message.role)">
                  <div class="flex gap-4">
                    <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl" :class="messageBadgeClass(message.role)">
                      <Truck v-if="message.role === 'courier'" class="h-5 w-5" />
                      <MessageSquare v-else-if="message.role === 'manager'" class="h-5 w-5" />
                      <User v-else class="h-5 w-5" />
                    </div>

                    <div class="min-w-0 flex-1">
                      <div class="flex flex-wrap items-center justify-between gap-2">
                        <span class="text-sm font-black text-secondary">{{ messageRoleLabel(message.role) }}</span>
                        <span class="text-xs font-semibold text-secondary/40">{{ message.time }}</span>
                      </div>
                      <p class="mt-2 text-sm font-semibold leading-6 text-secondary/75">{{ message.text }}</p>
                      <p class="mt-3 text-xs font-bold text-secondary/40">{{ message.status }}</p>
                    </div>
                  </div>
                </article>
              </div>

              <div class="mt-5 flex gap-3">
                <input v-model="messageDraft" type="text" class="message-input" placeholder="Напишите сообщение..." @keydown.enter.prevent="sendOrderMessage" />
                <button type="button" class="message-send" @click="sendOrderMessage">Отправить</button>
              </div>
            </div>
          </article>
        </section>

        <section class="mt-6 manager-surface p-6 sm:p-8">
          <div class="relative z-10">
            <div class="mb-6 flex items-center gap-3">
              <ShoppingBag class="h-5 w-5 text-primary" />
              <h3 class="text-2xl font-black tracking-tight text-secondary">Состав заказа</h3>
            </div>

            <div class="overflow-x-auto">
              <table class="table">
                <thead>
                  <tr class="text-[11px] font-black uppercase tracking-[0.2em] text-secondary/35">
                    <th>ID позиции</th>
                    <th>Товар</th>
                    <th>Количество</th>
                    <th>Цена на момент заказа</th>
                    <th class="text-right">Сумма по позиции</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in selectedOrderProfile.items" :key="item.id" class="font-semibold text-secondary">
                    <td class="py-4">{{ item.id }}</td>
                    <td class="py-4 font-black">{{ item.name }}</td>
                    <td class="py-4">{{ item.quantity }}</td>
                    <td class="py-4">{{ formatCurrency(item.unitPrice) }}</td>
                    <td class="py-4 text-right font-black">{{ formatCurrency(item.lineTotal) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="mt-6 flex items-center justify-end gap-4 border-t border-base-300 pt-5">
              <span class="text-sm font-black uppercase tracking-[0.18em] text-secondary/45">Итого по заказу</span>
              <span class="text-3xl font-black tracking-tight text-primary">{{ formatCurrency(selectedOrderProfile.total) }}</span>
            </div>
          </div>
        </section>
      </div>
    </AppModal>
  </div>
</template>

<style scoped>
.manager-shell {
  background:
    radial-gradient(circle at top left, rgba(255, 71, 87, 0.12), transparent 26%),
    radial-gradient(circle at top right, rgba(56, 189, 248, 0.1), transparent 24%),
    linear-gradient(180deg, #f8f7f4 0%, #f2f4f8 100%);
}

.manager-surface {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 2rem;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.96), rgba(255, 255, 255, 0.82));
  box-shadow: 0 28px 80px -42px rgba(15, 23, 42, 0.28);
}

.manager-surface::after {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at top left, rgba(255, 255, 255, 0.95), transparent 42%);
  pointer-events: none;
}

.metric-card {
  border-radius: 1.75rem;
  padding: 1.35rem;
  border: 1px solid rgba(148, 163, 184, 0.15);
  box-shadow: 0 24px 50px -38px rgba(15, 23, 42, 0.3);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.metric-card:hover,
.order-clickable-card:hover,
.order-row-button:hover {
  transform: translateY(-2px);
}

.metric-card-soft {
  background: rgba(255, 255, 255, 0.9);
}

.metric-card-featured {
  background: linear-gradient(135deg, #ff4757, #ff7f50);
}

.metric-card-primary {
  box-shadow: 0 24px 50px -38px rgba(255, 71, 87, 0.38);
}

.metric-card-green {
  box-shadow: 0 24px 50px -38px rgba(34, 197, 94, 0.35);
}

.metric-card-amber {
  box-shadow: 0 24px 50px -38px rgba(245, 158, 11, 0.35);
}

.metric-card-blue {
  box-shadow: 0 24px 50px -38px rgba(59, 130, 246, 0.35);
}

.metric-card-rose {
  box-shadow: 0 24px 50px -38px rgba(244, 63, 94, 0.35);
}

.metric-icon {
  display: inline-flex;
  height: 3.25rem;
  width: 3.25rem;
  align-items: center;
  justify-content: center;
  border-radius: 1.2rem;
}

.field-label {
  display: inline-block;
  padding-left: 0.2rem;
  font-size: 0.68rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.2em;
  color: rgba(15, 23, 42, 0.38);
}

.field-shell,
.select-shell,
.textarea-shell {
  width: 100%;
  border: 1px solid rgba(148, 163, 184, 0.22);
  background: rgba(255, 255, 255, 0.82);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.field-shell {
  display: flex;
  height: 3.25rem;
  align-items: center;
  gap: 0.75rem;
  border-radius: 1.2rem;
  padding: 0 1rem;
}

.select-shell {
  height: 3.25rem;
  border-radius: 1.2rem;
  padding: 0 1rem;
  font-size: 0.95rem;
  font-weight: 700;
  color: rgb(15 23 42 / 0.9);
  outline: none;
}

.textarea-shell {
  min-height: 6rem;
  border-radius: 1.2rem;
  padding: 1rem;
  font-size: 0.95rem;
  font-weight: 600;
  color: rgb(15 23 42 / 0.8);
  outline: none;
}

.export-button {
  display: inline-flex;
  height: 3.25rem;
  align-items: center;
  justify-content: center;
  gap: 0.65rem;
  border-radius: 1.2rem;
  padding: 0 1.4rem;
  background: linear-gradient(135deg, #ff5a5f, #ff7f50);
  color: white;
  font-size: 0.76rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.2em;
  box-shadow: 0 20px 40px -24px rgba(255, 71, 87, 0.65);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.export-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 22px 42px -22px rgba(255, 71, 87, 0.75);
}

.section-link {
  display: inline-flex;
  align-items: center;
  gap: 0.55rem;
  font-size: 0.75rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.18em;
  color: rgba(15, 23, 42, 0.55);
  transition: color 0.2s ease, transform 0.2s ease;
}

.section-link:hover {
  color: rgb(255 90 95);
  transform: translateX(2px);
}

.order-row-button {
  cursor: pointer;
  transition: background-color 0.2s ease, transform 0.2s ease;
}

.order-row-button:hover {
  background: rgba(255, 90, 95, 0.04);
}

.order-clickable-card {
  cursor: pointer;
  text-align: left;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.order-clickable-card:hover {
  border-color: rgba(255, 90, 95, 0.28);
  box-shadow: 0 28px 50px -34px rgba(255, 71, 87, 0.25);
}

.chat-button {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  border-radius: 999px;
  background: rgba(255, 90, 95, 0.12);
  padding: 0.7rem 1rem;
  font-size: 0.72rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.16em;
  color: rgb(255 90 95);
}

.order-detail-row {
  display: flex;
  min-height: 4.25rem;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  border-bottom: 1px solid rgba(148, 163, 184, 0.14);
  padding: 0 1.25rem;
}

.order-detail-row:last-child {
  border-bottom: none;
}

.order-detail-value {
  font-size: 0.95rem;
  font-weight: 900;
}

.order-detail-close {
  display: inline-flex;
  height: 2.75rem;
  width: 2.75rem;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.06);
  color: rgba(15, 23, 42, 0.75);
  transition: background-color 0.2s ease, color 0.2s ease;
}

.order-detail-close:hover {
  background: rgba(255, 90, 95, 0.12);
  color: rgb(255 90 95);
}

.message-input {
  height: 3.35rem;
  flex: 1;
  border-radius: 999px;
  border: 1px solid rgba(148, 163, 184, 0.22);
  background: rgba(255, 255, 255, 0.84);
  padding: 0 1.2rem;
  font-size: 0.95rem;
  font-weight: 600;
  color: rgb(15 23 42 / 0.8);
  outline: none;
}

.message-send {
  display: inline-flex;
  min-width: 7.75rem;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: linear-gradient(135deg, #ff5a5f, #ff7f50);
  padding: 0 1.2rem;
  font-size: 0.75rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.18em;
  color: white;
  box-shadow: 0 18px 32px -22px rgba(255, 71, 87, 0.7);
}
</style>

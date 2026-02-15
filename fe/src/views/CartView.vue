<script setup lang="ts">
import { useCartStore } from '../store/cart'
import { useRouter } from 'vue-router'
import { Trash2, Minus, Plus, ShoppingBag, CreditCard, QrCode, X, MapPin, Building2, Hash, MessageSquare, Tag } from 'lucide-vue-next'
import { inject, ref, computed } from 'vue'
import { ordersApi } from '../api'
import { formatPrice, handleImageError } from '../utils/format'
import AppModal from '../components/shared/AppModal.vue'

import type { PromoCode } from '../types'

const cartStore = useCartStore()
const router = useRouter()
const addToast = inject('addToast') as (msg: string, type?: any) => void

const isPlacingOrder = ref(false)
const isPaying = ref(false)
const showPaymentModal = ref(false)
const createdOrderId = ref<string | null>(null)
const promoCode = ref('')
const promoApplied = ref(false)
const promoData = ref<PromoCode | null>(null)

const address = ref({ street: '', house: '', apartment: '', floor: '', entrance: '', comment: '' })

const isAddressValid = computed(() => address.value.street.length > 3 && address.value.house.length > 0)

const calculateTotalWithPromo = () => {
  if (!promoApplied.value || !promoData.value) return cartStore.totalPrice
  if (promoData.value.type === 'percent') {
    return cartStore.totalPrice * (1 - promoData.value.amount / 100)
  }
  return Math.max(0, cartStore.totalPrice - promoData.value.amount)
}

const applyPromo = async () => {
  if (!promoCode.value) return
  try {
    const res = await ordersApi.checkPromoCode(promoCode.value.toUpperCase())
    if (res.success && res.data) {
      promoData.value = res.data
      promoApplied.value = true
      addToast(`Промокод применен! Скидка ${promoData.value.amount}${promoData.value.type === 'percent' ? '%' : ' ₽'}`, 'success')
    } else {
      addToast('Неверный промокод', 'error')
    }
  } catch (err) {
    addToast('Промокод не найден', 'error')
  }
}

const handleCheckout = async () => {
  if (cartStore.items.length === 0 || !isAddressValid.value) return
  try {
    isPlacingOrder.value = true
    const res = await ordersApi.createOrder({
      items: cartStore.items.map(item => ({ product_id: item.id, quantity: item.quantity })),
      address: { city: 'Пицца-Сити', ...address.value },
      promo_code: promoApplied.value ? promoCode.value.toUpperCase() : ''
    })
    if (res.success && res.data) {
      createdOrderId.value = res.data.order_id
      showPaymentModal.value = true
    }
  } catch (error: any) {
    addToast(error.message || 'Ошибка при оформлении', 'error')
  } finally {
    isPlacingOrder.value = false
  }
}

const confirmPayment = async () => {
  if (!createdOrderId.value) return
  try {
    isPaying.value = true
    const res = await ordersApi.payOrder(createdOrderId.value)
    if (res.success) {
      cartStore.clearCart()
      addToast('Оплата успешна!', 'success')
      showPaymentModal.value = false
      router.push(`/order/${createdOrderId.value}/success`)
    }
  } catch (error: any) {
    addToast('Ошибка оплаты. Попробуйте еще раз.', 'error')
  } finally {
    isPaying.value = false
  }
}
</script>

<template>
  <div class="max-w-5xl mx-auto pb-20">
    <h1 class="text-5xl font-black mb-10 tracking-tight uppercase">Ваша корзина</h1>

    <div v-if="cartStore.items.length === 0" class="card bg-base-100 shadow-xl p-20 text-center border border-base-200 rounded-[3.5rem]">
      <div class="flex flex-col items-center gap-6">
        <div class="bg-base-200 p-10 rounded-full"><ShoppingBag class="w-20 h-20 text-base-content/10" /></div>
        <h2 class="text-3xl font-black">Корзина пуста</h2>
        <p class="text-base-content/50 max-w-sm font-medium">Похоже, вы еще не выбрали свою идеальную пиццу.</p>
        <router-link to="/" class="btn btn-primary btn-wide h-16 rounded-2xl font-black uppercase shadow-xl shadow-primary/20 mt-4">К меню</router-link>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-10">
      <div class="lg:col-span-8 space-y-10">
        <div class="space-y-4">
          <div v-for="item in cartStore.items" :key="item.id" class="card card-side bg-base-100 shadow-sm border border-base-200 p-5 transition-all hover:shadow-lg rounded-3xl overflow-hidden group">
            <figure class="w-28 h-28 rounded-2xl overflow-hidden flex-shrink-0 shadow-inner">
              <img :src="item.imageUrl" :alt="item.name" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" @error="handleImageError" />
            </figure>
            <div class="card-body py-0 px-6 justify-between">
              <div class="flex justify-between items-start">
                <div>
                  <h3 class="font-black text-xl leading-tight tracking-tight">{{ item.name }}</h3>
                  <p class="text-primary font-black text-lg mt-1">{{ formatPrice(item.price) }}</p>
                </div>
                <button @click="cartStore.removeFromCart(item.id)" class="btn btn-ghost btn-sm btn-circle text-error/30 hover:text-error hover:bg-error/10"><Trash2 class="w-5 h-5" /></button>
              </div>
              <div class="flex justify-between items-center mt-2">
                <div class="join bg-base-200/50 p-1 rounded-xl border border-base-300">
                  <button @click="cartStore.updateQuantity(item.id, item.quantity - 1)" class="btn btn-ghost btn-xs join-item px-3"><Minus class="w-4 h-4" /></button>
                  <span class="btn btn-ghost btn-xs join-item pointer-events-none font-black text-sm w-10">{{ item.quantity }}</span>
                  <button @click="cartStore.updateQuantity(item.id, item.quantity + 1)" class="btn btn-ghost btn-xs join-item px-3"><Plus class="w-4 h-4" /></button>
                </div>
                <p class="font-black text-xl tracking-tighter">{{ formatPrice(item.price * item.quantity) }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="card bg-base-100 shadow-xl border border-base-200 overflow-hidden rounded-[2.5rem]">
          <div class="bg-base-200/50 px-10 py-5 border-b border-base-200 flex items-center justify-between">
            <h2 class="font-black uppercase text-[10px] tracking-[0.3em] flex items-center gap-3 text-base-content/60"><MapPin class="w-4 h-4 text-primary" /> Адрес доставки</h2>
          </div>
          <div class="card-body p-10 space-y-8">
            <div class="form-control w-full">
              <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Улица</span></label>
              <div class="relative"><MapPin class="absolute left-4 top-3.5 w-5 h-5 opacity-20" />
                <input v-model="address.street" type="text" placeholder="пр. Ленина" class="input input-bordered w-full pl-12 rounded-2xl h-12 font-bold" />
              </div>
            </div>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
              <div class="form-control">
                <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Дом</span></label>
                <div class="relative"><Building2 class="absolute left-4 top-3.5 w-5 h-5 opacity-20" /><input v-model="address.house" type="text" class="input input-bordered w-full pl-12 rounded-2xl h-12 font-bold" /></div>
              </div>
              <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Кв/Офис</span></label><input v-model="address.apartment" type="text" class="input input-bordered w-full px-6 rounded-2xl h-12 font-bold" /></div>
              <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Этаж</span></label><input v-model="address.floor" type="text" class="input input-bordered w-full px-6 rounded-2xl h-12 font-bold" /></div>
              <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Подъезд</span></label><input v-model="address.entrance" type="text" class="input input-bordered w-full px-6 rounded-2xl h-12 font-bold" /></div>
            </div>
            <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Комментарий</span></label><div class="relative"><MessageSquare class="absolute left-4 top-4 w-5 h-5 opacity-20" /><textarea v-model="address.comment" class="textarea textarea-bordered w-full pl-12 rounded-2xl min-h-[120px] pt-4 font-bold"></textarea></div></div>
          </div>
        </div>
      </div>

      <div class="lg:col-span-4 space-y-6">
        <div class="card bg-base-100 shadow-2xl border border-primary/10 sticky top-24 overflow-hidden rounded-[2.5rem]">
          <div class="bg-primary/5 px-8 py-5 border-b border-primary/10 flex items-center gap-3">
            <CreditCard class="w-6 h-6 text-primary" /><h2 class="font-black uppercase text-[10px] tracking-widest">Итого к оплате</h2>
          </div>
          <div class="card-body p-8">
            <div class="form-control mb-6">
              <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Промокод</span></label>
              <div class="relative"><Tag class="absolute left-4 top-3.5 w-5 h-5 opacity-20" /><input v-model="promoCode" type="text" placeholder="CODE" class="input input-bordered w-full pl-12 pr-20 rounded-2xl h-12 font-black uppercase tracking-widest" :disabled="promoApplied" /><button @click="applyPromo" class="btn btn-sm btn-ghost absolute right-2 top-2 rounded-xl font-bold uppercase text-[10px] tracking-widest" :class="promoApplied ? 'text-success' : 'text-primary'" :disabled="promoApplied || !promoCode">{{ promoApplied ? 'OK' : 'Применить' }}</button></div>
            </div>
            <div class="space-y-4">
              <div class="flex justify-between items-center text-sm font-bold uppercase tracking-tighter opacity-40"><span>Сумма товаров</span><span>{{ formatPrice(cartStore.totalPrice) }}</span></div>
              <div class="flex justify-between items-center text-sm font-bold uppercase tracking-tighter opacity-40"><span>Доставка</span><span class="text-success">Бесплатно</span></div>
              <div v-if="promoApplied && promoData" class="flex justify-between items-center text-sm font-bold uppercase tracking-tighter text-secondary">
                <span>Скидка</span>
                <span>-{{ promoData.amount }}{{ promoData.type === 'percent' ? '%' : ' ₽' }}</span>
              </div>
              <div class="divider my-2 opacity-20"></div>
              <div class="flex justify-between items-end pt-2">
                <span class="font-black text-2xl uppercase tracking-tighter">Всего</span>
                <div class="text-right">
                  <p class="text-4xl font-black text-primary leading-none tracking-tighter">
                    {{ formatPrice(calculateTotalWithPromo()) }}
                  </p>
                </div>
              </div>
            </div>
            <div class="card-actions mt-10">
              <button @click="handleCheckout" class="btn btn-primary btn-block h-20 rounded-3xl shadow-2xl shadow-primary/30 font-black uppercase text-lg transition-all hover:scale-[1.02]" :disabled="isPlacingOrder || !isAddressValid">
                <template v-if="isPlacingOrder"><span class="loading loading-spinner"></span></template>
                <template v-else>Оформить заказ</template>
              </button>
              <p v-if="!isAddressValid && cartStore.totalItems > 0" class="text-[9px] text-error font-black uppercase text-center w-full mt-4 tracking-widest animate-pulse">Укажите улицу и номер дома</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <AppModal :show="showPaymentModal" @close="!isPaying && (showPaymentModal = false)">
      <div class="p-12 flex flex-col items-center">
        <div class="bg-primary/10 p-6 rounded-[2rem] mb-8"><QrCode class="w-14 h-14 text-primary" /></div>
        <h3 class="font-black text-4xl mb-2 text-center tracking-tighter uppercase">Оплата</h3>
        <p class="text-base-content/40 text-center mb-10 font-bold uppercase text-[10px] tracking-widest">К оплате: <span class="text-primary ml-1">{{ formatPrice(calculateTotalWithPromo()) }}</span></p>
        <div class="relative bg-white p-10 rounded-[3rem] shadow-inner mb-12 group overflow-hidden border-8 border-base-200 transition-all hover:border-primary/20"><svg viewBox="0 0 100 100" class="w-48 h-48 text-black fill-current relative z-10"><path d="M0 0h35v10H10v25H0V0zM65 0h35v35h-10V10H65V0zM0 65h10v25h25v10H0V65zM90 65h10v35H65v-10h25V65z" /><path d="M20 20h15v15H20V20zM65 20h15v15H65V20zM20 65h15v15H20V65z" /><path d="M42 20h16v16h-16V20zM20 42h16v16h-16V42zM42 42h16v16h-16V42zM64 42h16v16h-16V42zM42 64h16v16h-16V64zM64 64h8v8h-8V64zM72 72h8v8h-8V72z" /></svg></div>
        <div class="flex flex-col w-full gap-4">
          <button @click="confirmPayment" class="btn btn-primary btn-block h-20 rounded-3xl text-xl font-black uppercase shadow-2xl shadow-primary/30 tracking-tight transition-all hover:scale-[1.02]" :disabled="isPaying">
            <span v-if="isPaying" class="loading loading-spinner"></span><span v-else>Я оплатил</span>
          </button>
        </div>
      </div>
    </AppModal>
  </div>
</template>

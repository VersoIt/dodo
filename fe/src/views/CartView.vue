<script setup lang="ts">
import { useCartStore } from '../store/cart'
import { useRouter } from 'vue-router'
import { Trash2, Minus, Plus, ShoppingBag, CreditCard, QrCode, ShieldCheck, X, Check, MapPin, Building2, Hash, MessageSquare } from 'lucide-vue-next'
import axios from 'axios'
import { inject, ref, computed } from 'vue'

const cartStore = useCartStore()
const router = useRouter()
const isPlacingOrder = ref(false)
const isPaying = ref(false)
const showPaymentModal = ref(false)
const createdOrderId = ref<string | null>(null)
const addToast = inject('addToast') as (msg: string, type?: any) => void

// Address state
const address = ref({
  street: '',
  house: '',
  apartment: '',
  floor: '',
  comment: ''
})

const isAddressValid = computed(() => {
  return address.value.street.length > 3 && address.value.house.length > 0
})

const handleCheckout = async () => {
  if (cartStore.items.length === 0 || !isAddressValid.value) return
  
  try {
    isPlacingOrder.value = true
    const orderData = {
      items: cartStore.items.map(item => ({
        product_id: item.id,
        quantity: item.quantity
      })),
      address: {
        city: 'Пицца-Сити', // Default for now
        street: address.value.street,
        house: address.value.house,
        apartment: address.value.apartment,
        floor: address.value.floor,
        comment: address.value.comment
      }
    }
    
    const response = await axios.post('/api/v1/orders', orderData)
    if (response.data.success) {
      createdOrderId.value = response.data.data.order_id
      showPaymentModal.value = true
    } else {
      addToast(response.data.error || 'Не удалось оформить заказ', 'error')
    }
  } catch (error: any) {
    console.error('Checkout failed:', error)
    addToast(error.response?.data?.error || 'Ошибка при оформлении. Вы вошли в аккаунт?', 'error')
  } finally {
    isPlacingOrder.value = false
  }
}

const confirmPayment = async () => {
  if (!createdOrderId.value) return
  
  try {
    isPaying.value = true
    const response = await axios.post(`/api/v1/orders/${createdOrderId.value}/pay`)
    
    if (response.data.success) {
      cartStore.clearCart()
      addToast('Оплата прошла успешно!', 'success')
      showPaymentModal.value = false
      router.push(`/order/${createdOrderId.value}/success`)
    } else {
      addToast(response.data.error || 'Ошибка оплаты', 'error')
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
    <h1 class="text-4xl font-black mb-8 tracking-tight">Ваша корзина</h1>

    <div v-if="cartStore.items.length === 0" class="card bg-base-100 shadow-xl p-12 text-center border border-base-200">
      <div class="flex flex-col items-center gap-4">
        <div class="bg-base-200 p-6 rounded-full">
          <ShoppingBag class="w-16 h-16 text-base-content/20" />
        </div>
        <h2 class="text-2xl font-bold">Корзина пуста</h2>
        <p class="text-base-content/60">Похоже, вы еще не выбрали свою идеальную пиццу.</p>
        <router-link to="/" class="btn btn-primary btn-wide mt-4 rounded-xl">Перейти к меню</router-link>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-8">
      <!-- Items & Address -->
      <div class="lg:col-span-8 space-y-8">
        <!-- Items List -->
        <div class="space-y-4">
          <div v-for="item in cartStore.items" :key="item.id" class="card card-side bg-base-100 shadow-sm border border-base-200 p-4 transition-all hover:shadow-md">
            <figure class="w-24 h-24 rounded-2xl overflow-hidden flex-shrink-0 shadow-inner">
              <img :src="item.imageUrl" :alt="item.name" class="w-full h-full object-cover" />
            </figure>
            <div class="card-body py-0 px-4 justify-between">
              <div class="flex justify-between items-start">
                <div>
                  <h3 class="font-bold text-lg leading-tight">{{ item.name }}</h3>
                  <p class="text-primary font-black mt-1">{{ item.price?.toLocaleString() }} ₽</p>
                </div>
                <button @click="cartStore.removeFromCart(item.id)" class="btn btn-ghost btn-sm btn-circle text-error/40 hover:text-error hover:bg-error/10">
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
              
              <div class="flex justify-between items-center mt-2">
                <div class="join bg-base-200/50 p-0.5 rounded-xl border border-base-300">
                  <button @click="cartStore.updateQuantity(item.id, item.quantity - 1)" class="btn btn-ghost btn-xs join-item px-2">
                    <Minus class="w-3 h-3" />
                  </button>
                  <span class="btn btn-ghost btn-xs join-item pointer-events-none font-bold text-xs">{{ item.quantity }}</span>
                  <button @click="cartStore.updateQuantity(item.id, item.quantity + 1)" class="btn btn-ghost btn-xs join-item px-2">
                    <Plus class="w-3 h-3" />
                  </button>
                </div>
                <p class="font-black text-lg">{{ (item.price * item.quantity)?.toLocaleString() }} ₽</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Address Form -->
        <div class="card bg-base-100 shadow-xl border border-base-200 overflow-hidden rounded-3xl">
          <div class="bg-base-200/50 px-8 py-4 border-b border-base-200">
            <h2 class="font-black uppercase text-xs tracking-[0.2em] flex items-center gap-2">
              <MapPin class="w-4 h-4 text-primary" /> Адрес доставки
            </h2>
          </div>
          <div class="card-body p-8 space-y-6">
            <div class="form-control w-full">
              <label class="label"><span class="label-text font-bold">Улица *</span></label>
              <div class="relative">
                <MapPin class="absolute left-4 top-3.5 w-5 h-5 opacity-30" />
                <input v-model="address.street" type="text" placeholder="пр. Ленина" class="input input-bordered w-full pl-12 rounded-2xl h-12" required />
              </div>
            </div>

            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div class="form-control">
                <label class="label"><span class="label-text font-bold">Дом *</span></label>
                <div class="relative">
                  <Building2 class="absolute left-4 top-3.5 w-5 h-5 opacity-30" />
                  <input v-model="address.house" type="text" placeholder="10" class="input input-bordered w-full pl-12 rounded-2xl h-12" required />
                </div>
              </div>
              <div class="form-control">
                <label class="label"><span class="label-text font-bold">Кв/Офис</span></label>
                <div class="relative">
                  <Hash class="absolute left-4 top-3.5 w-5 h-5 opacity-30" />
                  <input v-model="address.apartment" type="text" placeholder="42" class="input input-bordered w-full pl-12 rounded-2xl h-12" />
                </div>
              </div>
              <div class="form-control">
                <label class="label"><span class="label-text font-bold">Этаж</span></label>
                <input v-model="address.floor" type="text" placeholder="5" class="input input-bordered w-full rounded-2xl h-12" />
              </div>
              <div class="form-control">
                <label class="label"><span class="label-text font-bold">Подъезд</span></label>
                <input v-model="address.floor" type="text" placeholder="2" class="input input-bordered w-full rounded-2xl h-12" />
              </div>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Комментарий курьеру</span></label>
              <div class="relative">
                <MessageSquare class="absolute left-4 top-3.5 w-5 h-5 opacity-30" />
                <textarea v-model="address.comment" class="textarea textarea-bordered w-full pl-12 rounded-2xl min-h-[100px] pt-3" placeholder="Например: вход со двора, код домофона 123..."></textarea>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Summary Card -->
      <div class="lg:col-span-4">
        <div class="card bg-base-100 shadow-2xl border border-primary/10 sticky top-24 overflow-hidden rounded-3xl">
          <div class="bg-primary/5 px-6 py-4 border-b border-primary/10">
            <h2 class="font-bold text-lg flex items-center gap-2"><CreditCard class="w-5 h-5 text-primary" /> Итого</h2>
          </div>
          <div class="card-body p-6">
            <div class="space-y-3">
              <div class="flex justify-between text-sm">
                <span class="opacity-60">Сумма товаров</span>
                <span class="font-bold">{{ cartStore.totalPrice?.toLocaleString() }} ₽</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="opacity-60">Доставка</span>
                <span class="text-success font-bold uppercase text-[10px] tracking-widest">Бесплатно</span>
              </div>
              <div class="divider my-1 opacity-50"></div>
              <div class="flex justify-between items-end pt-2">
                <span class="font-bold text-lg">Всего</span>
                <div class="text-right">
                  <p class="text-3xl font-black text-primary leading-none">{{ cartStore.totalPrice?.toLocaleString() }} ₽</p>
                  <p class="text-[10px] uppercase font-bold opacity-40 mt-1">НДС включен</p>
                </div>
              </div>
            </div>
            <div class="card-actions mt-8">
              <button 
                @click="handleCheckout" 
                class="btn btn-primary btn-block btn-lg rounded-2xl shadow-lg shadow-primary/30 gap-3 h-16 font-black uppercase"
                :disabled="isPlacingOrder || !isAddressValid"
              >
                <span v-if="isPlacingOrder" class="loading loading-spinner"></span>
                Оформить заказ
              </button>
              <p v-if="!isAddressValid && cartStore.totalItems > 0" class="text-[10px] text-error font-bold uppercase text-center w-full mt-2">
                Введите улицу и номер дома
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- PAYMENT MODAL -->
    <Transition name="modal-fade">
      <div v-if="showPaymentModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
        <div class="fixed inset-0 bg-black/80" @click="!isPaying && (showPaymentModal = false)"></div>
        <Transition name="modal-zoom" appear>
          <div class="relative bg-base-100 w-full max-w-md rounded-[2.5rem] shadow-2xl border border-white/5 overflow-hidden">
            <button v-if="!isPaying" @click="showPaymentModal = false" class="absolute top-6 right-6 btn btn-ghost btn-circle btn-sm bg-base-200/50"><X class="w-4 h-4" /></button>
            <div class="p-10 flex flex-col items-center">
              <div class="bg-primary/10 p-5 rounded-3xl mb-6"><QrCode class="w-12 h-12 text-primary" /></div>
              <h3 class="font-black text-3xl mb-2 text-center tracking-tighter uppercase italic">Оплата заказа</h3>
              <p class="text-base-content/50 text-center mb-8 font-medium">К оплате: <span class="text-primary font-black ml-1">{{ cartStore.totalPrice?.toLocaleString() }} ₽</span></p>
              <div class="relative bg-white p-8 rounded-[2.5rem] shadow-inner mb-10 group overflow-hidden border-4 border-base-200">
                <div class="absolute inset-0 bg-primary/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                <svg viewBox="0 0 100 100" class="w-44 h-44 text-black fill-current relative z-10">
                  <path d="M0 0h35v10H10v25H0V0zM65 0h35v35h-10V10H65V0zM0 65h10v25h25v10H0V65zM90 65h10v35H65v-10h25V65z" />
                  <path d="M20 20h15v15H20V20zM65 20h15v15H65V20zM20 65h15v15H20V65z" />
                  <path d="M42 20h16v16h-16V20zM20 42h16v16h-16V42zM42 42h16v16h-16V42zM64 42h16v16h-16V42zM42 64h16v16h-16V64zM64 64h8v8h-8V64zM72 72h8v8h-8V72z" />
                </svg>
              </div>
              <div class="flex flex-col w-full gap-3">
                <button @click="confirmPayment" class="btn btn-primary btn-block rounded-2xl h-16 text-lg font-black uppercase shadow-xl shadow-primary/20 tracking-tight" :disabled="isPaying">
                  <span v-if="isPaying" class="loading loading-spinner"></span>
                  <Check v-else class="w-6 h-6 mr-2" /> Я оплатил
                </button>
                <div class="flex items-center justify-center gap-2 mt-4 opacity-30">
                  <ShieldCheck class="w-4 h-4" />
                  <span class="text-[9px] font-black uppercase tracking-[0.2em]">Безопасный платеж</span>
                </div>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.15s ease; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
.modal-zoom-enter-active { transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1); }
.modal-zoom-leave-active { transition: all 0.15s ease-in; }
.modal-zoom-enter-from, .modal-zoom-leave-to { opacity: 0; transform: scale(0.97) translateY(8px); }
</style>

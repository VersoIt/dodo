<script setup lang="ts">
import { ref, onMounted, computed, inject } from 'vue'
import { Plus, PackagePlus, Image as ImageIcon, Tag, DollarSign, FileText, LogIn } from 'lucide-vue-next'
import { useCartStore } from '../store/cart'
import { useAuthStore } from '../store/auth'
import { catalogApi } from '../api'
import { CATEGORIES, CATEGORY_MAP, HERO_IMAGE } from '../constants'
import { formatPrice, handleImageError } from '../utils/format'
import AppModal from '../components/shared/AppModal.vue'

import type { Product } from '../types'

const cartStore = useCartStore()
const authStore = useAuthStore()
const addToast = inject('addToast') as (msg: string, type?: any) => void

const products = ref<Product[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const selectedCategory = ref('Все')
const menuSection = ref<HTMLElement | null>(null)

const showAddModal = ref(false)
const showAuthModal = ref(false)
const isSubmitting = ref(false)
const isEditing = ref(false)
const editingId = ref<string | null>(null)

const productForm = ref({ name: '', description: '', price: 0, imageUrl: '', categoryId: 0, isAvailable: true })

const filteredProducts = computed(() => {
  if (selectedCategory.value === 'Все') return products.value
  return products.value.filter(p => {
    // Protobuf JSON might omit 0 values due to omitempty.
    // If category_id is missing, it's 0 (Classic).
    const catId = p.category_id !== undefined ? p.category_id : 0
    const categoryName = CATEGORIES[catId + 1]
    return categoryName === selectedCategory.value
  })
})

const fetchProducts = async () => {
  try {
    loading.value = true
    const res = await catalogApi.getProducts()
    if (res.success && Array.isArray(res.data)) {
      products.value = res.data.map((p: any) => ({
        ...p,
        category_id: p.category_id ?? 0,
        image_url: p.image_url || HERO_IMAGE,
        is_available: p.is_available ?? true
      }))
    }
  } catch (err) { error.value = 'Не удалось загрузить товары.' } finally { loading.value = false }
}

const handleAddToCart = (product: Product) => {
  if (!authStore.isAuthenticated) { showAuthModal.value = true; return }
  // Map internal store format if needed, but for now assuming direct mapping
  cartStore.addToCart({
    id: product.id,
    name: product.name,
    price: product.price,
    imageUrl: product.image_url || HERO_IMAGE
  })
  addToast(`${product.name} добавлена в корзину!`, 'success')
}

const openAddModal = () => {
  isEditing.value = false; editingId.value = null
  productForm.value = { name: '', description: '', price: 0, imageUrl: '', categoryId: 0, isAvailable: true }
  showAddModal.value = true
}

const openEditModal = (product: Product) => {
  isEditing.value = true; editingId.value = product.id
  productForm.value = { 
    name: product.name, 
    description: product.description, 
    price: product.price, 
    imageUrl: product.image_url || '', 
    categoryId: product.category_id, 
    isAvailable: product.is_available 
  }
  showAddModal.value = true
}

const handleSubmit = async () => {
  try {
    isSubmitting.value = true
    const payload = { ...productForm.value, price: parseFloat(productForm.value.price.toString()) }
    const res = isEditing.value && editingId.value ? await catalogApi.updateProduct(editingId.value, payload) : await catalogApi.createProduct(payload)
    if (res.success) { addToast(isEditing.value ? 'Товар обновлен!' : 'Товар добавлен!', 'success'); showAddModal.value = false; await fetchProducts() }
  } catch (err: any) { addToast(err.response?.data?.error || 'Ошибка операции', 'error') } finally { isSubmitting.value = false }
}

onMounted(fetchProducts)
</script>

<template>
  <div class="space-y-12 pb-20">
    <!-- Premium Hero Section -->
    <section class="hero min-h-[45vh] rounded-[3rem] overflow-hidden shadow-2xl relative group">
      <div class="absolute inset-0 bg-cover bg-center transition-transform duration-[2s] group-hover:scale-105" :style="`background-image: url(${HERO_IMAGE});`"></div>
      <div class="absolute inset-0 bg-gradient-to-r from-black/90 via-black/50 to-transparent"></div>
      <div class="hero-content text-left text-neutral-content relative z-10 w-full justify-start px-12 md:px-20">
        <div class="max-w-2xl">
          <div v-if="authStore.user?.role === 'manager'" class="badge badge-accent mb-6 font-black tracking-widest p-4">ADMIN MODE</div>
          <h1 class="mb-6 text-6xl md:text-7xl font-black uppercase tracking-tighter text-white leading-none drop-shadow-xl">
            {{ authStore.user?.role === 'manager' ? 'Управление Меню' : 'Горячая. Свежая. Твоя.' }}
          </h1>
          <p class="mb-8 text-xl text-white/90 font-medium leading-relaxed max-w-lg drop-shadow-md">
            {{ authStore.user?.role === 'manager' ? 'Обновляйте каталог, добавляйте новинки и управляйте наличием в реальном времени.' : 'Итальянские традиции в каждом кусочке. Быстрая доставка прямо к двери.' }}
          </p>
          <div class="flex gap-4">
            <button v-if="authStore.user?.role !== 'manager'" @click="menuSection?.scrollIntoView({ behavior: 'smooth' })" class="btn btn-primary btn-lg h-16 rounded-[1.5rem] px-12 font-black uppercase tracking-widest shadow-xl shadow-primary/30 hover:shadow-primary/50 border-none hover:-translate-y-1 transition-all">Заказать сейчас</button>
            <button v-else @click="openAddModal" class="btn btn-secondary btn-lg h-16 rounded-[1.5rem] px-10 gap-3 font-black uppercase tracking-widest hover:-translate-y-1 transition-all"><Plus class="w-6 h-6" /> Добавить товар</button>
          </div>
        </div>
      </div>
    </section>

    <section ref="menuSection" class="pt-8">
      <div class="flex flex-col md:flex-row justify-between items-end mb-12 gap-6">
        <div class="flex flex-col gap-2">
          <span class="text-secondary/40 font-bold uppercase tracking-[0.2em] text-xs pl-1">Наше Меню</span>
          <h2 class="text-5xl font-black uppercase tracking-tighter text-secondary">Каталог</h2>
        </div>
        <div class="bg-base-100/50 backdrop-blur-md p-2 rounded-[2rem] shadow-soft border border-white/40 overflow-x-auto max-w-full">
          <div class="flex gap-2 min-w-max">
            <a v-for="cat in CATEGORIES" :key="cat" class="btn btn-sm h-10 px-6 rounded-full font-bold uppercase text-[10px] tracking-widest border-none transition-all duration-300" :class="selectedCategory === cat ? 'btn-primary shadow-lg shadow-primary/30 text-white' : 'btn-ghost text-secondary/60 hover:bg-base-200'" @click="selectedCategory = cat">{{ cat }}</a>
          </div>
        </div>
      </div>
      
      <div v-if="loading" class="flex justify-center py-32"><span class="loading loading-spinner loading-lg text-primary"></span></div>
      <div v-else-if="error" class="alert alert-error rounded-[2rem] shadow-xl">{{ error }}</div>
      
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        <div v-for="product in filteredProducts" :key="product.id" class="card bg-base-100 shadow-soft hover:shadow-glow transition-all duration-500 group rounded-[2.5rem] border border-transparent hover:border-primary/10 overflow-visible h-full hover:-translate-y-2">
          <figure class="relative h-64 overflow-hidden rounded-t-[2.5rem] m-2 mb-0 mask mask-squircle">
            <img :src="product.image_url" :alt="product.name" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700" @error="handleImageError" />
            <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent opacity-60"></div>
            <div class="absolute top-4 right-4 z-20"><span class="badge bg-white/90 backdrop-blur text-secondary font-black text-lg p-4 shadow-lg rounded-2xl">{{ formatPrice(product.price) }}</span></div>
            <div class="absolute bottom-4 left-4 z-20 flex flex-col gap-2 items-start">
              <span class="badge badge-primary border-none text-white text-[10px] font-black uppercase px-3 py-1 shadow-md">{{ CATEGORIES[product.category_id + 1] || 'Классика' }}</span>
              <span v-if="!product.is_available" class="badge badge-error text-white text-[10px] font-black uppercase px-3 py-1 shadow-md">SOLD OUT</span>
            </div>
          </figure>
          <div class="card-body p-8 pt-6 flex flex-col justify-between">
            <div>
              <h3 class="card-title text-2xl font-black tracking-tight leading-tight mb-2 group-hover:text-primary transition-colors">{{ product.name }}</h3>
              <p class="text-sm text-secondary/60 line-clamp-3 leading-relaxed font-medium">{{ product.description }}</p>
            </div>
            <div class="card-actions justify-end mt-6">
              <button v-if="authStore.user?.role !== 'manager'" @click="handleAddToCart(product)" class="btn btn-primary btn-block rounded-2xl gap-3 font-black uppercase shadow-xl shadow-primary/20 h-14 border-none hover:shadow-primary/40 group-hover:scale-[1.02] transition-all" :disabled="!product.is_available"><Plus class="w-5 h-5" /> В корзину</button>
              <button v-else @click="openEditModal(product)" class="btn btn-outline btn-block btn-secondary rounded-2xl gap-3 font-black uppercase h-14"><FileText class="w-5 h-5" /> Изменить</button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <AppModal :show="showAuthModal" @close="showAuthModal = false">
      <div class="p-10 flex flex-col items-center text-center">
        <div class="bg-primary/10 p-6 rounded-[2rem] mb-8"><LogIn class="w-12 h-12 text-primary" /></div>
        <h3 class="font-black text-3xl mb-3 uppercase tracking-tighter">Нужна авторизация</h3>
        <p class="text-base-content/50 mb-10 px-4 font-medium leading-relaxed">Вам нужно войти в аккаунт, чтобы делать заказы. Присоединяйтесь к нам!</p>
        <div class="flex flex-col w-full gap-3">
          <router-link to="/login" class="btn btn-primary btn-lg rounded-2xl font-black uppercase shadow-xl shadow-primary/20 h-16">Войти</router-link>
          <router-link to="/register" class="btn btn-outline btn-lg rounded-2xl font-bold border-2 h-16">Создать аккаунт</router-link>
        </div>
      </div>
    </AppModal>

    <AppModal :show="showAddModal" :title="isEditing ? 'Изменить товар' : 'Новый товар'" maxWidth="2xl" @close="showAddModal = false">
      <div class="p-8 md:p-10">
        <form @submit.prevent="handleSubmit" class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div class="space-y-5">
            <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Название</span></label><div class="relative"><Tag class="absolute left-4 top-3.5 w-5 h-5 opacity-20" /><input v-model="productForm.name" type="text" class="input input-bordered w-full pl-12 rounded-2xl h-12" required /></div></div>
            <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Цена (₽)</span></label><div class="relative"><span class="absolute left-4 top-3 font-bold opacity-20 text-lg">₽</span><input v-model="productForm.price" type="number" class="input input-bordered w-full pl-12 rounded-2xl h-12" required /></div></div>
            <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Категория</span></label><select v-model="productForm.categoryId" class="select select-bordered w-full rounded-2xl h-12 font-bold"><option v-for="(val, key) in CATEGORY_MAP" :key="key" :value="val">{{ key }}</option></select></div>
            <label class="label cursor-pointer justify-start gap-4 p-4 bg-base-200/50 rounded-2xl border border-base-300/50"><input v-model="productForm.isAvailable" type="checkbox" class="checkbox checkbox-primary rounded-lg" /><span class="label-text font-black uppercase text-[10px] tracking-widest">Доступен в меню</span></label>
          </div>
          <div class="space-y-5">
            <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">URL изображения</span></label><div class="relative"><ImageIcon class="absolute left-4 top-3.5 w-5 h-5 opacity-20" /><input v-model="productForm.imageUrl" type="text" class="input input-bordered w-full pl-12 rounded-2xl h-12" /></div></div>
            <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Описание</span></label><div class="relative"><FileText class="absolute left-4 top-3.5 w-5 h-5 opacity-20" /><textarea v-model="productForm.description" class="textarea textarea-bordered w-full pl-12 rounded-2xl min-h-[165px] pt-3"></textarea></div></div>
          </div>
          <div class="col-span-full mt-4 flex gap-4">
            <button type="submit" class="btn btn-primary flex-1 h-16 rounded-2xl font-black uppercase shadow-xl shadow-primary/20" :disabled="isSubmitting"><template v-if="isSubmitting"><span class="loading loading-spinner"></span></template><template v-else>{{ isEditing ? 'Сохранить изменения' : 'Опубликовать' }}</template></button>
            <button type="button" @click="showAddModal = false" class="btn btn-ghost h-16 rounded-2xl px-10 font-bold opacity-40">Отмена</button>
          </div>
        </form>
      </div>
    </AppModal>
  </div>
</template>

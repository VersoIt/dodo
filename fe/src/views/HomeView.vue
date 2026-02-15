<script setup lang="ts">
import { ref, onMounted, computed, inject } from 'vue'
import { Plus, PackagePlus, Image as ImageIcon, Tag, DollarSign, FileText, LogIn } from 'lucide-vue-next'
import { useCartStore } from '../store/cart'
import { useAuthStore } from '../store/auth'
import { catalogApi } from '../api'
import { CATEGORIES, CATEGORY_MAP, HERO_IMAGE } from '../constants'
import { formatPrice, handleImageError } from '../utils/format'
import AppModal from '../components/shared/AppModal.vue'

const cartStore = useCartStore()
const authStore = useAuthStore()
const addToast = inject('addToast') as (msg: string, type?: any) => void

// State
const products = ref<any[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const selectedCategory = ref('Все')
const menuSection = ref<HTMLElement | null>(null)

// Modal State
const showAddModal = ref(false)
const showAuthModal = ref(false)
const isSubmitting = ref(false)
const isEditing = ref(false)
const editingId = ref<string | null>(null)

const productForm = ref({
  name: '', description: '', price: 0, imageUrl: '', categoryId: 0, isAvailable: true
})

const filteredProducts = computed(() => {
  if (selectedCategory.value === 'Все') return products.value
  return products.value.filter(p => p.category === selectedCategory.value)
})

const fetchProducts = async () => {
  try {
    loading.value = true
    const data = await catalogApi.getProducts()
    const productsData = data.data.products || []
    
    products.value = productsData.map((p: any) => ({
      ...p,
      imageUrl: p.image_url || HERO_IMAGE,
      category: CATEGORIES[p.category_id + 1] || 'Классика',
      isAvailable: p.is_available ?? true
    }))
  } catch (err) {
    error.value = 'Не удалось загрузить товары.'
  } finally {
    loading.value = false
  }
}

const handleAddToCart = (product: any) => {
  if (!authStore.isAuthenticated) {
    showAuthModal.value = true
    return
  }
  cartStore.addToCart(product)
  addToast(`${product.name} добавлена в корзину!`, 'success')
}

const openAddModal = () => {
  isEditing.value = false
  editingId.value = null
  productForm.value = { name: '', description: '', price: 0, imageUrl: '', categoryId: 0, isAvailable: true }
  showAddModal.value = true
}

const openEditModal = (product: any) => {
  isEditing.value = true
  editingId.value = product.id
  productForm.value = {
    name: product.name,
    description: product.description,
    price: product.price,
    imageUrl: product.imageUrl,
    categoryId: product.category_id,
    isAvailable: product.isAvailable
  }
  showAddModal.value = true
}

const handleSubmit = async () => {
  try {
    isSubmitting.value = true
    const payload = {
      ...productForm.value,
      price: parseFloat(productForm.value.price.toString()),
      category_id: productForm.value.categoryId,
      image_url: productForm.value.imageUrl,
      is_available: productForm.value.isAvailable
    }

    const res = isEditing.value && editingId.value 
      ? await catalogApi.updateProduct(editingId.value, payload)
      : await catalogApi.createProduct(payload)

    if (res.success) {
      addToast(isEditing.value ? 'Товар обновлен!' : 'Товар добавлен!', 'success')
      showAddModal.value = false
      await fetchProducts()
    }
  } catch (err: any) {
    addToast(err.response?.data?.error || 'Ошибка операции', 'error')
  } finally {
    isSubmitting.value = false
  }
}

onMounted(fetchProducts)
</script>

<template>
  <div class="space-y-8 pb-20">
    <!-- Hero Section -->
    <section class="hero min-h-[40vh] rounded-[2.5rem] overflow-hidden shadow-2xl" :style="`background-image: url(${HERO_IMAGE});` ">
      <div class="hero-overlay bg-black bg-opacity-60"></div>
      <div class="hero-content text-center text-neutral-content">
        <div class="max-w-md">
          <h1 class="mb-5 text-5xl font-extrabold uppercase tracking-tighter text-white italic">
            {{ authStore.user?.role === 'manager' ? 'Управление Меню' : 'Горячая и Свежая' }}
          </h1>
          <p class="mb-5 text-lg text-white/80 font-medium">
            {{ authStore.user?.role === 'manager' 
               ? 'Обновляйте каталог, добавляйте новинки и управляйте наличием.' 
               : 'Лучшая пицца в городе, доставленная прямо к вашей двери.' }}
          </p>
          <button v-if="authStore.user?.role !== 'manager'" @click="menuSection?.scrollIntoView({ behavior: 'smooth' })" class="btn btn-primary btn-lg rounded-2xl px-10">Заказать сейчас</button>
          <button v-else @click="openAddModal" class="btn btn-secondary btn-lg gap-2 rounded-2xl px-10">
            <Plus class="w-6 h-6" /> Добавить товар
          </button>
        </div>
      </div>
    </section>

    <!-- Catalog Section -->
    <section ref="menuSection" class="pt-4">
      <div class="flex flex-col md:flex-row justify-between items-center mb-10 gap-6">
        <div class="flex items-center gap-4">
          <h2 class="text-4xl font-black uppercase tracking-tighter italic">Каталог</h2>
          <div v-if="authStore.user?.role === 'manager'" class="badge badge-secondary badge-lg font-bold px-4">ADMIN</div>
        </div>
        <div class="tabs tabs-boxed bg-base-200 p-1.5 rounded-2xl">
          <a 
            v-for="cat in CATEGORIES" :key="cat"
            class="tab rounded-xl transition-all font-bold uppercase text-[10px] tracking-widest px-6"
            :class="{ 'tab-active bg-primary text-primary-content shadow-lg shadow-primary/20': selectedCategory === cat }"
            @click="selectedCategory = cat"
          >
            {{ cat }}
          </a>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-32"><span class="loading loading-spinner loading-lg text-primary"></span></div>
      <div v-else-if="error" class="alert alert-error rounded-3xl">{{ error }}</div>
      
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        <div v-for="product in filteredProducts" :key="product.id" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-500 group rounded-[2rem] border border-base-200 overflow-hidden">
          <figure class="relative h-56 overflow-hidden">
            <img :src="product.imageUrl" :alt="product.name" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700" @error="handleImageError" />
            <div class="absolute top-4 right-4"><span class="badge badge-primary badge-lg font-black shadow-lg">{{ formatPrice(product.price) }}</span></div>
            <div class="absolute top-4 left-4 flex flex-col gap-2">
              <span class="badge badge-ghost bg-black/40 backdrop-blur-md text-white text-[9px] font-black uppercase border-none px-3 py-3">{{ product.category }}</span>
              <span v-if="!product.isAvailable" class="badge badge-error text-[9px] font-black uppercase px-3 py-3">SOLD OUT</span>
            </div>
          </figure>
          <div class="card-body p-8">
            <h3 class="card-title text-2xl font-black tracking-tight leading-tight">{{ product.name }}</h3>
            <p class="text-sm text-base-content/60 line-clamp-2 mt-2 leading-relaxed">{{ product.description }}</p>
            <div class="card-actions justify-end mt-8">
              <button v-if="authStore.user?.role !== 'manager'" @click="handleAddToCart(product)" class="btn btn-primary btn-block rounded-2xl gap-3 font-black uppercase shadow-lg shadow-primary/10 h-14" :disabled="!product.isAvailable">
                <Plus class="w-5 h-5" /> В корзину
              </button>
              <button v-else @click="openEditModal(product)" class="btn btn-outline btn-block btn-secondary rounded-2xl gap-3 font-black uppercase h-14">
                <FileText class="w-5 h-5" /> Изменить
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- AUTH MODAL -->
    <AppModal :show="showAuthModal" @close="showAuthModal = false">
      <div class="p-10 flex flex-col items-center text-center">
        <div class="bg-primary/10 p-6 rounded-[2rem] mb-8"><LogIn class="w-12 h-12 text-primary" /></div>
        <h3 class="font-black text-3xl mb-3 uppercase tracking-tighter">Нужна авторизация</h3>
        <p class="text-base-content/50 mb-10 px-4 font-medium leading-relaxed">Вам нужно войти в аккаунт, чтобы делать заказы. Присоединяйтесь к нам!</p>
        <div class="flex flex-col w-full gap-3">
          <router-link to="/login" class="btn btn-primary btn-lg rounded-2xl font-black uppercase shadow-xl shadow-primary/20 h-16">Войти</router-link>
          <router-link to="/register" class="btn btn-outline btn-lg rounded-2xl font-bold border-2 h-16">Создать аккаунт</router-link>
          <button @click="showAuthModal = false" class="btn btn-ghost btn-sm mt-4 opacity-40 font-bold">Может позже</button>
        </div>
      </div>
    </AppModal>

    <!-- ADD/EDIT MODAL -->
    <AppModal :show="showAddModal" :title="isEditing ? 'Изменить товар' : 'Новый товар'" maxWidth="2xl" @close="showAddModal = false">
      <div class="p-8 md:p-10">
        <form @submit.prevent="handleSubmit" class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div class="space-y-5">
            <div class="form-control">
              <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Название</span></label>
              <div class="relative"><Tag class="absolute left-4 top-3.5 w-5 h-5 opacity-20" />
                <input v-model="productForm.name" type="text" class="input input-bordered w-full pl-12 rounded-2xl h-12" required />
              </div>
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Цена (₽)</span></label>
              <div class="relative"><span class="absolute left-4 top-3 font-bold opacity-20 text-lg">₽</span>
                <input v-model="productForm.price" type="number" class="input input-bordered w-full pl-12 rounded-2xl h-12" required />
              </div>
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Категория</span></label>
              <select v-model="productForm.categoryId" class="select select-bordered w-full rounded-2xl h-12 font-bold">
                <option v-for="(val, key) in CATEGORY_MAP" :key="key" :value="val">{{ key }}</option>
              </select>
            </div>
            <label class="label cursor-pointer justify-start gap-4 p-4 bg-base-200/50 rounded-2xl border border-base-300/50">
              <input v-model="productForm.isAvailable" type="checkbox" class="checkbox checkbox-primary rounded-lg" />
              <span class="label-text font-black uppercase text-[10px] tracking-widest">Доступен в меню</span>
            </label>
          </div>
          <div class="space-y-5">
            <div class="form-control">
              <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">URL изображения</span></label>
              <div class="relative"><ImageIcon class="absolute left-4 top-3.5 w-5 h-5 opacity-20" />
                <input v-model="productForm.imageUrl" type="text" class="input input-bordered w-full pl-12 rounded-2xl h-12" />
              </div>
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40">Описание</span></label>
              <div class="relative"><FileText class="absolute left-4 top-3.5 w-5 h-5 opacity-20" />
                <textarea v-model="productForm.description" class="textarea textarea-bordered w-full pl-12 rounded-2xl min-h-[165px] pt-3"></textarea>
              </div>
            </div>
          </div>
          <div class="col-span-full mt-4 flex gap-4">
            <button type="submit" class="btn btn-primary flex-1 h-16 rounded-2xl font-black uppercase shadow-xl shadow-primary/20" :disabled="isSubmitting">
              <span v-if="isSubmitting" class="loading loading-spinner"></span>
              <span v-else>{{ isEditing ? 'Сохранить изменения' : 'Опубликовать' }}</span>
            </button>
            <button type="button" @click="showAddModal = false" class="btn btn-ghost h-16 rounded-2xl px-10 font-bold opacity-40">Отмена</button>
          </div>
        </form>
      </div>
    </AppModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Plus, PackagePlus, Image as ImageIcon, Tag, DollarSign, FileText, X, Check, LogIn } from 'lucide-vue-next'
import { useCartStore } from '../store/cart'
import { useAuthStore } from '../store/auth'
import { inject } from 'vue'

const router = useRouter()
const cartStore = useCartStore()
const authStore = useAuthStore()
const addToast = inject('addToast') as (msg: string, type?: any) => void

interface Product {
  id: string
  name: string
  description: string
  price: number
  imageUrl: string
  category: string
  categoryId: number
  isAvailable: boolean
}

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

const categories = ['Все', 'Классика', 'Премиум', 'Вегетарианская', 'Острая', 'Напитки', 'Десерты']
const categoryMap: Record<string, number> = {
  'Классика': 0,
  'Премиум': 1,
  'Вегетарианская': 2,
  'Острая': 3,
  'Напитки': 4,
  'Десерты': 5
}

const productForm = ref({
  name: '',
  description: '',
  price: 0,
  imageUrl: '',
  categoryId: 0,
  isAvailable: true
})

const filteredProducts = computed(() => {
  if (selectedCategory.value === 'Все') return products.value
  return products.value.filter(p => p.category === selectedCategory.value)
})

const scrollToMenu = () => {
  menuSection.value?.scrollIntoView({ behavior: 'smooth' })
}

const fetchProducts = async () => {
  try {
    loading.value = true
    const response = await axios.get('/api/v1/catalog/products')
    const productsData = response.data.data.products || []
    
    products.value = productsData.map((p: any) => ({
      id: p.id,
      name: p.name,
      description: p.description,
      price: p.price,
      imageUrl: p.image_url || `https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=500&auto=format&fit=crop&sig=${p.id}`,
      categoryId: p.category_id,
      category: categories[p.category_id + 1] || 'Классика',
      isAvailable: p.is_available ?? true
    }))
  } catch (err: any) {
    console.error('Error fetching products:', err)
    error.value = 'Не удалось загрузить товары.'
  } finally {
    loading.value = false
  }
}

const handleAddToCart = (product: Product) => {
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

const openEditModal = (product: Product) => {
  isEditing.value = true
  editingId.value = product.id
  productForm.value = {
    name: product.name,
    description: product.description,
    price: product.price,
    imageUrl: product.imageUrl,
    categoryId: product.categoryId,
    isAvailable: product.isAvailable
  }
  showAddModal.value = true
}

const handleSubmit = async () => {
  try {
    isSubmitting.value = true
    const payload = {
      name: productForm.value.name,
      description: productForm.value.description,
      price: parseFloat(productForm.value.price.toString()),
      category_id: productForm.value.categoryId,
      image_url: productForm.value.imageUrl,
      is_available: productForm.value.isAvailable
    }

    let response
    if (isEditing.value && editingId.value) {
      response = await axios.put(`/api/v1/catalog/products/${editingId.value}`, payload)
    } else {
      response = await axios.post('/api/v1/catalog/products', payload)
    }

    if (response.data.success) {
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
    <section class="hero min-h-[40vh] rounded-3xl overflow-hidden shadow-2xl" style="background-image: url(https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=1200&auto=format&fit=crop);">
      <div class="hero-overlay bg-black bg-opacity-60"></div>
      <div class="hero-content text-center text-neutral-content">
        <div class="max-w-md">
          <h1 class="mb-5 text-5xl font-extrabold uppercase tracking-tighter text-white">
            {{ authStore.user?.role === 'manager' ? 'Управление Меню' : 'Горячая и Свежая' }}
          </h1>
          <p class="mb-5 text-lg text-white/80">
            {{ authStore.user?.role === 'manager' 
               ? 'Обновляйте каталог, добавляйте новинки и управляйте наличием.' 
               : 'Лучшая пицца в городе, доставленная прямо к вашей двери.' }}
          </p>
          <button v-if="authStore.user?.role !== 'manager'" @click="scrollToMenu" class="btn btn-primary btn-lg">Заказать сейчас</button>
          <button v-else @click="openAddModal" class="btn btn-secondary btn-lg gap-2">
            <Plus class="w-6 h-6" /> Добавить товар
          </button>
        </div>
      </div>
    </section>

    <!-- Menu Section -->
    <section ref="menuSection">
      <div class="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
        <div class="flex items-center gap-3">
          <h2 class="text-3xl font-bold">Каталог</h2>
          <div v-if="authStore.user?.role === 'manager'" class="badge badge-secondary badge-outline">Режим Менеджера</div>
        </div>
        <div class="tabs tabs-boxed">
          <a 
            v-for="cat in categories" 
            :key="cat"
            class="tab"
            :class="{ 'tab-active': selectedCategory === cat }"
            @click="selectedCategory = cat"
          >
            {{ cat }}
          </a>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center items-center h-64">
        <span class="loading loading-spinner loading-lg text-primary"></span>
      </div>

      <div v-else-if="error" class="alert alert-error">
        <span>{{ error }}</span>
      </div>

      <div v-else-if="products.length === 0" class="text-center py-12">
        <p class="text-base-content/60">В этой категории пока нет товаров.</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div v-for="product in filteredProducts" :key="product.id" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 group">
          <figure class="relative h-48 overflow-hidden">
            <img :src="product.imageUrl" :alt="product.name" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" />
            <div class="absolute top-2 right-2">
              <span class="badge badge-secondary font-semibold">{{ product.price?.toLocaleString() }} ₽</span>
            </div>
            <div class="absolute top-2 left-2 flex flex-col gap-1">
              <span class="badge badge-ghost bg-black/40 text-white text-[10px] uppercase border-none">{{ product.category }}</span>
              <span v-if="!product.isAvailable" class="badge badge-error text-[10px] uppercase font-bold">Нет в наличии</span>
            </div>
          </figure>
          <div class="card-body p-6">
            <h3 class="card-title text-xl font-bold">{{ product.name }}</h3>
            <p class="text-sm text-base-content/70 line-clamp-2 mb-4">{{ product.description }}</p>
            <div class="card-actions justify-end mt-auto">
              <button 
                v-if="authStore.user?.role !== 'manager'"
                @click="handleAddToCart(product)" 
                class="btn btn-primary btn-sm gap-2"
                :disabled="!product.isAvailable"
              >
                <Plus class="w-4 h-4" />
                В корзину
              </button>
              <button 
                v-else
                @click="openEditModal(product)"
                class="btn btn-outline btn-sm btn-secondary gap-2"
              >
                <FileText class="w-4 h-4" /> Изменить
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- UNIFIED MODAL COMPONENT (RUSSIAN) -->
    <Transition name="modal-fade">
      <div v-if="showAuthModal || showAddModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
        <div 
          class="fixed inset-0 bg-black/70"
          @click="showAuthModal = false; showAddModal = false"
        ></div>
        
        <Transition name="modal-zoom" appear>
          <!-- Login Modal -->
          <div v-if="showAuthModal" class="relative bg-base-100 w-full max-w-md rounded-3xl border-t-4 border-primary shadow-2xl overflow-hidden">
            <div class="p-8 flex flex-col items-center text-center">
              <div class="bg-primary/10 p-4 rounded-full mb-6">
                <LogIn class="w-12 h-12 text-primary" />
              </div>
              <h3 class="font-bold text-3xl mb-2 tracking-tight">Нужна авторизация</h3>
              <p class="text-base-content/60 mb-8 px-2">
                Вам нужно войти в аккаунт, чтобы делать заказы. Присоединяйтесь к нашему сообществу любителей пиццы!
              </p>
              
              <div class="flex flex-col w-full gap-3">
                <router-link to="/login" class="btn btn-primary btn-lg btn-block gap-2 shadow-lg shadow-primary/20">
                  <LogIn class="w-5 h-5" /> Войти
                </router-link>
                <router-link to="/register" class="btn btn-outline btn-lg btn-block border-2">
                  Создать аккаунт
                </router-link>
                <button @click="showAuthModal = false" class="btn btn-ghost btn-sm mt-4 opacity-50 hover:opacity-100">
                  Может позже
                </button>
              </div>
            </div>
          </div>

          <!-- Add/Edit Product Modal -->
          <div v-else-if="showAddModal" class="relative bg-base-100 w-full max-w-2xl rounded-[2.5rem] shadow-2xl border border-base-200 overflow-hidden">
            <div class="p-8 md:p-10">
              <div class="flex justify-between items-center mb-8">
                <div class="flex items-center gap-3 text-secondary">
                  <PackagePlus class="w-8 h-8" />
                  <h3 class="font-black text-2xl uppercase tracking-tight">{{ isEditing ? 'Изменить товар' : 'Новый товар' }}</h3>
                </div>
                <button @click="showAddModal = false" class="btn btn-ghost btn-sm btn-circle bg-base-200/50"><X class="w-4 h-4" /></button>
              </div>
              
              <form @submit.prevent="handleSubmit" class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6">
                <div class="space-y-4">
                  <div class="form-control">
                    <label class="label"><span class="label-text font-bold uppercase text-[10px] opacity-50">Название товара</span></label>
                    <div class="relative">
                      <Tag class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                      <input v-model="productForm.name" type="text" placeholder="например, Пепперони Экстрим" class="input input-bordered w-full pl-10 rounded-xl" required />
                    </div>
                  </div>
                  <div class="form-control">
                    <label class="label"><span class="label-text font-bold uppercase text-[10px] opacity-50">Цена (₽)</span></label>
                    <div class="relative">
                      <span class="absolute left-3 top-2.5 font-bold opacity-40">₽</span>
                      <input v-model="productForm.price" type="number" step="1" placeholder="599" class="input input-bordered w-full pl-10 rounded-xl" required />
                    </div>
                  </div>
                  <div class="form-control">
                    <label class="label"><span class="label-text font-bold uppercase text-[10px] opacity-50">Категория</span></label>
                    <select v-model="productForm.categoryId" class="select select-bordered w-full rounded-xl">
                      <option v-for="(val, key) in categoryMap" :key="key" :value="val">{{ key }}</option>
                    </select>
                  </div>
                  <div class="form-control pt-2">
                    <label class="label cursor-pointer justify-start gap-4">
                      <input v-model="productForm.isAvailable" type="checkbox" class="checkbox checkbox-secondary checkbox-sm rounded-lg" />
                      <span class="label-text font-bold">Доступен для заказа</span>
                    </label>
                  </div>
                </div>
                <div class="space-y-4">
                  <div class="form-control">
                    <label class="label"><span class="label-text font-bold uppercase text-[10px] opacity-50">URL изображения</span></label>
                    <div class="relative">
                      <ImageIcon class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                      <input v-model="productForm.imageUrl" type="text" placeholder="https://..." class="input input-bordered w-full pl-10 rounded-xl" />
                    </div>
                  </div>
                  <div class="form-control">
                    <label class="label"><span class="label-text font-bold uppercase text-[10px] opacity-50">Описание</span></label>
                    <div class="relative">
                      <FileText class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                      <textarea v-model="productForm.description" class="textarea textarea-bordered w-full pl-10 h-32 rounded-xl" placeholder="Расскажите об этой вкусной пицце..."></textarea>
                    </div>
                  </div>
                </div>
                <div class="col-span-full mt-6 flex gap-3">
                  <button type="submit" class="btn btn-secondary flex-1 btn-lg rounded-2xl shadow-lg shadow-secondary/20 font-black uppercase" :disabled="isSubmitting">
                    <span v-if="isSubmitting" class="loading loading-spinner"></span>
                    {{ isEditing ? 'Сохранить изменения' : 'Опубликовать в меню' }}
                  </button>
                  <button type="button" @click="showAddModal = false" class="btn btn-ghost btn-lg rounded-2xl">Отмена</button>
                </div>
              </form>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.tab-active {
  @apply font-bold !important;
}
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.15s ease; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
.modal-zoom-enter-active { transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1); }
.modal-zoom-leave-active { transition: all 0.15s ease-in; }
.modal-zoom-enter-from, .modal-zoom-leave-to { opacity: 0; transform: scale(0.98) translateY(5px); }
</style>

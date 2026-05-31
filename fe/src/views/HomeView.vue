<script setup lang="ts">
import { computed, inject, onMounted, ref } from 'vue'
import { FileText, Image as ImageIcon, LogIn, Plus, Tag, X } from 'lucide-vue-next'
import { catalogApi } from '../api'
import AppModal from '../components/shared/AppModal.vue'
import { CATEGORY_MAP, CATEGORIES, HERO_IMAGE } from '../constants'
import { useAuthStore } from '../store/auth'
import { useCartStore } from '../store/cart'
import type { Product } from '../types'
import { formatPrice, handleImageError } from '../utils/format'

type ToastType = 'success' | 'error' | 'info'
type PublicationState = 'published' | 'draft'

interface ManagerProductMeta {
  unit: string
  cookingTime: number
  createdAt: string
  publicationState: PublicationState
}

interface ProductFormState {
  systemId: string
  name: string
  description: string
  price: number
  imageUrl: string
  categoryId: number
  unit: string
  cookingTime: number
  isAvailable: boolean
  createdAt: string
}

const cartStore = useCartStore()
const authStore = useAuthStore()
const addToast = inject<(msg: string, type?: ToastType) => void>('addToast', () => undefined)

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

const managerUnitOptions = ['шт.', 'порц.', 'бут.']
const managerProductMeta = ref<Record<string, ManagerProductMeta>>({})
const isManagerMode = computed(() => authStore.user?.role === 'manager')

const buildProductId = () => `PG-${Date.now().toString(36).toUpperCase()}`
const buildCreatedAt = () =>
  new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(new Date())

const buildEmptyProductForm = (): ProductFormState => ({
  systemId: buildProductId(),
  name: '',
  description: '',
  price: 1200,
  imageUrl: '',
  categoryId: 0,
  unit: 'шт.',
  cookingTime: 15,
  isAvailable: true,
  createdAt: buildCreatedAt()
})

const productForm = ref<ProductFormState>(buildEmptyProductForm())

const filteredProducts = computed(() => {
  const visibleProducts = products.value.filter((product) => (
    isManagerMode.value || managerProductMeta.value[product.id]?.publicationState !== 'draft'
  ))

  if (selectedCategory.value === 'Все') return visibleProducts

  return visibleProducts.filter((product) => {
    const categoryId = product.category_id !== undefined ? product.category_id : 0
    return CATEGORIES[categoryId + 1] === selectedCategory.value
  })
})

const fetchProducts = async () => {
  try {
    loading.value = true
    const res = await catalogApi.getProducts()

    if (res.success && Array.isArray(res.data)) {
      products.value = res.data.map((product: any) => ({
        ...product,
        category_id: product.category_id ?? 0,
        image_url: product.image_url || HERO_IMAGE,
        is_available: product.is_available ?? true
      }))
    }
  } catch (err) {
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

  cartStore.addToCart({
    id: product.id,
    name: product.name,
    price: product.price,
    imageUrl: product.image_url || HERO_IMAGE
  })

  addToast(`${product.name} добавлена в корзину!`, 'success')
}

const resolveProductMeta = (productId: string): ManagerProductMeta => (
  managerProductMeta.value[productId] || {
    unit: 'шт.',
    cookingTime: 15,
    createdAt: buildCreatedAt(),
    publicationState: 'published'
  }
)

const openAddModal = () => {
  isEditing.value = false
  editingId.value = null
  productForm.value = buildEmptyProductForm()
  showAddModal.value = true
}

const openEditModal = (product: Product) => {
  const meta = resolveProductMeta(product.id)

  isEditing.value = true
  editingId.value = product.id
  productForm.value = {
    systemId: product.id,
    name: product.name,
    description: product.description,
    price: product.price,
    imageUrl: product.image_url || '',
    categoryId: product.category_id ?? 0,
    unit: meta.unit,
    cookingTime: meta.cookingTime,
    isAvailable: product.is_available,
    createdAt: meta.createdAt
  }
  showAddModal.value = true
}

const isDraftProduct = (productId: string) => managerProductMeta.value[productId]?.publicationState === 'draft'

const saveProduct = async (publicationState: PublicationState) => {
  if (!productForm.value.name.trim()) {
    addToast('Заполните название товара.', 'error')
    return
  }

  try {
    isSubmitting.value = true

    const productId = editingId.value || productForm.value.systemId
    const normalizedProduct: Product = {
      id: productId,
      name: productForm.value.name.trim(),
      description: productForm.value.description.trim(),
      price: Number(productForm.value.price) || 0,
      category_id: productForm.value.categoryId,
      image_url: productForm.value.imageUrl.trim() || HERO_IMAGE,
      is_available: productForm.value.isAvailable
    }

    if (editingId.value) {
      products.value = products.value.map((product) => (
        product.id === editingId.value ? normalizedProduct : product
      ))
    } else {
      products.value = [normalizedProduct, ...products.value]
    }

    managerProductMeta.value[productId] = {
      unit: productForm.value.unit,
      cookingTime: Number(productForm.value.cookingTime) || 0,
      createdAt: productForm.value.createdAt,
      publicationState
    }

    showAddModal.value = false
    addToast(
      publicationState === 'draft'
        ? 'Черновик товара сохранен локально.'
        : isEditing.value
          ? 'Товар обновлен локально.'
          : 'Товар опубликован локально.',
      'success'
    )
  } finally {
    isSubmitting.value = false
  }
}

onMounted(fetchProducts)
</script>

<template>
  <div class="space-y-12 pb-20">
    <section class="hero relative min-h-[45vh] overflow-hidden rounded-[3rem] shadow-2xl group">
      <div class="absolute inset-0 bg-cover bg-center transition-transform duration-[2s] group-hover:scale-105" :style="`background-image: url(${HERO_IMAGE});`"></div>
      <div class="absolute inset-0 bg-gradient-to-r from-black/90 via-black/50 to-transparent"></div>
      <div class="hero-content relative z-10 w-full justify-start px-12 text-left text-neutral-content md:px-20">
        <div class="max-w-2xl">
          <div v-if="isManagerMode" class="badge badge-accent mb-6 p-4 font-black tracking-widest">ADMIN MODE</div>
          <h1 class="mb-6 text-6xl font-black uppercase leading-none tracking-tighter text-white drop-shadow-xl md:text-7xl">
            {{ isManagerMode ? 'Управление меню' : 'Горячая. Свежая. Твоя.' }}
          </h1>
          <p class="mb-8 max-w-lg text-xl font-medium leading-relaxed text-white/90 drop-shadow-md">
            {{ isManagerMode ? 'Обновляйте каталог, добавляйте новинки и ведите меню как плотный фронтовый mock-сценарий.' : 'Итальянские традиции в каждом кусочке. Быстрая доставка прямо к двери.' }}
          </p>
          <div class="flex gap-4">
            <button
              v-if="!isManagerMode"
              class="btn btn-primary btn-lg h-16 rounded-[1.5rem] border-none px-12 font-black uppercase tracking-widest shadow-xl shadow-primary/30 transition-all hover:-translate-y-1 hover:shadow-primary/50"
              @click="menuSection?.scrollIntoView({ behavior: 'smooth' })"
            >
              Заказать сейчас
            </button>
            <button
              v-else
              class="btn btn-secondary btn-lg h-16 gap-3 rounded-[1.5rem] px-10 font-black uppercase tracking-widest transition-all hover:-translate-y-1"
              @click="openAddModal"
            >
              <Plus class="h-6 w-6" />
              Добавить товар
            </button>
          </div>
        </div>
      </div>
    </section>

    <section ref="menuSection" class="pt-8">
      <div class="mb-12 flex flex-col items-end justify-between gap-6 md:flex-row">
        <div class="flex flex-col gap-2">
          <span class="pl-1 text-xs font-bold uppercase tracking-[0.2em] text-secondary/40">Наше меню</span>
          <h2 class="text-5xl font-black uppercase tracking-tighter text-secondary">Каталог</h2>
        </div>
        <div class="max-w-full overflow-x-auto rounded-[2rem] border border-white/40 bg-base-100/50 p-2 shadow-soft backdrop-blur-md">
          <div class="flex min-w-max gap-2">
            <a
              v-for="category in CATEGORIES"
              :key="category"
              class="btn btn-sm h-10 rounded-full border-none px-6 text-[10px] font-bold uppercase tracking-widest transition-all duration-300"
              :class="selectedCategory === category ? 'btn-primary text-white shadow-lg shadow-primary/30' : 'btn-ghost text-secondary/60 hover:bg-base-200'"
              @click="selectedCategory = category"
            >
              {{ category }}
            </a>
          </div>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-32">
        <span class="loading loading-spinner loading-lg text-primary"></span>
      </div>
      <div v-else-if="error" class="alert alert-error rounded-[2rem] shadow-xl">{{ error }}</div>

      <div v-else class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4">
        <div
          v-for="product in filteredProducts"
          :key="product.id"
          class="card group h-full overflow-visible rounded-[2.5rem] border border-transparent bg-base-100 shadow-soft transition-all duration-500 hover:-translate-y-2 hover:border-primary/10 hover:shadow-glow"
        >
          <figure class="mask mask-squircle relative m-2 mb-0 h-64 overflow-hidden rounded-t-[2.5rem]">
            <img :src="product.image_url" :alt="product.name" class="h-full w-full object-cover transition-transform duration-700 group-hover:scale-110" @error="handleImageError" />
            <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent opacity-60"></div>
            <div class="absolute right-4 top-4 z-20">
              <span class="badge rounded-2xl bg-white/90 p-4 text-lg font-black text-secondary shadow-lg backdrop-blur">{{ formatPrice(product.price) }}</span>
            </div>
            <div class="absolute bottom-4 left-4 z-20 flex flex-col items-start gap-2">
              <span class="badge badge-primary border-none px-3 py-1 text-[10px] font-black uppercase text-white shadow-md">
                {{ CATEGORIES[(product.category_id ?? 0) + 1] || 'Классика' }}
              </span>
              <span v-if="!product.is_available" class="badge badge-error px-3 py-1 text-[10px] font-black uppercase text-white shadow-md">SOLD OUT</span>
              <span v-if="isManagerMode && isDraftProduct(product.id)" class="badge border-none bg-slate-900 px-3 py-1 text-[10px] font-black uppercase text-white shadow-md">DRAFT</span>
            </div>
          </figure>

          <div class="card-body flex flex-col justify-between p-8 pt-6">
            <div>
              <h3 class="card-title mb-2 text-2xl font-black leading-tight tracking-tight transition-colors group-hover:text-primary">{{ product.name }}</h3>
              <p class="line-clamp-3 text-sm font-medium leading-relaxed text-secondary/60">{{ product.description }}</p>
            </div>

            <div class="card-actions mt-6 justify-end">
              <button
                v-if="!isManagerMode"
                class="btn btn-primary btn-block h-14 gap-3 rounded-2xl border-none font-black uppercase shadow-xl shadow-primary/20 transition-all group-hover:scale-[1.02] hover:shadow-primary/40"
                :disabled="!product.is_available"
                @click="handleAddToCart(product)"
              >
                <Plus class="h-5 w-5" />
                В корзину
              </button>
              <button
                v-else
                class="btn btn-outline btn-secondary btn-block h-14 gap-3 rounded-2xl font-black uppercase"
                @click="openEditModal(product)"
              >
                <FileText class="h-5 w-5" />
                Изменить
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <AppModal :show="showAuthModal" @close="showAuthModal = false">
      <div class="flex flex-col items-center p-10 text-center">
        <div class="mb-8 rounded-[2rem] bg-primary/10 p-6">
          <LogIn class="h-12 w-12 text-primary" />
        </div>
        <h3 class="mb-3 text-3xl font-black uppercase tracking-tighter">Нужна авторизация</h3>
        <p class="mb-10 px-4 font-medium leading-relaxed text-base-content/50">Вам нужно войти в аккаунт, чтобы делать заказы. Присоединяйтесь к нам!</p>
        <div class="flex w-full flex-col gap-3">
          <router-link to="/login" class="btn btn-primary btn-lg h-16 rounded-2xl font-black uppercase shadow-xl shadow-primary/20">Войти</router-link>
          <router-link to="/register" class="btn btn-outline btn-lg h-16 rounded-2xl font-bold border-2">Создать аккаунт</router-link>
        </div>
      </div>
    </AppModal>

    <AppModal :show="showAddModal" maxWidth="7xl" @close="showAddModal = false">
      <div class="px-6 pb-6 pt-6 md:px-8 md:pb-7 md:pt-7">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h3 class="text-4xl font-black uppercase tracking-[-0.08em] text-secondary">{{ isEditing ? 'Изменить товар' : 'Новый товар' }}</h3>
            <p class="mt-2 text-sm font-semibold text-secondary/45">
              Вы вошли как:
              <span class="font-black text-primary">Менеджер</span>
            </p>
          </div>

          <button class="manager-modal-close" type="button" @click="showAddModal = false">
            <X class="h-5 w-5" />
          </button>
        </div>

        <form class="mt-7 space-y-8" @submit.prevent="saveProduct('published')">
          <div class="grid gap-8 xl:grid-cols-2">
            <div class="space-y-5">
              <label class="manager-product-field">
                <span class="manager-product-label">Наименование</span>
                <span class="manager-product-hint">products.name</span>
                <input v-model="productForm.name" type="text" class="manager-product-input" placeholder="Маргарита" />
              </label>

              <label class="manager-product-field">
                <span class="manager-product-label">Категория</span>
                <span class="manager-product-hint">products.category_id</span>
                <select v-model="productForm.categoryId" class="manager-product-input">
                  <option v-for="(value, label) in CATEGORY_MAP" :key="label" :value="value">{{ label }}</option>
                </select>
              </label>

              <label class="manager-product-field">
                <span class="manager-product-label">Цена (₽)</span>
                <span class="manager-product-hint">products.price</span>
                <input v-model.number="productForm.price" type="number" class="manager-product-input" placeholder="1200" />
              </label>

              <label class="manager-product-field">
                <span class="manager-product-label">Единица измерения</span>
                <span class="manager-product-hint">products.unit_id</span>
                <select v-model="productForm.unit" class="manager-product-input">
                  <option v-for="unit in managerUnitOptions" :key="unit" :value="unit">{{ unit }}</option>
                </select>
              </label>

              <label class="manager-product-field">
                <span class="manager-product-label">Время приготовления, мин</span>
                <span class="manager-product-hint">products.cooking_time_min</span>
                <input v-model.number="productForm.cookingTime" type="number" class="manager-product-input" placeholder="15" />
              </label>

              <label class="manager-product-toggle">
                <span class="manager-product-label">Активен в меню</span>
                <span class="manager-product-hint">products.is_active</span>
                <span class="mt-3 flex items-center gap-3">
                  <input v-model="productForm.isAvailable" type="checkbox" class="toggle toggle-primary" />
                  <span class="text-base font-semibold text-secondary">{{ productForm.isAvailable ? 'Да' : 'Нет' }}</span>
                </span>
              </label>
            </div>

            <div class="space-y-5">
              <label class="manager-product-field">
                <span class="manager-product-label">Описание</span>
                <span class="manager-product-hint">products.description</span>
                <textarea v-model="productForm.description" class="manager-product-textarea" placeholder="Классическая итальянская пицца с тонким тестом..."></textarea>
              </label>

              <label class="manager-product-field">
                <span class="manager-product-label">URL изображения</span>
                <span class="manager-product-hint">products.image_url</span>
                <input v-model="productForm.imageUrl" type="text" class="manager-product-input" placeholder="https://pizzagood.ru/images/products/margarita.jpg" />
              </label>

              <label class="manager-product-field">
                <span class="manager-product-label">ID товара</span>
                <span class="manager-product-hint">products.product_id</span>
                <input :value="productForm.systemId" type="text" class="manager-product-input manager-product-input-readonly" readonly />
              </label>

              <label class="manager-product-field">
                <span class="manager-product-label">Дата создания</span>
                <span class="manager-product-hint">products.created_at</span>
                <input :value="productForm.createdAt" type="text" class="manager-product-input manager-product-input-readonly" readonly />
              </label>
            </div>
          </div>

          <div class="manager-product-actions">
            <button type="submit" class="manager-product-submit" :disabled="isSubmitting">
              <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
              <span v-else>Опубликовать</span>
            </button>
            <button type="button" class="manager-product-draft" :disabled="isSubmitting" @click="saveProduct('draft')">
              Сохранить как черновик
            </button>
            <button type="button" class="manager-product-cancel" @click="showAddModal = false">
              Отмена
            </button>
          </div>
        </form>
      </div>
    </AppModal>
  </div>
</template>

<style scoped>
.manager-modal-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 3rem;
  width: 3rem;
  border-radius: 9999px;
  background: rgba(15, 23, 42, 0.05);
  color: rgba(15, 23, 42, 0.72);
  transition: all 0.2s ease;
}

.manager-modal-close:hover {
  background: rgba(255, 91, 102, 0.1);
  color: rgb(255, 91, 102);
}

.manager-product-field,
.manager-product-toggle {
  display: flex;
  flex-direction: column;
}

.manager-product-label {
  color: rgba(15, 23, 42, 0.86);
  font-size: 0.98rem;
  font-weight: 900;
  letter-spacing: 0.02em;
  text-transform: uppercase;
}

.manager-product-hint {
  margin-top: 0.18rem;
  color: rgba(15, 23, 42, 0.32);
  font-size: 0.82rem;
  font-weight: 700;
}

.manager-product-input,
.manager-product-textarea {
  margin-top: 0.65rem;
  width: 100%;
  border-radius: 1rem;
  border: 1px solid rgba(203, 213, 225, 0.95);
  background: white;
  color: rgb(30, 41, 59);
  font-size: 1rem;
  font-weight: 600;
  outline: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.manager-product-input {
  min-height: 3.4rem;
  padding: 0.9rem 1rem;
}

.manager-product-textarea {
  min-height: 10rem;
  resize: vertical;
  padding: 1rem;
  line-height: 1.65;
}

.manager-product-input:focus,
.manager-product-textarea:focus {
  border-color: rgba(255, 91, 102, 0.55);
  box-shadow: 0 0 0 4px rgba(255, 91, 102, 0.08);
}

.manager-product-input-readonly {
  background: rgba(248, 250, 252, 0.9);
  color: rgba(15, 23, 42, 0.48);
}

.manager-product-actions {
  display: grid;
  gap: 1rem;
  border-top: 1px solid rgba(226, 232, 240, 0.95);
  padding-top: 1.5rem;
}

.manager-product-submit,
.manager-product-draft,
.manager-product-cancel {
  display: inline-flex;
  min-height: 3.6rem;
  align-items: center;
  justify-content: center;
  border-radius: 1rem;
  padding: 0.95rem 1.4rem;
  font-size: 0.92rem;
  font-weight: 900;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  transition: all 0.2s ease;
}

.manager-product-submit {
  border: 0;
  background: linear-gradient(135deg, #ff5b66 0%, #ff3f55 100%);
  color: white;
  box-shadow: 0 24px 50px -28px rgba(255, 63, 85, 0.75);
}

.manager-product-submit:hover,
.manager-product-draft:hover {
  transform: translateY(-1px);
}

.manager-product-draft {
  border: 1px solid rgba(148, 163, 184, 0.45);
  background: white;
  color: rgba(15, 23, 42, 0.82);
}

.manager-product-cancel {
  border: 0;
  background: transparent;
  color: rgba(15, 23, 42, 0.48);
}

.manager-product-submit:disabled,
.manager-product-draft:disabled {
  transform: none;
  opacity: 0.7;
}

@media (min-width: 960px) {
  .manager-product-actions {
    grid-template-columns: minmax(0, 1.05fr) minmax(0, 1.05fr) minmax(180px, 0.8fr);
    align-items: center;
  }
}
</style>

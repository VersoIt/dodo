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
const selectedCategory = ref('All')
const menuSection = ref<HTMLElement | null>(null)
const showAddModal = ref(false)
const showAuthModal = ref(false)
const isSubmitting = ref(false)
const isEditing = ref(false)
const editingId = ref<string | null>(null)

const categories = ['All', 'Classic', 'Premium', 'Veggie', 'Spicy', 'Drinks', 'Desserts']
const categoryMap: Record<string, number> = {
  'Classic': 0,
  'Premium': 1,
  'Veggie': 2,
  'Spicy': 3,
  'Drinks': 4,
  'Desserts': 5
}

// Form data for new product
const productForm = ref({
  name: '',
  description: '',
  price: 0,
  imageUrl: '',
  categoryId: 0,
  isAvailable: true
})

const filteredProducts = computed(() => {
  if (selectedCategory.value === 'All') return products.value
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
      category: categories[p.category_id + 1] || 'Classic',
      isAvailable: p.is_available ?? true
    }))
  } catch (err: any) {
    console.error('Error fetching products:', err)
    error.value = 'Failed to load products.'
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
  addToast(`${product.name} added to cart!`, 'success')
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
      addToast(isEditing.value ? 'Product updated!' : 'Product added!', 'success')
      showAddModal.value = false
      await fetchProducts()
    }
  } catch (err: any) {
    addToast(err.response?.data?.error || 'Operation failed', 'error')
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
            {{ authStore.user?.role === 'manager' ? 'Menu Management' : 'Hot & Fresh' }}
          </h1>
          <p class="mb-5 text-lg text-white/80">
            {{ authStore.user?.role === 'manager' 
               ? 'Update your catalog, add new seasonal pizzas and manage availability.' 
               : 'Experience the best pizza in town, delivered straight to your door.' }}
          </p>
          <button v-if="authStore.user?.role !== 'manager'" @click="scrollToMenu" class="btn btn-primary btn-lg">Order Now</button>
          <button v-else @click="openAddModal" class="btn btn-secondary btn-lg gap-2">
            <Plus class="w-6 h-6" /> Add New Product
          </button>
        </div>
      </div>
    </section>

    <!-- Menu Section -->
    <section ref="menuSection">
      <div class="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
        <div class="flex items-center gap-3">
          <h2 class="text-3xl font-bold">Catalog</h2>
          <div v-if="authStore.user?.role === 'manager'" class="badge badge-secondary badge-outline">Admin Mode</div>
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
        <p class="text-base-content/60">No products found in this category.</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div v-for="product in filteredProducts" :key="product.id" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 group">
          <figure class="relative h-48 overflow-hidden">
            <img :src="product.imageUrl" :alt="product.name" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" />
            <div class="absolute top-2 right-2">
              <span class="badge badge-secondary font-semibold">${{ product.price?.toFixed(2) }}</span>
            </div>
            <div class="absolute top-2 left-2 flex flex-col gap-1">
              <span class="badge badge-ghost bg-black/40 text-white text-[10px] uppercase border-none">{{ product.category }}</span>
              <span v-if="!product.isAvailable" class="badge badge-error text-[10px] uppercase font-bold">Out of stock</span>
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
                Add to Cart
              </button>
              <button 
                v-else
                @click="openEditModal(product)"
                class="btn btn-outline btn-sm btn-secondary gap-2"
              >
                <FileText class="w-4 h-4" /> Edit
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- AUTH MODAL (For guests) -->
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="showAuthModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div 
          class="fixed inset-0 bg-black/60 backdrop-blur-sm"
          @click="showAuthModal = false"
        ></div>
        
        <!-- Modal Content -->
        <Transition
          appear
          enter-active-class="transition duration-500 delay-100 ease-out"
          enter-from-class="opacity-0 translate-y-8 scale-95"
          enter-to-class="opacity-100 translate-y-0 scale-100"
        >
          <div class="relative bg-base-100 w-full max-w-md rounded-3xl border-t-4 border-primary shadow-2xl overflow-hidden">
            <div class="p-8 flex flex-col items-center text-center">
              <div class="bg-primary/10 p-4 rounded-full mb-6">
                <LogIn class="w-12 h-12 text-primary" />
              </div>
              <h3 class="font-bold text-3xl mb-2 tracking-tight">Login Required</h3>
              <p class="text-base-content/60 mb-8 px-2">
                You need to have an account to place orders. Join our pizza-loving community today!
              </p>
              
              <div class="flex flex-col w-full gap-3">
                <router-link to="/login" class="btn btn-primary btn-lg btn-block gap-2 shadow-lg shadow-primary/20">
                  <LogIn class="w-5 h-5" /> Sign In
                </router-link>
                <router-link to="/register" class="btn btn-outline btn-lg btn-block border-2">
                  Create New Account
                </router-link>
                <button @click="showAuthModal = false" class="btn btn-ghost btn-sm mt-4 opacity-50 hover:opacity-100">
                  Maybe later
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>

    <!-- FORM MODAL (Add/Edit) -->
    <dialog :class="{ 'modal-open': showAddModal }" class="modal modal-bottom sm:modal-middle transition-all">
      <div class="modal-box bg-base-100 max-w-2xl border border-base-300 shadow-2xl">
        <div class="flex justify-between items-center mb-6">
          <div class="flex items-center gap-3 text-secondary">
            <PackagePlus class="w-8 h-8" />
            <h3 class="font-bold text-2xl">{{ isEditing ? 'Edit Product' : 'Create New Item' }}</h3>
          </div>
          <button @click="showAddModal = false" class="btn btn-ghost btn-sm btn-circle"><X /></button>
        </div>
        
        <form @submit.prevent="handleSubmit" class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Left Col -->
          <div class="space-y-4">
            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Product Name</span></label>
              <div class="relative">
                <Tag class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                <input v-model="productForm.name" type="text" placeholder="e.g. Pepperoni Extreme" class="input input-bordered w-full pl-10" required />
              </div>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Price ($)</span></label>
              <div class="relative">
                <DollarSign class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                <input v-model="productForm.price" type="number" step="0.01" placeholder="15.99" class="input input-bordered w-full pl-10" required />
              </div>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Category</span></label>
              <select v-model="productForm.categoryId" class="select select-bordered w-full">
                <option v-for="(val, key) in categoryMap" :key="key" :value="val">{{ key }}</option>
              </select>
            </div>

            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-4">
                <input v-model="productForm.isAvailable" type="checkbox" class="checkbox checkbox-secondary" />
                <span class="label-text font-bold">Available for order</span>
              </label>
            </div>
          </div>

          <!-- Right Col -->
          <div class="space-y-4">
            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Image URL</span></label>
              <div class="relative">
                <ImageIcon class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                <input v-model="productForm.imageUrl" type="text" placeholder="https://images.unsplash.com/..." class="input input-bordered w-full pl-10" />
              </div>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Description</span></label>
              <div class="relative">
                <FileText class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                <textarea v-model="productForm.description" class="textarea textarea-bordered w-full pl-10 h-32" placeholder="Tell us about this delicious pizza..."></textarea>
              </div>
            </div>
          </div>

          <div class="col-span-full mt-4 flex gap-3">
            <button type="submit" class="btn btn-secondary flex-1" :disabled="isSubmitting">
              <span v-if="isSubmitting" class="loading loading-spinner"></span>
              {{ isEditing ? 'Save Changes' : 'Publish to Menu' }}
            </button>
            <button type="button" @click="showAddModal = false" class="btn btn-ghost">Cancel</button>
          </div>
        </form>
      </div>
      <div class="modal-backdrop bg-black/60 backdrop-blur-sm" @click="showAddModal = false"></div>
    </dialog>
  </div>
</template>

<style scoped>
.tab-active {
  @apply font-bold !important;
}
</style>

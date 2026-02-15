<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Plus, PackagePlus, Image as ImageIcon, Tag, DollarSign, FileText, X } from 'lucide-vue-next'
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
}

const products = ref<Product[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const selectedCategory = ref('All')
const menuSection = ref<HTMLElement | null>(null)
const showAddModal = ref(false)
const isSubmitting = ref(false)

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
const newProduct = ref({
  name: '',
  description: '',
  price: 0,
  imageUrl: '',
  categoryId: 0
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
    const responseData = response.data.data
    const productsData = responseData.products || []
    
    products.value = productsData.map((p: any) => ({
      id: p.id,
      name: p.name,
      description: p.description,
      price: p.price,
      imageUrl: p.imageUrl || `https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=500&auto=format&fit=crop&sig=${p.id}`,
      categoryId: p.categoryId,
      category: categories[p.categoryId + 1] || 'Classic'
    }))
  } catch (err: any) {
    console.error('Error fetching products:', err)
    error.value = 'Failed to load products. Please try again later.'
  } finally {
    loading.value = false
  }
}

const handleAddProduct = async () => {
  try {
    isSubmitting.value = true
    const response = await axios.post('/api/v1/catalog/products', {
      name: newProduct.value.name,
      description: newProduct.value.description,
      price: parseFloat(newProduct.value.price.toString()),
      category_id: newProduct.value.categoryId,
      image_url: newProduct.value.imageUrl
    })

    if (response.data.success) {
      addToast('Product added successfully!', 'success')
      showAddModal.value = false
      // Reset form
      newProduct.value = { name: '', description: '', price: 0, imageUrl: '', categoryId: 0 }
      await fetchProducts()
    }
  } catch (err: any) {
    addToast(err.response?.data?.error || 'Failed to add product', 'error')
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  fetchProducts()
})
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
          <button v-else @click="showAddModal = true" class="btn btn-secondary btn-lg gap-2">
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
            <div class="absolute top-2 left-2">
              <span class="badge badge-ghost bg-black/40 text-white text-[10px] uppercase border-none">{{ product.category }}</span>
            </div>
          </figure>
          <div class="card-body p-6">
            <h3 class="card-title text-xl font-bold">{{ product.name }}</h3>
            <p class="text-sm text-base-content/70 line-clamp-2 mb-4">{{ product.description }}</p>
            <div class="card-actions justify-end mt-auto">
              <!-- ONLY CLIENTS CAN SEE ADD TO CART -->
              <button 
                v-if="authStore.user?.role !== 'manager'"
                @click="cartStore.addToCart(product); addToast(`${product.name} added to cart!`)" 
                class="btn btn-primary btn-sm gap-2"
              >
                <Plus class="w-4 h-4" />
                Add to Cart
              </button>
              <!-- MANAGERS SEE EDIT BUTTON (UI PLACEHOLDER) -->
              <button 
                v-else
                class="btn btn-outline btn-sm btn-secondary"
              >
                Edit Product
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ADD PRODUCT MODAL -->
    <dialog :class="{ 'modal-open': showAddModal }" class="modal modal-bottom sm:modal-middle transition-all">
      <div class="modal-box bg-base-100 max-w-2xl border border-base-300 shadow-2xl">
        <div class="flex justify-between items-center mb-6">
          <div class="flex items-center gap-3 text-secondary">
            <PackagePlus class="w-8 h-8" />
            <h3 class="font-bold text-2xl">Create New Item</h3>
          </div>
          <button @click="showAddModal = false" class="btn btn-ghost btn-sm btn-circle"><X /></button>
        </div>
        
        <form @submit.prevent="handleAddProduct" class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Left Col -->
          <div class="space-y-4">
            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Product Name</span></label>
              <div class="relative">
                <Tag class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                <input v-model="newProduct.name" type="text" placeholder="e.g. Pepperoni Extreme" class="input input-bordered w-full pl-10" required />
              </div>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Price ($)</span></label>
              <div class="relative">
                <DollarSign class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                <input v-model="newProduct.price" type="number" step="0.01" placeholder="15.99" class="input input-bordered w-full pl-10" required />
              </div>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Category</span></label>
              <select v-model="newProduct.categoryId" class="select select-bordered w-full">
                <option v-for="(val, key) in categoryMap" :key="key" :value="val">{{ key }}</option>
              </select>
            </div>
          </div>

          <!-- Right Col -->
          <div class="space-y-4">
            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Image URL</span></label>
              <div class="relative">
                <ImageIcon class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                <input v-model="newProduct.imageUrl" type="text" placeholder="https://images.unsplash.com/..." class="input input-bordered w-full pl-10" />
              </div>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text font-bold">Description</span></label>
              <div class="relative">
                <FileText class="absolute left-3 top-3 w-5 h-5 opacity-40" />
                <textarea v-model="newProduct.description" class="textarea textarea-bordered w-full pl-10 h-32" placeholder="Tell us about this delicious pizza..."></textarea>
              </div>
            </div>
          </div>

          <div class="col-span-full mt-4 flex gap-3">
            <button type="submit" class="btn btn-secondary flex-1" :disabled="isSubmitting">
              <span v-if="isSubmitting" class="loading loading-spinner"></span>
              Publish to Menu
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

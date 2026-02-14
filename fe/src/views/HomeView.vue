<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { Plus } from 'lucide-vue-next'
import { useCartStore } from '../store/cart'

const cartStore = useCartStore()

interface Product {
  id: string
  name: string
  description: string
  price: number
  imageUrl: string
  category: string
}

const products = ref<Product[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const selectedCategory = ref('All')
const menuSection = ref<HTMLElement | null>(null)

const categories = ['All', 'Classic', 'Premium', 'Veggie']

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
    products.value = response.data.products || []
  } catch (err: any) {
    console.error('Error fetching products:', err)
    error.value = 'Failed to load products. Please try again later.'
    // Fallback data for demo if API fails
    products.value = [
      { id: '1', name: 'Margherita', description: 'Classic tomato sauce, mozzarella, and fresh basil', price: 12.99, imageUrl: 'https://images.unsplash.com/photo-1574071318508-1cdbad80ad50?q=80&w=500&auto=format&fit=crop', category: 'Classic' },
      { id: '2', name: 'Pepperoni', description: 'Tomato sauce, mozzarella, and spicy pepperoni slices', price: 14.99, imageUrl: 'https://images.unsplash.com/photo-1628840042765-356cda07504e?q=80&w=500&auto=format&fit=crop', category: 'Premium' },
      { id: '3', name: 'Hawaiian', description: 'Ham and pineapple on a bed of mozzarella and tomato sauce', price: 13.99, imageUrl: 'https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?q=80&w=500&auto=format&fit=crop', category: 'Classic' },
      { id: '4', name: 'BBQ Chicken', description: 'Grilled chicken, red onions, and tangy BBQ sauce', price: 15.99, imageUrl: 'https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=500&auto=format&fit=crop', category: 'Premium' },
    ]
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProducts()
})
</script>

<template>
  <div class="space-y-8">
    <!-- Hero Section -->
    <section class="hero min-h-[40vh] rounded-3xl overflow-hidden shadow-2xl" style="background-image: url(https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=1200&auto=format&fit=crop);">
      <div class="hero-overlay bg-black bg-opacity-60"></div>
      <div class="hero-content text-center text-neutral-content">
        <div class="max-w-md">
          <h1 class="mb-5 text-5xl font-extrabold uppercase tracking-tighter">Hot & Fresh</h1>
          <p class="mb-5 text-lg">Experience the best pizza in town, delivered straight to your door. Hand-crafted with love and the finest ingredients.</p>
          <button @click="scrollToMenu" class="btn btn-primary btn-lg">Order Now</button>
        </div>
      </div>
    </section>

    <!-- Menu Section -->
    <section ref="menuSection">
      <div class="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
        <h2 class="text-3xl font-bold">Our Menu</h2>
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

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div v-for="product in filteredProducts" :key="product.id" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 group">
          <figure class="relative h-48 overflow-hidden">
            <img :src="product.imageUrl" :alt="product.name" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" />
            <div class="absolute top-2 right-2">
              <span class="badge badge-secondary font-semibold">${{ product.price }}</span>
            </div>
          </figure>
          <div class="card-body p-6">
            <h3 class="card-title text-xl font-bold">{{ product.name }}</h3>
            <p class="text-sm text-base-content/70 line-clamp-2 mb-4">{{ product.description }}</p>
            <div class="card-actions justify-end mt-auto">
              <button 
                @click="cartStore.addToCart(product)" 
                class="btn btn-primary btn-sm gap-2"
              >
                <Plus class="w-4 h-4" />
                Add to Cart
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

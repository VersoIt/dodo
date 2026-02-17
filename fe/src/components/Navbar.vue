<script setup lang="ts">
import { useAuthStore } from '../store/auth'
import { useCartStore } from '../store/cart'
import { useRouter } from 'vue-router'
import { ShoppingCart, User, LogOut, Pizza, Menu, Info, ChefHat, Truck, BarChart3 } from 'lucide-vue-next'

const authStore = useAuthStore()
const cartStore = useCartStore()
const router = useRouter()

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}
</script>

<template>
  <div class="navbar bg-base-100/80 backdrop-blur-xl shadow-soft px-4 md:px-8 rounded-[2rem] mt-6 mx-auto max-w-7xl sticky top-6 z-50 border border-white/40 supports-[backdrop-filter]:bg-base-100/60 transition-all duration-500 hover:shadow-glow/20">
    <div class="flex-1">
      <router-link to="/" class="btn btn-ghost gap-3 px-2 hover:bg-transparent group">
        <div class="bg-primary p-2.5 rounded-2xl shadow-lg shadow-primary/30 group-hover:rotate-12 group-hover:scale-110 transition-all duration-300">
          <Pizza class="w-6 h-6 text-primary-content" />
        </div>
        <div class="flex flex-col items-start leading-none">
          <span class="text-xl font-black tracking-tight uppercase text-secondary group-hover:text-primary transition-colors duration-300">Pizza</span>
          <span class="text-[10px] font-bold tracking-[0.2em] text-secondary/40 uppercase pl-0.5">Good</span>
        </div>
      </router-link>
    </div>

    <div class="flex-none gap-2 md:gap-4">
      <div class="hidden md:flex items-center gap-2 mr-6">
        <router-link to="/" active-class="text-primary bg-primary/10" class="btn btn-ghost btn-sm h-10 px-4 rounded-xl gap-2 font-black uppercase text-[10px] tracking-widest hover:bg-base-200/50 transition-all">
          <Menu class="w-4 h-4" /> Меню
        </router-link>
        <template v-if="authStore.isAuthenticated">
          <router-link v-if="['chef', 'manager'].includes(authStore.user?.role)" to="/kitchen" active-class="text-accent bg-accent/10" class="btn btn-ghost btn-sm h-10 px-4 rounded-xl gap-2 font-black uppercase text-[10px] tracking-widest hover:bg-base-200/50 transition-all">
            <ChefHat class="w-4 h-4" /> Кухня
          </router-link>
          <router-link v-if="['courier', 'manager'].includes(authStore.user?.role)" to="/logistics" active-class="text-info bg-info/10" class="btn btn-ghost btn-sm h-10 px-4 rounded-xl gap-2 font-black uppercase text-[10px] tracking-widest hover:bg-base-200/50 transition-all">
            <Truck class="w-4 h-4" /> Доставка
          </router-link>
          <router-link v-if="authStore.user?.role === 'manager'" to="/manager" active-class="text-secondary bg-secondary/10" class="btn btn-ghost btn-sm h-10 px-4 rounded-xl gap-2 font-black uppercase text-[10px] tracking-widest hover:bg-base-200/50 transition-all">
            <BarChart3 class="w-4 h-4" /> Управление
          </router-link>
        </template>
      </div>

      <router-link to="/cart" class="btn btn-ghost btn-circle w-12 h-12 relative bg-base-200/50 hover:bg-primary/10 hover:text-primary group transition-all duration-300">
        <ShoppingCart class="w-5 h-5 transition-transform group-hover:scale-110" />
        <span v-if="cartStore.totalItems > 0" class="badge badge-primary badge-sm absolute top-0 right-0 font-bold animate-in zoom-in border-2 border-base-100 shadow-md">{{ cartStore.totalItems }}</span>
      </router-link>

      <div v-if="authStore.isAuthenticated" class="dropdown dropdown-end ml-2">
        <div tabindex="0" role="button" class="btn btn-ghost gap-3 pl-2 pr-1 h-12 rounded-2xl hover:bg-base-200/50 transition-all border border-transparent hover:border-base-200">
          <div class="flex flex-col items-end hidden md:flex"><span class="text-[9px] font-black uppercase text-secondary/40 tracking-wider mb-0.5">{{ authStore.user?.role }}</span><span class="text-sm font-bold leading-none text-secondary">{{ authStore.user?.name }}</span></div>
          <div class="avatar placeholder"><div class="bg-secondary text-secondary-content rounded-xl w-10 h-10 shadow-lg ring-2 ring-base-100 ring-offset-2 ring-offset-base-100"><span class="text-lg font-black">{{ authStore.user?.name?.charAt(0).toUpperCase() }}</span></div></div>
        </div>
        <ul tabindex="0" class="mt-4 z-[1] p-2 shadow-soft menu menu-sm dropdown-content bg-base-100/95 backdrop-blur-xl rounded-[1.5rem] w-60 border border-base-200">
          <li class="menu-title opacity-40 uppercase text-[9px] font-black tracking-[0.2em] px-4 py-3">Аккаунт</li>
          <li><router-link to="/profile" active-class="bg-base-200 text-primary" class="py-3 px-4 rounded-xl gap-3 font-bold hover:bg-base-200"><User class="w-4 h-4 opacity-50" /> Профиль</router-link></li>
          <div class="divider my-1 opacity-10 px-4"></div>
          <li><button @click="handleLogout" class="text-error py-3 px-4 rounded-xl gap-3 font-bold hover:bg-error/10 hover:text-error"><LogOut class="w-4 h-4" /> Выйти</button></li>
        </ul>
      </div>
      <router-link v-else to="/login" class="btn btn-primary h-12 rounded-2xl px-8 font-black uppercase tracking-widest shadow-xl shadow-primary/30 hover:shadow-primary/50 hover:-translate-y-0.5 transition-all ml-2">Войти</router-link>
    </div>
  </div>
</template>

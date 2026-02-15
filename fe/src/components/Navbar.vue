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
  <div class="navbar bg-base-100 shadow-lg px-4 md:px-8 rounded-2xl mt-4 sticky top-4 z-50 border border-base-200">
    <div class="flex-1">
      <router-link to="/" class="btn btn-ghost gap-2 px-2 hover:bg-transparent group">
        <div class="bg-primary p-2 rounded-xl group-hover:rotate-12 transition-transform duration-300">
          <Pizza class="w-6 h-6 text-primary-content" />
        </div>
        <span class="text-xl font-black tracking-tighter uppercase">Пиццерия</span>
      </router-link>
    </div>

    <div class="flex-none gap-2 md:gap-4">
      <div class="hidden md:flex items-center gap-1 mr-4">
        <template v-if="authStore.isAuthenticated">
          <router-link v-if="authStore.user?.role === 'chef' || authStore.user?.role === 'manager'" to="/kitchen" class="btn btn-ghost btn-sm gap-2 font-bold uppercase text-[10px] tracking-widest">
            <ChefHat class="w-4 h-4" /> Кухня
          </router-link>
          <router-link v-if="authStore.user?.role === 'courier' || authStore.user?.role === 'manager'" to="/logistics" class="btn btn-ghost btn-sm gap-2 font-bold uppercase text-[10px] tracking-widest">
            <Truck class="w-4 h-4" /> Доставка
          </router-link>
          <router-link v-if="authStore.user?.role === 'manager'" to="/manager" class="btn btn-ghost btn-sm gap-2 font-bold uppercase text-[10px] tracking-widest text-secondary">
            <BarChart3 class="w-4 h-4" /> Управление
          </router-link>
        </template>
        <template v-if="!authStore.isAuthenticated || authStore.user?.role === 'client' || authStore.user?.role === 'manager'">
          <router-link to="/" class="btn btn-ghost btn-sm gap-2 font-bold uppercase text-[10px] tracking-widest">
            <Menu class="w-4 h-4" /> Меню
          </router-link>
          <router-link to="/about" class="btn btn-ghost btn-sm gap-2 font-bold uppercase text-[10px] tracking-widest">
            <Info class="w-4 h-4" /> О нас
          </router-link>
        </template>
      </div>

      <router-link v-if="!authStore.isAuthenticated || authStore.user?.role === 'client'" to="/cart" class="btn btn-ghost btn-circle relative bg-base-200/50 hover:bg-primary/10 group transition-colors">
        <ShoppingCart class="w-5 h-5 group-hover:text-primary transition-colors" />
        <span v-if="cartStore.totalItems > 0" class="badge badge-primary badge-sm absolute -top-1 -right-1 font-bold animate-in zoom-in h-5 w-5 p-0">{{ cartStore.totalItems }}</span>
      </router-link>

      <div v-if="authStore.isAuthenticated" class="dropdown dropdown-end">
        <div tabindex="0" role="button" class="btn btn-ghost gap-3 px-2 md:px-4 rounded-xl hover:bg-base-200 transition-all">
          <div class="flex flex-col items-end hidden md:flex"><span class="text-[10px] font-black uppercase opacity-40 leading-none mb-1">{{ authStore.user?.role }}</span><span class="text-sm font-bold leading-none">{{ authStore.user?.name }}</span></div>
          <div class="avatar placeholder"><div class="bg-primary text-primary-content rounded-xl w-10 h-10 shadow-lg shadow-primary/20"><span class="text-xl font-black">{{ authStore.user?.name?.charAt(0).toUpperCase() }}</span></div></div>
        </div>
        <ul tabindex="0" class="mt-3 z-[1] p-2 shadow-2xl menu menu-sm dropdown-content bg-base-100 rounded-2xl w-52 border border-base-200">
          <li class="menu-title opacity-40 uppercase text-[9px] font-black tracking-[0.2em] px-4 py-2">Аккаунт</li>
          <li><router-link to="/profile" class="py-3 rounded-xl gap-3"><User class="w-4 h-4 opacity-50" /> Профиль</router-link></li>
          <div class="divider my-1 opacity-50 px-2"></div>
          <li><button @click="handleLogout" class="text-error py-3 rounded-xl gap-3 hover:bg-error/10"><LogOut class="w-4 h-4" /> Выйти</button></li>
        </ul>
      </div>
      <router-link v-else to="/login" class="btn btn-primary rounded-xl px-6 font-black uppercase shadow-lg shadow-primary/20">Войти</router-link>
    </div>
  </div>
</template>

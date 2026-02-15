<script setup lang="ts">
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { LogIn, Mail, Lock, AlertCircle } from 'lucide-vue-next'

const authStore = useAuthStore()
const router = useRouter()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const addToast = inject('addToast') as (msg: string, type?: any) => void

const handleLogin = async () => {
  if (!email.value || !password.value) {
    error.value = 'Пожалуйста, заполните все поля'
    return
  }

  loading.value = true
  error.value = ''
  
  try {
    const success = await authStore.login(email.value, password.value)
    if (success) {
      addToast('С возвращением!', 'success')
      // Redirect based on role
      const role = authStore.user?.role
      if (role === 'chef') router.push('/kitchen')
      else if (role === 'courier') router.push('/logistics')
      else router.push('/')
    } else {
      error.value = 'Неверный email или пароль'
    }
  } catch (err) {
    error.value = 'Произошла ошибка при входе'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center px-4">
    <div class="max-w-md w-full">
      <div class="text-center mb-10">
        <div class="bg-primary/10 w-20 h-20 rounded-3xl flex items-center justify-center mx-auto mb-6">
          <LogIn class="w-10 h-10 text-primary" />
        </div>
        <h1 class="text-4xl font-black tracking-tight uppercase italic">Вход</h1>
        <p class="text-base-content/50 font-bold uppercase text-[10px] tracking-[0.2em] mt-2">Добро пожаловать в Пиццерию</p>
      </div>

      <div class="card bg-base-100 shadow-2xl border border-base-200 overflow-hidden rounded-[2.5rem]">
        <div class="card-body p-10">
          <form @submit.prevent="handleLogin" class="space-y-6">
            <div v-if="error" class="alert alert-error rounded-2xl py-3 text-sm flex gap-2">
              <AlertCircle class="w-4 h-4" />
              <span>{{ error }}</span>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40 tracking-widest">Email адрес</span></label>
              <div class="relative">
                <Mail class="absolute left-4 top-3.5 w-5 h-5 opacity-30" />
                <input v-model="email" type="email" placeholder="name@example.com" class="input input-bordered w-full pl-12 rounded-2xl h-12" required />
              </div>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40 tracking-widest">Пароль</span></label>
              <div class="relative">
                <Lock class="absolute left-4 top-3.5 w-5 h-5 opacity-30" />
                <input v-model="password" type="password" placeholder="••••••••" class="input input-bordered w-full pl-12 rounded-2xl h-12" required />
              </div>
            </div>

            <button type="submit" class="btn btn-primary btn-block h-14 rounded-2xl font-black uppercase shadow-lg shadow-primary/20 mt-4" :disabled="loading">
              <span v-if="loading" class="loading loading-spinner"></span>
              <span v-else>Войти</span>
            </button>
          </form>

          <div class="divider opacity-30 my-8">ИЛИ</div>

          <div class="text-center">
            <p class="text-sm text-base-content/60">Нет аккаунта?</p>
            <router-link to="/register" class="link link-primary font-black uppercase text-[10px] tracking-widest mt-2 block">Создать новый аккаунт</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

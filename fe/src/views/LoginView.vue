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
  if (!email.value || !password.value) { error.value = 'Пожалуйста, заполните все поля'; return }
  loading.value = true; error.value = ''
  try {
    const success = await authStore.login(email.value, password.value)
    if (success) {
      addToast('С возвращением!', 'success')
      const role = authStore.user?.role
      if (role === 'chef') router.push('/kitchen')
      else if (role === 'courier') router.push('/logistics')
      else router.push('/')
    } else { error.value = 'Неверный email или пароль' }
  } catch (err) { error.value = 'Произошла ошибка при входе' } finally { loading.value = false }
}
</script>

<template>
  <div class="min-h-[85vh] flex items-center justify-center relative overflow-hidden">
    <!-- Decorative blobs -->
    <div class="absolute -top-40 -left-40 w-96 h-96 bg-primary/20 rounded-full blur-[100px] animate-pulse"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 bg-accent/20 rounded-full blur-[100px] animate-pulse delay-1000"></div>

    <div class="w-full max-w-md relative z-10">
      <div class="text-center mb-12">
        <div class="bg-gradient-to-br from-primary to-error w-24 h-24 rounded-[2rem] flex items-center justify-center mx-auto mb-8 shadow-glow shadow-primary/40 rotate-3 hover:rotate-6 transition-transform duration-500">
          <LogIn class="w-10 h-10 text-white" />
        </div>
        <h1 class="text-5xl font-black tracking-tight text-secondary mb-2">Вход</h1>
        <p class="text-secondary/50 font-bold uppercase text-[11px] tracking-[0.2em]">Добро пожаловать обратно</p>
      </div>

      <div class="card bg-base-100/80 backdrop-blur-xl shadow-soft border border-white/50 rounded-[2.5rem] overflow-hidden hover:shadow-glow/10 transition-shadow duration-500">
        <div class="card-body p-10">
          <form @submit.prevent="handleLogin" class="space-y-6">
            <div v-if="error" class="alert alert-error rounded-2xl py-3 text-sm flex gap-3 shadow-lg text-white font-bold"><AlertCircle class="w-5 h-5" /><span>{{ error }}</span></div>
            
            <div class="form-control group">
              <label class="label pl-4"><span class="label-text font-black uppercase text-[10px] text-secondary/40 tracking-widest group-focus-within:text-primary transition-colors">Email</span></label>
              <div class="relative transition-all duration-300 transform group-focus-within:scale-[1.02]">
                <Mail class="absolute left-5 top-4 w-5 h-5 text-secondary/30 group-focus-within:text-primary transition-colors" />
                <input v-model="email" type="email" placeholder="hello@example.com" class="input input-lg w-full pl-14 rounded-2xl bg-base-200/50 border-transparent focus:bg-white focus:border-primary/20 focus:shadow-lg focus:shadow-primary/10 transition-all font-bold text-secondary placeholder:text-secondary/20 h-14" required />
              </div>
            </div>

            <div class="form-control group">
              <label class="label pl-4"><span class="label-text font-black uppercase text-[10px] text-secondary/40 tracking-widest group-focus-within:text-primary transition-colors">Пароль</span></label>
              <div class="relative transition-all duration-300 transform group-focus-within:scale-[1.02]">
                <Lock class="absolute left-5 top-4 w-5 h-5 text-secondary/30 group-focus-within:text-primary transition-colors" />
                <input v-model="password" type="password" placeholder="••••••••" class="input input-lg w-full pl-14 rounded-2xl bg-base-200/50 border-transparent focus:bg-white focus:border-primary/20 focus:shadow-lg focus:shadow-primary/10 transition-all font-bold text-secondary placeholder:text-secondary/20 h-14" required />
              </div>
            </div>

            <button type="submit" class="btn btn-primary btn-block h-16 rounded-2xl font-black uppercase tracking-widest text-sm shadow-xl shadow-primary/30 mt-6 hover:shadow-primary/50 hover:-translate-y-1 transition-all border-none" :disabled="loading">
              <span v-if="loading" class="loading loading-spinner"></span>
              <span v-else>Войти в аккаунт</span>
            </button>
          </form>
          
          <div class="divider opacity-10 my-8 font-black text-xs text-secondary/40 tracking-widest">ИЛИ</div>
          
          <div class="text-center">
            <p class="text-sm font-medium text-secondary/60 mb-4">Впервые у нас?</p>
            <router-link to="/register" class="btn btn-ghost btn-outline border-2 border-base-300 hover:border-primary hover:bg-primary/5 hover:text-primary w-full h-14 rounded-2xl font-black uppercase text-[10px] tracking-widest transition-all">Создать аккаунт</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

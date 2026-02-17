<script setup lang="ts">
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { UserPlus, Mail, Lock, User, ShieldCheck, AlertCircle } from 'lucide-vue-next'
import axios from 'axios'

const authStore = useAuthStore()
const router = useRouter()
const name = ref('')
const email = ref('')
const password = ref('')
const role = ref('client')
const loading = ref(false)
const error = ref('')
const addToast = inject('addToast') as (msg: string, type?: any) => void

const handleRegister = async () => {
  if (!name.value || !email.value || !password.value) { error.value = 'Пожалуйста, заполните все поля'; return }
  loading.value = true; error.value = ''
  try {
    const response = await axios.post('/api/v1/auth/register', { name: name.value, email: email.value, password: password.value, role: role.value })
    if (response.data.success) {
      addToast('Регистрация успешна!', 'success')
      const loginSuccess = await authStore.login(email.value, password.value)
      if (loginSuccess) {
        if (role.value === 'chef') router.push('/kitchen')
        else if (role.value === 'courier') router.push('/logistics')
        else router.push('/')
      } else { router.push('/login') }
    }
  } catch (err: any) { error.value = err.response?.data?.error || 'Ошибка при регистрации' } finally { loading.value = false }
}
</script>

<template>
  <div class="min-h-[90vh] flex items-center justify-center relative overflow-hidden py-12">
    <!-- Decorative blobs -->
    <div class="absolute -top-40 -right-40 w-96 h-96 bg-secondary/10 rounded-full blur-[100px] animate-pulse"></div>
    <div class="absolute -bottom-40 -left-40 w-96 h-96 bg-primary/10 rounded-full blur-[100px] animate-pulse delay-700"></div>

    <div class="w-full max-w-md relative z-10">
      <div class="text-center mb-10">
        <div class="bg-gradient-to-br from-secondary to-neutral w-24 h-24 rounded-[2rem] flex items-center justify-center mx-auto mb-6 shadow-glow shadow-secondary/20 -rotate-3 hover:-rotate-6 transition-transform duration-500">
          <UserPlus class="w-10 h-10 text-white" />
        </div>
        <h1 class="text-4xl font-black tracking-tight text-secondary mb-2">Регистрация</h1>
        <p class="text-secondary/50 font-bold uppercase text-[11px] tracking-[0.2em]">Создайте свой профиль</p>
      </div>

      <div class="card bg-base-100/80 backdrop-blur-xl shadow-soft border border-white/50 rounded-[2.5rem] overflow-hidden hover:shadow-glow/10 transition-shadow duration-500">
        <div class="card-body p-10">
          <form @submit.prevent="handleRegister" class="space-y-5">
            <div v-if="error" class="alert alert-error rounded-2xl py-3 text-sm flex gap-3 shadow-lg text-white font-bold"><AlertCircle class="w-5 h-5" /><span>{{ error }}</span></div>
            
            <div class="form-control group">
              <label class="label pl-4"><span class="label-text font-black uppercase text-[10px] text-secondary/40 tracking-widest group-focus-within:text-primary transition-colors">Имя</span></label>
              <div class="relative transition-all duration-300 transform group-focus-within:scale-[1.02]">
                <User class="absolute left-5 top-4 w-5 h-5 text-secondary/30 group-focus-within:text-primary transition-colors" />
                <input v-model="name" type="text" placeholder="Иван Иванов" class="input input-lg w-full pl-14 rounded-2xl bg-base-200/50 border-transparent focus:bg-white focus:border-primary/20 focus:shadow-lg focus:shadow-primary/10 transition-all font-bold text-secondary placeholder:text-secondary/20 h-14" required />
              </div>
            </div>

            <div class="form-control group">
              <label class="label pl-4"><span class="label-text font-black uppercase text-[10px] text-secondary/40 tracking-widest group-focus-within:text-primary transition-colors">Email</span></label>
              <div class="relative transition-all duration-300 transform group-focus-within:scale-[1.02]">
                <Mail class="absolute left-5 top-4 w-5 h-5 text-secondary/30 group-focus-within:text-primary transition-colors" />
                <input v-model="email" type="email" placeholder="ivan@example.com" class="input input-lg w-full pl-14 rounded-2xl bg-base-200/50 border-transparent focus:bg-white focus:border-primary/20 focus:shadow-lg focus:shadow-primary/10 transition-all font-bold text-secondary placeholder:text-secondary/20 h-14" required />
              </div>
            </div>

            <div class="form-control group">
              <label class="label pl-4"><span class="label-text font-black uppercase text-[10px] text-secondary/40 tracking-widest group-focus-within:text-primary transition-colors">Пароль</span></label>
              <div class="relative transition-all duration-300 transform group-focus-within:scale-[1.02]">
                <Lock class="absolute left-5 top-4 w-5 h-5 text-secondary/30 group-focus-within:text-primary transition-colors" />
                <input v-model="password" type="password" placeholder="••••••••" class="input input-lg w-full pl-14 rounded-2xl bg-base-200/50 border-transparent focus:bg-white focus:border-primary/20 focus:shadow-lg focus:shadow-primary/10 transition-all font-bold text-secondary placeholder:text-secondary/20 h-14" required />
              </div>
            </div>
            
            <div class="form-control group">
              <label class="label pl-4"><span class="label-text font-black uppercase text-[10px] text-secondary/40 tracking-widest group-focus-within:text-primary transition-colors">Роль</span></label>
              <div class="relative transition-all duration-300 transform group-focus-within:scale-[1.02]">
                <ShieldCheck class="absolute left-5 top-4 w-5 h-5 text-secondary/30 group-focus-within:text-primary transition-colors" />
                <select v-model="role" class="select select-lg w-full pl-14 rounded-2xl bg-base-200/50 border-transparent focus:bg-white focus:border-primary/20 focus:shadow-lg focus:shadow-primary/10 transition-all font-bold text-secondary h-14 appearance-none cursor-pointer"><option value="client">Клиент</option><option value="chef">Повар</option><option value="courier">Курьер</option><option value="manager">Менеджер</option></select>
              </div>
            </div>

            <button type="submit" class="btn btn-secondary btn-block h-16 rounded-2xl font-black uppercase tracking-widest text-sm shadow-xl shadow-secondary/30 mt-6 hover:shadow-secondary/50 hover:-translate-y-1 transition-all border-none" :disabled="loading">
              <span v-if="loading" class="loading loading-spinner"></span>
              <span v-else>Присоединиться</span>
            </button>
          </form>
          
          <div class="divider opacity-10 my-8 font-black text-xs text-secondary/40 tracking-widest">ИЛИ</div>
          
          <div class="text-center">
            <p class="text-sm font-medium text-secondary/60 mb-4">Уже есть аккаунт?</p>
            <router-link to="/login" class="btn btn-ghost btn-outline border-2 border-base-300 hover:border-secondary hover:bg-secondary/5 hover:text-secondary w-full h-14 rounded-2xl font-black uppercase text-[10px] tracking-widest transition-all">Войти в систему</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

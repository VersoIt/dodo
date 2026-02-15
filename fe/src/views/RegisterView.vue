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
  <div class="min-h-[80vh] flex items-center justify-center px-4 py-12">
    <div class="max-w-md w-full">
      <div class="text-center mb-10">
        <div class="bg-secondary/10 w-20 h-20 rounded-3xl flex items-center justify-center mx-auto mb-6"><UserPlus class="w-10 h-10 text-secondary" /></div>
        <h1 class="text-4xl font-black tracking-tight uppercase text-secondary">Регистрация</h1>
        <p class="text-base-content/50 font-bold uppercase text-[10px] tracking-[0.2em] mt-2">Станьте частью нашей команды</p>
      </div>
      <div class="card bg-base-100 shadow-2xl border border-base-200 overflow-hidden rounded-[2.5rem]">
        <div class="card-body p-10">
          <form @submit.prevent="handleRegister" class="space-y-5">
            <div v-if="error" class="alert alert-error rounded-2xl py-3 text-sm flex gap-2"><AlertCircle class="w-4 h-4" /><span>{{ error }}</span></div>
            <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40 tracking-widest">Полное имя</span></label><div class="relative"><User class="absolute left-4 top-3.5 w-5 h-5 opacity-30" /><input v-model="name" type="text" placeholder="Иван Иванов" class="input input-bordered w-full pl-12 rounded-2xl h-12" required /></div></div>
            <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40 tracking-widest">Email адрес</span></label><div class="relative"><Mail class="absolute left-4 top-3.5 w-5 h-5 opacity-30" /><input v-model="email" type="email" placeholder="ivan@example.com" class="input input-bordered w-full pl-12 rounded-2xl h-12" required /></div></div>
            <div class="form-control"><label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40 tracking-widest">Пароль</span></label><div class="relative"><Lock class="absolute left-4 top-3.5 w-5 h-5 opacity-30" /><input v-model="password" type="password" placeholder="••••••••" class="input input-bordered w-full pl-12 rounded-2xl h-12" required /></div></div>
            <div class="form-control">
              <label class="label"><span class="label-text font-black uppercase text-[10px] opacity-40 tracking-widest">Ваша роль</span></label>
              <div class="relative"><ShieldCheck class="absolute left-4 top-3.5 w-5 h-5 opacity-30" /><select v-model="role" class="select select-bordered w-full pl-12 rounded-2xl h-12 font-bold"><option value="client">Клиент</option><option value="chef">Повар</option><option value="courier">Курьер</option><option value="manager">Менеджер</option></select></div>
            </div>
            <button type="submit" class="btn btn-secondary btn-block h-14 rounded-2xl font-black uppercase shadow-lg shadow-secondary/20 mt-6" :disabled="loading"><span v-if="loading" class="loading loading-spinner"></span><span v-else>Зарегистрироваться</span></button>
          </form>
          <div class="divider opacity-30 my-8">ИЛИ</div>
          <div class="text-center"><p class="text-sm text-base-content/60">Уже есть аккаунт?</p><router-link to="/login" class="link link-secondary font-black uppercase text-[10px] tracking-widest mt-2 block">Войти в систему</router-link></div>
        </div>
      </div>
    </div>
  </div>
</template>

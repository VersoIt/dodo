import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import { logger } from '../api/logger'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<any>(JSON.parse(localStorage.getItem('user') || 'null'))

  const isAuthenticated = computed(() => !!token.value)

  function setToken(newToken: string) {
    logger.info('Setting new auth token')
    token.value = newToken
    localStorage.setItem('token', newToken)
    axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
  }

  function logout() {
    logger.info('Logging out user', { userId: user.value?.id })
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    delete axios.defaults.headers.common['Authorization']
  }

  async function fetchMe() {
    if (!token.value) return
    logger.debug('Fetching user profile')
    try {
      const response = await axios.get('/api/v1/auth/me')
      if (response.data.success) {
        user.value = response.data.data
        localStorage.setItem('user', JSON.stringify(user.value))
        logger.info('User profile loaded', { email: user.value.email })
      }
    } catch (error) {
      logger.error('Failed to fetch user info', error)
      logout()
    }
  }

  async function login(email: string, password: string): Promise<boolean> {
    try {
      const response = await axios.post('/api/v1/auth/login', { email, password })
      if (response.data.success && response.data.data.token) {
        setToken(response.data.data.token)
        await fetchMe()
        return true
      }
      return false
    } catch (error) {
      logger.error('Login method failed', error)
      return false
    }
  }

  if (token.value) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return { token, user, isAuthenticated, setToken, logout, fetchMe, login }
})

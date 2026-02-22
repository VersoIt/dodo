import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../api'
import { logger } from '../api/logger'

export const useAuthStore = defineStore('auth', () => {
  const getSafeUser = () => {
    try {
      const stored = localStorage.getItem('user')
      if (stored && stored !== 'undefined') {
        return JSON.parse(stored)
      }
    } catch (e) {
      localStorage.removeItem('user')
    }
    return null
  }

  const token = ref(localStorage.getItem('token') || '')
  const user = ref<any>(getSafeUser())

  const isAuthenticated = computed(() => !!token.value)

  function setToken(newToken: string) {
    logger.info('Setting new auth token')
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUser(newUser: any) {
    if (newUser) {
      user.value = newUser
      localStorage.setItem('user', JSON.stringify(newUser))
    } else {
      user.value = null
      localStorage.setItem('user', 'null')
    }
  }

  function logout() {
    logger.info('Logging out user', { userId: user.value?.id })
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function fetchMe() {
    if (!token.value) return
    logger.debug('Fetching user profile')
    try {
      const response = await authApi.getMe()
      if (response.success && response.data) {
        setUser(response.data)
        logger.info('User profile loaded', { email: response.data.email })
      }
    } catch (error) {
      logger.error('Failed to fetch user info', error)
      logout()
    }
  }

  async function login(email: string, password: string): Promise<boolean> {
    try {
      const response = await authApi.login({ email, password })
      if (response.success && response.data && response.data.token) {
        setToken(response.data.token)
        
        if (response.data.user) {
          setUser(response.data.user)
        } else {
          // If the backend only returns a token, we need to fetch the user profile
          logger.info('Login response missing user data, fetching profile...')
          await fetchMe()
        }

        if (user.value) {
           return true
        } else {
           logger.error('Login succeeded but failed to load user profile')
           logout() // Clean up invalid state
           return false
        }
      }
      logger.warn('Login failed or API response is missing token', { response })
      return false
    } catch (error) {
      logger.error('Login method failed', error)
      return false
    }
  }

  return { token, user, isAuthenticated, setToken, logout, fetchMe, login, setUser }
})

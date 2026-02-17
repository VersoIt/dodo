<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch, computed } from 'vue'
import axios from 'axios'
import { Send, Check, CheckCheck, Clock, User, AlertCircle } from 'lucide-vue-next'
import { useAuthStore } from '../store/auth'

const props = defineProps<{
  orderId: string | number
}>()

const authStore = useAuthStore()
const myId = computed(() => authStore.user?.id || authStore.user?.user_id)
const myRole = computed(() => authStore.user?.role || 'client')

interface Message {
  message_id?: number
  request_id?: string
  order_id: string
  text: string
  role: string
  sender_id?: string
  created_at?: string
  status: 'pending' | 'sent' | 'read'
}

const messages = ref<Message[]>([])
const newMessage = ref('')
const isConnected = ref(false)
const lastSeenId = ref(0)
const scrollContainer = ref<HTMLElement | null>(null)
const socket = ref<WebSocket | null>(null)
const reconnectAttempts = ref(0)

const scrollToBottom = async () => {
  await nextTick()
  if (scrollContainer.value) {
    scrollContainer.value.scrollTop = scrollContainer.value.scrollHeight
  }
}

const fetchHistory = async () => {
  try {
    const response = await axios.get(`/api/v1/chat/history?order_id=${props.orderId}&limit=50`)
    messages.value = response.data.map((m: any) => ({
      ...m,
      status: m.is_read ? 'read' : 'sent'
    }))
    if (messages.value.length > 0) {
      lastSeenId.value = Math.max(...messages.value.map(m => m.message_id || 0))
    }
    scrollToBottom()
  } catch (err) {
    console.error('Failed to fetch chat history', err)
  }
}

const syncMessages = async () => {
  try {
    const response = await axios.get(`/api/v1/chat/sync?order_id=${props.orderId}&after_id=${lastSeenId.value}`)
    const newMsgs = response.data.map((m: any) => ({
      ...m,
      status: m.is_read ? 'read' : 'sent'
    }))
    if (newMsgs.length > 0) {
      messages.value.push(...newMsgs)
      lastSeenId.value = Math.max(lastSeenId.value, ...newMsgs.map((m: any) => m.message_id || 0))
      scrollToBottom()
    }
  } catch (err) {
    console.error('Failed to sync messages', err)
  }
}

const connectWebSocket = () => {
  const token = localStorage.getItem('token')
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const wsUrl = `${protocol}//${host}/ws/chat?token=${token}&order_id=${props.orderId}`
  
  socket.value = new WebSocket(wsUrl)

  socket.value.onopen = () => {
    console.log('Chat WebSocket connected')
    isConnected.value = true
    reconnectAttempts.value = 0
    if (lastSeenId.value > 0) {
      syncMessages()
    }
  }

  socket.value.onmessage = (event) => {
    const data = JSON.parse(event.data)
    
    if (data.event === 'message_ack') {
      const msg = messages.value.find(m => m.request_id === data.request_id)
      if (msg) {
        msg.message_id = data.message_id
        msg.status = 'sent'
        lastSeenId.value = Math.max(lastSeenId.value, data.message_id)
      }
    } else if (data.event === 'new_message') {
      if (!messages.value.some(m => m.message_id === data.message_id)) {
        messages.value.push({
          ...data,
          status: 'sent'
        })
        lastSeenId.value = Math.max(lastSeenId.value, data.message_id)
        scrollToBottom()
        
        socket.value?.send(JSON.stringify({
          action: 'read',
          message_id: data.message_id
        }))
      }
    }
  }

  socket.value.onclose = () => {
    isConnected.value = false
    const timeout = Math.min(1000 * Math.pow(2, reconnectAttempts.value), 30000)
    reconnectAttempts.value++
    setTimeout(connectWebSocket, timeout)
  }
}

const sendMessage = () => {
  if (!newMessage.value.trim()) return
  
  if (!isConnected.value) {
    console.warn('Cannot send: WebSocket is disconnected')
    return
  }

  const requestId = `req-${Date.now()}`
  const tempMsg: Message = {
    request_id: requestId,
    order_id: String(props.orderId),
    text: newMessage.value,
    role: myRole.value,
    sender_id: myId.value,
    status: 'pending',
    created_at: new Date().toISOString()
  }

  messages.value.push(tempMsg)
  newMessage.value = ''
  scrollToBottom()

  socket.value?.send(JSON.stringify({
    action: 'send_message',
    request_id: requestId,
    text: tempMsg.text
  }))
}

onMounted(() => {
  fetchHistory()
  connectWebSocket()
})

onUnmounted(() => {
  socket.value?.close()
})

const getRoleLabel = (role: string) => {
  switch (role) {
    case 'courier': return 'Курьер'
    case 'support': return 'Поддержка'
    case 'manager': return 'Менеджер'
    default: return ''
  }
}

const isOwn = (msg: Message) => {
  if (msg.sender_id && myId.value) return String(msg.sender_id) === String(myId.value)
  return msg.role === myRole.value
}
</script>

<template>
  <div class="card bg-base-100 shadow-xl border border-base-200 flex flex-col h-[500px]">
    <div class="card-body p-0 flex flex-col h-full overflow-hidden">
      <!-- Header -->
      <div class="p-4 border-b border-base-200 flex justify-between items-center bg-base-200/30">
        <h3 class="font-bold flex items-center gap-2">
          Чат по заказу
          <div v-if="isConnected" class="badge badge-success badge-xs"></div>
          <div v-else class="badge badge-error badge-xs animate-pulse"></div>
        </h3>
        <span class="text-[10px] uppercase tracking-widest opacity-50 font-black">Pizza Live Chat</span>
      </div>

      <!-- Messages Area -->
      <div ref="scrollContainer" class="flex-1 overflow-y-auto p-4 space-y-4 scroll-smooth">
        <div v-for="(msg, index) in messages" :key="msg.message_id || msg.request_id" 
          class="chat" :class="isOwn(msg) ? 'chat-end' : 'chat-start'">
          
          <div class="chat-header opacity-50 text-[10px] mb-1 font-bold uppercase tracking-tighter">
            {{ isOwn(msg) ? 'Вы' : getRoleLabel(msg.role) }}
          </div>

          <div class="chat-bubble shadow-sm text-sm" 
            :class="isOwn(msg) ? 'chat-bubble-primary' : 'chat-bubble-secondary'">
            {{ msg.text }}
          </div>

          <div class="chat-footer opacity-50 flex items-center gap-1 mt-1">
            <time class="text-[9px]">{{ msg.created_at ? new Date(msg.created_at).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) : '' }}</time>
            <template v-if="isOwn(msg)">
              <Clock v-if="msg.status === 'pending'" class="w-3 h-3" />
              <Check v-else-if="msg.status === 'sent'" class="w-3 h-3" />
              <CheckCheck v-else-if="msg.status === 'read'" class="w-3 h-3 text-primary" />
            </template>
          </div>
        </div>

        <div v-if="messages.length === 0" class="flex flex-col items-center justify-center h-full opacity-20 grayscale">
          <User class="w-12 h-12 mb-2" />
          <p class="font-black uppercase tracking-widest text-xs">Нет сообщений</p>
        </div>
      </div>

      <!-- Input Area -->
      <div class="p-4 border-t border-base-200 bg-base-200/30">
        <form @submit.prevent="sendMessage" class="flex gap-2">
          <input 
            v-model="newMessage"
            type="text" 
            placeholder="Напишите сообщение..." 
            class="input input-bordered input-sm flex-1 rounded-xl focus:ring-2 ring-primary/20"
          />
          <button 
            type="submit" 
            class="btn btn-primary btn-sm btn-square rounded-xl shadow-lg shadow-primary/20"
            :disabled="!isConnected || !newMessage.trim()"
          >
            <Send class="w-4 h-4" />
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-bubble {
  @apply min-h-0 py-2 px-4 rounded-2xl;
}
.chat-start .chat-bubble {
  @apply rounded-bl-none;
}
.chat-end .chat-bubble {
  @apply rounded-br-none;
}
</style>

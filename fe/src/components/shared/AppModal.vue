<script setup lang="ts">
import { X } from 'lucide-vue-next'

interface Props {
  show: boolean
  title?: string
  maxWidth?: 'sm' | 'md' | 'lg' | '2xl'
  closeOnBackdrop?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  show: false,
  maxWidth: 'md',
  closeOnBackdrop: true
})

const emit = defineEmits(['close'])

const close = () => {
  emit('close')
}

const handleBackdropClick = () => {
  if (props.closeOnBackdrop) close()
}

const maxWidthClass = {
  'sm': 'max-w-sm',
  'md': 'max-w-md',
  'lg': 'max-w-lg',
  '2xl': 'max-w-2xl'
}[props.maxWidth]
</script>

<template>
  <Transition name="modal-fade">
    <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
      <!-- Backdrop -->
      <div 
        class="fixed inset-0 bg-black/80 transition-opacity" 
        @click="handleBackdropClick"
      ></div>
      
      <!-- Content -->
      <Transition name="modal-zoom" appear>
        <div 
          class="relative bg-base-100 w-full rounded-[2.5rem] shadow-2xl border border-white/5 overflow-hidden transition-all"
          :class="maxWidthClass"
        >
          <!-- Optional Header -->
          <div v-if="title" class="px-8 pt-8 flex justify-between items-center">
            <h3 class="font-black text-2xl uppercase tracking-tighter">{{ title }}</h3>
            <button @click="close" class="btn btn-ghost btn-circle btn-sm bg-base-200/50">
              <X class="w-4 h-4" />
            </button>
          </div>

          <slot></slot>
        </div>
      </Transition>
    </div>
  </Transition>
</template>

<style scoped>
.modal-fade-enter-active, .modal-fade-leave-active {
  transition: opacity 0.15s ease;
}
.modal-fade-enter-from, .modal-fade-leave-to {
  opacity: 0;
}

.modal-zoom-enter-active {
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
.modal-zoom-leave-active {
  transition: all 0.15s ease-in;
}
.modal-zoom-enter-from, .modal-zoom-leave-to {
  opacity: 0;
  transform: scale(0.97) translateY(8px);
}
</style>

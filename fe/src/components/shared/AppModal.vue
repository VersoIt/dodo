<script setup lang="ts">
import { X } from 'lucide-vue-next'

interface Props {
  show: boolean
  title?: string
  maxWidth?: 'sm' | 'md' | 'lg' | '2xl' | '7xl'
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
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  '2xl': 'max-w-2xl',
  '7xl': 'max-w-7xl'
}[props.maxWidth]
</script>

<template>
  <Transition name="modal-fade">
    <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
      <div class="fixed inset-0 bg-black/80 transition-opacity" @click="handleBackdropClick"></div>

      <Transition name="modal-zoom" appear>
        <div
          class="relative w-full overflow-hidden rounded-[2.5rem] border border-white/5 bg-base-100 shadow-2xl transition-all"
          :class="maxWidthClass"
        >
          <div v-if="title" class="flex items-center justify-between px-8 pt-8">
            <h3 class="text-2xl font-black uppercase tracking-tighter">{{ title }}</h3>
            <button class="btn btn-ghost btn-circle btn-sm bg-base-200/50" @click="close">
              <X class="h-4 w-4" />
            </button>
          </div>

          <slot></slot>
        </div>
      </Transition>
    </div>
  </Transition>
</template>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.15s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-zoom-enter-active {
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.modal-zoom-leave-active {
  transition: all 0.15s ease-in;
}

.modal-zoom-enter-from,
.modal-zoom-leave-to {
  opacity: 0;
  transform: scale(0.97) translateY(8px);
}
</style>

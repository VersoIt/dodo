<script setup lang="ts">
import { RouterView } from 'vue-router'
import Navbar from './components/Navbar.vue'
</script>

<template>
  <div class="min-h-screen bg-base-100 flex flex-col">
    <Navbar />
    <main class="flex-grow container mx-auto px-4 py-8">
      <RouterView />
    </main>
    <footer class="footer footer-center p-4 bg-base-300 text-base-content">
      <aside>
        <p>Copyright © 2026 - All right reserved by Pizza Good</p>
      </aside>
    </footer>
  </div>
</template>

<style>
/* Global styles if needed */
</style>

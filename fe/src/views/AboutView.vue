<script setup lang="ts">
import { Pizza, Heart, Clock, Truck, ShieldCheck, Star } from 'lucide-vue-next'

const features = [
  { icon: Heart, title: 'Сделано с любовью', desc: 'Каждая пицца готовится вручную из самых свежих ингредиентов от местных фермеров.' },
  { icon: Clock, title: 'Быстрая готовка', desc: 'Наши мастера кухни начинают работу сразу после оплаты вашего заказа.' },
  { icon: Truck, title: 'Бережная доставка', desc: 'Наши курьеры используют термосумки, чтобы пицца приехала к вам горячей.' },
  { icon: ShieldCheck, title: 'Гарантия качества', desc: 'Мы следим за каждым этапом: от выбора муки до упаковки готового блюда.' }
]
</script>

<template>
  <div class="max-w-5xl mx-auto py-12 px-4 pb-20 space-y-20">
    <section class="text-center space-y-6">
      <div class="inline-flex bg-primary/10 p-4 rounded-3xl mb-4"><Pizza class="w-12 h-12 text-primary" /></div>
      <h1 class="text-5xl font-black tracking-tight uppercase">Наша история</h1>
      <p class="text-xl text-base-content/60 max-w-2xl mx-auto leading-relaxed">Мы — команда энтузиастов, которые верят, что идеальная пицца может сделать день лучше. Наша миссия — объединять людей через вкусную еду и безупречный сервис.</p>
    </section>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <div v-for="f in features" :key="f.title" class="card bg-base-100 shadow-xl border border-base-200 p-8 rounded-[2.5rem] hover:shadow-2xl transition-all group">
        <div class="bg-base-200 w-16 h-16 rounded-2xl flex items-center justify-center mb-6 group-hover:bg-primary group-hover:text-primary-content transition-colors duration-500"><component :is="f.icon" class="w-8 h-8" /></div>
        <h3 class="text-2xl font-black mb-3 uppercase tracking-tight">{{ f.title }}</h3>
        <p class="text-base-content/60 leading-relaxed">{{ f.desc }}</p>
      </div>
    </div>

    <div class="stats stats-vertical lg:stats-horizontal shadow-2xl w-full rounded-[3rem] border border-base-200 overflow-hidden">
      <div class="stat p-10 bg-base-100"><div class="stat-figure text-primary"><Star class="w-8 h-8" /></div><div class="stat-title font-bold uppercase text-[10px] tracking-widest">Довольных клиентов</div><div class="stat-value text-primary font-black">50,000+</div><div class="stat-desc font-medium mt-1">И это только начало</div></div>
      <div class="stat p-10 bg-base-100"><div class="stat-figure text-secondary"><Pizza class="w-8 h-8" /></div><div class="stat-title font-bold uppercase text-[10px] tracking-widest">Вариантов начинки</div><div class="stat-value text-secondary font-black">40+</div><div class="stat-desc font-medium mt-1">От классики до экзотики</div></div>
      <div class="stat p-10 bg-base-100"><div class="stat-figure text-success"><Truck class="w-8 h-8" /></div><div class="stat-title font-bold uppercase text-[10px] tracking-widest">Среднее время</div><div class="stat-value text-success font-black">25 мин</div><div class="stat-desc font-medium mt-1">Быстрее всех в городе</div></div>
    </div>
  </div>
</template>

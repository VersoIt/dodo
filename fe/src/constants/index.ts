export const ORDER_STATUS = {
  CREATED: 'created',
  PAID: 'paid',
  COOKING: 'cooking',
  READY: 'ready',
  DELIVERING: 'delivering',
  COMPLETED: 'completed',
  CANCELED: 'canceled'
} as const;

export const CATEGORIES = ['Все', 'Классика', 'Премиум', 'Вегетарианская', 'Острая', 'Напитки', 'Десерты'] as const;

export const CATEGORY_MAP: Record<string, number> = {
  'Классика': 0,
  'Премиум': 1,
  'Вегетарианская': 2,
  'Острая': 3,
  'Напитки': 4,
  'Десерты': 5
};

export const HERO_IMAGE = "https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=1200&auto=format&fit=crop";

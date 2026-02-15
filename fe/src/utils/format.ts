import { HERO_IMAGE } from '../constants'

/**
 * Formats number to Russian Ruble currency string
 */
export const formatPrice = (value: number | undefined | null): string => {
  if (value === undefined || value === null) return '0 ₽'
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0
  }).format(value)
}

/**
 * Fallback for broken images
 */
export const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  if (img.src !== HERO_IMAGE) {
    img.src = HERO_IMAGE
  }
}

/**
 * Get short ID for display
 */
export const shortId = (id: string): string => {
  if (!id) return ''
  return id.slice(0, 8).toUpperCase()
}

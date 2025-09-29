<template>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <!-- Total Recommendations -->
      <div class="card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Total Recommendations</p>
            <p class="text-3xl font-bold text-gray-900">
              {{ animatedTotal.toLocaleString() }}
            </p>
            <p v-if="store.loading" class="text-xs text-blue-600 mt-1">
              Loading {{ store.loadingPercentage }}%
            </p>
          </div>
          <div class="flex items-center text-gray-400">
            <BarChart3Icon class="h-8 w-8" />
          </div>
        </div>
      </div>
  
      <!-- Upgrades -->
      <div class="card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Upgrades</p>
            <p class="text-3xl font-bold text-green-600">
              {{ animatedUpgrades.toLocaleString() }}
            </p>
            <p v-if="store.currentRecommendationCount > 0" class="text-xs text-gray-500 mt-1">
              {{ Math.round((store.upgrades / store.currentRecommendationCount) * 100) }}% of total
            </p>
          </div>
          <div class="flex items-center text-green-500">
            <TrendingUpIcon class="h-8 w-8" />
          </div>
        </div>
      </div>
  
      <!-- Downgrades -->
      <div class="card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Downgrades</p>
            <p class="text-3xl font-bold text-red-600">
              {{ animatedDowngrades.toLocaleString() }}
            </p>
            <p v-if="store.currentRecommendationCount > 0" class="text-xs text-gray-500 mt-1">
              {{ Math.round((store.downgrades / store.currentRecommendationCount) * 100) }}% of total
            </p>
          </div>
          <div class="flex items-center text-red-500">
            <TrendingDownIcon class="h-8 w-8" />
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, watch, onMounted } from 'vue'
  import { useMainStore } from '@/stores/index'
  
  const BarChart3Icon = () => '📊'
  const TrendingUpIcon = () => '📈'
  const TrendingDownIcon = () => '📉'
  
  const store = useMainStore()
  
  // Animated counters
  const animatedTotal = ref(store.recommendations.length)
  const animatedUpgrades = ref(store.upgrades)
  const animatedDowngrades = ref(store.downgrades)
  
  // Animate numbers
  function animateValue(from: number, to: number, duration: number = 1000): Promise<void> {
    return new Promise(resolve => {
      const start = Date.now()
      const timer = setInterval(() => {
        const progress = Math.min((Date.now() - start) / duration, 1)
        const current = Math.floor(from + (to - from) * progress)
        
        if (progress === 1) {
          clearInterval(timer)
          resolve()
        }
      }, 16)
    })
  }
  
  // Watch for changes and animate
  watch(() => store.currentRecommendationCount, (newVal, oldVal) => {
    animateValue(animatedTotal.value, newVal, 500).then(() => {
      animatedTotal.value = newVal
    })
  })
  
  watch(() => store.upgrades, (newVal, oldVal) => {
    animateValue(animatedUpgrades.value, newVal, 500).then(() => {
      animatedUpgrades.value = newVal
    })
  })
  
  watch(() => store.downgrades, (newVal, oldVal) => {
    animateValue(animatedDowngrades.value, newVal, 500).then(() => {
      animatedDowngrades.value = newVal
    })
  })
  </script>
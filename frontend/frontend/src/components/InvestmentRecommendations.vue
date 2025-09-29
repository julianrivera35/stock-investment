<template>
    <div class="space-y-6">
      <!-- Header -->
      <div class="card p-6">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="text-xl font-bold text-gray-900">🚀 AI Investment Recommendations</h2>
            <p class="text-sm text-gray-600 mt-1">
              Based on analysis of {{ store.recommendations.length.toLocaleString() }} analyst recommendations
            </p>
          </div>
          <div class="text-right">
            <div class="text-sm text-gray-500">Last updated</div>
            <div class="text-sm font-medium text-gray-900">{{ lastUpdated }}</div>
          </div>
        </div>
        
        <!-- Algorithm Info -->
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <h3 class="text-sm font-medium text-blue-900 mb-2">📊 How it works:</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-xs text-blue-800">
            <div>• <strong>Recommendation Volume</strong> (15%): More analyst coverage = higher confidence</div>
            <div>• <strong>Upgrade Ratio</strong> (25%): Percentage of positive rating changes</div>
            <div>• <strong>Target Upside</strong> (20%): Potential price appreciation</div>
            <div>• <strong>Analyst Consensus</strong> (20%): Agreement among analysts</div>
            <div>• <strong>Recency</strong> (10%): How recent the recommendations are</div>
            <div>• <strong>Brokerage Quality</strong> (10%): Reputation of recommending firms</div>
          </div>
        </div>
      </div>
  
      <!-- Top Picks -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div 
          v-for="(stock, index) in topPicks" 
          :key="stock.ticker"
          class="card p-6 hover:shadow-lg transition-shadow"
        >
          <!-- Rank Badge -->
          <div class="flex items-start justify-between mb-4">
            <div class="flex items-center">
              <div 
                :class="getRankBadgeClass(index)"
                class="w-8 h-8 rounded-full flex items-center justify-center text-white font-bold text-sm mr-3"
              >
                {{ index + 1 }}
              </div>
              <div>
                <h3 class="font-bold text-lg text-gray-900">{{ stock.ticker }}</h3>
                <p class="text-sm text-gray-600">{{ stock.companyName }}</p>
              </div>
            </div>
            <div class="text-right">
              <div 
                :class="getRecommendationClass(stock.recommendation)"
                class="px-3 py-1 rounded-full text-xs font-medium"
              >
                {{ stock.recommendation }}
              </div>
              <div class="text-lg font-bold text-gray-900 mt-1">
                {{ stock.totalScore.toFixed(1) }}/100
              </div>
            </div>
          </div>
  
          <!-- Key Metrics -->
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div class="text-center p-3 bg-gray-50 rounded-lg">
              <div class="text-lg font-bold text-gray-900">{{ stock.details.totalRecommendations }}</div>
              <div class="text-xs text-gray-600">Recommendations</div>
            </div>
            <div class="text-center p-3 bg-gray-50 rounded-lg">
              <div class="text-lg font-bold text-green-600">{{ stock.details.upgrades }}</div>
              <div class="text-xs text-gray-600">Upgrades</div>
            </div>
            <div class="text-center p-3 bg-gray-50 rounded-lg">
              <div class="text-lg font-bold text-blue-600">
                ${{ stock.details.averageTargetPrice.toFixed(2) }}
              </div>
              <div class="text-xs text-gray-600">Avg Target</div>
            </div>
            <div class="text-center p-3 bg-gray-50 rounded-lg">
              <div class="text-lg font-bold text-purple-600">{{ stock.details.uniqueBrokerages }}</div>
              <div class="text-xs text-gray-600">Brokerages</div>
            </div>
          </div>
  
          <!-- Score Breakdown -->
          <div class="space-y-2">
            <h4 class="text-sm font-medium text-gray-700">Score Breakdown:</h4>
            
            <div class="space-y-1">
              <div class="flex justify-between items-center">
                <span class="text-xs text-gray-600">Upgrade Ratio</span>
                <span class="text-xs font-medium">{{ stock.metrics.upgradeRatio.toFixed(1) }}%</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-1.5">
                <div 
                  class="bg-blue-600 h-1.5 rounded-full transition-all duration-300"
                  :style="{ width: `${stock.metrics.upgradeRatio}%` }"
                ></div>
              </div>
            </div>
  
            <div class="space-y-1">
              <div class="flex justify-between items-center">
                <span class="text-xs text-gray-600">Analyst Consensus</span>
                <span class="text-xs font-medium">{{ stock.metrics.analystConsensus.toFixed(1) }}%</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-1.5">
                <div 
                  class="bg-green-600 h-1.5 rounded-full transition-all duration-300"
                  :style="{ width: `${stock.metrics.analystConsensus}%` }"
                ></div>
              </div>
            </div>
  
            <div class="space-y-1">
              <div class="flex justify-between items-center">
                <span class="text-xs text-gray-600">Brokerage Quality</span>
                <span class="text-xs font-medium">{{ stock.metrics.brokerageQuality.toFixed(1) }}%</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-1.5">
                <div 
                  class="bg-purple-600 h-1.5 rounded-full transition-all duration-300"
                  :style="{ width: `${stock.metrics.brokerageQuality}%` }"
                ></div>
              </div>
            </div>
          </div>
  
          <!-- Last Activity -->
          <div class="mt-4 pt-4 border-t border-gray-200">
            <div class="flex justify-between items-center text-xs text-gray-600">
              <span>Last recommendation:</span>
              <span>{{ stock.details.daysFromLastRecommendation }} days ago</span>
            </div>
          </div>
        </div>
      </div>
  
      <!-- Full Analysis Button -->
      <div class="text-center">
        <button 
          @click="showFullAnalysis = !showFullAnalysis"
          class="btn-secondary"
        >
          {{ showFullAnalysis ? 'Hide' : 'Show' }} Full Analysis ({{ allAnalysis.length }} stocks)
        </button>
      </div>
  
      <!-- Full Analysis Table -->
      <div v-if="showFullAnalysis" class="card overflow-hidden">
        <div class="p-4 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">Complete Stock Analysis</h3>
        </div>
        
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="table-header">Rank</th>
                <th class="table-header">Stock</th>
                <th class="table-header">Score</th>
                <th class="table-header">Recommendation</th>
                <th class="table-header">Recommendations</th>
                <th class="table-header">Upgrades</th>
                <th class="table-header">Avg Target</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="(stock, index) in allAnalysis" :key="stock.ticker" class="hover:bg-gray-50">
                <td class="table-cell">
                  <span 
                    :class="getRankTextClass(index)"
                    class="font-medium"
                  >
                    #{{ index + 1 }}
                  </span>
                </td>
                <td class="table-cell">
                  <div>
                    <div class="font-medium text-gray-900">{{ stock.ticker }}</div>
                    <div class="text-sm text-gray-500">{{ stock.companyName }}</div>
                  </div>
                </td>
                <td class="table-cell">
                  <span class="font-bold text-gray-900">{{ stock.totalScore.toFixed(1) }}</span>
                </td>
                <td class="table-cell">
                  <span 
                    :class="getRecommendationClass(stock.recommendation)"
                    class="px-2 py-1 rounded-full text-xs font-medium"
                  >
                    {{ stock.recommendation }}
                  </span>
                </td>
                <td class="table-cell">{{ stock.details.totalRecommendations }}</td>
                <td class="table-cell">
                  <span class="text-green-600 font-medium">{{ stock.details.upgrades }}</span>
                </td>
                <td class="table-cell">${{ stock.details.averageTargetPrice.toFixed(2) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed } from 'vue'
  import { useMainStore } from '@/stores'
  import { analyzeStockInvestments, getTopInvestmentPicks } from '@/utils/investmentAlgorithm'
  
  const store = useMainStore()
  const showFullAnalysis = ref(false)
  
  const allAnalysis = computed(() => {
    if (store.recommendations.length === 0) return []
    return analyzeStockInvestments(store.recommendations)
  })
  
  const topPicks = computed(() => {
    return allAnalysis.value.slice(0, 6)
  })
  
  const lastUpdated = computed(() => {
    return new Date().toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  })
  
  function getRankBadgeClass(index: number) {
    if (index === 0) return 'bg-yellow-500'
    if (index === 1) return 'bg-gray-400'
    if (index === 2) return 'bg-amber-600'
    return 'bg-blue-600'
  }
  
  function getRankTextClass(index: number) {
    if (index < 5) return 'text-green-600'
    if (index < 10) return 'text-blue-600'
    if (index < 20) return 'text-gray-600'
    return 'text-gray-400'
  }
  
  function getRecommendationClass(recommendation: string) {
    switch (recommendation) {
      case 'Strong Buy': return 'bg-green-100 text-green-800'
      case 'Buy': return 'bg-green-50 text-green-700'
      case 'Hold': return 'bg-yellow-100 text-yellow-800'
      case 'Sell': return 'bg-red-50 text-red-700'
      case 'Strong Sell': return 'bg-red-100 text-red-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }
  </script>
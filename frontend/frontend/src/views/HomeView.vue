<template>
    <div class="min-h-screen bg-gray-50">
        <div class="max-w-7xl mx-auto px-4 py-8">
            <h1 class="text-3xl font-bold text-gray-900 mb-8">AnalystHub Dashboard</h1>

            <!-- Loading State -->
            <LoadingSpinner v-if="store.loading && !store.isInitialLoadComplete" :percentage="store.loadingPercentage"
                :current="store.currentRecommendationCount" :total="store.totalRecommendations" />

            <!-- Content (shows when 50% loaded) -->
            <div v-else>
                <StatsCards />

                <!-- Loading indicator for remaining data -->
                <div v-if="store.loading" class="mb-6">
                    <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
                        <div class="flex items-center">
                            <div
                                class="animate-spin h-4 w-4 border-2 border-blue-600 border-t-transparent rounded-full mr-3">
                            </div>
                            <span class="text-blue-700 text-sm">
                                Loading remaining data... {{ store.loadingPercentage }}% complete
                            </span>
                        </div>
                    </div>
                </div>

                <!-- Investment Recommendations -->
                <InvestmentRecommendations v-if="store.recommendations.length > 0" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useMainStore } from '@/stores'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import StatsCards from '@/components/StatsCards.vue'
import InvestmentRecommendations from '@/components/InvestmentRecommendations.vue'

const store = useMainStore()

onMounted(() => {
    // Only fetch if no data exists
    if (!store.dataFetched && store.recommendations.length === 0) {
        store.fetchAllRecommendations()
    }
})
</script>

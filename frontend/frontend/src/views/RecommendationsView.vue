<template>
    <div class="max-w-7xl mx-auto px-4 py-8">
        <div class="mb-8">
            <h1 class="text-3xl font-bold text-gray-900">Recommendations</h1>
            <p class="text-gray-600 mt-2">Explore the latest analyst recommendations.</p>
        </div>

        <!-- Loading State -->
        <LoadingSpinner v-if="store.loading && !store.isInitialLoadComplete" :percentage="store.loadingPercentage"
            :current="store.currentRecommendationCount" :total="store.totalRecommendations" />

        <!-- Content (shows when 50% loaded) -->
        <div v-else>
            <StatsCards />

            <!-- Table/Heatmap Toggle -->
            <div class="mb-6">
                <div class="flex space-x-1 bg-gray-100 p-1 rounded-lg w-fit">
                    <button @click="viewMode = 'table'" :class="[
                        'px-4 py-2 text-sm font-medium rounded-md transition-colors',
                        viewMode === 'table'
                            ? 'bg-white text-gray-900 shadow-sm'
                            : 'text-gray-600 hover:text-gray-900'
                    ]">
                        📊 Table
                    </button>
                    <button @click="viewMode = 'heatmap'" :class="[
                        'px-4 py-2 text-sm font-medium rounded-md transition-colors',
                        viewMode === 'heatmap'
                            ? 'bg-white text-gray-900 shadow-sm'
                            : 'text-gray-600 hover:text-gray-900'
                    ]">
                        🔥 Heatmap
                    </button>
                </div>
            </div>

            <!-- Content based on view mode -->
            <div v-if="viewMode === 'table'">
                <RecommendationTable />
            </div>

            <div v-else>
                <Heatmap />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMainStore } from '@/stores'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import StatsCards from '@/components/StatsCards.vue'
import RecommendationTable from '@/components/Tables/RecommendationTable.vue'
import Heatmap from '@/components/Heatmap/Heatmap.vue'

const store = useMainStore()
const viewMode = ref<'table' | 'heatmap'>('table')

onMounted(() => {
    // Only fetch if no data exists
    if (!store.dataFetched && store.recommendations.length === 0) {
        store.fetchAllRecommendations()
    }
})

</script>

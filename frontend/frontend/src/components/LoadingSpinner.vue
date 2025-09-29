<template>
    <div class="flex flex-col items-center justify-center min-h-[400px]">
        <!-- Spinner -->
        <div class="relative">
            <div class="w-16 h-16 border-4 border-gray-200 border-t-blue-600 rounded-full animate-spin"></div>
            <div class="absolute inset-0 flex items-center justify-center">
                <span class="text-xs font-semibold text-blue-600">{{ percentage }}%</span>
            </div>
        </div>

        <!-- Progress Bar -->
        <div class="w-64 mt-6">
            <div class="flex justify-between text-sm text-gray-600 mb-2">
                <span>Loading recommendations...</span>
                <span>{{ current }} / {{ total }}</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
                <div class="bg-blue-600 h-2 rounded-full transition-all duration-300 ease-out"
                    :style="{ width: `${percentage}%` }">
                </div>
            </div>
        </div>

        <!-- Loading Text -->
        <p class="mt-4 text-gray-600 text-center max-w-md">
            <span v-if="percentage < 50">
                Fetching stock recommendations from analysts...
            </span>
            <span v-else-if="percentage < 100">
                Almost there! Loading remaining data...
            </span>
            <span v-else>
                Loading complete!
            </span>
        </p>
    </div>
</template>

<script setup lang="ts">
interface Props {
    percentage: number;
    current: number;
    total: number;
}

defineProps<Props>();
</script>
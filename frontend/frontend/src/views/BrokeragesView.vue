<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Brokerages</h1>
      <p class="text-gray-600 mt-2">Browse and search for brokerages.</p>
    </div>

    <!-- Search -->
    <div class="mb-6">
      <div class="relative max-w-md">
        <input
          v-model="searchTerm"
          type="text"
          placeholder="Search by name..."
          class="search-input"
        />
        <div
          class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
        >
          <span class="text-gray-400">🔍</span>
        </div>
      </div>
    </div>

    <!-- Brokerages List -->
    <div class="card">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="table-header">Name</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="brokerage in filteredBrokerages" :key="brokerage.id">
              <td class="table-cell">{{ brokerage.name }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="filteredBrokerages.length === 0" class="text-center py-8">
        <p class="text-gray-500">No brokerages found.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useMainStore } from "@/stores";

const store = useMainStore();
const searchTerm = ref("");

const filteredBrokerages = computed(() => {
  if (!searchTerm.value) return store.brokerages;

  const term = searchTerm.value.toLowerCase();
  return store.brokerages.filter((brokerage) =>
    brokerage.name.toLowerCase().includes(term)
  );
});

onMounted(() => {
  store.fetchBrokerages();
});
</script>

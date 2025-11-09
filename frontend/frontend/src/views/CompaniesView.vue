<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Companies</h1>
      <p class="text-gray-600 mt-2">Browse and search for companies.</p>
    </div>

    <!-- Search -->
    <div class="mb-6">
      <div class="relative max-w-md">
        <input
          v-model="searchTerm"
          type="text"
          placeholder="Search by name or ticker..."
          class="search-input"
        />
        <div
          class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
        >
          <span class="text-gray-400">🔍</span>
        </div>
      </div>
    </div>

    <!-- Companies List -->
    <div class="card">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="table-header">Company</th>
              <th class="table-header">Ticker</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="company in filteredCompanies" :key="company.id">
              <td class="table-cell">
                <div class="flex items-center">
                  <div
                    class="h-8 w-8 bg-blue-100 rounded-full flex items-center justify-center mr-3"
                  >
                    <span class="text-blue-600 font-semibold text-sm">
                      {{ company.name.charAt(0) }}
                    </span>
                  </div>
                  {{ company.name }}
                </div>
              </td>
              <td class="table-cell font-mono font-semibold">
                {{ company.ticker }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="filteredCompanies.length === 0" class="text-center py-8">
        <p class="text-gray-500">No companies found.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useMainStore } from "@/stores";

const store = useMainStore();
const searchTerm = ref("");

const filteredCompanies = computed(() => {
  if (!searchTerm.value) return store.companies;

  const term = searchTerm.value.toLowerCase();
  return store.companies.filter(
    (company) =>
      company.name.toLowerCase().includes(term) ||
      company.ticker.toLowerCase().includes(term)
  );
});

onMounted(() => {
  store.fetchCompanies();
});
</script>

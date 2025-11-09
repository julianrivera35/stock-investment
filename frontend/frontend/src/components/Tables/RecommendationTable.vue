<template>
  <div class="card overflow-hidden">
    <!-- Search and Filters -->
    <div class="p-6 border-b border-gray-200">
      <div class="flex flex-col sm:flex-row gap-4">
        <!-- Search -->
        <div class="relative flex-1">
          <input
            v-model="searchTerm"
            type="text"
            placeholder="Search by company or ticker..."
            class="search-input"
          />
          <div
            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
          >
            <span class="text-gray-400">🔍</span>
          </div>
        </div>

        <!-- Filters -->
        <div class="flex gap-2">
          <select
            v-model="selectedBrokerage"
            class="px-3 py-2 border border-gray-300 rounded-lg text-sm"
          >
            <option value="">All Brokerages</option>
            <option
              v-for="brokerage in uniqueBrokerages"
              :key="brokerage"
              :value="brokerage"
            >
              {{ brokerage }}
            </option>
          </select>

          <select
            v-model="selectedAction"
            class="px-3 py-2 border border-gray-300 rounded-lg text-sm"
          >
            <option value="">All Actions</option>
            <option value="upgrade">Upgrades</option>
            <option value="downgrade">Downgrades</option>
            <option value="reiterate">Reiterations</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th
              @click="sortBy('company')"
              class="table-header cursor-pointer hover:bg-gray-100"
            >
              <div class="flex items-center">
                Company
                <span class="ml-1">{{ getSortIcon("company") }}</span>
              </div>
            </th>
            <th
              @click="sortBy('brokerage')"
              class="table-header cursor-pointer hover:bg-gray-100"
            >
              <div class="flex items-center">
                Brokerage
                <span class="ml-1">{{ getSortIcon("brokerage") }}</span>
              </div>
            </th>
            <th
              @click="sortBy('action')"
              class="table-header cursor-pointer hover:bg-gray-100"
            >
              <div class="flex items-center">
                Action
                <span class="ml-1">{{ getSortIcon("action") }}</span>
              </div>
            </th>
            <th class="table-header">Rating Change</th>
            <th class="table-header">Price Target Change</th>
            <th
              @click="sortBy('time')"
              class="table-header cursor-pointer hover:bg-gray-100"
            >
              <div class="flex items-center">
                Date
                <span class="ml-1">{{ getSortIcon("time") }}</span>
              </div>
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr
            v-for="recommendation in paginatedRecommendations"
            :key="recommendation.id"
            class="hover:bg-gray-50"
          >
            <!-- Company -->
            <td class="table-cell">
              <div class="flex items-center">
                <div
                  class="h-10 w-10 bg-blue-100 rounded-full flex items-center justify-center mr-3"
                >
                  <span class="text-blue-600 font-semibold text-sm">
                    {{ recommendation.company.ticker.charAt(0) }}
                  </span>
                </div>
                <div>
                  <div class="font-medium text-gray-900">
                    {{ recommendation.company.name }}
                  </div>
                  <div class="text-sm text-gray-500 font-mono">
                    {{ recommendation.company.ticker }}
                  </div>
                </div>
              </div>
            </td>

            <!-- Brokerage -->
            <td class="table-cell">
              <span v-if="recommendation.brokerage" class="text-gray-900">
                {{ recommendation.brokerage.name }}
              </span>
              <span v-else class="text-gray-400 italic">Unknown</span>
            </td>

            <!-- Action -->
            <td class="table-cell">
              <span
                :class="getActionBadgeClass(recommendation.action)"
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
              >
                {{ formatAction(recommendation.action) }}
              </span>
            </td>

            <!-- Rating Change -->
            <td class="table-cell">
              <div
                v-if="recommendation.rating_from && recommendation.rating_to"
              >
                <span class="text-gray-600">{{
                  recommendation.rating_from
                }}</span>
                <span class="mx-2 text-gray-400">→</span>
                <span class="font-medium text-gray-900">{{
                  recommendation.rating_to
                }}</span>
              </div>
              <span v-else class="text-gray-400">N/A</span>
            </td>

            <!-- Price Target Change -->
            <td class="table-cell">
              <div
                v-if="recommendation.target_from && recommendation.target_to"
              >
                <div class="flex items-center">
                  <span class="text-gray-600"
                    >${{ recommendation.target_from.toFixed(2) }}</span
                  >
                  <span class="mx-2 text-gray-400">→</span>
                  <span class="font-medium text-gray-900"
                    >${{ recommendation.target_to.toFixed(2) }}</span
                  >
                  <span
                    :class="
                      getPriceChangeClass(
                        recommendation.target_from,
                        recommendation.target_to
                      )
                    "
                    class="ml-2 text-xs"
                  >
                    {{
                      getPriceChangeText(
                        recommendation.target_from,
                        recommendation.target_to
                      )
                    }}
                  </span>
                </div>
              </div>
              <span v-else class="text-gray-400">N/A</span>
            </td>

            <!-- Date -->
            <td class="table-cell">
              <div class="text-gray-900">
                {{ formatDate(recommendation.time) }}
              </div>
              <div class="text-xs text-gray-500">
                {{ formatTime(recommendation.time) }}
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div
      class="px-6 py-4 border-t border-gray-200 flex items-center justify-between"
    >
      <div class="text-sm text-gray-700">
        Showing {{ (currentPage - 1) * pageSize + 1 }} to
        {{ Math.min(currentPage * pageSize, filteredRecommendations.length) }}
        of {{ filteredRecommendations.length }} recommendations
      </div>

      <div class="flex items-center space-x-2">
        <button
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="px-3 py-2 text-sm border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
        >
          Previous
        </button>

        <span class="px-3 py-2 text-sm">
          Page {{ currentPage }} of {{ totalPages }}
        </span>

        <button
          @click="currentPage++"
          :disabled="currentPage === totalPages"
          class="px-3 py-2 text-sm border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
        >
          Next
        </button>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="filteredRecommendations.length === 0" class="text-center py-12">
      <div class="text-gray-400 text-lg mb-2">📊</div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">
        No recommendations found
      </h3>
      <p class="text-gray-600">Try adjusting your search or filter criteria.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useMainStore } from "@/stores";

const store = useMainStore();

// State
const searchTerm = ref("");
const selectedBrokerage = ref("");
const selectedAction = ref("");
const sortField = ref<string>("time");
const sortDirection = ref<"asc" | "desc">("desc");
const currentPage = ref(1);
const pageSize = ref(25);

// Computed
const uniqueBrokerages = computed(() => {
  const brokerages = store.recommendations
    .map((r) => r.brokerage?.name)
    .filter(Boolean)
    .filter((value, index, self) => self.indexOf(value) === index);
  return brokerages.sort();
});

const filteredRecommendations = computed(() => {
  let filtered = store.recommendations;

  // Search filter
  if (searchTerm.value) {
    const term = searchTerm.value.toLowerCase();
    filtered = filtered.filter(
      (r) =>
        r.company.name.toLowerCase().includes(term) ||
        r.company.ticker.toLowerCase().includes(term) ||
        r.brokerage?.name.toLowerCase().includes(term)
    );
  }

  // Brokerage filter
  if (selectedBrokerage.value) {
    filtered = filtered.filter(
      (r) => r.brokerage?.name === selectedBrokerage.value
    );
  }

  // Action filter
  if (selectedAction.value) {
    filtered = filtered.filter((r) => {
      const action = r.action.toLowerCase();
      switch (selectedAction.value) {
        case "upgrade":
          return action.includes("raised") || action.includes("upgrade");
        case "downgrade":
          return action.includes("lowered") || action.includes("downgrade");
        case "reiterate":
          return action.includes("reiterat");
        default:
          return true;
      }
    });
  }

  // Sort
  filtered.sort((a, b) => {
    let aVal, bVal;

    switch (sortField.value) {
      case "company":
        aVal = a.company.name;
        bVal = b.company.name;
        break;
      case "brokerage":
        aVal = a.brokerage?.name || "";
        bVal = b.brokerage?.name || "";
        break;
      case "action":
        aVal = a.action;
        bVal = b.action;
        break;
      case "time":
        aVal = new Date(a.time);
        bVal = new Date(b.time);
        break;
      default:
        return 0;
    }

    if (aVal < bVal) return sortDirection.value === "asc" ? -1 : 1;
    if (aVal > bVal) return sortDirection.value === "asc" ? 1 : -1;
    return 0;
  });

  return filtered;
});

const totalPages = computed(() =>
  Math.ceil(filteredRecommendations.value.length / pageSize.value)
);

const paginatedRecommendations = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return filteredRecommendations.value.slice(start, end);
});

// Methods
function sortBy(field: string) {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === "asc" ? "desc" : "asc";
  } else {
    sortField.value = field;
    sortDirection.value = "asc";
  }
  currentPage.value = 1;
}

function getSortIcon(field: string) {
  if (sortField.value !== field) return "↕️";
  return sortDirection.value === "asc" ? "↑" : "↓";
}

function getActionBadgeClass(action: string) {
  const actionLower = action.toLowerCase();
  if (actionLower.includes("raised") || actionLower.includes("upgrade")) {
    return "bg-green-100 text-green-800";
  } else if (
    actionLower.includes("lowered") ||
    actionLower.includes("downgrade")
  ) {
    return "bg-red-100 text-red-800";
  } else {
    return "bg-gray-100 text-gray-800";
  }
}

function formatAction(action: string) {
  return action.charAt(0).toUpperCase() + action.slice(1);
}

function getPriceChangeClass(from: number, to: number) {
  if (to > from) return "text-green-600";
  if (to < from) return "text-red-600";
  return "text-gray-600";
}

function getPriceChangeText(from: number, to: number) {
  const change = to - from;
  const percentage = ((change / from) * 100).toFixed(1);
  return change > 0 ? `+${percentage}%` : `${percentage}%`;
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

function formatTime(dateString: string) {
  return new Date(dateString).toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
  });
}

// Reset page when filters change
function resetPage() {
  currentPage.value = 1;
}

// Watch for filter changes
import { watch } from "vue";
watch([searchTerm, selectedBrokerage, selectedAction], resetPage);
</script>

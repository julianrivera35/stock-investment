<template>
  <div class="card">
    <div class="p-6 border-b border-gray-200">
      <h3 class="text-lg font-semibold text-gray-900 mb-2">
        Recommendations Heatmap
      </h3>
      <p class="text-sm text-gray-600">
        Number of recommendations by brokerage and company
      </p>
    </div>

    <div class="p-6">
      <!-- Controls -->
      <div class="mb-6 flex flex-wrap gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Show Top Companies
          </label>
          <select
            v-model="topCompaniesCount"
            class="px-3 py-2 border border-gray-300 rounded-lg text-sm"
          >
            <option :value="5">Top 5</option>
            <option :value="10">Top 10</option>
            <option :value="15">Top 15</option>
            <option :value="20">Top 20</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Show Top Brokerages
          </label>
          <select
            v-model="topBrokeragesCount"
            class="px-3 py-2 border border-gray-300 rounded-lg text-sm"
          >
            <option :value="5">Top 5</option>
            <option :value="10">Top 10</option>
            <option :value="15">Top 15</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Color Scale
          </label>
          <select
            v-model="colorScale"
            class="px-3 py-2 border border-gray-300 rounded-lg text-sm"
          >
            <option value="blue">Blue</option>
            <option value="green">Green</option>
            <option value="red">Red</option>
            <option value="purple">Purple</option>
          </select>
        </div>
      </div>

      <!-- HeatMap -->
      <div class="overflow-x-auto">
        <div class="inline-block min-w-full">
          <!-- Header with company tickers -->
          <div class="flex">
            <!-- Empty corner cell -->
            <div
              class="w-32 h-12 border border-gray-200 bg-gray-50 flex items-center justify-center"
            >
              <span class="text-xs font-medium text-gray-500">Brokerage</span>
            </div>

            <!-- Company headers -->
            <div
              v-for="company in topCompanies"
              :key="company.ticker"
              class="w-20 h-12 border border-gray-200 bg-gray-50 flex items-center justify-center"
            >
              <span
                class="text-xs font-medium text-gray-700 transform -rotate-45 origin-center"
                :title="company.name"
              >
                {{ company.ticker }}
              </span>
            </div>
          </div>

          <!-- Brokerage rows -->
          <div
            v-for="brokerage in topBrokerages"
            :key="brokerage.name"
            class="flex"
          >
            <!-- Brokerage name cell -->
            <div
              class="w-32 h-12 border border-gray-200 bg-gray-50 flex items-center px-2"
            >
              <span
                class="text-xs font-medium text-gray-700 truncate"
                :title="brokerage.name"
              >
                {{ brokerage.name }}
              </span>
            </div>

            <!-- Data cells -->
            <div
              v-for="company in topCompanies"
              :key="`${brokerage.name}-${company.ticker}`"
              class="w-20 h-12 border border-gray-200 flex items-center justify-center cursor-pointer transition-all hover:ring-2 hover:ring-blue-500"
              :style="getCellStyle(brokerage.name, company.ticker)"
              :title="getCellTooltip(brokerage.name, company.ticker)"
              @click="showCellDetails(brokerage.name, company.ticker)"
            >
              <span
                class="text-xs font-medium"
                :class="getCellTextClass(brokerage.name, company.ticker)"
              >
                {{ getCellValue(brokerage.name, company.ticker) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Legend -->
      <div class="mt-6 flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <span class="text-sm text-gray-600">Recommendations:</span>
          <div class="flex items-center space-x-2">
            <div
              class="w-4 h-4 border border-gray-300"
              :style="{ backgroundColor: getColorForValue(0) }"
            ></div>
            <span class="text-xs text-gray-500">0</span>
          </div>
          <div class="flex items-center space-x-2">
            <div
              class="w-4 h-4 border border-gray-300"
              :style="{ backgroundColor: getColorForValue(maxValue / 2) }"
            ></div>
            <span class="text-xs text-gray-500">{{
              Math.round(maxValue / 2)
            }}</span>
          </div>
          <div class="flex items-center space-x-2">
            <div
              class="w-4 h-4 border border-gray-300"
              :style="{ backgroundColor: getColorForValue(maxValue) }"
            ></div>
            <span class="text-xs text-gray-500">{{ maxValue }}+</span>
          </div>
        </div>

        <div class="text-sm text-gray-600">Click cells for details</div>
      </div>

      <!-- Cell Details Modal/Panel -->
      <div
        v-if="selectedCell"
        class="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg"
      >
        <div class="flex justify-between items-start mb-2">
          <h4 class="font-medium text-gray-900">
            {{ selectedCell.brokerage }} × {{ selectedCell.company }}
          </h4>
          <button
            @click="selectedCell = null"
            class="text-gray-400 hover:text-gray-600"
          >
            ✕
          </button>
        </div>

        <div class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span class="text-gray-600">Total Recommendations:</span>
            <span class="font-medium ml-2">{{ selectedCell.count }}</span>
          </div>
          <div>
            <span class="text-gray-600">Recent Activity:</span>
            <span class="font-medium ml-2">{{ selectedCell.recentDate }}</span>
          </div>
        </div>

        <div class="mt-3">
          <span class="text-gray-600 text-sm">Actions:</span>
          <div class="flex gap-2 mt-1">
            <span class="badge-upgrade"
              >{{ selectedCell.upgrades }} upgrades</span
            >
            <span class="badge-downgrade"
              >{{ selectedCell.downgrades }} downgrades</span
            >
            <span class="badge-reiterate"
              >{{ selectedCell.reiterations }} reiterations</span
            >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useMainStore } from "@/stores";

const store = useMainStore();

// State
const topCompaniesCount = ref(10);
const topBrokeragesCount = ref(5);
const colorScale = ref("blue");
const selectedCell = ref<{
  brokerage: string;
  company: string;
  count: number;
  recentDate: string;
  upgrades: number;
  downgrades: number;
  reiterations: number;
} | null>(null);

// Computed
const heatmapData = computed(() => {
  const data: Record<string, Record<string, any[]>> = {};

  store.recommendations.forEach((rec) => {
    const brokerage = rec.brokerage?.name || "Unknown";
    const ticker = rec.company.ticker;

    if (!data[brokerage]) {
      data[brokerage] = {};
    }
    if (!data[brokerage][ticker]) {
      data[brokerage][ticker] = [];
    }

    data[brokerage][ticker].push(rec);
  });

  return data;
});

const topCompanies = computed(() => {
  const companyCounts: Record<
    string,
    { ticker: string; name: string; count: number }
  > = {};

  store.recommendations.forEach((rec) => {
    const ticker = rec.company.ticker;
    if (!companyCounts[ticker]) {
      companyCounts[ticker] = {
        ticker,
        name: rec.company.name,
        count: 0,
      };
    }
    companyCounts[ticker].count++;
  });

  return Object.values(companyCounts)
    .sort((a: { count: number }, b: { count: number }) => b.count - a.count)
    .slice(0, topCompaniesCount.value);
});

const topBrokerages = computed(() => {
  const brokerageCounts: Record<string, { name: string; count: number }> = {};

  store.recommendations.forEach((rec) => {
    const name = rec.brokerage?.name || "Unknown";
    if (!brokerageCounts[name]) {
      brokerageCounts[name] = { name, count: 0 };
    }
    brokerageCounts[name].count++;
  });

  return Object.values(brokerageCounts)
    .sort((a: { count: number }, b: { count: number }) => b.count - a.count)
    .slice(0, topBrokeragesCount.value);
});

const maxValue = computed(() => {
  let max = 0;
  Object.values(heatmapData.value).forEach((brokerageData: any) => {
    Object.values(brokerageData).forEach((recommendations: any) => {
      max = Math.max(max, recommendations.length);
    });
  });
  return max;
});

// Methods
function getCellValue(brokerage: string, ticker: string): number {
  return heatmapData.value[brokerage]?.[ticker]?.length || 0;
}

function getCellStyle(brokerage: string, ticker: string) {
  const value = getCellValue(brokerage, ticker);
  return {
    backgroundColor: getColorForValue(value),
  };
}

function getCellTextClass(brokerage: string, ticker: string) {
  const value = getCellValue(brokerage, ticker);
  // Use white text for darker backgrounds
  return value > maxValue.value * 0.6 ? "text-white" : "text-gray-900";
}

function getColorForValue(value: number): string {
  if (value === 0) return "#f9fafb"; // gray-50

  const intensity = Math.min(value / maxValue.value, 1);

  const colors = {
    blue: {
      light: [239, 246, 255], // blue-50
      dark: [29, 78, 216], // blue-700
    },
    green: {
      light: [240, 253, 244], // green-50
      dark: [21, 128, 61], // green-700
    },
    red: {
      light: [254, 242, 242], // red-50
      dark: [185, 28, 28], // red-700
    },
    purple: {
      light: [250, 245, 255], // purple-50
      dark: [126, 34, 206], // purple-700
    },
  };

  const colorConfig = colors[colorScale.value as keyof typeof colors];

  const r = Math.round(
    colorConfig.light[0] +
      (colorConfig.dark[0] - colorConfig.light[0]) * intensity
  );
  const g = Math.round(
    colorConfig.light[1] +
      (colorConfig.dark[1] - colorConfig.light[1]) * intensity
  );
  const b = Math.round(
    colorConfig.light[2] +
      (colorConfig.dark[2] - colorConfig.light[2]) * intensity
  );

  return `rgb(${r}, ${g}, ${b})`;
}

function getCellTooltip(brokerage: string, ticker: string): string {
  const count = getCellValue(brokerage, ticker);
  return `${brokerage} → ${ticker}: ${count} recommendations`;
}

function showCellDetails(brokerage: string, ticker: string) {
  const recommendations = heatmapData.value[brokerage]?.[ticker] || [];

  if (recommendations.length === 0) {
    selectedCell.value = null;
    return;
  }

  const company =
    topCompanies.value.find((c: { ticker: string }) => c.ticker === ticker)
      ?.name || ticker;

  const upgrades = recommendations.filter(
    (r) =>
      r.action.toLowerCase().includes("raised") ||
      r.action.toLowerCase().includes("upgrade")
  ).length;

  const downgrades = recommendations.filter(
    (r) =>
      r.action.toLowerCase().includes("lowered") ||
      r.action.toLowerCase().includes("downgrade")
  ).length;

  const reiterations = recommendations.filter((r) =>
    r.action.toLowerCase().includes("reiterat")
  ).length;

  const recentDate = recommendations
    .map((r) => new Date(r.time))
    .sort((a, b) => b.getTime() - a.getTime())[0]
    .toLocaleDateString();

  selectedCell.value = {
    brokerage,
    company,
    count: recommendations.length,
    recentDate,
    upgrades,
    downgrades,
    reiterations,
  };
}
</script>

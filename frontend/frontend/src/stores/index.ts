import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { apiService } from "@/services/api";
import type { Company, Brokerage, Recommendation } from "@/types";

export const useMainStore = defineStore("main", () => {
  //State
  const companies = ref<Company[]>([]);
  const brokerages = ref<Brokerage[]>([]);
  const recommendations = ref<Recommendation[]>([]);
  const loading = ref(false);
  const loadingProgress = ref(0);
  const totalRecommendations = ref(0);
  const error = ref<string | null>(null);
  const dataFetched = ref(false);

  //Computed
  const currentRecommendationCount = computed(
    () => recommendations.value.length
  );
  const loadingPercentage = computed(() =>
    totalRecommendations.value > 0
      ? Math.round(
          (currentRecommendationCount.value / totalRecommendations.value) * 100
        )
      : 0
  );

  const upgrades = computed(
    () =>
      recommendations.value.filter(
        (r) =>
          r.action.toLowerCase().includes("raised") ||
          r.action.toLowerCase().includes("upgrade")
      ).length
  );

  const downgrades = computed(
    () =>
      recommendations.value.filter(
        (r) =>
          r.action.toLowerCase().includes("downgrade") ||
          r.action.toLowerCase().includes("lowered")
      ).length
  );

  const reiterations = computed(
    () =>
      recommendations.value.filter((r) =>
        r.action.toLowerCase().includes("reiterat")
      ).length
  );

  const isInitialLoadComplete = computed(() => loadingPercentage.value >= 50);

  //Actions
  async function fetchAllRecommendations() {
    // Prevent multiple fetches
    if (loading.value || dataFetched.value) {
      console.log('Already loading or data already fetched, skipping...')
      return;
    }

    try {
      console.log('🚀 Starting to fetch recommendations...')
      loading.value = true;
      error.value = null;
      recommendations.value = [];

      const limit = 100;
      let offset = 0;
      let hasMore = true;

      const firstResponse = await apiService.getRecommendations({
        limit,
        offset,
      });

      if (!firstResponse.success) {
        error.value = firstResponse.error || "Failed to fetch recommendations";
        return;
      }

      totalRecommendations.value = firstResponse.meta?.total || firstResponse.data.length;
      recommendations.value = [...firstResponse.data];
      offset += limit;

      console.log(`✅ Loaded ${recommendations.value.length}/${totalRecommendations.value} recommendations`)

      while (hasMore && offset < totalRecommendations.value) {
        const response = await apiService.getRecommendations({ limit, offset });

        if (response.success) {
          recommendations.value = [...recommendations.value, ...response.data];
          offset += limit;
          console.log(`📈 Progress: ${recommendations.value.length}/${totalRecommendations.value}`)

          await new Promise((resolve) => setTimeout(resolve, 100));
        } else {
          hasMore = false;
          error.value = response.error || "Failed to fetch recommendations";
        }
      }
      
      dataFetched.value = true;
      console.log('🎉 All data fetched successfully!')
      
    } catch (err) {
      error.value = "Network error: " + err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchCompanies() {
    try {
      const response = await apiService.getCompanies();
      if (response.success) {
        companies.value = [...response.data];
      }
    } catch (err) {
      console.error("Failed to fetch companies: ", err);
    }
  }

  async function fetchBrokerages() {
    try {
      const response = await apiService.getBrokerages();
      if (response.success) {
        brokerages.value = [...response.data];
      }
    } catch (err) {
      console.error("Failed to fetch brokerages: ", err);
    }
  }
  return {
    //State
    companies,
    brokerages,
    recommendations,
    loading,
    loadingProgress,
    totalRecommendations: computed(() => totalRecommendations.value),
    error,
    dataFetched,

    //Computed
    currentRecommendationCount,
    loadingPercentage,
    upgrades,
    downgrades,
    reiterations,
    isInitialLoadComplete,

    //Actions
    fetchAllRecommendations,
    fetchCompanies,
    fetchBrokerages
  }
});

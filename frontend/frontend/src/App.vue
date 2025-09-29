<script setup lang="ts">
// App.vue is now just the layout - views handle their own logic
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header with Navigation -->
    <header class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <div class="flex items-center">
            <span class="text-2xl">📊</span>
            <span class="ml-2 text-xl font-bold text-gray-900">AnalystHub</span>
          </div>
          
          <!-- Navigation -->
          <nav class="flex space-x-8">
            <router-link
              to="/"
              class="nav-link"
              :class="{ 'nav-link-active': $route.name === 'home' }"
            >
              Home
            </router-link>
            <router-link
              to="/recommendations"
              class="nav-link"
              :class="{ 'nav-link-active': $route.name === 'recommendations' }"
            >
              Recommendations
            </router-link>
            <router-link
              to="/companies"
              class="nav-link"
              :class="{ 'nav-link-active': $route.name === 'companies' }"
            >
              Companies
            </router-link>
            <router-link
              to="/brokerages"
              class="nav-link"
              :class="{ 'nav-link-active': $route.name === 'brokerages' }"
            >
              Brokerages
            </router-link>
          </nav>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main>
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.nav-link {
  @apply text-gray-700 hover:text-blue-600 px-3 py-2 text-sm font-medium transition-colors;
}

.nav-link-active {
  @apply text-blue-600 border-b-2 border-blue-600;
}
</style>

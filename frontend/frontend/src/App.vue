<script setup lang="ts">
// Test component
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-7xl mx-auto px-4 py-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-8">AnalystHub - Tailwind Test</h1>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="card">
          <h2 class="text-lg font-semibold text-gray-900">Total Recommendations</h2>
          <p class="text-3xl font-bold text-gray-600 mt-2">200</p>
        </div>
        
        <div class="card">
          <h2 class="text-lg font-semibold text-gray-900">Upgrades</h2>
          <p class="text-3xl font-bold text-green-600 mt-2">80</p>
        </div>
        
        <div class="card">
          <h2 class="text-lg font-semibold text-gray-900">Downgrades</h2>
          <p class="text-3xl font-bold text-red-600 mt-2">64</p>
        </div>
      </div>
      
      <button class="btn-primary mt-6">Test Button</button>
    </div>
  </div>
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>

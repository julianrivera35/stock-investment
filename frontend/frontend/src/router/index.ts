import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("../views/HomeView.vue"),
    },
    {
      path: "/recommendations",
      name: "recommendations",
      component: () => import("../views/RecommendationsView.vue"),
    },
    {
      path: "/companies",
      name: "companies",
      component: () => import("../views/CompaniesView.vue"),
    },
    {
      path: "/brokerages",
      name: "brokerages",
      component: () => import("../views/BrokeragesView.vue"),
    },
  ],
});

export default router;

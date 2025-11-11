import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import { createRouter, createMemoryHistory } from "vue-router";
import App from "@/App.vue";
import router from "../router/index";

const router = createRouter({
  history: createMemoryHistory(),
  routes: [
    { path: "/", name: "home", component: { template: "<div>Home</div>" } },
    {
      path: "/recommendations",
      name: "recommendations",
      component: { template: "<div>Recommendations</div>" },
    },
    {
      path: "/companies",
      name: "companies",
      component: { template: "<div>Companies</div>" },
    },
    {
      path: "/brokerages",
      name: "brokerages",
      component: { template: "<div>Brokerages</div>" },
    },
  ],
});

describe("App.vue", () => {
  it("renders header with logo and title", () => {
    const wrapper = mount(App, {
      global: {
        plugins: [router],
      },
    });

    expect(wrapper.text()).toContain("AnalystHub");
    expect(wrapper.text()).toContain("📊");
  });

  it("displays navigation links", () => {
    const wrapper = mount(App, {
      global: {
        plugins: [router],
      },
    });
    const links = wrapper.findAll("nav a");
    expect(links).toHaveLength(4);
    expect(links[0].text()).toBe("Home");
    expect(links[1].text()).toBe("Recommendations");
    expect(links[2].text()).toBe("Companies");
    expect(links[3].text()).toBe("Brokerages");
  });

  it("applies active class to current route link", async () => {
    await router.push("/recommendations");
    await router.isReady();

    const wrapper = mount(App, {
      global: {
        plugins: [router],
      },
    });
    const links = wrapper.findAll("nav a");
    expect(links[1].classes()).toContain("nav-link-active");
  });
});

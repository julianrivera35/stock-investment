import { expect } from "./../../../node_modules/@types/chai/index.d";
import { describe, it, expect, beforeEach, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import BrokeragesView from "../BrokeragesView.vue";
import { useMainStore } from "../../stores/index";

describe("BrokeragesView", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("renders the component", () => {
    const wrapper = mount(BrokeragesView);
    expect(wrapper.find("h1").text()).toBe("Brokerages");
  });

  it("displays brokerages from the store", async () => {
    const store = useMainStore();
    store.brokerages = [
      { id: 1, name: "Brokerage A" },
      { id: 2, name: "Brokerage B" },
    ];
    //Creates the component for testing -- mounts it
    const wrapper = mount(BrokeragesView);
    //.vm allows access to the component methods and data
    await wrapper.vm.$nextTick();
    // Function to find all table rows in the tbody
    const rows = wrapper.findAll("tbody tr");
    expect(rows).toHaveLength(2);
    expect(rows[0].text()).toContain("Brokerage A");
    expect(rows[1].text()).toContain("Brokerage B");
  });

  it("filters brokerages based on search term", async () => {
    const store = useMainStore();
    store.brokerages = [
      { id: 1, name: "Brokerage A" },
      { id: 2, name: "Brokerage B" },
      { id: 3, name: "Brokerage C" },
    ];

    const wrapper = mount(BrokeragesView);
    await wrapper.vm.$nextTick();

    // Find the search input
    const searchInput = wrapper.find('input[type="text"]');

    // Type the search term
    await searchInput.setValue("Brokerage A");
    await wrapper.vm.$nextTick();

    const rows = wrapper.findAll("tbody tr");
    expect(rows).toHaveLength(1);
    expect(rows[0].text()).toContain("Brokerage A");
  });

  it('shows "No brokerages found" message when filtered list is empty', async () => {
    const store = useMainStore();
    store.brokerages = [
      { id: 1, name: "Brokerage A" },
      { id: 2, name: "Brokerage B" },
      { id: 3, name: "Brokerage C" },
    ];

    const wrapper = mount(BrokeragesView);
    await wrapper.vm.$nextTick();

    const searchInput = wrapper.find('input[type="text"]');
    await searchInput.setValue("Brokerage D");
    await wrapper.vm.$nextTick();

    expect(wrapper.text()).toContain("No brokerages found.");
  });

  it("calls fetchBrokerages on mount", () => {
    const store = useMainStore();
    //Creates a mock function to track calls
    store.fetchBrokerages = vi.fn();
    mount(BrokeragesView);
    expect(store.fetchBrokerages).toHaveBeenCalledOnce();
  });
});

import { describe, it, expect, beforeEach, vi } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useMainStore } from "../index";

describe("Main Store", () => {
  // Create a fresh Pinia instance before each test
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("initializes with empty arrays", () => {
    const store = useMainStore();
    expect(store.brokerages).toEqual([]);
    expect(store.companies).toEqual([]);
    expect(store.recommendations).toEqual([]);
  });

  it("fetches brokerages successfully", async () => {
    const store = useMainStore();

    global.fetch = vi.fn(() => {
      Promise.resolve({
        ok: true,
        json: () =>
          Promise.resolve([
            { id: 1, name: "Brokerage A" },
            { id: 2, name: "Brokerage B" },
          ]),
      } as Response);
    });

    await store.fetchBrokerages();

    expect(store.brokerages).toHaveLength(2);
    expect(store.brokerages[0].name).toBe("Brokerage A");
  });

  it("handles fetch errors gracefully", async () => {
    const store = useMainStore();

    global.fetch = vi.fn(() => {
      Promise.reject(new Error("Network error"));
    });
    expect(store.brokerages).toEqual([]);
  });
});

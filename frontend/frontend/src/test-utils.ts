import { createPinia, setActivePinia } from "pinia";
import { beforeEach } from "vitest";

export function setupTestPinia() {
  beforeEach(() => {
    setActivePinia(createPinia());
  });
}

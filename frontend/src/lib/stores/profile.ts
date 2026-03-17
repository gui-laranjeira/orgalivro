import { writable } from "svelte/store";
import type { Profile } from "../types";

const STORAGE_KEY = "orgalivro_active_profile";

function createProfileStore() {
  const stored =
    typeof localStorage !== "undefined"
      ? localStorage.getItem(STORAGE_KEY)
      : null;
  const initial: Profile | null = stored ? (JSON.parse(stored) as Profile) : null;

  const { subscribe, set } = writable<Profile | null>(initial);

  return {
    subscribe,
    set(profile: Profile | null) {
      if (profile) {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(profile));
      } else {
        localStorage.removeItem(STORAGE_KEY);
      }
      set(profile);
    },
  };
}

export const activeProfile = createProfileStore();

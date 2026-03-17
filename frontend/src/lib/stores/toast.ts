import { writable } from "svelte/store";

export type ToastType = "success" | "error" | "info";

export interface Toast {
  id: number;
  message: string;
  type: ToastType;
}

let _id = 0;

function createToastStore() {
  const { subscribe, update } = writable<Toast[]>([]);

  function show(message: string, type: ToastType = "info") {
    const id = ++_id;
    update((toasts) => [...toasts, { id, message, type }]);
    setTimeout(() => {
      update((toasts) => toasts.filter((t) => t.id !== id));
    }, 3000);
  }

  return {
    subscribe,
    success: (msg: string) => show(msg, "success"),
    error: (msg: string) => show(msg, "error"),
    info: (msg: string) => show(msg, "info"),
  };
}

export const toast = createToastStore();

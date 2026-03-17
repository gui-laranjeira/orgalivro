export const STATUS_LABELS: Record<string, string> = {
  want_to_read: "Want to Read",
  reading: "Reading",
  read: "Read",
};

export const STATUS_COLORS: Record<string, string> = {
  want_to_read: "bg-gray-100 text-gray-700",
  reading: "bg-blue-100 text-blue-700",
  read: "bg-green-100 text-green-700",
};

export function debounce<T extends (...args: unknown[]) => unknown>(
  fn: T,
  ms: number
): (...args: Parameters<T>) => void {
  let timer: ReturnType<typeof setTimeout>;
  return (...args: Parameters<T>) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), ms);
  };
}

export function starRating(rating: number | null): string {
  if (!rating) return "—";
  return "★".repeat(rating) + "☆".repeat(5 - rating);
}

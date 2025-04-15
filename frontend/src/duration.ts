export function duration(minutes: number) {
  const hours = Math.floor((minutes / 60));
  const m = minutes % 60;
  return `${hours.toString().padStart(2, "0")}:${m.toString().padStart(2, "0")}`;
}
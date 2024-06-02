export function hoursToMilliseconds(hours: number) {
  return 1000 * 60 * 60 * hours;
}

function pad(number: number): string {
  return String(number).padStart(2, "0");
}

export function formatSeconds(totalSeconds: number): string {
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor(totalSeconds % 3600 / 60);
  const seconds = totalSeconds % 60;
  return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
}

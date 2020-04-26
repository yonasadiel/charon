export function durationText(startTime: Date, endTime: Date): string {
  const diff = endTime.getTime() - startTime.getTime();
  const durationDays = Math.floor(diff / (24 * 60 * 60 * 1000));
  const durationHours = Math.floor(diff / (60 * 60 * 1000)) % 24;
  const durationMinutes = Math.floor(diff / (60 * 1000)) % 60;
  const durationSeconds = Math.floor(diff / (1000)) % 60;
  let durationText = '';
  if (durationDays > 0) {
    durationText = `${durationDays} hari ${durationHours} jam ${durationMinutes} menit ${durationSeconds} detik`;
  } else if (durationHours > 0) {
    durationText = `${durationHours} jam ${durationMinutes} menit ${durationSeconds} detik`;
  } else if (durationMinutes > 0) {
    durationText = `${durationMinutes} menit ${durationSeconds} detik`;
  } else if (durationSeconds > 0) {
    durationText = `${durationSeconds} detik`;
  }
  return durationText;
}

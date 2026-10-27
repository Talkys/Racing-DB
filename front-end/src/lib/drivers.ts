import { get } from '../api/client.ts'
import type { Driver } from '../api/types.ts'

/** Fetches each distinct driver once. Unknown ids map to undefined so callers can show the raw id. */
export async function loadDrivers(ids: (string | null | undefined)[]): Promise<Map<string, Driver | undefined>> {
  const unique = [...new Set(ids.filter((id): id is string => !!id))]
  const drivers = await Promise.all(unique.map((id) => get('driver', id)))
  return new Map(unique.map((id, i) => [id, drivers[i] ?? undefined]))
}

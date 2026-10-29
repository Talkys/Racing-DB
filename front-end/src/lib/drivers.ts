import { getMany } from '../api/client.ts'
import type { Driver } from '../api/types.ts'

/** Fetches the given drivers in one request. Unknown ids are absent so callers can show the raw id. */
export const loadDrivers = (ids: (string | null | undefined)[]): Promise<Map<string, Driver>> => getMany('drivers', ids)

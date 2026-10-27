/** Null, undefined and empty values render as an em dash. */
export const dash = (value: string | number | null | undefined) => (value == null || value === '' ? '—' : String(value))

/** `2026-05-24` or `2026-05-24T12:00:00Z` -> `2026-05-24`. */
export const formatDate = (iso: string | null | undefined) => dash(iso?.slice(0, 10))

export const formatDateTime = (iso: string) => new Date(iso).toLocaleString()

/** `HYBRID_TURBO` -> `hybrid turbo`. */
export const formatEnum = (value: string) => value.replaceAll('_', ' ').toLowerCase()

export const driverName = (d: { first_name: string; last_name: string; code: string }) =>
  `${d.first_name} ${d.last_name} (${d.code})`

export const articleHref = (articleId: string, lang?: string) =>
  `article.html?id=${encodeURIComponent(articleId)}${lang ? `&lang=${encodeURIComponent(lang)}` : ''}`

export const raceHref = (raceId: string) => `race.html?id=${encodeURIComponent(raceId)}`

export const carHref = (carId: string) => `car.html?id=${encodeURIComponent(carId)}`

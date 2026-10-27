/**
 * Picks which of an article's existing texts to show. Never invents a language:
 * the result is always one of `available`, or null when there is none.
 *
 * Candidates are `requested`, then each browser language. For each candidate the first hit wins:
 * exact tag, then its primary subtag (`pt-br` -> `pt`), then the first available tag sharing that
 * primary subtag. With no hit, the first available tag is used.
 */
export function pickLanguage(
  available: readonly string[],
  requested?: string | null,
  browserLanguages: readonly string[] = navigator.languages,
): string | null {
  const tags = available.map((t) => t.toLowerCase())
  if (tags.length === 0) return null

  const candidates = [...(requested ? [requested] : []), ...browserLanguages].map((c) => c.toLowerCase())
  for (const candidate of candidates) {
    const primary = candidate.split('-')[0]
    const hit =
      tags.find((t) => t === candidate) ??
      tags.find((t) => t === primary) ??
      tags.find((t) => t.split('-')[0] === primary)
    if (hit) return hit
  }
  return tags[0]
}

import { raceHref } from '../lib/format.ts'
import { dash, initSearchPage, lookup } from '../lib/search-page.ts'

initSearchPage({
  path: 'race',
  async rows(races) {
    const [categories, circuits, drivers] = await Promise.all([lookup('categories'), lookup('circuit'), lookup('driver')])
    return races.map((r) => {
      const winner = r.winner_id ? drivers.get(r.winner_id) : undefined
      return {
        date: r.race_date.slice(0, 10),
        name: r.name,
        href: raceHref(r.id),
        category: categories.get(r.category_id)?.name ?? r.category_id,
        season: String(r.season),
        round: String(r.round),
        circuit: dash(circuits.get(r.circuit_id)?.name),
        winner: winner ? `${winner.first_name} ${winner.last_name}` : '—',
      }
    })
  },
})

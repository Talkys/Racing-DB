import { getMany } from '../api/client.ts'
import { raceHref } from '../lib/format.ts'
import { dash, initSearchPage } from '../lib/search-page.ts'

initSearchPage({
  path: 'races',
  async rows(races) {
    const [categories, circuits, drivers] = await Promise.all([
      getMany('categories', races.map((r) => r.category_id)),
      getMany('circuits', races.map((r) => r.circuit_id)),
      getMany('drivers', races.map((r) => r.winner_id)),
    ])
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
        winner: winner ? `${winner.first_name} ${winner.last_name}` : '-',
      }
    })
  },
})

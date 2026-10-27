import { search } from '../api/client.ts'
import type { Driver } from '../api/types.ts'
import { entityArticles, latestArticles } from '../lib/articles.ts'
import { fill, fillList, section } from '../lib/dom.ts'
import { formatDate, raceHref } from '../lib/format.ts'

const NONE = '—'
const BULLETS = ' ••• '

const fullName = (d: Driver) => `${d.first_name} ${d.last_name}`
const byId = <T extends { id: string }>(rows: T[]) => new Map(rows.map((r) => [r.id, r]))

// The API has no get-by-id or join endpoints, so related names are resolved from the first page of each list.
const [categories, drivers, circuits, races, articles] = await Promise.all([
  search('categories', { sort: 'name,asc' }),
  search('driver', { sort: 'last_name,asc' }),
  search('circuit'),
  search('race', { sort: 'race_date,desc' }),
  latestArticles(3),
])

const categoryName = byId(categories.items)
const driverById = byId(drivers.items)
const circuitName = byId(circuits.items)

const latestRaces = races.items.slice(0, 5).map((race) => {
  const winner = race.winner_id ? driverById.get(race.winner_id) : undefined
  return {
    date: formatDate(race.race_date),
    name: race.name,
    href: raceHref(race.id),
    category: categoryName.get(race.category_id)?.name ?? race.category_id,
    circuit: circuitName.get(race.circuit_id)?.name ?? NONE,
    winner: winner ? fullName(winner) : NONE,
  }
})

// Ticker
fill(document.body, {
  ticker:
    'LATEST UPDATES: ' +
    latestRaces
      .slice(0, 4)
      .map((r) => `[${r.date}] ${r.name}${r.winner === NONE ? '' : ` won by ${r.winner}`}`)
      .join(BULLETS),
})

// Browse vaults
fillList(
  section('categories'),
  categories.items.map((c) => ({ name: c.name, href: `races.html?col=category_id&op=eq&q=${encodeURIComponent(c.id)}` })),
)

// Latest races
fillList(section('latest-races'), latestRaces)

// Featured driver: rotates daily, with the newest article about them as teaser.
if (drivers.items.length > 0) {
  const dayOfYear = Math.floor((Date.now() - new Date(new Date().getFullYear(), 0, 0).getTime()) / 864e5)
  const featured = drivers.items[dayOfYear % drivers.items.length]
  const [bio] = await entityArticles('driver', featured.id)
  const text = bio?.text
  fill(section('featured'), {
    name: `${featured.last_name}, ${featured.first_name}`,
    code: featured.code,
    nationality: featured.nationality,
    photo: `[ PHOTO: ${fullName(featured)} ]`,
    teaser: text?.summary ?? text?.title,
    href: bio?.fields.href,
  })
}

// Latest articles
fillList(section('articles'), articles.map((a) => a.fields))

// Driver index
fillList(
  section('drivers'),
  drivers.items.slice(0, 5).map((d) => ({ name: fullName(d), code: d.code, href: `drivers.html?col=code&op=eq&q=${encodeURIComponent(d.code)}` })),
)

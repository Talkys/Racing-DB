import { getMany, search } from '../api/client.ts'
import type { Driver } from '../api/types.ts'
import { entityArticles, latestArticles } from '../lib/articles.ts'
import { fill, fillList, section, showError } from '../lib/dom.ts'
import { entityHref, formatDate, hrefOf, raceHref } from '../lib/format.ts'

const NONE = '-'
const BULLETS = ' ••• '

const fullName = (d: Driver) => `${d.first_name} ${d.last_name}`

async function load(): Promise<void> {
  const [categories, drivers, races, articles] = await Promise.all([
    // Inactive categories are listed last, with a badge.
    search('categories', { sort: ['is_active,desc', 'name'] }),
    search('drivers', { sort: 'last_name,asc' }),
    search('races', { sort: 'race_date,desc' }),
    latestArticles(3),
  ])

  const recent = races.items.slice(0, 5)
  const [raceCategories, circuits, winners] = await Promise.all([
    getMany('categories', recent.map((r) => r.category_id)),
    getMany('circuits', recent.map((r) => r.circuit_id)),
    getMany('drivers', recent.map((r) => r.winner_id)),
  ])

  const latestRaces = recent.map((race) => {
    const winner = race.winner_id ? winners.get(race.winner_id) : undefined
    return {
      date: formatDate(race.race_date),
      name: race.name,
      href: raceHref(race.id),
      category: raceCategories.get(race.category_id)?.name ?? race.category_id,
      categoryHref: hrefOf('category', race.category_id, raceCategories.has(race.category_id)),
      circuit: circuits.get(race.circuit_id)?.name ?? NONE,
      circuitHref: hrefOf('circuit', race.circuit_id, circuits.has(race.circuit_id)),
      winner: winner ? fullName(winner) : NONE,
      winnerHref: hrefOf('driver', race.winner_id, !!winner),
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
    categories.items.map((c) => ({ name: c.name, href: entityHref('category', c.id), inactive: c.is_active ? '' : '1' })),
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
      driverHref: entityHref('driver', featured.id),
      code: featured.code,
      nationality: featured.nationality,
      photo: `[ PHOTO: ${fullName(featured)} ]`,
      teaser: text?.summary ?? text?.title,
      bioHref: bio?.fields.href,
    })
  }

  // Latest articles
  fillList(section('articles'), articles.map((a) => a.fields))

  // Driver index
  fillList(
    section('drivers'),
    drivers.items.slice(0, 5).map((d) => ({ name: fullName(d), code: d.code, href: entityHref('driver', d.id) })),
  )
}

try {
  await load()
} catch (err) {
  fill(document.body, { ticker: 'LATEST UPDATES: unavailable' })
  showError(err)
}

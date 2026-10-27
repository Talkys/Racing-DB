import { get } from '../api/client.ts'
import { entityArticles } from '../lib/articles.ts'
import { fill, fillList, section, showError } from '../lib/dom.ts'
import { dash, formatEnum } from '../lib/format.ts'

// car.html?id=<car id>

async function load(id: string): Promise<void> {
  const car = await get('car', id)
  if (!car) throw new Error(`Car with ID '${id}' does not exist.`)

  const [category, team, engine, chassis, articles] = await Promise.all([
    get('categories', car.category_id),
    get('team', car.team_id),
    get('engine', car.engine_id),
    get('chassis', car.chassis_id),
    entityArticles('car', car.id),
  ])
  const { parts } = car

  document.title = `${car.name} (${car.season}) - RaceDB.org`
  fill(document.body, {
    loaded: '1',
    name: `${car.name} (${car.season})`,
    season: String(car.season),
    category: category?.name ?? car.category_id,
    categoryHref: `cars.html?col=category_id&op=eq&q=${encodeURIComponent(car.category_id)}`,

    teamName: team?.name ?? car.team_id,
    teamCountry: dash(team?.country),
    teamFounded: dash(team?.founded_year),

    engineName: engine ? `${engine.manufacturer} ${engine.name}` : car.engine_id,
    engineConfiguration: dash(engine?.configuration),
    engineDisplacement: dash(engine?.displacement_cc),
    engineCylinders: dash(engine?.cylinders),
    engineAspiration: engine ? formatEnum(engine.aspiration) : '—',
    enginePower: dash(engine?.power_kw),

    chassisName: chassis ? `${chassis.manufacturer} ${chassis.name}` : car.chassis_id,
    chassisMaterial: dash(chassis?.material),
    chassisWheelbase: dash(chassis?.wheelbase_mm),
    chassisWeight: dash(chassis?.weight_kg),

    gearboxType: formatEnum(parts.gearbox.type),
    gearboxSpeeds: String(parts.gearbox.speeds),
    brakesFront: parts.brakes.front,
    brakesRear: parts.brakes.rear,
    suspensionFront: parts.suspension.front,
    suspensionRear: parts.suspension.rear,
    tyreSupplier: parts.tyre_supplier,
  })
  fillList(
    section('parts'),
    Object.entries(parts.other).map(([key, value]) => ({ key: formatEnum(key), value })),
  )
  fillList(section('articles'), articles.map((a) => a.fields))
  section('articles').querySelector<HTMLElement>('[data-empty]')!.hidden = articles.length > 0
}

const id = new URLSearchParams(location.search).get('id')
try {
  if (!id) throw new Error('No car id given. Use car.html?id=<car id>.')
  await load(id)
} catch (err) {
  showError(err)
} finally {
  fill(document.body, { loading: '' })
}

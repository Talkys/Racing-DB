import { get } from '../api/client.ts'
import { initDetailPage } from '../lib/detail-page.ts'
import { fill, fillList, section } from '../lib/dom.ts'
import { dash, formatEnum, hrefOf } from '../lib/format.ts'
import { carPartsFields } from '../lib/payloads.ts'

// car.html?id=<car id>

initDetailPage({
  path: 'cars',
  noun: 'car',
  articles: 'car',
  async render(car) {
    const [category, team, engine, chassis] = await Promise.all([
      get('categories', car.category_id),
      get('teams', car.team_id),
      get('engines', car.engine_id),
      get('chassis', car.chassis_id),
    ])

    const parts = carPartsFields(car)

    document.title = `${car.name} (${car.season}) - RaceDB.org`
    fill(document.body, {
      name: `${car.name} (${car.season})`,
      season: String(car.season),
      category: category?.name ?? car.category_id,
      categoryHref: hrefOf('category', car.category_id, !!category),

      teamName: team?.name ?? car.team_id,
      teamHref: hrefOf('team', car.team_id, !!team),
      teamCountry: dash(team?.country),
      teamFounded: dash(team?.founded_year),

      engineName: engine ? `${engine.manufacturer} ${engine.name}` : car.engine_id,
      engineHref: hrefOf('engine', car.engine_id, !!engine),
      engineConfiguration: dash(engine?.configuration),
      engineDisplacement: dash(engine?.displacement_cc),
      engineCylinders: dash(engine?.cylinders),
      engineAspiration: engine ? formatEnum(engine.aspiration) : '-',
      enginePower: dash(engine?.power_kw),

      chassisName: chassis ? `${chassis.manufacturer} ${chassis.name}` : car.chassis_id,
      chassisHref: hrefOf('chassis', car.chassis_id, !!chassis),
      chassisMaterial: dash(chassis?.material),
      chassisWheelbase: dash(chassis?.wheelbase_mm),
      chassisWeight: dash(chassis?.weight_kg),

      ...parts.fields,
    })
    fillList(section('parts'), parts.other)
  },
})

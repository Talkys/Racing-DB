import { fillRelated, initDetailPage, listHref, relatedCars } from '../lib/detail-page.ts'
import { fill } from '../lib/dom.ts'
import { dash, formatEnum } from '../lib/format.ts'

// engine.html?id=<engine id>

initDetailPage({
  path: 'engines',
  noun: 'engine',
  articles: 'engine',
  async render(engine) {
    const name = `${engine.manufacturer} ${engine.name}`
    document.title = `${name} - RaceDB.org`
    fill(document.body, {
      name,
      manufacturer: engine.manufacturer,
      configuration: engine.configuration,
      displacement: String(engine.displacement_cc),
      cylinders: String(engine.cylinders),
      aspiration: formatEnum(engine.aspiration),
      power: dash(engine.power_kw),
    })
    fillRelated('cars', await relatedCars(`engine_id,eq,${engine.id}`), listHref('cars.html', 'engine_id', engine.id))
  },
})

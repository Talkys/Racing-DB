import { fillRelated, initDetailPage, listHref, relatedCars } from '../lib/detail-page.ts'
import { fill } from '../lib/dom.ts'
import { dash } from '../lib/format.ts'

// chassis-detail.html?id=<chassis id> (chassis.html is the list)

initDetailPage({
  path: 'chassis',
  noun: 'chassis',
  articles: 'chassis',
  async render(chassis) {
    const name = `${chassis.manufacturer} ${chassis.name}`
    document.title = `${name} - RaceDB.org`
    fill(document.body, {
      name,
      manufacturer: chassis.manufacturer,
      material: chassis.material,
      wheelbase: dash(chassis.wheelbase_mm),
      weight: dash(chassis.weight_kg),
    })
    fillRelated('cars', await relatedCars(`chassis_id,eq,${chassis.id}`), listHref('cars.html', 'chassis_id', chassis.id))
  },
})

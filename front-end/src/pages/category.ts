import { fillRelated, initDetailPage, listHref, relatedCars, relatedRaces } from '../lib/detail-page.ts'
import { fill } from '../lib/dom.ts'

// category.html?id=<category slug>
// Inactive categories are shown like the others, with an "inactive" badge.

initDetailPage({
  path: 'categories',
  noun: 'category',
  articles: 'category',
  async render(category) {
    document.title = `${category.name} - RaceDB.org`
    fill(document.body, {
      name: category.name,
      description: category.description,
      status: category.is_active ? 'Active' : 'Inactive',
      inactive: category.is_active ? '' : '1',
    })
    const [races, cars] = await Promise.all([relatedRaces(`category_id,eq,${category.id}`), relatedCars(`category_id,eq,${category.id}`)])
    fillRelated('races', races, listHref('races.html', 'category_id', category.id))
    fillRelated('cars', cars, listHref('cars.html', 'category_id', category.id))
  },
})

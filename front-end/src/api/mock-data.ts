import type { EntityMap, EntityPath } from './types.ts'

// Stand-in data in the Go back end's JSON shape. Categories come from ../init.sql; everything else is the
// spec's seed (spec5.md §8.2), which covers every case the UI must handle: the three race payload types,
// a race-less winner, nullable entity fields, multi-language articles, a draft and an article with no English text.
// Additions beyond the spec: a driver article (a5) for the home page's featured box, a future-dated article (a6)
// that must stay hidden, and enough extra drivers to need a second page (25 rows per page).
// ../../../seed.sql is generated from this file, so both stay identical.

const created = '2026-01-01T00:00:00Z'

const DRIVER = {
  ver: 'b11a43a8-4e36-47a8-8e68-301fc1cae6ef',
  bag: 'a22b54b9-5f47-48b9-9f79-412fd2dbf7f0',
  har: 'c33c65ca-6058-49ca-a080-5230e3ec0801',
  bac: 'd44d76db-7169-4adb-b191-6341f4fd1912',
}
const CIRCUIT = {
  monaco: '23fa3b45-df13-4f9e-a864-44aa4cf81190',
  mugello: '6fa85f64-5717-4562-b3fc-2c963f66afa6',
  sarthe: 'b39a4896-e26b-4e01-9a74-d4b94639e24f',
}
const RACE = {
  monaco: 'c39a7593-36c5-4ceb-8d07-6bc82f808381',
  mugello: 'd40b8604-47d6-5dfc-9e18-7cd93f919492',
  lemans: 'e51c9715-58e7-6e0d-af29-8de04fa205a3',
}
const seq = (prefix: number, n: number) => `${prefix}000000-0000-4000-8000-${String(n).padStart(12, '0')}`
const TEAM = [1, 2, 3].map((n) => seq(71, n))
const ENGINE = [1, 2, 3].map((n) => seq(72, n))
const CHASSIS = [1, 2, 3].map((n) => seq(73, n))
const CAR = [1, 2, 3].map((n) => seq(74, n))
const ARTICLE = [1, 2, 3, 4, 5, 6].map((n) => seq(75, n))

const MORE_DRIVERS = [
  ['Lewis', 'Hamilton', 'HAM', 'British'], ['Charles', 'Leclerc', 'LEC', 'Monegasque'], ['Lando', 'Norris', 'NOR', 'British'],
  ['Oscar', 'Piastri', 'PIA', 'Australian'], ['George', 'Russell', 'RUS', 'British'], ['Kimi', 'Antonelli', 'ANT', 'Italian'],
  ['Fernando', 'Alonso', 'ALO', 'Spanish'], ['Lance', 'Stroll', 'STR', 'Canadian'], ['Pierre', 'Gasly', 'GAS', 'French'],
  ['Esteban', 'Ocon', 'OCO', 'French'], ['Alexander', 'Albon', 'ALB', 'Thai'], ['Carlos', 'Sainz', 'SAI', 'Spanish'],
  ['Nico', 'Hulkenberg', 'HUL', 'German'], ['Gabriel', 'Bortoleto', 'BOR', 'Brazilian'], ['Yuki', 'Tsunoda', 'TSU', 'Japanese'],
  ['Liam', 'Lawson', 'LAW', 'New Zealander'], ['Isack', 'Hadjar', 'HAD', 'French'], ['Oliver', 'Bearman', 'BEA', 'British'],
  ['Franco', 'Colapinto', 'COL', 'Argentine'], ['Marc', 'Marquez', 'MAR', 'Spanish'], ['Jorge', 'Martin', 'MRT', 'Spanish'],
  ['Pedro', 'Acosta', 'ACO', 'Spanish'], ['Enea', 'Bastianini', 'BAS', 'Italian'], ['Fabio', 'Quartararo', 'QUA', 'French'],
  ['Kamui', 'Kobayashi', 'KOB', 'Japanese'], ['Sebastien', 'Buemi', 'BUE', 'Swiss'],
]

const text = (n: number, article: number, language: string, title: string, summary: string | null, body_md: string, at: string) => ({
  id: seq(76, n), article_id: ARTICLE[article - 1], language, title, summary, body_md, created_at: at, updated_at: at,
})

export const mockData: { [P in EntityPath]: EntityMap[P][] } = {
  categories: [
    { id: 'f1', name: 'Formula 1', description: 'The pinnacle of open-wheel single-seater motorsport, featuring hybrid power units and world-renowned constructors.', is_active: true },
    { id: 'motogp', name: 'MotoGP', description: 'The premier class of motorcycle road racing, showcasing high-performance 1000cc prototype bikes.', is_active: true },
    { id: 'wec', name: 'FIA World Endurance Championship', description: 'Multi-class sports car racing featuring Hypercars and GT cars battling in endurance events including the 24 Hours of Le Mans.', is_active: true },
    { id: 'indycar', name: 'IndyCar Series', description: 'North America premier open-wheel racing series featuring oval, road, and street circuit competition.', is_active: true },
    { id: 'wrc', name: 'World Rally Championship', description: 'The ultimate off-road competition, challenging drivers and co-drivers across gravel, snow, and asphalt stages worldwide.', is_active: true },
    { id: 'formula-e', name: 'Formula E', description: 'All-electric single-seater championship racing through street circuits in major global cities.', is_active: true },
    { id: 'nascar-cup', name: 'NASCAR Cup Series', description: 'The top racing series of stock car racing in North America, primarily contested on oval tracks.', is_active: true },
    { id: 'imsa', name: 'IMSA WeatherTech SportsCar Championship', description: 'Premier North American endurance and sports car racing series featuring GTP, LMP2, and GTD classes.', is_active: true },
    { id: 'f2', name: 'Formula 2', description: 'Official feeder series for Formula 1 featuring identical chassis and engine specifications to test driver skill.', is_active: true },
    { id: 'supercars', name: 'Repco Supercars Championship', description: 'Australia premier touring car racing series, featuring high-powered V8 race cars.', is_active: true },
    { id: 'superbike', name: 'World Superbike (WorldSBK)', description: 'Production-based motorcycle road racing championship featuring modified street bikes.', is_active: true },
    { id: 'dtm', name: 'Deutsche Tourenwagen Masters', description: 'Historic German touring car series operating under GT3 specification rules.', is_active: true },
  ],
  circuits: [
    { id: CIRCUIT.monaco, name: 'Circuit de Monaco', location: 'Monte Carlo', country: 'Monaco', length_km: 3.337, created_at: created },
    { id: CIRCUIT.mugello, name: 'Autodromo del Mugello', location: 'Scarperia e San Piero', country: 'Italy', length_km: 5.245, created_at: created },
    { id: CIRCUIT.sarthe, name: 'Circuit de la Sarthe', location: 'Le Mans', country: 'France', length_km: 13.626, created_at: created },
  ],
  drivers: [
    { id: DRIVER.ver, first_name: 'Max', last_name: 'Verstappen', code: 'VER', nationality: 'Dutch', created_at: created },
    { id: DRIVER.bag, first_name: 'Francesco', last_name: 'Bagnaia', code: 'BAG', nationality: 'Italian', created_at: created },
    { id: DRIVER.har, first_name: 'Brendon', last_name: 'Hartley', code: 'HAR', nationality: 'New Zealander', created_at: created },
    { id: DRIVER.bac, first_name: 'Klaus', last_name: 'Bachler', code: 'BAC', nationality: 'Austrian', created_at: created },
    ...MORE_DRIVERS.map(([first_name, last_name, code, nationality], i) => ({ id: seq(70, i + 1), first_name, last_name, code, nationality, created_at: created })),
  ],
  races: [
    {
      id: RACE.monaco, category_id: 'f1', circuit_id: CIRCUIT.monaco, winner_id: DRIVER.ver,
      name: 'Monaco Grand Prix', season: 2026, round: 8, race_date: '2026-05-24T00:00:00Z',
      payload: {
        fastest_lap: { driver_id: DRIVER.ver, lap_number: 52, lap_time: '1:14.123' },
        tyre_strategies: [{ driver_id: DRIVER.ver, stops: [{ lap: 30, compound: 'HARD' }] }],
        safety_car_deployments: 1,
        track_temperature_celsius: 42.1,
      },
      created_at: created,
    },
    {
      id: RACE.mugello, category_id: 'motogp', circuit_id: CIRCUIT.mugello, winner_id: DRIVER.bag,
      name: "Gran Premio d'Italia", season: 2026, round: 7, race_date: '2026-06-02T00:00:00Z',
      payload: {
        rider_of_the_day_id: DRIVER.bag,
        tyre_spec: { front: 'HARD', rear: 'SOFT' },
        holeshot_winner_id: DRIVER.bag,
        weather_conditions: 'DRY',
      },
      created_at: created,
    },
    {
      id: RACE.lemans, category_id: 'wec', circuit_id: CIRCUIT.sarthe, winner_id: DRIVER.har,
      name: '24 Hours of Le Mans', season: 2026, round: 4, race_date: '2026-06-14T00:00:00Z',
      payload: {
        total_laps_completed: 368,
        class_winners: { hypercar: DRIVER.har, lmgt3: DRIVER.bac },
        full_course_yellow_count: 4,
        red_flag_duration_minutes: 0,
      },
      created_at: created,
    },
  ],
  teams: [
    { id: TEAM[0], name: 'Oracle Red Bull Racing', short_name: 'Red Bull', country: 'United Kingdom', founded_year: 2005, created_at: created },
    { id: TEAM[1], name: 'Ducati Lenovo Team', short_name: 'Ducati', country: 'Italy', founded_year: null, created_at: created },
    { id: TEAM[2], name: 'Toyota Gazoo Racing', short_name: 'Toyota', country: 'Japan', founded_year: 2007, created_at: created },
  ],
  engines: [
    { id: ENGINE[0], manufacturer: 'Red Bull Ford', name: 'DM01', configuration: 'V6', displacement_cc: 1600, cylinders: 6, aspiration: 'HYBRID_TURBO', power_kw: 750.0, created_at: created },
    { id: ENGINE[1], manufacturer: 'Ducati', name: 'Desmosedici GP', configuration: 'V4', displacement_cc: 1000, cylinders: 4, aspiration: 'NATURALLY_ASPIRATED', power_kw: null, created_at: created },
    { id: ENGINE[2], manufacturer: 'Toyota', name: 'V35A', configuration: 'V6', displacement_cc: 3400, cylinders: 6, aspiration: 'TURBOCHARGED', power_kw: 520.0, created_at: created },
  ],
  chassis: [
    { id: CHASSIS[0], manufacturer: 'Red Bull Racing', name: 'RB22 monocoque', material: 'Carbon fibre monocoque', wheelbase_mm: 3600, weight_kg: 768.0, created_at: created },
    { id: CHASSIS[1], manufacturer: 'Ducati', name: 'Desmosedici frame', material: 'Aluminium', wheelbase_mm: null, weight_kg: null, created_at: created },
    { id: CHASSIS[2], manufacturer: 'Toyota', name: 'GR010 monocoque', material: 'Carbon fibre monocoque', wheelbase_mm: null, weight_kg: 1030.0, created_at: created },
  ],
  cars: [
    {
      id: CAR[0], category_id: 'f1', team_id: TEAM[0], engine_id: ENGINE[0], chassis_id: CHASSIS[0], name: 'RB22', season: 2026,
      parts: {
        gearbox: { type: 'SEQUENTIAL', speeds: 8 },
        brakes: { front: 'Carbon-carbon discs', rear: 'Carbon-carbon discs' },
        suspension: { front: 'Pushrod', rear: 'Pullrod' },
        tyre_supplier: 'Pirelli',
        other: { fuel_supplier: 'Petronas' },
      },
      created_at: created,
    },
    {
      id: CAR[1], category_id: 'motogp', team_id: TEAM[1], engine_id: ENGINE[1], chassis_id: CHASSIS[1], name: 'Desmosedici GP26', season: 2026,
      parts: {
        gearbox: { type: 'SEAMLESS', speeds: 6 },
        brakes: { front: 'Carbon discs', rear: 'Steel disc' },
        suspension: { front: 'Telescopic fork', rear: 'Monoshock swingarm' },
        tyre_supplier: 'Michelin',
        other: {},
      },
      created_at: created,
    },
    {
      id: CAR[2], category_id: 'wec', team_id: TEAM[2], engine_id: ENGINE[2], chassis_id: CHASSIS[2], name: 'GR010 Hybrid', season: 2026,
      parts: {
        gearbox: { type: 'SEQUENTIAL', speeds: 7 },
        brakes: { front: 'Carbon discs', rear: 'Carbon discs' },
        suspension: { front: 'Double wishbone', rear: 'Double wishbone' },
        tyre_supplier: 'Michelin',
        other: { hybrid_system: 'Front-axle MGU' },
      },
      created_at: created,
    },
  ],
  articles: [
    { id: ARTICLE[0], entity_type: 'race', entity_id: RACE.monaco, kind: 'report', published_at: '2026-05-25T12:00:00Z', created_at: created },
    { id: ARTICLE[1], entity_type: 'category', entity_id: 'f1', kind: 'overview', published_at: '2026-01-02T00:00:00Z', created_at: created },
    { id: ARTICLE[2], entity_type: 'car', entity_id: CAR[0], kind: 'technical', published_at: null, created_at: created },
    { id: ARTICLE[3], entity_type: 'team', entity_id: TEAM[0], kind: 'history', published_at: '2026-02-01T00:00:00Z', created_at: created },
    { id: ARTICLE[4], entity_type: 'driver', entity_id: DRIVER.ver, kind: 'overview', published_at: '2026-03-01T00:00:00Z', created_at: created },
    { id: ARTICLE[5], entity_type: 'driver', entity_id: DRIVER.bag, kind: 'preview', published_at: '2099-01-01T00:00:00Z', created_at: created },
  ],
  'article-texts': [
    text(1, 1, 'en', 'Verstappen masters Monaco', 'A controlled drive from pole.',
      '# Verstappen masters Monaco\n\nA **controlled** drive from pole, with a single pit stop on lap 30.\n\n* Fastest lap: 1:14.123\n* Safety cars: 1\n',
      '2026-05-25T12:00:00Z'),
    text(2, 1, 'pt-br', 'Verstappen domina em Mônaco', null,
      '## Resumo da corrida\n\nVitória sem sustos em Monte Carlo. Este texto é independente da versão em inglês.\n',
      '2026-05-25T12:00:00Z'),
    text(3, 2, 'en', 'Formula 1 at a glance', 'What the championship is.',
      '# Formula 1\n\nThe FIA Formula One World Championship. See the [FIA site](https://www.fia.com).\n',
      '2026-01-02T00:00:00Z'),
    text(4, 3, 'en', 'RB22 technical notes (draft)', null, 'Draft. Must never be served.', '2026-01-03T00:00:00Z'),
    text(5, 4, 'pt', 'A história da equipe', 'Uma breve história.', '# História\n\nTexto disponível somente em português.\n', '2026-02-01T00:00:00Z'),
    text(6, 5, 'en', 'Max Verstappen', 'Four-time world champion and Red Bull team leader.',
      '# Max Verstappen\n\nDutch driver, champion from 2021 to 2024.\n', '2026-03-01T00:00:00Z'),
    text(7, 6, 'en', 'Bagnaia season preview (scheduled)', null, 'Future-dated. Must never be served.', '2026-03-01T00:00:00Z'),
  ],
}

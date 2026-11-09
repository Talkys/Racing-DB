import type { CarParts, ComputedFields, EntityMap, EntityPath, RacePayload } from './types.ts'

// Stand-in data in the Go back end's JSON shape. Categories come from ../init.sql; everything else is the
// spec's seed (spec5.md §8.2), which covers every case the UI must handle: the three race payload types,
// a race-less winner, nullable entity fields, multi-language articles, a draft and an article with no English text.
// Additions beyond the spec: a driver article (a5) for the home page's featured box, a future-dated article (a6)
// that must stay hidden, enough extra drivers, races and cars to need a second page (25 rows per page), and one race
// and one car whose JSONB breaks its contract (RACE_BROKEN, CAR_BROKEN).
// The extra races are a 2025 season, so the spec's 2026 rows stay on top of the newest-first race list.
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
const TEAM = [1, 2, 3, 4, 5, 6, 7].map((n) => seq(71, n))
const ENGINE = [1, 2, 3, 4, 5, 6, 7].map((n) => seq(72, n))
const CHASSIS = [1, 2, 3, 4, 5, 6, 7].map((n) => seq(73, n))
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

const DRIVER_BY_CODE: Record<string, string> = {
  VER: DRIVER.ver, BAG: DRIVER.bag, HAR: DRIVER.har, BAC: DRIVER.bac,
  ...Object.fromEntries(MORE_DRIVERS.map(([, , code], i) => [code, seq(70, i + 1)])),
}

// Extra circuits for the 2025 season below: key, name, location, country, length_km.
const MORE_CIRCUITS: [string, string, string, string, number][] = [
  ['albertpark', 'Albert Park Circuit', 'Melbourne', 'Australia', 5.278], ['suzuka', 'Suzuka International Racing Course', 'Suzuka', 'Japan', 5.807],
  ['bahrain', 'Bahrain International Circuit', 'Sakhir', 'Bahrain', 5.412], ['redbullring', 'Red Bull Ring', 'Spielberg', 'Austria', 4.318],
  ['silverstone', 'Silverstone Circuit', 'Silverstone', 'United Kingdom', 5.891], ['spa', 'Circuit de Spa-Francorchamps', 'Stavelot', 'Belgium', 7.004],
  ['monza', 'Autodromo Nazionale Monza', 'Monza', 'Italy', 5.793], ['interlagos', 'Autódromo José Carlos Pace', 'São Paulo', 'Brazil', 4.309],
  ['losail', 'Lusail International Circuit', 'Lusail', 'Qatar', 5.38], ['jerez', 'Circuito de Jerez - Ángel Nieto', 'Jerez de la Frontera', 'Spain', 4.423],
  ['assen', 'TT Circuit Assen', 'Assen', 'Netherlands', 4.542], ['motegi', 'Mobility Resort Motegi', 'Motegi', 'Japan', 4.801],
  ['phillipisland', 'Phillip Island Grand Prix Circuit', 'Phillip Island', 'Australia', 4.448], ['fuji', 'Fuji Speedway', 'Oyama', 'Japan', 4.563],
  ['imola', 'Autodromo Enzo e Dino Ferrari', 'Imola', 'Italy', 4.909],
]
const CIRCUIT_BY_KEY: Record<string, string> = {
  monaco: CIRCUIT.monaco, mugello: CIRCUIT.mugello, sarthe: CIRCUIT.sarthe,
  ...Object.fromEntries(MORE_CIRCUITS.map(([key], i) => [key, seq(78, i + 1)])),
}

// 2025 season. Every payload follows its category's contract (types.ts), so the race page renders it.
// F1: round, name, circuit, date, winner, fastest lap, safety cars, track temperature.
const F1_2025: [number, string, string, string, string, string, number, number][] = [
  [1, 'Australian Grand Prix', 'albertpark', '2025-03-16', 'NOR', '1:22.167', 3, 19.8], [3, 'Japanese Grand Prix', 'suzuka', '2025-04-06', 'VER', '1:30.965', 0, 33.4],
  [4, 'Bahrain Grand Prix', 'bahrain', '2025-04-13', 'PIA', '1:35.140', 1, 31.2], [8, 'Monaco Grand Prix', 'monaco', '2025-05-25', 'NOR', '1:13.221', 0, 44.5],
  [11, 'Austrian Grand Prix', 'redbullring', '2025-06-29', 'NOR', '1:07.924', 1, 49.0], [12, 'British Grand Prix', 'silverstone', '2025-07-06', 'NOR', '1:29.337', 3, 21.6],
  [13, 'Belgian Grand Prix', 'spa', '2025-07-27', 'PIA', '1:44.861', 1, 24.1], [16, 'Italian Grand Prix', 'monza', '2025-09-07', 'VER', '1:20.901', 0, 41.3],
  [21, 'São Paulo Grand Prix', 'interlagos', '2025-11-09', 'NOR', '1:12.400', 2, 46.7], [23, 'Qatar Grand Prix', 'losail', '2025-11-30', 'VER', '1:22.996', 1, 27.5],
]
// MotoGP: round, name, circuit, date, winner, front tyre, rear tyre, weather.
const MOTOGP_2025: [number, string, string, string, string, 'SOFT' | 'MEDIUM' | 'HARD', 'SOFT' | 'MEDIUM' | 'HARD', 'DRY' | 'WET'][] = [
  [4, 'Grand Prix of Qatar', 'losail', '2025-04-13', 'MAR', 'MEDIUM', 'MEDIUM', 'DRY'], [5, 'Gran Premio de España', 'jerez', '2025-04-27', 'QUA', 'HARD', 'MEDIUM', 'DRY'],
  [6, 'British Grand Prix', 'silverstone', '2025-05-25', 'BAG', 'MEDIUM', 'SOFT', 'WET'], [9, "Gran Premio d'Italia", 'mugello', '2025-06-22', 'MAR', 'HARD', 'SOFT', 'DRY'],
  [10, 'Dutch TT', 'assen', '2025-06-29', 'MAR', 'MEDIUM', 'MEDIUM', 'DRY'], [13, 'Austrian Grand Prix', 'redbullring', '2025-08-17', 'ACO', 'HARD', 'MEDIUM', 'DRY'],
  [17, 'Grand Prix of Japan', 'motegi', '2025-09-28', 'BAG', 'MEDIUM', 'SOFT', 'DRY'], [19, 'Australian Grand Prix', 'phillipisland', '2025-10-19', 'BAS', 'SOFT', 'MEDIUM', 'WET'],
]
// WEC: round, name, circuit, date, hypercar winner, LMGT3 winner, laps, full course yellows.
const WEC_2025: [number, string, string, string, string, string, number, number][] = [
  [1, 'Qatar 1812 km', 'losail', '2025-02-28', 'KOB', 'BAC', 335, 5], [3, '6 Hours of Spa-Francorchamps', 'spa', '2025-05-10', 'BUE', 'BAC', 146, 3],
  [4, '24 Hours of Le Mans', 'sarthe', '2025-06-15', 'HAR', 'BAC', 387, 6], [5, '6 Hours of São Paulo', 'interlagos', '2025-07-13', 'KOB', 'BAC', 239, 2],
  [7, '6 Hours of Fuji', 'fuji', '2025-09-28', 'HAR', 'BAC', 212, 4], [8, '8 Hours of Bahrain', 'bahrain', '2025-11-08', 'BUE', 'BAC', 254, 1],
]

const race2025 = (n: number, category_id: string, round: number, name: string, circuit: string, date: string, winner: string, payload: RacePayload) => ({
  id: seq(77, n), category_id, circuit_id: CIRCUIT_BY_KEY[circuit], winner_id: DRIVER_BY_CODE[winner],
  name, season: 2025, round, race_date: `${date}T00:00:00Z`, payload, created_at: created,
})

const MORE_RACES = [
  ...F1_2025.map(([round, name, circuit, date, winner, lap_time, sc, temp]) => ({ category_id: 'f1', round, name, circuit, date, winner, payload: {
    fastest_lap: { driver_id: DRIVER_BY_CODE[winner], lap_number: 40 + round, lap_time },
    tyre_strategies: [{ driver_id: DRIVER_BY_CODE[winner], stops: [{ lap: 20 + round, compound: 'HARD' }] }],
    safety_car_deployments: sc,
    track_temperature_celsius: temp,
  } })),
  ...MOTOGP_2025.map(([round, name, circuit, date, winner, front, rear, weather]) => ({ category_id: 'motogp', round, name, circuit, date, winner, payload: {
    rider_of_the_day_id: DRIVER_BY_CODE[winner], tyre_spec: { front, rear }, holeshot_winner_id: DRIVER_BY_CODE[winner], weather_conditions: weather,
  } })),
  ...WEC_2025.map(([round, name, circuit, date, hypercar, lmgt3, laps, fcy]) => ({ category_id: 'wec', round, name, circuit, date, winner: hypercar, payload: {
    total_laps_completed: laps, class_winners: { hypercar: DRIVER_BY_CODE[hypercar], lmgt3: DRIVER_BY_CODE[lmgt3] },
    full_course_yellow_count: fcy, red_flag_duration_minutes: 0,
  } })),
].map((r, i) => race2025(i + 1, r.category_id, r.round, r.name, r.circuit, r.date, r.winner, r.payload))

// Rows as a broken ETL run could store them: an f1 race with an empty payload and a car with half its parts. The API
// returns them flagged (payload_valid/parts_valid false) and the race/car pages fall back to the raw-JSON panel.
// The Go integration tests and front-end tests rely on both staying invalid.
const RACE_BROKEN = seq(77, 90)
const CAR_BROKEN = seq(74, 90)

// Older seasons for every team, plus 2026 cars for the extra teams. Each car reuses its team's engine and chassis:
// enough for paging, not a technical history. team index, category, season, name.
const MORE_CARS: [number, string, number, string][] = [
  [0, 'f1', 2023, 'RB19'], [0, 'f1', 2024, 'RB20'], [0, 'f1', 2025, 'RB21'],
  [1, 'motogp', 2023, 'Desmosedici GP23'], [1, 'motogp', 2024, 'Desmosedici GP24'], [1, 'motogp', 2025, 'Desmosedici GP25'],
  [2, 'wec', 2023, 'GR010 Hybrid'], [2, 'wec', 2024, 'GR010 Hybrid'], [2, 'wec', 2025, 'GR010 Hybrid'],
  [3, 'f1', 2023, 'SF-23'], [3, 'f1', 2024, 'SF-24'], [3, 'f1', 2025, 'SF-25'], [3, 'f1', 2026, 'SF-26'],
  [4, 'f1', 2023, 'MCL60'], [4, 'f1', 2024, 'MCL38'], [4, 'f1', 2025, 'MCL39'], [4, 'f1', 2026, 'MCL40'],
  [5, 'motogp', 2023, 'RS-GP23'], [5, 'motogp', 2024, 'RS-GP24'], [5, 'motogp', 2025, 'RS-GP25'], [5, 'motogp', 2026, 'RS-GP26'],
  [6, 'wec', 2023, '963'], [6, 'wec', 2024, '963'], [6, 'wec', 2025, '963'], [6, 'wec', 2026, '963'],
]

// cars.parts per team; every car of a team shares them.
const TEAM_PARTS: CarParts[] = [
  {
    gearbox: { type: 'SEQUENTIAL', speeds: 8 },
    brakes: { front: 'Carbon-carbon discs', rear: 'Carbon-carbon discs' },
    suspension: { front: 'Pushrod', rear: 'Pullrod' },
    tyre_supplier: 'Pirelli',
    other: { fuel_supplier: 'Petronas' },
  },
  {
    gearbox: { type: 'SEAMLESS', speeds: 6 },
    brakes: { front: 'Carbon discs', rear: 'Steel disc' },
    suspension: { front: 'Telescopic fork', rear: 'Monoshock swingarm' },
    tyre_supplier: 'Michelin',
    other: {},
  },
  {
    gearbox: { type: 'SEQUENTIAL', speeds: 7 },
    brakes: { front: 'Carbon discs', rear: 'Carbon discs' },
    suspension: { front: 'Double wishbone', rear: 'Double wishbone' },
    tyre_supplier: 'Michelin',
    other: { hybrid_system: 'Front-axle MGU' },
  },
  {
    gearbox: { type: 'SEQUENTIAL', speeds: 8 },
    brakes: { front: 'Carbon-carbon discs', rear: 'Carbon-carbon discs' },
    suspension: { front: 'Pushrod', rear: 'Pushrod' },
    tyre_supplier: 'Pirelli',
    other: { fuel_supplier: 'Shell' },
  },
  {
    gearbox: { type: 'SEQUENTIAL', speeds: 8 },
    brakes: { front: 'Carbon-carbon discs', rear: 'Carbon-carbon discs' },
    suspension: { front: 'Pushrod', rear: 'Pullrod' },
    tyre_supplier: 'Pirelli',
    other: { fuel_supplier: 'Aramco' },
  },
  {
    gearbox: { type: 'SEAMLESS', speeds: 6 },
    brakes: { front: 'Carbon discs', rear: 'Steel disc' },
    suspension: { front: 'Telescopic fork', rear: 'Monoshock swingarm' },
    tyre_supplier: 'Michelin',
    other: {},
  },
  {
    gearbox: { type: 'SEQUENTIAL', speeds: 7 },
    brakes: { front: 'Carbon discs', rear: 'Carbon discs' },
    suspension: { front: 'Double wishbone', rear: 'Double wishbone' },
    tyre_supplier: 'Michelin',
    other: { hybrid_system: 'Rear-axle MGU' },
  },
]

// A live blog on the 2025 Monaco GP: more articles about one entity than fit on a page (25), so the race page has to
// follow the pagination. Dated 2025, so the newest-first article list still opens with the spec's articles.
const RACE_MONACO_2025 = seq(77, 4)
const LIVE_UPDATES = 27
const liveAt = (i: number) => `2025-05-25T13:${String(i * 2).padStart(2, '0')}:00Z`

const text = (n: number, article: number, language: string, title: string, summary: string | null, body_md: string, at: string) => ({
  id: seq(76, n), article_id: ARTICLE[article - 1], language, title, summary, body_md, created_at: at, updated_at: at,
})

/** Rows as stored: what the back end computes on read (payload_valid, ...) is added by mock.ts. */
export const mockData: { [P in EntityPath]: Omit<EntityMap[P], ComputedFields>[] } = {
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
    ...MORE_CIRCUITS.map(([key, name, location, country, length_km]) => ({ id: CIRCUIT_BY_KEY[key], name, location, country, length_km, created_at: created })),
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
    ...MORE_RACES,
    {
      id: RACE_BROKEN, category_id: 'f1', circuit_id: CIRCUIT_BY_KEY.imola, winner_id: DRIVER.ver,
      name: 'Emilia Romagna Grand Prix', season: 2025, round: 7, race_date: '2025-05-18T00:00:00Z',
      payload: {},
      created_at: created,
    },
  ],
  teams: [
    { id: TEAM[0], name: 'Oracle Red Bull Racing', short_name: 'Red Bull', country: 'United Kingdom', founded_year: 2005, created_at: created },
    { id: TEAM[1], name: 'Ducati Lenovo Team', short_name: 'Ducati', country: 'Italy', founded_year: null, created_at: created },
    { id: TEAM[2], name: 'Toyota Gazoo Racing', short_name: 'Toyota', country: 'Japan', founded_year: 2007, created_at: created },
    { id: TEAM[3], name: 'Scuderia Ferrari HP', short_name: 'Ferrari', country: 'Italy', founded_year: 1929, created_at: created },
    { id: TEAM[4], name: 'McLaren Formula 1 Team', short_name: 'McLaren', country: 'United Kingdom', founded_year: 1963, created_at: created },
    { id: TEAM[5], name: 'Aprilia Racing', short_name: 'Aprilia', country: 'Italy', founded_year: null, created_at: created },
    { id: TEAM[6], name: 'Porsche Penske Motorsport', short_name: 'Porsche', country: 'Germany', founded_year: 2022, created_at: created },
  ],
  engines: [
    { id: ENGINE[0], manufacturer: 'Red Bull Ford', name: 'DM01', configuration: 'V6', displacement_cc: 1600, cylinders: 6, aspiration: 'HYBRID_TURBO', power_kw: 750.0, created_at: created },
    { id: ENGINE[1], manufacturer: 'Ducati', name: 'Desmosedici GP', configuration: 'V4', displacement_cc: 1000, cylinders: 4, aspiration: 'NATURALLY_ASPIRATED', power_kw: null, created_at: created },
    { id: ENGINE[2], manufacturer: 'Toyota', name: 'V35A', configuration: 'V6', displacement_cc: 3400, cylinders: 6, aspiration: 'TURBOCHARGED', power_kw: 520.0, created_at: created },
    { id: ENGINE[3], manufacturer: 'Ferrari', name: '067', configuration: 'V6', displacement_cc: 1600, cylinders: 6, aspiration: 'HYBRID_TURBO', power_kw: 745.0, created_at: created },
    { id: ENGINE[4], manufacturer: 'Mercedes-AMG', name: 'F1 M17 E Performance', configuration: 'V6', displacement_cc: 1600, cylinders: 6, aspiration: 'HYBRID_TURBO', power_kw: 748.0, created_at: created },
    { id: ENGINE[5], manufacturer: 'Aprilia', name: 'RS-GP V4', configuration: 'V4', displacement_cc: 1000, cylinders: 4, aspiration: 'NATURALLY_ASPIRATED', power_kw: 220.0, created_at: created },
    { id: ENGINE[6], manufacturer: 'Porsche', name: '9RD', configuration: 'V8', displacement_cc: 4600, cylinders: 8, aspiration: 'TURBOCHARGED', power_kw: 500.0, created_at: created },
  ],
  chassis: [
    { id: CHASSIS[0], manufacturer: 'Red Bull Racing', name: 'RB22 monocoque', material: 'Carbon fibre monocoque', wheelbase_mm: 3600, weight_kg: 768.0, created_at: created },
    { id: CHASSIS[1], manufacturer: 'Ducati', name: 'Desmosedici frame', material: 'Aluminium', wheelbase_mm: null, weight_kg: null, created_at: created },
    { id: CHASSIS[2], manufacturer: 'Toyota', name: 'GR010 monocoque', material: 'Carbon fibre monocoque', wheelbase_mm: null, weight_kg: 1030.0, created_at: created },
    { id: CHASSIS[3], manufacturer: 'Ferrari', name: 'SF-26 monocoque', material: 'Carbon fibre monocoque', wheelbase_mm: 3600, weight_kg: 768.0, created_at: created },
    { id: CHASSIS[4], manufacturer: 'McLaren', name: 'MCL40 monocoque', material: 'Carbon fibre monocoque', wheelbase_mm: 3600, weight_kg: 768.0, created_at: created },
    { id: CHASSIS[5], manufacturer: 'Aprilia', name: 'RS-GP frame', material: 'Aluminium twin-spar', wheelbase_mm: null, weight_kg: 157.0, created_at: created },
    { id: CHASSIS[6], manufacturer: 'Multimatic', name: '963 LMDh', material: 'Carbon fibre monocoque', wheelbase_mm: 3148, weight_kg: 1030.0, created_at: created },
  ],
  cars: [
    {
      id: CAR[0], category_id: 'f1', team_id: TEAM[0], engine_id: ENGINE[0], chassis_id: CHASSIS[0], name: 'RB22', season: 2026,
      parts: TEAM_PARTS[0],
      created_at: created,
    },
    {
      id: CAR[1], category_id: 'motogp', team_id: TEAM[1], engine_id: ENGINE[1], chassis_id: CHASSIS[1], name: 'Desmosedici GP26', season: 2026,
      parts: TEAM_PARTS[1],
      created_at: created,
    },
    {
      id: CAR[2], category_id: 'wec', team_id: TEAM[2], engine_id: ENGINE[2], chassis_id: CHASSIS[2], name: 'GR010 Hybrid', season: 2026,
      parts: TEAM_PARTS[2],
      created_at: created,
    },
    ...MORE_CARS.map(([team, category_id, season, name], i) => ({
      id: seq(74, i + 4), category_id, team_id: TEAM[team], engine_id: ENGINE[team], chassis_id: CHASSIS[team], name, season,
      parts: TEAM_PARTS[team], created_at: created,
    })),
    {
      id: CAR_BROKEN, category_id: 'f1', team_id: TEAM[0], engine_id: ENGINE[0], chassis_id: CHASSIS[0], name: 'RB18', season: 2022,
      parts: { gearbox: { type: 'SEQUENTIAL' }, tyre_supplier: 'Pirelli' },
      created_at: created,
    },
  ],
  articles: [
    { id: ARTICLE[0], entity_type: 'race', entity_id: RACE.monaco, kind: 'report', published_at: '2026-05-25T12:00:00Z', created_at: created },
    { id: ARTICLE[1], entity_type: 'category', entity_id: 'f1', kind: 'overview', published_at: '2026-01-02T00:00:00Z', created_at: created },
    { id: ARTICLE[2], entity_type: 'car', entity_id: CAR[0], kind: 'technical', published_at: null, created_at: created },
    { id: ARTICLE[3], entity_type: 'team', entity_id: TEAM[0], kind: 'history', published_at: '2026-02-01T00:00:00Z', created_at: created },
    { id: ARTICLE[4], entity_type: 'driver', entity_id: DRIVER.ver, kind: 'overview', published_at: '2026-03-01T00:00:00Z', created_at: created },
    // Must stay future-dated: mock.test.ts and back-end/scripts/smoke.sh check that it is never served. Don't "fix" it to a near date.
    { id: ARTICLE[5], entity_type: 'driver', entity_id: DRIVER.bag, kind: 'preview', published_at: '2099-01-01T00:00:00Z', created_at: created },
    ...Array.from({ length: LIVE_UPDATES }, (_, i) => ({
      id: seq(75, 101 + i), entity_type: 'race' as const, entity_id: RACE_MONACO_2025, kind: 'live', published_at: liveAt(i), created_at: created,
    })),
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
    ...Array.from({ length: LIVE_UPDATES }, (_, i) => ({
      id: seq(76, 101 + i), article_id: seq(75, 101 + i), language: 'en', title: `Monaco 2025 live, update ${i + 1}`, summary: null,
      body_md: `Lap-by-lap coverage from Monte Carlo, update ${i + 1}.\n`, created_at: liveAt(i), updated_at: liveAt(i),
    })),
  ],
}

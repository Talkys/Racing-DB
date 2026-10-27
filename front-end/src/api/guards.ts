import type { F1Payload, MotoGPPayload, Race, WECPayload } from './types.ts'

// The payload shape is decided by the race's category.

export const isF1Race = (race: Race): race is Race<F1Payload> => race.category_id === 'f1'
export const isMotoGPRace = (race: Race): race is Race<MotoGPPayload> => race.category_id === 'motogp'
export const isWECRace = (race: Race): race is Race<WECPayload> => race.category_id === 'wec'

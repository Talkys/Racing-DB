import { describe, expect, test } from 'bun:test'
import { pickLanguage } from '../src/lib/language.ts'

describe('pickLanguage', () => {
  test('returns null when there is no text', () => {
    expect(pickLanguage([], 'en', ['en'])).toBeNull()
  })

  test('prefers the exact requested tag', () => {
    expect(pickLanguage(['en', 'pt', 'pt-br'], 'pt-br', ['en'])).toBe('pt-br')
  })

  test('falls back from a region tag to its primary subtag', () => {
    expect(pickLanguage(['en', 'pt'], 'pt-br', [])).toBe('pt')
  })

  test('falls back from a primary subtag to a sibling region', () => {
    expect(pickLanguage(['en', 'pt-br'], 'pt', [])).toBe('pt-br')
  })

  test('is case-insensitive and returns the lowercase tag', () => {
    expect(pickLanguage(['EN', 'PT-BR'], 'pt-BR', [])).toBe('pt-br')
  })

  test('tries the browser languages in order after the requested one', () => {
    expect(pickLanguage(['de', 'fr', 'en'], 'es', ['fr-CA', 'en'])).toBe('fr')
  })

  test('the requested language beats the browser languages', () => {
    expect(pickLanguage(['de', 'fr'], 'de', ['fr'])).toBe('de')
  })

  test('uses the first available tag when nothing matches', () => {
    expect(pickLanguage(['de', 'fr'], 'ja', ['ko'])).toBe('de')
    expect(pickLanguage(['de', 'fr'], null, [])).toBe('de')
  })

  test('never invents a language', () => {
    expect(pickLanguage(['pt'], 'en', ['en-US'])).toBe('pt')
  })
})

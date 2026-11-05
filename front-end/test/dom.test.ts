import { afterAll, beforeAll, beforeEach, describe, expect, test } from 'bun:test'
import { GlobalRegistrator } from '@happy-dom/global-registrator'
import { fill, fillList, section, showError } from '../src/lib/dom.ts'

beforeAll(() => GlobalRegistrator.register())
afterAll(() => GlobalRegistrator.unregister())

const html = (markup: string) => {
  document.body.innerHTML = markup
  return document.body
}

describe('fill', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
  })

  test('copies values into data-field, data-field-href and data-field-value', () => {
    const root = html(`
      <h1 data-field="name">placeholder</h1>
      <a data-field-href="href">link</a>
      <input data-field-value="q">`)
    fill(root, { name: '<b>Monaco</b>', href: '/race?id=1', q: 'ver' })

    expect(root.querySelector('h1')!.textContent).toBe('<b>Monaco</b>')
    expect(root.querySelector('h1')!.children.length).toBe(0) // text, never HTML
    expect(root.querySelector('a')!.getAttribute('href')).toBe('/race?id=1')
    expect(root.querySelector('input')!.value).toBe('ver')
  })

  test('leaves elements untouched when their key is absent or null', () => {
    const root = html('<span data-field="a">keep a</span><span data-field="b">keep b</span>')
    fill(root, { b: null })
    expect(root.textContent).toBe('keep akeep b')
  })

  test('writes empty strings', () => {
    const root = html('<span data-field="a">placeholder</span>')
    fill(root, { a: '' })
    expect(root.querySelector('span')!.textContent).toBe('')
  })

  test('data-show-if hides on empty values and only reacts to keys it is given', () => {
    const root = html('<p data-show-if="full">full</p><p data-show-if="empty">empty</p><p data-show-if="other" hidden>other</p>')
    fill(root, { full: 'x', empty: '' })
    const [full, empty, other] = root.querySelectorAll('p')
    expect(full.hidden).toBe(false)
    expect(empty.hidden).toBe(true)
    expect(other.hidden).toBe(true)
  })
})

describe('fillList', () => {
  test('clones the template once per row, after the template, in order', () => {
    const root = html('<ul><template><li data-field="name"></li></template></ul>')
    const list = root.querySelector('ul')!
    fillList(list, [{ name: 'a' }, { name: 'b' }, { name: 'c' }])
    expect([...list.querySelectorAll('li')].map((li) => li.textContent)).toEqual(['a', 'b', 'c'])
  })

  test('appends to the template parent, so tables keep their <tbody>', () => {
    // Built by hand: happy-dom's parser moves a <template> out of <tbody>, which browsers don't do.
    const root = html('<table><tbody></tbody></table>')
    const template = document.createElement('template')
    template.innerHTML = '<tr><td data-field="x"></td></tr>'
    root.querySelector('tbody')!.appendChild(template)
    fillList(root.querySelector('table')!, [{ x: '1' }])
    expect(root.querySelector('tbody > tr > td')!.textContent).toBe('1')
  })
})

describe('section and showError', () => {
  test('section finds a data-section element', () => {
    const root = html('<div data-section="articles">x</div>')
    expect(section('articles')).toBe(root.querySelector('div')!)
  })

  test('showError fills and reveals the error banner', () => {
    const root = html('<p data-error hidden></p>')
    showError(new Error('boom'))
    const box = root.querySelector<HTMLElement>('[data-error]')!
    expect(box.textContent).toBe('boom')
    expect(box.hidden).toBe(false)

    showError('plain string')
    expect(box.textContent).toBe('plain string')
  })
})

import fs from 'node:fs'
import { resolve } from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [
    tailwindcss(),
    {
      name: 'html-ext-fallback',
      configureServer(server) {
        return () => {
          server.middlewares.use((req, _res, next) => {
            if (req.url && !req.url.includes('.') && req.url !== '/') {
              const htmlPath = resolve(__dirname, `${req.url.slice(1)}.html`)
              if (fs.existsSync(htmlPath)) {
                req.url = `${req.url}.html`
              }
            }
            next()
          })
        }
      }
    }
  ],
  server: {
    // Back end listens on :5000. To work without it, run with VITE_API_MOCK=true (see src/api/client.ts).
    proxy: { '/api': process.env.VITE_API_TARGET ?? 'http://localhost:5000' },
  },
  build: {
    rollupOptions: {
      // Multi page app: every top-level .html file is an entry.
      input: Object.fromEntries(
        fs.readdirSync(__dirname)
          .filter((file) => file.endsWith('.html'))
          .map((file) => [file.slice(0, -5), resolve(__dirname, file)]),
      ),
    },
  },
})
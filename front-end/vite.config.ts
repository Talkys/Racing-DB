import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'path'
import fs from 'node:fs'

export default defineConfig({
  plugins: [
    tailwindcss(),
    {
      name: 'html-ext-fallback',
      configureServer(server) {
        return () => {
          server.middlewares.use((req, res, next) => {
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
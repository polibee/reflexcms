import tailwindcss from '@tailwindcss/vite'

// Read Node env without pulling @types/node into the app tsconfig.
const nodeEnv = (globalThis as { process?: { env?: Record<string, string | undefined> } })
  .process?.env ?? {}

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@pinia/nuxt',
    '@nuxt/eslint'
  ],

  components: [
    { path: '~/admin/ui', pathPrefix: false },
    { path: '~/admin/framework', pathPrefix: false }
  ],

  // Framework DSL auto-imports (builders available everywhere)
  imports: {
    dirs: [
      '~/stores',
      '~/admin/core',
      '~/admin/panel',
      '~/admin/navigation',
      '~/admin/permissions',
      '~/admin/schemas/builders',
      '~/admin/infolists',
      '~/admin/actions',
      '~/admin/widgets',
      '~/admin/modules',
      '~/admin/forms',
      '~/admin/tables',
      '~/admin/framework',
      '~/admin/notifications'
    ]
  },

  devtools: { enabled: true },

  css: ['~/assets/css/main.css'],

  runtimeConfig: {
    // Server-side only: where the BFF proxies admin requests to.
    // Override with BACKEND_URL in .env for non-local deployments.
    backendUrl: nodeEnv.BACKEND_URL || 'http://127.0.0.1:9000',
    // Session cookie Secure flag; derived here because server/utils has no
    // @types/node.
    cookieSecure: nodeEnv.NODE_ENV === 'production'
  },

  routeRules: {
    '/api/**': { cors: true }
  },

  compatibilityDate: '2026-06-30',

  vite: {
    plugins: [tailwindcss()]
  },

  eslint: {
    config: {
      stylistic: {
        commaDangle: 'never',
        braceStyle: '1tbs'
      }
    }
  }
})

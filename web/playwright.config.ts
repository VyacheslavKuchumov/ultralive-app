import { defineConfig, devices } from '@playwright/test'

function envIsTrue(value?: string) {
  return value === '1' || value === 'true'
}

const workers = process.env.CI ? 1 : process.env.PLAYWRIGHT_WORKERS

export default defineConfig({
  testDir: './e2e',
  timeout: 90_000,
  expect: {
    timeout: 10_000
  },
  fullyParallel: envIsTrue(process.env.PLAYWRIGHT_FULLY_PARALLEL),
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers,
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL: process.env.PLAYWRIGHT_BASE_URL || 'http://127.0.0.1:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    headless: true
  },
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome']
      }
    }
  ]
})

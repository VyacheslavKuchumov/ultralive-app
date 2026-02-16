import { expect, test, type Page } from '@playwright/test'

const BACKEND_URL = process.env.PLAYWRIGHT_BACKEND_URL || 'http://127.0.0.1:8000'

type TestUser = {
  firstName: string
  lastName: string
  email: string
  password: string
}

function uniqueSuffix() {
  return `${Date.now()}_${Math.floor(Math.random() * 100000)}`
}

function buildUser(prefix: string): TestUser {
  const suffix = uniqueSuffix()
  return {
    firstName: `QA_${prefix}_Name_${suffix}`,
    lastName: `QA_${prefix}_Last_${suffix}`,
    email: `qa_${prefix.toLowerCase()}_${suffix}@example.com`,
    password: `Pwd_${suffix}`
  }
}

async function waitForBackend() {
  const startedAt = Date.now()
  let lastError = ''

  while (Date.now() - startedAt < 60_000) {
    try {
      const response = await fetch(`${BACKEND_URL}/api/v1/login`, {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: '{}'
      })

      if (response.status < 500) return
      lastError = `status ${response.status}`
    } catch (error: any) {
      lastError = error?.message || String(error)
    }

    await new Promise((resolve) => setTimeout(resolve, 1000))
  }

  throw new Error(
    `Backend is not reachable at ${BACKEND_URL}. Start postgres+server+web before e2e run. Last error: ${lastError}`
  )
}

async function signUp(page: Page, user: TestUser) {
  await page.goto('/signup')

  await page.getByLabel(/^Имя/).fill(user.firstName)
  await page.getByLabel(/^Фамилия/).fill(user.lastName)
  await page.getByLabel(/^Email/).fill(user.email)
  await page.getByLabel(/^Пароль/).fill(user.password)

  await Promise.all([
    page.waitForURL('**/'),
    page.getByRole('button', { name: 'Зарегистрироваться' }).click()
  ])
}

test.beforeAll(async () => {
  await waitForBackend()
})

test('user can register and manage CRM dictionaries', async ({ page }) => {
  const user = buildUser('Flow')
  const suffix = uniqueSuffix()
  const setTypeName = `Тип_${suffix}`
  const projectTypeName = `Площадка_${suffix}`
  const warehouseName = `Склад_${suffix}`
  const draftName = `Шаблон_${suffix}`

  await signUp(page, user)
  await expect(page.getByRole('heading', { name: 'Ultralive CRM' })).toBeVisible()

  await page.goto('/set_types')
  const setTypeInput = page.getByPlaceholder('Название вида')
  await setTypeInput.fill(setTypeName)
  await expect(setTypeInput).toHaveValue(setTypeName)
  await page.locator('form').first().getByRole('button', { name: 'Добавить' }).click()
  await expect(page.getByText(setTypeName)).toBeVisible()

  await page.goto('/project_types')
  await page.getByPlaceholder('Название площадки').fill(projectTypeName)
  await page.getByPlaceholder('Neaktor ID (опционально)').fill(`N_${suffix}`)
  await page.getByRole('button', { name: 'Добавить' }).click()
  await expect(page.getByText(projectTypeName)).toBeVisible()

  await page.goto('/warehouses')
  await page.getByPlaceholder('Название склада').fill(warehouseName)
  await page.getByPlaceholder('Адрес').fill(`Адрес ${suffix}`)
  await page.getByRole('button', { name: 'Добавить' }).click()
  await expect(page.getByText(warehouseName)).toBeVisible()

  await page.goto('/drafts')
  await page.getByPlaceholder('Название шаблона').fill(draftName)
  await page.getByRole('button', { name: 'Добавить' }).click()
  await expect(page.getByText(draftName)).toBeVisible()

  const draftRow = page.locator('tr', { hasText: draftName }).first()
  await draftRow.getByRole('link', { name: 'Состав' }).click()
  await expect(page.getByRole('heading', { name: `Состав шаблона: ${draftName}` })).toBeVisible()
})

test('user can logout and login again', async ({ page }) => {
  const user = buildUser('Auth')
  const updatedLastName = `${user.lastName}_Updated`
  const nextPassword = `${user.password}_new`

  await signUp(page, user)
  await expect(page.getByRole('heading', { name: 'Ultralive CRM' })).toBeVisible()

  await page.getByRole('banner').getByRole('button', { name: 'Выйти' }).click()
  await expect(page).toHaveURL(/\/login$/)

  await page.getByLabel(/^Email/).fill(user.email)
  await page.getByLabel(/^Пароль/).fill(user.password)

  await Promise.all([
    page.waitForURL('**/'),
    page.getByRole('button', { name: 'Войти' }).click()
  ])

  await expect(page.getByRole('heading', { name: 'Ultralive CRM' })).toBeVisible()
  await page.goto('/profile')

  await expect(page.getByLabel(/^Email/)).toHaveValue(user.email)
  await page.getByLabel(/^Фамилия/).fill(updatedLastName)
  await page.getByLabel('Подтверждение текущим паролем').fill(user.password)
  await page.getByRole('button', { name: 'Сохранить профиль' }).click()
  await expect(page.getByLabel(/^Фамилия/)).toHaveValue(updatedLastName)

  await page.getByLabel(/^Текущий пароль$/).fill(user.password)
  await page.getByLabel(/^Новый пароль$/).fill(nextPassword)
  await page.getByLabel('Повтор нового пароля').fill(nextPassword)
  await page.getByRole('button', { name: 'Обновить пароль' }).click()

  await page.getByRole('banner').getByRole('button', { name: 'Выйти' }).click()
  await expect(page).toHaveURL(/\/login$/)

  await page.getByLabel(/^Email/).fill(user.email)
  await page.getByLabel(/^Пароль/).fill(nextPassword)
  await Promise.all([
    page.waitForURL('**/'),
    page.getByRole('button', { name: 'Войти' }).click()
  ])
  await expect(page.getByRole('heading', { name: 'Ultralive CRM' })).toBeVisible()
})

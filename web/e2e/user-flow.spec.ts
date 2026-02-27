import { expect, test, type Locator, type Page } from '@playwright/test'

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

function rowsByText(page: Page, text: string): Locator {
  return page.locator('tbody tr').filter({ hasText: text })
}

function rowByText(page: Page, text: string): Locator {
  return rowsByText(page, text).first()
}

function equipmentSetCards(page: Page, text: string): Locator {
  return page.locator('[data-testid^="equipment-set-card-"]').filter({ hasText: text })
}

function equipmentSetCard(page: Page, text: string): Locator {
  return equipmentSetCards(page, text).first()
}

function projectCards(page: Page, text: string): Locator {
  return page.locator('[data-testid^="project-card-"]').filter({ hasText: text })
}

function projectCard(page: Page, text: string): Locator {
  return projectCards(page, text).first()
}

async function expectRowAbsent(page: Page, text: string) {
  await expect(rowsByText(page, text)).toHaveCount(0)
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

async function openCreateDialog(page: Page) {
  await page.getByRole('button', { name: 'Добавить' }).first().click()
  const dialog = page.getByRole('dialog').last()
  await expect(dialog).toBeVisible()
  return dialog
}

async function openEditDialog(page: Page, rowText: string) {
  await rowByText(page, rowText).getByRole('button', { name: 'Изменить' }).click()
  const dialog = page.getByRole('dialog').last()
  await expect(dialog).toBeVisible()
  return dialog
}

async function selectInDialog(page: Page, dialog: Locator, triggerText: string, optionText: string) {
  const labelTrigger = dialog.getByLabel(triggerText).first()
  if (await labelTrigger.count()) {
    await labelTrigger.click()
  } else {
    const comboboxTrigger = dialog.getByRole('combobox', { name: triggerText }).first()
    if (await comboboxTrigger.count()) {
      await comboboxTrigger.click()
    } else {
      await dialog.locator('button').filter({ hasText: triggerText }).first().click()
    }
  }
  await page.getByRole('option', { name: optionText }).first().click()
}

async function createSetType(page: Page, name: string) {
  await page.goto('/set_types')
  const dialog = await openCreateDialog(page)
  await dialog.getByPlaceholder('Название вида').fill(name)
  await dialog.getByRole('button', { name: 'Создать' }).click()
  await expect(rowByText(page, name)).toBeVisible()
}

async function editSetType(page: Page, currentName: string, nextName: string) {
  await page.goto('/set_types')
  const dialog = await openEditDialog(page, currentName)
  await dialog.getByPlaceholder('Название вида').fill(nextName)
  await dialog.getByRole('button', { name: 'Сохранить' }).click()
  await expect(rowByText(page, nextName)).toBeVisible()
  await expectRowAbsent(page, currentName)
}

async function deleteSetType(page: Page, name: string) {
  await page.goto('/set_types')
  await rowByText(page, name).getByRole('button', { name: 'Удалить' }).click()
  await expectRowAbsent(page, name)
}

async function createProjectType(page: Page, name: string, neaktorId: string) {
  await page.goto('/project_types')
  const dialog = await openCreateDialog(page)
  await dialog.getByPlaceholder('Название площадки').fill(name)
  await dialog.getByPlaceholder('Neaktor ID (опционально)').fill(neaktorId)
  await dialog.getByRole('button', { name: 'Создать' }).click()
  await expect(rowByText(page, name)).toBeVisible()
}

async function editProjectType(page: Page, currentName: string, nextName: string, neaktorId: string) {
  await page.goto('/project_types')
  const dialog = await openEditDialog(page, currentName)
  await dialog.getByPlaceholder('Название площадки').fill(nextName)
  await dialog.getByPlaceholder('Neaktor ID (опционально)').fill(neaktorId)
  await dialog.getByRole('button', { name: 'Сохранить' }).click()
  await expect(rowByText(page, nextName)).toBeVisible()
  await expectRowAbsent(page, currentName)
}

async function deleteProjectType(page: Page, name: string) {
  await page.goto('/project_types')
  await rowByText(page, name).getByRole('button', { name: 'Удалить' }).click()
  await expectRowAbsent(page, name)
}

async function createWarehouse(page: Page, name: string, address: string) {
  await page.goto('/warehouses')
  const dialog = await openCreateDialog(page)
  await dialog.getByPlaceholder('Название склада').fill(name)
  await dialog.getByPlaceholder('Адрес').fill(address)
  await dialog.getByRole('button', { name: 'Создать' }).click()
  await expect(rowByText(page, name)).toBeVisible()
}

async function editWarehouse(page: Page, currentName: string, nextName: string, address: string) {
  await page.goto('/warehouses')
  const dialog = await openEditDialog(page, currentName)
  await dialog.getByPlaceholder('Название склада').fill(nextName)
  await dialog.getByPlaceholder('Адрес').fill(address)
  await dialog.getByRole('button', { name: 'Сохранить' }).click()
  await expect(rowByText(page, nextName)).toBeVisible()
  await expectRowAbsent(page, currentName)
}

async function deleteWarehouse(page: Page, name: string) {
  await page.goto('/warehouses')
  await rowByText(page, name).getByRole('button', { name: 'Удалить' }).click()
  await expectRowAbsent(page, name)
}

async function createDraft(page: Page, name: string) {
  await page.goto('/drafts')
  const dialog = await openCreateDialog(page)
  await dialog.getByPlaceholder('Название шаблона').fill(name)
  await dialog.getByRole('button', { name: 'Создать' }).click()
  await expect(rowByText(page, name)).toBeVisible()
}

async function editDraft(page: Page, currentName: string, nextName: string) {
  await page.goto('/drafts')
  const dialog = await openEditDialog(page, currentName)
  await dialog.getByPlaceholder('Название шаблона').fill(nextName)
  await dialog.getByRole('button', { name: 'Сохранить' }).click()
  await expect(rowByText(page, nextName)).toBeVisible()
  await expectRowAbsent(page, currentName)
}

async function deleteDraft(page: Page, name: string) {
  await page.goto('/drafts')
  await rowByText(page, name).getByRole('button', { name: 'Удалить' }).click()
  await expectRowAbsent(page, name)
}

async function createEquipmentSet(page: Page, name: string, description: string, setTypeName: string) {
  await page.goto('/equipment_sets')
  const dialog = await openCreateDialog(page)
  await dialog.getByPlaceholder('Название комплекта').fill(name)
  await dialog.getByPlaceholder('Описание').fill(description)
  await selectInDialog(page, dialog, 'Вид комплекта', setTypeName)
  await dialog.getByRole('button', { name: 'Создать' }).click()
  await expect(equipmentSetCard(page, name)).toBeVisible()
}

async function editEquipmentSet(page: Page, currentName: string, nextName: string, description: string) {
  await page.goto('/equipment_sets')
  await equipmentSetCard(page, currentName).getByRole('button', { name: 'Изменить' }).click()
  const dialog = page.getByRole('dialog').last()
  await expect(dialog).toBeVisible()
  await dialog.getByPlaceholder('Название комплекта').fill(nextName)
  await dialog.getByPlaceholder('Описание').fill(description)
  await dialog.getByRole('button', { name: 'Сохранить' }).click()
  await expect(equipmentSetCard(page, nextName)).toBeVisible()
  await expect(equipmentSetCards(page, currentName)).toHaveCount(0)
}

async function deleteEquipmentSet(page: Page, name: string) {
  await page.goto('/equipment_sets')
  await equipmentSetCard(page, name).getByRole('button', { name: 'Удалить' }).click()
  await expect(equipmentSetCards(page, name)).toHaveCount(0)
}

async function openEquipmentSet(page: Page, equipmentSetName: string) {
  await page.goto('/equipment_sets')
  await equipmentSetCard(page, equipmentSetName).getByRole('link', { name: 'Открыть комплект' }).click()
  await expect(page).toHaveURL(/\/equipment\?set=\d+/)
}

async function createEquipment(page: Page, name: string, serial: string, equipmentSetName: string, warehouseName: string) {
  await openEquipmentSet(page, equipmentSetName)
  const dialog = await openCreateDialog(page)
  await dialog.getByPlaceholder('Название').fill(name)
  await dialog.getByPlaceholder('Серийный номер').fill(serial)
  await selectInDialog(page, dialog, 'Склад', warehouseName)
  await dialog.getByRole('button', { name: 'Создать' }).click()
  await expect(rowByText(page, name)).toBeVisible()
}

async function editEquipment(page: Page, currentName: string, nextName: string, serial: string, equipmentSetName: string) {
  await openEquipmentSet(page, equipmentSetName)
  const dialog = await openEditDialog(page, currentName)
  await dialog.getByPlaceholder('Название').fill(nextName)
  await dialog.getByPlaceholder('Серийный номер').fill(serial)
  await dialog.getByRole('button', { name: 'Сохранить' }).click()
  await expect(rowByText(page, nextName)).toBeVisible()
  await expectRowAbsent(page, currentName)
}

async function deleteEquipment(page: Page, name: string, equipmentSetName: string) {
  await openEquipmentSet(page, equipmentSetName)
  await rowByText(page, name).getByRole('button', { name: 'Удалить' }).click()
  await expectRowAbsent(page, name)
}

async function createProject(
  page: Page,
  name: string,
  projectTypeName: string,
  chiefEngineerName: string,
  startDate: string,
  endDate: string
) {
  await page.goto('/projects')
  const dialog = await openCreateDialog(page)
  await dialog.getByPlaceholder('Название съёмки').fill(name)
  await selectInDialog(page, dialog, 'Площадка', projectTypeName)
  await selectInDialog(page, dialog, 'Главный инженер', chiefEngineerName)
  await dialog.getByPlaceholder('Дата начала YYYY-MM-DD').fill(startDate)
  await dialog.getByPlaceholder('Дата конца YYYY-MM-DD').fill(endDate)
  await dialog.getByRole('button', { name: 'Создать' }).click()
  const card = projectCard(page, name)
  await expect(card).toBeVisible()
  await expect(card).toContainText(chiefEngineerName)
}

async function editProject(page: Page, currentName: string, nextName: string, startDate: string, endDate: string) {
  await page.goto('/projects')
  await projectCard(page, currentName).getByRole('button', { name: 'Изменить' }).click()
  const dialog = page.getByRole('dialog').last()
  await expect(dialog).toBeVisible()
  await dialog.getByPlaceholder('Название съёмки').fill(nextName)
  await dialog.getByPlaceholder('Дата начала YYYY-MM-DD').fill(startDate)
  await dialog.getByPlaceholder('Дата конца YYYY-MM-DD').fill(endDate)
  await dialog.getByRole('button', { name: 'Сохранить' }).click()
  await expect(projectCard(page, nextName)).toBeVisible()
  await expect(projectCards(page, currentName)).toHaveCount(0)
}

async function deleteProject(page: Page, name: string) {
  await page.goto('/projects')
  await projectCard(page, name).getByRole('button', { name: 'Удалить' }).click()
  await expect(projectCards(page, name)).toHaveCount(0)
}

test.beforeAll(async () => {
  await waitForBackend()
})

async function signUpAndOpenDashboard(page: Page, prefix: string) {
  const user = buildUser(prefix)
  await signUp(page, user)
  await expect(page.getByRole('heading', { name: 'Ultralive CRM' })).toBeVisible()
  return user
}

test.describe('crud forms', () => {
  test('set types: create/edit/delete', async ({ page }) => {
    test.slow()
    const user = await signUpAndOpenDashboard(page, 'SetTypes')
    const suffix = user.email.split('@')[0]
    const setTypeCreate = `QA_SetType_Create_${suffix}`
    const setTypeEdit = `QA_SetType_Edit_${suffix}`

    await createSetType(page, setTypeCreate)
    await editSetType(page, setTypeCreate, setTypeEdit)
    await deleteSetType(page, setTypeEdit)
  })

  test('project types: create/edit/delete', async ({ page }) => {
    test.slow()
    const user = await signUpAndOpenDashboard(page, 'ProjectTypes')
    const suffix = user.email.split('@')[0]
    const projectTypeCreate = `QA_ProjectType_Create_${suffix}`
    const projectTypeEdit = `QA_ProjectType_Edit_${suffix}`

    await createProjectType(page, projectTypeCreate, `N_CREATE_${suffix}`)
    await editProjectType(page, projectTypeCreate, projectTypeEdit, `N_EDIT_${suffix}`)
    await deleteProjectType(page, projectTypeEdit)
  })

  test('warehouses: create/edit/delete', async ({ page }) => {
    test.slow()
    const user = await signUpAndOpenDashboard(page, 'Warehouses')
    const suffix = user.email.split('@')[0]
    const warehouseCreate = `QA_Warehouse_Create_${suffix}`
    const warehouseEdit = `QA_Warehouse_Edit_${suffix}`

    await createWarehouse(page, warehouseCreate, `Адрес создание ${suffix}`)
    await editWarehouse(page, warehouseCreate, warehouseEdit, `Адрес редактирование ${suffix}`)
    await deleteWarehouse(page, warehouseEdit)
  })

  test('drafts: create/edit/delete', async ({ page }) => {
    test.slow()
    const user = await signUpAndOpenDashboard(page, 'Drafts')
    const suffix = user.email.split('@')[0]
    const draftCreate = `QA_Draft_Create_${suffix}`
    const draftEdit = `QA_Draft_Edit_${suffix}`

    await createDraft(page, draftCreate)
    await editDraft(page, draftCreate, draftEdit)
    await deleteDraft(page, draftEdit)
  })

  test('equipment sets: create/edit/delete', async ({ page }) => {
    test.slow()
    const user = await signUpAndOpenDashboard(page, 'EquipmentSets')
    const suffix = user.email.split('@')[0]
    const dependencySetType = `QA_SetType_Dependency_${suffix}`
    const equipmentSetCreate = `QA_EqSet_Create_${suffix}`
    const equipmentSetEdit = `QA_EqSet_Edit_${suffix}`

    await createSetType(page, dependencySetType)
    await createEquipmentSet(page, equipmentSetCreate, `Описание создание ${suffix}`, dependencySetType)
    await editEquipmentSet(page, equipmentSetCreate, equipmentSetEdit, `Описание редактирование ${suffix}`)
    await deleteEquipmentSet(page, equipmentSetEdit)
    await deleteSetType(page, dependencySetType)
  })

  test('equipment in set: create/edit/delete', async ({ page }) => {
    test.slow()
    const user = await signUpAndOpenDashboard(page, 'Equipment')
    const suffix = user.email.split('@')[0]
    const equipmentSetType = `QA_SetType_Equipment_${suffix}`
    const equipmentWarehouse = `QA_Warehouse_Equipment_${suffix}`
    const equipmentSetName = `QA_EqSet_Equipment_${suffix}`
    const equipmentCreate = `QA_Equipment_Create_${suffix}`
    const equipmentEdit = `QA_Equipment_Edit_${suffix}`

    await createSetType(page, equipmentSetType)
    await createWarehouse(page, equipmentWarehouse, `Адрес оборудования ${suffix}`)
    await createEquipmentSet(page, equipmentSetName, `Комплект для оборудования ${suffix}`, equipmentSetType)

    await createEquipment(page, equipmentCreate, `SER_CREATE_${suffix}`, equipmentSetName, equipmentWarehouse)
    await editEquipment(page, equipmentCreate, equipmentEdit, `SER_EDIT_${suffix}`, equipmentSetName)
    await deleteEquipment(page, equipmentEdit, equipmentSetName)

    await deleteEquipmentSet(page, equipmentSetName)
    await deleteSetType(page, equipmentSetType)
    await deleteWarehouse(page, equipmentWarehouse)
  })

  test('projects: create/edit/delete', async ({ page }) => {
    test.slow()
    const user = await signUpAndOpenDashboard(page, 'Projects')
    const suffix = user.email.split('@')[0]
    const chiefEngineerName = `${user.firstName} ${user.lastName}`
    const projectTypeName = `QA_ProjectType_Project_${suffix}`
    const projectCreate = `QA_Project_Create_${suffix}`
    const projectEdit = `QA_Project_Edit_${suffix}`

    await createProjectType(page, projectTypeName, `N_PROJECT_${suffix}`)
    await createProject(page, projectCreate, projectTypeName, chiefEngineerName, '2030-01-10', '2030-01-11')
    await editProject(page, projectCreate, projectEdit, '2030-01-12', '2030-01-13')
    await deleteProject(page, projectEdit)
    await deleteProjectType(page, projectTypeName)
  })
})

test('auth: logout/login and profile/password update', async ({ page }) => {
  const user = await signUpAndOpenDashboard(page, 'Auth')
  const updatedLastName = `${user.lastName}_Updated`
  const nextPassword = `${user.password}_new`

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

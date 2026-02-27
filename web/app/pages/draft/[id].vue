<template>
  <div class="space-y-6">
    <UCard v-if="board?.draft">
      <template #header>
        <h1 class="text-xl font-semibold">Состав шаблона: {{ board.draft.draft_name }}</h1>
      </template>
    </UCard>

    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold">Добавить в шаблон</h2>
      </template>

      <div class="grid gap-4 md:grid-cols-2">
        <div class="space-y-2">
          <label class="text-sm text-gray-600">Оборудование</label>
          <select v-model.number="selectedEquipmentId" class="w-full rounded border border-gray-300 px-3 py-2">
            <option :value="0">Выберите оборудование</option>
            <option v-for="item in board?.available_equipment || []" :key="item.equipment_id" :value="item.equipment_id">
              {{ item.equipment_name }} ({{ item.serial_number }})
            </option>
          </select>
          <UButton color="primary" variant="soft" @click="addEquipment">Добавить оборудование</UButton>
        </div>

        <div class="space-y-2">
          <label class="text-sm text-gray-600">Комплект</label>
          <select v-model.number="selectedSetId" class="w-full rounded border border-gray-300 px-3 py-2">
            <option :value="0">Выберите комплект</option>
            <option v-for="item in board?.sets_in_draft || []" :key="item.equipment_set_id" :value="item.equipment_set_id">
              {{ item.equipment_set_name }}
            </option>
          </select>
          <UButton color="primary" variant="soft" @click="addSet">Добавить комплект</UButton>
        </div>
      </div>
    </UCard>

    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold">Оборудование в шаблоне</h2>
      </template>

      <div class="space-y-4">
        <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <UInput
            v-model="search"
            icon="i-lucide-search"
            placeholder="Поиск по оборудованию шаблона"
            class="md:max-w-sm"
          />

          <label class="flex items-center gap-2 text-sm text-gray-600">
            На странице
            <select v-model.number="perPage" class="rounded border border-gray-300 px-2 py-1 text-sm">
              <option v-for="option in perPageOptions" :key="option" :value="option">{{ option }}</option>
            </select>
          </label>
        </div>

        <div v-if="groupedEquipment.length" class="grid gap-4">
          <UCard v-for="group in groupedEquipment" :key="group.equipment_set_id">
            <template #header>
              <h3 class="text-base font-semibold">{{ group.equipment_set_name }}</h3>
            </template>

            <div class="overflow-x-auto">
              <table class="w-full min-w-max text-sm">
                <thead>
                  <tr class="text-left border-b border-gray-200 whitespace-nowrap">
                    <th class="py-2">ID</th>
                    <th class="py-2">Название</th>
                    <th class="py-2">Серия</th>
                    <th class="py-2">Склад</th>
                    <th class="py-2">Действие</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="item in group.items"
                    :key="item.equipment_id"
                    :class="[
                      'border-b border-gray-100',
                      item.is_in_draft ? '' : 'bg-gray-50 text-gray-400'
                    ]"
                  >
                    <td class="py-2">{{ item.equipment_id }}</td>
                    <td class="py-2">{{ item.equipment_name }}</td>
                    <td class="py-2">{{ item.serial_number }}</td>
                    <td class="py-2">{{ item.storage?.warehouse_name || '-' }}</td>
                    <td class="py-2">
                      <UButton
                        v-if="item.is_in_draft"
                        size="xs"
                        color="error"
                        variant="soft"
                        @click="removeEquipment(item)"
                      >
                        Удалить
                      </UButton>
                      <UButton
                        v-else
                        size="xs"
                        color="neutral"
                        variant="soft"
                        @click="addEquipmentByItem(item)"
                      >
                        Добавить
                      </UButton>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </UCard>
        </div>

        <p v-else class="text-sm text-gray-600">Ничего не найдено.</p>

        <div class="flex flex-col gap-3 border-t border-gray-100 pt-3 md:flex-row md:items-center md:justify-between">
          <p class="text-sm text-gray-600">Показано {{ from }}-{{ to }} из {{ total }}</p>

          <div class="flex items-center gap-2">
            <UButton size="xs" color="neutral" variant="soft" :disabled="page <= 1" @click="prevPage">Назад</UButton>
            <span class="text-sm text-gray-600">Стр. {{ page }} / {{ totalPages }}</span>
            <UButton size="xs" color="neutral" variant="soft" :disabled="page >= totalPages" @click="nextPage">Вперед</UButton>
          </div>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from '#imports'
import { useCRMStore } from '~/stores/crm'

const route = useRoute()
const crm = useCRMStore()
const draftId = computed(() => Number(route.params.id))

const selectedEquipmentId = ref(0)
const selectedSetId = ref(0)

const search = ref('')
const page = ref(1)
const perPage = ref(10)
const perPageOptions = [10, 20, 50]

await crm.fetchDraftBoard(draftId.value)

const board = computed(() => crm.draftBoard)

const allEquipment = computed(() => {
  const byId = new Map()

  for (const item of board.value?.available_equipment || []) {
    byId.set(item.equipment_id, {
      ...item,
      is_in_draft: false
    })
  }

  for (const item of board.value?.equipment_in_draft || []) {
    byId.set(item.equipment_id, {
      ...item,
      is_in_draft: true
    })
  }

  return Array.from(byId.values()).sort((a, b) => {
    const aSetName = (a.equipment_set?.equipment_set_name || '').toLowerCase()
    const bSetName = (b.equipment_set?.equipment_set_name || '').toLowerCase()
    if (aSetName < bSetName) return -1
    if (aSetName > bSetName) return 1
    return a.equipment_id - b.equipment_id
  })
})

const filteredEquipment = computed(() => {
  const term = search.value.trim().toLowerCase()
  if (!term) return allEquipment.value

  return allEquipment.value.filter((item) => {
    const values = [
      String(item.equipment_id),
      item.equipment_name,
      item.serial_number,
      item.equipment_set?.equipment_set_name || '',
      item.storage?.warehouse_name || ''
    ]

    return values.some((value) => String(value).toLowerCase().includes(term))
  })
})

const total = computed(() => filteredEquipment.value.length)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / perPage.value)))

const pagedEquipment = computed(() => {
  const start = (page.value - 1) * perPage.value
  const end = start + perPage.value
  return filteredEquipment.value.slice(start, end)
})

const groupedEquipment = computed(() => {
  const groups = new Map()

  for (const item of pagedEquipment.value) {
    const setID = item.equipment_set?.equipment_set_id || item.equipment_set_id || 0
    const setName = item.equipment_set?.equipment_set_name || 'Без комплекта'

    if (!groups.has(setID)) {
      groups.set(setID, {
        equipment_set_id: setID,
        equipment_set_name: setName,
        items: []
      })
    }

    groups.get(setID).items.push(item)
  }

  return Array.from(groups.values())
})

const from = computed(() => {
  if (!total.value) return 0
  return (page.value - 1) * perPage.value + 1
})

const to = computed(() => {
  if (!total.value) return 0
  return Math.min(page.value * perPage.value, total.value)
})

watch(search, () => {
  page.value = 1
})

watch(perPage, () => {
  page.value = 1
})

watch(totalPages, (value) => {
  if (page.value > value) {
    page.value = value
  }
})

function prevPage() {
  if (page.value > 1) {
    page.value -= 1
  }
}

function nextPage() {
  if (page.value < totalPages.value) {
    page.value += 1
  }
}

async function addEquipment() {
  if (!selectedEquipmentId.value) return

  const equipment = (board.value?.available_equipment || []).find((item) => item.equipment_id === selectedEquipmentId.value)
  if (!equipment) return

  await crm.addEquipmentToDraft({
    draft_id: draftId.value,
    equipment_id: equipment.equipment_id,
    equipment_set_id: equipment.equipment_set_id
  })

  selectedEquipmentId.value = 0
}

async function addEquipmentByItem(item) {
  await crm.addEquipmentToDraft({
    draft_id: draftId.value,
    equipment_id: item.equipment_id,
    equipment_set_id: item.equipment_set_id
  })
}

async function removeEquipment(item) {
  await crm.removeEquipmentFromDraft({
    draft_id: draftId.value,
    equipment_id: item.equipment_id
  })
}

async function addSet() {
  if (!selectedSetId.value) return

  await crm.addSetToDraft({
    draft_id: draftId.value,
    equipment_set_id: selectedSetId.value
  })

  selectedSetId.value = 0
}
</script>

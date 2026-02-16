<template>
  <div class="space-y-6">
    <UCard v-if="board?.project">
      <template #header>
        <h1 class="text-xl font-semibold">Состав съёмки: {{ board.project.project_name }}</h1>
      </template>
      <p class="text-sm text-gray-600">
        {{ board.project.shooting_start_date }} - {{ board.project.shooting_end_date }}
      </p>
    </UCard>

    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold">Добавить в съёмку</h2>
      </template>

      <div class="grid gap-4 md:grid-cols-3">
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
            <option v-for="item in board?.sets_in_project || []" :key="item.equipment_set_id" :value="item.equipment_set_id">
              {{ item.equipment_set_name }}
            </option>
          </select>
          <UButton color="primary" variant="soft" @click="addSet">Добавить комплект</UButton>
        </div>

        <div class="space-y-2">
          <label class="text-sm text-gray-600">Шаблон</label>
          <select v-model.number="selectedDraftId" class="w-full rounded border border-gray-300 px-3 py-2">
            <option :value="0">Выберите шаблон</option>
            <option v-for="item in crm.drafts" :key="item.draft_id" :value="item.draft_id">
              {{ item.draft_name }}
            </option>
          </select>
          <UButton color="primary" variant="soft" @click="addDraft">Добавить шаблон</UButton>
        </div>
      </div>

      <div class="mt-4 flex gap-3">
        <UButton color="error" variant="soft" @click="resetEquipment">Очистить состав</UButton>
        <UButton color="neutral" variant="soft" @click="loadConflicts">Проверить конфликты</UButton>
      </div>
    </UCard>

    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold">Оборудование в съёмке</h2>
      </template>

      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left border-b border-gray-200">
              <th class="py-2">ID</th>
              <th class="py-2">Название</th>
              <th class="py-2">Серия</th>
              <th class="py-2">Комплект</th>
              <th class="py-2 w-24">Действие</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in board?.equipment_in_project || []" :key="item.equipment_id" class="border-b border-gray-100">
              <td class="py-2">{{ item.equipment_id }}</td>
              <td class="py-2">{{ item.equipment_name }}</td>
              <td class="py-2">{{ item.serial_number }}</td>
              <td class="py-2">{{ item.equipment_set?.equipment_set_name || '-' }}</td>
              <td class="py-2">
                <UButton size="xs" color="error" variant="soft" @click="removeEquipment(item)">Удалить</UButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>

    <UCard v-if="conflictingEquipment.length">
      <template #header>
        <h2 class="text-lg font-semibold">Конфликтующее оборудование</h2>
      </template>
      <ul class="space-y-1 text-sm">
        <li v-for="item in conflictingEquipment" :key="`${item.equipment_id}-${item.project_id}`">
          {{ item.equipment_name }} / {{ item.equipment_set_name }} -> проект #{{ item.project_id }} {{ item.project_name }}
        </li>
      </ul>
    </UCard>

    <UCard v-if="conflictingProjects.length">
      <template #header>
        <h2 class="text-lg font-semibold">Конфликтующие проекты</h2>
      </template>
      <ul class="space-y-1 text-sm">
        <li v-for="item in conflictingProjects" :key="item.project_id">
          #{{ item.project_id }} {{ item.project_name }} ({{ item.conflicting_equipment_count }} конфликтов)
        </li>
      </ul>
    </UCard>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute } from '#imports'
import { useCRMStore } from '~/stores/crm'

const route = useRoute()
const crm = useCRMStore()
const projectId = computed(() => Number(route.params.id))

const selectedEquipmentId = ref(0)
const selectedSetId = ref(0)
const selectedDraftId = ref(0)
const conflictingEquipment = ref([])
const conflictingProjects = ref([])

await Promise.all([
  crm.fetchProjectBoard(projectId.value),
  crm.fetchDrafts()
])

const board = computed(() => crm.projectBoard)

async function refreshBoard() {
  await crm.fetchProjectBoard(projectId.value)
}

async function addEquipment() {
  if (!selectedEquipmentId.value) return

  const equipment = (board.value?.available_equipment || []).find((item) => item.equipment_id === selectedEquipmentId.value)
  if (!equipment) return

  await crm.addEquipmentToProject({
    project_id: projectId.value,
    equipment_id: equipment.equipment_id,
    equipment_set_id: equipment.equipment_set_id
  })

  selectedEquipmentId.value = 0
}

async function removeEquipment(item) {
  await crm.removeEquipmentFromProject({
    project_id: projectId.value,
    equipment_id: item.equipment_id
  })
}

async function addSet() {
  if (!selectedSetId.value) return

  await crm.addSetToProject({
    project_id: projectId.value,
    equipment_set_id: selectedSetId.value
  })

  selectedSetId.value = 0
}

async function addDraft() {
  if (!selectedDraftId.value) return

  await crm.addDraftToProject({
    project_id: projectId.value,
    draft_id: selectedDraftId.value
  })

  selectedDraftId.value = 0
}

async function resetEquipment() {
  await crm.resetEquipmentInProject(projectId.value)
  await refreshBoard()
}

async function loadConflicts() {
  conflictingEquipment.value = await crm.fetchConflictingEquipment(projectId.value)
  conflictingProjects.value = await crm.fetchConflictingProjects()
}
</script>

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

      <div class="overflow-x-auto">
        <table class="w-full min-w-max text-sm">
          <thead>
            <tr class="text-left border-b border-gray-200 whitespace-nowrap">
              <th class="py-2">ID</th>
              <th class="py-2">Название</th>
              <th class="py-2">Серия</th>
              <th class="py-2">Комплект</th>
              <th class="py-2 w-24">Действие</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in board?.equipment_in_draft || []" :key="item.equipment_id" class="border-b border-gray-100">
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
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute } from '#imports'
import { useCRMStore } from '~/stores/crm'

const route = useRoute()
const crm = useCRMStore()
const draftId = computed(() => Number(route.params.id))

const selectedEquipmentId = ref(0)
const selectedSetId = ref(0)

await crm.fetchDraftBoard(draftId.value)

const board = computed(() => crm.draftBoard)

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

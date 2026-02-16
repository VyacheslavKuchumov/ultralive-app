<template>
  <section class="space-y-6">
    <UCard>
      <template #header>
        <div class="space-y-1">
          <h1 class="text-xl font-semibold">Ultralive CRM</h1>
          <p class="text-sm text-gray-600">Учет оборудования и комплектации съёмок.</p>
        </div>
      </template>

      <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
        <UCard>
          <p class="text-sm text-gray-600">Съёмки</p>
          <p class="text-2xl font-semibold">{{ crm.projects.length }}</p>
        </UCard>
        <UCard>
          <p class="text-sm text-gray-600">Шаблоны</p>
          <p class="text-2xl font-semibold">{{ crm.drafts.length }}</p>
        </UCard>
        <UCard>
          <p class="text-sm text-gray-600">Комплекты</p>
          <p class="text-2xl font-semibold">{{ crm.equipmentSets.length }}</p>
        </UCard>
        <UCard>
          <p class="text-sm text-gray-600">Оборудование</p>
          <p class="text-2xl font-semibold">{{ crm.equipment.length }}</p>
        </UCard>
      </div>
    </UCard>

    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      <UButton to="/projects" color="primary" variant="soft" class="justify-center">Съёмки</UButton>
      <UButton to="/projects/archived" color="primary" variant="soft" class="justify-center">Архив съёмок</UButton>
      <UButton to="/drafts" color="primary" variant="soft" class="justify-center">Шаблоны</UButton>
      <UButton to="/equipment_sets" color="primary" variant="soft" class="justify-center">Комплекты оборудования</UButton>
      <UButton to="/equipment" color="primary" variant="soft" class="justify-center">Оборудование</UButton>
      <UButton to="/set_types" color="primary" variant="soft" class="justify-center">Виды комплектов</UButton>
      <UButton to="/warehouses" color="primary" variant="soft" class="justify-center">Склады</UButton>
      <UButton to="/project_types" color="primary" variant="soft" class="justify-center">Площадки</UButton>
    </div>
  </section>
</template>

<script setup>
const crm = useCRMStore()

onMounted(async () => {
  await Promise.all([
    crm.fetchProjects(),
    crm.fetchDrafts(),
    crm.fetchEquipmentSets(),
    crm.fetchEquipment()
  ])
})
</script>

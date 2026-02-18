<template>
  <header class="border-b border-(--ui-border)">
    <UContainer class="py-3 flex items-center justify-between gap-4">
      <NuxtLink to="/" class="font-semibold text-lg primary">
        Ultralive CRM
      </NuxtLink>

      <div v-if="auth.isAuthenticated" class="hidden md:flex items-center gap-2">
        <UButton size="sm" color="neutral" variant="ghost" to="/projects">Съёмки</UButton>
        <UButton size="sm" color="neutral" variant="ghost" to="/projects/archived">Архив</UButton>
        <UButton size="sm" color="neutral" variant="ghost" to="/drafts">Шаблоны</UButton>
        <UButton size="sm" color="neutral" variant="ghost" to="/equipment_sets">Комплекты</UButton>
        <UButton size="sm" color="neutral" variant="ghost" to="/set_types">Виды комплектов</UButton>
        <UButton size="sm" color="neutral" variant="ghost" to="/warehouses">Склады</UButton>
        <UButton size="sm" color="neutral" variant="ghost" to="/project_types">Площадки</UButton>
      </div>

      <div v-if="auth.isAuthenticated" class="flex items-center gap-3">
        <UButton
          size="sm"
          color="neutral"
          variant="ghost"
          class="md:hidden"
          :icon="menuOpen ? 'i-lucide-x' : 'i-lucide-menu'"
          @click="menuOpen = !menuOpen"
        />
        <UButton size="sm" color="neutral" variant="ghost" to="/profile" class="hidden md:inline-flex">Профиль</UButton>
        <span class="text-sm text-gray-600 hidden md:inline">{{ auth.displayName }}</span>
        <UButton size="sm" color="error" variant="soft" class="hidden md:inline-flex" @click="handleLogout">Выйти</UButton>
      </div>
    </UContainer>

    <div v-if="auth.isAuthenticated && menuOpen" class="md:hidden border-t border-(--ui-border)">
      <UContainer class="py-3 grid gap-2">
        <UButton color="neutral" variant="ghost" to="/projects" class="justify-start" @click="menuOpen = false">Съёмки</UButton>
        <UButton color="neutral" variant="ghost" to="/projects/archived" class="justify-start" @click="menuOpen = false">Архив</UButton>
        <UButton color="neutral" variant="ghost" to="/drafts" class="justify-start" @click="menuOpen = false">Шаблоны</UButton>
        <UButton color="neutral" variant="ghost" to="/equipment_sets" class="justify-start" @click="menuOpen = false">Комплекты</UButton>
        <UButton color="neutral" variant="ghost" to="/set_types" class="justify-start" @click="menuOpen = false">Виды комплектов</UButton>
        <UButton color="neutral" variant="ghost" to="/warehouses" class="justify-start" @click="menuOpen = false">Склады</UButton>
        <UButton color="neutral" variant="ghost" to="/project_types" class="justify-start" @click="menuOpen = false">Площадки</UButton>
        <UButton color="neutral" variant="ghost" to="/profile" class="justify-start" @click="menuOpen = false">Профиль</UButton>
        <UButton color="error" variant="soft" class="justify-start" @click="handleLogout">Выйти</UButton>
      </UContainer>
    </div>
  </header>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useAuthStore } from '~/stores/auth'

const auth = useAuthStore()
const route = useRoute()
const menuOpen = ref(false)

function handleLogout() {
  menuOpen.value = false
  auth.logout()
}

watch(
  () => route.fullPath,
  () => {
    menuOpen.value = false
  }
)
</script>

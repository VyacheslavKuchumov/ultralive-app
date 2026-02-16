<template>
  <section class="space-y-4">
    <UCard>
      <template #header>
        <h1 class="text-xl font-semibold">Профиль</h1>
      </template>

      <UForm :state="profileState" class="space-y-4" @submit="onProfileSubmit">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <UFormField label="Имя" required>
            <UInput v-model="profileState.firstName" class="w-full" placeholder="Имя" />
          </UFormField>

          <UFormField label="Фамилия" required>
            <UInput v-model="profileState.lastName" class="w-full" placeholder="Фамилия" />
          </UFormField>
        </div>

        <UFormField label="Email" required>
          <UInput v-model="profileState.email" type="email" class="w-full" placeholder="you@example.com" />
        </UFormField>

        <UFormField label="Подтверждение текущим паролем" required>
          <UInput
            v-model="profileState.currentPassword"
            type="password"
            class="w-full"
            placeholder="Введите текущий пароль"
          />
        </UFormField>

        <UFormField label="Роль">
          <UInput :model-value="auth.profile?.role || '-'" class="w-full" disabled />
        </UFormField>

        <p v-if="profileError" class="text-sm text-red-600">{{ profileError }}</p>

        <div class="flex flex-wrap gap-2">
          <UButton type="submit" color="primary" size="sm" :loading="profileLoading">Сохранить профиль</UButton>
          <UButton color="neutral" variant="soft" size="sm" :loading="profileLoading" @click="reloadProfile">Обновить</UButton>
          <UButton color="error" variant="soft" size="sm" @click="auth.logout()">Выйти</UButton>
        </div>
      </UForm>
    </UCard>

    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold">Смена пароля</h2>
      </template>

      <UForm :state="passwordState" class="space-y-4" @submit="onPasswordSubmit">
        <UFormField label="Текущий пароль" required>
          <UInput v-model="passwordState.currentPassword" type="password" class="w-full" placeholder="Текущий пароль" />
        </UFormField>

        <UFormField label="Новый пароль" required>
          <UInput v-model="passwordState.newPassword" type="password" class="w-full" placeholder="Минимум 3 символа" />
        </UFormField>

        <UFormField label="Повтор нового пароля" required>
          <UInput v-model="passwordState.confirmPassword" type="password" class="w-full" placeholder="Повторите новый пароль" />
        </UFormField>

        <p v-if="passwordError" class="text-sm text-red-600">{{ passwordError }}</p>

        <UButton type="submit" color="primary" size="sm" :loading="passwordLoading">
          Обновить пароль
        </UButton>
      </UForm>
    </UCard>
  </section>
</template>

<script setup>
const toast = useToast()
const auth = useAuthStore()
const profileLoading = ref(false)
const passwordLoading = ref(false)
const profileError = ref('')
const passwordError = ref('')

const profileState = reactive({
  firstName: '',
  lastName: '',
  email: '',
  currentPassword: ''
})

const passwordState = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

function extractErrorMessage(error, fallback) {
  return error?.data?.statusMessage || error?.data?.message || error?.message || fallback
}

function syncProfileState(profile) {
  profileState.firstName = profile?.firstName || ''
  profileState.lastName = profile?.lastName || ''
  profileState.email = profile?.email || ''
  profileState.currentPassword = ''
}

async function reloadProfile() {
  profileLoading.value = true
  profileError.value = ''

  try {
    const profile = await auth.fetchProfile()
    syncProfileState(profile)
  } catch (error) {
    profileError.value = extractErrorMessage(error, 'Не удалось загрузить профиль')
    toast.add({
      title: 'Ошибка профиля',
      description: profileError.value,
      color: 'error'
    })
  } finally {
    profileLoading.value = false
  }
}

async function onProfileSubmit() {
  profileLoading.value = true
  profileError.value = ''

  try {
    const payload = {
      firstName: profileState.firstName.trim(),
      lastName: profileState.lastName.trim(),
      email: profileState.email.trim().toLowerCase(),
      currentPassword: profileState.currentPassword
    }

    if (!payload.firstName || !payload.lastName || !payload.email || !payload.currentPassword) {
      throw new Error('Заполните имя, фамилию, email и текущий пароль')
    }

    const profile = await auth.updateProfile(payload)
    syncProfileState(profile)

    toast.add({
      title: 'Профиль обновлён',
      color: 'success'
    })
  } catch (error) {
    profileError.value = extractErrorMessage(error, 'Не удалось обновить профиль')
    toast.add({
      title: 'Ошибка обновления профиля',
      description: profileError.value,
      color: 'error'
    })
  } finally {
    profileLoading.value = false
  }
}

async function onPasswordSubmit() {
  passwordLoading.value = true
  passwordError.value = ''

  try {
    const payload = {
      currentPassword: passwordState.currentPassword,
      newPassword: passwordState.newPassword
    }

    if (!payload.currentPassword || !payload.newPassword) {
      throw new Error('Заполните все поля для смены пароля')
    }
    if (payload.newPassword.length < 3) {
      throw new Error('Новый пароль должен содержать минимум 3 символа')
    }
    if (passwordState.newPassword !== passwordState.confirmPassword) {
      throw new Error('Повтор нового пароля не совпадает')
    }

    await auth.updatePassword(payload)

    passwordState.currentPassword = ''
    passwordState.newPassword = ''
    passwordState.confirmPassword = ''

    toast.add({
      title: 'Пароль обновлён',
      color: 'success'
    })
  } catch (error) {
    passwordError.value = extractErrorMessage(error, 'Не удалось обновить пароль')
    toast.add({
      title: 'Ошибка смены пароля',
      description: passwordError.value,
      color: 'error'
    })
  } finally {
    passwordLoading.value = false
  }
}

onMounted(async () => {
  if (auth.profile) syncProfileState(auth.profile)
  if (!auth.profile) await reloadProfile()
})
</script>

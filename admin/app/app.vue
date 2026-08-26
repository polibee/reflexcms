<script setup lang="ts">
const ui = useUiStore()

useHead({
  htmlAttrs: {
    class: () => (ui.theme === 'dark' ? 'dark' : '')
  },
  titleTemplate: title => (title ? `${title} · Nuxt Admin` : 'Nuxt Admin')
})
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
const mobileOpen = ref(false)
const route = useRoute()
watch(() => route.path, () => {
  mobileOpen.value = false
})
</script>

<template>
  <div class="flex min-h-screen">
    <!-- desktop sidebar -->
    <div class="sticky top-0 hidden h-screen lg:block">
      <AdminSidebar />
    </div>

    <!-- mobile drawer -->
    <Teleport to="body">
      <template v-if="mobileOpen">
        <div
          class="fixed inset-0 z-50 bg-black/60 lg:hidden"
          @click="mobileOpen = false"
        />
        <div class="fixed inset-y-0 left-0 z-50 lg:hidden">
          <AdminSidebar />
        </div>
      </template>
    </Teleport>

    <div class="flex min-w-0 flex-1 flex-col">
      <AdminHeader @toggle-sidebar="mobileOpen = !mobileOpen" />
      <main class="flex-1 p-4 md:p-6 lg:p-8">
        <slot />
      </main>
    </div>

    <ActionHost />
    <ToastHost />
  </div>
</template>

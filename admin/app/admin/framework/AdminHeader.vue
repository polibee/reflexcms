<script setup lang="ts">
import { LogOutIcon, MenuIcon, MoonIcon, SunIcon } from 'lucide-vue-next'
import {
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuRoot,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from 'reka-ui'

defineEmits<{ 'toggle-sidebar': [] }>()

const ui = useUiStore()
const auth = useAuthStore()

const roleBadge: Record<string, string> = {
  admin: 'Full access',
  editor: 'Content editor',
  viewer: 'Read only'
}
</script>

<template>
  <header class="sticky top-0 z-40 flex h-14 items-center gap-3 border-b bg-background/95 px-4 backdrop-blur md:px-6">
    <UiButton
      variant="ghost"
      size="icon"
      class="lg:hidden"
      @click="$emit('toggle-sidebar')"
    >
      <MenuIcon class="h-5 w-5" />
    </UiButton>

    <div class="min-w-0 flex-1" />

    <!-- theme toggle -->
    <UiButton
      variant="ghost"
      size="icon"
      aria-label="Toggle theme"
      @click="ui.toggleTheme()"
    >
      <SunIcon
        v-if="ui.theme === 'dark'"
        class="h-4 w-4"
      />
      <MoonIcon
        v-else
        class="h-4 w-4"
      />
    </UiButton>

    <!-- user menu -->
    <DropdownMenuRoot v-if="auth.user">
      <DropdownMenuTrigger as-child>
        <button
          class="flex h-9 items-center gap-2 rounded-full border pl-1 pr-3 text-sm transition-colors hover:bg-accent"
        >
          <span class="flex h-7 w-7 items-center justify-center rounded-full bg-primary text-xs font-semibold text-primary-foreground">
            {{ auth.user.name.split(' ').map(p => p[0]).slice(0, 2).join('') }}
          </span>
          <span class="hidden sm:inline">{{ auth.user.name }}</span>
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuPortal>
        <DropdownMenuContent
          align="end"
          class="z-50 w-56 rounded-md border bg-popover p-1 shadow-md"
        >
          <DropdownMenuLabel class="px-2 py-1.5">
            <p class="text-sm font-medium">
              {{ auth.user.name }}
            </p>
            <p class="text-xs text-muted-foreground">
              {{ auth.user.email }}
            </p>
          </DropdownMenuLabel>
          <DropdownMenuSeparator class="my-1 h-px bg-border" />
          <div class="px-2 py-1 text-xs text-muted-foreground">
            {{ roleBadge[auth.user.role] ?? auth.user.role }}
          </div>
          <DropdownMenuSeparator class="my-1 h-px bg-border" />
          <DropdownMenuItem
            class="flex cursor-pointer items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent"
            @select="auth.logout()"
          >
            <LogOutIcon class="h-4 w-4" /> Sign out
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenuPortal>
    </DropdownMenuRoot>
  </header>
</template>

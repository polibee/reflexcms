import {
  ActivityIcon,
  BadgeCheckIcon,
  BellIcon,
  BoxIcon,
  CircleIcon,
  ClipboardListIcon,
  CreditCardIcon,
  FileTextIcon,
  FilesIcon,
  KeyRoundIcon,
  LayoutDashboardIcon,
  PackageIcon,
  PencilIcon,
  PlusIcon,
  SettingsIcon,
  ShieldIcon,
  ShoppingCartIcon,
  StarIcon,
  TagIcon,
  Trash2Icon,
  UserIcon,
  UsersIcon,
  EyeIcon
} from 'lucide-vue-next'
import type { Component } from 'vue'

const map: Record<string, Component> = {
  'activity': ActivityIcon,
  'badge-check': BadgeCheckIcon,
  'bell': BellIcon,
  'box': BoxIcon,
  'circle': CircleIcon,
  'clipboard': ClipboardListIcon,
  'credit-card': CreditCardIcon,
  'dashboard': LayoutDashboardIcon,
  'eye': EyeIcon,
  'file-text': FileTextIcon,
  'files': FilesIcon,
  'key': KeyRoundIcon,
  'package': PackageIcon,
  'pencil': PencilIcon,
  'plus': PlusIcon,
  'settings': SettingsIcon,
  'shield': ShieldIcon,
  'shopping-cart': ShoppingCartIcon,
  'star': StarIcon,
  'tag': TagIcon,
  'trash': Trash2Icon,
  'user': UserIcon,
  'users': UsersIcon
}

/** resolve a lucide icon by framework key; falls back to a plain circle */
export function getIcon(name?: string): Component {
  return (name && map[name]) || CircleIcon
}

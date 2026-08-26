import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: Array<string | false | Record<string, boolean> | undefined | null>): string {
  const classes: string[] = []
  for (const input of inputs) {
    if (!input) continue
    if (typeof input === 'string') classes.push(input)
    else {
      for (const [key, value] of Object.entries(input)) if (value) classes.push(key)
    }
  }
  return twMerge(clsx(classes))
}

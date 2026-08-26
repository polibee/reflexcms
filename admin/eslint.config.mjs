// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs'

export default withNuxt(
  // reference materials shipped under docs/ are not part of this project
  {
    ignores: ['docs/**', '.output/**']
  }
)

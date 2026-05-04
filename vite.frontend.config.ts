// vite.frontend.config.ts
//
// Vite configuration for building the standalone SPA frontend.
// This produces a self-contained static site in frontend-dist/ that can be
// deployed to GitHub Pages (or any static host) independently of the Gitea
// Go binary.
//
// Usage:
//   pnpm exec vite build --config vite.frontend.config.ts
//
// With a custom default API server:
//   VITE_DEFAULT_API_URL=https://gitea.example.com \
//     pnpm exec vite build --config vite.frontend.config.ts

import {defineConfig} from 'vite';
import vuePlugin from '@vitejs/plugin-vue';
import {join} from 'node:path';
import tailwindcss from 'tailwindcss';
import tailwindConfig from './tailwind.config.ts';

export default defineConfig({
  root: import.meta.dirname,
  // './' makes asset URLs relative so the site works in any subdirectory
  // (e.g. https://username.github.io/repo-name/).
  base: './',
  plugins: [
    vuePlugin({
      template: {
        compilerOptions: {
          // Keep parity with the main Vite config web-components list.
          isCustomElement: (tag) =>
            ['overflow-menu', 'relative-time', 'markdown-toolbar', 'text-expander'].includes(tag),
        },
      },
    }),
  ],
  css: {
    transformer: 'postcss',
    postcss: {plugins: [tailwindcss(tailwindConfig)]},
  },
  define: {
    __VUE_OPTIONS_API__: true,
    __VUE_PROD_DEVTOOLS__: false,
    __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
  },
  build: {
    outDir: join(import.meta.dirname, 'frontend-dist'),
    emptyOutDir: true,
    rollupOptions: {
      input: {
        index: join(import.meta.dirname, 'frontend/index.html'),
      },
    },
  },
});

// spaconfig.ts — centralised configuration for the SPA.
//
// Works in two modes:
//   Embedded: SPA served by the Gitea Go binary.  window.config is populated
//             by the Go template before the JS runs.
//   Standalone: SPA deployed to GitHub Pages (or any static host).
//               VITE_DEFAULT_API_URL is set at build time to point to the
//               Gitea server; window.config is initialised to empty defaults
//               by the standalone entry HTML.

// Build-time API URL injected via VITE_DEFAULT_API_URL env-var.
// Empty string when building in embedded (normal Gitea) mode.
const buildApiUrl: string = import.meta.env['VITE_DEFAULT_API_URL'] ?? '';

// Safely read window.config which may be absent in standalone mode.
function wConfig() {
  return (window as any).config as typeof window.config | undefined;
}

/** True when the SPA was built for standalone GitHub Pages deployment. */
export const isStandalone: boolean = !!buildApiUrl;

/**
 * Base URL of the Gitea instance, WITHOUT a trailing slash.
 * Examples:
 *   embedded  →  '' (same origin)
 *   standalone →  'https://gitea.example.com'
 */
export const appSubUrl: string = buildApiUrl
  ? buildApiUrl.replace(/\/$/, '')
  : (wConfig()?.appSubUrl ?? '');

/**
 * Base URL for all Gitea REST-API v1 calls.
 * Examples:
 *   embedded  →  '/api/v1'
 *   standalone →  'https://gitea.example.com/api/v1'
 */
export const apiBase: string = `${appSubUrl}/api/v1`;

/**
 * URL prefix for bundled static assets (logo, etc.).
 * In embedded mode this comes from window.config; in standalone mode
 * assets are bundled into the frontend itself so the prefix is empty
 * (relative URLs work).
 */
export const assetUrlPrefix: string = buildApiUrl
  ? ''
  : (wConfig()?.assetUrlPrefix ?? '');

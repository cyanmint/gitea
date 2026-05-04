import {createApp} from 'vue';
import {router} from './router/index.ts';
import App from './App.vue';

// Initialise window.config with safe defaults when running in standalone mode
// (GitHub Pages deployment) where the Go template does not inject config.
if (!(window as any).config) {
  (window as any).config = {
    appUrl: import.meta.env['VITE_DEFAULT_API_URL'] ?? '',
    appSubUrl: '',
    assetUrlPrefix: '',
    sharedWorkerUri: '',
    runModeIsProd: true,
    customEmojis: {},
    pageData: {},
    notificationSettings: {MinTimeout: 10000, TimeoutStep: 10000, MaxTimeout: 300000, EventSourceUpdateTime: 10000},
    enableTimeTracking: false,
    mermaidMaxSourceCharacters: 5000,
    i18n: {},
    frontendInited: false,
  };
}

const app = createApp(App);
app.use(router);
app.mount('#app');

window.config.frontendInited = true;

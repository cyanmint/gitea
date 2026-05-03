import {createRouter, createWebHistory, type RouteRecordRaw} from 'vue-router';
import HomePage from '../pages/HomePage.vue';
import ExplorePage from '../pages/ExplorePage.vue';
import LoginPage from '../pages/LoginPage.vue';
import RepoOverviewPage from '../pages/RepoOverviewPage.vue';
import NotFoundPage from '../pages/NotFoundPage.vue';

const {appSubUrl} = window.config;

const routes: RouteRecordRaw[] = [
  // Home / dashboard
  {path: '/', component: HomePage, meta: {title: 'Dashboard'}},

  // Explore
  {path: '/explore', redirect: '/explore/repos'},
  {path: '/explore/repos', component: ExplorePage, meta: {title: 'Explore Repositories', tab: 'repos'}},
  {path: '/explore/users', component: ExplorePage, meta: {title: 'Explore Users', tab: 'users'}},
  {path: '/explore/organizations', component: ExplorePage, meta: {title: 'Explore Organizations', tab: 'orgs'}},

  // Auth
  {path: '/user/login', component: LoginPage, meta: {title: 'Sign In', public: true}},

  // Repository routes — the SPA handles the top-level overview page.
  // Sub-paths (src/branch/…, issues/…, etc.) are still served by the
  // Go backend and will be gradually migrated to the SPA in future phases.
  {path: '/:owner/:repo', component: RepoOverviewPage, meta: {title: 'Repository'}},

  // 404 catch-all
  {path: '/:pathMatch(.*)*', component: NotFoundPage, meta: {title: 'Page Not Found', public: true}},
];

export const router = createRouter({
  history: createWebHistory(appSubUrl || '/'),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) return savedPosition;
    return {top: 0};
  },
});

// Capture the application name from the page's initial title (e.g. "Gitea")
// once at startup before any navigation modifies it.
const appName = document.title.split(' - ').pop() ?? 'Gitea';

// Update document title on route change
router.afterEach((to) => {
  const title = to.meta.title as string | undefined;
  if (title) {
    document.title = `${title} - ${appName}`;
  }
});

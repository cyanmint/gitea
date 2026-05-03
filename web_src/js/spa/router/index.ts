import {createRouter, createWebHistory, type RouteRecordRaw} from 'vue-router';
import HomePage from '../pages/HomePage.vue';
import ExplorePage from '../pages/ExplorePage.vue';
import LoginPage from '../pages/LoginPage.vue';
import RegisterPage from '../pages/RegisterPage.vue';
import UserProfilePage from '../pages/UserProfilePage.vue';
import RepoOverviewPage from '../pages/RepoOverviewPage.vue';
import IssueListPage from '../pages/IssueListPage.vue';
import IssueDetailPage from '../pages/IssueDetailPage.vue';
import PullRequestListPage from '../pages/PullRequestListPage.vue';
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
  {path: '/user/sign_up', component: RegisterPage, meta: {title: 'Register', public: true}},

  // Repository routes
  {path: '/:owner/:repo', component: RepoOverviewPage, meta: {title: 'Repository'}},
  {path: '/:owner/:repo/issues', component: IssueListPage, meta: {title: 'Issues'}},
  {path: '/:owner/:repo/issues/:id', component: IssueDetailPage, meta: {title: 'Issue'}},
  {path: '/:owner/:repo/pulls', component: PullRequestListPage, meta: {title: 'Pull Requests'}},

  // User / org profile — comes after repo routes so /:owner/:repo wins
  {path: '/:username', component: UserProfilePage, meta: {title: 'Profile', public: true}},

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

// Capture the application name once at startup. The initial page title has the
// form "Page Title - AppName"; if there's no separator, the whole string is used.
const titleParts = document.title.split(' - ');
const appName = titleParts.length > 1 ? titleParts.at(-1)! : titleParts[0] || 'Gitea';

// Update document title on route change
router.afterEach((to) => {
  const title = to.meta.title as string | undefined;
  if (title) {
    document.title = `${title} - ${appName}`;
  }
});

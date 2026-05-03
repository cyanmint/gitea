import {createRouter, createWebHistory, type RouteRecordRaw} from 'vue-router';
import HomePage from '../pages/HomePage.vue';
import ExplorePage from '../pages/ExplorePage.vue';
import LoginPage from '../pages/LoginPage.vue';
import RegisterPage from '../pages/RegisterPage.vue';
import UserProfilePage from '../pages/UserProfilePage.vue';
import UserIssuesPage from '../pages/UserIssuesPage.vue';
import UserSettingsPage from '../pages/UserSettingsPage.vue';
import NotificationsPage from '../pages/NotificationsPage.vue';
import RepoOverviewPage from '../pages/RepoOverviewPage.vue';
import IssueListPage from '../pages/IssueListPage.vue';
import IssueDetailPage from '../pages/IssueDetailPage.vue';
import PullRequestListPage from '../pages/PullRequestListPage.vue';
import RepoSourcePage from '../pages/RepoSourcePage.vue';
import RepoCommitsPage from '../pages/RepoCommitsPage.vue';
import RepoBranchesPage from '../pages/RepoBranchesPage.vue';
import RepoReleasesPage from '../pages/RepoReleasesPage.vue';
import RepoTagsPage from '../pages/RepoTagsPage.vue';
import RepoWikiPage from '../pages/RepoWikiPage.vue';
import RepoActivityPage from '../pages/RepoActivityPage.vue';
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

  // User settings
  {path: '/user/settings', component: UserSettingsPage, meta: {title: 'Settings'}},
  {path: '/user/settings/:tab', component: UserSettingsPage, meta: {title: 'Settings'}},

  // User issues/pulls/milestones dashboard
  {path: '/issues', component: UserIssuesPage, meta: {title: 'Issues'}},
  {path: '/pulls', component: UserIssuesPage, meta: {title: 'Pull Requests'}},
  {path: '/milestones', component: UserIssuesPage, meta: {title: 'Milestones'}},

  // Notifications
  {path: '/notifications', component: NotificationsPage, meta: {title: 'Notifications'}},
  {path: '/notifications/subscriptions', component: NotificationsPage, meta: {title: 'Subscriptions'}},
  {path: '/notifications/watching', component: NotificationsPage, meta: {title: 'Watching'}},

  // Repository routes — must come before /:username
  {path: '/:owner/:repo', component: RepoOverviewPage, meta: {title: 'Repository'}},
  {path: '/:owner/:repo/issues', component: IssueListPage, meta: {title: 'Issues'}},
  {path: '/:owner/:repo/issues/:id', component: IssueDetailPage, meta: {title: 'Issue'}},
  {path: '/:owner/:repo/pulls', component: PullRequestListPage, meta: {title: 'Pull Requests'}},
  {path: '/:owner/:repo/pulls/:index', component: IssueDetailPage, meta: {title: 'Pull Request'}},

  // Repo source browser
  {path: '/:owner/:repo/src/:refType/:ref', component: RepoSourcePage, meta: {title: 'Source'}},
  {path: '/:owner/:repo/src/:refType/:ref/:pathMatch(.*)', component: RepoSourcePage, meta: {title: 'Source'}},

  // Repo commits
  {path: '/:owner/:repo/commits/:refType/:ref', component: RepoCommitsPage, meta: {title: 'Commits'}},

  // Repo branches
  {path: '/:owner/:repo/branches', component: RepoBranchesPage, meta: {title: 'Branches'}},

  // Repo releases
  {path: '/:owner/:repo/releases', component: RepoReleasesPage, meta: {title: 'Releases'}},
  {path: '/:owner/:repo/releases/tag/:tag', component: RepoReleasesPage, meta: {title: 'Release'}},
  {path: '/:owner/:repo/releases/latest', component: RepoReleasesPage, meta: {title: 'Latest Release'}},

  // Repo tags
  {path: '/:owner/:repo/tags', component: RepoTagsPage, meta: {title: 'Tags'}},

  // Repo wiki
  {path: '/:owner/:repo/wiki', component: RepoWikiPage, meta: {title: 'Wiki'}},
  {path: '/:owner/:repo/wiki/:pathMatch(.*)', component: RepoWikiPage, meta: {title: 'Wiki'}},

  // Repo activity
  {path: '/:owner/:repo/activity', component: RepoActivityPage, meta: {title: 'Activity'}},
  {path: '/:owner/:repo/activity/:period', component: RepoActivityPage, meta: {title: 'Activity'}},

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

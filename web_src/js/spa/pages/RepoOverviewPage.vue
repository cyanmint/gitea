<template>
  <AppLayout>
    <div v-if="loading" class="ui container tw-py-8">
      <div class="ui active centered inline loader"/>
    </div>

    <div v-else-if="error" class="ui container tw-py-6">
      <div class="ui negative message">
        <h3>Repository not found</h3>
        <p>{{ error }}</p>
      </div>
    </div>

    <template v-else-if="repo">
      <!-- ── Secondary nav (repo header + tabs) ──────────────────────── -->
      <div class="secondary-nav">
        <div class="ui container">
          <div class="repo-header flex-left-right">
            <!-- left: owner/repo name + badges -->
            <div class="flex-text-block">
              <div class="flex-text-block tw-flex-wrap tw-text-18">
                <RouterLink class="muted tw-font-normal" :to="`/${owner}`">{{ owner }}</RouterLink>/<RouterLink class="muted" :to="`/${owner}/${repoName}`">{{ repoName }}</RouterLink>
              </div>
              <div class="flex-text-block tw-flex-wrap">
                <span v-if="repo.archived" class="ui basic label not-mobile">Archived</span>
                <span v-if="repo.private" class="ui basic label not-mobile">Private</span>
                <span v-if="repo.fork" class="ui basic label not-mobile">Fork</span>
              </div>
            </div>

            <!-- right: star/fork buttons -->
            <div class="flex-text-block tw-flex-wrap">
              <!-- star button -->
              <div v-if="currentUser" class="ui buttons">
                <button
                  class="ui compact small basic button"
                  :class="{loading: starLoading}"
                  :disabled="starLoading"
                  @click="toggleStar"
                >
                  <SvgIcon name="octicon-star" :size="16" />
                  {{ starred ? 'Unstar' : 'Star' }}
                </button>
                <RouterLink :to="`/${owner}/${repoName}/stargazers`" class="ui compact small basic button">
                  {{ repo.stars_count }}
                </RouterLink>
              </div>
              <div v-else class="ui buttons">
                <RouterLink to="/user/login" class="ui compact small basic button" rel="nofollow">
                  <SvgIcon name="octicon-star" :size="16" /> Star
                </RouterLink>
                <RouterLink :to="`/${owner}/${repoName}/stargazers`" class="ui compact small basic button">
                  {{ repo.stars_count }}
                </RouterLink>
              </div>
              <!-- fork button -->
              <div class="ui buttons tw-ml-1">
                <RouterLink :to="`/${owner}/${repoName}`" class="ui compact small basic button">
                  <SvgIcon name="octicon-repo-forked" :size="16" /> Fork
                </RouterLink>
                <RouterLink :to="`/${owner}/${repoName}/forks`" class="ui compact small basic button">
                  {{ repo.forks_count }}
                </RouterLink>
              </div>
            </div>
          </div>

          <!-- fork/mirror info -->
          <div v-if="repo.fork && repo.parent" class="secondary-info">
            Forked from <RouterLink :to="`/${repo.parent.full_name}`">{{ repo.parent.full_name }}</RouterLink>
          </div>
        </div>

        <!-- tab navigation -->
        <div class="ui container">
          <overflow-menu class="ui secondary pointing menu">
            <div class="overflow-menu-items">
              <RouterLink :to="`/${owner}/${repoName}`" class="item" :class="{active: activeTab === 'code'}">
                <SvgIcon name="octicon-code" :size="16" /> Code
              </RouterLink>
              <RouterLink :to="`/${owner}/${repoName}/issues`" class="item" :class="{active: activeTab === 'issues'}">
                <SvgIcon name="octicon-issue-opened" :size="16" /> Issues
                <span v-if="repo.open_issues_count" class="ui small label">{{ repo.open_issues_count }}</span>
              </RouterLink>
              <RouterLink :to="`/${owner}/${repoName}/pulls`" class="item" :class="{active: activeTab === 'pulls'}">
                <SvgIcon name="octicon-git-pull-request" :size="16" /> Pull Requests
              </RouterLink>
              <RouterLink :to="`/${owner}/${repoName}/releases`" class="item" :class="{active: activeTab === 'releases'}">
                <SvgIcon name="octicon-tag" :size="16" /> Releases
              </RouterLink>
              <RouterLink :to="`/${owner}/${repoName}/wiki`" class="item" :class="{active: activeTab === 'wiki'}">
                <SvgIcon name="octicon-book" :size="16" /> Wiki
              </RouterLink>
              <RouterLink :to="`/${owner}/${repoName}/activity`" class="item" :class="{active: activeTab === 'activity'}">
                <SvgIcon name="octicon-pulse" :size="16" /> Activity
              </RouterLink>
              <template v-if="currentUser && (currentUser.login === owner || currentUser.is_admin)">
                <span class="item-flex-space"/>
                <RouterLink :to="`/${owner}/${repoName}/settings`" class="item" :class="{active: activeTab === 'settings'}">
                  <SvgIcon name="octicon-tools" :size="16" /> Settings
                </RouterLink>
              </template>
            </div>
          </overflow-menu>
        </div>
        <div class="ui tabs divider"/>
      </div>

      <!-- ── Main content ─────────────────────────────────────────────── -->
      <div class="ui container">
        <div class="repo-grid-filelist-sidebar">
          <!-- file list + README area -->
          <div class="repo-home-filelist">
            <!-- empty repo notice -->
            <div v-if="repo.empty" class="ui segment">
              <p>This repository is empty.</p>
              <div class="tw-mt-4 tw-font-mono tw-text-sm">
                <p>Quick setup — if you've done this kind of thing before</p>
                <pre class="tw-bg-gray-100 tw-rounded tw-p-3">git clone {{ httpCloneUrl }}</pre>
              </div>
            </div>

            <template v-else>
              <!-- latest commit bar -->
              <div id="repo-files-table">
                <div v-if="contentsLoading" class="repo-file-line tw-py-2">
                  <div class="ui active inline loader"/>
                </div>
                <div v-else-if="contentsError" class="repo-file-line tw-text-red-600 tw-text-sm tw-px-3 tw-py-2">
                  {{ contentsError }}
                </div>
                <template v-else>
                  <div
                    v-for="item in dirContents"
                    :key="item.path"
                    class="repo-file-item"
                  >
                    <div class="repo-file-cell name muted-links">
                      <SvgIcon :name="item.type === 'dir' ? 'octicon-file-directory-fill' : 'octicon-file'" :size="16" />
                      <RouterLink :to="contentItemUrl(item)" class="entry-name" :title="item.name">
                        {{ item.name }}
                      </RouterLink>
                    </div>
                    <div class="repo-file-cell message"/>
                    <div class="repo-file-cell age"/>
                  </div>
                </template>
              </div>
            </template>
          </div>

          <!-- right sidebar -->
          <div class="repo-home-sidebar-top flex-relaxed-list">
            <!-- description -->
            <div class="flex-relaxed-list">
              <div class="repo-home-sidebar-header">About</div>
              <div class="repo-description tw-break-anywhere">
                {{ repo.description || 'No description provided.' }}
              </div>
              <a v-if="repo.website" :href="repo.website" class="flex-text-block" target="_blank" rel="noopener">
                <SvgIcon name="octicon-link" :size="16" /> {{ repo.website }}
              </a>
            </div>

            <!-- stats -->
            <div class="flex-relaxed-list tw-mt-4">
              <RouterLink :to="`/${owner}/${repoName}/stargazers`" class="flex-text-block muted">
                <SvgIcon name="octicon-star" :size="16" />
                <strong>{{ repo.stars_count }}</strong> Stars
              </RouterLink>
              <RouterLink :to="`/${owner}/${repoName}/forks`" class="flex-text-block muted">
                <SvgIcon name="octicon-repo-forked" :size="16" />
                <strong>{{ repo.forks_count }}</strong> Forks
              </RouterLink>
              <RouterLink :to="`/${owner}/${repoName}/issues`" class="flex-text-block muted">
                <SvgIcon name="octicon-issue-opened" :size="16" />
                <strong>{{ repo.open_issues_count }}</strong> Open Issues
              </RouterLink>
              <RouterLink :to="`/${owner}/${repoName}/watchers`" class="flex-text-block muted">
                <SvgIcon name="octicon-eye" :size="16" />
                <strong>{{ repo.watchers_count }}</strong> Watchers
              </RouterLink>
            </div>

            <!-- clone panel -->
            <div class="tw-mt-4">
              <div class="ui small action input tw-flex tw-w-full">
                <input
                  type="text"
                  :value="httpCloneUrl"
                  readonly
                  class="tw-font-mono"
                  @focus="($event.target as HTMLInputElement).select()"
                >
                <button class="ui button" @click="copyUrl(httpCloneUrl)">Copy</button>
              </div>
              <div v-if="sshCloneUrl" class="ui small action input tw-flex tw-w-full tw-mt-2">
                <input
                  type="text"
                  :value="sshCloneUrl"
                  readonly
                  class="tw-font-mono"
                  @focus="($event.target as HTMLInputElement).select()"
                >
                <button class="ui button" @click="copyUrl(sshCloneUrl)">Copy</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </AppLayout>
</template>

<script setup lang="ts">
import {ref, computed, onMounted} from 'vue';
import {RouterLink, useRoute} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {SvgIcon} from '../../../js/svg.ts';
import {
  getRepo, getRepoContents, getCurrentUser,
  isRepoStarred, starRepo, unstarRepo,
  type Repository, type ContentsResponse, type User,
} from '../api/index.ts';
import {rewriteToBackend} from '../spaconfig.ts';

const route = useRoute();
const owner = String(route.params.owner);
const repoName = String(route.params.repo);

const loading = ref(true);
const error = ref('');
const repo = ref<Repository | null>(null);
const currentUser = ref<User | null>(null);

const contentsLoading = ref(false);
const contentsError = ref('');
const dirContents = ref<ContentsResponse[]>([]);

const starred = ref(false);
const starLoading = ref(false);

// Active tab is always 'code' on this page
const activeTab = 'code';

const httpCloneUrl = computed(() => rewriteToBackend(repo.value?.clone_url ?? ''));
const sshCloneUrl = computed(() => repo.value?.ssh_url ?? '');

function contentItemUrl(item: ContentsResponse): string {
  const branch = repo.value?.default_branch ?? 'HEAD';
  return `/${owner}/${repoName}/src/branch/${branch}/${item.path}`;
}

function copyUrl(url: string) {
  navigator.clipboard?.writeText(url).catch(() => {});
}

async function toggleStar() {
  if (!repo.value) return;
  starLoading.value = true;
  try {
    if (starred.value) {
      await unstarRepo(owner, repoName);
      repo.value.stars_count = Math.max(0, repo.value.stars_count - 1);
      starred.value = false;
    } else {
      await starRepo(owner, repoName);
      repo.value.stars_count += 1;
      starred.value = true;
    }
  } catch {
    // silently ignore
  } finally {
    starLoading.value = false;
  }
}

onMounted(async () => {
  currentUser.value = await getCurrentUser();

  try {
    repo.value = await getRepo(owner, repoName);
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Repository not found';
    loading.value = false;
    return;
  }
  loading.value = false;

  if (currentUser.value) {
    isRepoStarred(owner, repoName).then(v => { starred.value = v; }).catch(() => {});
  }

  contentsLoading.value = true;
  try {
    const contents = await getRepoContents(owner, repoName, '', repo.value!.default_branch);
    if (Array.isArray(contents)) {
      dirContents.value = (contents as ContentsResponse[]).sort((a, b) => {
        if (a.type === 'dir' && b.type !== 'dir') return -1;
        if (a.type !== 'dir' && b.type === 'dir') return 1;
        return a.name.localeCompare(b.name);
      });
    }
  } catch (err) {
    contentsError.value = err instanceof Error ? err.message : 'Failed to load files';
  } finally {
    contentsLoading.value = false;
  }
});
</script>
      <!-- Repo header -->
      <div class="tw-flex tw-items-center tw-gap-2 tw-mb-4">
        <RouterLink :to="`/${owner}`" class="tw-text-blue-600 hover:tw-underline tw-text-lg">{{ owner }}</RouterLink>
        <span class="tw-text-gray-400">/</span>
        <span class="tw-font-semibold tw-text-lg">{{ repoName }}</span>
        <span v-if="repo.private" class="ui mini label">Private</span>
        <span v-if="repo.archived" class="ui mini label tw-ml-1">Archived</span>
        <span v-if="repo.fork" class="ui mini label tw-ml-1">Fork</span>

        <!-- Star button (only for signed-in users) -->
        <div v-if="currentUser" class="tw-ml-auto">
          <button
            class="ui small basic button"
            :class="{loading: starLoading}"
            :disabled="starLoading"
            @click="toggleStar"
          >
            {{ starred ? '⭐ Unstar' : '☆ Star' }}
            <span class="ui label tw-ml-1">{{ repo.stars_count }}</span>
          </button>
        </div>
      </div>

      <p v-if="repo.description" class="tw-text-gray-600 tw-mb-4">{{ repo.description }}</p>

      <!-- Stats bar -->
      <div class="tw-flex tw-items-center tw-gap-4 tw-text-sm tw-mb-6">
        <RouterLink :to="`/${owner}/${repoName}/stargazers`" class="tw-flex tw-items-center tw-gap-1 hover:tw-text-blue-600">
          <span>⭐</span> <span class="tw-font-medium">{{ repo.stars_count }}</span> Stars
        </RouterLink>
        <RouterLink :to="`/${owner}/${repoName}/forks`" class="tw-flex tw-items-center tw-gap-1 hover:tw-text-blue-600">
          <span>🍴</span> <span class="tw-font-medium">{{ repo.forks_count }}</span> Forks
        </RouterLink>
        <RouterLink :to="`/${owner}/${repoName}/issues`" class="tw-flex tw-items-center tw-gap-1 hover:tw-text-blue-600">
          <span>🔴</span> <span class="tw-font-medium">{{ repo.open_issues_count }}</span> Issues
        </RouterLink>
      </div>

      <!-- Clone URL bar -->
      <div class="tw-mb-6">
        <div class="tw-flex tw-items-center tw-gap-2 tw-mb-2">
          <span class="tw-text-sm tw-font-medium tw-text-gray-600">Clone</span>
        </div>
        <div class="tw-flex tw-flex-col tw-gap-1">
          <div class="tw-flex tw-items-center tw-gap-2">
            <span class="tw-text-xs tw-text-gray-500 tw-w-8">HTTPS</span>
            <input
              type="text"
              :value="httpCloneUrl"
              readonly
              class="tw-flex-1 tw-border tw-rounded tw-px-3 tw-py-1 tw-text-sm tw-font-mono tw-bg-gray-50"
              @focus="($event.target as HTMLInputElement).select()"
            >
          </div>
          <div v-if="sshCloneUrl" class="tw-flex tw-items-center tw-gap-2">
            <span class="tw-text-xs tw-text-gray-500 tw-w-8">SSH</span>
            <input
              type="text"
              :value="sshCloneUrl"
              readonly
              class="tw-flex-1 tw-border tw-rounded tw-px-3 tw-py-1 tw-text-sm tw-font-mono tw-bg-gray-50"
              @focus="($event.target as HTMLInputElement).select()"
            >
          </div>
        </div>
      </div>

      <!-- File browser -->
      <div class="tw-border tw-rounded">
        <div class="tw-bg-gray-50 tw-px-4 tw-py-2 tw-border-b tw-flex tw-items-center tw-gap-2">
          <span class="tw-font-medium">{{ repo.default_branch }}</span>
          <span class="tw-text-gray-400">·</span>
          <span class="tw-text-sm tw-text-gray-600">{{ contentsCount }} files</span>
        </div>

        <div v-if="contentsLoading" class="tw-py-8 tw-text-center">
          <div class="ui active centered inline loader"/>
        </div>

        <div v-else-if="contentsError" class="tw-px-4 tw-py-4 tw-text-red-600 tw-text-sm">
          {{ contentsError }}
        </div>

        <div v-else>
          <div
            v-for="item in dirContents"
            :key="item.path"
            class="tw-flex tw-items-center tw-gap-2 tw-px-4 tw-py-2 tw-border-b last:tw-border-0 hover:tw-bg-gray-50"
          >
            <span class="tw-text-gray-500 tw-w-4">{{ item.type === 'dir' ? '📁' : '📄' }}</span>
            <a
              :href="contentItemUrl(item)"
              class="tw-font-medium tw-text-blue-600 hover:tw-underline"
            >
              {{ item.name }}
            </a>
          </div>
        </div>
      </div>

      <!-- Issues section -->
      <div class="tw-mt-8">
        <h3 class="tw-font-semibold tw-text-lg tw-mb-3">
          Recent Issues
          <RouterLink :to="`/${owner}/${repoName}/issues`" class="tw-text-blue-600 tw-text-sm tw-font-normal tw-ml-2">View all →</RouterLink>
          <RouterLink v-if="currentUser" :to="`/${owner}/${repoName}/issues/new`" class="tw-text-blue-600 tw-text-sm tw-font-normal tw-ml-2">+ New issue</RouterLink>
        </h3>
        <div v-if="issuesLoading" class="ui active centered inline loader"/>
        <div v-else-if="issues.length === 0" class="tw-text-gray-500 tw-text-sm">No open issues.</div>
        <div v-else class="tw-border tw-rounded">
          <div
            v-for="issue in issues"
            :key="issue.id"
            class="tw-px-4 tw-py-3 tw-border-b last:tw-border-0 hover:tw-bg-gray-50"
          >
            <RouterLink :to="`/${owner}/${repoName}/issues/${issue.number}`" class="tw-font-medium tw-text-blue-600 hover:tw-underline">
              #{{ issue.number }} {{ issue.title }}
            </RouterLink>
            <p class="tw-text-xs tw-text-gray-500 tw-mt-1">
              Opened by {{ issue.user.login }} · {{ issue.comments }} comments
            </p>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import {ref, computed, onMounted} from 'vue';
import {RouterLink, useRoute} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {
  getRepo, getRepoContents, getRepoIssues, getCurrentUser,
  isRepoStarred, starRepo, unstarRepo,
  type Repository, type Issue, type ContentsResponse, type User,
} from '../api/index.ts';

import {appSubUrl, rewriteToBackend} from '../spaconfig.ts';

const route = useRoute();
const owner = String(route.params.owner);
const repoName = String(route.params.repo);

const loading = ref(true);
const error = ref('');
const repo = ref<Repository | null>(null);
const currentUser = ref<User | null>(null);

const contentsLoading = ref(false);
const contentsError = ref('');
const dirContents = ref<ContentsResponse[]>([]);
const contentsCount = ref(0);

const issuesLoading = ref(false);
const issues = ref<Issue[]>([]);

const starred = ref(false);
const starLoading = ref(false);

/** HTTP clone URL rewritten to the configured backend origin. */
const httpCloneUrl = computed(() => rewriteToBackend(repo.value?.clone_url ?? ''));
/** SSH clone URL (always points to the backend; no rewriting needed). */
const sshCloneUrl = computed(() => repo.value?.ssh_url ?? '');

/** Returns the SPA router path for a repository content item (file or dir). */
function contentItemUrl(item: ContentsResponse): string {
  const branch = repo.value?.default_branch ?? 'HEAD';
  return `/${owner}/${repoName}/src/branch/${branch}/${item.path}`;
}

async function toggleStar() {
  if (!repo.value) return;
  starLoading.value = true;
  try {
    if (starred.value) {
      await unstarRepo(owner, repoName);
      repo.value.stars_count = Math.max(0, repo.value.stars_count - 1);
      starred.value = false;
    } else {
      await starRepo(owner, repoName);
      repo.value.stars_count += 1;
      starred.value = true;
    }
  } catch {
    // silently ignore
  } finally {
    starLoading.value = false;
  }
}

onMounted(async () => {
  currentUser.value = await getCurrentUser();

  try {
    repo.value = await getRepo(owner, repoName);
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Repository not found';
    loading.value = false;
    return;
  }
  loading.value = false;

  // Check star status for signed-in users
  if (currentUser.value) {
    isRepoStarred(owner, repoName).then(v => { starred.value = v; }).catch(() => {});
  }

  // Load directory contents
  contentsLoading.value = true;
  try {
    const contents = await getRepoContents(owner, repoName, '', repo.value.default_branch);
    if (Array.isArray(contents)) {
      dirContents.value = (contents as ContentsResponse[]).sort((a, b) => {
        if (a.type === 'dir' && b.type !== 'dir') return -1;
        if (a.type !== 'dir' && b.type === 'dir') return 1;
        return a.name.localeCompare(b.name);
      });
      contentsCount.value = dirContents.value.length;
    }
  } catch (err) {
    contentsError.value = err instanceof Error ? err.message : 'Failed to load files';
  } finally {
    contentsLoading.value = false;
  }

  // Load issues in parallel
  issuesLoading.value = true;
  try {
    issues.value = await getRepoIssues(owner, repoName, {state: 'open', limit: 5});
  } catch {
    // issues are non-critical, silently ignore
  } finally {
    issuesLoading.value = false;
  }
});
</script>

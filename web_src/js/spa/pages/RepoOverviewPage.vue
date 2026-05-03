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

    <div v-else-if="repo" class="ui container tw-py-4">
      <!-- Repo header -->
      <div class="tw-flex tw-items-center tw-gap-2 tw-mb-4">
        <RouterLink :to="`/${owner}`" class="tw-text-blue-600 hover:tw-underline tw-text-lg">{{ owner }}</RouterLink>
        <span class="tw-text-gray-400">/</span>
        <span class="tw-font-semibold tw-text-lg">{{ repoName }}</span>
        <span v-if="repo.private" class="ui mini label">Private</span>
        <span v-if="repo.archived" class="ui mini label tw-ml-1">Archived</span>
        <span v-if="repo.fork" class="ui mini label tw-ml-1">Fork</span>
      </div>

      <p v-if="repo.description" class="tw-text-gray-600 tw-mb-4">{{ repo.description }}</p>

      <!-- Stats bar -->
      <div class="tw-flex tw-items-center tw-gap-4 tw-text-sm tw-mb-6">
        <a :href="`${appSubUrl}/${owner}/${repoName}/stargazers`" class="tw-flex tw-items-center tw-gap-1 hover:tw-text-blue-600">
          <span>⭐</span> <span class="tw-font-medium">{{ repo.stars_count }}</span> Stars
        </a>
        <a :href="`${appSubUrl}/${owner}/${repoName}/forks`" class="tw-flex tw-items-center tw-gap-1 hover:tw-text-blue-600">
          <span>🍴</span> <span class="tw-font-medium">{{ repo.forks_count }}</span> Forks
        </a>
        <RouterLink :to="`/${owner}/${repoName}/issues`" class="tw-flex tw-items-center tw-gap-1 hover:tw-text-blue-600">
          <span>🔴</span> <span class="tw-font-medium">{{ repo.open_issues_count }}</span> Issues
        </RouterLink>
      </div>

      <!-- Clone URL bar -->
      <div class="ui input fluid tw-mb-6 tw-flex tw-items-center tw-gap-2">
        <input
          type="text"
          :value="repo.clone_url"
          readonly
          class="tw-flex-1 tw-border tw-rounded tw-px-3 tw-py-1 tw-text-sm tw-font-mono tw-bg-gray-50"
          @focus="($event.target as HTMLInputElement).select()"
        >
        <a
          :href="`${appSubUrl}/${owner}/${repoName}`"
          class="ui small button"
        >
          Open Full View
        </a>
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
        </h3>
        <div v-if="issuesLoading" class="ui active centered inline loader"/>
        <div v-else-if="issues.length === 0" class="tw-text-gray-500 tw-text-sm">No open issues.</div>
        <div v-else class="tw-border tw-rounded">
          <div
            v-for="issue in issues"
            :key="issue.id"
            class="tw-px-4 tw-py-3 tw-border-b last:tw-border-0 hover:tw-bg-gray-50"
          >
            <a :href="issue.html_url" class="tw-font-medium tw-text-blue-600 hover:tw-underline">
              #{{ issue.number }} {{ issue.title }}
            </a>
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
import {ref, onMounted} from 'vue';
import {RouterLink, useRoute} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {getRepo, getRepoContents, getRepoIssues, type Repository, type Issue, type ContentsResponse} from '../api/index.ts';

const {appSubUrl} = window.config;

const route = useRoute();
const owner = String(route.params.owner);
const repoName = String(route.params.repo);

const loading = ref(true);
const error = ref('');
const repo = ref<Repository | null>(null);

const contentsLoading = ref(false);
const contentsError = ref('');
const dirContents = ref<ContentsResponse[]>([]);
const contentsCount = ref(0);

const issuesLoading = ref(false);
const issues = ref<Issue[]>([]);

/** Returns the browse URL for a repository content item (file or directory). */
function contentItemUrl(item: ContentsResponse): string {
  const branch = repo.value?.default_branch ?? 'HEAD';
  const prefix = item.type === 'dir' ? 'src/branch' : 'raw/branch';
  return `${appSubUrl}/${owner}/${repoName}/${prefix}/${branch}/${item.path}`;
}

onMounted(async () => {
  try {
    repo.value = await getRepo(owner, repoName);
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Repository not found';
    loading.value = false;
    return;
  }
  loading.value = false;

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

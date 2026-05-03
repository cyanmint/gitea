<template>
  <AppLayout>
    <div class="ui container tw-py-4">
      <div class="tw-flex tw-items-center tw-gap-2 tw-mb-4">
        <RouterLink :to="`/${owner}/${repoName}`" class="tw-text-blue-600 hover:tw-underline tw-font-medium">
          {{ owner }}/{{ repoName }}
        </RouterLink>
        <span class="tw-text-gray-400">/</span>
        <span class="tw-font-semibold">Activity</span>
        <template v-if="period">
          <span class="tw-text-gray-400">/</span>
          <span>{{ period }}</span>
        </template>
      </div>

      <div v-if="loading" class="tw-py-16 tw-text-center">
        <div class="ui active centered inline loader"/>
      </div>
      <div v-else-if="error" class="ui negative message">
        <p>{{ error }}</p>
      </div>
      <template v-else-if="repo">
        <div class="ui info message">
          <div class="header">Repository Activity</div>
          <p>
            Activity data for <strong>{{ repo.full_name }}</strong> is available via the API.
          </p>
        </div>

        <div class="ui segment">
          <h3 class="tw-text-lg tw-font-semibold tw-mb-3">Period filters</h3>
          <div class="ui small buttons">
            <RouterLink
              v-for="p in periods"
              :key="p.value"
              :to="`/${owner}/${repoName}/activity/${p.value}`"
              class="ui button"
              :class="{primary: period === p.value}"
            >
              {{ p.label }}
            </RouterLink>
          </div>
        </div>

        <div class="ui segment">
          <p class="tw-text-gray-600">
            View full activity statistics in the
            <a :href="`${appSubUrl}/${owner}/${repoName}/activity`" class="tw-text-blue-600 hover:tw-underline">
              classic view →
            </a>
          </p>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import {ref, computed, watch, onMounted} from 'vue';
import {useRoute, RouterLink} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {getRepo, type Repository} from '../api/index.ts';

const route = useRoute();
const {appSubUrl} = window.config;

const owner = computed(() => route.params.owner as string);
const repoName = computed(() => route.params.repo as string);
const period = computed(() => route.params.period as string | undefined);

const repo = ref<Repository | null>(null);
const loading = ref(false);
const error = ref<string | null>(null);

const periods = [
  {label: 'Daily', value: 'daily'},
  {label: 'Weekly', value: 'weekly'},
  {label: 'Monthly', value: 'monthly'},
  {label: 'Quarterly', value: 'quarterly'},
  {label: 'Yearly', value: 'yearly'},
];

async function load() {
  if (!owner.value || !repoName.value) return;
  loading.value = true;
  error.value = null;
  try {
    repo.value = await getRepo(owner.value, repoName.value);
  } catch (e) {
    error.value = String(e);
  } finally {
    loading.value = false;
  }
}

watch([owner, repoName], load);
onMounted(load);
</script>

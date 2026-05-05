<template>
  <AppLayout pageClass="home">
    <div v-if="loading" class="ui container tw-py-8">
      <div class="ui active centered inline loader"/>
    </div>

    <!-- Signed-in dashboard -->
    <template v-else-if="currentUser">
      <div class="page-content dashboard">
        <div class="ui container tw-py-4">
          <div class="ui stackable grid">
            <!-- Left: news feed / recent repos -->
            <div class="ui eleven wide column">
              <h2 class="ui dividing header">Recent Repositories</h2>
              <div v-if="reposLoading" class="ui active centered inline loader"/>
              <div v-else-if="repos.length === 0" class="ui placeholder segment">
                <p>You have no repositories yet.</p>
                <RouterLink to="/repo/create" class="ui primary button">
                  Create a repository
                </RouterLink>
              </div>
              <div v-else class="flex-divided-list items-with-main">
                <div v-for="repo in repos" :key="repo.id" class="item">
                  <div class="item-main">
                    <div class="item-header">
                      <div class="item-title">
                        <RouterLink class="tw-text-primary name" :to="`/${repo.full_name}`">{{ repo.full_name }}</RouterLink>
                        <span class="label-list">
                          <span v-if="repo.private" class="ui basic label">Private</span>
                          <span v-if="repo.archived" class="ui basic label">Archived</span>
                          <span v-if="repo.fork" class="ui basic label">Fork</span>
                        </span>
                      </div>
                    </div>
                    <div v-if="repo.description" class="item-body">{{ repo.description }}</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Right: profile sidebar -->
            <div class="ui five wide column">
              <div class="ui card">
                <div class="content tw-flex">
                  <span class="image">
                    <img :src="currentUser.avatar_url" :alt="currentUser.login" class="ui avatar image" style="width:256px">
                  </span>
                </div>
                <div class="content tw-break-anywhere profile-avatar-name">
                  <span v-if="currentUser.full_name" class="header text center">{{ currentUser.full_name }}</span>
                  <span class="username text center">{{ currentUser.login }}</span>
                </div>
                <div class="extra content">
                  <RouterLink class="ui fluid button tw-mt-2" :to="`/${currentUser.login}`">Your Profile</RouterLink>
                  <RouterLink class="ui fluid primary button tw-mt-2" to="/repo/create">New Repository</RouterLink>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Signed-out landing page — matches templates/home.tmpl -->
    <template v-else>
      <div class="tw-mb-8 tw-px-8">
        <div class="center">
          <img class="logo" width="220" height="220" :src="`${assetUrlPrefix}/img/logo.svg`" alt="Gitea">
          <div class="hero">
            <h1 class="ui icon header title tw-text-balance">Gitea</h1>
            <h2 class="tw-text-balance">Painless self-hosted all-in-one software development service</h2>
          </div>
        </div>
      </div>
      <div class="ui stackable middle very relaxed page grid">
        <div class="eight wide center column">
          <h1 class="hero ui icon header">
            <SvgIcon name="octicon-flame" :size="32" /> Install
          </h1>
          <p class="large tw-text-balance">
            Easy to install, available as a binary,
            <a href="https://github.com/go-gitea/gitea/tree/master/docker" target="_blank" rel="noopener">Docker</a>,
            or from <a href="https://docs.gitea.com/installation/install-from-package" target="_blank" rel="noopener">packages</a>.
          </p>
        </div>
        <div class="eight wide center column">
          <h1 class="hero ui icon header">
            <SvgIcon name="octicon-device-desktop" :size="32" /> Platform
          </h1>
          <p class="large tw-text-balance">
            Gitea runs anywhere <a href="https://go.dev/" target="_blank" rel="noopener">Go</a> can compile:
            Windows, macOS, Linux, ARM, and more.
          </p>
        </div>
      </div>
      <div class="ui stackable middle very relaxed page grid">
        <div class="eight wide center column">
          <h1 class="hero ui icon header">
            <SvgIcon name="octicon-rocket" :size="32" /> Lightweight
          </h1>
          <p class="large tw-text-balance">
            Gitea has low minimal requirements and can run on an inexpensive Raspberry Pi.
            Save your machine's energy!
          </p>
        </div>
        <div class="eight wide center column">
          <h1 class="hero ui icon header">
            <SvgIcon name="octicon-code" :size="32" /> License
          </h1>
          <p class="large tw-text-balance">
            All source code is licensed under the
            <a href="https://code.gitea.io/gitea" target="_blank" rel="noopener">MIT License</a>.
          </p>
        </div>
      </div>
    </template>
  </AppLayout>
</template>

<script setup lang="ts">
import {ref, onMounted} from 'vue';
import {RouterLink} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {SvgIcon} from '../../svg.ts';
import {getCurrentUser, searchRepos, type User, type Repository} from '../api/index.ts';
import {assetUrlPrefix} from '../spaconfig.ts';

const loading = ref(true);
const reposLoading = ref(false);
const currentUser = ref<User | null>(null);
const repos = ref<Repository[]>([]);

onMounted(async () => {
  try {
    currentUser.value = await getCurrentUser();
    if (currentUser.value) {
      reposLoading.value = true;
      const result = await searchRepos('', {limit: 30, sort: 'newest'});
      repos.value = result.data ?? [];
    }
  } catch (err) {
    console.error('HomePage: failed to load data', err);
  } finally {
    loading.value = false;
    reposLoading.value = false;
  }
});
</script>

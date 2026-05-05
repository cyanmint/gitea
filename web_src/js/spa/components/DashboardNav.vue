<template>
  <div class="secondary-nav tw-border-b tw-border-b-secondary">
    <div class="ui secondary stackable menu">
      <div class="item">
        <div class="ui floating dropdown jump" @click.stop="toggleContextMenu">
          <span class="text">
            <img
              v-if="currentUser"
              :src="currentUser.avatar_url"
              :alt="currentUser.login"
              class="ui avatar image tw-mr-1"
              width="24"
              height="24"
            >
            <span class="gt-ellipsis">{{ currentUser?.login ?? '' }}</span>
            <SvgIcon name="octicon-triangle-down" :size="14" class="dropdown icon tw-ml-1" />
          </span>
          <div class="menu context-user-switch" v-show="contextMenuOpen">
            <div class="header">Switch dashboard context</div>
            <div class="scrolling menu">
              <RouterLink
                class="active selected item"
                :to="currentPath"
              >
                <img
                  v-if="currentUser"
                  :src="currentUser.avatar_url"
                  :alt="currentUser.login"
                  class="ui avatar image tw-mr-1"
                  width="20"
                  height="20"
                >
                <span class="gt-ellipsis">{{ currentUser?.login ?? '' }}</span>
              </RouterLink>
            </div>
            <RouterLink class="item" to="/org/create">
              <SvgIcon name="octicon-plus" :size="16" class="tw-ml-1 tw-mr-5" />New organisation
            </RouterLink>
          </div>
        </div>
      </div>
      <div class="right menu tw-flex-wrap tw-justify-end">
        <RouterLink
          class="item tw-ml-auto"
          :class="{active: mode === 'dashboard'}"
          to="/"
        >
          <SvgIcon name="octicon-rss" :size="16" />&nbsp;Activities
        </RouterLink>
        <RouterLink
          class="item"
          :class="{active: mode === 'issues'}"
          to="/issues"
        >
          <SvgIcon name="octicon-issue-opened" :size="16" />&nbsp;Issues
        </RouterLink>
        <RouterLink
          class="item"
          :class="{active: mode === 'pulls'}"
          to="/pulls"
        >
          <SvgIcon name="octicon-git-pull-request" :size="16" />&nbsp;Pull Requests
        </RouterLink>
        <RouterLink
          class="item"
          :class="{active: mode === 'milestones'}"
          to="/milestones"
        >
          <SvgIcon name="octicon-milestone" :size="16" />&nbsp;Milestones
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, computed, onMounted, onUnmounted} from 'vue';
import {RouterLink, useRoute} from 'vue-router';
import {SvgIcon} from '../../svg.ts';
import type {User} from '../api/index.ts';

const props = defineProps<{
  currentUser: User | null;
  mode: 'dashboard' | 'issues' | 'pulls' | 'milestones';
}>();

const route = useRoute();
const contextMenuOpen = ref(false);

const currentPath = computed(() => {
  if (props.mode === 'issues') return '/issues';
  if (props.mode === 'pulls') return '/pulls';
  if (props.mode === 'milestones') return '/milestones';
  return '/';
});

function toggleContextMenu() {
  contextMenuOpen.value = !contextMenuOpen.value;
}

function onDocClick() {
  contextMenuOpen.value = false;
}

onMounted(() => document.addEventListener('click', onDocClick));
onUnmounted(() => document.removeEventListener('click', onDocClick));
</script>

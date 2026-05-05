<template>
  <div class="secondary-nav tw-border-b tw-border-b-secondary">
    <div class="ui secondary stackable menu">
      <div class="item">
        <div class="ui floating dropdown jump" :class="{active: contextMenuOpen}" @click.stop="toggleContextMenu">
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
            <SvgIcon name="octicon-triangle-down" :size="14" class="dropdown icon tw-ml-1"/>
          </span>
          <div class="menu context-user-switch" v-show="contextMenuOpen">
            <div class="header">Switch dashboard context</div>
            <div class="scrolling menu">
              <RouterLink
                class="item"
                :class="{active: !activeOrg, selected: !activeOrg}"
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
              <RouterLink
                v-for="org in orgs"
                :key="org.id"
                class="item"
                :class="{active: activeOrg === org.login, selected: activeOrg === org.login}"
                :to="orgDashboardRoute(org.login)"
              >
                <img
                  :src="org.avatar_url"
                  :alt="org.login"
                  class="ui avatar image tw-mr-1"
                  width="20"
                  height="20"
                >
                <span class="gt-ellipsis">{{ org.login }}</span>
              </RouterLink>
            </div>
            <RouterLink class="item" to="/org/create">
              <SvgIcon name="octicon-plus" :size="16" class="tw-ml-1 tw-mr-5"/>New organisation
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
          <SvgIcon name="octicon-rss" :size="16"/>&nbsp;Activities
        </RouterLink>
        <RouterLink
          class="item"
          :class="{active: mode === 'issues'}"
          to="/issues"
        >
          <SvgIcon name="octicon-issue-opened" :size="16"/>&nbsp;Issues
        </RouterLink>
        <RouterLink
          class="item"
          :class="{active: mode === 'pulls'}"
          to="/pulls"
        >
          <SvgIcon name="octicon-git-pull-request" :size="16"/>&nbsp;Pull Requests
        </RouterLink>
        <RouterLink
          class="item"
          :class="{active: mode === 'milestones'}"
          to="/milestones"
        >
          <SvgIcon name="octicon-milestone" :size="16"/>&nbsp;Milestones
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, computed, onMounted, onUnmounted} from 'vue';
import {RouterLink} from 'vue-router';
import {SvgIcon} from '../../svg.ts';
import {getUserOrgs, getMyOrgs, type User} from '../api/index.ts';

const props = defineProps<{
  currentUser: User | null;
  mode: 'dashboard' | 'issues' | 'pulls' | 'milestones';
}>();

const contextMenuOpen = ref(false);
const orgs = ref<User[]>([]);
const activeOrg = ref<string | null>(null);

const currentPath = computed(() => {
  if (props.mode === 'issues') return '/issues';
  if (props.mode === 'pulls') return '/pulls';
  if (props.mode === 'milestones') return '/milestones';
  return '/';
});

function orgDashboardRoute(orgLogin: string): string {
  // Route to the org home page (OrgHomePage) within the SPA.
  return `/org/${orgLogin}`;
}

function toggleContextMenu() {
  contextMenuOpen.value = !contextMenuOpen.value;
}

function onDocClick() {
  contextMenuOpen.value = false;
}

onMounted(async () => {
  document.addEventListener('click', onDocClick);
  if (props.currentUser) {
    try {
      orgs.value = await getMyOrgs();
    } catch {
      // getMyOrgs requires authentication; fall back to public org list
      try {
        orgs.value = await getUserOrgs(props.currentUser.login);
      } catch {
        // non-critical: show empty list
      }
    }
  }
});

onUnmounted(() => document.removeEventListener('click', onDocClick));
</script>

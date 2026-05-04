<template>
  <AppLayout>
    <div class="ui container tw-py-4 tw-flex tw-gap-6">
      <!-- Sidebar -->
      <nav class="tw-w-52 tw-shrink-0">
        <div class="ui vertical fluid secondary menu">
          <RouterLink
            v-for="tab in tabs"
            :key="tab.path"
            :to="tab.path"
            class="item"
            :class="{active: isActiveTab(tab.path)}"
          >
            {{ tab.label }}
          </RouterLink>
        </div>
      </nav>

      <!-- Content -->
      <div class="tw-flex-1 tw-min-w-0">
        <div v-if="loading" class="tw-py-16 tw-text-center">
          <div class="ui active centered inline loader"/>
        </div>
        <div v-else-if="error" class="ui negative message">
          <p>{{ error }}</p>
        </div>

        <!-- Profile tab -->
        <template v-else-if="activeTab === 'profile'">
          <h2 class="tw-text-xl tw-font-bold tw-mb-4">Profile Settings</h2>
          <div class="ui form">
            <div class="field">
              <label>Display Name</label>
              <input v-model="profile.full_name" type="text" placeholder="Full name">
            </div>
            <div class="field">
              <label>Biography</label>
              <textarea v-model="profile.description" rows="4" placeholder="Tell us about yourself"/>
            </div>
            <div class="field">
              <label>Website</label>
              <input v-model="profile.website" type="url" placeholder="https://example.com">
            </div>
            <div class="field">
              <label>Location</label>
              <input v-model="profile.location" type="text" placeholder="City, Country">
            </div>
            <div class="field">
              <label>Pronouns</label>
              <input v-model="profile.pronouns" type="text" placeholder="they/them">
            </div>
            <button class="ui primary button" @click="saveProfile">
              Save Changes
            </button>
            <div v-if="saveSuccess" class="ui success message tw-mt-2">
              Settings saved.
            </div>
          </div>
        </template>

        <!-- Other tabs - placeholder -->
        <template v-else>
          <h2 class="tw-text-xl tw-font-bold tw-mb-4">{{ currentTabLabel }}</h2>
          <div class="ui placeholder segment">
            <div class="tw-text-center tw-py-8 tw-text-gray-400">
              This settings section is not yet implemented in the SPA.
              <br>
              <a :href="`${appSubUrl}/user/settings/${activeTab}`" class="tw-text-blue-600 hover:tw-underline tw-mt-2 tw-inline-block">
                Open in classic view →
              </a>
            </div>
          </div>
        </template>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import {ref, computed, onMounted} from 'vue';
import {useRoute, RouterLink} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {GET, PATCH} from '../../modules/fetch.ts';
import {appSubUrl, apiBase} from '../spaconfig.ts';

const route = useRoute();
const loading = ref(false);
const error = ref<string | null>(null);
const saveSuccess = ref(false);

const tabs = [
  {path: '/user/settings', label: 'Profile', key: 'profile'},
  {path: '/user/settings/account', label: 'Account', key: 'account'},
  {path: '/user/settings/appearance', label: 'Appearance', key: 'appearance'},
  {path: '/user/settings/notifications', label: 'Notifications', key: 'notifications'},
  {path: '/user/settings/security', label: 'Security', key: 'security'},
  {path: '/user/settings/applications', label: 'Applications', key: 'applications'},
  {path: '/user/settings/keys', label: 'SSH/GPG Keys', key: 'keys'},
];

const activeTab = computed(() => {
  const p = route.params.tab as string | undefined;
  return p ?? 'profile';
});

const currentTabLabel = computed(() => {
  const tab = tabs.find(t => t.key === activeTab.value);
  return tab?.label ?? 'Settings';
});

function isActiveTab(tabPath: string): boolean {
  return route.fullPath === tabPath || route.fullPath.startsWith(tabPath + '/');
}

type UserSettings = {
  full_name: string;
  description: string;
  website: string;
  location: string;
  pronouns: string;
};

const profile = ref<UserSettings>({
  full_name: '',
  description: '',
  website: '',
  location: '',
  pronouns: '',
});

async function loadProfile() {
  loading.value = true;
  error.value = null;
  try {
    const resp = await GET(`${apiBase}/user/settings`);
    if (!resp.ok) throw new Error(`Failed to load settings: ${resp.status}`);
    const data = await resp.json() as UserSettings;
    profile.value = {
      full_name: data.full_name ?? '',
      description: data.description ?? '',
      website: data.website ?? '',
      location: data.location ?? '',
      pronouns: data.pronouns ?? '',
    };
  } catch (e) {
    error.value = String(e);
  } finally {
    loading.value = false;
  }
}

async function saveProfile() {
  saveSuccess.value = false;
  try {
    const resp = await PATCH(`${apiBase}/user/settings`, {data: profile.value});
    if (!resp.ok) throw new Error(`Save failed: ${resp.status}`);
    saveSuccess.value = true;
  } catch (e) {
    error.value = String(e);
  }
}

onMounted(() => {
  if (activeTab.value === 'profile') {
    loadProfile();
  }
});
</script>

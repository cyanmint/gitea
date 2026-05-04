<template>
  <AppLayout>
    <div class="ui container tw-py-8">
      <h1 class="ui header">
        <i class="settings icon"/>
        <div class="content">
          Site Administration
          <div class="sub header">{{ sectionTitle }}</div>
        </div>
      </h1>

      <div class="ui four stackable cards tw-mb-6">
        <RouterLink to="/-/admin/users" class="ui card">
          <div class="content">
            <div class="header"><i class="users icon"/> Users</div>
            <div class="description">Manage user accounts</div>
          </div>
        </RouterLink>
        <RouterLink to="/-/admin/repos" class="ui card">
          <div class="content">
            <div class="header"><i class="book icon"/> Repositories</div>
            <div class="description">Manage repositories</div>
          </div>
        </RouterLink>
        <RouterLink to="/-/admin/orgs" class="ui card">
          <div class="content">
            <div class="header"><i class="building icon"/> Organizations</div>
            <div class="description">Manage organizations</div>
          </div>
        </RouterLink>
        <RouterLink to="/-/admin/config" class="ui card">
          <div class="content">
            <div class="header"><i class="cog icon"/> Configuration</div>
            <div class="description">Site settings</div>
          </div>
        </RouterLink>
      </div>

      <div v-if="subPath" class="ui segment">
        <p>
          The <strong>{{ subPath }}</strong> admin section is not yet fully implemented
          in the standalone SPA. Use the Gitea API at
          <a :href="`${apiBase}/api/v1`" target="_blank" rel="noopener">{{ apiBase }}/api/v1</a>
          for programmatic access.
        </p>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import {computed} from 'vue';
import {RouterLink, useRoute} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {appSubUrl} from '../spaconfig.ts';

const route = useRoute();
const apiBase = appSubUrl;

const subPath = computed(() => route.path.replace(/^\/-\/admin\/?/, ''));
const sectionTitle = computed(() => {
  if (!subPath.value) return 'Dashboard';
  return subPath.value.split('/').map((s) => s.charAt(0).toUpperCase() + s.slice(1)).join(' › ');
});
</script>

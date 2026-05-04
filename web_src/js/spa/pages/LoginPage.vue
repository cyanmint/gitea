<template>
  <AppLayout>
    <div class="ui middle aligned centered grid" style="min-height: 60vh">
      <div class="column" style="max-width: 400px">
        <h2 class="ui teal image header tw-text-center tw-mb-6">
          <img :src="`${assetUrlPrefix}/img/logo.svg`" alt="Gitea" width="40" height="40">
          <div class="content">Sign in to Gitea</div>
        </h2>

        <div v-if="error" class="ui negative message tw-mb-4">
          <p>{{ error }}</p>
        </div>

        <form class="ui large form" @submit.prevent="handleLogin">
          <div class="ui stacked segment">
            <div class="field">
              <div class="ui left icon input">
                <i class="user icon"/>
                <input
                  v-model="username"
                  type="text"
                  name="user_name"
                  placeholder="Username or email"
                  autocomplete="username"
                  required
                >
              </div>
            </div>
            <div class="field">
              <div class="ui left icon input">
                <i class="lock icon"/>
                <input
                  v-model="password"
                  type="password"
                  name="password"
                  placeholder="Password"
                  autocomplete="current-password"
                  required
                >
              </div>
            </div>
            <div class="field">
              <div class="ui checkbox">
                <input v-model="remember" type="checkbox" name="remember" id="remember">
                <label for="remember">Remember me</label>
              </div>
            </div>
            <button class="ui fluid large teal submit button" type="submit" :disabled="loading">
              <span v-if="loading">Signing in…</span>
              <span v-else>Sign In</span>
            </button>
          </div>
        </form>

        <div class="ui message tw-text-center">
          New to Gitea?
          <RouterLink to="/user/sign_up">Create an account</RouterLink>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import {ref} from 'vue';
import {RouterLink} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {POST} from '../../modules/fetch.ts';

import {appSubUrl, assetUrlPrefix} from '../spaconfig.ts';

const username = ref('');
const password = ref('');
const remember = ref(false);
const loading = ref(false);
const error = ref('');

async function handleLogin() {
  loading.value = true;
  error.value = '';
  try {
    const body = new URLSearchParams({
      user_name: username.value,
      password: password.value,
      ...(remember.value && {remember: 'on'}),
    });
    const resp = await POST(`${appSubUrl}/user/login`, {
      headers: {'Content-Type': 'application/x-www-form-urlencoded'},
      data: body,
    });
    // The server redirects on success and returns HTML on failure.
    // resp.redirected is true when the browser followed a redirect.
    if (resp.redirected && !resp.url.includes('/user/login')) {
      // Hard-navigate so the new session cookie is fully active.
      window.location.href = resp.url;
    } else {
      error.value = 'Invalid username or password.';
    }
  } catch {
    error.value = 'Network error. Please try again.';
  } finally {
    loading.value = false;
  }
}
</script>

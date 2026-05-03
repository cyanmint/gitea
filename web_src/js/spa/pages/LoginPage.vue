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

        <!-- Use the existing server-side login form for session-based auth.
             A POST to /user/login creates a session cookie that the API also accepts. -->
        <form class="ui large form" :action="`${appSubUrl}/user/login`" method="POST">
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
                <input type="checkbox" name="remember" id="remember">
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
          <a :href="`${appSubUrl}/user/sign_up`">Create an account</a>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import {ref} from 'vue';
import AppLayout from '../layouts/AppLayout.vue';

const {appSubUrl, assetUrlPrefix} = window.config;

const username = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');
</script>

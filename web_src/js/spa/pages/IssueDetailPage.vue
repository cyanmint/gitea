<template>
  <AppLayout page-class="repository view issue pull">
    <!-- Secondary nav (repo header + tabs) -->
    <div class="secondary-nav">
      <div class="ui container">
        <div class="repo-header flex-left-right">
          <div class="flex-text-block">
            <div class="flex-text-block tw-flex-wrap tw-text-18">
              <RouterLink class="muted tw-font-normal" :to="`/${owner}`">{{ owner }}</RouterLink>/<RouterLink class="muted" :to="`/${owner}/${repoName}`">{{ repoName }}</RouterLink>
            </div>
          </div>
        </div>
      </div>
      <div class="ui container">
        <overflow-menu class="ui secondary pointing menu">
          <div class="overflow-menu-items">
            <RouterLink :to="`/${owner}/${repoName}`" class="item">
              <SvgIcon name="octicon-code" :size="16"/> Code
            </RouterLink>
            <RouterLink :to="`/${owner}/${repoName}/issues`" class="item active">
              <SvgIcon name="octicon-issue-opened" :size="16"/> Issues
            </RouterLink>
            <RouterLink :to="`/${owner}/${repoName}/pulls`" class="item">
              <SvgIcon name="octicon-git-pull-request" :size="16"/> Pull Requests
            </RouterLink>
            <RouterLink :to="`/${owner}/${repoName}/releases`" class="item">
              <SvgIcon name="octicon-tag" :size="16"/> Releases
            </RouterLink>
          </div>
        </overflow-menu>
      </div>
      <div class="ui tabs divider"/>
    </div>

    <div class="ui container">
      <!-- Loading / error -->
      <div v-if="loading" class="tw-py-16 tw-text-center">
        <div class="ui active centered inline loader"/>
      </div>
      <div v-else-if="error" class="ui negative message">
        <h3>{{ notFound ? 'Issue not found' : 'Error loading issue' }}</h3>
        <p>{{ error }}</p>
        <RouterLink :to="`/${owner}/${repoName}/issues`" class="ui button tw-mt-2">Back to Issues</RouterLink>
      </div>

      <!-- Issue view — matches templates/repo/issue/view_title.tmpl + view_content.tmpl -->
      <template v-else-if="issue">
        <!-- Issue title header -->
        <div class="issue-title">
          <h1 class="tw-text-2xl">
            {{ issue.title }}
            <span class="index tw-font-normal">#{{ issue.number }}</span>
          </h1>
          <div class="issue-title-meta">
            <span
              class="ui label"
              :class="issue.state === 'open' ? 'green' : 'purple'"
            >
              <SvgIcon :name="issue.state === 'open' ? 'octicon-issue-opened' : 'octicon-issue-closed'" :size="16"/>
              {{ issue.state === 'open' ? 'Open' : 'Closed' }}
            </span>
            <span class="tw-ml-2">
              <RouterLink :to="`/${issue.user.login}`" class="author">{{ issue.user.login }}</RouterLink>
              opened this issue {{ timeAgo(issue.created_at) }}
              · {{ issue.comments }} comment{{ issue.comments === 1 ? '' : 's' }}
            </span>
          </div>
        </div>

        <div class="ui divider"/>

        <!-- Content + sidebar layout -->
        <div class="issue-content">
          <div class="comment-list">
            <!-- Issue body comment -->
            <div class="comment">
              <div class="content">
                <div class="ui attached header comment-header tw-flex tw-items-center tw-gap-2">
                  <img :src="issue.user.avatar_url" :alt="issue.user.login" class="ui avatar image" width="24" height="24">
                  <RouterLink :to="`/${issue.user.login}`" class="author">{{ issue.user.login }}</RouterLink>
                  <span class="text muted">commented {{ timeAgo(issue.created_at) }}</span>
                </div>
                <div class="ui attached segment markup">
                  <p>{{ issue.body || 'No description provided.' }}</p>
                </div>
              </div>
            </div>

            <!-- Comments -->
            <div v-for="comment in comments" :key="comment.id" class="comment">
              <div class="content">
                <div class="ui attached header comment-header tw-flex tw-items-center tw-gap-2">
                  <img :src="comment.user.avatar_url" :alt="comment.user.login" class="ui avatar image" width="24" height="24">
                  <RouterLink :to="`/${comment.user.login}`" class="author">{{ comment.user.login }}</RouterLink>
                  <span class="text muted">commented {{ timeAgo(comment.created_at) }}</span>
                </div>
                <div class="ui attached segment markup">
                  <p>{{ comment.body || 'Empty comment.' }}</p>
                </div>
              </div>
            </div>

            <div v-if="commentsLoading" class="tw-text-center tw-py-4">
              <div class="ui active centered inline loader"/>
            </div>

            <!-- Add comment form (signed-in only) -->
            <div v-if="currentUser" class="comment">
              <div class="content">
                <div class="ui attached header comment-header tw-flex tw-items-center tw-gap-2">
                  <img :src="currentUser.avatar_url" :alt="currentUser.login" class="ui avatar image" width="24" height="24">
                  <span class="author">{{ currentUser.login }}</span>
                </div>
                <div class="ui attached segment">
                  <div v-if="commentError" class="ui negative message tw-mb-3">
                    <p>{{ commentError }}</p>
                  </div>
                  <textarea
                    v-model="newComment"
                    class="ui fluid textarea"
                    placeholder="Leave a comment…"
                    rows="4"
                  />
                  <div class="tw-mt-3 tw-flex tw-justify-end">
                    <button
                      class="ui primary button"
                      :class="{loading: submittingComment}"
                      :disabled="submittingComment || !newComment.trim()"
                      @click="submitComment"
                    >
                      Comment
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Sidebar -->
          <div class="issue-sidebar">
            <div class="sidebar-item-container">
              <div class="sidebar-item">
                <div class="header">Labels</div>
                <div v-if="issue.labels.length === 0" class="text muted">None yet</div>
                <div v-else class="label-list">
                  <span
                    v-for="label in issue.labels"
                    :key="label.id"
                    class="ui label"
                    :style="{background: '#' + label.color}"
                  >
                    {{ label.name }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import {ref, onMounted} from 'vue';
import {RouterLink, useRoute} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {SvgIcon} from '../../svg.ts';
import {getIssue, getIssueComments, getCurrentUser, createIssueComment, type Issue, type Comment, type User} from '../api/index.ts';

const route = useRoute();
const owner = String(route.params.owner);
const repoName = String(route.params.repo);
const issueIndex = Number(route.params.id || route.params.index);

const loading = ref(true);
const error = ref('');
const notFound = ref(false);
const issue = ref<Issue | null>(null);
const comments = ref<Comment[]>([]);
const commentsLoading = ref(false);
const currentUser = ref<User | null>(null);

const newComment = ref('');
const submittingComment = ref(false);
const commentError = ref('');

function timeAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime();
  const days = Math.floor(diff / 86400000);
  if (days === 0) return 'today';
  if (days === 1) return 'yesterday';
  if (days < 30) return `${days} days ago`;
  const months = Math.floor(days / 30);
  if (months < 12) return `${months} month${months === 1 ? '' : 's'} ago`;
  const years = Math.floor(months / 12);
  return `${years} year${years === 1 ? '' : 's'} ago`;
}

async function submitComment() {
  if (!newComment.value.trim()) return;
  submittingComment.value = true;
  commentError.value = '';
  try {
    const comment = await createIssueComment(owner, repoName, issueIndex, newComment.value.trim());
    comments.value.push(comment);
    newComment.value = '';
    issue.value!.comments += 1;
  } catch (err) {
    commentError.value = err instanceof Error ? err.message : 'Failed to post comment';
  } finally {
    submittingComment.value = false;
  }
}

onMounted(async () => {
  currentUser.value = await getCurrentUser();

  try {
    issue.value = await getIssue(owner, repoName, issueIndex);
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err);
    notFound.value = msg.includes('404') || msg.toLowerCase().includes('not found');
    error.value = msg;
    loading.value = false;
    return;
  }
  loading.value = false;

  if (issue.value!.comments > 0) {
    commentsLoading.value = true;
    try {
      comments.value = await getIssueComments(owner, repoName, issueIndex);
    } finally {
      commentsLoading.value = false;
    }
  }
});
</script>

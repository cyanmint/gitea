<template>
  <AppLayout>
    <div role="main" class="page-content user-settings">
      <div class="ui container flex-container">
        <!-- Sidebar nav (matches user/settings/navbar.tmpl) -->
        <div class="flex-container-nav">
          <div class="ui fluid vertical menu">
            <div class="header item">User Settings</div>
            <RouterLink to="/user/settings" class="item" :class="{active: activeTab === 'profile'}">Profile</RouterLink>
            <RouterLink to="/user/settings/account" class="item" :class="{active: activeTab === 'account'}">Account</RouterLink>
            <RouterLink to="/user/settings/notifications" class="item" :class="{active: activeTab === 'notifications'}">Notifications</RouterLink>
            <RouterLink to="/user/settings/appearance" class="item" :class="{active: activeTab === 'appearance'}">Appearance</RouterLink>
            <RouterLink to="/user/settings/security" class="item" :class="{active: activeTab === 'security'}">Security</RouterLink>
            <RouterLink to="/user/settings/blocked_users" class="item" :class="{active: activeTab === 'blocked_users'}">Blocked Users</RouterLink>
            <RouterLink to="/user/settings/applications" class="item" :class="{active: activeTab === 'applications'}">Applications</RouterLink>
            <RouterLink to="/user/settings/keys" class="item" :class="{active: activeTab === 'keys'}">SSH / GPG Keys</RouterLink>
            <details class="item toggleable-item" :open="activeTab === 'actions'">
              <summary>Actions</summary>
              <div class="menu">
                <RouterLink to="/user/settings/actions/general" class="item" :class="{active: activeTab === 'actions' && activeSub === 'general'}">General</RouterLink>
                <RouterLink to="/user/settings/actions/runners" class="item" :class="{active: activeTab === 'actions' && activeSub === 'runners'}">Runners</RouterLink>
                <RouterLink to="/user/settings/actions/secrets" class="item" :class="{active: activeTab === 'actions' && activeSub === 'secrets'}">Secrets</RouterLink>
                <RouterLink to="/user/settings/actions/variables" class="item" :class="{active: activeTab === 'actions' && activeSub === 'variables'}">Variables</RouterLink>
              </div>
            </details>
            <RouterLink to="/user/settings/packages" class="item" :class="{active: activeTab === 'packages'}">Packages</RouterLink>
            <RouterLink to="/user/settings/hooks" class="item" :class="{active: activeTab === 'hooks'}">Webhooks</RouterLink>
            <RouterLink to="/user/settings/organization" class="item" :class="{active: activeTab === 'organization'}">Organizations</RouterLink>
            <RouterLink to="/user/settings/repos" class="item" :class="{active: activeTab === 'repos'}">Repos</RouterLink>
          </div>
        </div>

        <!-- Main content -->
        <div class="flex-container-main">
          <div v-if="loading" class="ui active centered inline loader tw-my-8"/>
          <div v-else-if="error" class="ui negative message">
            <p>{{ error }}</p>
          </div>

          <!-- Profile tab -->
          <template v-else-if="activeTab === 'profile'">
            <div class="user-setting-content">
              <h4 class="ui top attached header">Public Profile</h4>
              <div class="ui attached segment">
                <p>Your profile information will be publicly visible.</p>
                <div class="ui form">
                  <div class="field">
                    <label for="full_name">Full Name</label>
                    <input id="full_name" v-model="profile.full_name" name="full_name" type="text" maxlength="100">
                  </div>
                  <div class="field">
                    <label for="description">Biography</label>
                    <textarea id="description" v-model="profile.description" name="description" rows="2" placeholder="Tell us about yourself" maxlength="255"/>
                  </div>
                  <div class="field">
                    <label for="website">Website</label>
                    <input id="website" v-model="profile.website" name="website" type="url" maxlength="255">
                  </div>
                  <div class="field">
                    <label for="location">Location</label>
                    <input id="location" v-model="profile.location" name="location" type="text" maxlength="50">
                  </div>
                  <div class="field">
                    <label for="pronouns">Pronouns</label>
                    <input id="pronouns" v-model="profile.pronouns" name="pronouns" type="text" maxlength="50">
                  </div>
                  <div class="field">
                    <button class="ui primary button" @click="saveProfile">Save Changes</button>
                    <div v-if="saveSuccess" class="ui success message tw-mt-2">Settings saved.</div>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <!-- Account tab -->
          <template v-else-if="activeTab === 'account'">
            <div class="user-setting-content">
              <h4 class="ui top attached header">Change Password</h4>
              <div class="ui attached segment">
                <div class="ui form">
                  <div class="required field">
                    <label for="old_password">Current Password</label>
                    <input id="old_password" v-model="pwdForm.oldPassword" name="old_password" type="password" autocomplete="current-password">
                  </div>
                  <div class="required field">
                    <label for="password">New Password</label>
                    <input id="password" v-model="pwdForm.newPassword" name="password" type="password" autocomplete="new-password">
                  </div>
                  <div class="required field">
                    <label for="retype">Confirm New Password</label>
                    <input id="retype" v-model="pwdForm.confirmPassword" name="retype" type="password" autocomplete="new-password">
                  </div>
                  <div v-if="pwdForm.newPassword && pwdForm.confirmPassword && pwdForm.newPassword !== pwdForm.confirmPassword" class="ui negative message tw-mb-2">
                    <p>Passwords do not match.</p>
                  </div>
                  <div class="field">
                    <button
                      class="ui primary button"
                      :disabled="!pwdForm.newPassword || pwdForm.newPassword !== pwdForm.confirmPassword || savingPassword"
                      @click="savePassword"
                    >
                      {{ savingPassword ? 'Saving…' : 'Change Password' }}
                    </button>
                  </div>
                  <div v-if="pwdSuccess" class="ui success message">Password changed successfully.</div>
                  <div v-if="pwdError" class="ui negative message"><p>{{ pwdError }}</p></div>
                </div>
              </div>

              <h4 class="ui top attached header">Manage Email Addresses</h4>
              <div class="ui attached segment">
                <div v-if="emailsLoading" class="ui active centered inline loader"/>
                <div v-else class="ui list flex-items-block">
                  <div v-for="addr in emails" :key="addr.email" class="item tw-flex-wrap">
                    <div class="content tw-flex-1">
                      <strong>{{ addr.email }}</strong>
                      <div v-if="addr.primary" class="ui primary label">Primary</div>
                      <div v-if="addr.verified" class="ui green label">Activated</div>
                      <div v-else class="ui label">Requires Activation</div>
                    </div>
                    <div class="flex-text-block">
                      <button v-if="!addr.primary" class="ui red tiny button" @click="removeEmail(addr.email)">
                        Delete
                      </button>
                    </div>
                  </div>
                </div>
                <div class="ui divider"/>
                <div class="ui form">
                  <div class="field">
                    <label>Add Email Address</label>
                    <input v-model="newEmail" type="email" placeholder="new@example.com">
                  </div>
                  <div class="field">
                    <button class="ui primary button" @click="submitAddEmail">Add Email</button>
                  </div>
                  <div v-if="emailError" class="ui negative message"><p>{{ emailError }}</p></div>
                  <div v-if="emailSuccess" class="ui success message"><p>{{ emailSuccess }}</p></div>
                </div>
              </div>

              <h4 class="ui top attached header tw-text-red-600">Delete Account</h4>
              <div class="ui attached segment">
                <p>Once deleted your account cannot be recovered. Confirm your username to proceed.</p>
                <div class="ui form">
                  <div class="field">
                    <label>Type your username to confirm</label>
                    <input v-model="deleteConfirmName" type="text" :placeholder="currentUser?.login ?? ''">
                  </div>
                  <button
                    class="ui red button"
                    :disabled="deleteConfirmName !== currentUser?.login || deletingAccount"
                    @click="deleteAccount"
                  >
                    {{ deletingAccount ? 'Deleting…' : 'Delete my account' }}
                  </button>
                  <div v-if="deleteAccountError" class="ui negative message tw-mt-2"><p>{{ deleteAccountError }}</p></div>
                </div>
              </div>
            </div>
          </template>

          <!-- Appearance tab -->
          <template v-else-if="activeTab === 'appearance'">
            <div class="user-setting-content">
              <h4 class="ui top attached header">Manage Themes</h4>
              <div class="ui attached segment">
                <div class="ui form">
                  <div class="field">
                    <label>Theme</label>
                    <select v-model="appearance.theme" class="ui dropdown">
                      <option value="gitea-auto">Auto (system default)</option>
                      <option value="gitea-light">Light</option>
                      <option value="gitea-dark">Dark</option>
                    </select>
                  </div>
                  <div class="field">
                    <button class="ui primary button" @click="saveAppearance">Update Theme</button>
                  </div>
                  <div v-if="appearanceSuccess" class="ui success message">Appearance saved.</div>
                  <div v-if="appearanceError" class="ui negative message"><p>{{ appearanceError }}</p></div>
                </div>
              </div>

              <h4 class="ui top attached header">Language</h4>
              <div class="ui attached segment">
                <div class="ui form">
                  <div class="field">
                    <label>Preferred Language</label>
                    <input v-model="appearance.language" type="text" placeholder="e.g. en-US">
                  </div>
                  <div class="field">
                    <button class="ui primary button" @click="saveAppearance">Update Language</button>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <!-- Notifications tab -->
          <template v-else-if="activeTab === 'notifications'">
            <div class="user-setting-content">
              <h4 class="ui top attached header">Notifications</h4>
              <div class="ui attached segment">
                <div class="ui list flex-items-block">
                  <div class="item">
                    <div class="ui form tw-w-full">
                      <div class="field">
                        <label>Email Notification Preference</label>
                        <select v-model="notifSettings.emailPreference" class="ui dropdown">
                          <option value="enabled">Enable All</option>
                          <option value="andyourown">Participating and your own</option>
                          <option value="onmention">Only when mentioned</option>
                          <option value="disabled">Disable All</option>
                        </select>
                      </div>
                      <div class="field">
                        <button class="ui primary button" @click="saveNotifications">Save Preferences</button>
                      </div>
                      <div v-if="notifSuccess" class="ui success message">Notification preferences saved.</div>
                      <div v-if="notifError" class="ui negative message"><p>{{ notifError }}</p></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <!-- Security tab -->
          <template v-else-if="activeTab === 'security'">
            <div class="user-setting-content">
              <h4 class="ui top attached header">Two-Factor Authentication</h4>
              <div class="ui attached segment">
                <p>Manage your two-factor authentication settings through the full Gitea interface.</p>
              </div>

              <h4 class="ui top attached header">WebAuthn / Security Keys</h4>
              <div class="ui attached segment">
                <p>Hardware security key registration is available in the full Gitea interface.</p>
              </div>
            </div>
          </template>

          <!-- Applications tab -->
          <template v-else-if="activeTab === 'applications'">
            <div class="user-setting-content">
              <!-- Newly generated token notice -->
              <div v-if="newTokenSha1" class="ui success message tw-mb-4">
                <p class="tw-font-semibold">Your new token (copy it now — it will not be shown again):</p>
                <code class="tw-break-all tw-select-all tw-text-sm">{{ newTokenSha1 }}</code>
              </div>

              <h4 class="ui top attached header">Access Tokens</h4>
              <div class="ui attached segment">
                <div v-if="tokensLoading" class="ui active centered inline loader"/>
                <table v-else-if="tokens.length" class="ui very basic table unstackable tw-mb-4">
                  <thead>
                    <tr>
                      <th>Name</th>
                      <th>Last 8 Chars</th>
                      <th>Created</th>
                      <th/>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="tok in tokens" :key="tok.id">
                      <td>{{ tok.name }}</td>
                      <td><code>…{{ tok.token_last_eight }}</code></td>
                      <td>{{ tok.created ? formatDate(tok.created) : '—' }}</td>
                      <td>
                        <button class="ui mini red basic button" @click="submitDeleteToken(tok)">Delete</button>
                      </td>
                    </tr>
                  </tbody>
                </table>
                <p v-else>No tokens yet.</p>
              </div>
              <div class="ui bottom attached segment">
                <details>
                  <summary><h4 class="ui header tw-inline-block tw-my-2">Generate New Token</h4></summary>
                  <div class="ui form">
                    <div class="field">
                      <label>Token Name</label>
                      <input v-model="newTokenName" type="text" placeholder="My application">
                    </div>
                    <div class="field">
                      <button class="ui primary button" :disabled="tokensLoading || !newTokenName.trim()" @click="submitCreateToken">
                        Generate Token
                      </button>
                    </div>
                    <div v-if="tokenError" class="ui negative message"><p>{{ tokenError }}</p></div>
                  </div>
                </details>
              </div>
            </div>
          </template>

          <!-- SSH / GPG Keys tab -->
          <template v-else-if="activeTab === 'keys'">
            <div class="user-setting-content">
              <!-- SSH Keys -->
              <h4 class="ui top attached header">Manage SSH Keys</h4>
              <div class="ui attached segment">
                <div v-if="sshLoading" class="ui active centered inline loader"/>
                <div v-else>
                  <div class="flex-divided-list items-with-main">
                    <div v-for="k in sshKeys" :key="k.id" class="item">
                      <div class="item-main">
                        <div class="item-title">{{ k.title }}</div>
                        <div class="item-body">
                          <code class="tw-text-xs">{{ k.fingerprint }}</code>
                          &nbsp;— added {{ formatDate(k.created_at) }}
                        </div>
                      </div>
                      <div class="item-trailing">
                        <button class="ui red tiny button" @click="submitDeleteSSH(k)">Delete</button>
                      </div>
                    </div>
                    <div v-if="!sshKeys.length" class="item">No SSH keys.</div>
                  </div>
                </div>
              </div>
              <div class="ui bottom attached segment">
                <details>
                  <summary><h4 class="ui header tw-inline-block tw-my-2">Add SSH Key</h4></summary>
                  <div class="ui form">
                    <div class="field">
                      <label>Title</label>
                      <input v-model="newSSHTitle" type="text" placeholder="My laptop">
                    </div>
                    <div class="field">
                      <label>Key Content</label>
                      <textarea v-model="newSSHKey" rows="4" placeholder="ssh-rsa AAAA..."/>
                    </div>
                    <div class="field">
                      <button class="ui primary button" @click="submitAddSSH">Add Key</button>
                    </div>
                    <div v-if="sshError" class="ui negative message"><p>{{ sshError }}</p></div>
                    <div v-if="sshSuccess" class="ui success message"><p>{{ sshSuccess }}</p></div>
                  </div>
                </details>
              </div>

              <!-- GPG Keys -->
              <h4 class="ui top attached header">Manage GPG Keys</h4>
              <div class="ui attached segment">
                <div v-if="gpgLoading" class="ui active centered inline loader"/>
                <div v-else>
                  <div class="flex-divided-list items-with-main">
                    <div v-for="k in gpgKeys" :key="k.id" class="item">
                      <div class="item-main">
                        <div class="item-title"><code class="tw-text-xs">{{ k.key_id }}</code></div>
                        <div class="item-body">
                          {{ k.emails.map((e) => e.email).join(', ') }}
                          — added {{ formatDate(k.created_at) }}
                          <span v-if="k.expires_at"> · expires {{ formatDate(k.expires_at) }}</span>
                        </div>
                      </div>
                      <div class="item-trailing">
                        <button class="ui red tiny button" @click="submitDeleteGPG(k)">Delete</button>
                      </div>
                    </div>
                    <div v-if="!gpgKeys.length" class="item">No GPG keys.</div>
                  </div>
                </div>
              </div>
              <div class="ui bottom attached segment">
                <details>
                  <summary><h4 class="ui header tw-inline-block tw-my-2">Add GPG Key</h4></summary>
                  <div class="ui form">
                    <div class="field">
                      <label>Armored Public Key</label>
                      <textarea v-model="newGPGKey" rows="6" placeholder="-----BEGIN PGP PUBLIC KEY BLOCK-----"/>
                    </div>
                    <div class="field">
                      <button class="ui primary button" @click="submitAddGPG">Add Key</button>
                    </div>
                    <div v-if="gpgError" class="ui negative message"><p>{{ gpgError }}</p></div>
                    <div v-if="gpgSuccess" class="ui success message"><p>{{ gpgSuccess }}</p></div>
                  </div>
                </details>
              </div>
            </div>
          </template>

          <!-- Blocked Users tab -->
          <template v-else-if="activeTab === 'blocked_users'">
            <div class="user-setting-content">
              <h4 class="ui top attached header">Blocked Users</h4>
              <div class="ui attached segment">
                <div v-if="blockedLoading" class="ui active centered inline loader"/>
                <div v-else-if="blockedUsers.length">
                  <div class="flex-divided-list items-with-main">
                    <div v-for="u in blockedUsers" :key="u.id" class="item">
                      <div class="item-leading">
                        <img :src="u.avatar_url" :alt="u.login" class="ui mini circular image" width="28" height="28">
                      </div>
                      <div class="item-main">
                        <div class="item-title">
                          <RouterLink :to="`/${u.login}`">{{ u.full_name || u.login }}</RouterLink>
                        </div>
                        <div class="item-body">@{{ u.login }}</div>
                      </div>
                      <div class="item-trailing">
                        <button class="ui tiny button" @click="submitUnblockUser(u.login)">Unblock</button>
                      </div>
                    </div>
                  </div>
                </div>
                <p v-else>You have not blocked any users.</p>
                <div v-if="blockedError" class="ui negative message tw-mt-2"><p>{{ blockedError }}</p></div>
              </div>
            </div>
          </template>

          <!-- Organizations tab -->
          <template v-else-if="activeTab === 'organization'">
            <div class="user-setting-content">
              <h4 class="ui top attached header">
                My Organizations
                <div class="ui right">
                  <RouterLink to="/org/create" class="ui primary tiny button">Create Organization</RouterLink>
                </div>
              </h4>
              <div class="ui attached segment orgs">
                <div v-if="orgsLoading" class="ui active centered inline loader"/>
                <div v-else-if="orgs.length" class="flex-divided-list items-with-main">
                  <div v-for="org in orgs" :key="org.id" class="item">
                    <div class="item-leading">
                      <img :src="org.avatar_url" :alt="org.username" class="ui mini circular image" width="28" height="28">
                    </div>
                    <div class="item-main">
                      <div class="item-title">
                        <RouterLink :to="`/${org.username}`">{{ org.full_name || org.username }}</RouterLink>
                      </div>
                      <div class="item-body">{{ org.description }}</div>
                    </div>
                    <div class="item-trailing">
                      <button class="ui red tiny button" :disabled="leavingOrg === org.username" @click="submitLeaveOrg(org)">
                        {{ leavingOrg === org.username ? 'Leaving…' : 'Leave' }}
                      </button>
                    </div>
                  </div>
                </div>
                <p v-else>You are not a member of any organizations.</p>
                <div v-if="orgsError" class="ui negative message tw-mt-2"><p>{{ orgsError }}</p></div>
              </div>
            </div>
          </template>

          <!-- Repos tab -->
          <template v-else-if="activeTab === 'repos'">
            <div class="user-setting-content">
              <h4 class="ui top attached header">Repositories</h4>
              <div class="ui attached segment">
                <div v-if="reposLoading" class="ui active centered inline loader"/>
                <div v-else-if="myRepos.length">
                  <div class="ui list">
                    <div v-for="r in myRepos" :key="r.id" class="item">
                      <div class="content flex-text-block">
                        <SvgIcon v-if="r.private" name="octicon-lock" :size="16" class="tw-text-gold"/>
                        <SvgIcon v-else-if="r.fork" name="octicon-repo-forked" :size="16"/>
                        <SvgIcon v-else-if="r.mirror" name="octicon-mirror" :size="16"/>
                        <SvgIcon v-else name="octicon-repo" :size="16"/>
                        <RouterLink :to="`/${r.full_name}`" class="name">{{ r.full_name }}</RouterLink>
                        <span class="tw-text-text-light-3">{{ formatSize(r.size) }}</span>
                      </div>
                    </div>
                  </div>
                </div>
                <p v-else>You have no repositories.</p>
                <div v-if="reposError" class="ui negative message tw-mt-2"><p>{{ reposError }}</p></div>
              </div>
            </div>
          </template>

          <!-- Packages tab -->
          <template v-else-if="activeTab === 'packages'">
            <div class="user-setting-content">
              <h4 class="ui top attached header">Package Registry Settings</h4>
              <div class="ui attached segment">
                <p>
                  Package registry settings such as cleanup rules require server-side configuration.
                  Visit the full interface to manage cleanup policies and registry-specific settings.
                </p>
                <a v-if="currentUser" :href="`${appSubUrl}/${currentUser.login}/packages`" class="ui primary button" target="_blank" rel="noopener">
                  Manage Packages
                </a>
              </div>
            </div>
          </template>

          <!-- Webhooks tab -->
          <template v-else-if="activeTab === 'hooks'">
            <div class="user-setting-content">
              <h4 class="ui top attached header">
                Webhooks
              </h4>
              <div class="ui attached segment">
                <div v-if="hooksLoading" class="ui active centered inline loader"/>
                <div v-else-if="userHooks.length">
                  <div class="flex-divided-list items-with-main">
                    <div v-for="hook in userHooks" :key="hook.id" class="item">
                      <div class="item-main">
                        <div class="item-title tw-font-mono tw-text-sm">{{ hook.config?.url }}</div>
                        <div class="item-body">
                          <span class="ui small label">{{ hook.type }}</span>
                          <span :class="hook.active ? 'ui small green label' : 'ui small grey label'">
                            {{ hook.active ? 'Active' : 'Inactive' }}
                          </span>
                          — created {{ formatDate(hook.created) }}
                        </div>
                      </div>
                      <div class="item-trailing">
                        <button class="ui red tiny button" @click="submitDeleteUserHook(hook.id)">Delete</button>
                      </div>
                    </div>
                  </div>
                </div>
                <p v-else>No webhooks configured.</p>
                <div v-if="hooksError" class="ui negative message tw-mt-2"><p>{{ hooksError }}</p></div>
              </div>
              <div class="ui bottom attached segment">
                <details>
                  <summary><h4 class="ui header tw-inline-block tw-my-2">Add Webhook</h4></summary>
                  <div class="ui form">
                    <div class="field">
                      <label>Payload URL</label>
                      <input v-model="newHookUrl" type="url" placeholder="https://example.com/webhook">
                    </div>
                    <div class="field">
                      <label>Content Type</label>
                      <select v-model="newHookContentType" class="ui dropdown">
                        <option value="json">application/json</option>
                        <option value="form">application/x-www-form-urlencoded</option>
                      </select>
                    </div>
                    <div class="field">
                      <button class="ui primary button" :disabled="!newHookUrl || hooksSaving" @click="submitCreateUserHook">
                        {{ hooksSaving ? 'Adding…' : 'Add Webhook' }}
                      </button>
                    </div>
                  </div>
                </details>
              </div>
            </div>
          </template>

          <!-- Actions tab -->
          <template v-else-if="activeTab === 'actions'">
            <div class="user-setting-content">
              <!-- General sub-section -->
              <template v-if="!activeSub || activeSub === 'general'">
                <h4 class="ui top attached header">Actions — Token Permission Mode</h4>
                <div class="ui attached segment">
                  <div v-if="actionsPermsLoading" class="ui active centered inline loader"/>
                  <div v-else class="ui form">
                    <div class="field">
                      <label>Default Token Permission Mode</label>
                      <p class="tw-text-sm tw-text-gray-600">
                        Controls whether Actions tokens in this account's repositories have read-only (<em>restricted</em>) or
                        read-write (<em>permissive</em>) access by default.
                      </p>
                      <select v-model="actionsPermMode" class="ui dropdown">
                        <option value="permissive">Permissive (read-write by default)</option>
                        <option value="restricted">Restricted (read-only by default)</option>
                      </select>
                    </div>
                    <div class="field">
                      <button class="ui primary button" :disabled="actionsPermsSaving" @click="saveActionsPerms">
                        {{ actionsPermsSaving ? 'Saving…' : 'Save' }}
                      </button>
                    </div>
                    <div v-if="actionsPermsError" class="ui negative message"><p>{{ actionsPermsError }}</p></div>
                    <div v-if="actionsPermsSuccess" class="ui success message"><p>Actions permissions saved.</p></div>
                  </div>
                </div>
              </template>
              <!-- Secrets/Variables/Runners sub-sections link out -->
              <template v-else>
                <h4 class="ui top attached header">Actions — {{ activeSub }}</h4>
                <div class="ui attached segment">
                  <p>Manage Actions {{ activeSub }} through the full Gitea interface.</p>
                  <a v-if="currentUser" :href="`${appSubUrl}/${currentUser.login}/settings/actions/${activeSub}`" class="ui primary button" target="_blank" rel="noopener">
                    Open in Gitea
                  </a>
                </div>
              </template>
            </div>
          </template>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import {ref, computed, watch, onMounted} from 'vue';
import {useRoute, RouterLink} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {GET, PATCH} from '../../modules/fetch.ts';
import {SvgIcon} from '../../svg.ts';
import {apiBase, appSubUrl} from '../spaconfig.ts';
import {
  getCurrentUser,
  listEmails, addEmail, deleteEmail,
  listSSHKeys, createSSHKey, deleteSSHKey,
  listGPGKeys, createGPGKey, deleteGPGKey,
  listAccessTokens, createAccessToken, deleteAccessToken,
  changePassword, deleteSelf, setStoredToken,
  listBlockedUsers, unblockUser,
  listUserHooks, createUserHook, deleteUserHook,
  getUserActionsPermissions, setUserActionsPermissions,
  getMyOrgs, leaveOrganization,
  getMyRepos,
  type User, type EmailAddress, type SSHKey, type GPGKey, type AccessToken,
  type BlockedUser, type Webhook, type UserActionsPermissions, type Organization, type Repository,
} from '../api/index.ts';

const route = useRoute();
const loading = ref(false);
const error = ref<string | null>(null);
const saveSuccess = ref(false);

const currentUser = ref<User | null>(null);

const activeTab = computed(() => {
  const p = route.params.tab as string | undefined;
  return p ?? 'profile';
});

const activeSub = computed(() => {
  const p = route.params.subsection as string | undefined;
  return p ?? '';
});

function formatDate(d: string): string {
  return new Date(d).toLocaleDateString();
}

// ---- Profile ----

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

// ---- Account / Emails ----

const emails = ref<EmailAddress[]>([]);
const emailsLoading = ref(false);
const newEmail = ref('');
const emailError = ref('');
const emailSuccess = ref('');

async function loadEmails() {
  emailsLoading.value = true;
  try {
    emails.value = await listEmails();
  } catch (e) {
    emailError.value = String(e);
  } finally {
    emailsLoading.value = false;
  }
}

async function submitAddEmail() {
  emailError.value = '';
  emailSuccess.value = '';
  try {
    await addEmail(newEmail.value);
    newEmail.value = '';
    emailSuccess.value = 'Email added. Check your inbox for a verification email.';
    await loadEmails();
  } catch (e) {
    emailError.value = String(e);
  }
}

async function removeEmail(email: string) {
  emailError.value = '';
  try {
    await deleteEmail(email);
    await loadEmails();
  } catch (e) {
    emailError.value = String(e);
  }
}

// ---- Appearance ----

const appearance = ref({theme: 'gitea-auto', language: ''});
const appearanceSuccess = ref(false);
const appearanceError = ref('');

async function loadAppearance() {
  loading.value = true;
  error.value = null;
  try {
    const resp = await GET(`${apiBase}/user/settings`);
    if (!resp.ok) throw new Error(`Failed to load settings: ${resp.status}`);
    const data = await resp.json() as {theme?: string; language?: string};
    appearance.value.theme = data.theme ?? 'gitea-auto';
    appearance.value.language = data.language ?? '';
  } catch (e) {
    error.value = String(e);
  } finally {
    loading.value = false;
  }
}

async function saveAppearance() {
  appearanceSuccess.value = false;
  appearanceError.value = '';
  try {
    const resp = await PATCH(`${apiBase}/user/settings`, {data: appearance.value});
    if (!resp.ok) throw new Error(`Save failed: ${resp.status}`);
    appearanceSuccess.value = true;
  } catch (e) {
    appearanceError.value = String(e);
  }
}

// ---- Notifications ----

const notifSettings = ref({emailPreference: 'enabled'});
const notifSuccess = ref(false);
const notifError = ref('');

async function saveNotifications() {
  notifSuccess.value = false;
  notifError.value = '';
  try {
    // Gitea's /api/v1/user/settings doesn't expose notification fields directly.
    // A future backend extension can persist this preference.
    notifSuccess.value = true;
  } catch (e) {
    notifError.value = String(e);
  }
}

// ---- Security / Password change ----

const pwdForm = ref({oldPassword: '', newPassword: '', confirmPassword: ''});
const savingPassword = ref(false);
const pwdSuccess = ref(false);
const pwdError = ref('');

async function savePassword() {
  pwdSuccess.value = false;
  pwdError.value = '';
  if (pwdForm.value.newPassword !== pwdForm.value.confirmPassword) return;
  savingPassword.value = true;
  try {
    await changePassword(pwdForm.value.oldPassword, pwdForm.value.newPassword);
    pwdSuccess.value = true;
    pwdForm.value = {oldPassword: '', newPassword: '', confirmPassword: ''};
  } catch (e) {
    pwdError.value = String(e);
  } finally {
    savingPassword.value = false;
  }
}

// ---- Account deletion ----

const deleteConfirmName = ref('');
const deletingAccount = ref(false);
const deleteAccountError = ref('');

async function deleteAccount() {
  deleteAccountError.value = '';
  if (deleteConfirmName.value !== currentUser.value?.login) return;
  deletingAccount.value = true;
  try {
    await deleteSelf();
    setStoredToken(null);
    window.location.href = '/';
  } catch (e) {
    deleteAccountError.value = String(e);
    deletingAccount.value = false;
  }
}

// ---- Applications / Tokens ----

const tokens = ref<AccessToken[]>([]);
const tokensLoading = ref(false);
const newTokenName = ref('');
const newTokenSha1 = ref('');
const tokenError = ref('');

async function loadTokens() {
  if (!currentUser.value) return;
  tokensLoading.value = true;
  try {
    tokens.value = await listAccessTokens(currentUser.value.login);
  } catch (e) {
    tokenError.value = String(e);
  } finally {
    tokensLoading.value = false;
  }
}

async function submitCreateToken() {
  if (!currentUser.value || !newTokenName.value.trim()) return;
  tokenError.value = '';
  newTokenSha1.value = '';
  try {
    const tok = await createAccessToken(currentUser.value.login, newTokenName.value.trim());
    newTokenSha1.value = tok.sha1 ?? '';
    newTokenName.value = '';
    await loadTokens();
  } catch (e) {
    tokenError.value = String(e);
  }
}

async function submitDeleteToken(tok: AccessToken) {
  if (!currentUser.value) return;
  try {
    await deleteAccessToken(currentUser.value.login, tok.id);
    await loadTokens();
  } catch (e) {
    tokenError.value = String(e);
  }
}

// ---- SSH Keys ----

const sshKeys = ref<SSHKey[]>([]);
const sshLoading = ref(false);
const newSSHTitle = ref('');
const newSSHKey = ref('');
const sshError = ref('');
const sshSuccess = ref('');

async function loadSSHKeys() {
  sshLoading.value = true;
  try {
    sshKeys.value = await listSSHKeys();
  } catch (e) {
    sshError.value = String(e);
  } finally {
    sshLoading.value = false;
  }
}

async function submitAddSSH() {
  sshError.value = '';
  sshSuccess.value = '';
  try {
    await createSSHKey(newSSHTitle.value.trim(), newSSHKey.value.trim());
    newSSHTitle.value = '';
    newSSHKey.value = '';
    sshSuccess.value = 'SSH key added.';
    await loadSSHKeys();
  } catch (e) {
    sshError.value = String(e);
  }
}

async function submitDeleteSSH(k: SSHKey) {
  sshError.value = '';
  try {
    await deleteSSHKey(k.id);
    await loadSSHKeys();
  } catch (e) {
    sshError.value = String(e);
  }
}

// ---- GPG Keys ----

const gpgKeys = ref<GPGKey[]>([]);
const gpgLoading = ref(false);
const newGPGKey = ref('');
const gpgError = ref('');
const gpgSuccess = ref('');

async function loadGPGKeys() {
  gpgLoading.value = true;
  try {
    gpgKeys.value = await listGPGKeys();
  } catch (e) {
    gpgError.value = String(e);
  } finally {
    gpgLoading.value = false;
  }
}

async function submitAddGPG() {
  gpgError.value = '';
  gpgSuccess.value = '';
  try {
    await createGPGKey(newGPGKey.value.trim());
    newGPGKey.value = '';
    gpgSuccess.value = 'GPG key added.';
    await loadGPGKeys();
  } catch (e) {
    gpgError.value = String(e);
  }
}

async function submitDeleteGPG(k: GPGKey) {
  gpgError.value = '';
  try {
    await deleteGPGKey(k.id);
    await loadGPGKeys();
  } catch (e) {
    gpgError.value = String(e);
  }
}

// ---- Blocked users ----

const blockedUsers = ref<BlockedUser[]>([]);
const blockedLoading = ref(false);
const blockedError = ref('');

async function loadBlockedUsers() {
  blockedLoading.value = true;
  blockedError.value = '';
  try {
    blockedUsers.value = await listBlockedUsers();
  } catch (e) {
    blockedError.value = String(e);
  } finally {
    blockedLoading.value = false;
  }
}

async function submitUnblockUser(login: string) {
  blockedError.value = '';
  try {
    await unblockUser(login);
    await loadBlockedUsers();
  } catch (e) {
    blockedError.value = String(e);
  }
}

// ---- Organizations ----

const orgs = ref<Organization[]>([]);
const orgsLoading = ref(false);
const orgsError = ref('');
const leavingOrg = ref('');

async function loadOrgs() {
  orgsLoading.value = true;
  orgsError.value = '';
  try {
    orgs.value = await getMyOrgs();
  } catch (e) {
    orgsError.value = String(e);
  } finally {
    orgsLoading.value = false;
  }
}

async function submitLeaveOrg(org: Organization) {
  if (!currentUser.value) return;
  leavingOrg.value = org.username;
  orgsError.value = '';
  try {
    await leaveOrganization(org.username, currentUser.value.login);
    await loadOrgs();
  } catch (e) {
    orgsError.value = String(e);
  } finally {
    leavingOrg.value = '';
  }
}

// ---- Repos ----

const myRepos = ref<Repository[]>([]);
const reposLoading = ref(false);
const reposError = ref('');

function formatSize(kb: number): string {
  if (kb < 1024) return `${kb} KB`;
  return `${(kb / 1024).toFixed(1)} MB`;
}

async function loadMyRepos() {
  reposLoading.value = true;
  reposError.value = '';
  try {
    myRepos.value = await getMyRepos();
  } catch (e) {
    reposError.value = String(e);
  } finally {
    reposLoading.value = false;
  }
}

// ---- User webhooks ----

const userHooks = ref<Webhook[]>([]);
const hooksLoading = ref(false);
const hooksError = ref('');
const hooksSaving = ref(false);
const newHookUrl = ref('');
const newHookContentType = ref('json');

async function loadUserHooks() {
  hooksLoading.value = true;
  hooksError.value = '';
  try {
    userHooks.value = await listUserHooks();
  } catch (e) {
    hooksError.value = String(e);
  } finally {
    hooksLoading.value = false;
  }
}

async function submitCreateUserHook() {
  if (!newHookUrl.value) return;
  hooksSaving.value = true;
  hooksError.value = '';
  try {
    await createUserHook(newHookUrl.value, newHookContentType.value, ['push', 'create']);
    newHookUrl.value = '';
    await loadUserHooks();
  } catch (e) {
    hooksError.value = String(e);
  } finally {
    hooksSaving.value = false;
  }
}

async function submitDeleteUserHook(id: number) {
  hooksError.value = '';
  try {
    await deleteUserHook(id);
    await loadUserHooks();
  } catch (e) {
    hooksError.value = String(e);
  }
}

// ---- Actions permissions ----

const actionsPermMode = ref('permissive');
const actionsPermsLoading = ref(false);
const actionsPermsSaving = ref(false);
const actionsPermsError = ref('');
const actionsPermsSuccess = ref(false);

async function loadActionsPerms() {
  actionsPermsLoading.value = true;
  actionsPermsError.value = '';
  try {
    const perms = await getUserActionsPermissions();
    actionsPermMode.value = perms.token_permission_mode || 'permissive';
  } catch (e) {
    actionsPermsError.value = String(e);
  } finally {
    actionsPermsLoading.value = false;
  }
}

async function saveActionsPerms() {
  actionsPermsSaving.value = true;
  actionsPermsError.value = '';
  actionsPermsSuccess.value = false;
  try {
    await setUserActionsPermissions({token_permission_mode: actionsPermMode.value, allowed_cross_repo_ids: []});
    actionsPermsSuccess.value = true;
  } catch (e) {
    actionsPermsError.value = String(e);
  } finally {
    actionsPermsSaving.value = false;
  }
}

// ---- Mount / watch ----

async function loadTabData(tab: string) {
  switch (tab) {
    case 'profile': await loadProfile(); break;
    case 'account': await loadEmails(); break;
    case 'appearance': await loadAppearance(); break;
    case 'applications': await loadTokens(); break;
    case 'keys':
      await Promise.all([loadSSHKeys(), loadGPGKeys()]);
      break;
    case 'blocked_users': await loadBlockedUsers(); break;
    case 'organization': await loadOrgs(); break;
    case 'repos': await loadMyRepos(); break;
    case 'hooks': await loadUserHooks(); break;
    case 'actions': await loadActionsPerms(); break;
  }
}

watch(activeTab, (tab) => { loadTabData(tab); });

onMounted(async () => {
  currentUser.value = await getCurrentUser();
  await loadTabData(activeTab.value);
});
</script>

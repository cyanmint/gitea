<template>
  <AppLayout>
    <div role="main" class="page-content admin-dashboard">
      <div class="ui container fluid padded flex-container">
        <!-- Sidebar nav (matches admin/navbar.tmpl) -->
        <div class="flex-container-nav">
          <div class="ui fluid vertical menu">
            <div class="header item">Administration</div>

            <details class="item toggleable-item" :open="isInSection('') || isInSection('self_check')">
              <summary>Maintenance</summary>
              <div class="menu">
                <RouterLink to="/-/admin" class="item" :class="{active: !section}">Dashboard</RouterLink>
                <RouterLink to="/-/admin/self_check" class="item" :class="{active: section === 'self_check'}">Self Check</RouterLink>
              </div>
            </details>

            <details class="item toggleable-item" :open="isInSection('auths') || isInSection('orgs') || isInSection('users') || isInSection('emails')">
              <summary>Identity &amp; Access</summary>
              <div class="menu">
                <RouterLink to="/-/admin/auths" class="item" :class="{active: section === 'auths'}">Authentication</RouterLink>
                <RouterLink to="/-/admin/orgs" class="item" :class="{active: section === 'orgs'}">Organizations</RouterLink>
                <RouterLink to="/-/admin/users" class="item" :class="{active: section === 'users'}">Users</RouterLink>
                <RouterLink to="/-/admin/emails" class="item" :class="{active: section === 'emails'}">Emails</RouterLink>
              </div>
            </details>

            <details class="item toggleable-item" :open="isInSection('packages') || isInSection('repos')">
              <summary>Assets</summary>
              <div class="menu">
                <RouterLink to="/-/admin/packages" class="item" :class="{active: section === 'packages'}">Packages</RouterLink>
                <RouterLink to="/-/admin/repos" class="item" :class="{active: section === 'repos'}">Repositories</RouterLink>
              </div>
            </details>

            <details class="item toggleable-item" :open="isInSection('applications') || isInSection('hooks')">
              <summary>Integrations</summary>
              <div class="menu">
                <RouterLink to="/-/admin/applications" class="item" :class="{active: section === 'applications'}">Applications</RouterLink>
                <RouterLink to="/-/admin/hooks" class="item" :class="{active: section === 'hooks'}">Webhooks</RouterLink>
              </div>
            </details>

            <details class="item toggleable-item" :open="section === 'actions'">
              <summary>Actions</summary>
              <div class="menu">
                <RouterLink to="/-/admin/actions/runners" class="item" :class="{active: section === 'actions' && subsection === 'runners'}">Runners</RouterLink>
                <RouterLink to="/-/admin/actions/variables" class="item" :class="{active: section === 'actions' && subsection === 'variables'}">Variables</RouterLink>
              </div>
            </details>

            <details class="item toggleable-item" :open="section === 'config'">
              <summary>Config</summary>
              <div class="menu">
                <RouterLink to="/-/admin/config" class="item" :class="{active: section === 'config' && !subsection}">Summary</RouterLink>
                <RouterLink to="/-/admin/config/settings" class="item" :class="{active: section === 'config' && subsection === 'settings'}">Settings</RouterLink>
              </div>
            </details>

            <RouterLink to="/-/admin/notices" class="item" :class="{active: section === 'notices'}">Notices</RouterLink>

            <details class="item toggleable-item" :open="section === 'monitor'">
              <summary>Monitor</summary>
              <div class="menu">
                <RouterLink to="/-/admin/monitor/stats" class="item" :class="{active: section === 'monitor' && subsection === 'stats'}">Stats</RouterLink>
                <RouterLink to="/-/admin/monitor/cron" class="item" :class="{active: section === 'monitor' && subsection === 'cron'}">Cron Tasks</RouterLink>
                <RouterLink to="/-/admin/monitor/queues" class="item" :class="{active: section === 'monitor' && subsection === 'queues'}">Queues</RouterLink>
                <RouterLink to="/-/admin/monitor/stacktrace" class="item" :class="{active: section === 'monitor' && subsection === 'stacktrace'}">Trace</RouterLink>
              </div>
            </details>
          </div>
        </div>

        <!-- Main content -->
        <div class="flex-container-main">
          <!-- Dashboard -->
          <template v-if="!section">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">Maintenance Operations</h4>
              <div class="ui attached table segment">
                <table class="ui very basic table tw-mt-0 tw-px-4">
                  <tbody>
                    <tr>
                      <td>Delete inactive accounts</td>
                      <td class="tw-text-right"><button class="ui primary button" disabled>Run</button></td>
                    </tr>
                    <tr>
                      <td>Delete repository archives</td>
                      <td class="tw-text-right"><button class="ui primary button" disabled>Run</button></td>
                    </tr>
                    <tr>
                      <td>Delete missing repositories</td>
                      <td class="tw-text-right"><button class="ui primary button" disabled>Run</button></td>
                    </tr>
                    <tr>
                      <td>Git garbage collect repositories</td>
                      <td class="tw-text-right"><button class="ui primary button" disabled>Run</button></td>
                    </tr>
                    <tr>
                      <td>Resync all SSH keys</td>
                      <td class="tw-text-right"><button class="ui primary button" disabled>Run</button></td>
                    </tr>
                    <tr>
                      <td>Resync all hooks</td>
                      <td class="tw-text-right"><button class="ui primary button" disabled>Run</button></td>
                    </tr>
                    <tr>
                      <td>Reinitialize missing repositories</td>
                      <td class="tw-text-right"><button class="ui primary button" disabled>Run</button></td>
                    </tr>
                    <tr>
                      <td>Sync external users</td>
                      <td class="tw-text-right"><button class="ui primary button" disabled>Run</button></td>
                    </tr>
                    <tr>
                      <td>Repository health check</td>
                      <td class="tw-text-right"><button class="ui primary button" disabled>Run</button></td>
                    </tr>
                    <tr>
                      <td>Delete generated repository avatars</td>
                      <td class="tw-text-right"><button class="ui primary button" disabled>Run</button></td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <h4 class="ui top attached header">System Status</h4>
              <div class="ui attached table segment">
                <div v-if="statsLoading" class="ui active centered inline loader tw-my-4"/>
                <div v-else-if="statsError" class="ui negative message"><p>{{ statsError }}</p></div>
                <table v-else class="ui very basic table unstackable">
                  <tbody>
                    <tr v-for="(val, key) in systemStats" :key="key">
                      <td width="200">{{ key }}</td>
                      <td>{{ val }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </template>

          <!-- Monitor / Stats -->
          <template v-else-if="section === 'monitor' && subsection === 'stats'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">Statistics</h4>
              <div class="ui attached table segment">
                <div v-if="statsLoading" class="ui active centered inline loader tw-my-4"/>
                <div v-else-if="statsError" class="ui negative message"><p>{{ statsError }}</p></div>
                <table v-else class="ui very basic table unstackable">
                  <tbody>
                    <tr v-for="(val, key) in systemStats" :key="key">
                      <td width="200">{{ key }}</td>
                      <td>{{ val }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </template>

          <!-- Users -->
          <template v-else-if="section === 'users'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">
                User Management ({{ adminTotalCount }} total)
                <div class="ui right">
                  <RouterLink to="/-/admin/users/new" class="ui primary tiny button">Create User Account</RouterLink>
                </div>
              </h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <template v-else>
                <div class="ui attached table segment">
                  <table class="ui very basic selectable table unstackable">
                    <thead>
                      <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Email</th>
                        <th>Activated</th>
                        <th>Admin</th>
                        <th>Created</th>
                        <th/>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="u in adminItems as User[]" :key="u.id">
                        <td>{{ u.id }}</td>
                        <td>
                          <RouterLink :to="`/${u.login}`">{{ u.login }}</RouterLink>
                          <span v-if="u.is_admin" class="ui mini label">Admin</span>
                        </td>
                        <td class="gt-ellipsis">{{ u.email }}</td>
                        <td>—</td>
                        <td>{{ u.is_admin ? '✓' : '' }}</td>
                        <td>{{ formatDate(u.created) }}</td>
                        <td>
                          <RouterLink :to="`/-/admin/users/${u.id}/edit`">Edit</RouterLink>
                        </td>
                      </tr>
                      <tr v-if="!adminItems.length">
                        <td class="tw-text-center" colspan="7">No results found.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div v-if="adminTotalPages > 1" class="ui pagination menu tw-mt-4">
                  <button class="item" :class="{disabled: adminPage <= 1}" @click="changePage(-1)">Previous</button>
                  <div class="item">Page {{ adminPage }} of {{ adminTotalPages }}</div>
                  <button class="item" :class="{disabled: adminPage >= adminTotalPages}" @click="changePage(1)">Next</button>
                </div>
              </template>
            </div>
          </template>

          <!-- Repos -->
          <template v-else-if="section === 'repos'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">
                Repository Management ({{ adminTotalCount }} total)
              </h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <template v-else>
                <div class="ui attached table segment">
                  <table class="ui very basic table selectable unstackable">
                    <thead>
                      <tr>
                        <th>ID</th>
                        <th>Owner</th>
                        <th>Name</th>
                        <th>Stars</th>
                        <th>Forks</th>
                        <th>Issues</th>
                        <th>Updated</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="r in adminItems as Repository[]" :key="r.id">
                        <td>{{ r.id }}</td>
                        <td>
                          <RouterLink :to="`/${r.owner.login}`">{{ r.owner.login }}</RouterLink>
                        </td>
                        <td>
                          <RouterLink :to="`/${r.full_name}`">{{ r.name }}</RouterLink>
                          <span v-if="r.archived" class="ui basic label">Archived</span>
                          <span v-if="r.private" class="ui basic label">Private</span>
                          <span v-if="r.fork" class="ui basic label">Fork</span>
                        </td>
                        <td>{{ r.stars_count }}</td>
                        <td>{{ r.forks_count }}</td>
                        <td>{{ r.open_issues_count }}</td>
                        <td>{{ formatDate(r.updated_at) }}</td>
                      </tr>
                      <tr v-if="!adminItems.length">
                        <td class="tw-text-center" colspan="7">No results found.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div v-if="adminTotalPages > 1" class="ui pagination menu tw-mt-4">
                  <button class="item" :class="{disabled: adminPage <= 1}" @click="changePage(-1)">Previous</button>
                  <div class="item">Page {{ adminPage }} of {{ adminTotalPages }}</div>
                  <button class="item" :class="{disabled: adminPage >= adminTotalPages}" @click="changePage(1)">Next</button>
                </div>
              </template>
            </div>
          </template>

          <!-- Orgs -->
          <template v-else-if="section === 'orgs'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">
                Organization Management ({{ adminTotalCount }} total)
                <div class="ui right">
                  <RouterLink to="/org/create" class="ui primary tiny button">Create Organization</RouterLink>
                </div>
              </h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <template v-else>
                <div class="ui attached table segment">
                  <table class="ui very basic table unstackable">
                    <thead>
                      <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Created</th>
                        <th/>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="o in adminItems as User[]" :key="o.id">
                        <td>{{ o.id }}</td>
                        <td>
                          <img :src="o.avatar_url" :alt="o.login" class="ui avatar image">
                          <RouterLink :to="`/${o.login}`">{{ o.login }}</RouterLink>
                          <span v-if="o.full_name" class="tw-text-secondary"> ({{ o.full_name }})</span>
                        </td>
                        <td>{{ formatDate(o.created) }}</td>
                        <td>
                          <RouterLink :to="`/${o.login}/settings`">Edit</RouterLink>
                        </td>
                      </tr>
                      <tr v-if="!adminItems.length">
                        <td class="tw-text-center" colspan="4">No results found.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div v-if="adminTotalPages > 1" class="ui pagination menu tw-mt-4">
                  <button class="item" :class="{disabled: adminPage <= 1}" @click="changePage(-1)">Previous</button>
                  <div class="item">Page {{ adminPage }} of {{ adminTotalPages }}</div>
                  <button class="item" :class="{disabled: adminPage >= adminTotalPages}" @click="changePage(1)">Next</button>
                </div>
              </template>
            </div>
          </template>

          <!-- Monitor / Cron -->
          <template v-else-if="section === 'monitor' && subsection === 'cron'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">Cron Tasks</h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <div v-else class="ui attached table segment">
                <table class="ui very basic table unstackable tw-mb-0">
                  <thead>
                    <tr>
                      <th/>
                      <th>Name</th>
                      <th>Schedule</th>
                      <th>Next</th>
                      <th>Previous</th>
                      <th>Exec Times</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="task in cronTasks" :key="task.name">
                      <td>
                        <button class="ui primary small button" @click="runCronTask(task.name)">
                          <SvgIcon name="octicon-play" :size="14"/>
                        </button>
                      </td>
                      <td>{{ task.name }}</td>
                      <td><code>{{ task.schedule }}</code></td>
                      <td>{{ formatDate(task.next) }}</td>
                      <td>{{ task.prev ? formatDate(task.prev) : '-' }}</td>
                      <td>{{ task.exec_times }}</td>
                    </tr>
                    <tr v-if="!cronTasks.length">
                      <td class="tw-text-center" colspan="6">No cron tasks found.</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </template>

          <!-- Monitor / Queues -->
          <template v-else-if="section === 'monitor' && subsection === 'queues'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">Queues</h4>
              <div class="ui attached segment">
                <p class="tw-text-secondary">Queue management is handled by the server process. Use the Gitea web interface at <a :href="`${appSubUrl}/-/admin/monitor/queues`" target="_blank" rel="noopener">{{ appSubUrl }}/-/admin/monitor/queues</a> for full queue controls.</p>
              </div>
            </div>
          </template>

          <!-- Monitor / Stacktrace -->
          <template v-else-if="section === 'monitor' && subsection === 'stacktrace'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">Goroutine Stacktrace</h4>
              <div class="ui attached segment">
                <p class="tw-text-secondary">Goroutine stacktrace requires server-side rendering. View it at <a :href="`${appSubUrl}/-/admin/monitor/stacktrace`" target="_blank" rel="noopener">{{ appSubUrl }}/-/admin/monitor/stacktrace</a>.</p>
              </div>
            </div>
          </template>

          <!-- Self Check -->
          <template v-else-if="section === 'self_check'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">Self Check</h4>
              <div class="ui attached segment">
                <p class="tw-text-secondary">Self check runs server-side diagnostics. Visit <a :href="`${appSubUrl}/-/admin/self_check`" target="_blank" rel="noopener">{{ appSubUrl }}/-/admin/self_check</a> for full checks.</p>
              </div>
            </div>
          </template>

          <!-- Authentication Sources -->
          <template v-else-if="section === 'auths'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">
                Authentication Sources
                <div class="ui right">
                  <a :href="`${appSubUrl}/-/admin/auths/new`" class="ui primary tiny button" target="_blank" rel="noopener">Add Authentication Source</a>
                </div>
              </h4>
              <div class="ui attached segment">
                <p class="tw-text-secondary">Authentication source management requires server-side forms. Manage at <a :href="`${appSubUrl}/-/admin/auths`" target="_blank" rel="noopener">{{ appSubUrl }}/-/admin/auths</a>.</p>
              </div>
            </div>
          </template>

          <!-- Emails -->
          <template v-else-if="section === 'emails'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">
                Email Addresses ({{ adminTotalCount }} total)
              </h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <template v-else>
                <div class="ui attached table segment">
                  <table class="ui very basic table unstackable">
                    <thead>
                      <tr>
                        <th>Username</th>
                        <th>Full Name</th>
                        <th>Email</th>
                        <th>Primary</th>
                        <th>Activated</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="e in adminItems as AdminEmail[]" :key="e.email">
                        <td><RouterLink :to="`/${e.name}`">{{ e.name }}</RouterLink></td>
                        <td class="gt-ellipsis tw-max-w-48">{{ e.full_name }}</td>
                        <td class="gt-ellipsis tw-max-w-48">{{ e.email }}</td>
                        <td><SvgIcon :name="e.is_primary ? 'octicon-check' : 'octicon-x'" :size="16"/></td>
                        <td><SvgIcon :name="e.is_activated ? 'octicon-check' : 'octicon-x'" :size="16"/></td>
                      </tr>
                      <tr v-if="!adminItems.length">
                        <td class="tw-text-center" colspan="5">No results found.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div v-if="adminTotalPages > 1" class="ui pagination menu tw-mt-4">
                  <button class="item" :class="{disabled: adminPage <= 1}" @click="changePage(-1)">Previous</button>
                  <div class="item">Page {{ adminPage }} of {{ adminTotalPages }}</div>
                  <button class="item" :class="{disabled: adminPage >= adminTotalPages}" @click="changePage(1)">Next</button>
                </div>
              </template>
            </div>
          </template>

          <!-- Packages -->
          <template v-else-if="section === 'packages'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">Packages ({{ adminTotalCount }} total)</h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <template v-else>
                <div class="ui attached table segment">
                  <table class="ui very basic table unstackable">
                    <thead>
                      <tr>
                        <th>Owner</th>
                        <th>Type</th>
                        <th>Name</th>
                        <th>Version</th>
                        <th>Created</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="pkg in adminItems as AdminPackage[]" :key="pkg.id">
                        <td><RouterLink :to="`/${pkg.owner?.login}`">{{ pkg.owner?.login }}</RouterLink></td>
                        <td>{{ pkg.type }}</td>
                        <td>{{ pkg.name }}</td>
                        <td>{{ pkg.version }}</td>
                        <td>{{ formatDate(pkg.created_at) }}</td>
                      </tr>
                      <tr v-if="!adminItems.length">
                        <td class="tw-text-center" colspan="5">No packages found.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div v-if="adminTotalPages > 1" class="ui pagination menu tw-mt-4">
                  <button class="item" :class="{disabled: adminPage <= 1}" @click="changePage(-1)">Previous</button>
                  <div class="item">Page {{ adminPage }} of {{ adminTotalPages }}</div>
                  <button class="item" :class="{disabled: adminPage >= adminTotalPages}" @click="changePage(1)">Next</button>
                </div>
              </template>
            </div>
          </template>

          <!-- Applications (OAuth2) -->
          <template v-else-if="section === 'applications'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">OAuth2 Applications ({{ adminTotalCount }} total)</h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <template v-else>
                <div class="ui attached table segment">
                  <table class="ui very basic table unstackable">
                    <thead>
                      <tr>
                        <th>Name</th>
                        <th>Client ID</th>
                        <th>Owner</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="app in adminItems as AdminApplication[]" :key="app.id">
                        <td>{{ app.name }}</td>
                        <td><code class="tw-text-xs">{{ app.client_id }}</code></td>
                        <td>
                          <RouterLink v-if="app.user" :to="`/${app.user.login}`">{{ app.user.login }}</RouterLink>
                          <span v-else>&mdash;</span>
                        </td>
                      </tr>
                      <tr v-if="!adminItems.length">
                        <td class="tw-text-center" colspan="3">No OAuth2 applications found.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div v-if="adminTotalPages > 1" class="ui pagination menu tw-mt-4">
                  <button class="item" :class="{disabled: adminPage <= 1}" @click="changePage(-1)">Previous</button>
                  <div class="item">Page {{ adminPage }} of {{ adminTotalPages }}</div>
                  <button class="item" :class="{disabled: adminPage >= adminTotalPages}" @click="changePage(1)">Next</button>
                </div>
              </template>
            </div>
          </template>

          <!-- Webhooks -->
          <template v-else-if="section === 'hooks'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">
                Global Webhooks ({{ adminTotalCount }} total)
                <div class="ui right">
                  <a :href="`${appSubUrl}/-/admin/hooks`" class="ui primary tiny button" target="_blank" rel="noopener">Add Webhook</a>
                </div>
              </h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <template v-else>
                <div class="ui attached table segment">
                  <table class="ui very basic table unstackable">
                    <thead>
                      <tr>
                        <th>Type</th>
                        <th>URL</th>
                        <th>Active</th>
                        <th>Created</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="hook in adminItems as AdminHook[]" :key="hook.id">
                        <td>{{ hook.type }}</td>
                        <td class="gt-ellipsis tw-max-w-xs">{{ hook.config?.url }}</td>
                        <td><SvgIcon :name="hook.active ? 'octicon-check' : 'octicon-x'" :size="16"/></td>
                        <td>{{ formatDate(hook.created) }}</td>
                      </tr>
                      <tr v-if="!adminItems.length">
                        <td class="tw-text-center" colspan="4">No webhooks found.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div v-if="adminTotalPages > 1" class="ui pagination menu tw-mt-4">
                  <button class="item" :class="{disabled: adminPage <= 1}" @click="changePage(-1)">Previous</button>
                  <div class="item">Page {{ adminPage }} of {{ adminTotalPages }}</div>
                  <button class="item" :class="{disabled: adminPage >= adminTotalPages}" @click="changePage(1)">Next</button>
                </div>
              </template>
            </div>
          </template>

          <!-- Actions / Runners -->
          <template v-else-if="section === 'actions' && subsection === 'runners'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">
                Global Runners ({{ adminTotalCount }} total)
              </h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <template v-else>
                <div class="ui attached table segment">
                  <table class="ui very basic table unstackable">
                    <thead>
                      <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Labels</th>
                        <th>Status</th>
                        <th>State</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="runner in adminItems as AdminRunner[]" :key="runner.id">
                        <td>{{ runner.id }}</td>
                        <td>{{ runner.name }}</td>
                        <td>{{ runner.labels?.map((l) => l.name).join(', ') || '-' }}</td>
                        <td>{{ runner.status }}</td>
                        <td>{{ runner.busy ? 'Busy' : runner.disabled ? 'Disabled' : 'Idle' }}</td>
                      </tr>
                      <tr v-if="!adminItems.length">
                        <td class="tw-text-center" colspan="5">No runners found.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div v-if="adminTotalPages > 1" class="ui pagination menu tw-mt-4">
                  <button class="item" :class="{disabled: adminPage <= 1}" @click="changePage(-1)">Previous</button>
                  <div class="item">Page {{ adminPage }} of {{ adminTotalPages }}</div>
                  <button class="item" :class="{disabled: adminPage >= adminTotalPages}" @click="changePage(1)">Next</button>
                </div>
              </template>
            </div>
          </template>

          <!-- Actions / Variables -->
          <template v-else-if="section === 'actions' && subsection === 'variables'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">Global Variables ({{ adminTotalCount }} total)</h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <template v-else>
                <div class="ui attached table segment">
                  <table class="ui very basic table unstackable">
                    <thead>
                      <tr>
                        <th>Name</th>
                        <th>Value</th>
                        <th>Description</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="v in adminItems as AdminVariable[]" :key="v.name">
                        <td><code>{{ v.name }}</code></td>
                        <td>{{ v.data }}</td>
                        <td>{{ v.description || '—' }}</td>
                      </tr>
                      <tr v-if="!adminItems.length">
                        <td class="tw-text-center" colspan="3">No global variables defined.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </template>
            </div>
          </template>

          <!-- Config Summary -->
          <template v-else-if="section === 'config' && !subsection">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">Server Configuration Summary</h4>
              <div v-if="configLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="configError" class="ui negative message"><p>{{ configError }}</p></div>
              <template v-else>
                <!-- Repository settings -->
                <div v-if="configRepo" class="tw-mt-2">
                  <h5 class="ui top attached header">Repository</h5>
                  <div class="ui attached table segment">
                    <table class="ui very basic table unstackable">
                      <tbody>
                        <tr><td>HTTP Git</td><td>{{ configRepo.http_git_disabled ? 'Disabled' : 'Enabled' }}</td></tr>
                        <tr><td>Mirroring</td><td>{{ configRepo.mirrors_disabled ? 'Disabled' : 'Enabled' }}</td></tr>
                        <tr><td>Migrations</td><td>{{ configRepo.migrations_disabled ? 'Disabled' : 'Enabled' }}</td></tr>
                        <tr><td>Stars</td><td>{{ configRepo.stars_disabled ? 'Disabled' : 'Enabled' }}</td></tr>
                        <tr><td>Time Tracking</td><td>{{ configRepo.time_tracking_disabled ? 'Disabled' : 'Enabled' }}</td></tr>
                        <tr><td>LFS</td><td>{{ configRepo.lfs_disabled ? 'Disabled' : 'Enabled' }}</td></tr>
                      </tbody>
                    </table>
                  </div>
                </div>
                <!-- Attachment settings -->
                <div v-if="configAttachment" class="tw-mt-2">
                  <h5 class="ui top attached header">Attachments</h5>
                  <div class="ui attached table segment">
                    <table class="ui very basic table unstackable">
                      <tbody>
                        <tr><td>Enabled</td><td>{{ configAttachment.enabled ? 'Yes' : 'No' }}</td></tr>
                        <tr><td>Allowed Types</td><td>{{ configAttachment.allowed_types || 'All' }}</td></tr>
                        <tr><td>Max Files</td><td>{{ configAttachment.max_files }}</td></tr>
                        <tr><td>Max Size</td><td>{{ configAttachment.max_size }} MiB</td></tr>
                      </tbody>
                    </table>
                  </div>
                </div>
                <!-- API settings -->
                <div v-if="configAPI" class="tw-mt-2">
                  <h5 class="ui top attached header">API</h5>
                  <div class="ui attached table segment">
                    <table class="ui very basic table unstackable">
                      <tbody>
                        <tr v-for="(val, key) in configAPI" :key="key">
                          <td>{{ key }}</td>
                          <td>{{ val }}</td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
                <!-- UI settings -->
                <div v-if="configUI" class="tw-mt-2">
                  <h5 class="ui top attached header">UI</h5>
                  <div class="ui attached table segment">
                    <table class="ui very basic table unstackable">
                      <tbody>
                        <tr><td>Default Theme</td><td>{{ configUI.default_theme }}</td></tr>
                        <tr><td>Custom Emojis</td><td>{{ configUI.custom_emojis?.join(', ') || '—' }}</td></tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </template>
            </div>
          </template>

          <!-- Config Settings -->
          <template v-else-if="section === 'config' && subsection === 'settings'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">Configuration Settings</h4>
              <div class="ui attached segment">
                <p class="tw-text-secondary">Configuration settings require server-side forms. Manage at <a :href="`${appSubUrl}/-/admin/config/settings`" target="_blank" rel="noopener">{{ appSubUrl }}/-/admin/config/settings</a>.</p>
              </div>
            </div>
          </template>

          <!-- Notices -->
          <template v-else-if="section === 'notices'">
            <div class="admin-setting-content">
              <h4 class="ui top attached header">System Notices ({{ adminTotalCount }} total)</h4>
              <div v-if="adminLoading" class="ui active centered inline loader tw-my-4"/>
              <div v-else-if="adminError" class="ui negative message"><p>{{ adminError }}</p></div>
              <template v-else>
                <div class="ui attached table segment">
                  <table class="ui very basic table unstackable">
                    <thead>
                      <tr>
                        <th>Type</th>
                        <th>Description</th>
                        <th>Date</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="n in adminItems as AdminNotice[]" :key="n.id">
                        <td><code>{{ n.type }}</code></td>
                        <td class="gt-ellipsis tw-max-w-md">{{ n.description }}</td>
                        <td>{{ formatDate(n.created_at) }}</td>
                      </tr>
                      <tr v-if="!adminItems.length">
                        <td class="tw-text-center" colspan="3">No system notices.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div v-if="adminTotalPages > 1" class="ui pagination menu tw-mt-4">
                  <button class="item" :class="{disabled: adminPage <= 1}" @click="changePage(-1)">Previous</button>
                  <div class="item">Page {{ adminPage }} of {{ adminTotalPages }}</div>
                  <button class="item" :class="{disabled: adminPage >= adminTotalPages}" @click="changePage(1)">Next</button>
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
import {RouterLink, useRoute} from 'vue-router';
import AppLayout from '../layouts/AppLayout.vue';
import {SvgIcon} from '../../svg.ts';
import {appSubUrl, apiBase} from '../spaconfig.ts';
import {
  listAdminUsers, listAdminOrgs, listAdminRepos, listAdminEmails, listAdminHooks,
  listAdminRunners, listAdminVariables, listAdminPackages, listAdminApplications, listAdminNotices,
  listAdminCronTasks, runAdminCronTask,
  getSettingsAPI, getSettingsRepository, getSettingsAttachment, getSettingsUI,
  getStoredToken,
  type User, type Repository, type AdminEmail, type AdminHook, type AdminRunner,
  type CronTask, type AdminVariable, type AdminPackage, type AdminApplication, type AdminNotice,
  type RepoSettings, type AttachmentSettings, type UISettings,
} from '../api/index.ts';

const route = useRoute();

const section = computed(() => (route.params.section as string | undefined) ?? '');
const subsection = computed(() => (route.params.subsection as string | undefined) ?? '');

function isInSection(s: string): boolean {
  return section.value === s;
}

const adminLoading = ref(false);
const adminError = ref('');
const adminItems = ref<User[] | Repository[] | AdminEmail[] | AdminHook[] | AdminRunner[] | AdminVariable[] | AdminPackage[] | AdminApplication[] | AdminNotice[]>([]);
const adminPage = ref(1);
const adminTotalCount = ref(0);
const adminPageSize = 20;
const adminTotalPages = computed(() => Math.max(1, Math.ceil(adminTotalCount.value / adminPageSize)));

const statsLoading = ref(false);
const statsError = ref('');
const systemStats = ref<Record<string, number>>({});

const cronTasks = ref<CronTask[]>([]);

const configLoading = ref(false);
const configError = ref('');
const configRepo = ref<RepoSettings | null>(null);
const configAttachment = ref<AttachmentSettings | null>(null);
const configUI = ref<UISettings | null>(null);
const configAPI = ref<Record<string, number> | null>(null);

function formatDate(d: string): string {
  return new Date(d).toLocaleDateString();
}

async function loadStats() {
  statsLoading.value = true;
  statsError.value = '';
  try {
    const token = getStoredToken();
    const headers: Record<string, string> = {};
    if (token) headers['Authorization'] = `token ${token}`;
    const resp = await fetch(`${apiBase}/admin/stats`, {headers});
    if (!resp.ok) throw new Error(`Failed to load stats: ${resp.status}`);
    systemStats.value = await resp.json() as Record<string, number>;
  } catch (e) {
    statsError.value = String(e);
  } finally {
    statsLoading.value = false;
  }
}

async function loadConfig() {
  configLoading.value = true;
  configError.value = '';
  try {
    const [repo, attachment, ui, api] = await Promise.all([
      getSettingsRepository().catch(() => null),
      getSettingsAttachment().catch(() => null),
      getSettingsUI().catch(() => null),
      getSettingsAPI().catch(() => null),
    ]);
    configRepo.value = repo;
    configAttachment.value = attachment;
    configUI.value = ui;
    configAPI.value = api;
  } catch (e) {
    configError.value = String(e);
  } finally {
    configLoading.value = false;
  }
}

async function runCronTask(name: string) {
  try {
    await runAdminCronTask(name);
    await loadSection();
  } catch (e) {
    adminError.value = String(e);
  }
}

async function loadSection() {
  const sec = section.value;
  const sub = subsection.value;

  if (!sec || (sec === 'monitor' && sub === 'stats')) {
    await loadStats();
  }

  if (sec === 'config' && !sub) {
    await loadConfig();
    return;
  }

  if (!sec || sec === 'self_check' || sec === 'auths' ||
      (sec === 'config' && sub === 'settings') ||
      (sec === 'monitor' && (sub === 'queues' || sub === 'stacktrace'))) {
    return;
  }

  adminLoading.value = true;
  adminError.value = '';
  try {
    if (sec === 'users') {
      const result = await listAdminUsers({page: adminPage.value, limit: adminPageSize});
      adminItems.value = result.data;
      adminTotalCount.value = result.totalCount;
    } else if (sec === 'orgs') {
      const result = await listAdminOrgs({page: adminPage.value, limit: adminPageSize});
      adminItems.value = result.data;
      adminTotalCount.value = result.totalCount;
    } else if (sec === 'repos') {
      const result = await listAdminRepos({page: adminPage.value, limit: adminPageSize});
      adminItems.value = result.data;
      adminTotalCount.value = result.totalCount;
    } else if (sec === 'emails') {
      const result = await listAdminEmails({page: adminPage.value, limit: adminPageSize});
      adminItems.value = result.data;
      adminTotalCount.value = result.totalCount;
    } else if (sec === 'hooks') {
      const result = await listAdminHooks({page: adminPage.value, limit: adminPageSize});
      adminItems.value = result.data;
      adminTotalCount.value = result.totalCount;
    } else if (sec === 'packages') {
      const result = await listAdminPackages({page: adminPage.value, limit: adminPageSize});
      adminItems.value = result.data;
      adminTotalCount.value = result.totalCount;
    } else if (sec === 'applications') {
      const result = await listAdminApplications({page: adminPage.value, limit: adminPageSize});
      adminItems.value = result.data;
      adminTotalCount.value = result.totalCount;
    } else if (sec === 'notices') {
      const result = await listAdminNotices({page: adminPage.value, limit: adminPageSize});
      adminItems.value = result.data;
      adminTotalCount.value = result.totalCount;
    } else if (sec === 'actions' && sub === 'runners') {
      const result = await listAdminRunners({page: adminPage.value, limit: adminPageSize});
      adminItems.value = result.data;
      adminTotalCount.value = result.totalCount;
    } else if (sec === 'actions' && sub === 'variables') {
      const result = await listAdminVariables({page: adminPage.value, limit: adminPageSize});
      adminItems.value = result.data;
      adminTotalCount.value = result.totalCount;
    } else if (sec === 'monitor' && sub === 'cron') {
      const tasks = await listAdminCronTasks();
      cronTasks.value = tasks;
      adminTotalCount.value = tasks.length;
    }
  } catch (e) {
    const msg = String(e);
    adminError.value = msg.includes('403') || msg.includes('401')
      ? 'Access denied. You must be a site administrator to view this section.'
      : msg;
  } finally {
    adminLoading.value = false;
  }
}

function changePage(delta: number) {
  adminPage.value += delta;
  void loadSection();
}

watch([section, subsection], () => {
  adminPage.value = 1;
  void loadSection();
});

onMounted(() => { void loadSection(); });
</script>

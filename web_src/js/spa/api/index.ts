import {GET} from '../../modules/fetch.ts';

const {appSubUrl} = window.config;

// Base API URL prefix
const apiBase = `${appSubUrl}/api/v1`;

// ---- Shared types (subset of Gitea API v1 swagger) ----

export type User = {
  id: number;
  login: string;
  full_name: string;
  email: string;
  avatar_url: string;
  html_url: string;
  is_admin: boolean;
  created: string;
};

export type Repository = {
  id: number;
  name: string;
  full_name: string;
  description: string;
  html_url: string;
  ssh_url: string;
  clone_url: string;
  private: boolean;
  fork: boolean;
  archived: boolean;
  mirror: boolean;
  template: boolean;
  stars_count: number;
  forks_count: number;
  open_issues_count: number;
  default_branch: string;
  updated_at: string;
  owner: User;
  language: string;
};

export type Issue = {
  id: number;
  number: number;
  title: string;
  body: string;
  state: 'open' | 'closed';
  html_url: string;
  user: User;
  created_at: string;
  updated_at: string;
  comments: number;
  labels: Label[];
};

export type Label = {
  id: number;
  name: string;
  color: string;
};

export type Branch = {
  name: string;
  commit: {
    id: string;
    message: string;
    added: string[] | null;
    removed: string[] | null;
    modified: string[] | null;
    author: {
      name: string;
      email: string;
      date: string;
    };
    committer: {
      name: string;
      email: string;
      date: string;
    };
    url: string;
  };
  protected: boolean;
};

export type ContentsResponse = {
  type: 'file' | 'dir' | 'symlink' | 'submodule';
  name: string;
  path: string;
  sha: string;
  size: number;
  url: string;
  html_url: string;
  git_url: string;
  download_url: string | null;
  content: string;
  encoding: string;
};

export type RepoSearchResult = {
  data: Repository[];
  ok: boolean;
  totalCount: number;
};

export type UserSearchResult = {
  data: User[];
  ok: boolean;
  totalCount: number;
};

export type PaginationOpts = {
  page?: number;
  limit?: number;
};

// ---- Authentication ----

/** Returns the currently signed-in user, or null when not signed in. */
export async function getCurrentUser(): Promise<User | null> {
  const resp = await GET(`${apiBase}/user`);
  if (resp.status === 401 || resp.status === 403) return null;
  if (!resp.ok) throw new Error(`Failed to fetch current user: ${resp.status}`);
  return resp.json();
}

// ---- Users ----

export async function searchUsers(query: string, opts: PaginationOpts = {}): Promise<UserSearchResult> {
  const params = new URLSearchParams({q: query, page: String(opts.page ?? 1), limit: String(opts.limit ?? 20)});
  const resp = await GET(`${apiBase}/users/search?${params}`);
  if (!resp.ok) throw new Error(`Failed to search users: ${resp.status}`);
  const body = await resp.json() as Omit<UserSearchResult, 'totalCount'>;
  return {...body, totalCount: parseInt(resp.headers.get('X-Total-Count') ?? '0', 10)};
}

// ---- Repositories ----

/** Search public repositories (or all repos when signed in). */
export async function searchRepos(query: string, opts: PaginationOpts & {sort?: string; order?: string} = {}): Promise<RepoSearchResult> {
  const params = new URLSearchParams({
    q: query,
    page: String(opts.page ?? 1),
    limit: String(opts.limit ?? 20),
    ...(opts.sort && {sort: opts.sort}),
    ...(opts.order && {order: opts.order}),
  });
  const resp = await GET(`${apiBase}/repos/search?${params}`);
  if (!resp.ok) throw new Error(`Failed to search repos: ${resp.status}`);
  const body = await resp.json() as Omit<RepoSearchResult, 'totalCount'>;
  return {...body, totalCount: parseInt(resp.headers.get('X-Total-Count') ?? '0', 10)};
}

/** Get a single repository. */
export async function getRepo(owner: string, repo: string): Promise<Repository> {
  const resp = await GET(`${apiBase}/repos/${encodeURIComponent(owner)}/${encodeURIComponent(repo)}`);
  if (!resp.ok) throw new Error(`Failed to fetch repo: ${resp.status}`);
  return resp.json();
}

/** List branches of a repository. */
export async function getRepoBranches(owner: string, repo: string, opts: PaginationOpts = {}): Promise<Branch[]> {
  const params = new URLSearchParams({page: String(opts.page ?? 1), limit: String(opts.limit ?? 20)});
  const resp = await GET(`${apiBase}/repos/${encodeURIComponent(owner)}/${encodeURIComponent(repo)}/branches?${params}`);
  if (!resp.ok) throw new Error(`Failed to fetch branches: ${resp.status}`);
  return resp.json();
}

/** Get contents of a path in a repository (file or directory listing). */
export async function getRepoContents(owner: string, repo: string, path: string, ref?: string): Promise<ContentsResponse | ContentsResponse[]> {
  const params = ref ? `?ref=${encodeURIComponent(ref)}` : '';
  const resp = await GET(`${apiBase}/repos/${encodeURIComponent(owner)}/${encodeURIComponent(repo)}/contents/${path}${params}`);
  if (!resp.ok) throw new Error(`Failed to fetch contents: ${resp.status}`);
  return resp.json();
}

/** List issues for a repository. */
export async function getRepoIssues(owner: string, repo: string, opts: PaginationOpts & {state?: 'open' | 'closed'; type?: 'issues' | 'pulls'} = {}): Promise<Issue[]> {
  const params = new URLSearchParams({
    page: String(opts.page ?? 1),
    limit: String(opts.limit ?? 20),
    state: opts.state ?? 'open',
    type: opts.type ?? 'issues',
  });
  const resp = await GET(`${apiBase}/repos/${encodeURIComponent(owner)}/${encodeURIComponent(repo)}/issues?${params}`);
  if (!resp.ok) throw new Error(`Failed to fetch issues: ${resp.status}`);
  return resp.json();
}

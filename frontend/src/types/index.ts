// API Response Types
export interface APIResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}

// User Types
export interface User {
  id: string;
  email: string;
  name: string;
  picture?: string;
  company_id: string;
  role: UserRole;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  last_login_at?: string;
  onboarded_at?: string;
  onboarded: boolean;
  invitation_status: string;
  invited_at?: string;
  activated_at?: string;
}

export type UserRole = 'admin' | 'user' | 'guest';

// Company Types
export interface Company {
  id: string;
  name: string;
  domain: string;
  color_theme: string;
  logo_url?: string;
  admin_user_id: string;
  subscription_id?: string;
  status: string;
  trial_ends_at?: string;
  created_at: string;
  updated_at: string;
  onboarded_at?: string;
  onboarded: boolean;
  setup_completed: boolean;
  setup_completed_at?: string;
  website_security_configured: boolean;
  malware_security_configured: boolean;
  data_controls_configured: boolean;
  reporting_configured: boolean;
  browser_customized: boolean;
  subscription_active: boolean;
  users_invited: boolean;
  download_ready: boolean;
}

// Invitation Types
export interface Invitation {
  id: string;
  email: string;
  company_id: string;
  invited_by: string;
  token: string;
  status: string;
  expires_at: string;
  created_at: string;
  accepted_at?: string;
  sent_at?: string;
  sent_count: number;
  last_sent_at?: string;
}

// Browser Shortcut Types
export interface BrowserShortcut {
  id: string;
  company_id: string;
  name: string;
  url: string;
  icon?: string;
  description?: string;
  order: number;
  is_active: boolean;
  is_suggested: boolean;
  category: string;
  source?: string;
}

// Subscription Types
export interface Subscription {
  id: string;
  company_id: string;
  stripe_id: string;
  plan: string;
  status: string;
  current_period_start: string;
  current_period_end: string;
  trial_start?: string;
  trial_end?: string;
  created_at: string;
  updated_at: string;
  max_users: number;
  active_users: number;
  invited_users: number;
  is_trial_active: boolean;
  trial_days_remaining: number;
}

// Setup Progress Types
export interface CompanySetupProgress {
  company_id: string;
  step: string;
  progress: number;
  domain_provided: boolean;
  customization_completed: boolean;
  invitations_sent: boolean;
  subscription_started: boolean;
  setup_completed: boolean;
  last_updated: string;
}

// Company Stats Types
export interface CompanyStatsResponse {
  total_users: number;
  active_users: number;
  invited_users: number;
  pending_invitations: number;
  max_users: number;
  setup_progress: number;
  configuration_status: Record<string, boolean>;
}

// Request Types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface CompanyCreateRequest {
  name: string;
  domain: string;
  color_theme: string;
}

export interface InviteUserRequest {
  emails: string[];
}

export interface UpdateCompanyRequest {
  name?: string;
  domain?: string;
  color_theme?: string;
  logo_url?: string;
}

export interface UpdateUserRequest {
  name?: string;
  role?: UserRole;
  is_active?: boolean;
}

export interface BrowserShortcutRequest {
  name: string;
  url: string;
  icon?: string;
  description?: string;
  order: number;
  category?: string;
}

export interface CompanySetupRequest {
  domain: string;
  color_theme: string;
  logo_url?: string;
}

export interface GenerateShortcutsRequest {
  domain: string;
}

export interface NudgeUsersRequest {
  user_ids: string[];
  message?: string;
}

export interface ConfigurationStatusRequest {
  feature: string;
  status: boolean;
}

export interface DownloadInfoResponse {
  download_url: string;
  version: string;
  release_date: string;
  file_size: string;
  supported_os: string[];
  installation_instructions: string;
}

// Pagination Types
export interface PaginationRequest {
  page: number;
  limit: number;
}

export interface PaginationResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

// User Context Types
export interface UserContext {
  user_id: string;
  email: string;
  company_id: string;
  role: string;
}

// Auth Types
export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

// UI Types
export interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

export interface SetupStep {
  id: string;
  title: string;
  description: string;
  completed: boolean;
  current: boolean;
}

// Configuration Features
export interface ConfigurationFeature {
  id: string;
  name: string;
  description: string;
  configured: boolean;
  required: boolean;
}

// Chart Data Types
export interface ChartData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    backgroundColor: string[];
    borderColor: string[];
    borderWidth: number;
  }[];
}

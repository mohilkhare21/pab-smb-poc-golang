import React, { useEffect, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import apiService from '../services/api';

// Simple types for Dashboard
interface CompanyStatsResponse {
  total_users: number;
  active_users: number;
  invited_users: number;
  pending_invitations: number;
  max_users: number;
  setup_progress: number;
  configuration_status: Record<string, boolean>;
}

interface CompanySetupProgress {
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
import {
  Users,
  Globe,
  Shield,
  Download,
  CheckCircle,
  Clock,
  AlertCircle,
  TrendingUp,
  Building2,
  Settings
} from 'lucide-react';

const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const [stats, setStats] = useState<CompanyStatsResponse | null>(null);
  const [setupProgress, setSetupProgress] = useState<CompanySetupProgress | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        // Mock data for now - replace with actual API calls later
        const mockStats: CompanyStatsResponse = {
          total_users: 5,
          active_users: 3,
          invited_users: 2,
          pending_invitations: 1,
          max_users: 50,
          setup_progress: 75,
          configuration_status: {
            website_security: true,
            malware_protection: false,
            data_controls: true,
            reporting: false
          }
        };

        const mockProgress: CompanySetupProgress = {
          company_id: 'company-123',
          step: 'invitations',
          progress: 75,
          domain_provided: true,
          customization_completed: true,
          invitations_sent: true,
          subscription_started: false,
          setup_completed: false,
          last_updated: new Date().toISOString()
        };

        setStats(mockStats);
        setSetupProgress(mockProgress);
      } catch (error) {
        console.error('Error fetching dashboard data:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  const getProgressColor = (progress: number) => {
    if (progress >= 80) return 'text-green-600';
    if (progress >= 60) return 'text-yellow-600';
    if (progress >= 40) return 'text-orange-600';
    return 'text-red-600';
  };

  const getProgressBarColor = (progress: number) => {
    if (progress >= 80) return 'bg-green-600';
    if (progress >= 60) return 'bg-yellow-600';
    if (progress >= 40) return 'bg-orange-600';
    return 'bg-red-600';
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Welcome Header */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              Welcome back, {user?.name}!
            </h1>
            <p className="text-gray-600 mt-1">
              Here's what's happening with your company setup
            </p>
          </div>
          <div className="flex items-center space-x-2">
            <div className="h-3 w-3 bg-green-400 rounded-full animate-pulse"></div>
            <span className="text-sm text-gray-500">System Online</span>
          </div>
        </div>
      </div>

      {/* Setup Progress */}
      {setupProgress && (
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900">Setup Progress</h2>
            <span className={`text-2xl font-bold ${getProgressColor(setupProgress.progress)}`}>
              {setupProgress.progress}%
            </span>
          </div>
          
          <div className="mb-4">
            <div className="w-full bg-gray-200 rounded-full h-3">
              <div
                className={`h-3 rounded-full transition-all duration-500 ${getProgressBarColor(setupProgress.progress)}`}
                style={{ width: `${setupProgress.progress}%` }}
              ></div>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="flex items-center space-x-3">
              {setupProgress.domain_provided ? (
                <CheckCircle className="h-5 w-5 text-green-500" />
              ) : (
                <Clock className="h-5 w-5 text-gray-400" />
              )}
              <span className="text-sm text-gray-700">Domain Setup</span>
            </div>
            
            <div className="flex items-center space-x-3">
              {setupProgress.customization_completed ? (
                <CheckCircle className="h-5 w-5 text-green-500" />
              ) : (
                <Clock className="h-5 w-5 text-gray-400" />
              )}
              <span className="text-sm text-gray-700">Customization</span>
            </div>
            
            <div className="flex items-center space-x-3">
              {setupProgress.invitations_sent ? (
                <CheckCircle className="h-5 w-5 text-green-500" />
              ) : (
                <Clock className="h-5 w-5 text-gray-400" />
              )}
              <span className="text-sm text-gray-700">User Invitations</span>
            </div>
            
            <div className="flex items-center space-x-3">
              {setupProgress.subscription_started ? (
                <CheckCircle className="h-5 w-5 text-green-500" />
              ) : (
                <Clock className="h-5 w-5 text-gray-400" />
              )}
              <span className="text-sm text-gray-700">Subscription</span>
            </div>
          </div>
        </div>
      )}

      {/* Key Metrics */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <Users className="h-8 w-8 text-blue-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-500">Total Users</p>
                <p className="text-2xl font-semibold text-gray-900">{stats.total_users}</p>
              </div>
            </div>
            <div className="mt-4">
              <div className="flex items-center text-sm">
                <span className="text-gray-500">Active: {stats.active_users}</span>
                <span className="mx-2 text-gray-300">|</span>
                <span className="text-gray-500">Invited: {stats.invited_users}</span>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <Globe className="h-8 w-8 text-green-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-500">User Limit</p>
                <p className="text-2xl font-semibold text-gray-900">{stats.max_users}</p>
              </div>
            </div>
            <div className="mt-4">
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-green-600 h-2 rounded-full"
                  style={{ width: `${(stats.total_users / stats.max_users) * 100}%` }}
                ></div>
              </div>
              <p className="text-xs text-gray-500 mt-1">
                {Math.round((stats.total_users / stats.max_users) * 100)}% capacity used
              </p>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <AlertCircle className="h-8 w-8 text-orange-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-500">Pending Invitations</p>
                <p className="text-2xl font-semibold text-gray-900">{stats.pending_invitations}</p>
              </div>
            </div>
            <div className="mt-4">
              {stats.pending_invitations > 0 ? (
                <p className="text-sm text-orange-600">Action required</p>
              ) : (
                <p className="text-sm text-green-600">All caught up</p>
              )}
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <TrendingUp className="h-8 w-8 text-purple-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-500">Setup Progress</p>
                <p className="text-2xl font-semibold text-gray-900">{stats.setup_progress}%</p>
              </div>
            </div>
            <div className="mt-4">
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-purple-600 h-2 rounded-full transition-all duration-500"
                  style={{ width: `${stats.setup_progress}%` }}
                ></div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Quick Actions */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <button className="flex items-center justify-center px-4 py-3 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
            <Users className="h-5 w-5 mr-2" />
            Invite Users
          </button>
          
          <button className="flex items-center justify-center px-4 py-3 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
            <Globe className="h-5 w-5 mr-2" />
            Manage Shortcuts
          </button>
          
          <button className="flex items-center justify-center px-4 py-3 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
            <Settings className="h-5 w-5 mr-2" />
            Configuration
          </button>
        </div>
      </div>

      {/* Configuration Status */}
      {stats && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Configuration Status</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {Object.entries(stats.configuration_status).map(([feature, configured]) => (
              <div key={feature} className="flex items-center space-x-3">
                {configured ? (
                  <CheckCircle className="h-5 w-5 text-green-500" />
                ) : (
                  <Clock className="h-5 w-5 text-gray-400" />
                )}
                <span className="text-sm text-gray-700 capitalize">
                  {feature.replace(/_/g, ' ')}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;

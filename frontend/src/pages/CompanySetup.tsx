import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import apiService from '../services/api';

// Simple types for CompanySetup
interface Company {
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
  Building2,
  Globe,
  Palette,
  Users,
  CreditCard,
  CheckCircle,
  ArrowRight,
  ArrowLeft,
  Loader
} from 'lucide-react';

interface SetupStep {
  id: string;
  title: string;
  description: string;
  icon: React.ComponentType<any>;
  completed: boolean;
  current: boolean;
}

const CompanySetup: React.FC = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [currentStep, setCurrentStep] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [company, setCompany] = useState<Company | null>(null);
  const [setupProgress, setSetupProgress] = useState<CompanySetupProgress | null>(null);

  const [formData, setFormData] = useState({
    domain: '',
    colorTheme: '#3B82F6',
    logoUrl: '',
  });

  const steps: SetupStep[] = [
    {
      id: 'domain',
      title: 'Company Domain',
      description: 'Set up your company domain and basic information',
      icon: Globe,
      completed: false,
      current: true,
    },
    {
      id: 'customization',
      title: 'Browser Customization',
      description: 'Choose colors and branding for your custom browser',
      icon: Palette,
      completed: false,
      current: false,
    },
    {
      id: 'invitations',
      title: 'Invite Team Members',
      description: 'Send invitations to your team members',
      icon: Users,
      completed: false,
      current: false,
    },
    {
      id: 'subscription',
      title: 'Start Free Trial',
      description: 'Begin your 14-day free trial',
      icon: CreditCard,
      completed: false,
      current: false,
    },
  ];

  useEffect(() => {
    fetchCompanyData();
  }, []);

  const fetchCompanyData = async () => {
    try {
      const [companyResponse, progressResponse] = await Promise.all([
        apiService.getCompany(),
        apiService.getSetupProgress()
      ]);

      if (companyResponse.success && companyResponse.data) {
        setCompany(companyResponse.data);
        setFormData({
          domain: companyResponse.data.domain || '',
          colorTheme: companyResponse.data.color_theme || '#3B82F6',
          logoUrl: companyResponse.data.logo_url || '',
        });
      }

      if (progressResponse.success && progressResponse.data) {
        setSetupProgress(progressResponse.data);
        updateStepsFromProgress(progressResponse.data);
      }
    } catch (error) {
      console.error('Error fetching company data:', error);
    }
  };

  const updateStepsFromProgress = (progress: CompanySetupProgress) => {
    const updatedSteps = steps.map((step, index) => {
      let completed = false;
      let current = false;

      switch (step.id) {
        case 'domain':
          completed = progress.domain_provided;
          current = !progress.domain_provided;
          break;
        case 'customization':
          completed = progress.customization_completed;
          current = progress.domain_provided && !progress.customization_completed;
          break;
        case 'invitations':
          completed = progress.invitations_sent;
          current = progress.customization_completed && !progress.invitations_sent;
          break;
        case 'subscription':
          completed = progress.subscription_started;
          current = progress.invitations_sent && !progress.subscription_started;
          break;
      }

      return { ...step, completed, current };
    });

    // Find current step index
    const currentStepIndex = updatedSteps.findIndex(step => step.current);
    if (currentStepIndex !== -1) {
      setCurrentStep(currentStepIndex);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleColorChange = (color: string) => {
    setFormData({
      ...formData,
      colorTheme: color,
    });
  };

  const handleNext = async () => {
    if (currentStep === steps.length - 1) {
      // Complete setup
      await completeSetup();
      return;
    }

    // Save current step
    await saveCurrentStep();
    setCurrentStep(currentStep + 1);
  };

  const handlePrevious = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const saveCurrentStep = async () => {
    try {
      setIsLoading(true);
      
      switch (currentStep) {
        case 0: // Domain
          if (company) {
            await apiService.updateCompany({
              domain: formData.domain,
              color_theme: formData.colorTheme,
              logo_url: formData.logoUrl,
            });
            await apiService.updateSetupStep('domain', 25);
          }
          break;
        case 1: // Customization
          if (company) {
            await apiService.updateCompany({
              color_theme: formData.colorTheme,
              logo_url: formData.logoUrl,
            });
            await apiService.updateSetupStep('customization', 50);
          }
          break;
        case 2: // Invitations
          await apiService.updateSetupStep('invitations', 75);
          break;
        case 3: // Subscription
          await apiService.updateSetupStep('subscription', 100);
          break;
      }

      // Refresh progress
      await fetchCompanyData();
    } catch (error) {
      console.error('Error saving step:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const completeSetup = async () => {
    try {
      setIsLoading(true);
      await apiService.updateSetupStep('complete', 100);
      navigate('/dashboard');
    } catch (error) {
      console.error('Error completing setup:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const renderStepContent = () => {
    switch (currentStep) {
      case 0:
        return (
          <div className="space-y-6">
            <div>
              <label htmlFor="domain" className="block text-sm font-medium text-gray-700">
                Company Domain
              </label>
              <div className="mt-1">
                <input
                  type="text"
                  name="domain"
                  id="domain"
                  value={formData.domain}
                  onChange={handleInputChange}
                  className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                  placeholder="example.com"
                />
              </div>
              <p className="mt-2 text-sm text-gray-500">
                This will be used to generate suggested browser shortcuts for your team.
              </p>
            </div>

            <div>
              <label htmlFor="logoUrl" className="block text-sm font-medium text-gray-700">
                Company Logo URL (Optional)
              </label>
              <div className="mt-1">
                <input
                  type="url"
                  name="logoUrl"
                  id="logoUrl"
                  value={formData.logoUrl}
                  onChange={handleInputChange}
                  className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                  placeholder="https://example.com/logo.png"
                />
              </div>
            </div>
          </div>
        );

      case 1:
        return (
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-3">
                Browser Color Theme
              </label>
              <div className="grid grid-cols-6 gap-3">
                {['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#EC4899'].map((color) => (
                  <button
                    key={color}
                    onClick={() => handleColorChange(color)}
                    className={`w-12 h-12 rounded-full border-2 ${
                      formData.colorTheme === color ? 'border-gray-900' : 'border-gray-300'
                    }`}
                    style={{ backgroundColor: color }}
                  />
                ))}
              </div>
              <div className="mt-3">
                <input
                  type="color"
                  value={formData.colorTheme}
                  onChange={(e) => handleColorChange(e.target.value)}
                  className="h-10 w-20 rounded border border-gray-300"
                />
                <span className="ml-3 text-sm text-gray-500">{formData.colorTheme}</span>
              </div>
            </div>

            <div>
              <label htmlFor="logoUrl" className="block text-sm font-medium text-gray-700">
                Company Logo URL
              </label>
              <div className="mt-1">
                <input
                  type="url"
                  name="logoUrl"
                  id="logoUrl"
                  value={formData.logoUrl}
                  onChange={handleInputChange}
                  className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                  placeholder="https://example.com/logo.png"
                />
              </div>
            </div>
          </div>
        );

      case 2:
        return (
          <div className="space-y-6">
            <div className="text-center">
              <Users className="mx-auto h-16 w-16 text-gray-400" />
              <h3 className="mt-4 text-lg font-medium text-gray-900">Invite Team Members</h3>
              <p className="mt-2 text-sm text-gray-500">
                You can invite team members from the User Management section after completing setup.
              </p>
            </div>
            <div className="bg-blue-50 border border-blue-200 rounded-md p-4">
              <p className="text-sm text-blue-700">
                <strong>Tip:</strong> You can invite up to {company?.subscription_id ? 'your plan limit' : '20 users'} 
                during your free trial. Additional users can be added by upgrading your subscription.
              </p>
            </div>
          </div>
        );

      case 3:
        return (
          <div className="space-y-6">
            <div className="text-center">
              <CreditCard className="mx-auto h-16 w-16 text-gray-400" />
              <h3 className="mt-4 text-lg font-medium text-gray-900">Start Your Free Trial</h3>
              <p className="mt-2 text-sm text-gray-500">
                Your 14-day free trial will begin immediately. No credit card required to start.
              </p>
            </div>
            <div className="bg-green-50 border border-green-200 rounded-md p-4">
              <h4 className="text-sm font-medium text-green-800">What's included in your trial:</h4>
              <ul className="mt-2 text-sm text-green-700 space-y-1">
                <li>• Up to 20 team members</li>
                <li>• Custom browser branding</li>
                <li>• Browser shortcut management</li>
                <li>• User invitation system</li>
                <li>• Basic security features</li>
              </ul>
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      {/* Progress Steps */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          {steps.map((step, index) => {
            const Icon = step.icon;
            return (
              <div key={step.id} className="flex items-center">
                <div className="flex items-center">
                  <div
                    className={`flex items-center justify-center w-10 h-10 rounded-full border-2 ${
                      step.completed
                        ? 'bg-green-600 border-green-600 text-white'
                        : step.current
                        ? 'bg-blue-600 border-blue-600 text-white'
                        : 'bg-gray-200 border-gray-300 text-gray-500'
                    }`}
                  >
                    {step.completed ? (
                      <CheckCircle className="h-5 w-5" />
                    ) : (
                      <Icon className="h-5 w-5" />
                    )}
                  </div>
                  <div className="ml-3">
                    <h3 className="text-sm font-medium text-gray-900">{step.title}</h3>
                    <p className="text-xs text-gray-500">{step.description}</p>
                  </div>
                </div>
                {index < steps.length - 1 && (
                  <div
                    className={`w-16 h-0.5 mx-4 ${
                      step.completed ? 'bg-green-600' : 'bg-gray-300'
                    }`}
                  />
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Step Content */}
      <div className="bg-white rounded-lg shadow p-8">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-gray-900">{steps[currentStep].title}</h2>
          <p className="text-gray-600 mt-1">{steps[currentStep].description}</p>
        </div>

        {renderStepContent()}

        {/* Navigation */}
        <div className="flex items-center justify-between mt-8 pt-6 border-t border-gray-200">
          <button
            onClick={handlePrevious}
            disabled={currentStep === 0}
            className="flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Previous
          </button>

          <button
            onClick={handleNext}
            disabled={isLoading}
            className="flex items-center px-6 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
          >
            {isLoading ? (
              <Loader className="h-4 w-4 mr-2 animate-spin" />
            ) : currentStep === steps.length - 1 ? (
              'Complete Setup'
            ) : (
              <>
                Next
                <ArrowRight className="h-4 w-4 ml-2" />
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default CompanySetup;

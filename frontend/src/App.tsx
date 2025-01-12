import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import './App.css'
import RegisterForm from './pages/Register'
import LoginForm from './pages/Login'
import HomePage from './pages/Home'
import CreatePage from './pages/Create'
import { AuthProvider } from './pages/AuthContext'
import AccountPage from './pages/AccountPage'
import { NonExistantRoute } from './pages/NonExistantRoute'
import DownloadMobileApk from './pages/ClientMobile'
import ChangePasswordPage from './pages/ChangePassword'
import MyWorkflow from './pages/MyWorkflows'
import ProtectedRoute from './ProtectedRoute'
import RedirectRoute from './RedirectRoute'

function App() {
  return (
    <AuthProvider>
        <Router>
          <Routes>
              <Route
                path="/register"
                element={
                  <RegisterForm />
                }
              />

              <Route
                path="/login"
                element={
                  <LoginForm />
                }
              />

              <Route
                path="/client.apk"
                element={
                  <DownloadMobileApk />
                }
              />

              <Route
                path="/"
                element={
                  <RedirectRoute dest='/explore'/>
                }
              />

              <Route
                path="/explore"
                element={
                  <HomePage />
                }
              />

              <Route
                path="/my_workflows"
                element={
                  <ProtectedRoute>
                    <MyWorkflow />
                  </ProtectedRoute>
                }
              />

              <Route
                path="/create"
                element={
                  <ProtectedRoute>
                    <CreatePage />
                  </ProtectedRoute>
                }
              />

              <Route
                path="/settings"
                element={
                  <ProtectedRoute>
                    <AccountPage />
                  </ProtectedRoute>
                }
              />

              <Route
                path="/settings/change_password"
                element={
                  <ProtectedRoute>
                    <ChangePasswordPage />
                  </ProtectedRoute>
                }
              />

              <Route
                path="*"
                element={
                  <NonExistantRoute/>
                }
              />
          </Routes>
        </Router>
    </AuthProvider>
  )
}

export default App

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

function App() {
  return (
    <AuthProvider>
        <Router>
          <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/create" element={<CreatePage />} />
              <Route path="/register" element={<RegisterForm />} />
              <Route path="/login" element={<LoginForm />} />
              <Route path="/settings" element={<AccountPage />} />
              <Route path="/settings/change_password" element={<ChangePasswordPage />} />
              <Route path="/client.apk" element={<DownloadMobileApk />} />
              <Route path="*" element={<NonExistantRoute/>} />
          </Routes>
        </Router>
    </AuthProvider>
  )
}

export default App

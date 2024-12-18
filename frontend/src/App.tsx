import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import './App.css'
import RegisterForm from './pages/Register'
import LoginForm from './pages/Login'
import HomePage from './pages/Home'
import CreatePage from './pages/Create'
import SelectServicePage from './pages/SelectService'
import ServicePage from './pages/ServicePage'
import { AuthProvider } from './pages/AuthContext'
import ConnectService from './pages/ConnectService'
import TriggerActionPage from './pages/TriggerAction'
import AccountPage from './pages/AccountPage'
import DownloadMobileApk from './pages/DownloadMobileApk'

function App() {
  return (
    <AuthProvider>
        <Router>
          <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/create" element={<CreatePage />} />
              <Route path="/register" element={<RegisterForm />} />
              <Route path="/login" element={<LoginForm />} />
              <Route path="/account" element={<AccountPage />} />
              <Route path="/select-service" element={<SelectServicePage />} />
              <Route path="/select-service/:name" element={<ServicePage />} />
              <Route path="/select-action/:name" element={<TriggerActionPage />} />
              <Route path="/connect-service/:name" element={<ConnectService />} />
              <Route path="/client.apk" element={<DownloadMobileApk />} />
          </Routes>
        </Router>
    </AuthProvider>
  )
}

export default App

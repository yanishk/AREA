import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import RegisterForm from './pages/Register'
import LoginForm from './pages/Login'
import HomePage from './pages/Home'
import CreatePage from './pages/Create'
import SelectActionPage from './pages/SelectAction'
import SelectReactionPage from './pages/SelectReaction'
import SelectServicePage from './pages/SelectService'
import ServicePage from './pages/ServicePage'
import { AuthProvider } from './pages/AuthContext'

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/create" element={<CreatePage />} />
          <Route path="/register" element={<RegisterForm />} />
          <Route path="/login" element={<LoginForm />} />
          <Route path="/select-service" element={<SelectServicePage />} />
          <Route path="/select-service/:id" element={<ServicePage />} />
          <Route path="/select-action" element={<SelectActionPage />} />
          <Route path="/select-reaction" element={<SelectReactionPage />} />
        </Routes>
      </Router>
    </AuthProvider>
  )
}

export default App

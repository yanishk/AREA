import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import './App.css'
import RegisterForm from './pages/Register'
import LoginForm from './pages/Login'
import HomePage from './pages/Home'
import CreatePage from './pages/Create'
import SelectActionPage from './pages/SelectAction'

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/create" element={<CreatePage />} />
        <Route path="/register" element={<RegisterForm />} />
        <Route path="/login" element={<LoginForm />} />
        <Route path="/select-action" element={<SelectActionPage />} />
      </Routes>
    </Router>
  )
}

export default App

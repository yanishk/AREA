import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import './App.css'
import RegisterForm from './pages/Register'
import LoginForm from './pages/Login'
import HomePage from './pages/Home'

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/register" element={<RegisterForm />} />
        <Route path="/login" element={<LoginForm />} />
      </Routes>
    </Router>
  )
}

export default App

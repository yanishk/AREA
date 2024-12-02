import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "./App.css";
import RegisterForm from "./pages/register";
import LoginForm from "./pages/login";

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/register" element={<RegisterForm />} />
                <Route path="/login" element={<LoginForm />} />
            </Routes>
        </Router>
    );
}

export default App;

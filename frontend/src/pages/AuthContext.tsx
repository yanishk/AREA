import axios from "axios";
import { createContext, useContext, useState, ReactNode, useEffect, useMemo } from "react";

type AuthContextType = {
    isAuthenticated: boolean;
    login: () => void;
    logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    const login = () => {
        setIsAuthenticated(true);
    };

    const logout = () => {
        setIsAuthenticated(false);
    };

    const checkAuthStatus = async () => {
        try {
            const response = await axios.get(`${import.meta.env.VITE_API_URL}user/email`, {
                withCredentials: true,
            });

            if (response.status === 200) {
                setIsAuthenticated(true);
            } else {
                setIsAuthenticated(false);
            }
        } catch (error) {
            console.error("error: ", error);
            setIsAuthenticated(false);
        }
    };

    useEffect(() => {
        checkAuthStatus();
    }, []);

    return (
        <AuthContext.Provider value={useMemo(() => ({ isAuthenticated, login, logout }), [isAuthenticated, login, logout])}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("error with AuthProvider");
    }
    return context;
};

import { createContext, useContext, useState, ReactNode, useEffect, useMemo } from "react";
import isAuthenticated from "../utils/checkAuthStatus";

type AuthContextType = {
    isUserConnected: boolean;
    login: () => void;
    logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [isUserConnected, setIsUserConnected] = useState(false);

    const checkAuthStatus = async () => {
        const authStatus = await isAuthenticated();
        setIsUserConnected(authStatus);
    };

    const login = () => {
        setIsUserConnected(true);
    };

    const logout = () => {
        setIsUserConnected(false);
    };

    useEffect(() => {
        checkAuthStatus();
    }, []);

    return (
        <AuthContext.Provider value={useMemo(() => ({ isUserConnected, login, logout }), [isUserConnected])}>
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

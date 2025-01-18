import { createContext, useContext, useState, ReactNode, useEffect, useMemo } from "react";
import isAuthenticated from "../utils/CheckAuthStatus";

type AuthContextType = {
    isUserConnected: boolean;
    isLoading: boolean;
    login: () => void;
    logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [isUserConnected, setIsUserConnected] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    const checkAuthStatus = async () => {
        const authStatus = await isAuthenticated();
        setIsUserConnected(authStatus);
        setIsLoading(false);
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
        <AuthContext.Provider value={useMemo(() => ({ isUserConnected, login, logout, isLoading }), [isUserConnected, isLoading])}>
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

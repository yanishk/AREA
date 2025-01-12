import React from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "./pages/AuthContext";

interface ProtectedRouteProps {
    children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
    const { isUserConnected, isLoading } = useAuth();

    if (isLoading) {
        return;
    }

    if (!isUserConnected) {
        return <Navigate to="/login" replace />;
    }

    return <>
        {children}
    </>;
};

export default ProtectedRoute;

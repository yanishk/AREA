import React from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "./pages/AuthContext";

interface RedirectRouteProps {
    children: React.ReactNode;
    dest: string,
}

const RedirectRoute: React.FC<RedirectRouteProps> = ({ children, dest }) => {
    const { isAuthenticated } = useAuth();

    if (isAuthenticated) {
        return <Navigate to={`${dest}`} replace />;
    }

    return <>
        {children}
    </>;
};

export default RedirectRoute;

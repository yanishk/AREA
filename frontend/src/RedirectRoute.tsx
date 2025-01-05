import React from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "./pages/AuthContext";

interface RedirectRouteProps {
    children: React.ReactNode;
    dest: string,
}

const RedirectRoute: React.FC<RedirectRouteProps> = ({ children, dest }) => {
    const { isUserConnected } = useAuth();

    if (!isUserConnected) {
        return <>
            {children}
        </>;
    }

    return <Navigate to={`${dest}`} replace />;
};

export default RedirectRoute;

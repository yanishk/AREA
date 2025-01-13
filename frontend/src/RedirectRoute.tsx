import React from "react";
import { Navigate } from "react-router-dom";

interface RedirectRouteProps {
    dest: string,
}

const RedirectRoute: React.FC<RedirectRouteProps> = ({ dest }) => {
    return <Navigate to={`${dest}`} replace />;
};

export default RedirectRoute;

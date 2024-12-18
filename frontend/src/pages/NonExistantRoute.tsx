import { Navigate } from 'react-router-dom';

function NonExistantRoute() {
    const error = "Service doesn't exist or this account doesn't have access to it";
    const errorTrigger = Date.now();

    return (
        <Navigate
            to="/"
            state={{ error, errorTrigger }}
        />
    )
}

export { NonExistantRoute }

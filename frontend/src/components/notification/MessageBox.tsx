import React, { useEffect, useState } from "react";

interface ErrorMessageProps {
    message: string;
    type: string;
    trigger: number;
    timeout: number;
    bgcolor?: string;
    onClose: () => void;
}

const MessageBox: React.FC<ErrorMessageProps> = ({
    message,
    trigger,
    onClose,
    timeout,
    type,
    bgcolor = type === "error" ? "bg-[#D0011B]" : "bg-[#2CBE60]"
}) => {
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        if (message) {
            setIsVisible(true);

            const timer = setTimeout(() => {
                setIsVisible(false);
                onClose?.();
            }, timeout);

            return () => clearTimeout(timer);
        }
    }, [message, trigger]);

    if (!message) {
        return null;
    }

    return (
        <div
            className={`fixed top-0 left-0 w-full p-5 ${bgcolor} text-white text-center font-bold text-[30px] z-50 transition-transform duration-500 ${
                isVisible ? "translate-y-0" : "-translate-y-full"
            }`}
        >
            {message}
        </div>
    )
};

export default MessageBox;

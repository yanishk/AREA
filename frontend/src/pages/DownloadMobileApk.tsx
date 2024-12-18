import { useNavigate } from "react-router-dom";
import { useEffect, useRef } from "react";

const DownloadMobileApk = () => {
    const navigate = useNavigate();

    const downloadRef = useRef<HTMLAnchorElement>(null);

    useEffect(() => {
        if (downloadRef.current) {
            downloadRef.current.click();
        }

        const timer = setTimeout(() => {
            navigate(-1);
        }, 500);

        return () => clearTimeout(timer);
    }, [navigate]);

    return (
        <div className="flex items-center justify-center flex-col px-10 pt-[100px] py-5">
            <a href="/client.apk" download="client.apk" ref={downloadRef} style={{ display: "none" }}></a>
        </div>
    );
};

export default DownloadMobileApk
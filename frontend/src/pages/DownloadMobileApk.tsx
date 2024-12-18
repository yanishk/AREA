import { useNavigate } from "react-router-dom";
import { useEffect , useRef } from "react";

const DownloadMobileApk = () => {
    const navigate = useNavigate();

    const downloadRef = useRef<HTMLAnchorElement>(null);

    useEffect(() => {
        if (downloadRef.current) {
            downloadRef.current.click();
            const timer = setTimeout(() => {
                navigate(-1);
            }, 200);

            return () => clearTimeout(timer);
        }
    }, []);

    return (
        <div className="flex items-center justify-center flex-col px-10 pt-[100px] py-5">
            <h2>If the download has not started, tap download</h2>
            <a href="/area.apk" ref={downloadRef}>Download</a>
        </div>
    );
};

export default DownloadMobileApk
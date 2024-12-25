import { useNavigate } from "react-router-dom";
import { useEffect , useRef } from "react";
import Navbar from "../components/navbar/Navbar";
import TriggerButton from "../components/buttons/TriggerButton";

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

    const handleClick = () => {
        if (downloadRef.current) {
            downloadRef.current.click();
        }
    };

    return (
        <>
            <Navbar />
            <div className="flex items-center justify-center flex-col px-10 pt-[100px] py-5">
                <h2
                    className="font-[900] text-[55px]"
                >
                    If the download has not started, tap download
                </h2>
                <TriggerButton
                    text="Download"
                    textColor="text-white mt-[100px]"
                    bgColor="bg-[#222222]"
                    hoverColor="hover:bg-[#333333]"
                    onClick={handleClick}
                />
                <a
                    href="/area.apk"
                    ref={downloadRef}
                    className="hidden"
                >
                    Download
                </a>
            </div>
        </>
    );
};

export default DownloadMobileApk
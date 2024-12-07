import { useLocation, useNavigate } from "react-router-dom";
import CancelButton from "../components/buttons/CancelButton";

const SelectServicePage = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const goBack = () => {
        navigate(-1);
    };

    const type = location.state?.type;

    const services = [
        { id: 1, name: "Blogger", color: "#F48121", iconPath: "/images/Blogger.webp"},
        { id: 2, name: "Date & Time", color: "#333333", iconPath: "/images/DateLogo.webp"},
        { id: 3, name: "Faceebook Pages", color: "#3B579D", iconPath: "/images/Facebook.webp"},
        { id: 4, name: "Google Calendar", color: "#2C6EFC", iconPath: "/images/GoogleCalendar.webp"},
        { id: 5, name: "RSS Feed", color: "#EC7603", iconPath: "/images/RSS.webp"},
    ];

    const services2 = [
        { id: 1, name: "Spotify", color: "#1CB954", iconPath: "/images/Spotify.webp"},
        { id: 2, name: "Twitter", color: "#000000", iconPath: "/images/Twitter.webp"},
        { id: 3, name: "Weather", color: "#1E2024", iconPath: "/images/Weather.webp"},
        { id: 4, name: "Whebooks", color: "#0099FF", iconPath: "/images/Webhooks.webp"},
        { id: 5, name: "Youtube", color: "#FF0200", iconPath: "/images/Youtube.webp"},
    ];

    return (
        <div className="min-h-screen bg-gray-100">
            <div className="flex items-center justify-between px-10 py-5">
                <div>
                    <CancelButton
                        text="Cancel"
                        onClick={goBack}
                    />
                </div>

                <h1 className="text-[3.2rem] font-extrabold text-center flex-1">
                    Choose a service
                </h1>

                <div className="w-[100px]"></div>
            </div>

            <div className="mt-10 px-[140px]">
                <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
                    {services.map((service) => (
                        <div
                            key={service.id}
                            style={{ backgroundColor: service.color }}
                            className="h-[240px] flex flex-col items-center "
                        >
                            <img src={service.iconPath} className="size-28 mt-9" />
                            <h2 className="text-white text-lg text-[20.5px] font-bold mt-6">{service.name}</h2>
                        </div>
                    ))}
                </div>
            </div>

            <div className="px-[140px]">
                <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
                    {services2.map((service2) => (
                        <div
                            key={service2.id}
                            style={{ backgroundColor: service2.color }}
                            className="h-[240px] flex flex-col items-center "
                        >
                            <img src={service2.iconPath} className="size-28 mt-9" />
                            <h2 className="text-white text-lg text-[20.5px] font-bold mt-6">{service2.name}</h2>
                        </div>
                    ))}
                </div>
            </div>

        </div>
    );
};

export default SelectServicePage;

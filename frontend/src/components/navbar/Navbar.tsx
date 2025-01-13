import NavbarLink from "./NavbarLink"
import NavbarButton from "./NavbarButton"
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../pages/AuthContext";
import axios from "axios";

interface NavbarProps {
    bgColor?: string;
    isWhiteMode?: boolean;
}

const Navbar: React.FC<NavbarProps> = ({ bgColor, isWhiteMode }) => {
    const { isUserConnected } = useAuth();
    const navigate = useNavigate();

    const handleGetStarted = () => {
        navigate("/register");
    }

    const handleAccount = () => {
        navigate("/settings");
    };

    const handleCreate = () => {
        navigate("/create");
    }

    const handleLogout = async () => {
        try {
            await axios.post(`${import.meta.env.VITE_API_URL}logout`, {}, {
                withCredentials: true
            });
            navigate("/explore");
            window.location.reload();
        } catch (error) {
            console.error("Error during logout:", error);
        }
    }

    return (
        <div className={`${bgColor}`}>
            <div className="pt-6 pb-6 flex justify-between items-center px-6">
                {isWhiteMode ? (
                    <h1
                        className="text-4xl text-white hover:opacity-90 font-extrabold ml-4">
                            <a
                                href="/explore"
                            >
                                AREA
                            </a>
                    </h1>
                ) : (
                    <h1
                        className="text-4xl text-[#222222] hover:text-[#333333] font-extrabold ml-4">
                            <a
                                href="/explore"
                            >
                                AREA
                            </a>
                    </h1>
                )}

                <ul
                    className="flex flex-row items-center gap-x-6 mr-4"
                >
                    {!isUserConnected ? (
                        <>
                            <li>
                                <NavbarLink
                                    href="/explore"
                                    text="Explore"
                                    textColor="text-[#222222]"
                                    hoverColor="hover:text-[#444444]"
                                />
                            </li>
                            <li>
                                <NavbarLink
                                    href="/login"
                                    text="Log in"
                                    textColor="text-[#222222]"
                                    hoverColor="hover:text-[#444444]"
                                   />
                                </li>
                            <li>
                                <NavbarButton
                                    onClick={handleGetStarted}
                                    text="Get Started"
                                 />
                            </li>
                        </>
                    ) : (
                        <>
                            {isWhiteMode ? (
                                <>
                                    <li>
                                        <NavbarLink
                                            href="/explore"
                                            text="Explore"
                                            textColor="text-white"
                                            hoverColor="hover:opacity-90"
                                        />
                                    </li>
                                    <li>
                                        <NavbarLink
                                            href="/my_workflows"
                                            text="My workflows"
                                            textColor="text-white"
                                            hoverColor="hover:opacity-90"
                                        />
                                    </li>
                                    <li>
                                        <NavbarButton
                                            onClick={handleCreate}
                                            text="Create"
                                            isWhiteMode={true}
                                        />
                                    </li>
                                    <li>
                                        <button
                                            onClick={handleLogout}
                                            className="text-white font-bold text-[20px]"
                                        >
                                            Logout
                                        </button>
                                    </li>
                                    <li>
                                        <button
                                            className="bg-white hover:opacity-90 w-[55px] h-[55px] rounded-[30px] flex flex-col items-center justify-center text-center border-[3px] border-white cursor-pointer"
                                            onClick={handleAccount}
                                        >
                                            <img
                                                src="/iconAvatar.svg"
                                                alt="Avatar"
                                            />
                                        </button>
                                    </li>
                                </>
                            ) : (
                                <>
                                    <li>
                                        <NavbarLink
                                            href="/explore"
                                            text="Explore"
                                            textColor="text-[#222222]"
                                            hoverColor="hover:text-[#444444]"
                                        />
                                    </li>
                                    <li>
                                        <NavbarLink
                                            href="/my_workflows"
                                            text="My workflows"
                                            textColor="text-[#222222]"
                                            hoverColor="hover:text-[#444444]"
                                        />
                                    </li>
                                    <li>
                                        <NavbarButton
                                            onClick={handleCreate}
                                            text="Create"
                                        />
                                    </li>
                                    <li>
                                        <button
                                            onClick={handleLogout}
                                            className="text-[#222222] font-bold text-[20px]"
                                        >
                                            Logout
                                        </button>
                                    </li>
                                    <li>
                                        <button
                                            className="bg-[#222222] hover:bg-[#333333] w-[55px] h-[55px] rounded-[30px] flex flex-col items-center justify-center text-center border-[3px] border-[#222222] cursor-pointer"
                                            onClick={handleAccount}
                                        >
                                            <img
                                                src="/iconAvatar.svg"
                                                alt="Avatar"
                                            />
                                        </button>
                                    </li>
                                </>
                            )}
                        </>
                    )}
                </ul>

            </div>
        </div>
    )
}

export default Navbar
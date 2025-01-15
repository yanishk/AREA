import NavbarLink from "./NavbarLink"
import NavbarButton from "./NavbarButton"
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../pages/AuthContext";
import axios from "axios";
import {
    TEDropdown,
    TEDropdownToggle,
    TEDropdownMenu,
    TEDropdownItem,
} from "tw-elements-react";

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
                        className="text-white text-4xl hover:opacity-90 font-extrabold ml-4">
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
                    className="lg:flex flex-row items-center gap-x-6 mr-4 hidden lg:visible"
                >
                    {!isUserConnected ? (
                        <>
                            {!isWhiteMode ? (
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
                                            href="/login"
                                            text="Log in"
                                            textColor="text-white"
                                            hoverColor="hover:opacity-90"
                                        />
                                    </li>
                                    <li>
                                        <NavbarButton
                                            onClick={handleGetStarted}
                                            text="Get Started"
                                            isWhiteMode={true}
                                        />
                                    </li>
                                </>
                            )}
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
                                            Log out
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
                                            Log out
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

                <TEDropdown className="pt-6 pb-6 flex-1 flex justify-end items-center px-6 lg:hidden">
                    <TEDropdownToggle className="flex items-center px-6 pb-2 pt-2.5 text-xs font-medium leading-normal border-2 border-[#DEE3ED] rounded-lg bg-white">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
                        </svg>
                        <TEDropdownMenu className={`flex flex-col gap-y-4 text-sm pl-3 pr-3 items-start border-gray-200 border-2 py-2`}>
                            {!isUserConnected ? (
                                <>
                                    <TEDropdownItem>
                                        <NavbarLink
                                            href="/explore"
                                            text="Explore"
                                            textColor="text-[#222222]"
                                            hoverColor="hover:text-[#444444]"
                                        />
                                    </TEDropdownItem>
                                    <TEDropdownItem>
                                        <NavbarLink
                                            href="/login"
                                            text="Log in"
                                            textColor="text-[#222222]"
                                            hoverColor="hover:text-[#444444]"
                                        />
                                        </TEDropdownItem>
                                    <TEDropdownItem>
                                        <NavbarLink
                                            href="/register"
                                            text="Get Started"
                                            textColor="text-[#222222]"
                                            hoverColor="hover:text-[#444444]"
                                        />
                                    </TEDropdownItem>
                                </>
                            ) : (
                                <>
                                    <TEDropdownItem>
                                        <NavbarLink
                                            href="/explore"
                                            text="Explore"
                                            textColor="text-[#222222]"
                                            hoverColor="hover:text-[#444444]"
                                        />
                                    </TEDropdownItem>
                                    <TEDropdownItem>
                                        <NavbarLink
                                            href="/my_workflows"
                                            text="My workflows"
                                            textColor="text-[#222222]"
                                            hoverColor="hover:text-[#444444]"
                                        />
                                    </TEDropdownItem>
                                    <TEDropdownItem>
                                        <NavbarLink
                                            href="/create"
                                            text="Create"
                                            textColor="text-[#222222]"
                                            hoverColor="hover:text-[#444444]"
                                        />
                                    </TEDropdownItem>
                                    <TEDropdownItem>
                                        <button
                                            onClick={handleLogout}
                                            className="text-[#222222] font-bold text-[20px] ml-3"
                                        >
                                            Log out
                                        </button>
                                    </TEDropdownItem>
                                    <TEDropdownItem>
                                        <NavbarLink
                                            href="/settings"
                                            text="Account"
                                            textColor="text-[#222222]"
                                            hoverColor="hover:text-[#444444]"
                                        />
                                    </TEDropdownItem>
                                </>
                            )}
                        </TEDropdownMenu>
                    </TEDropdownToggle>
                </TEDropdown>
            </div>
        </div>
    )
}

export default Navbar
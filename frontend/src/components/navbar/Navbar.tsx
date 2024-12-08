import NavbarLink from "./NavbarLink"
import NavbarButton from "./NavbarButton"
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../pages/AuthContext";

const Navbar = () => {
    const { isAuthenticated, logout } = useAuth();
    const navigate = useNavigate();

    const handleGetStarted = () => {
        navigate("/register");
    }

    const handleLogout = () => {
        logout();
        navigate("/");
    };
    const handleCreate = () => {
        navigate("/create");
    }

    return (
        <>
            <div className="bg-white">
                <div className="pt-6 pb-6 flex justify-between items-center px-6">

                    <h1
                        className="text-4xl text-[#222222] hover:text-[#333333] font-extrabold ml-4">
                            <a
                                href="/"
                            >
                                AREA
                            </a>
                    </h1>

                    <ul
                        className="flex flex-row items-center gap-x-6 mr-4"
                    >
                        <li>
                            <NavbarLink
                                href="/"
                                text="Home"
                            />
                        </li>
                        <li>
                            <NavbarLink
                                href="/explore"
                                text="Explore"
                            />
                        </li>
                        {!isAuthenticated ? (
                            <>
                                <li>
                                    <NavbarLink
                                        href="/login"
                                        text="Log in"
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
                                    <NavbarButton
                                        onClick={handleCreate}
                                        text="Create"
                                    />
                                </li>
                                <li>
                                    <NavbarButton
                                        onClick={handleLogout}
                                        text="Log out"
                                    />
                                </li>
                            </>
                        )}
                    </ul>

                </div>
            </div>
        </>
    )
}

export default Navbar
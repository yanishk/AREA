import NavbarLink from "./NavbarLink"
import NavbarButton from "./NavbarButton"
import { useNavigate } from "react-router-dom";

const Navbar = () => {
    const navigate = useNavigate();

    const handleGetStarted = () => {
        navigate("/register");
    }

    return (
        <>
            <div className="bg-white h-[75px] flex justify-between items-center px-6">

                <h1
                    className="text-4xl font-extrabold ml-4">
                        AREA
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
                </ul>

            </div>
        </>
    )
}

export default Navbar
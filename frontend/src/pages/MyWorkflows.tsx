import MaxWidthWrapper from "../components/MaxWidthWrapper"
import Navbar from "../components/navbar/Navbar"

const MyWorkflow = () => {

    return (
        <>
            <Navbar />

            <MaxWidthWrapper
                maxWidth="1120px"
            >
                <div
                    className="flex items-center justify-center mt-10"
                >
                    <div
                        className="flex flex-col items-center text-center w-full"
                    >
                        <h1
                            className="text-5xl text-[#222222] font-[900]"
                        >
                            My Workflows
                        </h1>

                    </div>
                </div>
                <div className="mt-16 ml-[35px]">
                    <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3">
                        {/* {availableItems?.map((element) => ( */}
                            <button
                                // key={element.name}
                                // style={{ backgroundColor: "#333333"}}
                                className="h-[420px] flex flex-col items-start pt-2 cursor-pointer rounded-[10px] mr-6 mb-5 bg-[#333333] hover:bg-[#474747]"
                                onClick={() => ("")}
                            >
                                <div className="text-left w-full pl-[30px]">
                                    <h2 className="text-white text-lg text-[1.6em] font-[800] leading-[1.125] mt-6 m-0 p-0">aaa</h2>
                                    <h2 className="text-white text-lg text-[1.3em] font-[400] mt-6 mr-4">aaa</h2>
                                </div>
                            </button>

                        {/* ))} */}
                    </div>
                </div>
            </MaxWidthWrapper>
        </>
    )
}

export default MyWorkflow
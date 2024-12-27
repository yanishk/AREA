import { useNavigate, useLocation, useSearchParams } from "react-router-dom";
import ButtonIfCustom from "../components/buttons/ButtonIfCustom";
import { useEffect, useState } from "react";
import ButtonIfDefault from "../components/buttons/ButtonIfDefault";
import CancelButton from "../components/buttons/CancelButton";
import ButtonThenCustom from "../components/buttons/ButtonThenCustom";
import ButtonThenDefault from "../components/buttons/ButtonThenDefault";
import TriggerButton from "../components/buttons/TriggerButton";
import axios from "axios";

enum Steps {
    HOME,
    SELECT_SERVICE_ACTION,
    SELECT_SERVICE_REACTION,
    CONNECT_SERVICE,
    TRIGGER_ACTION,
    TRIGGER_REACTION,
    SUMMARY,
}

const WorkflowBuilder = () => {
    const navigate = useNavigate();

    const [selectedAction, setSelectedAction] = useState(null);
    const [selectedReaction, setSelectedReaction] = useState(null);
    const [services, setServices] = useState([]);

    const [currentStep, setCurrentStep] = useState<Steps>(Steps.HOME);

    const goToStep = (step: Steps) => {
        setCurrentStep(step);
    };

    const handleOnEditAction = () => {
        return;
    }

    const handleDeleteAction = () => {
        return;
    }

    const goHome = () => {
        return;
    }

    const handleContinue = () => {
        return;
    }

    const getServicesWithActions = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}actions-services`, {
                    withCredentials: true,
            });

            if (result) {
                const propertiesServices = result.data.map((service: any) => ({
                    name: service.name,
                    color: service.color,
                    iconPath: service.logo,
                }));
                setServices(propertiesServices);
            }
        } catch (err) {
            console.error("Error fetching services:", err);
        }
    };

    useEffect(() => {
        if (currentStep === Steps.SELECT_SERVICE_ACTION) {
            getServicesWithActions();
        }
    }, [currentStep]);

    const renderStep = () => {
        switch (currentStep) {
            case Steps.HOME:
                return (
                    <div>
                    <div className="flex items-center justify-between px-10 py-5">
                        <div>
                            <CancelButton
                                text="Cancel"
                                textColor="text-[#222222]"
                                borderColor="border-[#222222]"
                                bgColor="bg-[#F9FAFB]"
                                onClick={goHome}
                            />
                        </div>

                        <h1 className="text-[3.2rem] font-[900] text-center flex-1">
                            Create
                        </h1>

                        <div className="w-[140px]"></div>
                    </div>

                    <div className="flex items-center justify-center flex-col px-10 pt-[100px] py-5">
                        <div className="pt-10">
                            {selectedAction ? (
                                <ButtonIfDefault
                                    onClick={() =>
                                        goHome()
                                    }
                                    text="If This"
                                />
                            ) : (
                                <ButtonIfDefault
                                    onClick={() =>
                                        goToStep(Steps.SELECT_SERVICE_ACTION)
                                    }
                                    text="If This"
                                />
                            )}
                        </div>

                        <div className="h-[70px] w-[5px] bg-[#EEEEEE]" />

                        <div className="mb-[50px]">
                            {selectedReaction ? (
                                <ButtonThenDefault
                                onClick={() =>
                                    goHome()
                                }
                                text="Then That"
                                disabled={!selectedAction}
                            />
                            ) : (
                                <ButtonThenDefault
                                    onClick={() =>
                                        goHome()
                                    }
                                    text="Then That"
                                    disabled={!selectedAction}
                                />
                            )}
                        </div>

                        <div
                            className="w-full max-w-[350px] flex flex-col items-center justify-center"
                        >
                            {selectedAction && selectedReaction ? (
                                <TriggerButton
                                    text="Create Workflow"
                                    textColor="text-white"
                                    bgColor="bg-[#222222]"
                                    hoverColor="hover:bg-[#333333]"
                                    onClick={handleContinue}
                                />
                            ) : (
                                <></>
                            )}
                        </div>
                    </div>
                </div>
            )

            case Steps.SELECT_SERVICE_ACTION:
                return (
                    <div className="min-h-screen bg-gray-100">
                        <div className="flex items-center justify-between px-10 py-5">
                            <div>
                                <CancelButton
                                    text="Cancel"
                                    textColor="text-[#222222]"
                                    borderColor="border-[#222222]"
                                    bgColor="bg-[#F9FAFB]"
                                    onClick={() => goToStep(Steps.HOME)}
                                />
                            </div>

                            <h1 className="text-[3.2rem] font-extrabold text-center flex-1">
                                Choose a service
                            </h1>

                            <div className="w-[100px]"></div>
                        </div>

                        <div className="mt-10 px-[140px]">
                            <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
                                {services.map((service: any) => (
                                    <button
                                        key={service.id}
                                        style={{ backgroundColor: service.color }}
                                        className="h-[240px] flex flex-col items-center cursor-pointer"
                                        onClick={() => goHome()}
                                    >
                                        <img alt="icon" src={`${import.meta.env.VITE_HOST_URL}${service.iconPath}`} className="size-28 mt-9"/>
                                        <h2 className="text-white text-lg text-[20.5px] font-bold mt-6">{service.name}</h2>
                                    </button>
                                ))}
                            </div>
                        </div>
                    </div>
                );
        }
    }

    return (
        <>
            {renderStep()}
        </>
    );
};

export default WorkflowBuilder

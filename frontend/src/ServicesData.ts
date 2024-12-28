export type ServiceDataType = {
    [serviceName: string]: {
        actions?: {
            [actionName: string]: {
                id: string;
                label: string;
                type: string;
                options?: string[];
            }[];
        };
        reactions?: {
            [reactionName: string]: {
                serviceConnection: string;
                triggerReaction: string;
            }[];
        };
    };
};

export const servicesData: ServiceDataType = {
    "Time & Date": {
        actions: {
            "Every day at": [
                {
                    id: "time",
                    label: "Time",
                    type: "select",
                    options: ["0 h", "1 h", "2 h", "3 h", "4 h", "5 h", "6 h", "7 h", "8 h", "9 h", "10 h", "11 h", "12 h", "13 h", "14 h", "15 h", "16 h", "17 h", "18 h", "19 h", "20 h", "21 h", "22 h", "23 h"],
                },
                {
                    id: "quarter",
                    label: "Quarter Hour",
                    type: "select",
                    options: ["0 minutes", "15 minutes", "30 minutes", "45 minutes"],
                },
            ],
        },
    },
    "Spotify": {
        reactions: {
            "Start playback": [
                {
                    serviceConnection: "yes",
                    triggerReaction: "no",
                },
            ],
        },
    },
};
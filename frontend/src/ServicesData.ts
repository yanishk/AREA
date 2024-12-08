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
                id?: string;
                label?: string;
                type?: string;
                options?: string[];
            }[];
        }
    };
};

export const servicesData: ServiceDataType = {
    "Time & Date": {
        actions: {
            "Every day at": [
                { id: "time", label: "Time", type: "select", options: ["12 AM", "01 AM", "02 AM", "03 AM", "04 AM", "05 AM", "06 AM", "07 AM", "08 AM", "09 AM", "10 AM", "11 AM", "12 PM", "01 PM", "02 PM", "03 PM", "04 PM", "05 PM", "06 PM", "07 PM", "08 PM", "09 PM", "10 PM", "11 PM",]},
                { id: "quarter", label: "Quarter Hour", type: "select", options: ["00 minutes", "15 mintues", "30 mintues", "45 minutes"] },
            ],
        },
    },
    "Spotify": {
        reactions: {
            "Start playback": [
                {serviceConnection: "yes", triggerReaction: "no"},
            ]
        }
    }
};

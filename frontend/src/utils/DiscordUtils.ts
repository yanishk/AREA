import axios from "axios";

export interface DiscordItem {
    id: string;
    name: string;
}
export interface DiscordServer {
    id: string;
    name: string;
}

let Servers: DiscordServer[] = [];
let selectedServerId: string | null = null;
let selectedServerName: string | null = null;

const setDefaultServers = (servers: DiscordServer[]) => {
    if (servers.length === 0) {
        return;
    }
    selectedServerId = servers[0].id;
    selectedServerName = servers[0].name;
};

let isFetchingServers = false;

export const getDiscordServers = async (): Promise<DiscordServer[]> => {
    if (isFetchingServers) {
        console.warn("Discord servers already fetched");
        return Servers;
    }

    isFetchingServers = true;

    try {
        const response = await axios.get(`${import.meta.env.VITE_API_URL}discord/user/servers`, {
            withCredentials: true,
        });

        if (Array.isArray(response.data?.servers)) {
            Servers = response.data.servers.map((server: { id: string; name: string }) => ({
                id: server.id,
                name: server.name,
            }));

            setDefaultServers(Servers);

            return Servers;
        }

        console.warn("No servers found");
        return [];
    } catch (error) {
        console.error("Error fetching Discord servers : ", error);
        return [];
    }
};

export const getServerIdByName = async (workspaceName: string): Promise<string | null> => {
    const servers = await getDiscordServers();
    const server = servers.find((ws) => ws.name === workspaceName);
    return server ? server.id : null;
};

export const getSelectedServerId = (): string | null => {
    return selectedServerId;
};

export const getSelectedServerName = (): string | null => {
    return selectedServerName;
};

export const setSelectedServer = (serverId: string, serverName: string) => {
    selectedServerId = serverId;
    selectedServerName = serverName;
};

export const getDiscordChannels = async (serverId: string): Promise<DiscordItem[]> => {
    try {
        const response = await axios.get(
            `${import.meta.env.VITE_API_URL}discord/server/channels?id=${serverId}`,
            {
                withCredentials: true,
            }
        );

        if (Array.isArray(response.data?.channels)) {
                return response.data.channels.map((channel: { id: string; name: string }) => ({
                    id: channel.id,
                    name: channel.name,
                }));
        }

        console.warn("No channels found : ", serverId);
        return [];
    } catch (error) {
        console.error("Error fetching Discord channels : ", error);
        return [];
    }
};

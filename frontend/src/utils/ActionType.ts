import { Parameters } from "./ActionReactionParameters";

export type ActionType = {
    name: string;
    id: string;
    nbparam: number;
    parameters: Parameters[];
};
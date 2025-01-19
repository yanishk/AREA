import { Parameters } from "./ActionReactionParameters";

export type ReactionType = {
    name: string;
    id: string;
    nbparam: number;
    parameters: Parameters[];
};
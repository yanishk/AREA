import { Parameters } from "./ActionReactionParameters";

export const getFieldValues = (param: Parameters): string[] | undefined => {
    if (param.values[0] === null) {
        return []
    }
    if (param.isexhaustive === true && param.values.length > 1) {
        return param.values;
    }

    if (param.isexhaustive === false && param.values.length === 1) {
        const value = param.values[0];
        const [start, end] = value.split("-").map(Number);
        const newArray = Array.from({ length: end - start + 1 }, (_, i) => (start + i).toString());
        return newArray;
    }

    return undefined;
};

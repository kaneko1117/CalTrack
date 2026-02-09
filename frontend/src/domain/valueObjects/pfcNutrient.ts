import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type PfcNutrientValue = "protein" | "fat" | "carbs";

export type PfcNutrientErrorCode = "PFC_NUTRIENT_INVALID";

export type PfcNutrientError = DomainError<PfcNutrientErrorCode>;

export type PfcNutrient = Readonly<{
  value: PfcNutrientValue;
  getLabel: () => string;
  getShortLabel: () => string;
  equals: (other: PfcNutrient) => boolean;
}>;

/** 選択肢の型 */
export type PfcNutrientOption = {
  value: PfcNutrientValue;
  label: string;
  shortLabel: string;
};

/** PFC栄養素の選択肢（UI用） */
export const PFC_NUTRIENT_OPTIONS: PfcNutrientOption[] = [
  { value: "protein", label: "タンパク質", shortLabel: "P" },
  { value: "fat", label: "脂質", shortLabel: "F" },
  { value: "carbs", label: "炭水化物", shortLabel: "C" },
];

const VALID_PFC_NUTRIENTS: PfcNutrientValue[] = PFC_NUTRIENT_OPTIONS.map(
  (o) => o.value
);

const ERROR_MESSAGE_PFC_NUTRIENT_INVALID = "PFC栄養素を選択してください";

export const newPfcNutrient = (
  value: string
): Result<PfcNutrient, PfcNutrientError> => {
  if (!VALID_PFC_NUTRIENTS.includes(value as PfcNutrientValue)) {
    return err(
      domainError("PFC_NUTRIENT_INVALID", ERROR_MESSAGE_PFC_NUTRIENT_INVALID)
    );
  }

  const option = PFC_NUTRIENT_OPTIONS.find((o) => o.value === value)!;

  const pfcNutrient: PfcNutrient = Object.freeze({
    value: value as PfcNutrientValue,
    getLabel: () => option.label,
    getShortLabel: () => option.shortLabel,
    equals: (other: PfcNutrient) => value === other.value,
  });
  return ok(pfcNutrient);
};

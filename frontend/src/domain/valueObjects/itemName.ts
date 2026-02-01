import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type ItemNameErrorCode = "ITEM_NAME_REQUIRED";

export type ItemNameError = DomainError<ItemNameErrorCode>;

export type ItemName = Readonly<{
  value: string;
  equals: (other: ItemName) => boolean;
}>;

const ERROR_MESSAGE_ITEM_NAME_REQUIRED = "食品名を入力してください";

export const newItemName = (value: string): Result<ItemName, ItemNameError> => {
  if (!value || value.trim() === "") {
    return err(domainError("ITEM_NAME_REQUIRED", ERROR_MESSAGE_ITEM_NAME_REQUIRED));
  }

  const itemName: ItemName = Object.freeze({
    value,
    equals: (other: ItemName) => value === other.value,
  });
  return ok(itemName);
};

export type DomainError<T extends string> = Readonly<{
  code: T;
  message: string;
}>;

export const domainError = <T extends string>(
  code: T,
  message: string
): DomainError<T> => Object.freeze({ code, message });

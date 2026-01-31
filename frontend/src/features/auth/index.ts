/**
 * Auth Feature
 * Re-exports all auth-related types, API functions, and hooks
 */

// Types
export type {
  Gender,
  ActivityLevel,
  ErrorCode,
  RegisterUserRequest,
  RegisterUserResponse,
  ApiErrorResponse,
} from "./types";

// API
export { registerUser, ApiError } from "./api";

// Hooks
export { useRegisterUser } from "./hooks";
export type { UseRegisterUserReturn } from "./hooks";

// Components
export { RegisterForm, RegisterPage } from "./components";
export type { RegisterFormProps, RegisterPageProps } from "./components";

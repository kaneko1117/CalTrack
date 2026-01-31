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
  LoginRequest,
  LoginResponse,
  LogoutResponse,
  ApiErrorResponse,
} from "./types";

export {
  ERROR_CODE_INTERNAL_ERROR,
  ERROR_CODE_INVALID_CREDENTIALS,
  ERROR_MESSAGE_UNEXPECTED,
  ERROR_MESSAGE_INVALID_CREDENTIALS,
} from "./types";

// API
export { registerUser, login, logout, ApiError } from "./api";

// Hooks
export { useRegisterUser, useLogin, useLogout } from "./hooks";
export type { UseRegisterUserReturn, UseLoginReturn, UseLogoutReturn } from "./hooks";

// Components
export { RegisterForm, RegisterPage } from "./components";
export type { RegisterFormProps, RegisterPageProps } from "./components";

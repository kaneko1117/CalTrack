import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useRegisterUser } from "./index";
import * as authApi from "../api";
import type { RegisterUserRequest } from "../types";
import { ERROR_CODE_INTERNAL_ERROR } from "../types";

// Mock the api module
vi.mock("../api", async () => {
  const actual = await vi.importActual<typeof authApi>("../api");
  return {
    ...actual,
    registerUser: vi.fn(),
  };
});

const mockRegisterUser = vi.mocked(authApi.registerUser);

describe("useRegisterUser", () => {
  const validRequest: RegisterUserRequest = {
    email: "test@example.com",
    password: "password123",
    nickname: "TestUser",
    weight: 70,
    height: 175,
    birthDate: "1990-01-01",
    gender: "male",
    activityLevel: "moderate",
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.resetAllMocks();
  });

  describe("initial state", () => {
    it("should have correct initial state", () => {
      const { result } = renderHook(() => useRegisterUser());

      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(result.current.isSuccess).toBe(false);
      expect(typeof result.current.register).toBe("function");
      expect(typeof result.current.reset).toBe("function");
    });
  });

  describe("successful registration", () => {
    it("should set isSuccess to true on successful registration", async () => {
      mockRegisterUser.mockResolvedValueOnce({ userId: "user-123" });

      const { result } = renderHook(() => useRegisterUser());

      await act(async () => {
        await result.current.register(validRequest);
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.error).toBeNull();
      });

      expect(mockRegisterUser).toHaveBeenCalledWith(validRequest);
      expect(mockRegisterUser).toHaveBeenCalledTimes(1);
    });

    it("should set isLoading to true during registration", async () => {
      let resolvePromise: (value: { userId: string }) => void;
      const pendingPromise = new Promise<{ userId: string }>((resolve) => {
        resolvePromise = resolve;
      });
      mockRegisterUser.mockReturnValueOnce(pendingPromise);

      const { result } = renderHook(() => useRegisterUser());

      act(() => {
        result.current.register(validRequest);
      });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(true);
      });

      await act(async () => {
        resolvePromise!({ userId: "user-123" });
      });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });
    });
  });

  describe("validation error", () => {
    it("should set error on validation error", async () => {
      const validationError = new authApi.ApiError(
        "VALIDATION_ERROR",
        "Validation failed",
        400,
        ["email: invalid format", "password: too short"]
      );
      mockRegisterUser.mockRejectedValueOnce(validationError);

      const { result } = renderHook(() => useRegisterUser());

      await act(async () => {
        await result.current.register(validRequest);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
        expect(result.current.error?.code).toBe("VALIDATION_ERROR");
        expect(result.current.error?.details).toEqual([
          "email: invalid format",
          "password: too short",
        ]);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.isSuccess).toBe(false);
      });
    });
  });

  describe("email already exists error", () => {
    it("should set error when email already exists", async () => {
      const emailExistsError = new authApi.ApiError(
        "EMAIL_ALREADY_EXISTS",
        "Email already registered",
        409
      );
      mockRegisterUser.mockRejectedValueOnce(emailExistsError);

      const { result } = renderHook(() => useRegisterUser());

      await act(async () => {
        await result.current.register(validRequest);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
        expect(result.current.error?.code).toBe("EMAIL_ALREADY_EXISTS");
        expect(result.current.isLoading).toBe(false);
        expect(result.current.isSuccess).toBe(false);
      });
    });
  });

  describe("server error", () => {
    it("should set error on server error", async () => {
      const serverError = new authApi.ApiError(
        ERROR_CODE_INTERNAL_ERROR,
        "Internal server error",
        500
      );
      mockRegisterUser.mockRejectedValueOnce(serverError);

      const { result } = renderHook(() => useRegisterUser());

      await act(async () => {
        await result.current.register(validRequest);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
        expect(result.current.error?.code).toBe(ERROR_CODE_INTERNAL_ERROR);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.isSuccess).toBe(false);
      });
    });

    it("should handle unexpected error", async () => {
      mockRegisterUser.mockRejectedValueOnce(new Error("Network error"));

      const { result } = renderHook(() => useRegisterUser());

      await act(async () => {
        await result.current.register(validRequest);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
        expect(result.current.error?.code).toBe(ERROR_CODE_INTERNAL_ERROR);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.isSuccess).toBe(false);
      });
    });
  });

  describe("reset", () => {
    it("should reset all state", async () => {
      const validationError = new authApi.ApiError(
        "VALIDATION_ERROR",
        "Validation failed",
        400
      );
      mockRegisterUser.mockRejectedValueOnce(validationError);

      const { result } = renderHook(() => useRegisterUser());

      await act(async () => {
        await result.current.register(validRequest);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
      });

      act(() => {
        result.current.reset();
      });

      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(result.current.isSuccess).toBe(false);
    });
  });
});

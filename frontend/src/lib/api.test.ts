import { describe, it, expect, vi, beforeEach } from "vitest";
import { AxiosError, AxiosHeaders } from "axios";

/** routerのモック（vi.hoistedでホイスト対応） */
const { mockNavigate, mockPathname } = vi.hoisted(() => ({
  mockNavigate: vi.fn(),
  mockPathname: { value: "/dashboard" },
}));

vi.mock("@/routes", () => ({
  router: {
    state: {
      location: {
        get pathname() {
          return mockPathname.value;
        },
      },
    },
    navigate: mockNavigate,
  },
}));

/** インターセプターのエラーハンドラを保持 */
type ErrorHandler = (error: AxiosError) => Promise<never>;
const { capturedHandler } = vi.hoisted(() => ({
  capturedHandler: { current: null as ErrorHandler | null },
}));

vi.mock("axios", async () => {
  const actual = await vi.importActual("axios");
  return {
    ...actual,
    default: {
      create: vi.fn(() => ({
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
        interceptors: {
          response: {
            use: vi.fn((_onFulfilled: unknown, onRejected: ErrorHandler) => {
              capturedHandler.current = onRejected;
            }),
          },
        },
      })),
    },
  };
});

describe("apiClient", () => {
  describe("設定", () => {
    it("apiClientが正しく作成されている", async () => {
      const { apiClient } = await import("./api");
      expect(apiClient).toBeDefined();
    });
  });
});

describe("401エラー時のリダイレクト", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockPathname.value = "/dashboard";
    capturedHandler.current = null;
  });

  it("401エラー時にログイン画面へリダイレクトする", async () => {
    vi.resetModules();
    await import("./api");

    const axiosError = new AxiosError("Unauthorized");
    axiosError.response = {
      data: { code: "UNAUTHORIZED", message: "Unauthorized" },
      status: 401,
      statusText: "Unauthorized",
      headers: {},
      config: { headers: new AxiosHeaders() },
    };

    expect(capturedHandler.current).not.toBeNull();
    await expect(capturedHandler.current!(axiosError)).rejects.toEqual(
      axiosError
    );
    expect(mockNavigate).toHaveBeenCalledWith("/");
  });

  it("401以外のエラー時はリダイレクトしない", async () => {
    vi.resetModules();
    await import("./api");

    const axiosError = new AxiosError("Bad Request");
    axiosError.response = {
      data: { code: "VALIDATION_ERROR", message: "Invalid input" },
      status: 400,
      statusText: "Bad Request",
      headers: {},
      config: { headers: new AxiosHeaders() },
    };

    expect(capturedHandler.current).not.toBeNull();
    await expect(capturedHandler.current!(axiosError)).rejects.toEqual(
      axiosError
    );
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  it("既にログイン画面にいる場合はリダイレクトしない", async () => {
    mockPathname.value = "/";
    vi.resetModules();
    await import("./api");

    const axiosError = new AxiosError("Unauthorized");
    axiosError.response = {
      data: { code: "UNAUTHORIZED", message: "Unauthorized" },
      status: 401,
      statusText: "Unauthorized",
      headers: {},
      config: { headers: new AxiosHeaders() },
    };

    expect(capturedHandler.current).not.toBeNull();
    await expect(capturedHandler.current!(axiosError)).rejects.toEqual(
      axiosError
    );
    expect(mockNavigate).not.toHaveBeenCalled();
  });
});

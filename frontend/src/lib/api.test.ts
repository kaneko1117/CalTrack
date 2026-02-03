import { describe, it, expect, vi, beforeEach } from "vitest";
import { request, get, post, put, del, apiClient } from "./api";
import { AxiosError, AxiosHeaders, AxiosResponse } from "axios";

/** routerのモック（vi.hoistedでホイスト対応） */
const { mockNavigate } = vi.hoisted(() => ({
  mockNavigate: vi.fn(),
}));

vi.mock("@/routes", () => ({
  router: {
    state: {
      location: {
        pathname: "/dashboard",
      },
    },
    navigate: mockNavigate,
  },
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
            use: vi.fn(),
          },
        },
      })),
    },
  };
});

/**
 * AxiosResponseのモックを生成
 */
function createMockResponse<T>(data: T): AxiosResponse<T> {
  return {
    data,
    status: 200,
    statusText: "OK",
    headers: {},
    config: { headers: new AxiosHeaders() },
  };
}

describe("request", () => {
  describe("正常系", () => {
    it("レスポンスのdataを返す", async () => {
      const mockData = { id: 1, name: "test" };
      const mockPromise = Promise.resolve(createMockResponse(mockData));

      const result = await request(mockPromise);

      expect(result).toEqual(mockData);
    });
  });

  describe("異常系", () => {
    it("AxiosErrorをApiErrorResponse形式に変換する", async () => {
      const apiError = { code: "VALIDATION_ERROR", message: "Invalid input" };
      const axiosError = new AxiosError("Request failed");
      axiosError.response = {
        data: apiError,
        status: 400,
        statusText: "Bad Request",
        headers: {},
        config: { headers: new AxiosHeaders() },
      };
      const mockPromise = Promise.reject(axiosError);

      await expect(request(mockPromise)).rejects.toEqual(apiError);
    });

    it("不明なエラーはINTERNAL_ERRORを返す", async () => {
      const mockPromise = Promise.reject(new Error("Unknown error"));

      await expect(request(mockPromise)).rejects.toEqual({
        code: "INTERNAL_ERROR",
        message: "予期しないエラーが発生しました",
      });
    });
  });
});

describe("get", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("apiClient.getを呼び出してdataを返す", async () => {
    const mockData = { id: 1 };
    vi.mocked(apiClient.get).mockResolvedValue({ data: mockData });

    const result = await get<typeof mockData>("/test");

    expect(apiClient.get).toHaveBeenCalledWith("/test");
    expect(result).toEqual(mockData);
  });
});

describe("post", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("apiClient.postを呼び出してdataを返す", async () => {
    const mockData = { id: 1 };
    const requestData = { name: "test" };
    vi.mocked(apiClient.post).mockResolvedValue({ data: mockData });

    const result = await post<typeof mockData>("/test", requestData);

    expect(apiClient.post).toHaveBeenCalledWith("/test", requestData);
    expect(result).toEqual(mockData);
  });
});

describe("put", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("apiClient.putを呼び出してdataを返す", async () => {
    const mockData = { id: 1 };
    const requestData = { name: "updated" };
    vi.mocked(apiClient.put).mockResolvedValue({ data: mockData });

    const result = await put<typeof mockData>("/test", requestData);

    expect(apiClient.put).toHaveBeenCalledWith("/test", requestData);
    expect(result).toEqual(mockData);
  });
});

describe("del", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("apiClient.deleteを呼び出してdataを返す", async () => {
    const mockData = { success: true };
    vi.mocked(apiClient.delete).mockResolvedValue({ data: mockData });

    const result = await del<typeof mockData>("/test");

    expect(apiClient.delete).toHaveBeenCalledWith("/test");
    expect(result).toEqual(mockData);
  });
});

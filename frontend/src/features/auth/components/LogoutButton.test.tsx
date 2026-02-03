import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { LogoutButton } from "./LogoutButton";
import * as swr from "@/lib/swr";

vi.mock("@/lib/swr", () => ({
  fetcher: vi.fn(),
  mutate: vi.fn(),
}));

describe("LogoutButton", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("ログアウトボタンが表示される", () => {
    render(<LogoutButton />);
    expect(screen.getByRole("button", { name: "ログアウト" })).toBeInTheDocument();
  });

  it("クリックするとログアウトAPIが呼ばれる", async () => {
    const mockMutate = vi.mocked(swr.mutate);
    mockMutate.mockResolvedValue({ message: "ログアウトしました" });

    const user = userEvent.setup();
    render(<LogoutButton />);
    await user.click(screen.getByRole("button", { name: "ログアウト" }));

    await waitFor(() => {
      expect(mockMutate).toHaveBeenCalledWith(
        "/api/v1/auth/logout",
        { arg: { method: "POST", data: undefined } }
      );
    });
  });

  it("成功時にonSuccessが呼ばれる", async () => {
    const mockMutate = vi.mocked(swr.mutate);
    mockMutate.mockResolvedValue({ message: "ログアウトしました" });
    const onSuccess = vi.fn();

    const user = userEvent.setup();
    render(<LogoutButton onSuccess={onSuccess} />);
    await user.click(screen.getByRole("button", { name: "ログアウト" }));

    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalled();
    });
  });
});

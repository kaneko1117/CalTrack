import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { LogoutButton } from "./LogoutButton";
import * as api from "@/lib/api";

vi.mock("@/lib/api", () => ({
  post: vi.fn(),
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
    const mockPost = vi.mocked(api.post);
    mockPost.mockResolvedValue({ message: "ログアウトしました" });

    const user = userEvent.setup();
    render(<LogoutButton />);
    await user.click(screen.getByRole("button", { name: "ログアウト" }));

    await waitFor(() => {
      expect(mockPost).toHaveBeenCalledWith("/api/v1/auth/logout");
    });
  });

  it("成功時にonSuccessが呼ばれる", async () => {
    const mockPost = vi.mocked(api.post);
    mockPost.mockResolvedValue({ message: "ログアウトしました" });
    const onSuccess = vi.fn();

    const user = userEvent.setup();
    render(<LogoutButton onSuccess={onSuccess} />);
    await user.click(screen.getByRole("button", { name: "ログアウト" }));

    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalled();
    });
  });
});

/**
 * API Client Configuration
 * Shared axios instance for all API calls
 */

import axios from "axios";

export const apiClient = axios.create({
  baseURL: "http://localhost:8080",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

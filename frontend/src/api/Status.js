import { request } from "./Client";

export function getStatus() {
  return request("/status");
}
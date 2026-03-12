import "axios";

declare module "axios" {
  export interface AxiosRequestConfig {
    authMode?: "required" | "optional" | "none";
  }
}

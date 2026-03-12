const TOKEN_KEY = "neohome.access_token";
const LOGOUT_EVENT = "neohome.logout";

export const authStore = {
  getToken: () => localStorage.getItem(TOKEN_KEY),
  setToken: (token: string) => localStorage.setItem(TOKEN_KEY, token),
  clearToken: () => localStorage.removeItem(TOKEN_KEY),
};

export const authEvents = {
  emitLogout: () => window.dispatchEvent(new Event(LOGOUT_EVENT)),
  onLogout: (callback: () => void) => {
    window.addEventListener(LOGOUT_EVENT, callback);
    return () => window.removeEventListener(LOGOUT_EVENT, callback);
  },
};

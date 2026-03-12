import { AppRouter } from "./router";
import { AuthProvider, QueryProvider } from "./providers";

export const App = () => {
  return (
    <QueryProvider>
      <AuthProvider>
        <AppRouter />
      </AuthProvider>
    </QueryProvider>
  );
};

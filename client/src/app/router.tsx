import { Navigate, Outlet, Route, Routes } from "react-router-dom";
import { useAuth } from "./providers";
import { RootLayout } from "./layouts/RootLayout";
import { Spinner } from "@/shared/ui/Spinner/Spinner";
import { AuthPage, DashboardPage, DevicePage, ProfilePage, SimulatorPage } from "@/pages";

const ProtectedRoute = () => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="page-center">
        <Spinner label="Loading profile..." />
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/auth" replace />;
  }

  return <Outlet />;
};

const HomeRoute = () => {
  const { isAuthenticated } = useAuth();
  return <Navigate to={isAuthenticated ? "/dashboard" : "/auth"} replace />;
};

const NotFoundPage = () => {
  return (
    <div className="container page">
      <h2>Page not found</h2>
    </div>
  );
};

export const AppRouter = () => {
  return (
    <Routes>
      <Route element={<RootLayout />}>
        <Route index element={<HomeRoute />} />
        <Route path="/auth" element={<AuthPage />} />
        <Route path="/simulator" element={<SimulatorPage />} />

        <Route element={<ProtectedRoute />}>
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/devices/:deviceId" element={<DevicePage />} />
          <Route path="/profile" element={<ProfilePage />} />
        </Route>

        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  );
};

import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../providers";
import { env } from "@/shared/config";
import { Button } from "@/shared/ui/Button/Button";
import styles from "./RootLayout.module.scss";

export const RootLayout = () => {
  const navigate = useNavigate();
  const { isAuthenticated, logout } = useAuth();

  const handleLogout = () => {
    logout();
    navigate("/auth");
  };

  return (
    <div className={styles.shell}>
      <header className={styles.header}>
        <div className="container">
          <div className={styles.headerInner}>
            <h1 className={styles.logo}>{env.appName}</h1>
            <nav className={styles.nav}>
              {isAuthenticated ? (
                <>
                  <NavLink to="/dashboard" className={styles.link}>
                    Dashboard
                  </NavLink>
                  <NavLink to="/profile" className={styles.link}>
                    Profile
                  </NavLink>
                  <NavLink to="/simulator" className={styles.link}>
                    Simulator
                  </NavLink>
                  <Button variant="ghost" onClick={handleLogout}>
                    Logout
                  </Button>
                </>
              ) : (
                <>
                  <NavLink to="/auth" className={styles.link}>
                    Auth
                  </NavLink>
                  <NavLink to="/simulator" className={styles.link}>
                    Simulator
                  </NavLink>
                </>
              )}
            </nav>
          </div>
        </div>
      </header>

      <main className={styles.main}>
        <Outlet />
      </main>
    </div>
  );
};

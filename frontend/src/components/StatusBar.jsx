import { useEffect, useState } from "react";
import { getStatus } from "../api/Status";

const REFRESH_INTERVAL_MS = 10000;

function StatusBar() {
  const [status, setStatus] = useState(null);
  const [offline, setOffline] = useState(false);

  useEffect(() => {
    let ignore = false;

    function fetchStatus() {
      getStatus()
        .then((data) => {
          if (!ignore) {
            setStatus(data);
            setOffline(false);
          }
        })
        .catch(() => {
          if (!ignore) {
            setOffline(true);
          }
        });
    }

    fetchStatus();
    const intervalId = setInterval(fetchStatus, REFRESH_INTERVAL_MS);

    return () => {
      ignore = true;
      clearInterval(intervalId);
    };
  }, []);

  let message = "Conectando...";

  if (offline) {
    message = "API indisponível";
  } else if (status) {
    message = `API online · Uptime: ${status.uptime} · Consultas: ${status.request_count}`;
  }

  return (
    <footer className="status-bar">
      <span
        className={`status-bar__dot ${
          offline ? "status-bar__dot--offline" : "status-bar__dot--online"
        }`}
      />
      {message}
    </footer>
  );
}

export default StatusBar;
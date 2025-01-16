import React, {useEffect} from 'react';
import getPathOnVersion from "../services/versions.jsx";

function ProcessList({ processes, setProcesses, ws, serverInfo, url, buttonEnable }) {

    const handleStop = (pid) => {
        if (ws && pid) {
            ws.send(JSON.stringify({ action: 'stop', pid: String(pid) }));
        }
    };

    useEffect(() => {
        const stoppedProcesses = processes.filter(process => process.status === "stopped");

        stoppedProcesses.forEach(process => {
            const timer = setTimeout(() => {
                setProcesses((prevProcesses) =>
                    prevProcesses.filter(p => p.pid !== process.pid)
                );
            }, 5000); // 5 секунд

            return () => clearTimeout(timer); // Очистка таймера при размонтировании
        });
    }, [processes, setProcesses]);

    return (
        <div className="process-list">
            <h2>Процессы</h2>
            <div>
                {processes.map((process) => (
                    <div key={process.pid} className="process-item">
                        <div className={`process-status-${process.status === "running" ? "running" : "stopped"}`}>{process.status === "running" ? "Запущен": "Остановлен"}</div>
                        <div className="process-info">
                            <div className="process-url-param">{process.urlParam}</div>
                            <div>{process.version}</div>
                            <div>{process.edition === 'default' ? 'Office' : 'Chain'}</div>
                        </div>
                        <button onClick={() => handleStop(process.pid)} className="delete-button">X</button>
                    </div>
                ))}
            </div>
        </div>
    );
    }

    export default ProcessList;